package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TokenReturn struct {
	Code    int
	Message string
	Data    struct {
		Token  string
		Expire string
	}
}

//获取token
/***
*return token,expire,code
 */
func GetToken(url string, param string) (string, int64, int) {
	newurl := fmt.Sprintf("%s?%s", url, param)
	resp, _ := http.Get(newurl)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, -1
	}
	var tokenReturn = new(TokenReturn)
	json.Unmarshal(body, &tokenReturn)
	code := tokenReturn.Code
	if tokenReturn.Code == 1 {
		token := tokenReturn.Data.Token
		expire := tokenReturn.Data.Expire
		tokenExpire, _ := time.Parse("2006-01-02 15:04:05", expire)
		return token, tokenExpire.Unix(), code
	} else {
		return "", 0, -1
	}

}
