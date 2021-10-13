package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"leapsy.com/packages/configurations"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/packages/paths"

	"leapsy.com/records"
)

// APIServer - 軟體更新API伺服器
type APIServer struct {
	server *http.Server // 伺服器指標
}

// APIResponse - API回應
type APIResponse struct {
	IsSuccess bool          `json:"isSuccess"` // 是否成功
	Results   []interface{} `json:"results"`   // 訊息
}

const ()

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
// func (apiServer *APIServer) GetConfigValueOrPanic(key string) string {
// 	return configurations.GetConfigValueOrPanic(reflect.TypeOf(*apiServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
// }

// GetConfigPositiveIntValueOrPanic - 取得設定整數值否則結束程式
/**
 * @param  string key  關鍵字
 * @return int 設定資料區塊下關鍵字對應的整數值
 */
// func (apiServer *APIServer) GetConfigPositiveIntValueOrPanic(key string) int {
// 	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(*apiServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
// }

// start - 啟動環控API伺服器
func (apiServer *APIServer) start() {

	address := fmt.Sprintf(`%s:%d`,
		apiServer.GetConfigValueOrPanic(`host`),
		apiServer.GetConfigPositiveIntValueOrPanic(`port`),
	) // 預設主機

	network.SetAddressAlias(address, `軟體更新API伺服器`) // 設定預設主機別名

	// 後來加上的
	gin.SetMode(gin.ReleaseMode)

	enginePointer := gin.Default()

	// 驗證並取得所有 apps info
	enginePointer.POST(
		`/appsUpdate/postAllAppsInfo`,
		func(ginContextPointer *gin.Context) {
			postAllAppsInfoAPIHandler(apiServer, ginContextPointer)
		},
	)

	// 上傳APK
	enginePointer.POST(
		`/appUpdate/file/uploadSingle`,
		func(ginContextPointer *gin.Context) {
			UploadSingleApkAPIHandler(apiServer, ginContextPointer)
		},
	)

	// 下載APK
	enginePointer.GET(
		// `/appUpdate/download/:downloadKeyword`,
		`/appUpdate/download/:labelName`,
		func(ginContextPointer *gin.Context) {
			getAPPsAPIHandler(apiServer, ginContextPointer)
		},
	)

	// 已經有上傳功能直接棄用
	// 重新解析APK某資料夾的檔案
	// enginePointer.POST(
	// 	`/appUpdate/postReanalyse`,
	// 	func(ginContextPointer *gin.Context) {
	// 		postReanalyseAPIHandler(apiServer, ginContextPointer)
	// 	},
	// )

	// 取得所有 apps info
	// enginePointer.GET(
	// 	//`/appsUpdate/appsInfo/all`,
	// 	`/appsUpdate/allAppsInfo`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAllAppsInfoAPIHandler(APIServer, ginContextPointer)
	// 	},
	// )

	// 取得 app info
	// enginePointer.GET(
	// 	`/appsUpdate/appsInfo/:projectName/:appName`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAppsInfoAPIHandler(APIServer, ginContextPointer)
	// 	},
	// )

	//驗證
	// enginePointer.POST(
	// 	`/appsUpdate/authentication`,
	// 	func(ginContextPointer *gin.Context) {
	// 		postAuthenticationAPIHandler(APIServer, ginContextPointer)
	// 	},
	// )

	// enginePointer.PUT(
	// 	`/alerts/:alertEventID`,
	// 	func(ginContextPointer *gin.Context) {
	// 		putAlertAPIHandler(APIServer, ginContextPointer)
	// 	},
	// )

	apiServerPointer := &http.Server{
		Addr:    address,
		Handler: enginePointer,
	} // 設定伺服器
	apiServer.server = apiServerPointer // 儲存伺服器指標

	var apiServerPtrListenAndServeError error // 伺服器啟動錯誤

	go func() {
		apiServerPtrListenAndServeError = apiServerPointer.ListenAndServe() // 啟動伺服器或回傳伺服器啟動錯誤
	}()

	<-time.After(time.Second * 3) // 等待伺服器啟動結果

	logings.SendLog(
		[]string{`%s %s 啟動`},
		network.GetAliasAddressPair(address),
		apiServerPtrListenAndServeError,
		logrus.PanicLevel,
	)

}

// stop - 結束API伺服器
func (apiServer *APIServer) stop() {

	address := fmt.Sprintf(`%s:%d`,
		apiServer.GetConfigValueOrPanic(`host`),
		apiServer.GetConfigPositiveIntValueOrPanic(`port`),
	) // 預設主機

	if nil == apiServer || nil == apiServer.server {

		logings.SendLog(
			[]string{`%s %s 結束`},
			network.GetAliasAddressPair(address),
			nil,
			0,
		)

		return
	}

	eCAPIServerServerShutdownError := apiServer.server.Shutdown(context.TODO()) // 結束伺服器

	logings.SendLog(
		[]string{`%s %s 結束`},
		network.GetAliasAddressPair(address),
		eCAPIServerServerShutdownError,
		logrus.PanicLevel,
	)

}

