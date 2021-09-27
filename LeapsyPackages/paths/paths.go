package paths

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"leapsy.com/packages/logings"
)

// IsPathNotExisted - 判斷路徑是否不存在
/**
 * @param  string inputFileNameString  輸入檔名字串
 * @return bool  路徑是否不存在
 */
func IsFileNotExisted(inputFileNameString string) bool {

	_, osStatError := os.Stat(inputFileNameString) // 取得檔案資訊

	return os.IsNotExist(osStatError) // 回傳判斷是否為檔案不存在錯誤
}

// AppendSlashIfNotEndWithOne - 若字串沒以"/"結尾則加"/"
/**
 * @param  string inputString  輸入字串
 * @return string outputString  輸出字串
 */
func AppendSlashIfNotEndWithOne(inputString string) (outputString string) {

	if !strings.HasSuffix(inputString, `/`) { // 若輸入字串不以"/"結尾，則回傳輸入字串加"/"
		outputString = fmt.Sprintf(`%s/`, inputString) // 回傳輸入字串加"/"
	} else { // 若輸入字串以"/"結尾，則回傳輸入字串
		outputString = inputString // 回傳輸入字串
	}

	return // 回傳
}

// CreateIfPathNotExisted - 若路徑不存在則建立路徑
/**
 * @param  string inputPathString  輸入路徑字串
 */
func CreateIfPathNotExisted(inputPathString string) {

	if pathString := AppendSlashIfNotEndWithOne(inputPathString); IsFileNotExisted(pathString) { // 若路徑不存在

		logings.SendLog(
			[]string{`建立路徑 %s `},
			[]interface{}{pathString},
			os.MkdirAll(pathString, 0755),
			logrus.PanicLevel,
		)

	}
}
