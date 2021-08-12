package databases

import (
	"time"

	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getLookupPMKwhThisHourPrimitiveD - 取得查找該小時PM累計電度primitive.D
/**
 * @param time.Time dateTime 時間
 * @param string outputFieldName 輸出欄位名稱
 * @return primitive.D result 計算結果
 */
func (mongoDB *MongoDB) getLookupPMKwhThisHourPrimitiveD(dateTime time.Time, outputFieldName string) (result primitive.D) {

	// formatSlice := `取得查找時間 %+v 該小時PM累計電度為欄位 '%s' 的primitive.D`
	// defaultArgs := []interface{}{dateTime, outputFieldName}

	// if !dateTime.IsZero() {

	// 	low, _ := times.GetHourlyBounds(dateTime) // 取得上下限小時

	// 	result = mongoDB.getLookupSecondRecordFieldsDifferenceByTimePrimitiveD(low, dateTime, `pmkwh`, outputFieldName)

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice + ` %+v `},
	// 		append(defaultArgs, result),
	// 		nil,
	// 	)

	// 	// logger.Infof(formatString, args...) // 記錄資訊

	// } else {
	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice},
	// 		defaultArgs,
	// 		errors.New(timeIsZeroConstString),
	// 	)

	// 	// logger.Infof(formatString, args...) // 記錄資訊
	// }

	return // 回傳
}

// getLookupPMKwhTodayForHourlyRecordPrimitiveD - 取得查找今日的PM累計電度primitive.D
/**
 * @param time.Time dateTime 時間
 * @param string outputFieldName 輸出欄位名稱
 * @return primitive.D result 計算結果
 */
func (mongoDB *MongoDB) getLookupPMKwhTodayForHourlyRecordPrimitiveD(dateTime time.Time, outputFieldName string) (result primitive.D) {

	// formatSlice := `取得查找時間 %+v 小時紀錄當日PM累計電度為欄位 '%s' 的primitive.D`
	// defaultArgs := []interface{}{dateTime, outputFieldName}

	// if !dateTime.IsZero() {

	// 	low, _ := times.GetDailyBounds(dateTime) // 取得上下限日

	// 	result = mongoDB.getLookupSecondRecordFieldsDifferenceByTimePrimitiveD(low, dateTime, `pmkwh`, outputFieldName)

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice + ` %+v `},
	// 		append(defaultArgs, result),
	// 		nil,
	// 	)

	// 	// logger.Infof(formatString, args...) // 記錄資訊

	// } else {
	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice},
	// 		defaultArgs,
	// 		errors.New(timeIsZeroConstString),
	// 	)

	// 	// logger.Infof(formatString, args...) // 記錄資訊
	// }

	return // 回傳
}

// RepsertHourlyRecord - 代添小時紀錄
/**
 * @param records.HourlyRecord hourlyRecord 小時紀錄
 */
func (mongoDB *MongoDB) RepsertHourlyRecord(hourlyRecord records.HourlyRecord) {

	// formatSlice := `代添小時紀錄 %+v `
	// defaultArgs := []interface{}{hourlyRecord}
	// hourlyRecordTime := hourlyRecord.Time

	// if !hourlyRecordTime.IsZero() {

	// 	currentDateTime := times.ConvertToHourlyDateTime(hourlyRecordTime) // 時間
	// 	hourlyRecord.Time = currentDateTime

	// 	// 代添紀錄
	// 	mongoDB.repsertOneHourlyRecord(bson.M{`time`: currentDateTime}, hourlyRecord.PrimitiveM())

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice},
	// 		defaultArgs,
	// 		nil,
	// 	)

	// 	logger.Infof(formatString, args...) // 記錄資訊

	// } else {

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice},
	// 		defaultArgs,
	// 		errors.New(timeIsZeroConstString),
	// 	)

	// 	logger.Infof(formatString, args...) // 記錄資訊
	// }
}
