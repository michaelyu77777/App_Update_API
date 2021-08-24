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
func postAllAppsInfoAPIHandler(apiServer *APIServer, ginContextPointer *gin.Context) {

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
					apiServer.GetConfigValueOrPanic(`host`),
					apiServer.GetConfigPositiveIntValueOrPanic(`port`),
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
			result := mongoDB.FindAllAppsInfo()

			fmt.Printf("找到appsInfo結果 %d個", len(result))
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
			fmt.Println("密碼錯誤")

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

	// 搜尋帳號比對密碼
	allAccount := mongoDB.FindAllAccounts()

	fmt.Printf("找到 %d 個 %+v\n", len(allAccount), allAccount)

	account := mongoDB.FindAccountByUserID(userID)

	fmt.Printf("找到 %d 個\n", len(account))

	//若有找到結果
	if 1 > len(account) {
		fmt.Printf("找不到帳號\n")
		return false
	} else {

		// 密碼正確
		if userPassword == account[0].UserPassword {
			fmt.Printf("驗證密碼正確 %+v \n", account)
			return true
		} else {
			fmt.Printf("驗證密碼錯誤 %+v \n", account)
			return false
		}
	}

}
