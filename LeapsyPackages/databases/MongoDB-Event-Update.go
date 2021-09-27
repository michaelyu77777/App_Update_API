package databases

import (
	"context"
	"fmt"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findOneAndReplaceEvent - 查代事件
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndReplaceOptions 選項
 * @return *mongo.SingleResult returnSingleResultPointer 更添結果
 */
func (mongoDB *MongoDB) findOneAndReplaceEvent(
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
			Collection(mongoDB.GetConfigValueOrPanic(`event-table`))

		eventRWMutex.Lock() // 寫鎖

		collection.
			Indexes().
			CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys: eventSortPrimitiveM,
				},
			)

		// 更新事件
		singleResultPointer := collection.
			FindOneAndReplace(
				context.TODO(),
				filter,
				replacement,
				opts...,
			)

		eventRWMutex.Unlock() // 寫解鎖

		findOneAndReplaceError := singleResultPointer.Err() // 錯誤

		if mongo.ErrNoDocuments == findOneAndReplaceError {
			findOneAndReplaceError = nil
		}

		logings.SendLog(
			[]string{`%s %s 更添事件 %+v 選項 %+v 為 %+v `},
			append(network.GetAliasAddressPair(address), filter, opts, replacement),
			findOneAndReplaceError,
			0,
		)

		if nil != findOneAndReplaceError { // 若代添事件 錯誤
			return // 回傳
		}

		returnSingleResultPointer = singleResultPointer // 回傳結果指標

	}

	return // 回傳
}

// repsertOneEvent - 代添事件
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return []records.Event results 更新結果
 */
func (mongoDB *MongoDB) repsertOneEvent(filter, replacement primitive.M) (results []records.Event) {

	var replacedEvent records.Event // 更新的紀錄

	resultPointer :=
		mongoDB.
			findOneAndReplaceEvent(
				filter,
				replacement,
				options.FindOneAndReplace().SetUpsert(true),
			)

	if nil != resultPointer &&
		nil ==
			resultPointer.
				Decode(&replacedEvent) { // 若更新沒錯誤
		results = append(results, replacedEvent) // 回傳結果
	}

	return
}

// RepsertEvent - 代添事件
/**
 * @param records.Event event 事件
 * @return []records.Event results 回傳結果
 */
func (mongoDB *MongoDB) RepsertEvent(event records.Event) (
	results []records.Event) {

	// 代添紀錄
	results = mongoDB.repsertOneEvent(
		bson.M{
			`time`: event.Time,
		},
		event.PrimitiveM(),
	)

	return
}
