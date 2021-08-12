package databases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
)

// deleteManyAlert - 刪除一筆警報紀錄
/**
 * @param primitive.M filter 過濾器
 */
func (mongoDB *MongoDB) deleteManyAlert(filter primitive.M) {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			mongoDB.GetConfigValueOrPanic(`server`),
			mongoDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		alertRWMutex.Lock() // 寫鎖

		// 刪除警報記錄
		_, deleteManyError := mongoClientPointer.
			Database(mongoDB.GetConfigValueOrPanic(`database`)).
			Collection(mongoDB.GetConfigValueOrPanic(`alert-table`)).
			DeleteMany(context.TODO(), filter)

		alertRWMutex.Unlock() // 寫解鎖

		logings.SendLog(
			[]string{`%s %s 刪除警報記錄 %v `},
			network.GetAliasAddressPair(address),
			deleteManyError,
			logrus.ErrorLevel,
		)

	}

}

// DeleteAllAlertRecords - 刪除警報紀錄
/**
 */
func (mongoDB *MongoDB) DeleteAllAlertRecords() {

	mongoClientPointer := mongoDB.Connect() // 資料庫指標

	if nil != mongoClientPointer { // 若資料庫指標不為空
		defer mongoDB.Disconnect(mongoClientPointer) // 記得關閉資料庫指標

		// mongoDB.deleteManyAlert(bson.M{}) // 刪除所有警報紀錄

	}

}
