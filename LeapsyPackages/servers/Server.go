package servers

import (
	"leapsy.com/packages/logings"
)

// Server - 伺服器
type Server struct {
}

// StartServer - 啟動伺服器
func StartServer() {

	// log
	go logings.StartLogging()

	// 啟動下載事件紀錄
	go StartUpdatingEvents()

	var (
		apiServer APIServer // API伺服器
		// ecsDB               databases.ECSDB      // 來源資料庫
		// ecsAlertDB          databases.ECSAlertDB // 警報來源資料庫
		// periodicallyMongoDB databases.MongoDB // 記錄用資料庫
	)

	// defer: StartServer()整個函數結束之前，執行此函數
	defer func() {
		apiServer.stop() // 結束API伺服器
		StopServer()     // 結束整個伺服器
		// stopPeriodicallyRecord(&periodicallyMongoDB) // 結束周期性記錄
	}()

	// log
	logings.SendLog(
		[]string{`啟動 ` + apiServer.GetConfigValueOrPanic("nameOfServer") + ` `},
		[]interface{}{},
		nil,
		0,
	)

	// 啟動API伺服器
	go apiServer.start()
	// startPeriodicallyRecord(&periodicallyMongoDB) // 開始週期性記錄

	// (主線無限回圈)
	select {}

}

// StopServer - 結束伺服器
func StopServer() {

	var (
		apiServer APIServer // API伺服器
	)
	logings.SendLog(
		[]string{`結束 ` + apiServer.GetConfigValueOrPanic("nameOfServer") + ` `},
		[]interface{}{},
		nil,
		0,
	)
}
