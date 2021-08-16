package servers

import (
	"leapsy.com/databases"
	"leapsy.com/packages/logings"
)

// Server - 伺服器
type Server struct {
}

// StartServer - 啟動伺服器
func StartServer() {

	go logings.StartLogging()

	var (
		apiServer APIServer // API伺服器
		// ecsDB               databases.ECSDB      // 來源資料庫
		// ecsAlertDB          databases.ECSAlertDB // 警報來源資料庫
		periodicallyMongoDB databases.MongoDB // 記錄用資料庫
	)

	defer func() {
		apiServer.stop()                             // 結束API伺服器
		stopPeriodicallyRecord(&periodicallyMongoDB) // 結束周期性記錄
		StopServer()                                 // 結束環控伺服器
	}()

	logings.SendLog(
		[]string{`啟動 軟體更新API伺服器 `},
		[]interface{}{},
		nil,
		0,
	)

	go apiServer.start() // 啟動環控API伺服器

	startPeriodicallyRecord(&periodicallyMongoDB) // 開始週期性記錄

}

// StopServer - 結束伺服器
func StopServer() {
	logings.SendLog(
		[]string{`結束 軟體更新API伺服器 `},
		[]interface{}{},
		nil,
		0,
	)
}
