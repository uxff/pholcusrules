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
	"sync"
	"time"
)

const (
	STATUS_NONE = iota
	STATUS_DOING
	STATUS_OK
	STATUS_FAIL
)

const (
	BAIDU_API_URL       = "http://api.fanyi.baidu.com/api/trans/vip/translate"
	BAIDU_API_URL_HTTPS = "https://api.fanyi.baidu.com/api/trans/vip/translate"
)

const MAX_TRY_TIMES = 5

type BaiduTransTask struct {
	id        int
	fromLang  string
	toLang    string
	queryStr  string
	retStr    string
	status    int
	failMsg   string
	failTimes int
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

var baiduTransTaskNextId = 0
var baidu_appid string
var baidu_appsecret string

//var taskFinishChan chan *BaiduTransTask

func init() {
	//if this.tasks == nil {
	//	this.tasks = make(map[int]*BaiduTransTask, 0)
	//}
}

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
	if this.tasks == nil {
		this.tasks = make(map[int]*BaiduTransTask, 0)
	}

	theChan := this.AsyncTranslate(str)

	theRes := <-theChan

	//	this.tasks[taskId].Wait(time.Second * 10)
	//	if this.tasks[taskId].status == STATUS_OK {
	//		return this.tasks[taskId].retStr, nil
	//	}

	return theRes.Res, theRes.Err
	//return this.tasks[taskId].retStr, fmt.Errorf(this.tasks[taskId].failMsg)
	//return this.retStr, err

}
func (this *BaiduTranslator) AsyncTranslate(str string) <-chan *TransRes {

	baiduTransTaskNextId++
	task := &BaiduTransTask{id: baiduTransTaskNextId, queryStr: str, status: STATUS_NONE, fromLang: this.fromLang, toLang: this.toLang}

	theChan := make(chan *TransRes, 1)

	go func() {
		mu := &sync.Mutex{}
		mu.Lock()
		defer mu.Unlock()
		this.tasks[baiduTransTaskNextId] = task
		task.Start()

		tranRes := &TransRes{}
		tranRes.Res, tranRes.Err = task.GetResult()

		theChan <- tranRes

	}()

	//this.tasks = append(this.tasks, task)
	//task.Start()
	return theChan
}

func (this *BaiduTranslator) GetTransResult(taskId int) (ret string, err error) {
	theTask, ok := this.tasks[taskId]
	if !ok {
		err = fmt.Errorf("task id not exist:%v", taskId)
		return "", err
	}

	theTask.Wait(time.Second * 10)
	return theTask.GetResult()
	//return ret, err
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

func (this *BaiduTransTask) Start() {
	if this.status == STATUS_NONE {
		this.status = STATUS_DOING
	}

	if this.status == STATUS_OK {
		return
	}

	if this.failTimes > MAX_TRY_TIMES {
		return
	}

	for this.failTimes = 0; this.failTimes < MAX_TRY_TIMES; this.failTimes++ {
		_, err := this.RunTrans()
		if err == nil {
			this.status = STATUS_OK
			break
		}
		this.failMsg = err.Error()
	}
}

func (this *BaiduTransTask) RunTrans() (string, error) {

	if this.fromLang == "" {
		this.fromLang = "auto"
	}
	if this.toLang == "" {
		this.fromLang = "auto"
	}

	//this.queryStr

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

// you should use channel to wait, you should use context to timeout
func (this *BaiduTransTask) Wait(timeout time.Duration) bool {
	tChan := time.After(timeout)

	for {
		select {
		case <-tChan:
			return false
		default:
			if this.status == STATUS_FAIL || this.status == STATUS_OK {
				return true
			}
			time.Sleep(1 * time.Second)
		}
	}
	return false
}

func (this *BaiduTransTask) GetStatus() int {
	return this.status
}

func (this *BaiduTransTask) GetResult() (string, error) {

	switch this.status {
	case STATUS_NONE:
		if this.failTimes < MAX_TRY_TIMES {
			go this.Start()
		}
		return "", fmt.Errorf("its doing, not done yet")
	case STATUS_DOING:
		return "", fmt.Errorf("its doing, not done yet")
	case STATUS_OK:
		return this.retStr, nil
	case STATUS_FAIL:
		return "", fmt.Errorf("its failed %d times, now trying again", this.failTimes)
	}

	return "", nil
}
