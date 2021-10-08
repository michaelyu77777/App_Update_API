package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 尋找所有 apps info
func (mongoDB *MongoDB) FindAllAppsInfo() (results []records.AppsInfo) {

	// 回傳結果
	results = mongoDB.findAppsInfo(
		bson.M{},
	)

	return // 回傳
}

// 尋找符合的專案名稱,app名稱
func (mongoDB *MongoDB) FindAppsInfoByProjectNameAndAppName(projectName string, appName string) (results []records.AppsInfo) {

	// 回傳結果
	results = mongoDB.findAppsInfo(
		bson.M{
			"projectname": projectName,
			"appname":     appName,
		},
	)

	return // 回傳
}

/**以下為複製過來的函數**/


// 尋找符合的專案名稱,app名稱
func (mongoDB *MongoDB) FindAppsInfoByLabelName(labelName string) (results []records.AppsInfo) {

	// 回傳結果
	results = mongoDB.findAppsInfo(
		bson.M{
			"labelname": labelName,
		},
	)

	return // 回傳
}

// 查找[資料庫-軟體更新]appsinfo(軟體資訊)
func (mongoDB *MongoDB) findAppsInfo(filter primitive.M, opts ...*options.FindOptions) (results []records.AppsInfo) {

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

		appsInfoRWMutex.RLock() // 讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`appsInfo-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		appsInfoRWMutex.RUnlock() // 讀解鎖

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 查找資料庫-軟體更新 %+v `},
			append(defaultArgs, filter),
			findError,
			logrus.ErrorLevel,
		)

		if nil != findError { // 若查找錯誤
			return // 回傳
		}

		defer cursor.Close(context.TODO()) // 記得關閉

		for cursor.Next(context.TODO()) { // 拜訪每筆查詢

			var appsInfo records.AppsInfo

			cursorDecodeError := cursor.Decode(&appsInfo) // 解析單筆結果，放到appsInfo變數中

			// log 針對查出的每筆紀錄寫log
			logings.SendLog(
				[]string{`%s %s 取得資料庫-軟體更新 %+v `},
				append(defaultArgs, appsInfo),
				cursorDecodeError,
				logrus.ErrorLevel,
			)

			if nil != cursorDecodeError { // 若解析記錄錯誤
				return // 回傳
			}

			// appsInfo.Time = appsInfo.Time.Local() // 儲存為本地時間格式

			results = append(results, appsInfo) // 將單筆查詢結果，加入到results結果中
		}

		cursorErrError := cursor.Err() // 結果游標錯誤

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 查找資料庫-軟體資訊 游標`},
			defaultArgs,
			cursorErrError,
			logrus.ErrorLevel,
		)

		if nil != cursorErrError { // 若遊標錯誤
			return // 回傳
		}

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 取得資料庫-軟體更新 %+v `},
			append(defaultArgs, results),
			nil,
			0,
		)

	}

	return // 回傳
}

// findOneAndUpdateAccountSET - 提供更新部分欄位
/**
 * @param primitive.M filter 過濾器
 * @param primitive.M update 更新
 * @param ...*options.FindOneAndUpdateOptions 選項
 * @return result []records.AlertRecord 更添結果
 */
func (mongoDB *MongoDB) FindOneAndUpdateAppsInfoSET(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (results []records.AppsInfo) {

	return mongoDB.findOneAndUpdateAppsInfo(
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
func (mongoDB *MongoDB) findOneAndUpdateAppsInfo(
	filter, update primitive.M,
	opts ...*options.FindOneAndUpdateOptions) (results []records.AppsInfo) {

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

		appsInfoRWMutex.Lock() // 寫鎖

		// 更新
		singleResultPointer := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`appsInfo-table`)).
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

		appsInfoRWMutex.Unlock() // 寫解鎖

		findOneAndUpdateError := singleResultPointer.Err() // 更添錯誤

		if nil != findOneAndUpdateError { // 若更添警報紀錄錯誤且非檔案不存在錯誤

			// log 紀錄有查詢動作
			logings.SendLog(
				[]string{`%s %s 解析APK並更新資料庫AppsInfo，Error= %+v `},
				append(defaultArgs, update),
				findOneAndUpdateError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 解析APK並更新資料庫AppsInfo `},
			append(defaultArgs, update),
			nil,
			logrus.InfoLevel,
		)

		results = mongoDB.findAppsInfo(filter)

	}

	return
}

// insertOneAppsInfo - 新增一筆AppsInfo
func (mongoDB *MongoDB) InsertOneAppsInfo(appsInfo records.AppsInfo) (insertOneError error) {

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

		appsInfoRWMutex.Lock() // 寫鎖

		// 添加一筆軟體訊息
		_, insertOneError = mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`appsInfo-table`)).
			InsertOne(
				context.TODO(),
				appsInfo,
			)

		appsInfoRWMutex.Unlock() // 寫解鎖

		if nil != insertOneError { // 若更添警報紀錄錯誤且非檔案不存在錯誤

			// log 紀錄有查詢動作
			logings.SendLog(
				[]string{`%s %s 新增一筆資料到資料庫AppsInfo，Error= %+v `},
				append(defaultArgs, appsInfo),
				insertOneError,
				logrus.ErrorLevel,
			)

			return // 回傳
		}

		// log 紀錄有查詢動作
		logings.SendLog(
			[]string{`%s %s 新增一筆資料到資料庫AppsInfo `},
			append(defaultArgs, appsInfo),
			nil,
			logrus.InfoLevel,
		)

		if nil != insertOneError { // 若添寫公司紀錄錯誤
			return // 回傳
		}

		// results = mongoDB.findAppsInfo(filter)
		fmt.Println("已新增一筆資料到appsInfo")

	}

	return
}
