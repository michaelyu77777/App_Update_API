package servers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
)

// putAlertIsReadAPIHandler - 處理PUT設定警報資料已讀網頁
/**
 * @param  *ECAPIServer eCAPIServer 環控API伺服器指標
 * @param  *gin.Context ginContextPointer  gin Context 指標
 */
func putAlertAPIHandler(eCAPIServer *ECAPIServer, ginContextPointer *gin.Context) {

	type Parameters struct {
		AlertEventID int  `uri:"alertEventID" json:"alertEventID"`
		IsRead       bool `form:"isRead" json:"isRead"`
		IsHidden     bool `form:"isHidden" json:"isHidden"`
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

	parametersAlertEventID := parameters.AlertEventID

	if bindJSONError == nil &&
		bindURIError == nil {

		const (
			isReadConstString   = `isRead`   // 參數名
			isHiddenConstString = `isHidden` // 參數名
		)

		var (
			isReadMissing, // 是否已讀遺失
			isHiddenMissing bool // 是否被隱藏遺失
		)

		isToUpdate := false // 是否更新
		setPrimitiveM := bson.M{}

		isReadParameterValue := ginContextPointer.PostForm(isReadConstString) // 取得已讀
		isReadMissing = isReadParameterValue == ``                            // 是否已讀遺失

		isHiddenParameterValue := ginContextPointer.PostForm(isHiddenConstString) // 取得被隱藏
		isHiddenMissing = isHiddenParameterValue == ``                            // 是否被隱藏遺失

		defaultBoolFormatString := `網址參數 '%s' 值格式 '%+v' 或 '%+v'
`

		defaultBoolFormatArgs := []interface{}{true, false}

		if defaultFormatString := `網址參數 '%s' 遺失
`; isReadMissing && isHiddenMissing { // 若已讀遺失且已被隱藏遺失

			ginContextPointer.String(
				http.StatusOK,
				defaultFormatString,
				isReadConstString,
			)

			ginContextPointer.String(
				http.StatusOK,
				defaultFormatString,
				isHiddenConstString,
			)

		}

		if isRead, strconvParseBoolError := strconv.ParseBool(isReadParameterValue); !isReadMissing {

			if nil != strconvParseBoolError {

				ginContextPointer.String(
					http.StatusOK,
					defaultBoolFormatString,
					append([]interface{}{isReadConstString}, defaultBoolFormatArgs...)...,
				)

			} else {
				setPrimitiveM[`isread`] = isRead
				isToUpdate = true
			}

		}

		if isHidden, strconvParseBoolError := strconv.ParseBool(isHiddenParameterValue); !isHiddenMissing {

			if nil != strconvParseBoolError {

				ginContextPointer.String(
					http.StatusOK,
					defaultBoolFormatString,
					append([]interface{}{isHiddenConstString}, defaultBoolFormatArgs...)...,
				)

			} else {
				setPrimitiveM[`ishidden`] = isHidden
				isToUpdate = true
			}

		}

		result := []records.AlertRecord{}

		if isToUpdate { // 若要更新

			result =
				mongoDB.UpdateOneAlertRecord(
					bson.M{
						`alerteventid`: parametersAlertEventID,
					},
					bson.M{
						`$set`: setPrimitiveM,
					},
				)

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
