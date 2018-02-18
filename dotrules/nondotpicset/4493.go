package nondotpicset

/*
curl http://www.4493.com/xilie.html
需求： 下载静态网站中的图集
记录图库资源
PICSETNAME,IMG_OF_PICSET


*/

// 基础包
import (
	"fmt"
	//"io/ioutil"
	"net/http"
	//"strconv"
	"strings"

	//"golang.org/x/net/html"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	. "github.com/henrylee2cn/pholcus/app/spider" //必需

	helper "github.com/uxff/pholcusrules/consts"
)

func init() {
	config := &helper.AirConfig{
		Name:      "4493.com",
		Domain:    "www.4493.com",
		HomePage:  "http://www.4493.com/",
		FirstPage: "http://www.4493.com/xilie.html",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)

	helper.AIR_CONFIGS[config.Name] = config
	The4493.Name = config.Name
	The4493.Description = "[Auto Page] " + config.Domain

	The4493.Register()
}

var The4493 = &Spider{
	//Name:         "4493.com",
	//Description:  "[Auto Page] [www.4493.com]",
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
				},
			},

			"TAGLIST": {
				ItemFields: []string{
					"TAGNAME",
					"TAGURL",
					"TAGTHUMB",
				},
				ParseFunc: func(ctx *Context) {
					// tag list like: https://www.4493.com/star/mihuanmeinv/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					lis := query.Find(".fx_new").Find("ul").Find("li") // 不能写 ".thumb a"
					logs.Log.Warning("the nav li len=%v", lis.Length())

					lis.Each(func(i int, s *goquery.Selection) {
						if i > 10 {
							return
						}

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
					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					lis := query.Find(".piclist").Find("ul").Find("li") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 10 {
							return
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

						// download in disk // save to local
						helper.MakeDir(saveDir + "/" + picsetName)

						//logs.Log.Warning("extract picset url, img=%v, %v, %v", url, img, picsetName)

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
						if i > 1 {
							// only one page
							return
						}

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
