package kuaidaili

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

	// 字符串处理包
	//"regexp"
	"strconv"
	"strings"
	// 其他包
	"fmt"
	// "math"
	// "time"
)

const (
	HOME_URL       = "https://www.kuaidaili.com"
	FIRST_PAGE_URL = "https://www.kuaidaili.com/free/inha/1/"
	PAGE_URL_FMT   = "https://www.kuaidaili.com/free/inha/%d/"
	PUBLIC_AGENT   = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"
)

func init() {
	Kuaidaili.Register()
}

var Kuaidaili = &Spider{
	Name:        "Kuaidaili",
	Description: "Kuaidaili KEYIN: [cookie], set pause time: 2s or bigger",
	// Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//cookie := "channelid=0; sid=1516785851105118; Hm_lvt_7ed65b1cc4b810e9fd37959c9bb51b31=1516785861; Hm_lpvt_7ed65b1cc4b810e9fd37959c9bb51b31=1516785874; _ga=GA1.2.2140228708.1516785861; _gid=GA1.2.1413734529.1516785861"
			cookie := "channelid=0; sid=1566311104277914; Hm_lvt_7ed65b1cc4b810e9fd37959c9bb51b31=1566311907; _ga=GA1.2.1804874949.1566311907; _gid=GA1.2.1007360879.1566311907; Hm_lpvt_7ed65b1cc4b810e9fd37959c9bb51b31=1566311918"
			keyIn := strings.Trim(ctx.GetKeyin(), "\r\n\t ")

			if keyIn != "" {
				cookie = keyIn
			}

			ctx.AddQueue(&request.Request{
				Url:  FIRST_PAGE_URL,
				Rule: "FREE_INHA",
				Header: http.Header{
					"User-Agent":                []string{PUBLIC_AGENT},
					"Referer":                   []string{HOME_URL},
					"Accept-Language":           []string{"zh-CN,zh;q=0.8"},
					"Cookie":                    []string{cookie},
					"Upgrade-Insecure-Requests": []string{"1"},
				},
				Temp: map[string]interface{}{
					"page": 1,
				},
			})
		},

		Trunk: map[string]*Rule{
			"FREE_INHA": {
				ItemFields: []string{
					"HTTP_PROXY",
					"PROTO",
					"IP",
					"PORT",
				},
				ParseFunc: func(ctx *Context) {
					cookie := ctx.GetCookie()
					curPageNo := ctx.GetTemp("page", 1).(int)
					query := ctx.GetDom()
					query.Find("#list").Find("table").Find("tbody").Find("tr").Each(func(j int, s *goquery.Selection) {
						ip := ""
						port := ""
						proto := ""

						s.Find("td").Each(func(ei int, es *goquery.Selection) {

							dataTitleVal, dataTitleExist := es.Attr("data-title")

							if dataTitleExist {
								switch strings.ToUpper(dataTitleVal) {
								case "IP":
									ip = es.Text()
								case "PORT":
									port = es.Text()
								case "类型":
									proto = es.Text()
								}
							}

							//logs.Log.Warning("find a td: %v", es.Text())

						})
						ctx.Output(map[int]interface{}{
							0: fmt.Sprintf("%s://%v:%v", strings.ToLower(proto), ip, port),
							1: proto,
							2: ip,
							3: port,
						})

						//logs.Log.Warning("tr ok:%v", elem.Text())
					})

					// page to next

					allPageLis := query.Find("#listnav").Find("ul").Find("a")
					allPageLiNo := allPageLis.Length()
					lastPageNo := allPageLis.Eq(allPageLiNo - 1).Text()

					lastPageNoInt, _ := strconv.Atoi(lastPageNo)

					logs.Log.Warning("find cur page:%v lastPageNo=%v", curPageNo, lastPageNoInt)

					nextPageNoInt := curPageNo + 1

					if nextPageNoInt <= lastPageNoInt {

						theNextPageUrl := fmt.Sprintf(PAGE_URL_FMT, nextPageNoInt)
						ctx.AddQueue(&request.Request{
							Url:  theNextPageUrl,
							Rule: "FREE_INHA",
							Header: http.Header{
								"User-Agent":                []string{PUBLIC_AGENT},
								"Referer":                   []string{HOME_URL},
								"Accept-Language":           []string{"zh-CN,zh;q=0.8"},
								"Cookie":                    []string{cookie},
								"Upgrade-Insecure-Requests": []string{"1"},
							},
							Temp: map[string]interface{}{
								"page": nextPageNoInt,
							},
						})
					}

				},
			},
		},
	},
}
