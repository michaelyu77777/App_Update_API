package databases

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"leapsy.com/records"
)

// findOneAndReplaceHourlyRecord - 代添小時記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndReplaceOptions 選項
 * @return *mongo.SingleResult returnSingleResultPointer 更添結果
 */
func (mongoDB *MongoDB) findOneAndReplaceHourlyRecord(
	filter, replacement primitive.M,
	opts ...*options.FindOneAndReplaceOptions) (returnSingleResultPointer *mongo.SingleResult) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// 預設主機
		// address := fmt.Sprintf(
		// 	`%s:%d`,
		// 	mongoDB.GetConfigValueOrPanic(`server`),
		// 	mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
		// )

		// defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		// hourlyRWMutex.Lock() // 寫鎖

		// // 更新小時記錄
		// singleResultPointer := mongoClientPointer.
		// 	Database(mongoDB.GetConfigValueOrPanic(`database`)).
		// 	Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
		// 	FindOneAndReplace(
		// 		context.TODO(),
		// 		filter,
		// 		replacement,
		// 		opts...,
		// 	)

		// hourlyRWMutex.Unlock() // 寫解鎖

		// var findOneAndReplaceError error // 更添錯誤

		// singleResultPointerError := singleResultPointer.Err() // 錯誤

		// if mongo.ErrNoDocuments != singleResultPointerError { // 若非檔案不存在錯誤
		// 	findOneAndReplaceError = singleResultPointerError // 更添錯誤
		// }

		// // 取得記錄器格式和參數
		// formatString, args := logings.GetLogFuncFormatAndArguments(
		// 	[]string{`%s %s 更添小時記錄 %+v `},
		// 	append(defaultArgs, replacement),
		// 	findOneAndReplaceError,
		// )

		// if nil != findOneAndReplaceError { // 若代添小時紀錄錯誤
		// 	logger.Errorf(formatString, args...) // 記錄錯誤
		// 	return                               // 回傳
		// }

		// go logger.Infof(formatString, args...)          // 記錄資訊
		// returnSingleResultPointer = singleResultPointer // 回傳結果指標

	}

	return // 回傳
}

// repsertOneHourlyRecord - 代添小時記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return []records.HourlyRecord results 更新結果
 */
func (mongoDB *MongoDB) repsertOneHourlyRecord(filter, replacement primitive.M) (results []records.HourlyRecord) {

	var replacedHourlyRecord records.HourlyRecord // 更新的紀錄

	resultPointer :=
		mongoDB.
			findOneAndReplaceHourlyRecord(
				filter,
				replacement,
				options.FindOneAndReplace().SetUpsert(true),
			)

	if nil != resultPointer &&
		nil ==
			resultPointer.
				Decode(&replacedHourlyRecord) { // 若更新沒錯誤
		results = append(results, replacedHourlyRecord) // 回傳結果
	}

	return
}
