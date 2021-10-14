package servers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"leapsy.com/packages/configurations"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
)

// 處理APK上傳
/**
 * @param  apiServer *APIServer
 * @param  ginContextPointer *gin.Context  gin Context 指標
 */
func UploadSingleApkAPIHandler(apiServer *APIServer, ginContextPointer *gin.Context) {

	// For logings
	defaultArgs :=
		append(
			network.GetAliasAddressPair(
				fmt.Sprintf(`%s:%d`,
					apiServer.GetConfigValueOrPanic(`host`),
					apiServer.GetConfigPositiveIntValueOrPanic(`port`),
				),
			),
			ginContextPointer.ClientIP(),
			ginContextPointer.FullPath(),
		)

	logings.SendLog(
		[]string{`%s %s 接受 %s 請求 %s `},
		defaultArgs,
		nil,
		0,
	)

	// 收檔案、表頭
	file, header, err := ginContextPointer.Request.FormFile("file")
	message := ""

	if nil != err {
		s := fmt.Sprintf("[錯誤]apk收檔時發生錯誤，錯誤訊息如下，Error：%s。", err.Error())
		// fmt.Println(s)

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		message += s
		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})
		return
	}

	// 暫存檔案
	tempFileName := header.Filename // 取出檔名
	saveTempPath := configurations.GetConfigValueOrPanic(`local`, `pathTemp`)
	err, msg := saveFileToPath(file, saveTempPath, tempFileName)
	if nil != err {
		s := fmt.Sprintf("[錯誤]儲存暫存apk檔時發生錯誤，錯誤訊息如下，Error：%s，Msg:%s。", err.Error(), msg)

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		message += s
		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})
		return
	}

	fmt.Println("取出檔名＝" + tempFileName)

	// 取副檔名
	fileExtension := filepath.Ext(tempFileName)

	// 檢核副檔名是否為APK（case insensitive）
	if !strings.EqualFold(fileExtension, ".apk") {
		s := "[錯誤]非apk檔,判斷副檔名為" + fileExtension + "。"
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})
		return
	}

	// 解析暫存APK
	err, msg, packageName, labelName, versionCode, versionName := getApkDetailsInApkTempDirectory(tempFileName)

	if nil != err {
		s := fmt.Sprintf("[錯誤]解析APK時發生錯誤，錯誤訊息如下，Error：%s，Msg：%s。", err.Error(), msg)

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		message += s
		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})

		return
	}

	// 查找是否已建檔
	result := mongoDB.FindAppsInfoByLabelName(labelName)

	// 若沒建檔，則先建檔
	if 1 > len(result) {

		fmt.Println("未建檔")

		// 建立一筆新的appsInfo
		appsInfoCommonStruct := records.AppsInfoCommonStruct{
			PackageName:     packageName,
			LabelName:       labelName,
			LastVersionCode: versionCode,
			LastVersionName: versionName,
		}

		document := records.AppsInfo{
			AppsInfoCommonStruct: appsInfoCommonStruct,
		}

		err := mongoDB.InsertOneAppsInfo(document)

		if err == nil {
			s := "[完成初次建檔]解析apk之Label識別為初次上傳。"
			message += s

			// log
			logings.SendLog(
				[]string{`%s %s 接受 %s 請求 %s %s`},
				append(defaultArgs, s),
				err,
				logrus.InfoLevel,
			)

		} else {
			s := fmt.Sprintf("[錯誤]資料庫初次建檔時發生錯誤，錯誤訊息如下，Error： %s。", err.Error())
			message += s

			// log
			logings.SendLog(
				[]string{`%s %s 接受 %s 請求 %s %s`},
				append(defaultArgs, s),
				err,
				logrus.WarnLevel,
			)

			ginContextPointer.JSON(http.StatusInternalServerError, gin.H{
				"issuccess": false,
				"message":   message,
			})
			return
		}
	}

	// 確認apk儲存資料夾下是否存在名稱為「labelName」的資料夾
	apkPath := configurations.GetConfigValueOrPanic(`local`, `path`) + labelName + "/"
	// apkPath := "apk/" + labelName + "/"

	isFileExist, err := isExists(apkPath)
	if err != nil {
		s := fmt.Sprintf("[錯誤]判斷存檔Lable資料夾是否存在時發生錯誤，錯誤訊息如下，Error： %s。", err.Error())
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

	}

	// 「labelName」的資料夾若不存，則建立「labelName」資料夾
	if !isFileExist {

		//創建目錄
		err = os.Mkdir(configurations.GetConfigValueOrPanic(`local`, `path`)+labelName, os.ModePerm)

		if err != nil {
			s := fmt.Sprintf("[錯誤]創建名為Lable的資料夾時發生錯誤，錯誤訊息如下，Error： %s。", err.Error())
			message += s

			// log
			logings.SendLog(
				[]string{`%s %s 接受 %s 請求 %s %s`},
				append(defaultArgs, s),
				err,
				logrus.WarnLevel,
			)

			ginContextPointer.JSON(http.StatusInternalServerError, gin.H{
				"issuccess": false,
				"message":   message,
			})
			return

		} else {
			s := "[創建資料夾]已建立名為Label的資料夾。"
			message += s

			// log
			logings.SendLog(
				[]string{`%s %s 接受 %s 請求 %s %s`},
				append(defaultArgs, s),
				err,
				logrus.InfoLevel,
			)
		}
	}

	// 再收一次檔案,歸檔到正式LableName的資料夾
	file, header, err = ginContextPointer.Request.FormFile("file")
	if err != nil {
		s := fmt.Sprintf("[錯誤]正式儲存apk檔時發生錯誤，錯誤訊息如下，Error：%s", err.Error())
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})
		return
	}

	// 存檔：命名為 label + V_ + versionName
	apkName := labelName + "_v" + versionName + ".apk"
	savePath := configurations.GetConfigValueOrPanic(`local`, `path`) + labelName + "/"
	err, msg = saveFileToPath(file, savePath, apkName)

	if nil != err {
		s := fmt.Sprintf("[錯誤]儲存正式apk檔時發生錯誤，錯誤訊息如下，Error：%s，Msg:%s。", err.Error(), msg)
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		ginContextPointer.JSON(http.StatusBadRequest, gin.H{
			"issuccess": false,
			"message":   message,
		})
		return
	}

	// 刪除tempAPK檔
	err = os.Remove(saveTempPath + tempFileName)
	if err != nil {
		// Log待補
		s := fmt.Sprintf("[錯誤]刪除暫存Apk檔案時發生錯誤，錯誤訊息如下，Error：%s", err.Error())
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

	} else {
		s := "[已刪除暫存檔案]"
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.InfoLevel,
		)

	}

	// 將解析後的資訊，全部更新到對應app的appsinfo中
	results := mongoDB.FindOneAndUpdateAppsInfoSET(
		bson.M{
			"labelname": labelName,
		},
		bson.M{
			"packagename":       packageName,
			"lastversioncode":   versionCode,
			"lastversionname":   versionName,
			"apkfilename":       apkName,
			"lastapkupdatetime": time.Now(),
		})

	// Response Client

	if 1 > len(results) {
		//查無結果
		s := "[錯誤]資料庫查不到您上傳的APK建檔資料"
		message += s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.WarnLevel,
		)

		ginContextPointer.JSON(http.StatusInternalServerError, gin.H{
			"issuccess": false,
			"message":   message,
			"appsinfo":  results,
		})
	} else {
		//有查到結果
		s := "[您已完成APK檔案上傳,並於資料庫建檔或更新資料]"
		message = s

		// log
		logings.SendLog(
			[]string{`%s %s 接受 %s 請求 %s %s`},
			append(defaultArgs, s),
			err,
			logrus.InfoLevel,
		)

		ginContextPointer.JSON(http.StatusOK, gin.H{
			"issuccess": true,
			"message":   message,
			"appsinfo":  results[0],
		})
	}

}

