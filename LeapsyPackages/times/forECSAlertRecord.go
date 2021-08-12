package times

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"leapsy.com/packages/logings"
)

const (
	alerteventtimeConstPattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$` // 環控警報系統時間格式
)

// IsALERTEVENTTIMEString - 判斷是否為環控警報紀錄時間
/**
 * @param  string inputString 輸入字串
 * @return bool 判斷是否為環控警報紀錄時間
 */
func IsALERTEVENTTIMEString(inputString string) bool {

	result, _ := regexp.MatchString(alerteventtimeConstPattern, inputString) //比對字串

	return result //回傳結果
}

// ALERTEVENTTIMEStringToTime - 環控警報紀錄時間字串轉時間
/**
 * @param  string inputString 輸入字串
 * @return time.Time returnTime 時間
 */
func ALERTEVENTTIMEStringToTime(inputString string) (returnTime time.Time) {

	defaultArgs := []interface{}{inputString} // 預設參數

	if IsALERTEVENTTIMEString(inputString) { // 若輸入字串為環控紀錄時間字串

		inputSlices := regexp.MustCompile(`[:\-TZ]`).Split(inputString, -1) // 分割輸入字串

		year, _ := strconv.Atoi(inputSlices[0])   // 年
		month, _ := strconv.Atoi(inputSlices[1])  // 月
		day, _ := strconv.Atoi(inputSlices[2])    // 日
		hour, _ := strconv.Atoi(inputSlices[3])   // 時
		minute, _ := strconv.Atoi(inputSlices[4]) // 分
		second, _ := strconv.Atoi(inputSlices[5]) // 秒

		returnTime = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local) // 回傳時間

		logings.SendLog(
			[]string{`環控警報系統時間 %s 轉成時間 %+v `},
			append(defaultArgs, returnTime),
			nil,
			0,
		)

	} else {

		logings.SendLog(
			[]string{`%s 轉成時間`},
			defaultArgs,
			errors.New(`非環控警報系統時間格式 `+alerteventtimeConstPattern),
			0,
		)

	}

	return // 回傳
}
