package servers

import (
	"leapsy.com/databases"
	"leapsy.com/packages/logings"
)

// startPeriodicallyRecord - 開始週期性記錄
/**
 * @param  *databases.ECSDB eCSDB 來源資料庫
 * @param *databases.ECSAlertDB ecsAlertDB 警報來源資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */

// func startPeriodicallyRecord(mongoDB *databases.MongoDB) {

// 	logings.SendLog(
// 		[]string{`啟動 週期性記錄 `},
// 		[]interface{}{},
// 		nil,
// 		0,
// 	)

// 	var updateAlertsMutex sync.Mutex

// 	newCron := cron.New() // 新建一個定時任務物件

// 	newCron.AddFunc(
// 		`@every 15s`,
// 		func() {
// 			updateAlertsMutex.Lock()

// 			// ecsAlertDBRecordsCount := ecsAlertDB.CountAll()
// 			// mongoDBAllAlertRecordsCount := mongoDB.CountAllAlertRecords()

// 			// if ecsAlertDBRecordsCount > mongoDBAllAlertRecordsCount { // 若環控資料庫警報數 > 警報數

// 			// 	for _, ecsAlertRecord := range ecsAlertDB.ReadLast(ecsAlertDBRecordsCount - mongoDBAllAlertRecordsCount) {
// 			// 		mongoDB.RepsertAlertRecord(ecsAlertRecord.AlertRecord())
// 			// 	}

// 			// }

// 			updateAlertsMutex.Unlock()
// 		},
// 	) // 給物件增加定時任務

// 	mongoDB.DeleteAllAlertRecords()

// 	newCron.Start()

// 	newCron2 := cron.New() // 新建一個定時任務物件

// 	countOfSeconds := 15

// 	newCron2.AddFunc(
// 		fmt.Sprintf(`@every %ds`, countOfSeconds),
// 		func() {
// 			currentTime := time.Now()
// 			currentTime =
// 				time.Date(
// 					currentTime.Year(),
// 					currentTime.Month(),
// 					currentTime.Day(),
// 					currentTime.Hour(),
// 					currentTime.Minute(),
// 					(currentTime.Second()/countOfSeconds)*countOfSeconds,
// 					0,
// 					time.Local,
// 				) // 修改時間

// 			go func() {
// 				// secondRecord := eCSDB.Read().SecondRecord() // 讀取秒紀錄
// 				// secondRecord.Time = currentTime             // 儲存時間

// 				// mongoDB.InsertOneSecondRecordIfNotExisted(secondRecord) // 若不存在則加添秒紀錄

// 				// if times.IsHour(currentTime) { // 若為整點時間
// 				// 	mongoDB.AggregateRepsertHourlyRecordByTime(currentTime) // 依據時間聚集輸出小時紀錄
// 				// }

// 				// if times.IsDay(currentTime) { // 若為整日時間
// 				// 	mongoDB.AggregateRepsertDailyRecordByTime(currentTime) // 依據時間聚集輸出日紀錄
// 				// }

// 			}()
// 		},
// 	) // 給物件增加定時任務

// 	newCron2.Start()

// 	select {}
// }

// stopPeriodicallyRecord - 結束週期性記錄
/**
 * @param  *databases.ECSDB eCSDB 來源資料庫
 * @param  *databases.ECSAlertDB eCSAlertDB 來源警報資料庫
 * @param  *databases.MongoDB mongoDB 目的資料庫
 */
func stopPeriodicallyRecord(mongoDB *databases.MongoDB) {

	logings.SendLog(
		[]string{`結束 週期性記錄 `},
		[]interface{}{},
		nil,
		0,
	)

}
