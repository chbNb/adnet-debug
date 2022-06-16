package helpers

import (
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func handleLanguage(language string) string {
	// 去掉空格
	language = strings.Replace(language, " ", "", -1)
	// 如果没有-，则直接返回
	if !strings.Contains(language, "-") {
		return language
	}
	// 如果是中文
	if strings.Contains(language, "zh-Hant") {
		return "zh-Hant"
	}
	if strings.Contains(language, "zh-Hans") {
		return "zh-Hans"
	}
	if strings.Contains(language, "zh-TW") {
		return "zh-Hant"
	}
	arr := strings.Split(language, "-")
	return arr[0]
}

func RenderCtaText(linkType int32, language string) string {
	lang := handleLanguage(language)
	cta := "install"
	switch linkType {
	case 1, 2, 3:
		if cfgCta, ok := constant.CTALang[lang]; ok {
			cta = cfgCta
		}
	case 4, 5, 6:
		if cfgCta, ok := constant.CTAViewLang[lang]; ok {
			cta = cfgCta
		}
	}
	return cta
}
