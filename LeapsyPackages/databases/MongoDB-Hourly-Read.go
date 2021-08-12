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

// countHourlyRecords - 計算小時紀錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 小時紀錄個數
 */
func (mongoDB *MongoDB) countHourlyRecords(filter primitive.M) (returnCount int) {

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

		hourlyRWMutex.RLock() // 讀鎖

		// 取得小時紀錄個數
		count, countError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
			CountDocuments(context.TODO(), filter)

		hourlyRWMutex.RUnlock() // 讀解鎖

		if nil != countError && mongo.ErrNilDocument != countError { // 若取得小時紀錄個數錯誤，且不為空資料表錯誤

			logings.SendLog(
				[]string{`%s %s 取得小時紀錄個數 %d `},
				append(defaultArgs, count),
				countError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 取得小時紀錄個數 %d `},
			append(defaultArgs, count),
			countError,
			0,
		)

		returnCount = int(count)

	}

	return // 回傳
}

// CountHourlyRecordByTime - 依據時間計算小時記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountHourlyRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.countHourlyRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// CountHourlyRecordsBetweenTimes - 依據時間區間計算小時記錄數
/**
 * * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountHourlyRecordsBetweenTimes(
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
		result = mongoDB.countHourlyRecords(
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

// findHourlyRecords - 取得小時紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) findHourlyRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.HourlyRecord) {

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

		hourlyRWMutex.RLock() // 讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		hourlyRWMutex.RUnlock() // 讀解鎖

		logings.SendLog(
			[]string{`%s %s 查找小時記錄 %+v `},
			append(defaultArgs, filter),
			findError,
			logrus.ErrorLevel,
		)

		if nil != findError { // 若查找小時紀錄錯誤
			return // 回傳
		}

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var hourlyRecord records.HourlyRecord

			cursorDecodeError := cursor.Decode(&hourlyRecord) // 解析紀錄

			logings.SendLog(
				[]string{`%s %s 取得小時記錄 %+v `},
				append(defaultArgs, hourlyRecord),
				cursorDecodeError,
				logrus.ErrorLevel,
			)

			if nil != cursorDecodeError { // 若解析記錄錯誤
				return // 回傳
			}

			hourlyRecord.Time = hourlyRecord.Time.Local() // 儲存為本地時間格式

			results = append(results, hourlyRecord) // 儲存紀錄
		}

		cursorErrError := cursor.Err() // 游標錯誤

		logings.SendLog(
			[]string{`%s %s 查找小時記錄遊標運作`},
			defaultArgs,
			cursorErrError,
			logrus.ErrorLevel,
		)

		if nil != cursorErrError { // 若遊標錯誤
			return // 回傳
		}

		logings.SendLog(
			[]string{`%s %s 取得小時紀錄 %+v `},
			append(defaultArgs, results),
			nil,
			0,
		)

	}

	return // 回傳
}

// 查找[資料庫-軟體更新]appsinfo(軟體資訊)
// func (mongoDB *MongoDB) findAppsInfo(filter primitive.M, opts ...*options.FindOptions) (results []records.AppsInfo) {

// 	mongoClientPointer := mongoDB.Connect() // 資料庫指標

// 	if nil != mongoClientPointer { // 若資料庫指標不為空
// 		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

// 		// 預設主機
// 		address := fmt.Sprintf(
// 			`%s:%d`,
// 			mongoDB.GetConfigValueOrPanic(`server`),
// 			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
// 		)

// 		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

// 		hourlyRWMutex.RLock() // 讀鎖

// 		// 查找紀錄
// 		cursor, findError := mongoClientPointer.
// 			Database(mongoDB.GetConfigValueOrPanic(`database`)).
// 			Collection(mongoDB.GetConfigValueOrPanic(`appsInfo-table`)).
// 			Find(
// 				context.TODO(),
// 				filter,
// 				opts...,
// 			)

// 		hourlyRWMutex.RUnlock() // 讀解鎖

// 		// log 紀錄有查詢動作
// 		logings.SendLog(
// 			[]string{`%s %s 查找資料庫-軟體更新 %+v `},
// 			append(defaultArgs, filter),
// 			findError,
// 			logrus.ErrorLevel,
// 		)

// 		if nil != findError { // 若查找錯誤
// 			return // 回傳
// 		}

// 		defer cursor.Close(context.TODO()) // 記得關閉

// 		for cursor.Next(context.TODO()) { // 拜訪每筆查詢

// 			var appsInfo records.AppsInfo

// 			cursorDecodeError := cursor.Decode(&appsInfo) // 解析單筆結果，放到appsInfo變數中

// 			// log 針對查出的每筆紀錄寫log
// 			logings.SendLog(
// 				[]string{`%s %s 取得資料庫-軟體更新 %+v `},
// 				append(defaultArgs, appsInfo),
// 				cursorDecodeError,
// 				logrus.ErrorLevel,
// 			)

// 			if nil != cursorDecodeError { // 若解析記錄錯誤
// 				return // 回傳
// 			}

// 			// appsInfo.Time = appsInfo.Time.Local() // 儲存為本地時間格式

// 			results = append(results, appsInfo) // 將單筆查詢結果，加入到results結果中
// 		}

// 		cursorErrError := cursor.Err() // 結果游標錯誤

// 		// log 紀錄有查詢動作
// 		logings.SendLog(
// 			[]string{`%s %s 查找資料庫-軟體資訊 游標`},
// 			defaultArgs,
// 			cursorErrError,
// 			logrus.ErrorLevel,
// 		)

// 		if nil != cursorErrError { // 若遊標錯誤
// 			return // 回傳
// 		}

// 		// log 紀錄有查詢動作
// 		logings.SendLog(
// 			[]string{`%s %s 取得資料庫-軟體更新 %+v `},
// 			append(defaultArgs, results),
// 			nil,
// 			0,
// 		)

// 	}

// 	return // 回傳
// }

// FindHourlyRecordsBetweenTimes - 依據時間區間取得小時紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindHourlyRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.HourlyRecord) {

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
		results = mongoDB.findHourlyRecords(
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

// // 尋找所有 apps info
// func (mongoDB *MongoDB) FindAllAppsInfoByProjectNameAndAppName() (results []records.AppsInfo) {

// 	// if !lowerTime.IsZero() && !upperTime.IsZero() { //若上下限時間不為零時間

// 	// var (
// 	// 	greaterThanKeyword, lessThanKeyword string // 比較關鍵字
// 	// )

// 	// if !isLowerTimeIncluded { // 若不包含下限時間
// 	// 	greaterThanKeyword = greaterThanConstString // >
// 	// } else {
// 	// 	greaterThanKeyword = greaterThanEqualToConstString // >=
// 	// }

// 	// if !isUpperTimeIncluded { // 若不包含上限時間
// 	// 	lessThanKeyword = lessThanConstString // <
// 	// } else {
// 	// 	lessThanKeyword = lessThanEqualToConstString // <=
// 	// }

// 	// 回傳結果
// 	results = mongoDB.findAppsInfo(
// 		bson.M{},
// 	)

// 	// }

// 	return // 回傳
// }

// // 尋找符合的專案名稱,app名稱
// func (mongoDB *MongoDB) FindAppsInfoByProjectNameAndAppName(projectName string, appName string) (results []records.AppsInfo) {

// 	fmt.Println("測試中：取得參數projectName=", projectName, "appName=", appName)

// 	// 回傳結果
// 	results = mongoDB.findAppsInfo(
// 		bson.M{
// 			"projectName": projectName,
// 			"appName":     appName,
// 		},
// 	)

// 	// }

// 	return // 回傳
// }

// findOneHourlyRecord - 取得一筆小時紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return records.HourlyRecord result 取得結果
 */
func (mongoDB *MongoDB) findOneHourlyRecord(filter primitive.M, opts ...*options.FindOneOptions) (result records.HourlyRecord) {

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
			Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
			FindOne(
				context.TODO(),
				filter,
				opts...,
			).
			Decode(&result)

		periodicallyRWMutex.RUnlock() // 讀解鎖

		var hourlyRecord records.HourlyRecord

		logings.SendLog(
			[]string{`取得小時記錄 %+v`},
			append(defaultArgs, hourlyRecord),
			decodeError,
			logrus.ErrorLevel,
		)

		if nil != decodeError { // 若解析記錄錯誤
			return // 回傳
		}

		result.Time = result.Time.Local() // 儲存為本地時間格式

		logings.SendLog(
			[]string{`%s %s 取得小時資料 %+v `},
			append(defaultArgs, result),
			nil,
			0,
		)

	}

	return // 回傳
}

// FindHourlyRecordEdgesOfTimes - 依據時間區間兩端取得小時紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) FindHourlyRecordEdgesOfTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.HourlyRecord) {

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
			mongoDB.findOneHourlyRecord(
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
			mongoDB.findOneHourlyRecord(
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
