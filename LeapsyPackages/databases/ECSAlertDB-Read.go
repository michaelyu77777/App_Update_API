package databases

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
	"leapsy.com/packages/network"
	"leapsy.com/records"
)

// CountAll - 計算所有紀錄個數
/**
 * @return int returnCount 紀錄個數
 */
func (ecsAlertDB *ECSAlertDB) CountAll() (returnCount int) {

	sqlDBPointer := ecsAlertDB.Connect() // 資料庫指標

	if nil != sqlDBPointer { // 若資料庫指標不為空

		defer ecsAlertDB.Disconnect(sqlDBPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			ecsAlertDB.GetConfigValueOrPanic(`server`),
			ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		// 查詢紀錄數
		row := sqlDBPointer.QueryRow(
			fmt.Sprintf(
				`select count(*) from %s`,
				ecsAlertDB.GetConfigValueOrPanic(`table`),
			),
		)

		// 審視環控警報紀錄個數
		scanError := row.Scan(&returnCount)

		logings.SendLog(
			[]string{`%s %s 取得環控警報紀錄個數 %d `},
			append(defaultArgs, returnCount),
			scanError,
			logrus.ErrorLevel,
		)

	}

	return // 回傳
}

// Read - 讀一筆紀錄
/**
 * @return []records.ECSAlertRecord ecsAlertRecord 紀錄
 */
func (ecsAlertDB *ECSAlertDB) Read() (ecsAlertRecords []records.ECSAlertRecord) {

	ecsAlertRecords = ecsAlertDB.getRecords(
		fmt.Sprintf(
			`select * from %s`,
			ecsAlertDB.GetConfigValueOrPanic(`table`),
		),
	)

	return // 回傳
}

// ReadLast - 讀末N筆紀錄
/**
 * @param int n 紀錄個數
 * @return []records.ECSAlertRecord ecsAlertRecord 紀錄
 */
func (ecsAlertDB *ECSAlertDB) ReadLast(n int) (ecsAlertRecords []records.ECSAlertRecord) {

	ecsAlertRecords = ecsAlertDB.getRecords(
		fmt.Sprintf(
			`select top %d * from %s order by ALERTEVENTID desc`,
			n,
			ecsAlertDB.GetConfigValueOrPanic(`table`),
		),
	)

	return // 回傳
}

//  getRecords - 取得紀錄
/**
 * @param  string sqlCommand SQL指令
 * @return []records.ECSAlertRecord record 紀錄
 */
func (ecsAlertDB *ECSAlertDB) getRecords(sqlCommand string) (ecsAlertRecords []records.ECSAlertRecord) {

	sqlDBPointer := ecsAlertDB.Connect() // 資料庫指標

	if nil != sqlDBPointer { // 若資料庫指標不為空

		defer ecsAlertDB.Disconnect(sqlDBPointer) // 記得關閉資料庫指標

		// 預設主機
		address := fmt.Sprintf(
			`%s:%d`,
			ecsAlertDB.GetConfigValueOrPanic(`server`),
			ecsAlertDB.GetConfigPositiveIntValueOrPanic(`port`),
		)

		defaultArgs := network.GetAliasAddressPair(address) // 預設參數

		// 查詢紀錄
		rows, queryError := sqlDBPointer.Query(sqlCommand)
		defer rows.Close()

		logings.SendLog(
			[]string{`%s %s 查詢環控警報紀錄`},
			defaultArgs,
			queryError,
			logrus.ErrorLevel,
		)

		if nil != queryError { // 若查詢紀錄錯誤
			return // 回傳
		}

		for rows.Next() {

			var ecsAlertRecord records.ECSAlertRecord

			// 審視環控警報紀錄
			scanError := rows.Scan(
				&ecsAlertRecord.ALERTEVENTID,
				&ecsAlertRecord.ALERTEVENTTIME,
				&ecsAlertRecord.VARTAG,
				&ecsAlertRecord.COMMENT,
				&ecsAlertRecord.ALERTTYPE,
				&ecsAlertRecord.LINETEXT,
			)

			logings.SendLog(
				[]string{`%s %s 取得環控警報紀錄 %+v `},
				append(defaultArgs, ecsAlertRecord),
				scanError,
				logrus.ErrorLevel,
			)

			if nil != scanError { // 若審視環控警報紀錄錯誤
				return // 回傳
			}

			ecsAlertRecords = append(ecsAlertRecords, ecsAlertRecord)

		}

		logings.SendLog(
			[]string{`%s %s 取得環控警報紀錄 %+v `},
			append(defaultArgs, ecsAlertRecords),
			nil,
			0,
		)

	}

	return // 回傳
}
