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
	TRANSLATOR_BAIDU  TRANSLATOR_ID = iota
	TRANSLATOR_GOOGLE TRANSLATOR_ID = iota
	TRANSLATOR_YOUDAO TRANSLATOR_ID = iota
)

type Translator interface {
	SetApiConfig(map[string]interface{})
	SetFromLang(lang string)
	SetToLang(lang string)
	Translate(str string) (string, error)
	AsyncTranslate(str string) int
	GetTransResult(int) (string, error)
	//AddTranslateTask(str string) *TranslatorTask
}

/*
   //auto splite long string to litle string, get result each and gother finally
*/

func SelectTranslator(id TRANSLATOR_ID) Translator {
	switch id {
	case TRANSLATOR_BAIDU:
		return &BaiduTranslator{}
	case TRANSLATOR_GOOGLE:
		//return &GoogleTranslator{}
	case TRANSLATOR_YOUDAO:
		//return &YoudaoTranslator{}

	}
	return nil
}
