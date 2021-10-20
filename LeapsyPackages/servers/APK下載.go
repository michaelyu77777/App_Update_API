package servers

import (
	"fmt"
	"net/http"
	"time"

	"leapsy.com/packages/configurations"

	// "gopkg.in/mgo.v2/bson" //Michael
	"leapsy.com/records"

	"github.com/gin-gonic/gin"
	// "github.com/sirupsen/logrus"
	// "leapsy.com/packages/logings"
	// "leapsy.com/packages/network"

	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
	// "go.mongodb.org/mongo-driver/bson" //Kevin
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// getAPPsAPIHandler - 取得APK檔
/**
 * @param  *APIServer apiServer API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func getAPPsAPIHandler(apiServer *APIServer, ginContextPointer *gin.Context) {

	eventTime := time.Now()

	isStatusBadRequestErrorChannel := make(chan bool, 1)

	isStatusForbiddenErrorChannel := make(chan bool, 1)

	isStatusNotFoundErrorChannel := make(chan bool, 1)

	httpStatusChannel := make(chan int, 1)

	var parameters Parameters

	bindError := ginContextPointer.ShouldBind(&parameters)

	bindURIError := ginContextPointer.ShouldBindUri(&parameters)

	isError := nil != bindError || nil != bindURIError
	isStatusBadRequestErrorChannel <- isError

	if !isError {

		isToWorkChannel := make(chan bool, 1)

		// parametersLabelName := parameters.LabelName
		parametersPackageName := parameters.PackageName

		// go func() {

		// 	isError = !isLowerCaseOrDigit(parametersDownladKeyword)

		// 	isToWorkChannel <- !isError

		// 	isStatusBadRequestErrorChannel <- isError

		// }()

		// go func() {

		// 	isError = !isAuthorized(ginContextPointer)

		// 	isToWorkChannel <- !isError
		// 	isStatusForbiddenErrorChannel <- isError

		// }()

		apkFileName := "" //apk檔名

		// 看檔案是否存在
		go func() {

			// 檔案存在會取得APK檔名
			// isError, apkFileName = isFileNotExistedAndGetApkFileNameByLabelName(parametersPackageName)
			isError, apkFileName = isFileNotExistedAndGetApkFileNameByPackageName(parametersPackageName)
			// isError = isFileNotExisted(parametersDownladKeyword)

			isToWorkChannel <- !isError

			isStatusNotFoundErrorChannel <- isError

		}()

		go func() {

			isToWork := true

			for counter := 1; counter <= 1; counter++ {
				isToWork = isToWork && <-isToWorkChannel
			}

			if isToWork {

				// 下載檔案
				attachApkFile(ginContextPointer, parametersPackageName, apkFileName)
				// attachCybLicenseBin(ginContextPointer, parametersDownladKeyword)

				httpStatusChannel <- http.StatusOK
			}

		}()

	}

	go func() {

		for {

			if <-isStatusBadRequestErrorChannel {
				httpStatusChannel <- http.StatusBadRequest
			}

		}

	}()

	go func() {

		for {

			if <-isStatusForbiddenErrorChannel {
				httpStatusChannel <- http.StatusForbidden
			}

		}

	}()

	go func() {

		for {

			if <-isStatusNotFoundErrorChannel {
				httpStatusChannel <- http.StatusNotFound
			}

		}

	}()

	for {

		httpStatus := <-httpStatusChannel

		SendEvent(
			eventTime,
			ginContextPointer,
			parameters,
			nil,
			httpStatus,
			APIResponse{},
		)

		ginContextPointer.Status(httpStatus)

		return

	}

}

type ReanalyseAPI struct {
	IsSuccess bool
	Results   string
	Data      records.AppsInfo
}

func getApkDetailsByLabelName(labelName string, apkFileName string) (pkgName string, labelNameOut string, versionCode int, versionName string) {

	// 讀取apk
	pkg, _ := apk.OpenFile(configurations.GetConfigValueOrPanic("local", "pathTemp") + labelName + "/" + apkFileName)
	defer pkg.Close()

	// icon image to base64 string
	// icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image
	// fmt.Println("圖標：icon", icon)

	// buf := new(bytes.Buffer)

	// // Option.Quality壓縮品質:範圍1~100 (大小約1kb ~ 10kb)
	// jpeg.Encode(buf, icon, &jpeg.Options{100})
	// // jpeg.Encode(buf, icon, &jpeg.Options{35})

	// imageBit := buf.Bytes()
	// /*Defining the new image size*/

	// photoBase64 := b64.StdEncoding.EncodeToString([]byte(imageBit))
	// fmt.Println("Photo Base64.............................:" + photoBase64)

	// pkgName
	pkgName = pkg.PackageName() // returns the package name
	fmt.Println("pkgName=<" + pkgName + ">")

	resConfigEN := &androidbinary.ResTableConfig{
		Language: [2]uint8{uint8('e'), uint8('n')},
	}

	// appLabel
	labelNameOut, _ = pkg.Label(resConfigEN) // get app label for en translation
	fmt.Println("labelNameOut=<" + labelNameOut + ">")

	// versionCode
	mainfest := pkg.Manifest()
	fmt.Printf("versionCode=<%+v>\n", mainfest.VersionCode)
	vCode, err := mainfest.VersionCode.Int32()
	versionCode = int(vCode) // int32轉成int
	fmt.Printf("versionCode value=<%d>\n", vCode)
	fmt.Println("err=", err)

	// VersionName
	fmt.Printf("VersionName=<%+v> \n", mainfest.VersionName)
	versionName, err = mainfest.VersionName.String()
	fmt.Printf("VersionName value=<%s> \n", versionName)
	fmt.Println("err=", err)

	// mainActivity
	// mainActivity, err := pkg.MainActivity()
	// fmt.Printf("mainActivity = %+v \n", mainActivity)

	return
}

