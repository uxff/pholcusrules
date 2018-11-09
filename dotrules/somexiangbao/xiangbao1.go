package somexiangbao

// 基础包
import (
	"encoding/base64"
	"net/http" //设置http.Header
	"net/url"
	"strings"
	// "log"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用
)

//思路
// 抓取资讯 抓取发布者 自动注册 自动发布到对应板块
// done

const (
	HOME_URL = "http://www.baigou.net"

	AGENT_PUBLIC  = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"
	AGENT_WX      = "Mozilla/5.0 (Linux; Android 6.0; 1503-M02 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/37.0.0.0 Mobile MQQBrowser/6.2 TBS/036558 Safari/537.36 MicroMessenger/6.5.7.1041 NetType/WIFI Language/zh_CN"
	AGENT_WX_3G   = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12A365 MicroMessenger/6.0 NetType/3G+"
	AGENT_WX_WIFI = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12A365 MicroMessenger/6.0 NetType/WIFI"
	AGENT_WX_IOS  = "Mozilla/5.0 (iPhone; CPU iPhone OS 10_2_1 like Mac OS X) AppleWebKit/602.4.6 (KHTML, like Gecko) Mobile/14D27 MicroMessenger/6.5.6 NetType/4G Language/zh_CN"
	AGENT_WX_AND  = "Mozilla/5.0 (Linux; Android 5.1; OPPO R9tm Build/LMY47I; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.49 Mobile MQQBrowser/6.2 TBS/043220 Safari/537.36 MicroMessenger/6.5.7.1041 NetType/4G Language/zh_CN"
)

func init() {
	Xiangbao1.Register()
}

var Xiangbao1 = &Spider{
	Name:        "Xiangbao1",
	Description: "http://www.baigou.net",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  HOME_URL,
				Rule: "HOMEPAGE",
				Header: http.Header{
					"User-Agent": []string{AGENT_PUBLIC},
				},
			})
		},

		Trunk: map[string]*Rule{
			"HOMEPAGE": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".bx-newl").Find("li").Each(func(j int, s *goquery.Selection) {
						cate := s.Find("a").Eq(0)

						cateName := cate.Text()

						th := s.Find("a").Eq(1)
						thLink, ok := th.Attr("href")
						thName := th.Text()

						if ok {
							ctx.AddQueue(&request.Request{
								Url:    thLink,
								Rule:   "TH",
								Header: http.Header{"Referer": []string{HOME_URL}, "User-Agent": []string{AGENT_PUBLIC}},
								Temp: request.Temp{
									"cate": cateName,
									"th":   thName,
								},
							})
						}

					})

				},
			},
			"TH": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"Title",
					"User",
					"UserAuthed",
					"Contactor",
					"Content",
					"Thumb",
					"Time",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					// pubtime
					pubtime := query.Find(".ben-info").Find(".cen").Text()
					pubtimeIdx := strings.Index(pubtime, "发布时间：")
					if pubtimeIdx >= 0 && len(pubtime) > len("发布时间：")+19 {
						pubtime = pubtime[pubtimeIdx+len("发布时间：") : pubtimeIdx+len("发布时间：")+19]
					}

					// 获取内容
					content, _ := query.Find(".ben-sx").Html()

					// Title
					title := ctx.GetTemp("th", "")
					// Author
					contactor := query.Find(".ben-zone")
					//contactorHtml, _ := contactor.Html()
					contactorNoUrl, _ := contactor.Find("img").Attr("src")
					contactorNo := ""
					if uo, err := url.Parse(contactorNoUrl); err == nil {
						contactorNo = uo.Query().Get("vid")
						if b, err := base64.URLEncoding.DecodeString(contactorNo); err == nil {
							contactorNo = string(b)
						}
					}

					userNameText := query.Find(".bx-ben-r-a").Text()
					userAuthed := ""
					userAuthedImgs := query.Find(".bx-ben-r-a").Find("img")
					userAuthedImgs.Each(func(i int, s *goquery.Selection) {
						if u, ok := s.Attr("title"); ok == true {
							userAuthed += u + ";"
						}
					})

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: title,
						1: userNameText,
						2: userAuthed,
						3: contactorNo,
						4: content,
						5: "",
						6: pubtime,
					})
				},
			},
		},
	},
}
