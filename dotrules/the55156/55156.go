package the55156

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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	//"io/ioutil"
	"net/http"
	"os"
	//"strconv"
	"strings"

	//"golang.org/x/net/html"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	. "github.com/henrylee2cn/pholcus/app/spider" //必需

	"github.com/uxff/pholcusrules/consts"
)

const (
	PUBLIC_COOKIE = ""
	DOWNLOAD_ROOT = "./55156.com/"
	HOME_URL      = "http://www.55156.com/"
	FIRST_PAGE    = "http://www.55156.com/"
)

func init() {
	The55156.Register()
}

var The55156 = &Spider{
	Name:         "55156.com",
	Description:  "[Auto Page] [www.55156.com]",
	Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			entranceUrl := FIRST_PAGE
			keyIn := ctx.GetKeyin()
			if len(keyIn) > 4 {
				entranceUrl = keyIn
			}

			logs.Log.Warning("start with url:%v", entranceUrl)
			MakeDir(DOWNLOAD_ROOT)

			ctx.AddQueue(&request.Request{
				Url:  entranceUrl,
				Rule: "HOMEPAGE",
				Header: http.Header{
					"Cookie":     []string{PUBLIC_COOKIE},
					"User-Agent": []string{consts.AGENT_PUBLIC},
					"Referer":    []string{HOME_URL},
				},
			})
		},

		Trunk: map[string]*Rule{
			"HOMEPAGE": {
				ParseFunc: func(ctx *Context) {
					// home page like: http://www.55156.com/
					//logs.Log.Warning("content len of list=%v err=%v", ctx.Response.ContentLength, ctx.GetError())

					query := ctx.GetDom()
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
							url = FixUrl(url, ctx.GetUrl())

							// download in disk
							// save to local
							MakeDir(DOWNLOAD_ROOT + tagName)

							//logs.Log.Warning("extract tag url, img=%v, %v, %v", url, img, tagName)

							// cookie
							cookies := ""
							cookie := ctx.Response.Cookies()
							for _, c := range cookie {
								cookies += c.Name + "=" + c.Value + "; "
							}

							logs.Log.Warning("will request home->picsetlist: %v", url)

							// queue request the picset detail
							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: "PICSETLIST",
									Temp: map[string]interface{}{"DIR": DOWNLOAD_ROOT + tagName, "TAGNAME": tagName},
									Header: http.Header{
										//"Accept-Language":           []string{"zh-CN,zh"},
										"Cookie":     []string{cookies},
										"User-Agent": []string{consts.AGENT_PUBLIC},
										"Referer":    []string{HOME_URL},
										//"Upgrade-Insecure-Requests": []string{"1"},
										//"Cache-Control":             []string{"no-cache"},
									},
									DownloaderID: 0,
								},
							)

							//DownloadObject(img, DOWNLOAD_ROOT+tagName, "thumb")
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
					lis := query.Find(".labelbox_pic").Find("ul").Find("li") // 不能写 ".thumb a"
					lis.Each(func(i int, s *goquery.Selection) {
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

						img = FixUrl(img, ctx.GetUrl())

						if len(url) > 0 {
							//logs.Log.Warning("get a set url:%v", url)
							url = FixUrl(url, ctx.GetUrl())

							// record this picset
							ctx.Output(map[int]interface{}{
								0: tagName,
								1: url,
								2: img,
							})

							// download in disk
							// save to local
							MakeDir(DOWNLOAD_ROOT + tagName)

							//logs.Log.Warning("extract tag url, img=%v, %v, %v", url, img, tagName)

							// cookie
							cookies := ""
							cookie := ctx.Response.Cookies()
							for _, c := range cookie {
								cookies += c.Name + "=" + c.Value + "; "
							}

							//logs.Log.Warning("cookie=%s ", cookies)
							//cookies = "25a6da5acf7fde759f79e8c23ab0dc76d53f8=cGxmUnIwNWQ1T3JvTVRVeE16QTVNamsyTlMwd0xTRXcb; f6d9f67e5f2c5b9dd954e79d40588261f042e31abda5d=bkpHaHAwMU1ZV0U0TlRsbFpqRmlaREF4TTJaa1pXVXlORFZpTlRRd01ETXhNamRqWkRNPQc; 9353eb=1513092965; _ga=GA1.2.1711782959.1513095156; _gid=GA1.2.1692799327.1513095156; 34843e6d0d96c28940bc888267e9b3=ekxwRzExN1FTRVpoVXVuRHlKV3VMOTB1ZERRNU9UUTVPQT09a; 073273c2a4e3c0d936022720d=SzZlRE4xNlZCdk00QUdleWQ5NlVHMWVNaTB3a; 9353e=bm9yZWZ8fGRlZmF1bHR8MXwyfDJ8bm9uZXwwOmpwZ3JhdnVyZS5jb20%3D; __atuvc=3%7C50; __atuvs=5a2fffe7e5348de2002"

							logs.Log.Warning("will request tag->picsetlist: %v", url)

							// queue request the picset detail
							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: "PICSETLIST",
									Temp: map[string]interface{}{"DIR": DOWNLOAD_ROOT + tagName + "/", "TAGNAME": tagName},
									Header: http.Header{
										//"Accept-Language":           []string{"zh-CN,zh"},
										"Cookie":     []string{cookies},
										"User-Agent": []string{consts.AGENT_PUBLIC},
										"Referer":    []string{HOME_URL},
										//"Upgrade-Insecure-Requests": []string{"1"},
										//"Cache-Control":             []string{"no-cache"},
									},
									DownloaderID: 0,
								},
							)

							DownloadObject(img, DOWNLOAD_ROOT+tagName, "thumb")
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

					saveDir := ctx.GetTemp("DIR", DOWNLOAD_ROOT).(string)

					query := ctx.GetDom()
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

						img = FixUrl(img, ctx.GetUrl())

						if len(url) > 0 {
							//logs.Log.Warning("get a set url:%v", url)
							url = FixUrl(url, ctx.GetUrl())

							// record this picset
							ctx.Output(map[int]interface{}{
								0: picsetName,
								1: url,
								2: img,
							})

							// download in disk
							// save to local
							MakeDir(saveDir + "/" + picsetName)

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
										"User-Agent": []string{consts.AGENT_PUBLIC},
										"Referer":    []string{HOME_URL},
										//"Upgrade-Insecure-Requests": []string{"1"},
										//"Cache-Control":             []string{"no-cache"},
									},
									DownloaderID: 0,
								},
							)

							DownloadObject(img, saveDir+"/"+picsetName, "thumb")
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
					saveDir := ctx.GetTemp("DIR", DOWNLOAD_ROOT+picsetName).(string)

					imgUrl, _ := query.Find("#picBody").Find("a img").Attr("src")
					imgUrl = FixUrl(imgUrl, ctx.GetUrl())

					DownloadObject(imgUrl, saveDir, "")
					//logs.Log.Warning("IN %v imgUrl=%v", picsetName, imgUrl)

					query.Find(".pages").Find("ul").Find("li").Each(func(i int, s *goquery.Selection) {
						nextPageUrl, _ := s.Find("a").Attr("href")
						nextPageUrl = FixUrl(nextPageUrl, ctx.GetUrl())
						ctx.AddQueue(
							&request.Request{
								Url:  nextPageUrl,
								Rule: "PICSET",
								Temp: map[string]interface{}{"DIR": saveDir, "PICSETNAME": picsetName},
								Header: http.Header{
									//"Accept-Language":           []string{"zh-CN,zh"},
									//"Cookie":     []string{cookies},
									"User-Agent": []string{consts.AGENT_PUBLIC},
									"Referer":    []string{HOME_URL},
									//"Upgrade-Insecure-Requests": []string{"1"},
									//"Cache-Control":             []string{"no-cache"},
								},
								DownloaderID: 0,
							},
						)

					})

					return

					/*
						pageNo := ctx.GetTemp("PAGE_NO", 1).(int)

						allPageMark := query.Find(".pages").Find("ul").Find("li").Eq(0).Text()
						//allPageMarkRune := rune(allPageMark)
						allPageMark = strings.Trim(allPageMark, " \t\r\n ")
						allPageMark = strings.Trim(allPageMark, "共页：:")
						allPageNo, _ := strconv.Atoi(allPageMark)

						logs.Log.Warning("the res.len=%v status=%v header=%v", ctx.GetResponse().ContentLength, ctx.GetResponse().Status, ctx.GetResponse().Header)
					*/
				},
			},
		},
	},
}

