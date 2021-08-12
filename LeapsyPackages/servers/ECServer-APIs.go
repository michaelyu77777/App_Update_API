package servers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"leapsy.com/packages/configurations"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
)

// ECAPIServer - 環控API伺服器
type ECAPIServer struct {
	server *http.Server // 伺服器指標
}

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (eCAPIServer *ECAPIServer) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(*eCAPIServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// GetConfigPositiveIntValueOrPanic - 取得設定整數值否則結束程式
/**
 * @param  string key  關鍵字
 * @return int 設定資料區塊下關鍵字對應的整數值
 */
func (eCAPIServer *ECAPIServer) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(*eCAPIServer).String(), key) // 回傳取得的設定檔區塊下關鍵字對應的值
}

// start - 啟動環控API伺服器
func (eCAPIServer *ECAPIServer) start() {

	address := fmt.Sprintf(`%s:%d`,
		eCAPIServer.GetConfigValueOrPanic(`host`),
		eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`),
	) // 預設主機

	network.SetAddressAlias(address, `軟體更新API伺服器`) // 設定預設主機別名

	enginePointer := gin.Default()

	// enginePointer.GET(
	// 	`/record`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getRecordAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// enginePointer.GET(
	// 	`/records/:year/:month/:day`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getRecordsDailyAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// 取得所有 apps info
	// enginePointer.GET(
	// 	//`/appsUpdate/appsInfo/all`,
	// 	`/appsUpdate/allAppsInfo`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAllAppsInfoAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// 取得 app info
	// enginePointer.GET(
	// 	`/appsUpdate/appsInfo/:projectName/:appName`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAppsInfoAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	//驗證
	// enginePointer.POST(
	// 	`/appsUpdate/authentication`,
	// 	func(ginContextPointer *gin.Context) {
	// 		postAuthenticationAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// 驗證並取得所有 apps info
	enginePointer.POST(
		`/appsUpdate/postAllAppsInfo`,
		func(ginContextPointer *gin.Context) {
			postAllAppsInfoAPIHandler(eCAPIServer, ginContextPointer)
		},
	)

	// enginePointer.GET(
	// 	`/records/:year/:month`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getRecordsMonthlyAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// enginePointer.GET(
	// 	`/alerts/page/records/count/:pageCount/pages/count`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAlertsPagesPageCountAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// enginePointer.GET(
	// 	`/alerts/pages/:pageNumber/count/:pageCount`,
	// 	func(ginContextPointer *gin.Context) {
	// 		getAlertsPageNumberPageCountAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	// enginePointer.PUT(
	// 	`/alerts/:alertEventID`,
	// 	func(ginContextPointer *gin.Context) {
	// 		putAlertAPIHandler(eCAPIServer, ginContextPointer)
	// 	},
	// )

	apiServerPointer := &http.Server{Addr: address, Handler: enginePointer} // 設定伺服器
	eCAPIServer.server = apiServerPointer                                   // 儲存伺服器指標

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

// stop - 結束環控API伺服器
func (eCAPIServer *ECAPIServer) stop() {

	address := fmt.Sprintf(`%s:%d`,
		eCAPIServer.GetConfigValueOrPanic(`host`),
		eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`),
	) // 預設主機

	if nil == eCAPIServer || nil == eCAPIServer.server {

		logings.SendLog(
			[]string{`%s %s 結束`},
			network.GetAliasAddressPair(address),
			nil,
			0,
		)

		return
	}

	eCAPIServerServerShutdownError := eCAPIServer.server.Shutdown(context.TODO()) // 結束伺服器

	logings.SendLog(
		[]string{`%s %s 結束`},
		network.GetAliasAddressPair(address),
		eCAPIServerServerShutdownError,
		logrus.PanicLevel,
	)

}