func getApkDetailsInApkTempDirectory(apkFileName string) (err error, message string, pkgName string, labelName string, versionCode int, versionName string) {

	// 讀取apk
	path := configurations.GetConfigValueOrPanic("local", "pathTemp") + apkFileName
	pkg, err := apk.OpenFile(path)
	defer pkg.Close()
	if nil != err {
		//log
		s := fmt.Sprintf("[錯誤]開啟APK錯誤，路徑：%s，Error：%s", path, err.Error())
		message += s
		fmt.Println(s)
		return
	}

	// // icon image to base64 string
	// icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image
	// fmt.Println("圖標：icon", icon)

	// buf := new(bytes.Buffer)

	// // Option.Quality壓縮品質:範圍1~100 (大小約1kb ~ 10kb)
	// jpeg.Encode(buf, icon, &jpeg.Options{100})
	// // jpeg.Encode(buf, icon, &jpeg.Options{35})

	// imageBit := buf.Bytes()
	// /*Defining the new image size*/

	// photoBase64 := b64.StdEncoding.EncodeToString([]byte(imageBit))
	// fmt.Println("Photo Base64.............................:" + photoBase64)

	// pkgName
	pkgName = pkg.PackageName() // returns the package name
	fmt.Println("pkgName=<" + pkgName + ">")

	resConfigEN := &androidbinary.ResTableConfig{
		Language: [2]uint8{uint8('e'), uint8('n')},
	}

	// labelName
	labelName, err = pkg.Label(resConfigEN) // get app label for en translation
	if nil != err {
		//log
		s := fmt.Sprintf("[錯誤]解析APK-Label值，發生錯誤，錯誤：%s。", err.Error())
		message += s
		fmt.Println(s)
		return
	} else {
		s := "[取得labelName]" + labelName
		fmt.Println(s)
		message += s
	}

	// versionCode
	mainfest := pkg.Manifest()
	fmt.Printf("versionCode=<%+v>\n", mainfest.VersionCode)
	vCode, err := mainfest.VersionCode.Int32()
	if nil != err {
		//log
		s := fmt.Sprintf("[錯誤]解析APK-VersionCode值，發生錯誤，錯誤：%s。", err.Error())
		fmt.Println(s)
		message += s
		return
	} else {
		s := "[取得VersionCode]" + fmt.Sprint(vCode)
		fmt.Println(s)
		message += s
	}

	// versionCode:int32轉成int for return
	versionCode = int(vCode)

	// VersionName
	versionName, err = mainfest.VersionName.String()
	if nil != err {
		//log
		s := fmt.Sprintf("[錯誤]解析APK-VersionName值，發生錯誤，錯誤：%s。", err.Error())
		fmt.Println(s)
		message += s
		return
	} else {
		s := "[取得VersionName]" + versionName
		fmt.Println(s)
		message += s
	}

	// mainActivity
	// mainActivity, err := pkg.MainActivity()
	// fmt.Printf("mainActivity = %+v \n", mainActivity)

	return
}

