/*
   百度翻译
   优点：每月200万字符内免费
   缺点：长篇幅，只翻译第一行，需要自行拆开 单个句子不能长于2000字符 接口每秒不能超过5次
    结果相对谷歌来说，太垃圾
*/
package langtranslate

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	BAIDU_API_URL       = "http://api.fanyi.baidu.com/api/trans/vip/translate"
	BAIDU_API_URL_HTTPS = "https://api.fanyi.baidu.com/api/trans/vip/translate"
)

const MAX_TRY_TIMES = 5

type BaiduTransTask struct {
	fromLang  string
	toLang    string
	queryStr  string
	retStr    string
	status    int
	failMsg   string
	failTimes int
}

// this type implement Translator
type BaiduTranslator struct {
	fromLang string
	toLang   string
}

// trans res defined by baidu
type BaiduTransRes struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Trans_result []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

var baidu_appid string
var baidu_appsecret string

func init() {
}

func (this *BaiduTranslator) SetApiConfig(conf map[string]interface{}) {
	if _, ok := conf["appid"]; ok {
		baidu_appid, _ = conf["appid"].(string)
	}
	if _, ok := conf["appsecret"]; ok {
		baidu_appsecret, _ = conf["appsecret"].(string)
	}
}

func (this *BaiduTranslator) SetFromLang(lang string) {
	this.fromLang = lang
}

func (this *BaiduTranslator) SetToLang(lang string) {
	this.toLang = lang
}

func (this *BaiduTranslator) Translate(str string) (res string, err error) {
	failTimes := 0
	task := &BaiduTransTask{queryStr: str, fromLang: this.fromLang, toLang: this.toLang}

	for failTimes = 0; failTimes < MAX_TRY_TIMES; failTimes++ {
		res, err = task.RunTrans()
		if err == nil {
			break
		}
	}

	return res, err
}

/*
   params = ["q"=>"words for translate"]
*/
func makeSignOfBaidu(params map[string]string) string {
	//:= fmt.Sprintf("%d", time.Now().Unix())

	longStr := baidu_appid + params["q"] + params["salt"] + baidu_appsecret

	h := md5.New()
	h.Write([]byte(longStr))
	signByte := h.Sum(nil)

	sign := fmt.Sprintf("%x", signByte)
	return sign
}

func (this *BaiduTransTask) RunTrans() (string, error) {

	if this.fromLang == "" {
		this.fromLang = "auto"
	}
	if this.toLang == "" {
		this.fromLang = "auto"
	}

	contentType := "application/x-www-form-urlencoded"

	salt := fmt.Sprintf("%d", time.Now().Unix())
	sign := makeSignOfBaidu(map[string]string{"q": this.queryStr, "salt": salt})
	body := "q=" + this.queryStr + "&from=" + this.fromLang + "&to=" + this.toLang + "&appid=" + baidu_appid + "&salt=" + salt + "&sign=" + sign
	q, _ := url.ParseQuery(body)
	bodyEncoded := q.Encode()

	//fmt.Println("TRANS BY BAIDU:", body, bodyEncoded)

	res, err := http.Post(BAIDU_API_URL, contentType, bytes.NewReader([]byte(bodyEncoded)))
	if err != nil {
		return "", err
	}

	// do translate
	//this.retStr = str
	allRetBytes, err := ioutil.ReadAll(res.Body)

	//fmt.Println("getTransRes:", string(allRetBytes))

	transRes := new(BaiduTransRes)
	err = json.Unmarshal(allRetBytes, transRes)
	if err != nil {
		return "", err
	}

	if len(transRes.Trans_result) > 0 {
		this.retStr = transRes.Trans_result[0].Dst
	}

	//fmt.Println("TRANS OVER:", this.retStr, err)

	return this.retStr, err

}
