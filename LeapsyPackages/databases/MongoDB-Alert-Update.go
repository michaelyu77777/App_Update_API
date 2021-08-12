package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
)

// findOneAndUpdateAlertRecord - 更添警報記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndUpdateOptions 選項
 * @return result []records.AlertRecord 更添結果
 */
func (mongoDB *MongoDB) findOneAndUpdateAlertRecord(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (results []records.AlertRecord) {

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

		alertRWMutex.Lock() // 寫鎖

		// 更新警報記錄
		singleResultPointer := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
			FindOneAndUpdate(
				context.TODO(),
				filter,
				update,
				opts...,
			)

		alertRWMutex.Unlock() // 寫解鎖

		findOneAndUpdateError := singleResultPointer.Err() // 更添錯誤

		if nil != findOneAndUpdateError && mongo.ErrNoDocuments != findOneAndUpdateError { // 若更添警報紀錄錯誤且非檔案不存在錯誤

			logings.SendLog(
				[]string{`%s %s 更添警報記錄 %+v `},
				append(defaultArgs, update),
				findOneAndUpdateError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 更添警報記錄 %+v `},
			append(defaultArgs, update),
			findOneAndUpdateError,
			0,
		)

		results = mongoDB.findAlertRecords(filter)

	}

	return
}

// UpdateOneAlertRecord - 更新警報記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAlertRecord(filter, update primitive.M) (results []records.AlertRecord) {

	updatedAlertRecords := mongoDB.findOneAndUpdateAlertRecord(filter, update) // 更新的紀錄

	if nil != updatedAlertRecords { // 若更新沒錯誤
		results = append(results, updatedAlertRecords...) // 回傳結果
	}

	return
}

/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndReplaceOptions 選項
 * @return *mongo.SingleResult returnSingleResultPointer 更添結果
 */
func (mongoDB *MongoDB) findOneAndReplaceAlertRecord(
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

		// 表格
		collection := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`))

		periodicallyRWMutex.Lock() // 寫鎖

		collection.
			Indexes().
			CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys: alertRecordSortPrimitiveM,
				},
			)

		// 更新秒記錄
		singleResultPointer := collection.
			FindOneAndReplace(
				context.TODO(),
				filter,
				replacement,
				opts...,
			)

		periodicallyRWMutex.Unlock() // 寫解鎖

		var findOneAndReplaceError error // 更添錯誤

		singleResultPointerError := singleResultPointer.Err() // 錯誤

		if mongo.ErrNoDocuments != singleResultPointerError { // 若非檔案不存在錯誤
			findOneAndReplaceError = singleResultPointerError // 更添錯誤
		}

		logings.SendLog(
			[]string{`%s %s 更添秒記錄`},
			network.GetAliasAddressPair(address),
			findOneAndReplaceError,
			logrus.ErrorLevel,
		)

		if nil != findOneAndReplaceError { // 若代添秒紀錄錯誤
			return // 回傳
		}

		returnSingleResultPointer = singleResultPointer // 回傳結果指標

	}

	return // 回傳
}

// repsertOneAlertRecord - 代添秒記錄
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return []records.AlertRecord results 更新結果
 */
func (mongoDB *MongoDB) repsertOneAlertRecord(filter, replacement primitive.M) (results []records.AlertRecord) {

	var replacedAlertRecord records.AlertRecord // 更新的紀錄

	resultPointer :=
		mongoDB.
			findOneAndReplaceAlertRecord(
				filter,
				replacement,
				options.FindOneAndReplace().SetUpsert(true),
			)

	if nil != resultPointer &&
		nil ==
			resultPointer.
				Decode(&replacedAlertRecord) { // 若更新沒錯誤
		results = append(results, replacedAlertRecord) // 回傳結果
	}

	return
}

// RepsertAlertRecord - 代添秒紀錄
/**
 * @param records.AlertRecord alertRecord 秒紀錄
 * @return []records.AlertRecord results 回傳結果
 */
func (mongoDB *MongoDB) RepsertAlertRecord(alertRecord records.AlertRecord) (results []records.AlertRecord) {

	// 代添紀錄
	results = mongoDB.repsertOneAlertRecord(
		bson.M{
			`alerteventid`: alertRecord.AlertEventID,
		},
		alertRecord.PrimitiveM(),
	)

	return
}
