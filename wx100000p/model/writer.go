package model

import xorm "github.com/go-xorm/xorm"
import core "github.com/go-xorm/core"

import (

	// _ "github.com/mattn/go-sqlite3"
	// h "github.com/m3ng9i/go-utils/http"

	// "crypto/md5"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/henrylee2cn/pholcus/logs" //信息输出
)

var Orm *xorm.Engine
var OrmEngine string
var OrmDB string

//var NormalFetcher *h.Fetcher

var MysqlConnectionStr string = "www:123x456@tcp(127.0.0.1:3306)/xahoo?charset=utf8"

type ArticleWriter struct{}

func (this *ArticleWriter) Write(buf []byte) (n int, err error) {
	// json unmarshal from buf to entities

	var aList []ArticleEntity
	err = json.Unmarshal(buf, &aList)
	if err != nil {
		return -1, err
	}

	for i, _ := range aList {
		aList[i].Create_time = time.Now()

		if aList[i].Pubdate.IsZero() {
			aList[i].Pubdate = time.Now()
		}
	}

	SaveArticles(aList, "wx100000p")

	// open session
	// write in
	return len(buf), nil
}

func (this *ArticleWriter) Close() error {
	//Orm.CloseSession()
	return nil
}

func init() {

	var err error

	OrmDB = "xahoo"

	Orm, err = xorm.NewEngine("mysql", MysqlConnectionStr) //"www:123x456@tcp(127.0.0.1:3306)/xahoo?charset=utf8")
	if err != nil {
		logs.Log.Error("orm init error:%v", err)
		return
	}
	Orm.SetMapper(core.SameMapper{})
}

func SaveArticles(items []ArticleEntity, origin string) (succNum int, err error) {
	succNum = 0

	//session := Orm.NewSession()
	for _, item := range items {

		if len(item.Origin) == 0 {
			item.Origin = origin
		}

		if len(strings.Trim(item.Content, " \t\r\n")) == 0 {
			logs.Log.Warning("empty content for save:%v", item.Outer_url)
			continue
		}
		if len(strings.Trim(item.Title, " \t\r\n")) == 0 {
			logs.Log.Warning("empty content for save:%v", item.Outer_url)
			continue
		}
		if len(item.Author) == 0 {
			item.Author = "wx100000p(uxff)"
		}

		_, err = Orm.Insert(item)
		if err != nil {
			logs.Log.Warning("insert Article error:%v item=%v", err, item)
			continue
		}
		//fmt.Println("insert success: num=", num, "all=", succNum, "id=", item.Id)

		succNum++

		continue
		// save as hot article, so show
		hotItem := new(HotArticleEntity)
		hotItem.Title = item.Title
		hotItem.Is_local_url = 1
		hotItem.Status = 2
		hotItem.Surface_url = item.Surface_url
		hotItem.Create_time = item.Create_time
		hotItem.Last_modified = item.Last_modified
		hotItem.Admin_id = item.Admin_id
		hotItem.Admin_name = item.Admin_name + "(gohead)"
		hotItem.Url = MakeArticleUrl(&item)

		_, err = Orm.Insert(hotItem)
		if err != nil {
			fmt.Println("insert hotArticle error:", err)
			continue
		}

	}

	logs.Log.Debug("all %v saved", succNum)
	return
}

func MakeArticleUrl(a *ArticleEntity) string {
	//strings.a.Id
	sign := "ignorethesestrings"
	str := "http://xahoo.xenith.top/index.php?r=article/show&id=" + fmt.Sprintf("%d", a.Id) + "&sign=" + sign
	return str
}
