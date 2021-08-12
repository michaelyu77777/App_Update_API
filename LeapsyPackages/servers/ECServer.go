package servers

import (
	"leapsy.com/databases"
	"leapsy.com/packages/logings"
)

// ECServer - 環控伺服器
type ECServer struct {
}

// StartECServer - 啟動環控伺服器
func StartECServer() {

	go logings.StartLogging()

	var (
		eCAPIServer ECAPIServer // 環控API伺服器
		// ecsDB               databases.ECSDB      // 來源資料庫
		// ecsAlertDB          databases.ECSAlertDB // 警報來源資料庫
		periodicallyMongoDB databases.MongoDB // 記錄用資料庫
	)

	defer func() {
		eCAPIServer.stop()                           // 結束環控API伺服器
		stopPeriodicallyRecord(&periodicallyMongoDB) // 結束周期性記錄
		StopECServer()                               // 結束環控伺服器
	}()

	logings.SendLog(
		[]string{`啟動 軟體更新API伺服器 `},
		[]interface{}{},
		nil,
		0,
	)

	go eCAPIServer.start() // 啟動環控API伺服器

	startPeriodicallyRecord(&periodicallyMongoDB) // 開始週期性記錄

}

// StopECServer - 結束環控伺服器
func StopECServer() {
	logings.SendLog(
		[]string{`結束 環控伺服器 `},
		[]interface{}{},
		nil,
		0,
	)
}
