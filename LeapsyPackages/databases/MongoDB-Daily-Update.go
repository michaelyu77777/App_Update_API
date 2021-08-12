package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
)

// findOneAndReplaceDailyRecord - 代添日記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndReplaceOptions 選項
 * @return *mongo.SingleResult returnSingleResultPointer 更添結果
 */
func (mongoDB *MongoDB) findOneAndReplaceDailyRecord(
	filter, replacement primitive.M,
	opts ...*options.FindOneAndReplaceOptions) (returnSingleResultPointer *mongo.SingleResult) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			mongoDB.GetConfigValueOrPanic(`server`),
			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		dailyRWMutex.Lock() // 寫鎖

		// 更新日記錄
		singleResultPointer := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
			FindOneAndReplace(
				context.TODO(),
				filter,
				replacement,
				opts...,
			)

		dailyRWMutex.Unlock() // 寫解鎖

		var findOneAndReplaceError error // 更添錯誤

		singleResultPointerError := singleResultPointer.Err() // 錯誤

		if mongo.ErrNoDocuments != singleResultPointerError { // 若非檔案不存在錯誤
			findOneAndReplaceError = singleResultPointerError // 更添錯誤
		}

		logings.SendLog(
			[]string{`%s %s 更添日記錄 %+v `},
			append(defaultArgs, replacement),
			findOneAndReplaceError,
			logrus.ErrorLevel,
		)

		if nil != findOneAndReplaceError { // 若代添日紀錄錯誤
			return // 回傳
		}

		returnSingleResultPointer = singleResultPointer // 回傳結果指標

	}

	return // 回傳
}

// repsertOneDailyRecord - 代添日記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return []records.DailyRecord results 更新結果
 */
func (mongoDB *MongoDB) repsertOneDailyRecord(filter, replacement primitive.M) (results []records.DailyRecord) {

	var replacedDailyRecord records.DailyRecord // 更新的紀錄

	resultPointer :=
		mongoDB.
			findOneAndReplaceDailyRecord(
				filter,
				replacement,
				options.FindOneAndReplace().SetUpsert(true),
			)

	if nil != resultPointer &&
		nil ==
			resultPointer.
				Decode(&replacedDailyRecord) { // 若更新沒錯誤
		results = append(results, replacedDailyRecord) // 回傳結果
	}

	return // 回傳
}
