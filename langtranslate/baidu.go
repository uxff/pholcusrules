/*
   百度翻译
   优点：每月200万字符内免费
   缺点：长篇幅，只翻译第一行，需要自行拆开
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
	STATUS_NONE  = 0
	STATUS_DOING = 1
	STATUS_OK    = 2
	STATUS_FAIL  = 3
)

const (
	BAIDU_API_URL       = "http://api.fanyi.baidu.com/api/trans/vip/translate"
	BAIDU_API_URL_HTTPS = "https://api.fanyi.baidu.com/api/trans/vip/translate"
)

type BaiduTransTask struct {
	id       int
	queryStr string
	retStr   string
	status   int
	failMsg  string
}

type BaiduTranslator struct {
	fromLang string
	toLang   string
	queryStr string
	retStr   string
	tasks    map[int]*BaiduTransTask
}

type BaiduTransRes struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Trans_result []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

var transTaskNextId = 1
var baidu_appid string
var baidu_appsecret string

func (this *BaiduTranslator) SetApiConfig(conf map[string]interface{}) {
	baidu_appid, _ = conf["appid"].(string)
	baidu_appsecret, _ = conf["appsecret"].(string)
}

func (this *BaiduTranslator) SetFromLang(lang string) {
	this.fromLang = lang
}
func (this *BaiduTranslator) SetToLang(lang string) {
	this.toLang = lang
}
func (this *BaiduTranslator) Translate(str string) (string, error) {

	if this.fromLang == "" {
		this.fromLang = "auto"
	}
	if this.toLang == "" {
		this.fromLang = "auto"
	}

	this.queryStr = str

	contentType := "application/x-www-form-urlencoded"

	salt := fmt.Sprintf("%d", time.Now().Unix())
	sign := makeSign(map[string]string{"q": str, "salt": salt})
	body := "q=" + str + "&from=" + this.fromLang + "&to=" + this.toLang + "&appid=" + baidu_appid + "&salt=" + salt + "&sign=" + sign
	q, _ := url.ParseQuery(body)
	bodyEncoded := q.Encode()

	//fmt.Println("TRANS BY BAIDU:", body, bodyEncoded)

	res, err := http.Post(BAIDU_API_URL, contentType, bytes.NewReader([]byte(bodyEncoded)))
	if err != nil {
		return "", err
	}

	// to do translate
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
func (this *BaiduTranslator) AsyncTranslate(str string) int {
	task := &BaiduTransTask{id: transTaskNextId, queryStr: this.queryStr, status: STATUS_NONE}
	transTaskNextId++
	if this.tasks == nil {
		this.tasks = make(map[int]*BaiduTransTask, 0)
	}

	//this.tasks = append(this.tasks, task)
	this.tasks[transTaskNextId] = task
	//task.Start()
	return task.id
}
func (this *BaiduTranslator) GetTransResult(taskId int) (ret string, err error) {
	theTask, ok := this.tasks[taskId]
	if !ok {
		err = fmt.Errorf("task id not exist:%v", taskId)
		return "", err
	}

	theTask.Wait()
	return theTask.GetResult()
	//return ret, err
}

/*
   params = ["q"=>"words for translate"]
*/
func makeSign(params map[string]string) string {
	//:= fmt.Sprintf("%d", time.Now().Unix())

	longStr := baidu_appid + params["q"] + params["salt"] + baidu_appsecret

	h := md5.New()
	h.Write([]byte(longStr))
	signByte := h.Sum(nil)

	sign := fmt.Sprintf("%x", signByte)
	return sign
}

func (this *BaiduTransTask) Wait() bool {
	return false
}

func (this *BaiduTransTask) GetResult() (string, error) {
	return "", nil
}
