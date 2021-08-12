package servers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
	"leapsy.com/times"
)

// getRecordsMonthlyAPIHandler - 處理GET一月內日資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func getRecordsMonthlyAPIHandler(eCAPIServer *ECAPIServer, ginContextPointer *gin.Context) {

	type Parameters struct {
		Year  int `uri:"year" json:"year"`
		Month int `uri:"month" json:"month"`
		Day   int `uri:"day" json:"day"`
	}

	var parameters Parameters

	bindJSONError := ginContextPointer.ShouldBindJSON(&parameters)

	bindURIError := ginContextPointer.ShouldBindUri(&parameters)

	defaultArgs :=
		append(
			network.GetAliasAddressPair(
				fmt.Sprintf(`%s:%d`,
					eCAPIServer.GetConfigValueOrPanic(`host`),
					eCAPIServer.GetConfigPositiveIntValueOrPanic(`port`),
				),
			),
			ginContextPointer.ClientIP(),
			ginContextPointer.FullPath(),
			parameters,
		)

	logings.SendLog(
		[]string{`%s %s 接受 %s 請求 %s %+v `},
		defaultArgs,
		nil,
		0,
	)

	parametersYear := parameters.Year
	parametersMonth := parameters.Month

	if bindJSONError == nil &&
		bindURIError == nil &&
		parametersMonth >= 1 && parametersMonth <= 12 {

		low, upper :=
			times.
				GetMonthlyBounds(
					time.Date(
						parametersYear,
						time.Month(parametersMonth),
						1,
						0,
						0,
						0,
						0,
						time.Local,
					).
						AddDate(
							0,
							1,
							0,
						),
				) // 取得上下限月

		var (
			endDateTime time.Time // 結束時間
		)

		if now := time.Now(); upper.After(now) { // 若上限時間在現在時間之後
			endDateTime = now
		} else {
			endDateTime = upper
		}

		if daysCount := int(times.ConvertToDailyDateTime(endDateTime).Sub(low).Hours()) / 24; daysCount != mongoDB.CountDailyRecordsBetweenTimes(low, false, upper, true) { // 若缺資料

			dailyRecordsDoneChannel := make(chan bool, 1)

			for dateTime := low.AddDate(0, 0, 1); !dateTime.After(endDateTime); dateTime = dateTime.AddDate(0, 0, 1) { // 針對每一小時

				thisDateTime := dateTime

				go func() {

					if mongoDB.CountDailyRecordByTime(thisDateTime) == 0 {
						mongoDB.AggregateRepsertDailyRecordByTime(thisDateTime) // 代添紀錄到現在時間

						if mongoDB.CountDailyRecordByTime(thisDateTime) == 0 {
							mongoDB.RepsertDailyRecord(records.DailyRecord{Time: thisDateTime})
						}

					}

					dailyRecordsDoneChannel <- true

				}()

			}

			for dateTime := low.AddDate(0, 0, 1); !dateTime.After(endDateTime); dateTime = dateTime.AddDate(0, 0, 1) { // 針對每一小時
				<-dailyRecordsDoneChannel
			}

		}

		result := mongoDB.FindDailyRecordsBetweenTimes(low, false, upper, true)

		ginContextPointer.JSON(http.StatusOK, result)

		logings.SendLog(
			[]string{`%s %s 回應 %s 請求 %s %+v : %+v `},
			append(
				defaultArgs,
				result,
			),
			nil,
			0,
		)

	} else {
		ginContextPointer.Status(http.StatusNotFound)

		logings.SendLog(
			[]string{`%s %s 拒絕回應 %s 請求 %s %+v `},
			defaultArgs,
			nil,
			0,
		)

	}

}
