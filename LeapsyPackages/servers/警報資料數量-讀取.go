package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
)

// getAlertsCountAPIHandler - 處理GET警報數量資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func getAlertsPagesPageCountAPIHandler(eCAPIServer *ECAPIServer, ginContextPointer *gin.Context) {

	alertRecordsCountChannel := make(chan int, 1)

	go func() {

		result := mongoDB.CountAllAlertRecords()

		for {
			alertRecordsCountChannel <- result
		}

	}()

	type Parameters struct {
		PageCount int `uri:"pageCount" json:"pageCount"`
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

	parametersPageCount := parameters.PageCount

	if bindJSONError == nil && bindURIError == nil {

		result := parametersPageCount

		if parametersPageCount > 1 && <-alertRecordsCountChannel%parametersPageCount != 0 {
			result++
		}

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