// 確認檔案或資料夾是否存在
/**
 * @param  path string 指定路徑
 * return bool, error 結果,錯誤訊息
 */
func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 儲存檔案,用指定檔名,存到指定路徑
/**
 * @param file multipart.File 目標檔案
 * @param path string 指定路徑
 * @param fileName string 指定檔名
 * return err, msg 錯誤,相關訊息
 */
func saveFileToPath(file multipart.File, path string, fileName string) (err error, msg string) {

	// //測試碼 取得大小
	// switch t := file.(type) {
	// case *os.File:
	// 	fi, _ := t.Stat()
	// 	fmt.Println(fi.Size())
	// 	fmt.Println("測試:儲存的檔案大小")
	// default:
	// 	fmt.Println("測試:預設")
	// 	// Do Something
	// }

	out, err := os.Create(path + fileName) // 建立空檔
	defer out.Close()

	if err != nil {
		s := fmt.Sprintf("[錯誤]建立空檔，發生錯誤，錯誤：%s。", err.Error())
		msg += s
		return
	}

	_, err = io.Copy(out, file) // 將file數據複製到空檔

	if err != nil {
		// Log待補
		s := fmt.Sprintf("[錯誤]將file數據複製到空檔，發生錯誤，錯誤：%s。", err.Error())
		msg += s
		fmt.Println(s)
		return
	}

	return
}
