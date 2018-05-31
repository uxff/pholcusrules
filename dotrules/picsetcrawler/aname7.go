package picsetcrawler

/*
curl http://sexygirlcity.com/index.php?gal=0_1
需求：
 - 下载静态网站中的图集
 - 记录图库资源
dev:undown
download:undown


*/

// 基础包
import (
	"fmt"
	"net/http"
	"net/url"
	//"strconv"

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
		Name:      "ANAME7",
		Domain:    "sexygirlcity.com",
		HomePage:  "http://www.sexygirlcity.com/index.php",
		FirstPage: "http://www.sexygirlcity.com/index.php?gal=0_1",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)
	// save: tagname/picsetname/pics*.jpg

	helper.AIR_CONFIGS[config.Name] = config
	Aname7.Name = config.Name
	Aname7.Description = config.Domain

	Aname7.Register()
}

var Aname7 = &Spider{
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
				Rule: "HOMEPAGE",
				Header: http.Header{
					//"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
					"User-Agent": []string{helper.AGENT_PUBLIC},
					//"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
				},
				Temp: map[string]interface{}{
					"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot,
				},
			})
		},

		Trunk: map[string]*Rule{
			"HOMEPAGE": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					// first page picset
					ctx.AddQueue(&request.Request{
						Url:  helper.AIR_CONFIGS[ctx.GetName()].FirstPage,
						Rule: "PICSETLIST",
						Header: http.Header{
							//"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
							"User-Agent": []string{helper.AGENT_PUBLIC},
							//"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
						},
						Temp: map[string]interface{}{
							"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot,
						},
					})

					// page list, without first page
					query.Find(".pager-next").Each(func(i int, s *goquery.Selection) {
						pageUrl, _ := s.Attr("href")
						pageUrl = helper.FixUrl(pageUrl, ctx.GetUrl())
						ctx.AddQueue(&request.Request{
							Url:  pageUrl,
							Rule: "PICSETLIST",
							Header: http.Header{
								//"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
								"User-Agent": []string{helper.AGENT_PUBLIC},
								//"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
							},
							Temp: map[string]interface{}{
								"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot,
							},
						})
					})
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
					trs := query.Find(".pager").Parent().Find("table").Find("tr")
					logs.Log.Warning("the trs.lengh=%d", trs.Length())

					lis := query.Find(".pager").Parent().Find("table").Find("tr").Find("a")
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 3 {
							//return
						}

						targetUrl, _ := s.Attr("href")
						img, _ := s.Find("img").Eq(0).Attr("src")
						picsetName, _ := s.Find("img").Eq(0).Attr("alt")

						if len(targetUrl) == 0 || targetUrl[0] == '#' {
							return
						}

						img = helper.FixUrl(img, ctx.GetUrl())
						//logs.Log.Warning("get a set url:%v", targetUrl)
						targetUrl = helper.FixUrl(targetUrl, ctx.GetUrl())

						urlParsed, _ := url.Parse(targetUrl)
						picsetId := urlParsed.Query().Get("viewg")
						picsetName = picsetId + "-" + picsetName

						// download in disk, save to local
						//helper.MakeDir(saveDir + "/" + picsetName)

						logs.Log.Warning("will request picsetlsit->picset: %v picsetName=%s", targetUrl, picsetName)

						// queue request the picset detail
						ctx.AddQueue(
							&request.Request{
								Url:  targetUrl,
								Rule: "PICSET",
								Temp: map[string]interface{}{
									"DIR":        saveDir + "/",
									"THUMB_URL":  img,
									"PICSETNAME": picsetName,
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

					// other

				},
			},

			"PICSET": {
				ItemFields: []string{
					"PICSETNAME",
					"SAVEPATH",
					"IMAGEURL",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//title := query.Find("title").Text()
					picsetName := ctx.GetTemp("PICSETNAME", "").(string)

					saveDir := helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + picsetName
					helper.MakeDir(saveDir)
					helper.MakeDir(saveDir + "/thumbs")

					thumbUrl := ctx.GetTemp("THUMB_URL", "").(string)
					helper.DownloadObject(thumbUrl, saveDir, "thumb")

					writeConfig := map[string]string{
						"title":   picsetName,
						"url":     ctx.GetUrl(),
						"tags":    "",
						"pubdate": "",
					}

					helper.WritePicsetConfig(writeConfig, saveDir)

					query.Find(".links2").Each(func(imgi int, imga *goquery.Selection) {

						largePic, _ := imga.Attr("href")
						littlePic, _ := imga.Find("img").Attr("src")

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
