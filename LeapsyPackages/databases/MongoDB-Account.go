package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAllAlertRecords - 取得所有帳號
/**
 * @return results []records.Account 取得結果
 */
func (mongoDB *MongoDB) FindAllAccounts() (results []records.Account) {

	// 取得警報紀錄
	// results = mongoDB.findAlertRecords(bson.M{}, options.Find().SetSort(bson.M{`alerteventtime`: -1}).SetBatchSize(int32(batchSize)))

	results = mongoDB.findAccounts(bson.M{})

	return // 回傳
}

func (mongoDB *MongoDB) FindAccountByUserID(userID string) (results []records.Account) {

	// 取得警報紀錄
	// results = mongoDB.findAlertRecords(bson.M{}, options.Find().SetSort(bson.M{`alerteventtime`: -1}).SetBatchSize(int32(batchSize)))

	results = mongoDB.findAccounts(bson.M{`userID`: userID})

	return // 回傳
}

// findAlertRecords - 查找帳號
/**
 * @param bson.M filter 過濾器
 * @param ...*options.FindOptions opts 選項
 * @return results []records.Account 取得結果
 */
func (mongoDB *MongoDB) findAccounts(filter primitive.M, opts ...*options.FindOptions) (results []records.Account) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空

		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			mongoDB.GetConfigValueOrPanic(`server`),
			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		accountRWMutex.RLock() //讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`account-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		accountRWMutex.RUnlock() //讀解鎖

		// log
		logings.SendLog(
			[]string{`%s %s 查找帳戶 %+v `},
			append(defaultArgs, filter),
			findError,
			logrus.ErrorLevel,
		)

		if nil != findError { // 若查找警報紀錄錯誤
			return // 回傳
		}

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 針對每一紀錄

			var account records.Account

			cursorDecodeError := cursor.Decode(&account) // 解析紀錄

			// log
			logings.SendLog(
				[]string{`%s %s 取得帳戶 %+v `},
				append(defaultArgs, account),
				cursorDecodeError,
				logrus.ErrorLevel,
			)

			// go logger.Errorf(`%+v 取得帳戶 %+s`, append(defaultArgs, account), cursorDecodeError)

			if nil != cursorDecodeError { // 若解析記錄錯誤
				return // 回傳
			}

			// device.AlertEventTime = device.AlertEventTime.Local() // 儲存為本地時間格式

			results = append(results, account) // 儲存紀錄
		}

		cursorErrError := cursor.Err() // 游標錯誤

		// log
		logings.SendLog(
			[]string{`%s %s 查找帳戶遊標運作 `},
			defaultArgs,
			cursorErrError,
			logrus.ErrorLevel,
		)

		// go logger.Errorf(`%+v %s 查找裝置遊標運作`, defaultArgs, cursorErrError)

		if nil != cursorErrError { // 若遊標錯誤
			return // 回傳
		}

		// log
		logings.SendLog(
			[]string{`%s %s  取得裝置 %+v `},
			append(defaultArgs, results),
			nil,
			0,
		)

		// go logger.Infof(`%+v 取得裝置`, append(defaultArgs, results))

	}

	return // 回傳
}

/*** 以下為備用的範例-目前沒用到**/

// findOneAndUpdateAccountSET - 提供更新部分欄位
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndUpdateOptions 選項
 * @return result []records.AlertRecord 更添結果
 */
