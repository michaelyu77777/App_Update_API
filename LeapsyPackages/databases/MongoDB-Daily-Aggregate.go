package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/records"
	"leapsy.com/times"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// aggregateDailyRecordsFromSecondRecords - 從秒紀錄聚集日紀錄
/**
 * @param []primitive.D pipeline 管線
 * @param ...*options.AggregateOptions opts 選項
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) aggregateDailyRecordsFromSecondRecords(pipeline []primitive.D, opts ...*options.AggregateOptions) (results []records.DailyRecord) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		database := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`))

		periodicallyRWMutex.Lock() // 寫鎖

		database.
			Collection(mongoDB.GetConfigValueOrPanic(`daily-table`)).
			Indexes().
			CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys:    dailyRecordSortPrimitiveD,
					Options: options.Index().SetUnique(true),
				},
			)

		periodicallyRWMutex.Unlock() // 寫解鎖

		periodicallyRWMutex.RLock() // 讀鎖

		// 聚集紀錄
		cursor, aggregateError := database.
			Collection(mongoDB.GetConfigValueOrPanic(`second-table`)).
			Aggregate(
				context.TODO(),
				pipeline,
				opts...,
			)

		periodicallyRWMutex.RUnlock() // 讀解鎖

		logings.SendLog(
			[]string{`從秒紀錄聚集小時記錄`},
			[]interface{}{},
			aggregateError,
			logrus.ErrorLevel,
		)

		if nil != aggregateError { // 若聚集小時紀錄錯誤
			return // 回傳
		}

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var dailyRecord records.DailyRecord

			cursorDecodeError := cursor.Decode(&dailyRecord) // 解析紀錄

			// 取得記錄器格式和參數

			logings.SendLog(
				[]string{`取得日記錄 %+v `},
				[]interface{}{dailyRecord},
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

		logings.SendLog(
			[]string{`查找小時記錄遊標運作`},
			[]interface{}{},
			cursorErrError,
			logrus.ErrorLevel,
		)

	}

	return // 回傳
}

// AggregateRepsertDailyRecordByTime - 依據時間代添輸出日紀錄
/**
 * @param @param time.Time dateTime 時間
 * @return []records.DailyRecord results 取得結果
 */
func (mongoDB *MongoDB) AggregateRepsertDailyRecordByTime(dateTime time.Time) (results []records.DailyRecord) {

	if !dateTime.IsZero() { //若時間不為零時間

		const (
			lookupPMKwhThisMonthFieldName           = `pmwkhm`
			lookupPMKwhTodayForDailyRecordFieldName = `pmwkhd`
		)

		currentDateTime := times.ConvertToDailyDateTime(dateTime) // 時間

		// 回傳結果
		results = mongoDB.aggregateDailyRecordsFromSecondRecords(
			mongo.Pipeline{
				mongoDB.getLookupPMKwhThisMonthPrimitiveD(currentDateTime, lookupPMKwhThisMonthFieldName),
				{
					{
						unwindConstString,
						fmt.Sprintf(`$%s`, lookupPMKwhThisMonthFieldName),
					},
				},
				mongoDB.getLookupPMKwhTodayForDailyRecordPrimitiveD(currentDateTime, lookupPMKwhTodayForDailyRecordFieldName),
				{
					{
						unwindConstString,
						fmt.Sprintf(`$%s`, lookupPMKwhTodayForDailyRecordFieldName),
					},
				},
				{
					{
						groupConstString,
						bson.D{
							{
								`_id`,
								nil,
							},
							{
								`pmkwhthismonth`,
								bson.D{
									{
										firstConstString,
										fmt.Sprintf(`$%s.result`, lookupPMKwhThisMonthFieldName),
									},
								},
							},
							{
								`pmkwhtoday`,
								bson.D{
									{
										firstConstString,
										fmt.Sprintf(`$%s.result`, lookupPMKwhTodayForDailyRecordFieldName),
									},
								},
							},
						},
					},
				},
				{
					{
						setConstString,
						bson.D{
							{
								`time`,
								currentDateTime,
							},
						},
					},
				},
				{
					{
						unsetConstString,
						[]string{
							`_id`,
						},
					},
				},
				{
					{
						mergeConstString,
						bson.D{
							{
								intoConstString,
								mongoDB.GetConfigValueOrPanic(`daily-table`),
							},
							{
								onConstString,
								`time`,
							},
							{
								whenMatchedConstString,
								`replace`,
							},
							{
								whenNotMatchedConstString,
								`insert`,
							},
						},
					},
				},
			},
			options.
				Aggregate().
				SetBatchSize(int32(batchSize)).
				SetAllowDiskUse(true),
		)

	}

	return // 回傳
}
