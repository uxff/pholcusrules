/*
   有道翻译 利用浏览器翻译插件接口
   优点：永久免费 可以POST
   缺点：长句子要拆开翻译 只能en->zh
*/
package langtranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	YOUDAO_API_URL      = "http://fanyi.youdao.com/translate"
	YOUDAO_API_URL_FULL = "http://fanyi.youdao.com/translate?client=deskdict&keyfrom=chrome.extension&xmlVersion=1.1&dogVersion=1.0&ue=utf8&i=implements%20hiveserver2(thrift%20rpc)%20client%20in%20golang&doctype=xml"
	YOUDAO_API_COOKIE   = "OUTFOX_SEARCH_USER_ID=1960626865@111.199.186.18; JSESSIONID=abcz62zd30LRdhL_7q3ew"
)

var youdao_cookie = YOUDAO_API_COOKIE

type YoudaoTranslator struct {
	fromLang string
	toLang   string
	queryStr string
	retStr   string
}

type YoudaoTransRes struct {
	Type            string `json:"type"`
	ErrorCode       int    `json:"errorCode"`
	ElapsedTime     int    `json:"elapsedTime"`
	TranslateResult [][]struct {
		Src string `json:"src"`
		Tgt string `json:"tgt"`
	} `json:"translateResult"`
}

func (this *YoudaoTranslator) SetApiConfig(conf map[string]interface{}) {
	cookie, ok := conf["cookie"].(string)
	if ok {
		youdao_cookie = cookie + "; timestamp=" + fmt.Sprintf("%d", time.Now().Unix())
	}
}

func (this *YoudaoTranslator) SetFromLang(lang string) {
	this.fromLang = lang
}
func (this *YoudaoTranslator) SetToLang(lang string) {
	this.toLang = lang
}
func (this *YoudaoTranslator) Translate(str string) (string, error) {

	params := "client=deskdict&keyfrom=chrome.extension&jsonVersion=1&dogVersion=1.0&ue=utf8&doctype=json&i=" + str
	contentType := "application/x-www-form-urlencoded"

	transRetOfApi, err := http.Post(YOUDAO_API_URL, contentType, bytes.NewReader([]byte(params)))
	if err != nil {
		return "", err
	}

	willContentType := "application/json"
	realContentType := transRetOfApi.Header.Get("Content-Type")
	if realContentType[:len(willContentType)] != "application/json" {
		return "", fmt.Errorf("when youdao translate, error content type from api:%v, EXPECT application/json", realContentType)
	}

	targetRes, err := ioutil.ReadAll(transRetOfApi.Body)
	if err != nil {
		return "", err
	}

	youdaoTransRes := new(YoudaoTransRes)
	err = json.Unmarshal(targetRes, youdaoTransRes)
	if err != nil {
		return "", err
	}

	if len(youdaoTransRes.TranslateResult) > 0 {
		if len(youdaoTransRes.TranslateResult[0]) > 0 {
			return youdaoTransRes.TranslateResult[0][0].Tgt, nil
		}
	}

	contentLineAbsLen := len(targetRes)
	if contentLineAbsLen > 20 {
		contentLineAbsLen = 20
	}
	return "", fmt.Errorf("response could not unmarshal:%s...", targetRes[:contentLineAbsLen])
}
