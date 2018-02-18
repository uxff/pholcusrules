package nondotpicset

/*
curl http://www.airtmv.com/
需求： 下载静态网站中的图集
记录图库资源
PICSETNAME,IMG_OF_PICSET


*/

// 基础包
import (
	"fmt"
	"net/http"
	//"strconv"
	"strings"
	//"sync"

	//"golang.org/x/net/html"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	. "github.com/henrylee2cn/pholcus/app/spider" //必需

	helper "github.com/uxff/pholcusrules/consts"
)

func init() {
	config := &helper.AirConfig{
		Name:      "airtmv.com",
		Domain:    "airtmv.com",
		HomePage:  "http://www.airtmv.com/",
		FirstPage: "http://www.airtmv.com/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)

	helper.AIR_CONFIGS[config.Name] = config
	TheAirtmv.Name = config.Name
	TheAirtmv.Description = config.Domain

	TheAirtmv.Register()
}

var TheAirtmv = &Spider{
	//Name:         THE_DOMAIN,
	//Description:  THE_DOMAIN + " no need input",
	Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			entranceUrl := helper.AIR_CONFIGS[ctx.GetName()].FirstPage
			keyIn := ctx.GetKeyin()
			if len(keyIn) > 4 {
				entranceUrl = keyIn
			}

			logs.Log.Warning("start with url:%v", entranceUrl)
			helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot)

			ctx.AddQueue(&request.Request{
				Url:  entranceUrl,
				Rule: "TAGLIST",
				Header: http.Header{
					"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
					"User-Agent": []string{helper.AGENT_PUBLIC},
					"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
				},
			})
		},

		Trunk: map[string]*Rule{
			"HOMEPAGE": {
				ParseFunc: func(ctx *Context) {
					// home page
					// tag list like: https://www.4493.com/star/mihuanmeinv/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
					lis := query.Find(".nav").Find("ul").Find("li") // 不能写 ".thumb a"
					logs.Log.Warning("the nav li =%v", lis.Length())

					lis.Each(func(i int, s *goquery.Selection) {

						url, _ := s.Find("a").Eq(0).Attr("href")
						tagName := s.Find("a").Eq(0).Text()
						tagName = strings.Trim(tagName, " \t")

						if len(url) == 0 {
							return
						}

						//logs.Log.Warning("find a picset list(%v):%v", tagName, url)
						url = helper.FixUrl(url, ctx.GetUrl())

						imgThumb, _ := s.Find("img").Attr("src")
						imgThumb = helper.FixUrl(imgThumb, ctx.GetUrl())

						// download in disk , save to local
						helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName)

						//SaveConfig()

						//logs.Log.Warning("extract tag url, img=%v, %v, %v", url, img, tagName)

						// cookie
						cookies := ""
						cookie := ctx.Response.Cookies()
						for _, c := range cookie {
							cookies += c.Name + "=" + c.Value + "; "
						}

						logs.Log.Warning("will request taglist->picsetlist: %v", url)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: "PICSETLIST",
								Temp: map[string]interface{}{"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName, "TAGNAME": tagName},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									"Cookie":     []string{cookies},
									"User-Agent": []string{helper.AGENT_PUBLIC},
									"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
									//"Upgrade-Insecure-Requests": []string{"1"},
									//"Cache-Control":             []string{"no-cache"},
								},
								DownloaderID: 0,
							},
						)

						helper.DownloadObject(imgThumb, helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+tagName, "thumb")

					})
				},
			},

			"TAGLIST": {
				ParseFunc: func(ctx *Context) {
					// tag list like: https://www.4493.com/star/mihuanmeinv/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
					lis := query.Find(".nav").Find("ul").Find("li") // 不能写 ".thumb a"
					logs.Log.Warning("the nav li =%v", lis.Length())

					lis.Each(func(i int, s *goquery.Selection) {

						url, _ := s.Find("a").Eq(0).Attr("href")
						tagName := s.Find("a").Eq(0).Text()
						tagName = strings.Trim(tagName, " \t")

						if len(url) == 0 {
							return
						}

						//logs.Log.Warning("find a picset list(%v):%v", tagName, url)
						url = helper.FixUrl(url, ctx.GetUrl())

						imgThumb, _ := s.Find("img").Attr("src")
						imgThumb = helper.FixUrl(imgThumb, ctx.GetUrl())

						// download in disk , save to local
						helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName)

						//SaveConfig()

						//logs.Log.Warning("extract tag url, img=%v, %v, %v", url, img, tagName)

						// cookie
						cookies := ""
						cookie := ctx.Response.Cookies()
						for _, c := range cookie {
							cookies += c.Name + "=" + c.Value + "; "
						}

						logs.Log.Warning("will request taglist->picsetlist: %v", url)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: "PICSETLIST",
								Temp: map[string]interface{}{"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName, "TAGNAME": tagName},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									"Cookie":     []string{cookies},
									"User-Agent": []string{helper.AGENT_PUBLIC},
									"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
									//"Upgrade-Insecure-Requests": []string{"1"},
									//"Cache-Control":             []string{"no-cache"},
								},
								DownloaderID: 0,
							},
						)

						helper.DownloadObject(imgThumb, helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+tagName, "thumb")

					})
				},
			},

			"PICSETLIST": {
				ItemFields: []string{
					"PICSETNAME",
					"PICSETURL",
					"PICSETTHUMB",
				},
				ParseFunc: func(ctx *Context) {
					// picset list like: https://www.4493.com/star/mihuanmeinv/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					saveDir := ctx.GetTemp("DIR", helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot).(string)

					query := ctx.GetDom()
					lis := query.Find(".piclist").Find("ul").Find("li") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 10 {
							//return
						}

						url, _ := s.Find("a").Eq(0).Attr("href")
						img, _ := s.Find("img").Eq(0).Attr("src")
						picsetName := s.Find("a").Eq(0).Find("span").Text()
						picsetName = strings.Trim(picsetName, " \t")
						if len(picsetName) == 0 {
							picsetName = fmt.Sprintf("%v", i)
						}

						if len(url) == 0 {
							return
						}

						img = helper.FixUrl(img, ctx.GetUrl())
						//logs.Log.Warning("get a set url:%v", url)
						url = helper.FixUrl(url, ctx.GetUrl())

						// record this picset
						ctx.Output(map[int]interface{}{
							0: picsetName,
							1: url,
							2: img,
						})

						// download in disk
						// save to local
						helper.MakeDir(saveDir + "/" + picsetName)

						//logs.Log.Warning("extract picset url, img=%v, %v, %v", url, img, picsetName)

						// cookie
						cookies := ""
						cookie := ctx.Response.Cookies()
						for _, c := range cookie {
							cookies += c.Name + "=" + c.Value + "; "
						}

						//logs.Log.Warning("cookie=%s ", cookies)
						//cookies = "25a6da5acf7fde759f79e8c23ab0dc76d53f8=cGxmUnIwNWQ1T3JvTVRVeE16QTVNamsyTlMwd0xTRXcb; f6d9f67e5f2c5b9dd954e79d40588261f042e31abda5d=bkpHaHAwMU1ZV0U0TlRsbFpqRmlaREF4TTJaa1pXVXlORFZpTlRRd01ETXhNamRqWkRNPQc; 9353eb=1513092965; _ga=GA1.2.1711782959.1513095156; _gid=GA1.2.1692799327.1513095156; 34843e6d0d96c28940bc888267e9b3=ekxwRzExN1FTRVpoVXVuRHlKV3VMOTB1ZERRNU9UUTVPQT09a; 073273c2a4e3c0d936022720d=SzZlRE4xNlZCdk00QUdleWQ5NlVHMWVNaTB3a; 9353e=bm9yZWZ8fGRlZmF1bHR8MXwyfDJ8bm9uZXwwOmpwZ3JhdnVyZS5jb20%3D; __atuvc=3%7C50; __atuvs=5a2fffe7e5348de2002"

						logs.Log.Warning("will request picsetlsit->picset: %v", url)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: "PICSET",
								Temp: map[string]interface{}{"DIR": saveDir + "/" + picsetName + "/", "PICSETNAME": picsetName},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									"Cookie":     []string{cookies},
									"User-Agent": []string{helper.AGENT_PUBLIC},
									"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
									//"Upgrade-Insecure-Requests": []string{"1"},
									//"Cache-Control":             []string{"no-cache"},
								},
								DownloaderID: 0,
							},
						)

						helper.DownloadObject(img, saveDir+"/"+picsetName, "thumb")
					})
				},
			},

			"PICSET": {
				ItemFields: []string{
					"PICSETNAME",
					"SAVEPATH",
					"IMAGEURL",
				},
				ParseFunc: func(ctx *Context) {
					// picset like: https://www.4493.com/gaoqingmeinv/134254/1.htm
					query := ctx.GetDom()

					picsetName := ctx.GetTemp("PICSETNAME", "").(string)
					saveDir := ctx.GetTemp("DIR", helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+picsetName).(string)

					imgUrl, _ := query.Find(".picsboxcenter").Find("img").Eq(0).Attr("src")
					imgUrl = helper.FixUrl(imgUrl, ctx.GetUrl())

					helper.DownloadObject(imgUrl, saveDir, "")
					//logs.Log.Warning("IN %v imgUrl=%v", picsetName, imgUrl)

					pages := query.Find(".page").Find("a")
					pages.Each(func(i int, s *goquery.Selection) {

						if i == 0 || i == pages.Length()-1 {
							return
						}

						nextPageUrl, _ := s.Attr("href")
						nextPageUrl = strings.Trim(nextPageUrl, " \t\r\n")
						if nextPageUrl == "" {
							return
						}

						nextPageUrl = helper.FixUrl(nextPageUrl, ctx.GetUrl())
						ctx.AddQueue(
							&request.Request{
								Url:  nextPageUrl,
								Rule: "PICSET",
								Temp: map[string]interface{}{"DIR": saveDir, "PICSETNAME": picsetName},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									//"Cookie":     []string{cookies},
									"User-Agent": []string{helper.AGENT_PUBLIC},
									"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
									//"Upgrade-Insecure-Requests": []string{"1"},
									//"Cache-Control":             []string{"no-cache"},
								},
								DownloaderID: 0,
							},
						)

					})

				},
			},
		},
	},
}
