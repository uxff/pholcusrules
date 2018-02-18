package nondotpicset

/*
curl http://www.55156.com
需求： 下载静态网站中的图集
记录图库资源
PICSETNAME,IMG_OF_PICSET
TODO::
    - auto recognize tag url or picset url

*/

// 基础包
import (
	"fmt"
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
		Name:      "55156.com",
		Domain:    "www.55156.com",
		HomePage:  "http://www.55156.com/",
		FirstPage: "http://www.55156.com/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)

	helper.AIR_CONFIGS[config.Name] = config
	The55156.Name = config.Name
	The55156.Description = "[Auto Page] " + config.Domain

	The55156.Register()
}

var The55156 = &Spider{
	//Name:         "55156.com",
	//Description:  "[Auto Page] [www.55156.com]",
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
				Rule: "HOMEPAGE",
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
					// home page like: http://www.55156.com/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					lis := query.Find(".nav").Find("ul").Find("li") // 不能写 ".thumb a"
					logs.Log.Warning("the nav li =%v:%v", lis.Length(), lis.Text())
					lis.Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							// 0==homepage
							return
						}

						url, _ := s.Find("a").Eq(0).Attr("href")
						tagName := s.Find("a").Eq(0).Text()
						tagName = strings.Trim(tagName, " \t")

						//logs.Log.Warning("find a picset list(%v):%v", tagName, url)

						if len(url) > 0 {
							//logs.Log.Warning("get a set url:%v", url)
							url = helper.FixUrl(url, ctx.GetUrl())

							// download in disk
							// save to local
							helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName)

							//logs.Log.Warning("extract tag url, img=%v, %v, %v", url, img, tagName)

							logs.Log.Warning("will request home->picsetlist: %v", url)

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

							//helper.DownloadObject(img, helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+tagName, "thumb")
						}

					})
				},
			},

			"TAGLIST": {
				ItemFields: []string{
					"TAGNAME",
					"TAGURL",
					"TAGTHUMB",
				},
				ParseFunc: func(ctx *Context) {
					// like: http://www.55156.com/meinvtupian/bijinimeinv.html
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					lis := query.Find(".labelbox_pic").Find("ul").Find("li") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
						// limit 10
						if i > 10 {
							//return
						}

						url, _ := s.Find("a").Eq(0).Attr("href")
						img, _ := s.Find("a img").Attr("src")

						tagName := s.Find(".pic_tit").Text()
						tagName = strings.Trim(tagName, " \t")

						if len(tagName) == 0 {
							tagName = fmt.Sprintf("tag%v", i)
						}

						img = helper.FixUrl(img, ctx.GetUrl())

						if len(url) > 0 {
							//logs.Log.Warning("get a set url:%v", url)
							url = helper.FixUrl(url, ctx.GetUrl())

							// record this picset
							ctx.Output(map[int]interface{}{
								0: tagName,
								1: url,
								2: img,
							})

							// download in disk // save to local
							helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName)

							logs.Log.Warning("will request tag->picsetlist: %v", url)

							// queue request the picset detail
							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: "PICSETLIST",
									Temp: map[string]interface{}{"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + tagName + "/", "TAGNAME": tagName},
									Header: http.Header{
										//"Accept-Language":           []string{"zh-CN,zh"},
										"Cookie":     []string{cookies},
										"User-Agent": []string{helper.AGENT_PUBLIC},
										"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
									},
									DownloaderID: 0,
								},
							)

							helper.DownloadObject(img, helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+tagName, "thumb")
						}

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

					// like: http://www.55156.com/weimeiyijing/fengjingtupian/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					saveDir := ctx.GetTemp("DIR", helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot).(string)

					query := ctx.GetDom()
					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					lis := query.Find("#imgList").Find("ul").Find("li") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 10 {
							//return
						}

						url, _ := s.Find("a").Eq(0).Attr("href")
						img, _ := s.Find("img").Eq(0).Attr("src")
						picsetName, _ := s.Find("a").Eq(0).Attr("title")
						picsetName = strings.Trim(picsetName, " \t")
						if len(picsetName) == 0 {
							picsetName = fmt.Sprintf("%v", i)
						}

						img = helper.FixUrl(img, ctx.GetUrl())

						if len(url) > 0 {
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
						}

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
					// picset like: http://www.55156.com/weimeiyijing/fengjingtupian/206208.html
					query := ctx.GetDom()

					picsetName := ctx.GetTemp("PICSETNAME", "").(string)
					saveDir := ctx.GetTemp("DIR", helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot+picsetName).(string)

					imgUrl, _ := query.Find("#picBody").Find("a img").Attr("src")
					imgUrl = helper.FixUrl(imgUrl, ctx.GetUrl())

					helper.DownloadObject(imgUrl, saveDir, "")
					//logs.Log.Warning("IN %v imgUrl=%v", picsetName, imgUrl)

					query.Find(".pages").Find("ul").Find("li").Each(func(i int, s *goquery.Selection) {
						nextPageUrl, _ := s.Find("a").Attr("href")
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
