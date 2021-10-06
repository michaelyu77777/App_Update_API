package records

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*[改用]新的結構*/
// 軟體資訊
// type AppsInfo struct {
// 	AppNameCht string `json:"appNameCht"` //軟體名稱 正體
// 	AppNameChs string `json:"appNameChs"` //軟體名稱 簡體
// 	AppNameEng string `json:"appNameEng"` //軟體名稱 英文
// 	AppNameJpn string `json:"appNameJpn"` //軟體名稱 日文
// 	AppNameKor string `json:"appNameKor"` //軟體名稱 韓文

// 	LastVersionCode int    `json:"lastVersionCode"` //最新版本號
// 	LastVersionName string `json:"lastVersionName"` //最新版本名
// 	PackageName     string `json:"packageName"`     //封包名稱
// 	DownloadPath    string `json:"downloadPath"`    //下載APK網址
// 	PublishDate     string `json:"publishDate"`     //發佈日期

// 	ChangeDetailCht string `json:"changeDetailCht"` //更新內容 詳述 正
// 	ChangeDetailChs string `json:"changeDetailChs"` //更新內容 詳述 簡
// 	ChangeDetailEng string `json:"changeDetailEng"` //更新內容 詳述 英
// 	ChangeDetailJpn string `json:"changeDetailJpn"` //更新內容 詳述 日
// 	ChangeDetailKor string `json:"changeDetailKor"` //更新內容 詳述 韓

// 	ChangeBriefCht string `json:"changeBriefCht"` //更新內容 簡述 正
// 	ChangeBriefChs string `json:"changeBriefChs"` //更新內容 簡述 簡
// 	ChangeBriefEng string `json:"changeBriefEng"` //更新內容 簡述 英
// 	ChangeBriefJpn string `json:"changeBriefJpn"` //更新內容 簡述 日
// 	ChangeBriefKor string `json:"changeBriefKor"` //更新內容 簡述 韓

// 	ApkDirectoryName string `json:"apkDirectoryName"` //存放APK資料夾名稱
// 	ApkFileName      string `json:"apkFileName"`      //APK檔案名稱
// }

// // 包成回給前端<取AppsInfo格式>
// type AppsInfoResponse struct {
// 	IsSuccess bool       `json:"isSuccess"` //錯誤代碼
// 	Results   string     `json:"results"`   //錯誤訊息
// 	Data      []AppsInfo `json:"data"`      //查詢結果
// }

// // 回給前端<一般格式>
// type APIResponse struct {
// 	IsSuccess bool   `json:"isSuccess"` //錯誤代碼
// 	Results   string `json:"results"`   //錯誤訊息
// }

/* [改用] 下面新的函數 */
// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
// func (appsInfo AppsInfo) PrimitiveM() (returnPrimitiveM primitive.M) {

// 	returnPrimitiveM = bson.M{
// 		`appNameCht`: appsInfo.AppNameCht,
// 		`appNameChs`: appsInfo.AppNameChs,
// 		`appNameEng`: appsInfo.AppNameEng,
// 		`appNameJpn`: appsInfo.AppNameJpn,
// 		`appNameKor`: appsInfo.AppNameKor,

// 		`lastVersionCode`: appsInfo.LastVersionCode,
// 		`lastVersionName`: appsInfo.LastVersionName,
// 		`packageName`:     appsInfo.PackageName,
// 		`downloadPath`:    appsInfo.DownloadPath,
// 		`publishDate`:     appsInfo.PublishDate,

// 		`changeDetailCht`: appsInfo.ChangeDetailCht,
// 		`changeDetailChs`: appsInfo.ChangeDetailChs,
// 		`changeDetailEng`: appsInfo.ChangeDetailEng,
// 		`changeDetailJpn`: appsInfo.ChangeDetailJpn,
// 		`changeDetailKor`: appsInfo.ChangeDetailKor,

// 		`changeBriefCht`: appsInfo.ChangeBriefCht,
// 		`changeBriefChs`: appsInfo.ChangeBriefChs,
// 		`changeBriefEng`: appsInfo.ChangeBriefEng,
// 		`changeBriefJpn`: appsInfo.ChangeBriefJpn,
// 		`changeBriefKor`: appsInfo.ChangeBriefKor,

