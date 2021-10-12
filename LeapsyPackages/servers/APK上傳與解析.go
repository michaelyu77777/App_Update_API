package servers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"leapsy.com/packages/configurations"
	"leapsy.com/records"
)

//func UploadSingleIndex(ctx *gin.Context) {
func UploadSingleIndex(apiServer *APIServer, ginContextPointer *gin.Context) {

	// 收檔案、表頭
	file, header, err := ginContextPointer.Request.FormFile("file")
	if err != nil {
		ginContextPointer.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	// 暫存檔案
	tempFileName := header.Filename // 取出檔名
	saveTempPath := configurations.GetConfigValueOrPanic(`local`, `pathTemp`)
	saveFileToPath(file, saveTempPath, tempFileName)

	fmt.Println("取出檔名＝" + tempFileName)

	// 解析暫存APK
	packageName, labelName, versionCode, versionName := getApkDetailsInApkTempDirectory(tempFileName)
	//_, labelName, _, versionName := getApkDetailsInApkTempDirectory(tempFileName)

	// 查找是否已建檔
	result := mongoDB.FindAppsInfoByLabelName(labelName)

	// 回傳給Client的訊息
	message := ""

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

		fmt.Println("已完成初次建檔")

		if err == nil {
			message += "此 Apk Label 識別為「初次上傳」，已為您完成建檔。"
		} else {
			message += fmt.Sprintf("此 Apk Label 識別為「初次上傳」，建檔時發生錯誤，Error： %s", err.Error())
		}
	}

	// 看apk儲存資料夾下是否存在名稱為「labelName」的資料夾
	apkPath := configurations.GetConfigValueOrPanic(`local`, `path`) + labelName + "/"
	// apkPath := "apk/" + labelName + "/"

	isFileExist, err := isExists(apkPath)
	if err != nil {
		fmt.Println("判斷資料夾是否存在時發生錯誤：", err)
		log.Fatal(err) // Log待補
	}

	// 「labelName」的資料夾若不存，則建立「labelName」資料夾
	if !isFileExist {
		fmt.Println("LabelName資料夾不存在") // Log待補

		//創建目錄
		err = os.Mkdir(configurations.GetConfigValueOrPanic(`local`, `path`)+labelName, os.ModePerm)

		if err != nil {
			fmt.Println("創建資料夾時發生錯誤：", err)
			log.Fatal(err) // Log待補

			message += "創建資料夾時發生錯誤"

			ginContextPointer.JSON(http.StatusOK, gin.H{
				"issuccess": false,
				"message":   message,
			})

		} else {
			fmt.Println("已創建新資料夾") // Log待補
		}
	}

	// 再收一次檔案,歸檔到正式LableName的資料夾
	file, header, err = ginContextPointer.Request.FormFile("file")
	if err != nil {
		ginContextPointer.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	// 存檔：命名為 label + V_ + versionName
	apkName := labelName + "_v" + versionName + ".apk"
	savePath := configurations.GetConfigValueOrPanic(`local`, `path`) + labelName + "/"
	saveFileToPath(file, savePath, apkName)

	// 刪除tempAPK檔
	err = os.Remove(saveTempPath + tempFileName)
	if err != nil {
		log.Fatal(err) // Log待補
		fmt.Println(err)
		message += fmt.Sprintf("刪除檔案時發生錯誤，Error : %s ", err.Error())
	} else {
		fmt.Println("已刪除暫存檔案")
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
	message += "您已完成上傳檔案"

	if 1 > len(results) {
		//查無結果
		ginContextPointer.JSON(http.StatusOK, gin.H{
			"issuccess": true,
			"message":   message,
			"appsinfo":  results,
		})
	} else {
		//有查到結果
		ginContextPointer.JSON(http.StatusOK, gin.H{
			"issuccess": true,
			"message":   message,
			"appsinfo":  results[0],
		})
	}

}

// 確認檔案或資料夾是否存在
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
func saveFileToPath(file multipart.File, path string, fileName string) {

	//測試碼 取得大小
	switch t := file.(type) {
	case *os.File:
		fi, _ := t.Stat()
		fmt.Println(fi.Size())
		fmt.Println("測試:儲存的檔案大小")
	default:
		fmt.Println("測試:預設")
		// Do Something
	}

	out, err := os.Create(path + fileName) // 建立空檔

	if err != nil {
		log.Fatal(err) // Log待補
		fmt.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file) // 將file數據複製到空檔

	if err != nil {
		log.Fatal(err) // Log待補
		fmt.Println(err)
	}

}