func MakeDir(dirpath string) bool {
	err := os.Mkdir(dirpath, os.ModeDir)
	if err != nil {
		logs.Log.Error("mkdir error:%v", err)
		return false
	}
	return true
}

func DownloadObject(url string, saveDir string, saveName string) (savedPath string) {
	res, err := http.Get(url)
	if err != nil {
		logs.Log.Warning("download failed: url=%v err=%v", url, err)
		return
	}

	contentType := res.Header.Get("Content-Type")
	exts, err := mime.ExtensionsByType(contentType)

	//logs.Log.Warning("got mime type from %v=%v err=%v", contentType, exts, err)
	ext := ""
	if len(exts) > 0 {
		ext = exts[0]
	} else {
		ext = ".jpg"
	}

	if saveName == "" {

		for i := len(url) - 1; i > 0; i-- {
			if url[i] == '/' {
				saveName = url[i+1:]
				break
			}
		}

		if len(saveName) == 0 {
			md5er := md5.New()
			md5er.Write([]byte(url))
			saveName = hex.EncodeToString(md5er.Sum(nil)) + ext
		}
	}

	if strings.IndexByte(saveName, '.') < 0 {
		saveName = saveName + ext
	}

	savedPath = saveDir + "/" + saveName

	fhandle, err := os.Create(savedPath) //, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		logs.Log.Warning("create file failed: path=%v err=%v", savedPath, err)
		return ""
	}

	io.Copy(fhandle, res.Body)

	fhandle.Close()

	return savedPath
}

// you should consider url begin with '#'
func FixUrl(url string, route string) (finalUrl string) {

	finalUrl = url

	if len(url) > 4 && url[:4] == "http" {
		// legal
		return
	} else if len(url) > 0 && url[0] == '/' {
		finalUrl = HOME_URL + url[1:]
	} else {
		for i := len(route) - 1; i > 0; i-- {
			if route[i] == '/' {
				route = route[:i+1]
				break
			}
		}
		finalUrl = route + url
	}
	return
}
