package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	// "leapsy.com/records"
)

// 回應給客戶端有「大小」寫的json欄位
type AppsInfoWithDownloadPath struct {
	AppNameCht string `json:"appNameCht"` //軟體名稱 正體
	AppNameChs string `json:"appNameChs"` //軟體名稱 簡體
	AppNameEng string `json:"appNameEng"` //軟體名稱 英文
	AppNameJpn string `json:"appNameJpn"` //軟體名稱 日文
	AppNameKor string `json:"appNameKor"` //軟體名稱 韓文

	LabelName       string `json:"labelName"`       // APK Label名稱
	LastVersionCode int    `json:"lastVersionCode"` //最新版本號
	LastVersionName string `json:"lastVersionName"` //最新版本名
	PackageName     string `json:"packageName"`     //封包名稱
	PublishDate     string `json:"publishDate"`     //發佈日期

	ChangeDetailCht string `json:"changeDetailCht"` //更新內容 詳述 正
	ChangeDetailChs string `json:"changeDetailChs"` //更新內容 詳述 簡
	ChangeDetailEng string `json:"changeDetailEng"` //更新內容 詳述 英
	ChangeDetailJpn string `json:"changeDetailJpn"` //更新內容 詳述 日
	ChangeDetailKor string `json:"changeDetailKor"` //更新內容 詳述 韓

	ChangeBriefCht string `json:"changeBriefCht"` //更新內容 簡述 正
	ChangeBriefChs string `json:"changeBriefChs"` //更新內容 簡述 簡
	ChangeBriefEng string `json:"changeBriefEng"` //更新內容 簡述 英
	ChangeBriefJpn string `json:"changeBriefJpn"` //更新內容 簡述 日
	ChangeBriefKor string `json:"changeBriefKor"` //更新內容 簡述 韓

	//ApkFileName  string `json:"apkFileName"`  // APK檔名
	DownloadPath string `json:"downloadPath"` // 組合出下載APK網址
}

// 回應給客戶端有「大小」寫的json欄位
// type AppsInfoResponse struct {
// 	AppsInfoCommonStructResponse `bson:",inline"` //共用參數：會從DB拿、也會Response回Client的參數

// 	ApkFileName string `json:"apkFileName"` // APK檔案名稱
// 	// LabelName        string `json:"labelname"`        // APK Label名稱
// }

// 包成回給前端<取AppsInfo格式>
type AppsInfoResponse struct {
	IsSuccess bool   `json:"isSuccess"` //錯誤代碼
	Results   string `json:"results"`   //錯誤訊息
	// Data      []AppsInfo `json:"data"`      //查詢結果
	Data []AppsInfoWithDownloadPath `json:"data"` //查詢結果
}

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

	fmt.Println("已取得參數 parametersUserID=", parametersUserID, ",parametersUserPassword=", parametersUserPassword, ",parametersDeviceID=", parametersDeviceID, ",parametersDeviceBrand=", parametersDeviceBrand)

	// 若順利取出 則進行密碼驗證
	if bindJSONError == nil && bindURIError == nil {

		fmt.Println("取參數正確")

		//checkPassword(parametersUserID,parametersUserPassword)

		// 密碼正確
		if checkPassword(parametersUserID, parametersUserPassword) {

			// 找所有AppsInfo
			result := mongoDB.FindAllAppsInfo()

			// 複製結果到Response格式中(依照原本順序)
			var resultWithDownloadPath []AppsInfoWithDownloadPath

			// 若有結果
			if 0 < len(result) {

				// 結果轉存成，客戶端所使用的json欄位（有大小寫）
				for _, r := range result {

					tempObject := AppsInfoWithDownloadPath{
						AppNameCht: r.AppNameCht,
						AppNameChs: r.AppNameChs,
						AppNameEng: r.AppNameEng,
						AppNameJpn: r.AppNameJpn,
						AppNameKor: r.AppNameKor,

						LabelName:       r.LabelName,
						LastVersionCode: r.LastVersionCode,
						LastVersionName: r.LastVersionName,
						PackageName:     r.PackageName,
						PublishDate:     r.PublishDate,

						ChangeDetailCht: r.ChangeDetailCht,
						ChangeDetailChs: r.ChangeDetailChs,
						ChangeDetailEng: r.ChangeDetailEng,
						ChangeDetailJpn: r.ChangeDetailJpn,
						ChangeDetailKor: r.ChangeDetailKor,

						ChangeBriefCht: r.ChangeBriefCht,
						ChangeBriefChs: r.ChangeBriefChs,
						ChangeBriefEng: r.ChangeBriefEng,
						ChangeBriefJpn: r.ChangeBriefJpn,
						ChangeBriefKor: r.ChangeBriefKor,
					}

					resultWithDownloadPath = append(resultWithDownloadPath, tempObject)
				}
			}

			// 以下方法為兩個物件內的值互相複製
			// if jsonBytes, jsonMarshalError := json.Marshal(result); jsonMarshalError == nil {

			// 	if jsonUnmarshalError := json.Unmarshal(jsonBytes, &resultWithDownloadPath); jsonUnmarshalError != nil {

			// 		logings.SendLog(
			// 			[]string{`將JSON字串 %s 轉成 物件 %+v `},
			// 			[]interface{}{string(jsonBytes), resultWithDownloadPath},
			// 			jsonUnmarshalError,
			// 			logrus.PanicLevel,
			// 		)
			// 	}
			// } else {
			// 	logings.SendLog(
			// 		[]string{`將物件 %+v 轉成 JSON字串 %s `},
			// 		[]interface{}{result, string(jsonBytes)},
			// 		jsonMarshalError,
			// 		logrus.PanicLevel,
			// 	)
			// }

			// fmt.Printf("檢測點：DB結果 %+v", result)
			// fmt.Printf("檢測點：複製的結果 %+v", resultWithDownloadPath)

			// 取APK下載的設定值
			// apkDownloadURLBase := "http://192.168.1.190:63997/appUpdate/download/"
			apkDownloadHost := apiServer.GetConfigValueOrPanic(`apkDownloadHost`)
			apkDownloadPort := apiServer.GetConfigValueOrPanic(`port`) //下載APK port與API port一樣
			apkDownloadURLBase := apiServer.GetConfigValueOrPanic(`apkDownloadURLBase`)

			for i, _ := range resultWithDownloadPath {

				// 為每個結果，組出各APK的「下載網址」
				downloadPath := "http://" + apkDownloadHost + ":" + apkDownloadPort + apkDownloadURLBase + result[i].PackageName //downloadPath
				resultWithDownloadPath[i].DownloadPath = downloadPath                                                            //寫回array

				fmt.Printf("組出downloadPath= %s", downloadPath)
			}

			// fmt.Printf("找到appsInfo結果 %d個", len(result))
			// 包成前端格式
			myResult := AppsInfoResponse{
				IsSuccess: true,
				Results:   "",
				Data:      resultWithDownloadPath,
			}

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
			myResult := AppsInfoResponse{
				IsSuccess: false,
				Results:   "驗證失敗",
				Data:      nil,
			}

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
		myResult := AppsInfoResponse{
			IsSuccess: false,
			Results:   "驗證失敗",
			Data:      nil,
		}

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

	fmt.Println("比對密碼userID=", userID, ",userPassword=", userPassword)

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
