package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
)

// getAlertsAPIHandler - 處理GET警報資料網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func getAlertsPageNumberPageCountAPIHandler(eCAPIServer *ECAPIServer, ginContextPointer *gin.Context) {

	type Parameters struct {
		PageNumber int `uri:"pageNumber" json:"pageNumber"`
		PageCount  int `uri:"pageCount" json:"pageCount"`
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

	if bindJSONError == nil &&
		bindURIError == nil {

		result := mongoDB.FindAllAlertRecordsOfPage(parameters.PageNumber, parameters.PageCount)

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