// getMacAddressCybLicenseBinAPIHandler - 取得授權檔
/**
 * @param  *APIServer apiServer API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
// func getMacAddressCybLicenseBinAPIHandler(apiServer *APIServer, ginContextPointer *gin.Context) {

// 	eventTime := time.Now()

// 	isStatusBadRequestErrorChannel := make(chan bool, 1)

// 	isStatusForbiddenErrorChannel := make(chan bool, 1)

// 	isStatusNotFoundErrorChannel := make(chan bool, 1)

// 	httpStatusChannel := make(chan int, 1)

// 	var parameters Parameters

// 	bindError := ginContextPointer.ShouldBind(&parameters)

// 	bindURIError := ginContextPointer.ShouldBindUri(&parameters)

// 	isError := nil != bindError || nil != bindURIError
// 	isStatusBadRequestErrorChannel <- isError

// 	if !isError {

// 		isToWorkChannel := make(chan bool, 1)

// 		parametersMacAddress := parameters.MacAddress

// 		go func() {

// 			isError = !isLowerCaseOrDigit(parametersMacAddress)

// 			isToWorkChannel <- !isError

// 			isStatusBadRequestErrorChannel <- isError

// 		}()

// 		go func() {

// 			isError = !isAuthorized(ginContextPointer)

// 			isToWorkChannel <- !isError
// 			isStatusForbiddenErrorChannel <- isError

// 		}()

// 		go func() {

// 			isError = isFileNotExisted(parametersMacAddress)

// 			isToWorkChannel <- !isError

// 			isStatusNotFoundErrorChannel <- isError

// 		}()

// 		go func() {

// 			isToWork := true

// 			for counter := 1; counter <= 3; counter++ {
// 				isToWork = isToWork && <-isToWorkChannel
// 			}

// 			if isToWork {
// 				attachCybLicenseBin(ginContextPointer, parametersMacAddress)
// 				httpStatusChannel <- http.StatusOK
// 			}

// 		}()

// 	}

// 	go func() {

// 		for {

// 			if <-isStatusBadRequestErrorChannel {
// 				httpStatusChannel <- http.StatusBadRequest
// 			}

// 		}

// 	}()

// 	go func() {

// 		for {

// 			if <-isStatusForbiddenErrorChannel {
// 				httpStatusChannel <- http.StatusForbidden
// 			}

// 		}

// 	}()

// 	go func() {

// 		for {

// 			if <-isStatusNotFoundErrorChannel {
// 				httpStatusChannel <- http.StatusNotFound
// 			}

// 		}

// 	}()

// 	for {

// 		httpStatus := <-httpStatusChannel

// 		SendEvent(
// 			eventTime,
// 			ginContextPointer,
// 			parameters,
// 			nil,
// 			httpStatus,
// 			APIResponse{},
// 		)

// 		ginContextPointer.Status(httpStatus)

// 		return

// 	}

// }
