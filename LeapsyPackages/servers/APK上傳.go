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

	// 暫存APK
	tempFileName := header.Filename // 取出檔名
	saveTempPath := "apkTemp/"
	saveFileToPath(file, saveTempPath, tempFileName)

	fmt.Println("取出檔名＝" + tempFileName)

	// 解析暫存APK
	packageName, labelName, versionCode, versionName := getApkDetailsInApkTempDirectory(tempFileName)
	//_, labelName, _, versionName := getApkDetailsInApkTempDirectory(tempFileName)

	// 查找是否已建檔
	result := mongoDB.FindAppsInfoByLabelName(labelName)

	// 若已建檔，有找到結果，儲存正式APK檔
	if 1 > len(result) {
		//	沒建檔，刪除暫存檔
		fmt.Println("未建檔")

		// 刪除tempAPK檔
		err = os.Remove(saveTempPath + tempFileName)
		if err != nil {
			log.Fatal(err) // Log待補
			fmt.Println(err)
		} else {
			fmt.Println("已刪除暫存檔案")
		}

		// Response Client
		message := ""

		if err == nil {
			message = "您尚未建檔，無法上傳Apk檔，請先進行AppsInfo初次建檔"
		} else {
			message = fmt.Sprintf("Error : %s", err.Error())
		}

		ginContextPointer.JSON(http.StatusOK, gin.H{
			"isSuccess": false,
			"results":   message,
			"appsInfo": records.AppsInfoCommonStruct{
				LabelName:       labelName,
				PackageName:     packageName,
				LastVersionCode: versionCode,
				LastVersionName: versionName,
			},
		})

	} else {
		fmt.Println("已建檔")

		// 看路徑資料夾是否存在
		apkPath := "apk/" + labelName + "/"
		isFileExist, err := isExists(apkPath)
		if err != nil {
			log.Fatal(err) // Log待補
			fmt.Println(err)
		}

		// 若不存在則建立資料夾
		if !isFileExist {

			err = os.Mkdir("apk/"+labelName, os.ModePerm) //創建目錄
			if err != nil {
				log.Fatal(err) // Log待補
				fmt.Println(err)
			} else {
				fmt.Println("資料夾不存在，已創建資料夾") // Log待補
			}
		}

		// 存檔：命名為 label + V_ + versionName
		apkName := labelName + "_v" + versionName + ".apk"
		savePath := "apk/" + labelName + "/"
		saveFileToPath(file, savePath, apkName)

		// 刪除tempAPK檔
		err = os.Remove(saveTempPath + tempFileName)
		if err != nil {
			log.Fatal(err) // Log待補
			fmt.Println(err)
		} else {
			fmt.Println("已刪除暫存檔案")
		}

		// 將解析後的資訊，全部更新到對應app的appsinfo中
		results := mongoDB.FindOneAndUpdateAppsInfoSET(
			bson.M{
				"labelName": labelName,
			},
			bson.M{
				"packageName":       packageName,
				"lastVersionCode":   versionCode,
				"lastVersionName":   versionName,
				"apkFileName":       apkName,
				"lastApkUpdateTime": time.Now(),
			})

		// Response Client
		isSuccess := true
		message := ""

		if err == nil {
			message = "您已完成上傳檔案"
		} else {
			isSuccess = false
			message = fmt.Sprintf("Error : %s", err.Error())
		}

		ginContextPointer.JSON(http.StatusOK, gin.H{
			"isSuccess": isSuccess,
			"results":   message,
			"appsInfo":  results[0],
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

	out, err := os.Create(path + fileName) // 建立空檔
	defer out.Close()

	if err != nil {
		log.Fatal(err) // Log待補
		fmt.Println(err)
	}

	_, err = io.Copy(out, file) // 將file數據複製到空檔

	if err != nil {
		log.Fatal(err) // Log待補
		fmt.Println(err)
	}

}
