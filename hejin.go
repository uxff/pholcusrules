package otherrule

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	//"github.com/henrylee2cn/pholcus/common/goquery" //DOM解析
	"github.com/henrylee2cn/pholcus/logs" //信息输出

	"bufio"
	"bytes"
	"io"
	//"io/ioutil"
	// net包
	"net/http" //设置http.Header
	"net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"

	// 其他包
	//"fmt"
	// "math"
	"math/rand"
	"time"
)

func init() {
	Hejinx.Register()
}

var Hejinx = &Spider{
	Name:         "HEJINx",
	Description:  `HEJINx 自定义输入格式 url`,
	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//ctx.Request is nil, dont use it here
			param := ctx.GetKeyin()
			if len(param) <= 12 {
				logs.Log.Warning("自定义输入的url参数不正确！ use default")
				//return
				param = `http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=20`
			}
			//logs.Log.Warning("param=%v", param)

			//http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=20
			urlParsed, _ := url.Parse(param)
			urlParams, _ := url.ParseQuery(urlParsed.RawQuery)
			logs.Log.Error("host=%v script=%v query=%v", urlParsed.Host, urlParsed.Path, urlParsed.RawPath)
			urlModel, modelExist := urlParams["model"]
			if !modelExist || len(urlModel) == 0 {
				logs.Log.Error("不是有效的url,model not exist,有效的url应该类似：http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=1")
				return
			}
			if urlParsed.Path != "/plugin.php" {
				logs.Log.Error("不是有效的url,plugin.php not exist,有效的url应该类似：http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=1 %v", urlParsed.Path)
				return
			}

			pluginId, pluginIdExist := urlParams["id"]
			if !pluginIdExist || len(pluginId) == 0 {
				logs.Log.Error("不是有效的url,pluginId not exist,有效的url应该类似：http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=1 ")
				return
				///} else if string(pluginId[0])[:5] != "hejin" {
				//logs.Log.Error("不是有效的url,pluginId != hejin*,有效的url应该类似：http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=detail&zid=1 %v", pluginId[0])
				//return
			}
			logs.Log.Error("pluginId=%v urlModel=%v", pluginId, urlModel)
			vid, vidExist := urlParams["vid"]
			if !vidExist || len(vid) == 0 {
				vid = make([]string, 1)
				vid[0] = "1"
			}

			logs.Log.Error("vid=%v", vid)
			zid, zidExist := urlParams["zid"]
			if !zidExist || len(zid) == 0 {
				logs.Log.Error("没有匹配到要投票的用户 请输入带zid的url %v", zid[0])
				return
			}

			//logs.Log.Error("zid=%v", zid)
			//ctx.SetTemp("pluginId", pluginId[0]) // cannot use ctx.SetTemp in root

			/* 思路：获取所有zid；根据zid导出openid；唯一openid；请求ticket带cookie */
			// urlTop300 := "http://tzxts.lzyjdzsw.com/plugin.php?id=hejin_toupiao&model=top300&vid=1#top300"
			// formhash 在 top300 中
			urlPre := urlParsed.Scheme + "://" + urlParsed.Host + urlParsed.Path + "?id=" + pluginId[0]
			//ctx.SetTemp("urlPre", urlPre)
			urlTop300 := urlPre + "&model=top300&vid=" + vid[0]
			logs.Log.Warning("will top300: %v", urlTop300)

			ctx.AddQueue(&request.Request{
				Url:  urlTop300,
				Rule: "top300",
				Header: http.Header{
					"Cookie":     []string{},
					"Referer":    []string{param},
					"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"},
				},
				DownloaderID: 0,
				Temp: map[string]interface{}{
					"zid":      zid[0],
					"model":    urlModel[0],
					"pluginId": pluginId[0],
					"urlPre":   urlPre,
					"vid":      vid[0],
				},
			})
		},

		Trunk: map[string]*Rule{
			"top300": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Warning("start top300: url=%v", ctx.GetUrl())
					//query := ctx.GetDom()
					//logs.Log.Warning("got query when top300")
					if ctx.Response == nil {
						logs.Log.Error("no response!")
						return
					}
					textContent := ctx.GetText()
					tempContent := []byte(textContent)
					logs.Log.Warning("the textContent len=%v %s", len(textContent), string(tempContent[:32]))

					// find formhash
					formhash := ""
					formhashIdx := bytes.Index(tempContent, []byte("&formhash="))
					if formhashIdx > 0 {
						formhash = string(tempContent[formhashIdx+10 : formhashIdx+10])
						formhashEndIdx := bytes.IndexByte(tempContent[formhashIdx+10:formhashIdx+30], '\'')
						if formhashEndIdx > 0 {
							formhash = string(tempContent[formhashIdx+10 : formhashIdx+10+formhashEndIdx])
						}
						//strings.TrimRight()
						formhash = strings.TrimRight(formhash, ` \r\n\t\'\"=\&`)
						logs.Log.Warning("find formhash=%v", formhash)
					}

					// find rank list
					rankIdx := strings.Index(textContent, `<div class="rank300" id="top300">`)
					var rankContent string
					if rankIdx > 0 {
						logs.Log.Warning("find div rank300 at: %d", rankIdx)
						tempContent = tempContent[rankIdx:]
						textContent = string(tempContent)
						rankEndIdx := strings.Index(textContent, "</div>")
						if rankEndIdx > 0 {
							logs.Log.Warning("find </div> at: %d", rankEndIdx)
							//tempContent = []byte(textContent)
							tempContent = tempContent[:rankEndIdx]
							rankContent = string(tempContent)
						}
					}
					logs.Log.Warning("the rankContent len=%v %v", len(rankContent), string(tempContent[:32]))
					if len(rankContent) == 0 {
						logs.Log.Error("no rank content in text:%v", string(tempContent))
						return
					}

					// 因为该网站的代码比较垃圾 编码混乱，gb2312和utf8混排，导致goquery无法解析，只能手动
					var uids []int
					tempContent = []byte(rankContent)
					for {
						spanIdx := strings.Index(rankContent, "</span><span>1")
						uid := 0
						if spanIdx > 0 {
							logs.Log.Warning("find a span: %v", string(tempContent[spanIdx:spanIdx+18]))
							uid, _ = strconv.Atoi(string(tempContent[spanIdx+14 : spanIdx+18]))
							tempContent = tempContent[spanIdx+18:]
							rankContent = string(tempContent)

							if uid > 0 {
								uids = append(uids, uid)
							}

							continue
						}
						break
					}

					logs.Log.Warning("uids: len=%d %v", len(uids), uids)
					randSlice(uids)

					for _, uid := range uids {
						url := ctx.GetTemp("urlPre", "").(string) + "&model=dcexcel&zid=" + strconv.FormatInt(int64(uid), 10)
						logs.Log.Warning("will dcexcel: %v", url)

						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "dcexcel",
							Header: http.Header{
								"Cookie":     []string{},
								"Referer":    []string{url},
								"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"},
							},
							DownloaderID: 0,
							Temp: map[string]interface{}{
								"urlPre":   ctx.GetTemp("urlPre", ""),
								"vid":      ctx.GetTemp("vid", ""),
								"zid":      ctx.GetTemp("zid", ""),
								"formhash": formhash,
							},
						})

					}
					ctx.SetTemp("uids", uids)
				},
			},
			"dcexcel": {
				ItemFields: []string{
					"zid",
					"openid",
				},
				ParseFunc: func(ctx *Context) {

					text := ctx.GetText()
					var url = ctx.GetUrl()
					logs.Log.Warning("start dcexcel:%v len=%v", url, len(text))
					openIdNum := 0
					openIdMap := map[string]string{}
					doTicket := false

					if len(text) > 0 {
						br := bufio.NewReader(strings.NewReader(text))
						for {
							a, _, c := br.ReadLine()
							openIdNum++
							if c == io.EOF {
								break
							}
							commaPos := strings.Index(string(a), ",")
							zid := ""
							openId := ""
							if commaPos > 0 {
								zid = string(a[:commaPos])
								a = a[commaPos+1:]
								commaPos2 := strings.Index(string(a), ",")
								if commaPos2 > 0 {
									openId = string(a[:commaPos2])
								} else {
									openId = string(a)
								}
								if len(openId) > 0 && openId[0] != 'o' {
									// not valid openid
									continue
								}

								if _, oExist := openIdMap[openId]; !oExist {
									openIdMap[openId] = zid
								}
								//logs.Log.Warning("got a openid:%v %v", zid, openId)

								// 中间被停止
								//time.Sleep(time.Second * 3)
							}
							if openIdNum > 10 {
								//break
							}
						}

						logs.Log.Warning("openIdNum=%v len(map)=%v", openIdNum, len(openIdMap))
						for openId, zidVal := range openIdMap {
							rowRet := map[int]interface{}{
								0: zidVal,
								1: openId,
							}
							ctx.Output(rowRet)

							if doTicket {
								// start ticket
								url = ctx.GetTemp("urlPre", "").(string) + "&model=ticket&zid=" + ctx.GetTemp("zid", "").(string) + "&formhash=" + ctx.GetTemp("formhash", "").(string)
								ctx.AddQueue(&request.Request{
									Url:  url,
									Rule: "ticket",
									Header: http.Header{
										"Cookie":     []string{"hjbox_openid=" + ctx.GetTemp("openid", "").(string)},
										"Referer":    []string{url},
										"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"},
									},
									DownloaderID: 0,
									Temp: map[string]interface{}{
										"urlPre":   ctx.GetTemp("urlPre", ""),
										"vid":      ctx.GetTemp("vid", ""),
										"zid":      ctx.GetTemp("zid", ""),
										"openid":   openId,
										"formhash": ctx.GetTemp("formhash", ""),
									},
								})
							}
						}
					}
					return
				},
			},
			"ticket": {
				ItemFields: []string{
					"zid",
					"openid",
					"status",
				},
				ParseFunc: func(ctx *Context) {
					text := ctx.GetText()
					status := []byte(text)[:16]
					zid := ctx.GetTemp("zid", "")
					openId := ctx.GetTemp("openid", "")
					logs.Log.Warning("TICKET zid=%v openid=%v len(text)=%v status=%s", zid, openId, len(text), string(status))

					rowRet := map[int]interface{}{
						0: zid,
						1: openId,
						2: string(status),
					}

					// 结果输出
					ctx.Output(rowRet)
				},
			},
		},
	},
}

func randSlice(a []int) []int {
	rand.Seed(time.Now().UnixNano())
	thelen := len(a)
	for i := 0; i < 1000; i++ {
		r1 := rand.Int() % thelen
		r2 := rand.Int() % thelen
		t := a[r2]
		a[r2] = a[r1]
		a[r1] = t
	}
	return a
}
