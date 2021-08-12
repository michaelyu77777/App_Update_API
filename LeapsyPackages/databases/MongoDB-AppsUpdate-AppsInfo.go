package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/model"
	"leapsy.com/packages/network"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 查找[資料庫-軟體更新]appsinfo(軟體資訊)
func (mongoDB *MongoDB) findAppsInfo(filter primitive.M, opts ...*options.FindOptions) (results []model.AppsInfo) {

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

		hourlyRWMutex.RLock() // 讀鎖

		// 查找紀錄
		cursor, findError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`appsInfo-table`)).
			Find(
				context.TODO(),
				filter,
				opts...,
			)

		hourlyRWMutex.RUnlock() // 讀解鎖

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

			var appsInfo model.AppsInfo

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

// 尋找所有 apps info
func (mongoDB *MongoDB) FindAllAppsInfoByProjectNameAndAppName() (results []model.AppsInfo) {

	// if !lowerTime.IsZero() && !upperTime.IsZero() { //若上下限時間不為零時間

	// var (
	// 	greaterThanKeyword, lessThanKeyword string // 比較關鍵字
	// )

	// if !isLowerTimeIncluded { // 若不包含下限時間
	// 	greaterThanKeyword = greaterThanConstString // >
	// } else {
	// 	greaterThanKeyword = greaterThanEqualToConstString // >=
	// }

	// if !isUpperTimeIncluded { // 若不包含上限時間
	// 	lessThanKeyword = lessThanConstString // <
	// } else {
	// 	lessThanKeyword = lessThanEqualToConstString // <=
	// }

	// 回傳結果
	results = mongoDB.findAppsInfo(
		bson.M{},
	)

	// }

	return // 回傳
}

// 尋找符合的專案名稱,app名稱
func (mongoDB *MongoDB) FindAppsInfoByProjectNameAndAppName(projectName string, appName string) (results []model.AppsInfo) {

	fmt.Println("測試中：取得參數projectName=", projectName, "appName=", appName)

	// 回傳結果
	results = mongoDB.findAppsInfo(
		bson.M{
			"projectName": projectName,
			"appName":     appName,
		},
	)

	// }

	return // 回傳
}
