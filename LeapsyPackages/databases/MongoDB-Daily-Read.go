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

// countDailyRecords - 計算日紀錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 日紀錄個數
 */
func (mongoDB *MongoDB) countDailyRecords(filter primitive.M) (returnCount int) {

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

		dailyRWMutex.RLock() // 讀鎖

		// 取得日紀錄個數
		count, countError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
			CountDocuments(context.TODO(), filter)

		dailyRWMutex.RUnlock() // 讀解鎖

		if nil != countError && mongo.ErrNilDocument != countError { // 若取得日紀錄個數錯誤，且不為空資料表錯誤

			logings.SendLog(
				[]string{`%s %s 取得日紀錄個數 %d `},
				append(defaultArgs, count),
				countError,
				logrus.ErrorLevel,
			)

			return // 回傳

		}

		logings.SendLog(
			[]string{`%s %s 取得日紀錄個數 %d `},
			append(defaultArgs, count),
			countError,
			0,
		)

		returnCount = int(count)

	}

	return // 回傳
}

// CountDailyRecordByTime - 依據時間計算日記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountDailyRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.countDailyRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// CountDailyRecordsBetweenTimes - 依據時間區間計算日紀錄數
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountDailyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (result int) {

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

		// 回傳結果
		result = mongoDB.countDailyRecords(
			bson.M{
				`time`: bson.M{
					greaterThanKeyword: lowerTime,
					lessThanKeyword:    upperTime,
				},
			},
		)

	}

	return // 回傳
}

// findDailyRecords - 取得日紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) findDailyRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.DailyRecord) {

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

		dailyRWMutex.RLock() // 讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		dailyRWMutex.RUnlock() // 讀解鎖

		logings.SendLog(
			[]string{`%s %s 查找日記錄 %+v `},
			append(defaultArgs, filter),
			findError,
			logrus.ErrorLevel,
		)

		if nil != findError { // 若查找日紀錄錯誤
			return // 回傳
		} // 若查找日記錄成功

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var dailyRecord records.DailyRecord

			cursorDecodeError := cursor.Decode(&dailyRecord) // 解析紀錄

			logings.SendLog(
				[]string{`%s %s 取得日記錄 %+v `},
				append(defaultArgs, dailyRecord),
				cursorDecodeError,
				logrus.ErrorLevel,
			)

			if nil != cursorDecodeError { // 若解析記錄錯誤
				return // 回傳
			}

			dailyRecord.Time = dailyRecord.Time.Local() // 儲存為本地時間格式

			results = append(results, dailyRecord) // 儲存紀錄
		}

		cursorErrError := cursor.Err() // 游標錯誤

		// 取得記錄器格式和參數

		logings.SendLog(
			[]string{`%s %s 查找日記錄遊標運作`},
			defaultArgs,
			cursorErrError,
			logrus.ErrorLevel,
		)

		if nil != cursorErrError { // 若遊標錯誤
			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 取得日記錄 %+v `},
			append(defaultArgs, results),
			nil,
			0,
		)

	}

	return // 回傳
}

// FindDailyRecordsBetweenTimes - 依據時間區間取得日紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindDailyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.DailyRecord) {

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

		// 回傳結果
		results = mongoDB.findDailyRecords(
			bson.M{
				`time`: bson.M{
					greaterThanKeyword: lowerTime,
					lessThanKeyword:    upperTime,
				},
			},
			options.
				Find().
				SetSort(
					bson.M{
						`time`: 1,
					},
				).
				SetBatchSize(int32(batchSize)),
		)

	}

	return // 回傳
}

// findOneDailyRecord - 取得一筆日紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return records.DailyRecord result 取得結果
 */
func (mongoDB *MongoDB) findOneDailyRecord(filter primitive.M, opts ...*options.FindOneOptions) (result records.DailyRecord) {

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
			Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
			FindOne(
				context.TODO(),
				filter,
				opts...,
			).
			Decode(&result)

		periodicallyRWMutex.RUnlock() // 讀解鎖

		var dailyRecord records.DailyRecord

		logings.SendLog(
			[]string{`取得日記錄 %+v`},
			append(defaultArgs, dailyRecord),
			decodeError,
			logrus.ErrorLevel,
		)

		if nil != decodeError { // 若解析記錄錯誤
			return // 回傳
		}

		result.Time = result.Time.Local() // 儲存為本地時間格式

		logings.SendLog(
			[]string{`%s %s 取得日資料 %+v `},
			append(defaultArgs, result),
			nil,
			0,
		)

	}

	return // 回傳
}

// FindDailyRecordEdgesOfTimes - 依據時間區間兩端取得日紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindDailyRecordEdgesOfTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.DailyRecord) {

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
			mongoDB.findOneDailyRecord(
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
			mongoDB.findOneDailyRecord(
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