// 		`apkDirectoryName`: appsInfo.ApkDirectoryName,
// 		`apkFileName`:      appsInfo.ApkFileName,
// 	}

// 	return
// }

/**從model搬過來的AppsInfo**/
// 放 DB 與 Respone Client 共同會用到的參數
type AppsInfoCommonStruct struct {
	AppNameCht string `json:"appNameCht"` //軟體名稱 正體
	AppNameChs string `json:"appNameChs"` //軟體名稱 簡體
	AppNameEng string `json:"appNameEng"` //軟體名稱 英文
	AppNameJpn string `json:"appNameJpn"` //軟體名稱 日文
	AppNameKor string `json:"appNameKor"` //軟體名稱 韓文

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

}

/*[改用] 新的結構*/
// 軟體資訊(DB用)
type AppsInfo struct {
	AppsInfoCommonStruct `bson:",inline"` //共用參數：會從DB拿、也會Response回Client的參數

	ApkDirectoryName string `json:"apkDirectoryName"` // 存放APK資料夾名稱
	ApkFileName      string `json:"apkFileName"`      // APK檔案名稱
	LabelName        string `json:"labelName"`        // APK Label名稱
}

// 軟體資訊(Response client用)
type AppsInfoWithDownloadPath struct {
	AppsInfoCommonStruct `bson:",inline"` //共用參數：會從DB拿、也會Response回Client的參數

	DownloadPath string `json:"downloadPath"` // 組合出下載APK網址
}

// type AppsInfo struct {
// 	ProjectName  string //專案名稱
// 	AppName      string //軟體名稱
// 	LastVersion  string //最新版本號
// 	PackageName  string //封包名稱
// 	DownloadPath string //下載URL
// }

// 包成回給前端<取AppsInfo格式>
type AppsInfoResponse struct {
	IsSuccess bool   `json:"isSuccess"` //錯誤代碼
	Results   string `json:"results"`   //錯誤訊息
	// Data      []AppsInfo `json:"data"`      //查詢結果
	Data []AppsInfoWithDownloadPath `json:"data"` //查詢結果
}

// 回給前端<一般格式>
type APIResponse struct {
	IsSuccess bool   `json:"isSuccess"` //錯誤代碼
	Results   string `json:"results"`   //錯誤訊息
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (appsInfo AppsInfo) PrimitiveM() (returnPrimitiveM primitive.M) {

	returnPrimitiveM = bson.M{
		`appNameCht`: appsInfo.AppNameCht,
		`appNameChs`: appsInfo.AppNameChs,
		`appNameEng`: appsInfo.AppNameEng,
		`appNameJpn`: appsInfo.AppNameJpn,
		`appNameKor`: appsInfo.AppNameKor,

		`lastVersionCode`: appsInfo.LastVersionCode,
		`lastVersionName`: appsInfo.LastVersionName,
		`packageName`:     appsInfo.PackageName,
		`publishDate`:     appsInfo.PublishDate,

		`changeDetailCht`: appsInfo.ChangeDetailCht,
		`changeDetailChs`: appsInfo.ChangeDetailChs,
		`changeDetailEng`: appsInfo.ChangeDetailEng,
		`changeDetailJpn`: appsInfo.ChangeDetailJpn,
		`changeDetailKor`: appsInfo.ChangeDetailKor,

		`changeBriefCht`: appsInfo.ChangeBriefCht,
		`changeBriefChs`: appsInfo.ChangeBriefChs,
		`changeBriefEng`: appsInfo.ChangeBriefEng,
		`changeBriefJpn`: appsInfo.ChangeBriefJpn,
		`changeBriefKor`: appsInfo.ChangeBriefKor,

		`apkDirectoryName`: appsInfo.ApkDirectoryName, // 存放APK資料夾名稱
		`apkFileName`:      appsInfo.ApkFileName,      // APK檔案名稱
	}

	return
}
