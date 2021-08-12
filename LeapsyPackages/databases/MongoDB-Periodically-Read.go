package databases

import (
	"time"

	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CountSecondRecords - 計算秒紀錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 秒紀錄個數
 */
func (mongoDB *MongoDB) countSecondRecords(filter primitive.M) (returnCount int) {

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

	// 	periodicallyRWMutex.RLock() // 讀鎖

	// 	// 取得秒紀錄個數
	// 	count, countError := mongoClientPointer.
	// 		Database(mongoDB.GetConfigValueOrPanic(`database`)).
	// 		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
	// 		CountDocuments(context.TODO(), filter)

	// 	periodicallyRWMutex.RUnlock() // 讀解鎖

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 取得秒紀錄個數 %+v `},
	// 		append(defaultArgs, count),
	// 		countError,
	// 	)

	// 	if nil != countError && mongo.ErrNilDocument != countError { // 若取得秒紀錄個數錯誤，且不為空資料表錯誤
	// 		logger.Errorf(formatString, args...) // 記錄錯誤
	// 		return                               // 回傳
	// 	}

	// 	go logger.Infof(formatString, args...) // 記錄資訊
	// 	returnCount = int(count)

	// }

	return // 回傳
}

// CountSecondRecordByTime - 依據時間計算秒記錄數
/**
 * @param time.Time dateTime 時間
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountSecondRecordByTime(dateTime time.Time) (result int) {

	if !dateTime.IsZero() {
		result = mongoDB.countSecondRecords(
			bson.M{
				`time`: bson.M{
					equalToConstString: dateTime,
				},
			},
		)
	}

	return // 回傳
}

// findSecondRecords - 取得秒紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return []records.SecondRecord results 取得結果
 */
func (mongoDB *MongoDB) findSecondRecords(filter primitive.M, opts ...*options.FindOptions) (results []records.SecondRecord) {

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

	// 	periodicallyRWMutex.RLock() // 讀鎖

	// 	// 查找紀錄
	// 	cursor, findError := mongoClientPointer.
	// 		Database(mongoDB.GetConfigValueOrPanic(`database`)).
	// 		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
	// 		Find(
	// 			context.TODO(),
	// 			filter,
	// 			opts...,
	// 		)

	// 	periodicallyRWMutex.RUnlock() // 讀解鎖

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 查找秒記錄 %+v `},
	// 		append(defaultArgs, filter),
	// 		findError,
	// 	)

	// 	if nil != findError { // 若查找秒紀錄錯誤
	// 		logger.Errorf(formatString, args...) // 記錄錯誤
	// 		return                               // 回傳
	// 	}

	// 	defer cursor.Close(context.TODO()) // 記得關閉

	// 	for cursor.Next(context.TODO()) { // 針對每一紀錄

	// 		var secondRecord records.SecondRecord

	// 		cursorDecodeError := cursor.Decode(&secondRecord) // 解析紀錄

	// 		// 取得記錄器格式和參數
	// 		formatString, args = logings.GetLogFuncFormatAndArguments(
	// 			[]string{`%s %s 取得秒記錄 %+v`},
	// 			append(defaultArgs, secondRecord),
	// 			cursorDecodeError,
	// 		)

	// 		if nil != cursorDecodeError { // 若解析記錄錯誤
	// 			logger.Errorf(formatString, args...) // 記錄錯誤
	// 			return                               // 回傳
	// 		}

	// 		go logger.Infof(formatString, args...) // 記錄資訊

	// 		secondRecord.Time = secondRecord.Time.Local() // 儲存為本地時間格式

	// 		results = append(results, secondRecord) // 儲存紀錄
	// 	}

	// 	cursorErrError := cursor.Err() // 游標錯誤

	// 	// 取得記錄器格式和參數
	// 	formatString, args = logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 查找秒記錄遊標運作`},
	// 		defaultArgs,
	// 		cursorErrError,
	// 	)

	// 	if nil != cursorErrError { // 若遊標錯誤
	// 		logger.Errorf(formatString, args...) // 記錄錯誤
	// 		return                               // 回傳
	// 	}

	// 	logger.Infof(formatString, args...) // 記錄資訊

	// 	// 取得記錄器格式和參數
	// 	formatString, args = logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 取得秒資料 %+v `},
	// 		append(defaultArgs, results),
	// 		nil,
	// 	)

	// 	logger.Infof(formatString, args...) // 記錄資訊

	// }

	return // 回傳
}

// FindSecondRecordsBetweenTimes - 依據時間區間取得秒紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.SecondRecord results 取得結果
 */
func (mongoDB *MongoDB) FindSecondRecordsBetweenTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.SecondRecord) {

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
		results = mongoDB.findSecondRecords(
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

// findOneSecondRecord - 取得一筆秒紀錄
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return records.SecondRecord result 取得結果
 */
func (mongoDB *MongoDB) findOneSecondRecord(filter primitive.M, opts ...*options.FindOneOptions) (result records.SecondRecord) {

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

	// 	periodicallyRWMutex.RLock() // 讀鎖

	// 	// 查找紀錄
	// 	decodeError := mongoClientPointer.
	// 		Database(mongoDB.GetConfigValueOrPanic(`database`)).
	// 		Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
	// 		FindOne(
	// 			context.TODO(),
	// 			filter,
	// 			opts...,
	// 		).
	// 		Decode(&result)

	// 	periodicallyRWMutex.RUnlock() // 讀解鎖

	// 	var secondRecord records.SecondRecord

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{`取得秒記錄 %+v`},
	// 		append(defaultArgs, secondRecord),
	// 		decodeError,
	// 	)

	// 	if nil != decodeError { // 若解析記錄錯誤
	// 		logger.Errorf(formatString, args...) // 記錄錯誤
	// 		return                               // 回傳
	// 	}

	// 	go logger.Infof(formatString, args...) // 記錄資訊

	// 	result.Time = result.Time.Local() // 儲存為本地時間格式

	// 	// 取得記錄器格式和參數
	// 	formatString, args = logings.GetLogFuncFormatAndArguments(
	// 		[]string{`%s %s 取得秒資料 %+v `},
	// 		append(defaultArgs, result),
	// 		nil,
	// 	)

	// 	logger.Infof(formatString, args...) // 記錄資訊

	// }

	return // 回傳
}

// FindSecondRecordEdgesOfTimes - 依據時間區間兩端取得秒紀錄
/**
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @return []records.SecondRecord results 取得結果
 */
func (mongoDB *MongoDB) FindSecondRecordEdgesOfTimes(
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
) (results []records.SecondRecord) {

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
			mongoDB.findOneSecondRecord(
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
			mongoDB.findOneSecondRecord(
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
