package databases

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getLookupSecondRecordFieldsDifferenceByTimePrimitiveD - 取得查找秒紀錄欄位差primitive.D
/**
 * @param time.Time low 下限時間
 * @param time.Time upper 上限時間
 * @param string inputFieldName 輸入欄位名稱
 * @param string outputFieldName 輸出欄位名稱
 * @return primitive.D result 計算結果
 */
func (mongoDB *MongoDB) getLookupSecondRecordFieldsDifferenceByTimePrimitiveD(low, upper time.Time, inputFieldName, outputFieldName string) (result primitive.D) {

	// formatSlice := `取得 查找 時間 %+v 到 %+v 秒紀錄 欄位 '%s' 差 為欄位 '%s' 的primitive.D `
	// defaultArgs := []interface{}{low, upper, inputFieldName, outputFieldName}

	// if !low.IsZero() && !upper.IsZero() && `` != inputFieldName && `` != outputFieldName {

	// 	result = bson.D{
	// 		{
	// 			lookupConstString, bson.D{
	// 				{
	// 					fromConstString,
	// 					mongoDB.GetConfigValueOrPanic(`second-table`),
	// 				},
	// 				{
	// 					pipelineConstString,
	// 					[]primitive.D{
	// 						bson.D{
	// 							{
	// 								matchConstString, bson.D{
	// 									{
	// 										`time`,
	// 										bson.D{
	// 											{greaterThanEqualToConstString, low},
	// 											{lessThanEqualToConstString, upper},
	// 										},
	// 									},
	// 								},
	// 							},
	// 						},
	// 						bson.D{
	// 							{
	// 								sortConstString,
	// 								bson.D{
	// 									{
	// 										`time`,
	// 										-1,
	// 									},
	// 								},
	// 							},
	// 						},
	// 						bson.D{
	// 							{
	// 								groupConstString,
	// 								bson.D{
	// 									{
	// 										`_id`,
	// 										nil,
	// 									},
	// 									{
	// 										`_firstItem`,
	// 										bson.D{
	// 											{
	// 												firstConstString,
	// 												`$` + inputFieldName,
	// 											},
	// 										},
	// 									},
	// 									{
	// 										`_lastItem`,
	// 										bson.D{
	// 											{
	// 												lastConstString,
	// 												`$` + inputFieldName,
	// 											},
	// 										},
	// 									},
	// 								},
	// 							},
	// 						},
	// 						bson.D{
	// 							{
	// 								setConstString,
	// 								bson.D{
	// 									{
	// 										`result`,
	// 										bson.D{
	// 											{
	// 												subtractConstString,
	// 												[]string{
	// 													`$_firstItem`,
	// 													`$_lastItem`,
	// 												},
	// 											},
	// 										},
	// 									},
	// 								},
	// 							},
	// 						},
	// 						bson.D{
	// 							{
	// 								unsetConstString,
	// 								[]string{
	// 									`_id`,
	// 									`_firstItem`,
	// 									`_lastItem`,
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 				{
	// 					asConstString,
	// 					outputFieldName,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	// 取得記錄器格式和參數
	// 	formatString, args := logings.GetLogFuncFormatAndArguments(
	// 		[]string{formatSlice + ` %+v `},
	// 		append(defaultArgs, result),
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

	return // 回傳
}
