package servers

import (
	"fmt"
	"net/http"
	"time"

	"leapsy.com/records"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"

	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
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

		parametersDownladKeyword := parameters.DownloadKeyword

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

		apkDirectoryName := parametersDownladKeyword // apk資料夾名稱
		apkFileName := ""                            //apk檔名

		// 看檔案是否存在
		go func() {

			// 檔案存在會取得APK檔名
			isError, apkFileName = isFileNotExistedAndGetApkFileNameByDirectoryName(apkDirectoryName)
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
				attachApkFile(ginContextPointer, apkDirectoryName, apkFileName)
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

// postReanalyseAPIHandler 重新解析某資料夾下的APK
/**
 * @param  *APIServer apiServer API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func postReanalyseAPIHandler(apiServer *APIServer, ginContextPointer *gin.Context) {

	// 參數格式
	type Parameters struct {

		// 指定資料夾名稱
		ApkDirectoryName string `form:"apkDirectoryName" json:"apkDirectoryName" binding:"required"`
	}

	// 接收客戶端之參數
	var parameters Parameters

	// 轉譯json參數
	bindJSONError := ginContextPointer.ShouldBindJSON(&parameters)

	bindURIError := ginContextPointer.ShouldBindUri(&parameters)

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
			parameters,
		)

	// log
	logings.SendLog(
		[]string{`%s %s 接受 %s 請求 %s %+v `},
		defaultArgs,
		nil,
		0,
	)

	// 取得各參數值
	parametersApkDirectoryName := parameters.ApkDirectoryName

	fmt.Println("測試：已取得參數 parametersApkDirectoryName=", parametersApkDirectoryName)

	// 若順利取出 則進行密碼驗證
	if bindJSONError == nil && bindURIError == nil {

		fmt.Println("取參數正確")

		isNotExitsted, apkFileName := isFileNotExistedAndGetApkFileNameByDirectoryName(parametersApkDirectoryName)

		// 檔案若不在
		if isNotExitsted {

			// Response
			responseResult := ReanalyseAPI{
				IsSuccess: false,
				Results:   "此資料夾或APK檔不存在",
				Data:      records.AppsInfo{},
			}

			ginContextPointer.JSON(http.StatusOK, responseResult)

		} else {
			// 檔案存在

			// 進行APK分析
			pkgName, appLabel, versionCode, VersionName := getApkDetails(parametersApkDirectoryName, apkFileName)
			fmt.Printf("解析並取得APK詳細資料：pkgName=%s, appLabel=%s, versionCode=%d, VersionName=%s \n\n", pkgName, appLabel, versionCode, VersionName)

			// 將分析資料更新到資料庫中

			// 重找新找此筆資料
			results := mongoDB.FindAppsInfoByApkDirectoryName(parametersApkDirectoryName)

			// Response
			responseResult := ReanalyseAPI{
				IsSuccess: true,
				Results:   "解析成功",
				Data:      results[0],
			}

			ginContextPointer.JSON(http.StatusOK, responseResult)
		}

	} else if bindJSONError != nil {

		fmt.Println("取參數錯誤,錯誤訊息:bindJSONError=", bindJSONError, ",bindURIError=", bindURIError)

		// 包成回給前端的格式
		myResult := ReanalyseAPI{
			IsSuccess: false,
			Results:   "解析失敗",
			Data:      records.AppsInfo{},
		}

		// 回應給前端
		ginContextPointer.JSON(http.StatusNotFound, myResult)

		// log

		logings.SendLog(
			[]string{`%s %s 回應 %s 請求 %s %+v: 驗證失敗-取參數錯誤(參數有少或格式錯誤), bindJSONError=%s, bindURIError=%s`},
			append(
				defaultArgs,
				bindJSONError,
				bindURIError,
			),
			nil,              // 無錯誤
			logrus.InfoLevel, // info等級的log
		)
	}

}

func getApkDetails(apkDirectoryName string, apkFileName string) (pkgName string, appLabel string, versionCode int, versionName string) {

	// 讀取apk
	pkg, _ := apk.OpenFile("./apk/" + apkDirectoryName + "/" + apkFileName)
	defer pkg.Close()

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

	// appLabel
	appLabel, _ = pkg.Label(resConfigEN) // get app label for en translation
	fmt.Println("appLabel=<" + appLabel + ">")

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
