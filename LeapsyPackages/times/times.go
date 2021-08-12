package times

import (
	"errors"
	"time"

	"leapsy.com/packages/logings"
)

// IsSecond - 判斷時間是否為整秒
/**
 * @param  time.Time dateTime  時間
 * @return bool 判斷是否為整秒
 */
func IsSecond(dateTime time.Time) bool {
	return dateTime.Second() == 0 && dateTime.Nanosecond() == 0 // 分秒微秒為零
}

// GetHourlyBounds - 取得時間上下限小時
/**
 * @param  time.Time dateTime  時間
 * @return time.Time lower 下限小時
 * @return time.Time upper 上限小時
 */
func GetHourlyBounds(dateTime time.Time) (low time.Time, upper time.Time) {

	duration := time.Hour

	if IsHour(dateTime) {
		low = dateTime.Add(-duration)
		upper = dateTime
	} else {
		low = ConvertToHourlyDateTime(dateTime)
		upper = low.Add(duration)
	}

	low = time.Date(low.Year(), low.Month(), low.Day(), low.Hour(), 0, 0, 0, time.Local)
	upper = time.Date(upper.Year(), upper.Month(), upper.Day(), upper.Hour(), 0, 0, 0, time.Local)

	logings.SendLog(
		[]string{`取得時間 %+v 下限小時 %+v 上限小時 %+v`},
		[]interface{}{dateTime, low, upper},
		nil,
		0,
	)

	return
}

// IsHour - 判斷時間是否為整點
/**
 * @param  time.Time dateTime  時間
 * @return bool 判斷是否為整點
 */
func IsHour(dateTime time.Time) bool {
	return dateTime.Minute() == 0 && IsSecond(dateTime) // 分秒微秒為零
}

// ConvertToHourlyDateTime - 轉成整點時間
/**
 * @param time.Time dateTime 時間
 * @return time.Time returnHourlyDateTime 回傳小時時間
 */
func ConvertToHourlyDateTime(dateTime time.Time) (returnHourlyDateTime time.Time) {

	formatSlice := `將時間 %+v 轉成整點時間`
	defaultArgs := []interface{}{dateTime}

	if !dateTime.IsZero() {

		if !IsHour(dateTime) { // 若非整點時間
			// 修改時間
			returnHourlyDateTime = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), dateTime.Hour(), 0, 0, 0, time.Local)
		} else {
			returnHourlyDateTime = dateTime // 回傳
		}

		logings.SendLog(
			[]string{formatSlice + ` %+v `},
			append(defaultArgs, returnHourlyDateTime),
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

	return
}

// GetDailyBounds - 取得時間上下限日
/**
 * @param  time.Time dateTime  時間
 * @return  time.Time lower  下限日
 * @return  time.Time upper  上限日
 */
func GetDailyBounds(dateTime time.Time) (low time.Time, upper time.Time) {

	if IsDay(dateTime) { // 若為整日
		low = dateTime.AddDate(0, 0, -1) // 昨日
		upper = dateTime                 // 今日
	} else {
		low = ConvertToDailyDateTime(dateTime) // 今日
		upper = low.AddDate(0, 0, 1)           // 明日
	}

	logings.SendLog(
		[]string{`取得時間 %+v 下限日 %+v 上限日 %+v`},
		[]interface{}{dateTime, low, upper},
		nil,
		0,
	)

	return // 回傳
}

// IsDay - 判斷是否為整日
/**
 * @param  time.Time dateTime  時間
 * @return  bool 判斷是否為整日
 */
func IsDay(dateTime time.Time) bool {
	return dateTime.Hour() == 0 && IsHour(dateTime) // 時分秒微秒為零
}

// ConvertToDailyDateTime - 轉成整日時間
/**
 * @param time.Time dateTime 時間
 * @return time.Time returnDailyDateTime 回傳小時時間
 */
func ConvertToDailyDateTime(dateTime time.Time) (returnDailyDateTime time.Time) {
	formatSlice := `將時間 %+v 轉成整日時間`
	defaultArgs := []interface{}{dateTime}

	if !dateTime.IsZero() {

		if !IsDay(dateTime) { // 若非整日時間
			// 修改時間
			returnDailyDateTime = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.Local)
		} else {
			returnDailyDateTime = dateTime // 回傳
		}

		logings.SendLog(
			[]string{formatSlice + ` %+v `},
			append(defaultArgs, returnDailyDateTime),
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

	return
}

// GetMonthlyBounds - 取得時間上下限月
/**
 * @param  time.Time dateTime  時間
 * @return  time.Time lower  下限月
 * @return  time.Time upper  上限月
 */
func GetMonthlyBounds(dateTime time.Time) (low time.Time, upper time.Time) {

	if IsMonth(dateTime) { // 若為整月
		low = dateTime.AddDate(0, -1, 0) // 上個月
		upper = dateTime                 // 這個月
	} else {
		low = ConvertToMonthlyDateTime(dateTime) // 這個月
		upper = low.AddDate(0, 1, 0)             // 下個月
	}

	logings.SendLog(
		[]string{`取得時間 %+v 下限月 %+v 上限月 %+v`},
		[]interface{}{dateTime, low, upper},
		nil,
		0,
	)

	return // 回傳
}

// IsMonth - 判斷是否為整月
/**
 * @param  time.Time dateTime  時間
 * @return  bool 判斷是否為整月
 */
func IsMonth(dateTime time.Time) bool {
	// 日為一、時分秒微秒為零
	return dateTime.Day() == 1 && IsDay(dateTime)
}

// ConvertToMonthlyDateTime - 轉成整月時間
/**
 * @param time.Time dateTime 時間
 * @return time.Time returnMonthlyDateTime 回傳小時時間
 */
func ConvertToMonthlyDateTime(dateTime time.Time) (returnMonthlyDateTime time.Time) {
	formatSlice := `將時間 %+v 轉成整月時間`
	defaultArgs := []interface{}{dateTime}

	if !dateTime.IsZero() {

		if !IsMonth(dateTime) { // 若非整月時間
			// 修改時間
			returnMonthlyDateTime = time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, time.Local)
		} else {
			returnMonthlyDateTime = dateTime // 回傳
		}

		logings.SendLog(
			[]string{formatSlice + ` %+v `},
			append(defaultArgs, returnMonthlyDateTime),
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

	return
}
