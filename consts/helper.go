package consts

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/henrylee2cn/pholcus/logs" //信息输出
)

type AirConfig struct {
	Name         string
	Domain       string
	HomePage     string
	FirstPage    string
	DownloadRoot string
	Cookie       string
}

var (
	AIR_CONFIGS map[string]*AirConfig
)

func init() {
	AIR_CONFIGS = make(map[string]*AirConfig, 0)
}

func MakeDir(dirpath string) bool {
	err := os.Mkdir(dirpath, os.ModeDir)
	if err != nil && err != os.ErrExist {
		logs.Log.Error("mkdir error:%v", err)
		return false
	}
	return true
}

func DownloadObject(targetUrl string, saveDir string, saveName string) (savedPath string) {
	res, err := http.Get(targetUrl)
	if err != nil {
		logs.Log.Warning("download failed: targetUrl=%v err=%v", targetUrl, err)
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

		for i := len(targetUrl) - 1; i > 0; i-- {
			if targetUrl[i] == '/' {
				saveName = targetUrl[i+1:]
				break
			}
		}

		if len(saveName) == 0 {
			md5er := md5.New()
			md5er.Write([]byte(targetUrl))
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
func FixUrl(targetUrl string, route string) (finalUrl string) {

	finalUrl = targetUrl
	urlParsed, err := url.Parse(route)
	if err != nil {
		logs.Log.Warning("parse route url(%v) error:%v", route, err)
	}

	// is absolute url, start with http
	if len(targetUrl) > 4 && targetUrl[:4] == "http" {
		// legal
		return
	} else if len(targetUrl) > 0 && targetUrl[0] == '/' {
		// is start with /
		homePage := urlParsed.Scheme + "://" + urlParsed.Hostname() + "/"
		if urlParsed.Port() != "" && urlParsed.Port() != "80" && urlParsed.Port() != "443" {
			homePage += ":" + urlParsed.Port()
		}
		finalUrl = homePage + targetUrl[1:]
	} else {
		// is a relative url
		for i := len(route) - 1; i > 0; i-- {
			if route[i] == '/' {
				route = route[:i+1]
				break
			}
		}
		finalUrl = route + targetUrl
	}
	return
}
