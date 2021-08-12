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

// aggregateHourlyRecordsFromSecondRecords - 從秒紀錄聚集小時紀錄
/**
 * @param []primitive.D pipeline 管線
 * @param ...*options.AggregateOptions opts 選項
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) aggregateHourlyRecordsFromSecondRecords(pipeline []primitive.D, opts ...*options.AggregateOptions) (results []records.HourlyRecord) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		database := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`))

		periodicallyRWMutex.Lock() // 寫鎖

		database.
			Collection(mongoDB.GetConfigValueOrPanic(`hourly-table`)).
			Indexes().
			CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys:    hourlyRecordSortPrimitiveD,
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

		// 取得記錄器格式和參數

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

			var hourlyRecord records.HourlyRecord

			cursorDecodeError := cursor.Decode(&hourlyRecord) // 解析紀錄

			logings.SendLog(
				[]string{`取得小時記錄 %+v `},
				[]interface{}{hourlyRecord},
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
			[]string{`查找小時記錄遊標運作`},
			[]interface{}{},
			cursorErrError,
			logrus.ErrorLevel,
		)

	}

	return // 回傳
}

// AggregateRepsertHourlyRecordByTime - 依據時間聚集代添小時紀錄
/**
 * @param @param time.Time dateTime 時間
 * @return []records.HourlyRecord results 取得結果
 */
func (mongoDB *MongoDB) AggregateRepsertHourlyRecordByTime(dateTime time.Time) (results []records.HourlyRecord) {

	if !dateTime.IsZero() { //若時間不為零時間

		const (
			lookupPMKwhThisHourFieldName             = `pmwkhh`
			lookupPMKwhTodayForHourlyRecordFieldName = `pmwkhd`
		)

		currentDateTime := times.ConvertToHourlyDateTime(dateTime) // 時間

		// 回傳結果
		results = mongoDB.aggregateHourlyRecordsFromSecondRecords(
			mongo.Pipeline{
				mongoDB.getLookupPMKwhThisHourPrimitiveD(currentDateTime, lookupPMKwhThisHourFieldName),
				{
					{
						unwindConstString,
						fmt.Sprintf(`$%s`, lookupPMKwhThisHourFieldName),
					},
				},
				mongoDB.getLookupPMKwhTodayForHourlyRecordPrimitiveD(currentDateTime, lookupPMKwhTodayForHourlyRecordFieldName),
				{
					{
						unwindConstString,
						fmt.Sprintf(`$%s`, lookupPMKwhTodayForHourlyRecordFieldName),
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
								`pmkwhthishour`,
								bson.D{
									{
										firstConstString,
										fmt.Sprintf(`$%s.result`, lookupPMKwhThisHourFieldName),
									},
								},
							},
							{
								`pmkwhtoday`,
								bson.D{
									{
										firstConstString,
										fmt.Sprintf(`$%s.result`, lookupPMKwhTodayForHourlyRecordFieldName),
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
								mongoDB.GetConfigValueOrPanic(`hourly-table`),
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
