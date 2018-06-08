package library

import (
	"encoding/json"
)

//判断字符串是否是书中的一个值
func InArray(str string, a []string) bool {
	for _, v := range a {
		if v == str {
			return true
		}
	}
	return false
}

//吧json转换成map
func Json2map(jsonStr string) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}
