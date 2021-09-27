package jsons

import (
	"encoding/json"

	"leapsy.com/packages/logings"
)

var (
	logger = logings.GetLogger() // 記錄器
)

// JSONBytes - 取得JSON位元組陣列
/**
 * @param interface{} inputObject 輸入物件
 * @return []byte returnJSONBytes 取得JSON位元組陣列
 */
func JSONBytes(inputObject interface{}) (returnJSONBytes []byte) {

	returnJSONBytes, jsonMarshalError := json.Marshal(inputObject) // 轉成JSON

	// 取得記錄器格式和參數
	formatString, args := logings.GetLogFuncFormatAndArguments(
		[]string{`%+v 轉成JSON位元組陣列 %s `},
		[]interface{}{inputObject, string(returnJSONBytes)},
		jsonMarshalError,
	)

	if nil != jsonMarshalError { // 若轉成JSON錯誤
		logger.Errorf(formatString, args...) // 記錄錯誤
		return                               // 回傳
	}

	logger.Infof(formatString, args...) // 記錄資訊

	return // 回傳
}

// JSONString - 取得JSON字串
/**
 * @param interface{} inputObject 輸入物件
 * @return string returnJSONString 取得JSON字串
 */
func JSONString(inputObject interface{}) (returnJSONString string) {
	returnJSONString = string(JSONBytes(inputObject)) // 回傳JSON字串
	return                                            // 回傳
}
