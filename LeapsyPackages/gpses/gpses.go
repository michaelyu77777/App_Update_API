package gpses

import "regexp"

// IsValid - 是否為合法的GPS字串
/*
 * @param inputGPSString string 輸入GPS字串
 * @return isValid bool 是否合法
 */
func IsValid(inputGPSString string) (isValid bool) {

	isValid = regexp.MustCompile(`^[\+\-]?\d+\.\d+\s*,\s*[\+\-]?\d+\.\d+$`).MatchString(inputGPSString)

	return // 回傳
}
