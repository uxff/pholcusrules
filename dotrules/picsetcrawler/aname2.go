package picsetcrawler

/*
curl http://japanfuckpics.com/
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
		Name:      "ANAME2",
		Domain:    "Auto japanfuckpics.com",
		HomePage:  "http://japanfuckpics.com/",
		FirstPage: "http://japanfuckpics.com/latest/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)
	// save: tagname/picsetname/pics*.jpg

	helper.AIR_CONFIGS[config.Name] = config
	Aname2.Name = config.Name
	Aname2.Description = config.Domain

	Aname2.Register()
}

var Aname2 = &Spider{
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
				Rule: "PICSETLIST",
				Header: http.Header{
					"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
					"User-Agent": []string{helper.AGENT_PUBLIC},
					"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
				},
				Temp: map[string]interface{}{
					"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot,
				},
			})
		},

		Trunk: map[string]*Rule{
			"HOMEPAGE": {
				ParseFunc: func(ctx *Context) {
				},
			},

			"TAGLIST": {
				ParseFunc: func(ctx *Context) {
				},
			},

			"PICSETLIST": {
				ParseFunc: func(ctx *Context) {
					// picset list like:
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					saveDir := ctx.GetTemp("DIR", helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot).(string)

					query := ctx.GetDom()

					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					// <div id="id_container" class="thumbs container xx">
					lis := query.Find("#th_container").Find(".thumb")
					lis.Each(func(i int, s *goquery.Selection) {
						url, _ := s.Find("a").Eq(0).Attr("href")
						img, _ := s.Find("img").Eq(0).Attr("src")

						if len(url) == 0 || url[0] == '#' {
							return
						}

						img = helper.FixUrl(img, ctx.GetUrl())
						//logs.Log.Warning("get a set url:%v", url)
						url = helper.FixUrl(url, ctx.GetUrl())

						// download in disk, save to local
						//helper.MakeDir(saveDir + "/" + picsetName)

						//logs.Log.Warning("extract picset url, img=%v, %v, %v", url, img, picsetName)

						logs.Log.Warning("will request picsetlsit->picset: %v", url)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: "PICSET",
								Temp: map[string]interface{}{
									"DIR":       saveDir + "/",
									"THUMB_URL": img,
									//"PICSETNAME": picsetName,
								},
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

						//helper.DownloadObject(img, saveDir+"/"+picsetName, "thumb")
					})

					// page next
					query.Find(".pager").Find("a").Each(func(pi int, s *goquery.Selection) {
						nextUrl, _ := s.Attr("href")
						if nextUrl == "" || nextUrl[0] == '#' {
							return
						}

						nextUrl = helper.FixUrl(nextUrl, ctx.GetUrl())

						ctx.AddQueue(
							&request.Request{
								Url:  nextUrl,
								Rule: "PICSETLIST",
								Temp: map[string]interface{}{
									//"DIR": saveDir + "/",
								},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									"Cookie":     []string{cookies},
									"User-Agent": []string{helper.AGENT_PUBLIC},
									"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
								},
								DownloaderID: 0,
							},
						)

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

					urlArr := strings.Split(ctx.GetUrl(), helper.AIR_CONFIGS[ctx.GetName()].HomePage)
					picsetName := ""
					title := ""
					if len(urlArr) > 1 {
						picsetName = urlArr[1]
						picsetName = strings.Trim(picsetName, "/")
					} else {
						picsetName := query.Find("#header").Find("h1").Text()
						picsetName = strings.Trim(picsetName, " \t")
						title = picsetName
					}

					saveDir := helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + picsetName
					helper.MakeDir(saveDir)
					helper.MakeDir(saveDir + "/thumbs")

					thumbUrl := ctx.GetTemp("THUMB_URL", "").(string)
					helper.DownloadObject(thumbUrl, saveDir, "thumb")

					if title == "" {
						title = picsetName
					}

					writeConfig := map[string]string{
						"title":   title,
						"url":     ctx.GetUrl(),
						"tags":    "",
						"pubdate": "",
					}

					helper.WritePicsetConfig(writeConfig, saveDir)

					query.Find(".picthumbs").Find("a").Each(func(gi int, s *goquery.Selection) {
						largePic, _ := s.Attr("href")
						littlePic, _ := s.Find("img").Attr("src")

						largePic = helper.FixUrl(largePic, ctx.GetUrl())
						littlePic = helper.FixUrl(littlePic, ctx.GetUrl())

						helper.DownloadObject(largePic, saveDir, "")
						helper.DownloadObject(littlePic, saveDir+"/thumbs", "")
					})

				},
			},
		},
	},
}
