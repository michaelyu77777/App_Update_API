package records

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**從model搬過來的AppsInfo**/
// 放 DB 與 Respone Client 共同會用到的參數
type AppsInfoCommonStruct struct {
	AppNameCht string `json:"appnamecht"` //軟體名稱 正體
	AppNameChs string `json:"appnamechs"` //軟體名稱 簡體
	AppNameEng string `json:"appnameeng"` //軟體名稱 英文
	AppNameJpn string `json:"appnamejpn"` //軟體名稱 日文
	AppNameKor string `json:"appnamekor"` //軟體名稱 韓文

	LabelName       string `json:"labelname"`       // APK Label名稱
	LastVersionCode int    `json:"lastversioncode"` //最新版本號
	LastVersionName string `json:"lastversionname"` //最新版本名
	PackageName     string `json:"packagename"`     //封包名稱
	PublishDate     string `json:"publishdate"`     //發佈日期

	ChangeDetailCht string `json:"changedetailcht"` //更新內容 詳述 正
	ChangeDetailChs string `json:"changedetailchs"` //更新內容 詳述 簡
	ChangeDetailEng string `json:"changedetaileng"` //更新內容 詳述 英
	ChangeDetailJpn string `json:"changedetailjpn"` //更新內容 詳述 日
	ChangeDetailKor string `json:"changedetailkor"` //更新內容 詳述 韓

	ChangeBriefCht string `json:"changebriefcht"` //更新內容 簡述 正
	ChangeBriefChs string `json:"changebriefchs"` //更新內容 簡述 簡
	ChangeBriefEng string `json:"changebriefeng"` //更新內容 簡述 英
	ChangeBriefJpn string `json:"changebriefjpn"` //更新內容 簡述 日
	ChangeBriefKor string `json:"changebriefkor"` //更新內容 簡述 韓
}

/*[改用] 新的結構*/
// 軟體資訊(DB用)
type AppsInfo struct {
	AppsInfoCommonStruct `bson:",inline"` //共用參數：會從DB拿、也會Response回Client的參數

	ApkFileName string `json:"apkfilename"` // APK檔案名稱
	// LabelName        string `json:"labelname"`        // APK Label名稱
}

// 軟體資訊(Response client用)
type AppsInfoWithDownloadPath struct {
	AppsInfoCommonStruct `bson:",inline"` //共用參數：會從DB拿、也會Response回Client的參數

	DownloadPath string `json:"downloadpath"` // 組合出下載APK網址
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
func (appsInfo AppsInfo) PrimitiveM() (returnPrimitiveM primitive.M) {

	returnPrimitiveM = bson.M{
		`appnamecht`: appsInfo.AppNameCht,
		`appnamechs`: appsInfo.AppNameChs,
		`appnameeng`: appsInfo.AppNameEng,
		`appnamejpn`: appsInfo.AppNameJpn,
		`appnamekor`: appsInfo.AppNameKor,

		`lastversioncode`: appsInfo.LastVersionCode,
		`lastversionname`: appsInfo.LastVersionName,
		`packagename`:     appsInfo.PackageName,
		`publishdate`:     appsInfo.PublishDate,

		`changedetailcht`: appsInfo.ChangeDetailCht,
		`changedetailchs`: appsInfo.ChangeDetailChs,
		`changedetaileng`: appsInfo.ChangeDetailEng,
		`changedetailjpn`: appsInfo.ChangeDetailJpn,
		`changedetailkor`: appsInfo.ChangeDetailKor,

		`changebriefcht`: appsInfo.ChangeBriefCht,
		`changebriefchs`: appsInfo.ChangeBriefChs,
		`changebriefeng`: appsInfo.ChangeBriefEng,
		`changebriefjpn`: appsInfo.ChangeBriefJpn,
		`changebriefkor`: appsInfo.ChangeBriefKor,

		`apkfilename`: appsInfo.ApkFileName, // APK檔案名稱

	}

	return
}

// 包成回給前端<取AppsInfo格式>
type AppsInfoResponse struct {
	IsSuccess bool   `json:"issuccess"` //錯誤代碼
	Results   string `json:"results"`   //錯誤訊息
	// Data      []AppsInfo `json:"data"`      //查詢結果
	Data []AppsInfoWithDownloadPath `json:"data"` //查詢結果
}

// 回給前端<一般格式>
type APIResponse struct {
	IsSuccess bool   `json:"issuccess"` //錯誤代碼
	Results   string `json:"results"`   //錯誤訊息
}
