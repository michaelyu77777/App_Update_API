package records

// "leapsy.com/packages/networkHub"

// Device - 警報紀錄
type Account struct {
	UserID       string `json:"userid"`// 使用者登入帳號
	UserPassword string `json:"userpassword"`// 使用者登入密碼
}

// PrimitiveM - 轉成primitive.M
/*
 * @return primitive.M returnPrimitiveM 回傳結果
 */
// func (modelAccount Account) Account() (networkHubAccount networkHub.Account) {

// 	networkHubAccount.UserID = modelAccount.UserID

// 	return
// }

// modelAccount := mongoDB.Find....
// modelAccount.Account()
