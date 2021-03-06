package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// countAlertRecords - 計算警報紀錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 警報紀錄個數
 */
func (mongoDB *MongoDB) countAlertRecords(filter primitive.M) (returnCount int) {

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

		alertRWMutex.RLock() // 讀鎖

		// 取得警報紀錄個數
		count, countError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
			CountDocuments(context.TODO(), filter)

		alertRWMutex.RUnlock() // 讀解鎖

		if nil != countError && mongo.ErrNilDocument != countError { // 若取得警報紀錄個數錯誤，且不為空資料表錯誤

			logings.SendLog(
				[]string{`%s %s 取得警報紀錄個數 %d `},
				append(defaultArgs, count),
				countError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 取得警報紀錄個數 %d `},
			append(defaultArgs, count),
			countError,
			0,
		)

		returnCount = int(count)

	}

	return // 回傳
}

// CountAllAlertRecords - 計算所有警報紀錄個數
/**
 * @retrun int returnCount 警報紀錄個數
 */
func (mongoDB *MongoDB) CountAllAlertRecords() (returnCount int) {

	returnCount = mongoDB.countAlertRecords(bson.M{})

	return // 回傳
}

// findAlertRecords - 查找警報紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.AlertRecord results 取得結果
 */
func (mongoDB *MongoDB) findAlertRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.AlertRecord) {

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

		alertRWMutex.RLock() // 讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		alertRWMutex.RUnlock() // 讀解鎖

		logings.SendLog(
			[]string{`%s %s 查找警報紀錄 %+v `},
			append(defaultArgs, filter),
			findError,
			logrus.ErrorLevel,
		)

		if nil != findError { // 若查找警報紀錄錯誤
			return // 回傳
		}

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var alertRecord records.AlertRecord

			cursorDecodeError := cursor.Decode(&alertRecord) // 解析紀錄

			logings.SendLog(
				[]string{`%s %s 取得警報記錄 %+v `},
				append(defaultArgs, alertRecord),
				cursorDecodeError,
				logrus.ErrorLevel,
			)

			if nil != cursorDecodeError { // 若解析記錄錯誤
				return // 回傳
			}

			alertRecord.AlertEventTime = alertRecord.AlertEventTime.Local() // 儲存為本地時間格式

			results = append(results, alertRecord) // 儲存紀錄
		}

		cursorErrError := cursor.Err() // 游標錯誤

		logings.SendLog(
			[]string{`%s %s 查找警報記錄遊標運作`},
			defaultArgs,
			cursorErrError,
			logrus.ErrorLevel,
		)

		if nil != cursorErrError { // 若遊標錯誤
			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 取得警報紀錄 %+v `},
			append(defaultArgs, results),
			nil,
			0,
		)

	}

	return // 回傳
}

// FindAllAlertRecords - 取得所有警報紀錄
/**
 * @return []records.AlertRecord results 取得結果
 */
func (mongoDB *MongoDB) FindAllAlertRecords() (results []records.AlertRecord) {

	// 取得警報紀錄
	results = mongoDB.findAlertRecords(bson.M{}, options.Find().SetSort(bson.M{`alerteventtime`: -1}).SetBatchSize(int32(batchSize)))

	return // 回傳
}

// FindAllAlertRecordsOfPage - 取得所有某頁警報紀錄
/**
 * @return []records.AlertRecord results 取得結果
 */
func (mongoDB *MongoDB) FindAllAlertRecordsOfPage(pageNumber, pageCount int) (results []records.AlertRecord) {

	// 取得警報紀錄
	results =
		mongoDB.
			findAlertRecords(
				bson.M{},
				options.
					Find().
					SetSort(
						bson.M{
							`alerteventtime`: -1,
						},
					).
					SetSkip(int64((pageNumber-1)*pageCount)).
					SetLimit(int64(pageCount)).
					SetBatchSize(int32(batchSize)),
			)

	return // 回傳
}

// FindAlertRecordsBetweenTimes - 依據時間區間取得警報紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.AlertRecord results 取得結果
 */
func (mongoDB *MongoDB) FindAlertRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.AlertRecord) {

	if !lowerTime.IsZero() && !upperTime.IsZero() { //若上下限時間不為零時間

		var (
			greaterThanKeyword, lessThanKeyword string // 比較關鍵字
		)

		if !isLowerTimeIncluded { // 若不包含下限時間
			greaterThanKeyword = greaterThanConstString // >
		} else {
			greaterThanKeyword = greaterThanEqualToConstString // >=
		}

		if !isUpperTimeIncluded { // 若不包含上限時間
			lessThanKeyword = lessThanConstString // <
		} else {
			lessThanKeyword = lessThanEqualToConstString // <=
		}

		results = mongoDB.findAlertRecords(
			bson.M{
				`alerteventtime`: bson.M{
					greaterThanKeyword: lowerTime,
					lessThanKeyword:    upperTime,
				},
			},
			options.
				Find().
				SetSort(
					bson.M{
						`alerteventtime`: -1,
					},
				).
				SetBatchSize(int32(batchSize)),
		)

	}

	return // 回傳
}

// findOneAlertRecord - 取得一筆警報紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return records.AlertRecord result 取得結果
 */
func (mongoDB *MongoDB) findOneAlertRecord(filter primitive.M, opts ...*options.FindOneOptions) (result records.AlertRecord) {

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

		periodicallyRWMutex.RLock() // 讀鎖

		// 查找紀錄
		decodeError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
			FindOne(
				context.TODO(),
				filter,
				opts...,
			).
			Decode(&result)

		periodicallyRWMutex.RUnlock() // 讀解鎖

		var alertRecord records.AlertRecord

		logings.SendLog(
			[]string{`取得警報記錄 %+v`},
			append(defaultArgs, alertRecord),
			decodeError,
			logrus.ErrorLevel,
		)

		if nil != decodeError { // 若解析記錄錯誤
			return // 回傳
		}

		result.AlertEventTime = result.AlertEventTime.Local() // 儲存為本地時間格式

		logings.SendLog(
			[]string{`%s %s 取得警報資料 %+v `},
			append(defaultArgs, result),
			nil,
			0,
		)

	}

	return // 回傳
}

// FindAlertRecordEdgesOfTimes - 依據時間區間兩端取得警報紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.AlertRecord results 取得結果
 */
func (mongoDB *MongoDB) FindAlertRecordEdgesOfTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.AlertRecord) {

	if !lowerTime.IsZero() && !upperTime.IsZero() { //若上下限時間不為零時間

		var (
			greaterThanKeyword, lessThanKeyword string // 比較關鍵字
		)

		if !isLowerTimeIncluded { // 若不包含下限時間
			greaterThanKeyword = greaterThanConstString // >
		} else {
			greaterThanKeyword = greaterThanEqualToConstString // >=
		}

		if !isUpperTimeIncluded { // 若不包含上限時間
			lessThanKeyword = lessThanConstString // <
		} else {
			lessThanKeyword = lessThanEqualToConstString // <=
		}

		defaultFilter := bson.M{
			`time`: bson.M{
				greaterThanKeyword: lowerTime,
				lessThanKeyword:    upperTime,
			},
		}

		defaultOption := options.
			FindOne().
			SetBatchSize(int32(batchSize))

		results = append(
			results,
			mongoDB.findOneAlertRecord(
				defaultFilter,
				defaultOption.
					SetSort(
						bson.M{
							`time`: 1,
						},
					),
			),
		)

		results = append(
			results,
			mongoDB.findOneAlertRecord(
				defaultFilter,
				defaultOption.
					SetSort(
						bson.M{
							`time`: -1,
						},
					),
			),
		)

	}

	return // 回傳
}
