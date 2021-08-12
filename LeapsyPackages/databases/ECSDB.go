package databases

import (
	"database/sql"
	"fmt"
	"reflect"

	"leapsy.com/packages/configurations"
	"leapsy.com/packages/network"
	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
)

// ECSDB - 環控系統資料庫
type ECSDB struct {
}

// GetConfigValueOrPanic - 取得設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的值
 */
func (eCSDB *ECSDB) GetConfigValueOrPanic(key string) string {
	return configurations.GetConfigValueOrPanic(reflect.TypeOf(*eCSDB).String(), key)
}

// GetConfigPositiveIntValueOrPanic - 取得正整數設定值否則結束程式
/**
 * @param  string key  關鍵字
 * @return string 設定資料區塊下關鍵字對應的正整數值
 */
func (eCSDB *ECSDB) GetConfigPositiveIntValueOrPanic(key string) int {
	return configurations.GetConfigPositiveIntValueOrPanic(reflect.TypeOf(*eCSDB).String(), key)
}

// Connect - 連接資料庫
/**
 * @return *sql.DB returnDB 資料庫指標
 */
func (eCSDB *ECSDB) Connect() (returnDB *sql.DB) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		eCSDB.GetConfigValueOrPanic(`server`),
		eCSDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	network.SetAddressAlias(address, `環控資料庫`) // 設定預設主機別名

	// 連接預設主機
	db, sqlOpenError := sql.Open(
		"mssql",
		fmt.Sprintf(
			`server=%s;user id=%s;password=%s;database=%s`,
			eCSDB.GetConfigValueOrPanic(`server`),
			eCSDB.GetConfigValueOrPanic(`userid`),
			eCSDB.GetConfigValueOrPanic(`password`),
			eCSDB.GetConfigValueOrPanic(`database`),
		),
	)

	logings.SendLog(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		sqlOpenError,
		logrus.ErrorLevel,
	)

	if nil != sqlOpenError { // 若連接預設主機錯誤
		return // 回傳
	}

	dbPingError := db.Ping() // 確認主機可連接

	logings.SendLog(
		[]string{`%s %s 連接`},
		network.GetAliasAddressPair(address),
		dbPingError,
		logrus.ErrorLevel,
	)

	if nil != dbPingError { // 若確認主機可連接錯誤
		return // 回傳
	}

	returnDB = db // 回傳資料庫指標

	return // 回傳
}

// Disconnect - 中斷與資料庫的連線
/*
 * @params *sql.DB sqlDBPointer 資料庫指標
 */
func (eCSDB *ECSDB) Disconnect(sqlDBPointer *sql.DB) {

	// 預設主機
	address := fmt.Sprintf(
		`%s:%d`,
		eCSDB.GetConfigValueOrPanic(`server`),
		eCSDB.GetConfigPositiveIntValueOrPanic(`port`),
	)

	dbCloseError := sqlDBPointer.Close() // 斷接主機

	logings.SendLog(
		[]string{`%s %s 斷接`},
		network.GetAliasAddressPair(address),
		dbCloseError,
		logrus.ErrorLevel,
	)

}