/**以下函數為軟體下載所需函數**/

// Parameters - URL參數
type Parameters struct {
	// MacAddress string `uri:"macAddress"`
	// DownloadKeyword string `uri:"downloadKeyword"`
	LabelName string `uri:"labelName"`
}

// // APIResponse - API回應
// type APIResponse struct {
// 	IsSuccess bool          `json:"isSuccess"` // 是否成功
// 	Results   []interface{} `json:"results"`   // 訊息
// }

// // APIServer - API伺服器
// type APIServer struct {
// 	server *http.Server // 伺服器指標
// }

// EventData - 事件資料
type EventData struct {
	Time              time.Time
	GinContextPointer *gin.Context `json:"-"`
	Input, Output     interface{}
	Status            int
	APIResponse       APIResponse
}

const ()

var (
	eventDataChannel = make(chan EventData, configurations.GetConfigPositiveIntValueOrPanic(`local`, `channel-size`))
)

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (apiServer *APIServer) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(*apiServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// GetConfigPositiveIntValueOrPanic - 取得設定整數值否則結束程式
/**
 * @param  string key  關鍵字
 * @return int 設定資料區塊下關鍵字對應的整數值
 */
func (apiServer *APIServer) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(*apiServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// start - 啟動API伺服器
// func (apiServer *APIServer) start() {

// 	address := fmt.Sprintf(`%s:%d`,
// 		apiServer.GetConfigValueOrPanic(`host`),
// 		apiServer.GetConfigPositiveIntValueOrPanic(`port`),
// 	) // 預設主機

// 	network.SetAddressAlias(address, `API伺服器`) // 設定預設主機別名

// 	gin.SetMode(gin.ReleaseMode)

// 	enginePointer := gin.Default()

// 	enginePointer.GET(
// 		`/appUpdate/download/:downloadKeyword`,
// 		func(ginContextPointer *gin.Context) {
// 			getAPPsAPIHandler(apiServer, ginContextPointer)
// 		},
// 	)

// 	// 重新解析APK某資料夾的檔案
// 	enginePointer.POST(
// 		`/appUpdate/postReanalyse`,
// 		func(ginContextPointer *gin.Context) {
// 			postReanalyseAPIHandler(apiServer, ginContextPointer)
// 		},
// 	)

// 	// enginePointer.GET(
// 	// 	`/:macAddress/CybLicense.bin`,
// 	// 	func(ginContextPointer *gin.Context) {
// 	// 		getMacAddressCybLicenseBinAPIHandler(apiServer, ginContextPointer)
// 	// 	},
// 	// )

// 	// enginePointer.PUT(
// 	// 	`/:macAddress/CybLicense.bin`,
// 	// 	func(ginContextPointer *gin.Context) {
// 	// 		putMacAddressCybLicenseBinAPIHandler(apiServer, ginContextPointer)
// 	// 	},
// 	// )

// 	// enginePointer.DELETE(
// 	// 	`/:macAddress/CybLicense.bin`,
// 	// 	func(ginContextPointer *gin.Context) {
// 	// 		deleteMacAddressCybLicenseBinAPIHandler(apiServer, ginContextPointer)
// 	// 	},
// 	// )

// 	apiServerPointer := &http.Server{
// 		Addr:    address,
// 		Handler: enginePointer,
// 	} // 設定伺服器

// 	apiServer.server = apiServerPointer // 儲存伺服器指標

// 	var apiServerPtrListenAndServeError error // 伺服器啟動錯誤

// 	go func() {
// 		apiServerPtrListenAndServeError = enginePointer.Run(address) // 啟動伺服器或回傳伺服器啟動錯誤

// 		// const (
// 		// 	certPEMFileName = `cert.pem`
// 		// 	keyPEMFileName  = `key.pem`
// 		// )

// 		// tlss.CreateCertAndKeyPEMFiles(certPEMFileName, keyPEMFileName)

// 		// apiServerPtrListenAndServeError = enginePointer.RunTLS(address, certPEMFileName, keyPEMFileName) // 啟動伺服器或回傳伺服器啟動錯誤

// 		// apiServerPtrListenAndServeError = enginePointer.RunTLS(address, configurations.GetConfigValueOrPanic(`local`, `cert-pem-path`), configurations.GetConfigValueOrPanic(`local`, `private-key-pem-path`))
// 	}()

// 	<-time.After(time.Second * 3) // 等待伺服器啟動結果

// 	logings.SendLog(
// 		[]string{`%s %s 啟動 `},
// 		network.GetAliasAddressPair(address),
// 		apiServerPtrListenAndServeError,
// 		logrus.PanicLevel,
// 	)

// 	select {}

// }

// // stop - 結束API伺服器
// func (apiServer *APIServer) stop() {

// 	address := fmt.Sprintf(`%s:%d`,
// 		apiServer.GetConfigValueOrPanic(`host`),
// 		apiServer.GetConfigPositiveIntValueOrPanic(`port`),
// 	) // 預設主機

// 	logings.SendLog(
// 		[]string{`%s %s 結束 `},
// 		network.GetAliasAddressPair(address),
// 		nil,
// 		logrus.InfoLevel,
// 	)

// 	if nil == apiServer || nil == apiServer.server {
// 		return
// 	}

// 	apiServerServerShutdownError := apiServer.server.Shutdown(context.TODO()) // 結束伺服器

// 	logings.SendLog(
// 		[]string{`%s %s 結束 `},
// 		network.GetAliasAddressPair(address),
// 		apiServerServerShutdownError,
// 		logrus.PanicLevel,
// 	)

// }

// SendEvent - 傳送事件
/**
 * @param time.Time time 事件時間
 * @param *gin.Context eventGinContextPointer  事件gin Context指標
 * @param []interface{} eventInput 事件輸入
 * @param []interface{} eventOutput 事件輸出
 * @param  APIResponse eventAPIResponse 事件API回應
 */
func SendEvent(
	eventTime time.Time,
	eventGinContextPointer *gin.Context,
	eventInput, eventOutput interface{},
	eventStatus int,
	eventAPIResponse APIResponse,
) {
	eventDataChannel <- EventData{
		Time:              eventTime,
		GinContextPointer: eventGinContextPointer,
		Input:             eventInput,
		Output:            eventOutput,
		Status:            eventStatus,
		APIResponse:       eventAPIResponse,
	}
}

// StartUpdatingEvents - 開始更新事件
func StartUpdatingEvents() {

	logings.SendLog(
		[]string{`啟動 %s 更新事件 `},
		[]interface{}{`servers`},
		nil,
		0,
	)

	for {

		eventData := <-eventDataChannel

		go func() {

			eventDataGinContextPointer := eventData.GinContextPointer

			var event records.Event

			eventDataGinContextPointer.ShouldBind(&event)

			if jsonBytes, jsonMarshalError := json.Marshal(eventData); jsonMarshalError == nil {

				if jsonUnmarshalError := json.Unmarshal(jsonBytes, &event); jsonUnmarshalError != nil {

					logings.SendLog(
						[]string{`將JSON字串 %s 轉成 物件 %+v `},
						[]interface{}{string(jsonBytes), event},
						jsonUnmarshalError,
						logrus.PanicLevel,
					)

				}

			} else {

				logings.SendLog(
					[]string{`將物件 %+v 轉成 JSON字串 %s `},
					[]interface{}{eventData, string(jsonBytes)},
					jsonMarshalError,
					logrus.PanicLevel,
				)

			}

			eventDataGinContextPointerRequest := eventDataGinContextPointer.Request
			event.Method = eventDataGinContextPointerRequest.Method
			event.Header = eventDataGinContextPointerRequest.Header

			event.IP = eventDataGinContextPointer.ClientIP()
			event.URL = eventDataGinContextPointer.FullPath()

			event.Results = eventData.APIResponse.Results

			mongoDB.RepsertEvent(event)

		}()

	}

}

// isLowerCaseOrDigit - 判斷是否小寫或數字
/**
 * @param  string inputString 輸入字串
 * @return bool result 結果
 */
func isLowerCaseOrDigit(inputString string) (result bool) {

	result, _ = regexp.MatchString(`^[0-9a-z]+$`, inputString)

	return
}

// isFileNotExisted - 判斷./apk/下的指定的檔案是否「不存在」
/**
 * @param  string downloadKeyword apps代號
 * @return bool result 結果
 */
func isFileNotExistedByLabelName(labelName string) (result bool) {

	// 查資料夾是否不在
	directoryName := labelName // 指定apk資料夾名稱

	// 去資料庫查此資料夾名稱所對應的APK檔名
	appsInfo := mongoDB.FindAppsInfoByLabelName(labelName)

	// 找不到此APP
	if 1 > len(appsInfo) {
		fmt.Printf("找不到名稱為[ %s ]的資料夾之APP\n", directoryName)
		result = true
		return

	} else {

		// 取出APK檔名
		apkFileName := appsInfo[0].ApkFileName

		// 看APK檔案存不存在
		result = paths.IsFileNotExisted(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + directoryName + "/" + apkFileName)
		fmt.Println("檔案或路徑是否不存在？", result)
		return
	}

	// 查檔案是否不在
	// result = paths.IsFileNotExisted(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + downloadKeyword + "/camera.apk")
	// fmt.Print("路徑=", paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`))+downloadKeyword+"/camera.apk")

}

// isFileNotExisted - 判斷是否檔案不存在，回傳結果與檔名
/**
 * @param  string downloadKeyword apps代號
 * @return bool result 結果 apkFileName APK檔名
 */
func isFileNotExistedAndGetApkFileNameByLabelName(labelName string) (result bool, apkFileName string) {

	// 去資料庫查此資料夾名稱所對應的APK檔名
	appsInfo := mongoDB.FindAppsInfoByLabelName(labelName)

	// 找不到此APP
	if 1 > len(appsInfo) {

		detail := "資料庫中找不到存放資料夾名稱為<" + labelName + ">的APP"

		// log
		logings.SendLog(
			[]string{detail},
			[]interface{}{labelName},
			nil,
			logrus.WarnLevel,
		)

		fmt.Printf(detail+"\n", labelName)
		result = true
		return

	} else {

		// 取出APK檔名
		apkFileName = appsInfo[0].ApkFileName

		fmt.Println("查詢labelName=" + labelName)
		fmt.Println("取出ApkFileName=" + apkFileName)
		fmt.Printf("取出appsInfo[0]= %+v", appsInfo[0])

		// 看APK檔案存不存在
		result = paths.IsFileNotExisted(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + labelName + "/" + apkFileName)
		// path := paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + "相簿/" + apkFileName
		// result = paths.IsFileNotExisted(path)
		// fmt.Println("path=", path)
		fmt.Println("檔案或路徑是否不存在？", result)
		return
	}

	// 查檔案是否不在
	// result = paths.IsFileNotExisted(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + downloadKeyword + "/camera.apk")
	// fmt.Print("路徑=", paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`))+downloadKeyword+"/camera.apk")

}

// isAuthorized - 判斷是否認證通過
/**
 * @param  *gin.Context ginContextPointer  gin Context 指標
 * @return bool result 結果
 */
func isAuthorized(ginContextPointer *gin.Context) (result bool) {

	result, _ = regexp.MatchString(`^Basic\s+TGVhcHN5Vm9pY2VTZXJ2aWNlOjZXbUJuZ2R1SEZwc1I0eGRiSnU0ajZZdlNuV2VZYzdq$`, ginContextPointer.Request.Header.Get(`Authorization`))

	return
}

// attachCybLicenseBin - 附加apk檔案（檔案下載）
/**
 * @param  *gin.Context ginContextPointer  gin Context 指標
 * @param  string fileNameString 檔名字串
 */
func attachCybLicenseBin(ginContextPointer *gin.Context, labelName string) {

	// 去資料庫查此資料夾名稱所對應的APK檔名
	result := mongoDB.FindAppsInfoByLabelName(labelName)
	// apkFileName := `camera.apk`

	// 沒找到
	if 1 > len(result) {
		fmt.Printf("找不到此資料夾 \n")
	} else {

		// 取得APK檔名
		apkFileName := result[0].ApkFileName

		// 設定下載後檔名
		downloadFileName := apkFileName

		// 下載
		ginContextPointer.FileAttachment(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`))+labelName+`/`+apkFileName, downloadFileName)
	}

}

// attachCybLicenseBin - 附加apk檔案（檔案下載）
/**
 * @param  *gin.Context ginContextPointer  gin Context 指標 directoryName 想下載的apk資料夾名稱 apkFileName 想下載的APK檔案名稱
 * @param  string fileNameString 檔名字串
 */
func attachApkFile(ginContextPointer *gin.Context, labelName string, apkFileName string) {

	// 下載檔案名稱
	downloadFileName := apkFileName

	// 下載
	ginContextPointer.FileAttachment(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`))+labelName+`/`+apkFileName, downloadFileName)

}

// upsertCybLicenseBin - 更添授權檔
/**
 * @param  *gin.Context ginContextPointer  gin Context 指標
 * @param  string fileNameString 檔名字串
 * @return bool isSuccess 是否成功
 */
func upsertCybLicenseBin(ginContextPointer *gin.Context, fileNameString string) (isSuccess bool) {

	path := paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`))
	paths.CreateIfPathNotExisted(path)

	if destinationFilePointer, osCreateError := os.Create(path + fileNameString); osCreateError == nil {

		defer destinationFilePointer.Close()

		_, ioCopyError := io.Copy(destinationFilePointer, ginContextPointer.Request.Body)
		isSuccess = ioCopyError == nil

	}

	return
}

// deleteCybLicenseBin - 刪除授權檔
/**
 * @param  string fileNameString 檔名字串
 */
func deleteCybLicenseBin(fileNameString string) {
	os.Remove(paths.AppendSlashIfNotEndWithOne(configurations.GetConfigValueOrPanic(`local`, `path`)) + fileNameString)
}
