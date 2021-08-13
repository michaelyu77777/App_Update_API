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

// APIServer - 環控API伺服器
type APIServer struct {
	server *http.Server // 伺服器指標
}

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

// start - 啟動環控API伺服器
func (apiServer *APIServer) start() {

	address := fmt.Sprintf(`%s:%d`,
		apiServer.GetConfigValueOrPanic(`host`),
		apiServer.GetConfigPositiveIntValueOrPanic(`port`),
	) // 預設主機

	network.SetAddressAlias(address, `軟體更新API伺服器`) // 設定預設主機別名

	enginePointer := gin.Default()

	// 驗證並取得所有 apps info
	enginePointer.POST(
		`/appsUpdate/postAllAppsInfo`,
		func(ginContextPointer *gin.Context) {
			postAllAppsInfoAPIHandler(apiServer, ginContextPointer)
		},
	)

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

	apiServerPointer := &http.Server{Addr: address, Handler: enginePointer} // 設定伺服器
	apiServer.server = apiServerPointer                                     // 儲存伺服器指標

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
