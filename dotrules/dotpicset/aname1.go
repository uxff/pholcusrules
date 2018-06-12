package picsetcrawler

/*
curl http://highasianporn.com
需求： 下载静态网站中的图集
记录图库资源
PICSETNAME,IMG_OF_PICSET
dev:done
download:unknown

*/

// 基础包
import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"net/url"
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
		Name:      "ANAME1",
		Domain:    "Auto highasianporn.com",
		HomePage:  "http://highasianporn.com/",
		FirstPage: "http://highasianporn.com/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)

	helper.AIR_CONFIGS[config.Name] = config
	Aname1.Name = config.Name
	Aname1.Description = config.Domain

	Aname1.Register()
}

var Aname1 = &Spider{
	//Name:         "ANAME1",
	//Description:  "[Auto Page] [highasianporn.com]",
	Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			entranceUrl := helper.AIR_CONFIGS[ctx.GetName()].HomePage
			keyIn := ctx.GetKeyin()
			if len(keyIn) > 4 {
				entranceUrl = keyIn
			}

			logs.Log.Warning("start with url:%v", entranceUrl)

			ctx.AddQueue(&request.Request{
				Url:  entranceUrl,
				Rule: "PICSETLIST",
				Header: http.Header{
					"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
					"User-Agent": []string{helper.AGENT_PUBLIC},
					"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
				},
			})
		},

		Trunk: map[string]*Rule{

			"PICSETLIST": {
				ItemFields: []string{
					"PICSETNAME",
					"PICSETURL",
					"PICSETTHUMB",
				},
				ParseFunc: func(ctx *Context) {

					helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot)

					logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					// cookie
					cookies := ""
					cookie := ctx.Response.Cookies()
					for _, c := range cookie {
						cookies += c.Name + "=" + c.Value + "; "
					}

					query := ctx.GetDom()
					lis := query.Find(".thumb") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
						if i > 10 {
							return
						}

						nextUrl, _ := s.Find("a").Attr("href")
						img, _ := s.Find("a img").Attr("src")
						picsetName := "" //s.Find("a").Text()

						img = FixUrl(img, ctx.GetUrl())

						if len(nextUrl) > 0 {
							//logs.Log.Warning("get a set url:%v", nextUrl)
							nextUrl = FixUrl(nextUrl, ctx.GetUrl())

							urlTemp := nextUrl
							if urlTemp[len(urlTemp)-1] == '/' {
								urlTemp = urlTemp[:len(urlTemp)-2]
							}

							for ii := len(urlTemp) - 1; ii > 0; ii-- {
								if urlTemp[ii] == '/' {
									picsetName = urlTemp[ii+1:]
									break
								}
							}

							picsetName = strings.Trim(picsetName, " \t")
							if len(picsetName) == 0 {
								picsetName = fmt.Sprintf("%v", i)
							}

							// record this picset
							ctx.Output(map[int]interface{}{
								0: picsetName,
								1: nextUrl,
								2: img,
							})

							// download in disk
							// save to local
							helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + picsetName)

							logs.Log.Warning("extract picset url, img=%v, %v, %v", nextUrl, img, picsetName)

							logs.Log.Warning("will request: %v", nextUrl)

							// queue request the picset detail
							ctx.AddQueue(
								&request.Request{
									Url:  nextUrl,
									Rule: "PICSET",
									Temp: map[string]interface{}{"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + picsetName, "PICSETNAME": picsetName},
									Header: http.Header{
										//"Accept-Language":           []string{"zh-CN,zh"},
										"Cookie":                    []string{cookies},
										"User-Agent":                []string{helper.AGENT_PUBLIC},
										"Referer":                   []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
										"Upgrade-Insecure-Requests": []string{"1"},
										//"Cache-Control":             []string{"no-cache"},
									},
									DownloaderID: 0,
								},
							)

							helper.DownloadObject(img, helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot, "thumb")
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
					query := ctx.GetDom()

					picsetName := ctx.GetTemp("PICSETNAME", "").(string)
					saveDir := ctx.GetTemp("DIR", "").(string)

					query.Find(".thumb_g").Each(func(i int, s *goquery.Selection) {
						imgUrl, _ := s.Find("a").Attr("href")
						logs.Log.Warning("IN %v imgUrl=%v", picsetName, imgUrl)

						imgUrl = FixUrl(imgUrl, ctx.GetUrl())

						savedPath := helper.DownloadObject(imgUrl, saveDir, "")

						ctx.Output(map[int]string{
							0: picsetName,
							1: savedPath,
							2: imgUrl,
						})
					})

					logs.Log.Warning("the res.len=%v status=%v header=%v", ctx.GetResponse().ContentLength, ctx.GetResponse().Status, ctx.GetResponse().Header)
				},
			},
		},
	},
}

func FixUrl(targetUrl string, route string) (finalUrl string) {

	defer func() {
		if len(finalUrl) > 0 {
			replacer := strings.NewReplacer("//www.", "//")
			finalUrl = replacer.Replace(finalUrl)
		}
	}()

	urlParsed, err := url.Parse(route)
	if err != nil {
		logs.Log.Warning("parse route targetUrl(%v) error:%v", route, err)
	}

	hasHttpInParam := strings.Index(targetUrl, "=http")
	if hasHttpInParam > 0 {
		finalUrl = targetUrl[hasHttpInParam+1:]
		return
	}

	finalUrl = targetUrl

	if len(targetUrl) > 4 && targetUrl[:4] == "http" {
		// legal
		return
	} else if len(targetUrl) > 0 && targetUrl[0] == '/' {
		homePage := urlParsed.Scheme + "://" + urlParsed.Hostname() + "/"
		if urlParsed.Port() != "" && urlParsed.Port() != "80" && urlParsed.Port() != "443" {
			homePage += ":" + urlParsed.Port()
		}
		finalUrl = homePage + targetUrl[1:]
	} else {
		finalUrl = route + targetUrl
	}
	return
}
