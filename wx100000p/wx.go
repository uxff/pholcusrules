package wx100000p

// 基础包
import (
	// "log"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"         //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	//"regexp"
	//"strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

const (
	HOME_URL     = "http://100000p.com"
	PUBLIC_AGENT = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"
)

func init() {
	Wx100000p.Register()
}

var Wx100000p = &Spider{
	Name:        "Wx100000p",
	Description: "Wx100000p",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  HOME_URL,
				Rule: "TIMELINE",
				Header: http.Header{
					"User-Agent": []string{PUBLIC_AGENT},
				},
			})
		},

		Trunk: map[string]*Rule{
			"TIMELINE": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".site-content").Find("article").Each(func(j int, s *goquery.Selection) {
						a := s.Find("a")
						if url, ok := a.Attr("href"); ok {
							// log.Print(url)
							p := a.Find("p")

							ctx.AddQueue(&request.Request{Url: HOME_URL + url, Rule: "DETAIL", Header: http.Header{"Referer": []string{HOME_URL}, "User-Agent": []string{PUBLIC_AGENT}}, Temp: request.Temp{"abstract": p.Text()}})
						}
					})

				},
			},
			"DETAIL": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"Title",
					"Author",
					"Thumb",
					"Time",
					"Abstract",
					"Keywords",
					"Content",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					// 获取内容
					content, _ := query.Find(".entry-content").Html()

					// 过滤标签
					//re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					//contentText := re.ReplaceAllString(content, "")
					// 内容中如果图片不是

					// Title
					title := query.Find(".entry-title").Text()
					// Author

					// Addresses & Address
					addresses, ok := query.Find(".post-cover-title-image").Attr("style")
					if ok {
						attrArr := strings.Split(addresses, ";")
						for _, subAttr := range attrArr {
							if subAttr[:20] == "background-image:url" {
								addresses = subAttr[20:]
								addresses = strings.Trim(addresses, " ()'\"")
								break
							}
						}

					}

					// Time
					pubtime := query.Find(".entry-date").Text()

					// Abstract
					author := ""
					abstract := ctx.GetTemp("abstract", "").(string)

					// Keywords
					keywords := ""

					query.Find(".tag").Each(func(ti int, q *goquery.Selection) {
						keywords = keywords + "," + q.Text()
					})

					// 输出到

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: title,
						1: author,
						2: addresses,
						3: pubtime,
						4: abstract,
						5: keywords,
						6: content,
					})
				},
			},
		},
	},
}
