package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
)

// countEvents - 計算事件紀錄個數
/**
 * @param primitive.M filter 過濾器
 * @retrun int returnCount 事件紀錄個數
 */
func (mongoDB *MongoDB) countEvents(filter primitive.M) (returnCount int) {

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

		eventRWMutex.RLock() // 讀鎖

		// 取得事件紀錄個數
		count, countError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`event-table`)).
			CountDocuments(context.TODO(), filter)

		eventRWMutex.RUnlock() // 讀解鎖

		if mongo.ErrNilDocument == countError {
			countError = nil
		}

		logings.SendLog(
			[]string{`%s %s 取得事件紀錄 %+v 個數 %+v `},
			append(defaultArgs, filter, count),
			countError,
			logrus.InfoLevel,
		)

		if nil != countError { // 若取得事件紀錄個數錯誤，且不為空資料表錯誤
			return // 回傳
		}

		returnCount = int(count)

	}

	return // 回傳
}

// CountEventsByDeviceMethodURLIsSuccessResultsLength1 - 依據裝置、方法、網址、是否成功、結果長度1計算事件記錄數
/**
 * @param string deviceString 裝置字串
 * @param string methodString 方法字串
 * @param string urlString 網址字串
 * @param bool isSuccessBool 是否成功布林
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountEventsByDeviceMethodURLIsSuccessResultsLength1(deviceString, methodString, urlString string, isSuccessBool bool) (result int) {

	result = mongoDB.countEvents(
		bson.M{
			`device`:    deviceString,
			`method`:    methodString,
			`url`:       urlString,
			`issuccess`: isSuccessBool,
			`results`: bson.M{
				sizeConstString: 1,
			},
		},
	)

	return // 回傳
}

// CountEventsByEmployeeIDDeviceMethodURLIsSuccessResultsLength1 - 依據員工編號、裝置、方法、網址、是否成功、結果長度1計算事件記錄數
/**
 * @param string employeeIDString 員工編號字串
 * @param string deviceString 裝置字串
 * @param string methodString 方法字串
 * @param string urlString 網址字串
 * @param bool isSuccessBool 是否成功布林
 * @return int result 取得結果
 */
func (mongoDB *MongoDB) CountEventsByEmployeeIDDeviceMethodURLIsSuccessResultsLength1(
	employeeIDString, deviceString, methodString, urlString string, isSuccessBool bool) (result int) {

	result = mongoDB.countEvents(
		bson.M{
			`employeeid`: employeeIDString,
			`device`:     deviceString,
			`method`:     methodString,
			`url`:        urlString,
			`issuccess`:  isSuccessBool,
			`results`: bson.M{
				sizeConstString: 1,
			},
		},
	)

	return // 回傳
}
