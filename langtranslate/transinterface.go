package langtranslate

type LANG_TYPE string

const (
	LANG_AUTO LANG_TYPE = "auto"
	LANG_ZH   LANG_TYPE = "zh"
	LANG_EN   LANG_TYPE = "en"
	LANG_JP   LANG_TYPE = "jp"
	LANG_KR   LANG_TYPE = "kr"
	LANG_DE   LANG_TYPE = "de"
	LANG_FR   LANG_TYPE = "fr"
	LANG_SP   LANG_TYPE = "sp"
	LANG_RU   LANG_TYPE = "ru"
)

type TRANSLATOR_ID int

const (
	// baidu: http://api.fanyi.baidu.com/api/trans/product/apidoc
	TRANSLATOR_BAIDU TRANSLATOR_ID = iota
	// google:
	TRANSLATOR_GOOGLE TRANSLATOR_ID = iota
	// youdao chrome extension:
	TRANSLATOR_YOUDAO_FREE TRANSLATOR_ID = iota
	// youdao:
	TRANSLATOR_YOUDAO TRANSLATOR_ID = iota
	// microsoft azure: https://docs.microsoft.com/en-us/azure/cognitive-services/translator/
	TRANSLATOR_MICROSOFT TRANSLATOR_ID = iota
	// xunfei: http://www.xfyun.cn/services/mtranslate
	TRANSLATOR_XFYUN TRANSLATOR_ID = iota
	// tencent: https://cloud.tencent.com/document/product/551
	TRANSLATOR_TENCENT TRANSLATOR_ID = iota
	// aliyun: no translator api support, but third part supported
	//TRANSLATOR_ALIYUN TRANSLATOR_ID = iota // use iciba youdao baidu//
	//
	TRANSLATOR_ICIBA TRANSLATOR_ID = iota
	// yeecloud: https://market.aliyun.com/products/57124001/cmapi014709.html?spm=5176.730005.0.0.mUzPbY
	TRANSLATOR_YEECLOUD TRANSLATOR_ID = iota
)

type Translator interface {
	SetApiConfig(map[string]interface{})
	SetFromLang(lang string)
	SetToLang(lang string)
	Translate(str string) (string, error)
	AsyncTranslate(str string) <-chan *TransRes
	GetTransResult(int) (string, error)
	//AddTranslateTask(str string) *TranslatorTask
}

type TransRes struct {
	Res string
	Err error
}

func SelectTranslator(id TRANSLATOR_ID) Translator {
	switch id {
	case TRANSLATOR_BAIDU:
		return &BaiduTranslator{}
	case TRANSLATOR_GOOGLE:
		//return &GoogleTranslator{}
	case TRANSLATOR_YOUDAO_FREE:
		return &YoudaoFreeTranslator{}
	case TRANSLATOR_YOUDAO:
		return &YoudaoTranslator{}
	}
	return nil
}