func (mongoDB *MongoDB) findOneAndUpdateAccountSET(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (results []records.Account) {

	return mongoDB.findOneAndUpdateAccount(
		filter,
		bson.M{
			`$set`: update,
		},
	)

}

// findOneAndUpdateAccount - 提供可以丟整個物件的更新(primitive.M)
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndUpdateOptions 選項
 * @return result []records.AlertRecord 更添結果
 */
func (mongoDB *MongoDB) findOneAndUpdateAccount(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (results []records.Account) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			mongoDB.GetConfigValueOrPanic(`server`),
			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		accountRWMutex.Lock() // 寫鎖

		// 更新警報記錄
		singleResultPointer := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`account-table`)).
			FindOneAndUpdate(
				context.TODO(),
				filter,
				update,
				opts...,
			)

			/*
				FindOneAndUpdate(
					context.TODO(),
					filter,
						bson.M{
							`$set`:update,
						},
						opts...,
					)
			*/

		accountRWMutex.Unlock() // 寫解鎖

		findOneAndUpdateError := singleResultPointer.Err() // 更添錯誤

		if nil != findOneAndUpdateError { // 若更添警報紀錄錯誤且非檔案不存在錯誤

			// log 紀錄有查詢動作
			logings.SendLog(
				[]string{`%s %s 修改帳號密碼，Error= %+v `},
				append(defaultArgs, update),
				findOneAndUpdateError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 修改帳號密碼 `},
			append(defaultArgs, update),
			nil,
			logrus.InfoLevel,
		)

		results = mongoDB.findAccounts(filter)

	}

	return
}

// UpdateOneAccountPassword - 更新帳戶驗證碼
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAccountPassword(userPassword string, userID string) (results []records.Account) {

	updatedModelAccount := mongoDB.findOneAndUpdateAccountSET(
		bson.M{
			`userID`: userID,
		},
		bson.M{
			`userPassword`: userPassword,
		},
	) // 更新的紀錄
	// fmt.Println("標記：", bson.M{`deviceID`: deviceID, `deviceBrand`: deviceBrand}, bson.M{`area.$[]`: newAreaID})

	if nil != updatedModelAccount { // 若更新沒錯誤
		results = append(results, updatedModelAccount...) // 回傳結果
	}

	return
}

// UpdateOneArea - 更新帳戶驗證碼
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAccountVerificationCode(verificationCode string, userID string) (results []records.Account) {

	updatedModelAccount := mongoDB.findOneAndUpdateAccountSET(
		bson.M{
			`userID`: userID,
		},
		bson.M{
			`verificationCode`: verificationCode,
		},
	) // 更新的紀錄
	// fmt.Println("標記：", bson.M{`deviceID`: deviceID, `deviceBrand`: deviceBrand}, bson.M{`area.$[]`: newAreaID})

	if nil != updatedModelAccount { // 若更新沒錯誤
		results = append(results, updatedModelAccount...) // 回傳結果
	}

	return
}

// UpdateOneArea - 更新帳戶有效時間
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAccountPasswordAndVerificationCodeValidPeriod(userPassword string, validPeriod time.Time, userID string) (results []records.Account) {

	updatedModelAccount := mongoDB.findOneAndUpdateAccountSET(
		bson.M{
			`userID`: userID,
		},
		bson.M{
			`userPassword`:                userPassword,
			`verificationCodeValidPeriod`: validPeriod,
			// `verificationCodeTime`: validPeriod,
		},
	) // 更新的紀錄
	// fmt.Println("標記：", bson.M{`deviceID`: deviceID, `deviceBrand`: deviceBrand}, bson.M{`area.$[]`: newAreaID})

	if nil != updatedModelAccount { // 若更新沒錯誤
		results = append(results, updatedModelAccount...) // 回傳結果
	}

	return
}

// UpdateOneArea - 更新帳戶驗證碼與有效時間
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @return *mongo.UpdateResult returnUpdateResult 更新結果
 */
func (mongoDB *MongoDB) UpdateOneAccountVerificationCodeAndVerificationCodeValidPeriod(verificationCode string, validPeriod time.Time, userID string) (results []records.Account) {

	updatedModelAccount := mongoDB.findOneAndUpdateAccountSET(
		bson.M{
			`userID`: userID,
		},
		bson.M{
			// `userPassword`:         userPassword,
			// `verificationCodeTime`: validPeriod,
			`verificationCode`:            verificationCode,
			`verificationCodeValidPeriod`: validPeriod,
		},
	) // 更新的紀錄
	// fmt.Println("標記：", bson.M{`deviceID`: deviceID, `deviceBrand`: deviceBrand}, bson.M{`area.$[]`: newAreaID})

	if nil != updatedModelAccount { // 若更新沒錯誤
		results = append(results, updatedModelAccount...) // 回傳結果
	}

	return
}