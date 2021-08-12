package records

import (
	"reflect"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"leapsy.com/packages/logings"
	"leapsy.com/times"
)

// AlertRecord - 警報紀錄
type AlertRecord struct {
	AlertEventID, // 警報編號
	AlertType int // 警報群組
	AlertEventTime time.Time // 日期時間
	VarTag,        // 點名稱
	Comment, // 說明
	LineText string // 行文字
	IsRead, // 是否已讀
	IsHidden bool // 是否隱藏
}

var (
	alertRecordToECSAlertRecordMap = map[string]string{ // alertRecordToECSAlertRecordMap - 警報紀錄與環控資料庫警報紀錄欄位對照
		`AlertEventID`:   `ALERTEVENTID`,   // 警報編號 int
		`AlertEventTime`: `ALERTEVENTTIME`, // 日期時間	datetime
		`VarTag`:         `VARTAG`,         // 點名稱	nvarchar(50)
		`Comment`:        `COMMENT`,        // 說明	nvarchar(max)
		`AlertType`:      `ALERTTYPE`,      // 警報群組	int
		`LineText`:       `LINETEXT`,       // 行文字	nvarchar(max)
	}
)

// getMappedToECSAlertRecordFieldName - 取得警報紀錄對應的環控警報紀錄欄位名
/**
 * @param  string alertRecordFieldName 警報紀錄欄位名
 * @return string 環控警報紀錄欄位名
 */
func getMappedToECSAlertRecordFieldName(alertRecordFieldName string) string {
	return alertRecordToECSAlertRecordMap[alertRecordFieldName] // 回傳警報紀錄對應的環控警報紀錄欄位名
}

// AlertRecord - 將ECSAlertRecord轉成AlertRecord
/**
 * @return AlertRecord 警報紀錄
 */
func (ecsAlertRecord ECSAlertRecord) AlertRecord() (alertRecord AlertRecord) {

	valueOfECSAlertRecord := reflect.ValueOf(ecsAlertRecord)   // 環控警報紀錄的值
	typeOfAlertRecord := reflect.TypeOf(alertRecord)           // 警報紀錄的資料型別
	valueOfAlertRecord := reflect.ValueOf(&alertRecord).Elem() // 警報紀錄的值

	for index := 0; index < typeOfAlertRecord.NumField(); index++ { // 針對警報紀錄每一個欄位

		alertRecordFieldName := typeOfAlertRecord.Field(index).Name                         // 警報紀錄欄位名
		alertRecordFieldValue := valueOfAlertRecord.Field(index)                            // 警報紀錄欄位值
		ecsAlertRecordFieldName := getMappedToECSAlertRecordFieldName(alertRecordFieldName) // 環控警報紀錄欄位名

		if `` != ecsAlertRecordFieldName { // 若有對應的環控警報紀錄欄位名

			ecsAlertRecordFieldValue := valueOfECSAlertRecord.FieldByName(ecsAlertRecordFieldName) // 環控警報紀錄欄位值

			switch typeOfAlertRecord.Field(index).Type.String() { // 若警報紀錄欄位型別為

			case `int`: // 整數

				integer, strconvAtoiError := strconv.Atoi(ecsAlertRecordFieldValue.String()) // 將環控警報紀錄欄位值字串轉為整數

				logings.SendLog(
					[]string{`環控警報紀錄欄位 %s 值轉成整數`},
					[]interface{}{ecsAlertRecordFieldName},
					strconvAtoiError,
					logrus.ErrorLevel,
				)

				alertRecordFieldValue.SetInt(int64(integer)) // 設定警報紀錄欄位值為環控警報紀錄欄位轉化後的整數值

			case `time.Time`: // 時間
				alertRecordFieldValue.Set(reflect.ValueOf(times.ALERTEVENTTIMEStringToTime(ecsAlertRecordFieldValue.String()))) // 設定警報紀錄欄位值為環控警報紀錄欄位的時間值

			case `bool`: // 布林值
				alertRecordFieldValue.SetBool(ecsAlertRecordFieldValue.Bool()) // 設定警報紀錄欄位值為環控警報紀錄欄位的布林值

			default: // 預設
				alertRecordFieldValue.SetString(ecsAlertRecordFieldValue.String()) // 設定警報紀錄欄位值為環控警報紀錄欄位的字串值

			}

		}

	}

	return // 回傳
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (alertRecord AlertRecord) PrimitiveM() (returnPrimitiveM primitive.M) {

	returnPrimitiveM = bson.M{
		`alerteventid`:   alertRecord.AlertEventID,   // 警報編號
		`alerttype`:      alertRecord.AlertType,      // 警報群組
		`alerteventtime`: alertRecord.AlertEventTime, // 日期時間
		`vartag`:         alertRecord.VarTag,         // 點名稱
		`comment`:        alertRecord.Comment,        // 說明
		`linetext`:       alertRecord.LineText,       // 行文字
		`isread`:         alertRecord.IsRead,         // 是否已讀
		`ishidden`:       alertRecord.IsHidden,       // 是否隱藏
	}

	return
}
