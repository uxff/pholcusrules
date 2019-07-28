package qichacha

// 基础包
import (
	// "log"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	"strings"

	// "math"
	// "time"
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	HOME_URL       = "https://www.qichacha.com"
	FIRST_PAGE_URL = "https://www.qichacha.com/search?key=%E5%8C%97%E4%BA%AC%E5%B0%8F%E6%A1%94%E7%A7%91%E6%8A%80%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8"
	PAGE_URL_FMT   = "https://www.qichacha.com/search?key=%s"
	PUBLIC_AGENT   = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"
)

func init() {
	Qichacha.Register()
}

var Qichacha = &Spider{
	Name:        "Qichacha",
	Description: "Qichacha KEYIN: [cookie], set pause time: 2s or bigger",
	// Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			cookie := "acw_tc=2ff6129915642819433274164e3eff1cefe53f7a23ba05e26a95ab2c5e; QCCSESSID=jcd9fa3m0ts0511rt4b6glehe0; UM_distinctid=16c36784d711bf-0980ea6494f30c-3f385c06-100200-16c36784d7231c; CNZZDATA1254842228=1891776064-1564280713-https%253A%252F%252Fwww.google.com%252F%7C1564280713; _uab_collina=156428195163611994688374; zg_did=%7B%22did%22%3A%20%2216c36785c133dd-0cd4190bc5fa52-3f385c06-100200-16c36785c146f8%22%7D; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=%7B%22sid%22%3A%201564281953304%2C%22updated%22%3A%201564281953309%2C%22info%22%3A%201564281953307%2C%22superProperty%22%3A%20%22%7B%7D%22%2C%22platform%22%3A%20%22%7B%7D%22%2C%22utm%22%3A%20%22%7B%7D%22%2C%22referrerDomain%22%3A%20%22www.google.com%22%7D; Hm_lvt_3456bee468c83cc63fb5147f119f1075=1564281953; Hm_lpvt_3456bee468c83cc63fb5147f119f1075=1564281953"
			keyIn := strings.Trim(ctx.GetKeyin(), "\r\n\t ")

			if keyIn != "" {
				cookie = keyIn
			}

			comListFile := keyIn

			//
			RangeComListFile(comListFile, func(line string, lineNo int) {
				if lineNo > 5 {
					//return
				}

				comName := strings.Trim(line, " \t\r\n")

				logs.Log.Debug(" line(%d)== %s ", lineNo, comName)
				ctx.AddQueue(&request.Request{
					Url:  fmt.Sprintf(PAGE_URL_FMT+"&__ci=%d", comName, lineNo),
					Rule: "SEARCHLIST",
					Header: http.Header{
						"User-Agent":                []string{PUBLIC_AGENT},
						"Referer":                   []string{HOME_URL},
						"Accept-Language":           []string{"zh-CN,zh;q=0.8"},
						"Cookie":                    []string{cookie},
						"Upgrade-Insecure-Requests": []string{"1"},
					},
					Temp: map[string]interface{}{
						"comIdx":  lineNo,
						"comName": comName,
					},
				})
			})

		},

		Trunk: map[string]*Rule{
			"SEARCHLIST": {
				ParseFunc: func(ctx *Context) {
					cookie := ctx.GetCookie()
					//curPageNo := ctx.GetTemp("page", 1).(int)
					query := ctx.GetDom()
					query.Find("#search-result").Find(".ma_h1").Each(func(j int, s *goquery.Selection) {

						theNextUrl, exist := s.Attr("href")
						if !exist {
							return
						}

						if len(theNextUrl) > 4 && theNextUrl[:4] != "http" {
							theNextUrl = HOME_URL + theNextUrl
						}

						ctx.AddQueue(&request.Request{
							Url:  theNextUrl,
							Rule: "SEARCHITEM",
							Header: http.Header{
								"User-Agent":                []string{PUBLIC_AGENT},
								"Referer":                   []string{HOME_URL},
								"Accept-Language":           []string{"zh-CN,zh;q=0.8"},
								"Cookie":                    []string{cookie},
								"Upgrade-Insecure-Requests": []string{"1"},
							},
						})
					})

				},
			},
			"SEARCHITEM": {
				ItemFields: []string{
					"公司名",
					"简介",
					"主营业务",
				},
				ParseFunc: func(ctx *Context) {
					//cookie := ctx.GetCookie()
					query := ctx.GetDom()
					descDom := query.Find("#jianjieModal").Find(".modal-body").Find("pre")
					descText := descDom.Text()

					titleDom := query.Find("#company-top").Find(".content").Find("h1")
					titleText := titleDom.Text()

					infoSectionDom := query.Find("#Cominfo")
					bizIntro := ""

					infoSectionDom.Find("tr").Each(func(j int, s *goquery.Selection) {
						tableColName := s.Find("tb").Text()
						if tableColName == "经营范围" {
							logs.Log.Debug("find 经营范围:%s", titleText)
							bizIntro = s.Find("td").Eq(1).Text()
						}
					})

					ctx.Output(map[int]interface{}{
						0: titleText,
						1: descText,
						2: bizIntro,
					})
				},
			},
		},
	},
}

func RangeComListFile(comListFile string, lineHandler func(line string, lineNo int)) error {
	fhandle, err := os.Open(comListFile)
	if err != nil {
		logs.Log.Error("open comListFile(%s) error:%v", comListFile, err)
		return err
	}

	defer fhandle.Close()

	br := bufio.NewReader(fhandle)

	lineNo := 0
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineNo++

		lineHandler(string(line), lineNo)
	}

	return nil
}
