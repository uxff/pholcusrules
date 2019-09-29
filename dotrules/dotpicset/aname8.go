package picsetcrawler

/*
curl http://sexygirlcity.com/index.php?gal=0_1
需求：
 - 下载静态网站中的图集
 - 记录图库资源
dev:done
download:done,12.3GB


*/

// 基础包
import (
	"fmt"
	"net/http"
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
		Name:      "ANAME8",
		Domain:    "kinghost.com",
		HomePage:  "http://www.kinghost.com/asian/manpe/pacific/841/",
		FirstPage: "http://www.kinghost.com/asian/manpe/pacific/",
	}

	config.DownloadRoot = fmt.Sprintf("./%s/", config.Name)
	// save: tagname/picsetname/pics*.jpg

	helper.AIR_CONFIGS[config.Name] = config
	Aname8.Name = config.Name
	Aname8.Description = config.Domain

	Aname8.Register()
}

var Aname8 = &Spider{
	//Name:         THE_DOMAIN,
	//Description:  THE_DOMAIN + " no need input",
	Pausetime:    300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			entranceUrl := helper.AIR_CONFIGS[ctx.GetName()].FirstPage
			//			keyIn := ctx.GetKeyin()
			//			if len(keyIn) > 4 {
			//				entranceUrl = keyIn
			//			}

			logs.Log.Warning("start with url:%v", entranceUrl)
			helper.MakeDir(helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot)

			for i := 370; i < 2000; i++ {
				ctx.AddQueue(&request.Request{
					Url:  entranceUrl + fmt.Sprintf("%d/", i),
					Rule: "PICSET",
					Header: http.Header{
						//"Cookie":     []string{helper.AIR_CONFIGS[ctx.GetName()].Cookie},
						"User-Agent": []string{helper.AGENT_PUBLIC},
						//"Referer":    []string{helper.AIR_CONFIGS[ctx.GetName()].HomePage},
					},
					Temp: map[string]interface{}{
						"DIR": helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot + fmt.Sprintf("/%d", i),
					},
				})

			}
		},

		Trunk: map[string]*Rule{

			"PICSET": {
				ItemFields: []string{
					"PICSETNAME",
					"SAVEPATH",
					"IMAGEURL",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//title := query.Find("title").Text()
					//picsetName := ctx.GetTemp("PICSETNAME", "").(string)
					tags := ""
					picsetName := query.Find("title").Text()
					query.Find("meta").Each(func(mi int, mo *goquery.Selection) {
						if _, keywordExist := mo.Attr("keywords"); keywordExist {
							tags, _ = mo.Attr("content")
						}
					})

					saveDir := helper.AIR_CONFIGS[ctx.GetName()].DownloadRoot // + picsetName
					helper.MakeDir(saveDir)
					helper.MakeDir(saveDir + "/thumbs")

					//thumbUrl := //ctx.GetTemp("THUMB_URL", "").(string)
					//helper.DownloadObject(thumbUrl, saveDir, "thumb")

					writeConfig := map[string]string{
						"title":   picsetName,
						"url":     ctx.GetUrl(),
						"tags":    "",
						"pubdate": "",
					}

					helper.WritePicsetConfig(writeConfig, saveDir)

					query.Find("table").Eq(0).Find("td").Each(func(tdi int, tdo *goquery.Selection) {

						colSpan, colSpanExist := tdo.Attr("colspan")
						if colSpanExist || len(colSpan) > 0 {
							return
						}

						largePic, _ := tdo.Find("a").Attr("href")
						littlePic, _ := tdo.Find("a").Find("img").Attr("src")

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
