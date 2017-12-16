package otherrule

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/henrylee2cn/pholcus/common/goquery" //DOM解析
	"github.com/henrylee2cn/pholcus/logs"           //信息输出

	// net包
	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"

	// 其他包
	"fmt"
	// "math"
	"time"
	// "io/ioutil"
)

func init() {
	QzoneArticlesx.Register()
}

var QzoneArticlesx = &Spider{
	Name:         "QZONE",
	Description:  `QZONE [自定义输入格式 "ID"::"Cookie"][最多支持250页，内设定时1~2s]`,
	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			param := strings.Split(ctx.GetKeyin(), "::")
			if len(param) != 2 {
				logs.Log.Error("自定义输入的参数不正确！")
				return
			}
			id := strings.Trim(param[0], " ")
			cookie := strings.Trim(param[1], " ")

			var i = 0
			for _, blogid := range BlogIds {
				i++
				if i > 20000 {
					break
				}
				ctx.AddQueue(&request.Request{
					//https://h5.qzone.qq.com/proxy/domain/b.qzone.qq.com/cgi-bin/blognew/blog_output_data?uin=545845496&blogid=1266762662&styledm=qzonestyle.gtimg.cn&imgdm=qzs.qq.com&bdm=b.qzone.qq.com&mode=2&numperpage=15&timestamp=1502100358&dprefix=&blogseed=0.17796215554699302&inCharset=gb2312&outCharset=gb2312&ref=qzone&entertime=1502100360571&cdn_use_https=1
					Url:  "https://h5.qzone.qq.com/proxy/domain/b.qzone.qq.com/cgi-bin/blognew/blog_output_data?uin=" + id + "&blogid=" + strconv.Itoa(blogid) + "&styledm=qzonestyle.gtimg.cn&imgdm=qzs.qq.com&bdm=b.qzone.qq.com&mode=2&numperpage=15&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&dprefix=&inCharset=gb2312&outCharset=gb2312&ref=qzone&page=3&refererurl=https%3A%2F%2Fqzs.qq.com%2Fqzone%2Fapp%2Fblog%2Fv6%2Fbloglist.html%23nojump%3D1%26page%3D3%26catalog%3Dlist",
					Rule: "文章详情",
					Header: http.Header{
						"Cookie":     []string{cookie},
						"Referer":    []string{"https://user.qzone.qq.com/" + id},
						"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"},
					},
					DownloaderID: 0,
				})
			}
		},

		Trunk: map[string]*Rule{
			"文章列表": {
				ItemFields: []string{
					"文章名",
					"blogid",
					"url",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					// 找到文章链接 加入队列
					结果 := map[int]interface{}{
						0: ctx.GetTemp("好友名", ""),
						1: ctx.GetTemp("好友ID", ""),
						2: ctx.GetTemp("认证", ""),
					}
					//var i int = 0
					query.Find(".article").Each(func(i int, s *goquery.Selection) {
						logs.Log.Error("this is eq %d", i)
						if i >= 3 {
							return
						}

						artLink := s.Find(".c_tx2 a")
						title, _ := artLink.Attr("title")
						name := artLink.Find("span").Text()
						fmt.Println(name)
						url, _ := artLink.Attr("href")
						blogid, _ := artLink.Attr("blogid")

						logs.Log.Error("i=%d title=%s name=%v url=%v blogid=%v\n", i, title, name, url, blogid)
						结果[i] = title

						/*
							x := &request.Request{
								Url:          url,
								Rule:         "文章详情",
								DownloaderID: 0,
								Temp: map[string]interface{}{
									"好友名":  name,
									"好友ID": uid,
									"认证":   认证,
									"关注":   关注,
									"粉丝":   粉丝,
									"微博":   微博,
								},
							}
							ctx.AddQueue(x)
						*/
					})
					// 结果输出
					ctx.Output(结果)
				},
			},
			"文章详情": {
				ItemFields: []string{
					"blogId",
					"文章标题",
					"文章内容",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var url = ctx.GetUrl()

					//var blogId int64 = 0
					var blogIdStr string
					if blogIdIdx := strings.Index(url, "blogid="); blogIdIdx > 0 {
						blogIdStr = url[blogIdIdx+9 : blogIdIdx+30]
						if commaIdx := strings.Index(blogIdStr, "&"); commaIdx > 0 {
							blogIdStr = blogIdStr[:commaIdx]
						}
						//fmt.Sscanf(string(blogIdStr), "%d", &blogId)
					}

					var detail = query.Find("#blogDetailDiv").Text()
					var title = query.Find(".blog_tit_detail").Eq(0).Text()
					logs.Log.Error("blogid=%v title=%v len(detail)=%v ", blogIdStr, title, len(detail))
					rowRet := map[int]interface{}{
						0: blogIdStr,
						1: title,
						2: detail,
					}

					// 结果输出
					ctx.Output(rowRet)
				},
			},
		},
	},
}
