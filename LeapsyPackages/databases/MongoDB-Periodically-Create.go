package databases

import (
	"leapsy.com/records"
)

// insertOneSecondRecord - 加添一筆秒紀錄
/**
 * @param  records.SecondRecord secondRecord  秒紀錄
 */
func (mongoDB *MongoDB) insertOneSecondRecord(secondRecord records.SecondRecord) {

	// mongoClientPointer := mongoDB.Connect() // 資料庫指標

	// if nil != mongoClientPointer { // 若資料庫指標不為空
	// 	defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

	// 	// 預設主機
	// 	address := fmt.Sprintf(
	// 		`%s:%d`,
	// 		mongoDB.GetConfigValueOrPanic(`server`),
	// 		mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	// 	)

	// 	defaultArgs := network.GetAliasAddressPair(address) // 預設參數

	// 	periodicallyRWMutex.Lock() // 寫鎖

	// 	// 添寫秒記錄
	// 	_, insertOneError := mongoClientPointer.
	// 		Database(mongoDB.GetConfigValueOrPanic(`database`)).
	// 		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
	// 		InsertOne(context.TODO(), secondRecord)

	// 	periodicallyRWMutex.Unlock() // 寫解鎖

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 添寫秒記錄 %+v `},
	// 		append(defaultArgs, secondRecord),
	// 		insertOneError,
	// 	)

	// 	if nil != insertOneError { // 若添寫秒紀錄錯誤
	// 		logger.Errorf(formatString, args...) // 記錄錯誤
	// 		return                               // 回傳
	// 	}

	// 	logger.Infof(formatString, args...) // 記錄資訊

	// }

}

// InsertOneSecondRecordIfNotExisted - 若不存在則加添秒紀錄
/**
 * @param  records.SecondRecord secondRecord  秒紀錄
 */
func (mongoDB *MongoDB) InsertOneSecondRecordIfNotExisted(secondRecord records.SecondRecord) {

	// mongoClientPointer := mongoDB.Connect() // 資料庫指標

	// if nil != mongoClientPointer { // 若資料庫指標不為空
	// 	defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

	// 	if 0 == mongoDB.CountSecondRecordByTime(secondRecord.Time) {
	// 		mongoDB.insertOneSecondRecord(secondRecord) // 加添一筆秒紀錄
	// 	} else {

	// 		// 預設主機
	// 		address := fmt.Sprintf(
	// 			`%s:%d`,
	// 			mongoDB.GetConfigValueOrPanic(`server`),
	// 			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
	// 		)

	// 		defaultArgs := network.GetAliasAddressPair(address) // 預設參數
	// 		insertOneDuplicateSecondRecordError := errors.New(`秒紀錄重複加添錯誤`)

	// 		// 取得記錄器格式和參數
	// 		// formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		// 	[]string{`%s %s 添寫秒記錄 %+v `},
	// 		// 	append(defaultArgs, secondRecord),
	// 		// 	insertOneDuplicateSecondRecordError,
	// 		// )

	// 		if nil != insertOneDuplicateSecondRecordError { // 若添寫秒紀錄錯誤
	// 			// logger.Errorf(formatString, args...) // 記錄錯誤
	// 			return // 回傳
	// 		}
	// 	}

	// }

}
