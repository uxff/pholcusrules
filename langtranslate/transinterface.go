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

const (
	// baidu: http://api.fanyi.baidu.com/api/trans/product/apidoc
	TRANSLATOR_BAIDU = iota + 1
	// google:
	TRANSLATOR_GOOGLE
	// youdao chrome extension:
	TRANSLATOR_YOUDAO_FREE
	// youdao:
	TRANSLATOR_YOUDAO
	// microsoft azure: https://docs.microsoft.com/en-us/azure/cognitive-services/translator/
	TRANSLATOR_MICROSOFT
	// xunfei: http://www.xfyun.cn/services/mtranslate
	TRANSLATOR_XFYUN
	// tencent: https://cloud.tencent.com/document/product/551
	TRANSLATOR_TENCENT
	// aliyun: no translator api support, but third part supported
	//TRANSLATOR_ALIYUN  = iota // use iciba youdao baidu//
	// aiciba
	TRANSLATOR_ICIBA
	// yeecloud: https://market.aliyun.com/products/57124001/cmapi014709.html?spm=5176.730005.0.0.mUzPbY
	TRANSLATOR_YEECLOUD
)

type Translator interface {
	// todo: SetApiConfig(map[interface{}]interface{})
	SetApiConfig(map[string]interface{})
	SetFromLang(lang string)
	SetToLang(lang string)
	Translate(str string) (string, error)
}

type TransRes struct {
	Res            string
	Err            error
	AssumeFromLang string
	AssumeToLang   string
}

func SelectTranslator(translatorId int) Translator {
	switch translatorId {
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
