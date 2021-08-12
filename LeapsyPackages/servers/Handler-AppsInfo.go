package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/model"
	"leapsy.com/packages/network"
)

// 驗證並取得所有apps info
func postAllAppsInfoAPIHandler(eCAPIServer *ECAPIServer, ginContextPointer *gin.Context) {

	// 客戶端參數格式
	type Parameters struct {

		//帳戶資訊
		UserID       string `form:"userID" json:"userID" binding:"required"`
		UserPassword string `form:"userPassword" json:"userPassword" binding:"required"`
		DeviceID     string `form:"deviceID" json:"deviceID" binding:"required"`
		DeviceBrand  string `form:"deviceBrand" json:"deviceBrand" binding:"required"`

		// ProjectName string `form:"projectName" json:"projectName" binding:"required"`
		// AppName string `form:"appName" json:"appName" binding:"required"`
	}

	// 接收客戶端之參數
	var parameters Parameters

	// 轉譯json參數
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

	// log
	logings.SendLog(
		[]string{`%s %s 接受 %s 請求 %s %+v `},
		defaultArgs,
		nil,
		0,
	)

	// 取得各參數值
	parametersUserID := parameters.UserID
	parametersUserPassword := parameters.UserPassword
	parametersDeviceID := parameters.DeviceID
	parametersDeviceBrand := parameters.DeviceBrand

	fmt.Println("測試：已取得參數 parametersUserID=", parametersUserID, ",parametersUserPassword=", parametersUserPassword, ",parametersDeviceID=", parametersDeviceID, ",parametersDeviceBrand=", parametersDeviceBrand)

	// 若順利取出 則進行密碼驗證
	if bindJSONError == nil && bindURIError == nil {

		fmt.Println("取參數正確")

		//checkPassword(parametersUserID,parametersUserPassword)

		// 密碼正確
		if checkPassword(parametersUserID, parametersUserPassword) {

			// 查資料庫
			result := mongoDB.FindAllAppsInfoByProjectNameAndAppName()

			// 包成前端格式
			myResult := model.AppsInfoResponse{
				IsSuccess: true,
				Results:   "",
				Data:      result,
			}
			// myResult := model.AppsInfoResponse{
			// 	Code:    "200",
			// 	Message: "",
			// 	Data:    result,
			// }

			// 回應給前端
			ginContextPointer.JSON(http.StatusOK, myResult)

			logings.SendLog(
				[]string{`%s %s 回應 %s 請求 %s %+v: 密碼正確-查詢結果為 %+v`},
				append(
					defaultArgs,
					result,
					//[]interface{}{`a`},
				),
				nil,
				0,
			)
		} else {
			//密碼錯誤
			fmt.Println("取參數錯誤,錯誤訊息:bindJSONError=", bindJSONError, ",bindURIError=", bindURIError)

			// 包成回給前端的格式
			myResult := model.AppsInfoResponse{
				IsSuccess: false,
				Results:   "驗證失敗",
				Data:      nil,
			}

			// myResult := model.APIResponse{
			// 	IsSuccess: false,
			// 	Results:   "驗證失敗",
			// 	Data:      nil,
			// }

			// 回應給前端
			ginContextPointer.JSON(http.StatusNotFound, myResult)

			// log

			logings.SendLog(
				[]string{`%s %s 回應 %s 請求 %s %+v: 驗證失敗-帳號或密碼錯誤 `},
				append(
					defaultArgs,
				),
				nil,              // 無錯誤
				logrus.InfoLevel, // info等級的log
			)

		}

	} else if bindJSONError != nil {

		fmt.Println("取參數錯誤,錯誤訊息:bindJSONError=", bindJSONError, ",bindURIError=", bindURIError)

		// 包成回給前端的格式
		myResult := model.AppsInfoResponse{
			IsSuccess: false,
			Results:   "驗證失敗",
			Data:      nil,
		}

		// myResult := model.APIResponse{
		// 	IsSuccess: false,
		// 	Results:   "驗證失敗",
		// 	Data:      nil,
		// }

		// 回應給前端
		ginContextPointer.JSON(http.StatusNotFound, myResult)

		// log

		logings.SendLog(
			[]string{`%s %s 回應 %s 請求 %s %+v: 驗證失敗-取參數錯誤(參數有少或格式錯誤), bindJSONError=%s, bindURIError=%s`},
			append(
				defaultArgs,
				bindJSONError,
				bindURIError,
			),
			nil,              // 無錯誤
			logrus.InfoLevel, // info等級的log
		)
	}

}

func checkPassword(userID string, userPassword string) (result bool) {

	// 愈設帳號
	if userID == "default" && userPassword == "default" {
		return true
	} else {
		return false
	}

}
