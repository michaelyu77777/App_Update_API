package databases

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"leapsy.com/packages/configurations"
)

const (
	equalToConstString            = `$eq`  // =
	greaterThanConstString        = `$gt`  // >
	greaterThanEqualToConstString = `$gte` // >=
	lessThanConstString           = `$lt`  // <
	lessThanEqualToConstString    = `$lte` // <=
)

const (
	includedConstString = `包含`
	excludedConstString = `不含`
)

const (
	nonIntFieldConstString = `非整數型態欄位錯誤`
)

const (
	timeIsZeroConstString = `零時間錯誤`
)

const (
	lookupConstString         = `$lookup`
	fromConstString           = `from`
	pipelineConstString       = `pipeline`
	matchConstString          = `$match`
	sortConstString           = `$sort`
	groupConstString          = `$group`
	setConstString            = `$set`
	unsetConstString          = `$unset`
	asConstString             = `as`
	unwindConstString         = `$unwind`
	firstConstString          = `$first`
	lastConstString           = `$last`
	outConstString            = `$out`
	subtractConstString       = `$subtract`
	projectConstString        = `$project`
	mergeConstString          = `$merge`
	intoConstString           = `into`
	onConstString             = `on`
	whenMatchedConstString    = `whenMatched`
	whenNotMatchedConstString = `whenNotMatched`
)

var (
	alertRecordSortPrimitiveM = bson.M{
		`alerteventid`: 1,
	}
)

var (
	hourlyRecordSortPrimitiveD = bson.D{
		{`time`, 1},
	}
	dailyRecordSortPrimitiveD = bson.D{
		{`time`, 1},
	}
	alertRecordSortPrimitiveD = bson.D{
		{`alerteventid`, 1},
	}
)

var (
	batchSize                                                      = configurations.GetConfigPositiveIntValueOrPanic(`local`, `batch-size`) // 取得預設批次大小
	periodicallyRWMutex, hourlyRWMutex, dailyRWMutex, alertRWMutex sync.RWMutex                                                             // 讀寫鎖
)

// getProcessRecordsBetweenTimesLogFormatSlice - 取得處理兩時間內紀錄的日誌格式片段
/**
 * @param string processVerb 處理動詞
 * @param time.Time lowerTime 下限時間
 * @param bool isLowerTimeIncluded 是否包含下限時間
 * @param time.Time upperTime 上限時間
 * @param bool isUpperTimeIncluded 是否包含上限時間
 * @param string recordName 紀錄名
 * @return string logFormatSlice 日誌格式片段
 */
func getProcessRecordsBetweenTimesLogFormatSlice(
	processVerb string,
	lowerTime time.Time,
	isLowerTimeIncluded bool,
	upperTime time.Time,
	isUpperTimeIncluded bool,
	recordName string,
) (logFormatSlice string) {

	logFormatSlice = `%s %s ` + processVerb + ` 時間 %+v ( ` // 格式字串

	if isLowerTimeIncluded { // 若包含下限
		logFormatSlice += includedConstString
	} else {
		logFormatSlice += excludedConstString
	}

	logFormatSlice += ` ) 到 %+v ( `

	if isUpperTimeIncluded { // 若包含上限
		logFormatSlice += includedConstString
	} else {
		logFormatSlice += excludedConstString
	}

	logFormatSlice += ` ) 的 ` + recordName

	return // 回傳
}
