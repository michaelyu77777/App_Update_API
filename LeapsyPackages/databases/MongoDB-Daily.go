package databases

import (
	"errors"
	"time"

	"leapsy.com/packages/logings"
	"leapsy.com/records"
	"leapsy.com/times"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getLookupPMKwhThisMonthPrimitiveD - 取得查找當月PM累計電度primitive.D
/**
 * @param time.Time dateTime 時間
 * @param string outputFieldName 輸出欄位名稱
 * @return primitive.D result 計算結果
 */
func (mongoDB *MongoDB) getLookupPMKwhThisMonthPrimitiveD(dateTime time.Time, outputFieldName string) (result primitive.D) {

	formatSlice := `取得查找時間 %+v 當月PM累計電度為欄位 '%s' 的primitive.D`
	defaultArgs := []interface{}{dateTime, outputFieldName}

	if !dateTime.IsZero() {

		low, _ := times.GetMonthlyBounds(dateTime) // 取得上下限月

		result = mongoDB.getLookupSecondRecordFieldsDifferenceByTimePrimitiveD(low, dateTime, `pmkwh`, outputFieldName)

		logings.SendLog(
			[]string{formatSlice + ` %+v `},
			append(defaultArgs, result),
			nil,
			0,
		)

	} else {

		logings.SendLog(
			[]string{formatSlice},
			defaultArgs,
			errors.New(timeIsZeroConstString),
			0,
		)

	}

	return // 回傳
}

// getLookupPMKwhTodayForDailyRecordPrimitiveD - 取得查找當日PM累計電度primitive.D
/**
 * @param time.Time dateTime 時間
 * @param string outputFieldName 輸出欄位名稱
 * @return primitive.D result 計算結果
 */
func (mongoDB *MongoDB) getLookupPMKwhTodayForDailyRecordPrimitiveD(dateTime time.Time, outputFieldName string) (result primitive.D) {

	formatSlice := `取得查找時間 %+v 日紀錄當日PM累計電度為欄位 '%s' 的primitive.D`
	defaultArgs := []interface{}{dateTime, outputFieldName}

	if !dateTime.IsZero() {

		low, _ := times.GetDailyBounds(dateTime) // 取得上下限日

		result = mongoDB.getLookupSecondRecordFieldsDifferenceByTimePrimitiveD(low, dateTime, `pmkwh`, outputFieldName)

		logings.SendLog(
			[]string{formatSlice + ` %+v `},
			append(defaultArgs, result),
			nil,
			0,
		)

	} else {

		logings.SendLog(
			[]string{formatSlice},
			defaultArgs,
			errors.New(timeIsZeroConstString),
			0,
		)

	}

	return // 回傳
}

// RepsertDailyRecord - 代添日紀錄
/**
 * @param records.DailyRecord dailyRecord 日紀錄
 */
func (mongoDB *MongoDB) RepsertDailyRecord(dailyRecord records.DailyRecord) {

	formatSlice := `代添日紀錄 %+v `
	defaultArgs := []interface{}{dailyRecord}
	dailyRecordTime := dailyRecord.Time

	if !dailyRecordTime.IsZero() {

		currentDateTime := times.ConvertToDailyDateTime(dailyRecordTime) // 時間
		dailyRecord.Time = currentDateTime

		// 代添紀錄
		mongoDB.repsertOneDailyRecord(bson.M{`time`: currentDateTime}, dailyRecord.PrimitiveM())

		logings.SendLog(
			[]string{formatSlice},
			defaultArgs,
			nil,
			0,
		)

	} else {

		logings.SendLog(
			[]string{formatSlice},
			defaultArgs,
			errors.New(timeIsZeroConstString),
			0,
		)

	}
}
