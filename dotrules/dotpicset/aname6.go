package picsetcrawler

/*
curl http://www.asianamateurpussy.com/ with redirect url=xxx
需求： 下载静态网站中的图集
记录图库资源
PICSETNAME,IMG_OF_PICSET
dev:unknown
download:unknown
tags: bad albumn structure

*/

// 基础包
import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

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
		Name:      "ANAME6",
		Domain:    "asianamateurpussy.com",
		HomePage:  "http://www.asianamateurpussy.com/",
		FirstPage: "http://www.asianamateurpussy.com/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)
	// save: tagname/picsetname/pics*.jpg

	helper.AIR_CONFIGS[config.Name] = config
	Aname6.Name = config.Name
	Aname6.Description = config.Domain

	Aname6.Register()
}

var Aname6 = &Spider{
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
					// has base64 in url, you need decode:
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
					lis := query.Find(".thumbs").Find("a")
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 3 {
							//return
						}

						targetUrl, _ := s.Attr("href")
						img, _ := s.Find("img").Eq(0).Attr("src")

						if len(targetUrl) == 0 || targetUrl[0] == '#' {
							return
						}

						img = helper.FixUrl(img, ctx.GetUrl())
						//logs.Log.Warning("get a set url:%v", targetUrl)
						targetUrl = helper.FixUrl(targetUrl, ctx.GetUrl())

						// base64 decode
						urlParsed, _ := url.Parse(targetUrl)
						q := urlParsed.Query()
						targetUrlEnc := q.Get("url")

						var targetUrlOri []byte = make([]byte, 1024)

						urlLen, _ := base64.NewDecoder(base64.StdEncoding, strings.NewReader(targetUrlEnc)).Read(targetUrlOri)

						targetUrl = string(targetUrlOri[:urlLen])

						// download in disk, save to local
						//helper.MakeDir(saveDir + "/" + picsetName)

						logs.Log.Warning("will request picsetlsit->picset: %v ori=%v", targetUrl, targetUrlEnc)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  targetUrl,
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
								},
								DownloaderID: 0,
							},
						)

						//helper.DownloadObject(img, saveDir+"/"+picsetName, "thumb")
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

					title := query.Find("title").Text()

					saveDir := helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + title
					helper.MakeDir(saveDir)
					helper.MakeDir(saveDir + "/thumbs")

					thumbUrl := ctx.GetTemp("THUMB_URL", "").(string)
					helper.DownloadObject(thumbUrl, saveDir, "thumb")

					writeConfig := map[string]string{
						"title":   title,
						"url":     ctx.GetUrl(),
						"tags":    "",
						"pubdate": "",
					}

					helper.WritePicsetConfig(writeConfig, saveDir)

					query.Find(".thumbs").Find("tr").Each(func(tri int, trs *goquery.Selection) {
						colspan, _ := trs.Attr("colspan")
						if colspan != "" {
							return
						}

						trs.Find("a").Each(func(tdi int, tds *goquery.Selection) {

							largePic, _ := tds.Attr("href")
							littlePic, _ := tds.Find("img").Attr("src")

							largePic = helper.FixUrl(largePic, ctx.GetUrl())
							littlePic = helper.FixUrl(littlePic, ctx.GetUrl())

							helper.DownloadObject(largePic, saveDir, "")
							helper.DownloadObject(littlePic, saveDir+"/thumbs", "")
						})

					})

				},
			},
		},
	},
}
