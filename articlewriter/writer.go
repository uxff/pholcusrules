package articlewriter

import xorm "github.com/go-xorm/xorm"
import core "github.com/go-xorm/core"

import (

	// _ "github.com/mattn/go-sqlite3"
	// h "github.com/m3ng9i/go-utils/http"

	// "crypto/md5"
	"time"

	// 去掉driver初始化 在别处已经引用 只能引用一次
	//_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/henrylee2cn/pholcus/logs" //信息输出
)

type ArticleWriter struct{}

var Orm *xorm.Engine
var OrmEngine string
var OrmDB string

//var NormalFetcher *h.Fetcher

var MysqlConnectionStr string = "www:123x456@tcp(127.0.0.1:3306)/xahoo?charset=utf8"
var defaultWriter = &ArticleWriter{}
var defaultAuthor string = "xahoo"
var defaultOrigin string = "wx100000p"
var defaultUrlMaker func(*ArticleEntity) string

func (this *ArticleWriter) Write(buf []byte) (n int, err error) {
	// json unmarshal from buf to entities

	var aList []ArticleEntity
	err = json.Unmarshal(buf, &aList)
	if err != nil {
		logs.Log.Error("json unmarshal error when write:%v", err)
		return -1, err
	}

	for i, _ := range aList {
		aList[i].Create_time = time.Now()

		if aList[i].Pubdate.IsZero() {
			aList[i].Pubdate = time.Now()
		}
	}

	SaveArticles(aList, defaultOrigin)

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

	insertedIds := make([]int, 0)

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
			logs.Log.Warning("empty title for save:%v", item.Outer_url)
			continue
		}
		if len(item.Author) == 0 {
			item.Author = defaultAuthor
		}

		// you should search, if url exist, do not save
		var exist bool

		if len(item.Outer_url) > 1 {
			var queryArticle = ArticleEntity{Outer_url: item.Outer_url}
			//Orm.QueryRow("select * from fh_article where ").Scan(&existArticle)
			rows, err := Orm.Rows(&queryArticle)
			if err != nil {
				logs.Log.Warning("could not query by outer_url, err:%v", err)
			} else {
				defer rows.Close()
				for rows.Next() {
					err = rows.Scan(&queryArticle)
					if err != nil {
						logs.Log.Warning("could not scan, err:%v", err)
					} else {
						exist = true
						break
					}
				}
			}

			if exist {
				logs.Log.Warning("outer_url already exist in db:%v", item.Outer_url)
				continue
			}
		}

		_, err = Orm.Insert(item)
		if err != nil {
			logs.Log.Error("insert Article error:%v item=%v", err, item)
			continue
		}
		//fmt.Println("insert success: num=", num, "all=", succNum, "id=", item.Id)
		insertedIds = append(insertedIds, item.Id)

		succNum++

	}

	logs.Log.Debug("all %v saved, ids=%v", succNum, insertedIds)
	return
}

func MakeArticleUrl(a *ArticleEntity) string {
	if defaultUrlMaker != nil {
		return defaultUrlMaker(a)
	}
	//strings.a.Id
	sign := "ignorethesestrings"
	str := "http://xahoo.xenith.top/index.php?r=article/show&id=" + fmt.Sprintf("%d", a.Id) + "&sign=" + sign
	return str
}

func Write(buf []byte) (n int, err error) {
	return defaultWriter.Write(buf)
}

func Close() error {
	return defaultWriter.Close()
}

func SetOrigin(v string) {
	defaultOrigin = v
}

func SetAuthor(v string) {
	defaultAuthor = v
}

func SetUrlMaker(maker func(*ArticleEntity) string) {
	defaultUrlMaker = maker
}
