package constant

const (
	BaseCampaignId = 50000000000
	DefaultFCA     = 5
)

const (
	DefaultCtaText = "install"
)

const (
	CampaignSourceOfferSync     = 0
	CampaignSourceAdnPortal     = 1
	CampaignSourceSSPlatform    = 2
	CampaignSourceSSAdvPlatform = 3
)

const (
	NormalAgency = 903
)

const (
	ATTR_BRAND_OFFER = 1
	ATTR_VTA_OFFER   = 2
	ATTR_CITY_OFFER  = 4
)

const (
	THIRD_PARTY_APPSFLYER = "appsflyer"
	THIRD_PARTY_ADJUST    = "adjust"
	THIRD_PARTY_S2S       = "s2s"
)

var CTALang = map[string]string{
	"ar":      "تثبيت",
	"en":      "Install",
	"fr":      "Installer",
	"de":      "Installieren",
	"id":      "Memasang",
	"it":      "Installare",
	"ja":      "インストール",
	"ko":      "설치",
	"nb":      "Installere",
	"pt":      "Instalar",
	"ru":      "Установить",
	"zh":      "安装",
	"es":      "Instalar",
	"sv":      "Installera",
	"th":      "ติดตั้ง",
	"zh-Hant": "安裝",
	"zh-Hans": "安装",
	"tr":      "İNDİR",
	"vi":      "cài đặt, dựng lên",
}

var CTAViewLang = map[string]string{
	"ar":      "للعرض",
	"en":      "View",
	"fr":      "Voir plus",
	"de":      "mehr sehen",
	"id":      "lihat lebih",
	"it":      "vedi dettagli",
	"ja":      "もっと",
	"ko":      "더 많은",
	"nb":      "se mer ",
	"pt":      "Veja mais ",
	"ru":      "Больше",
	"zh":      "查看",
	"es":      "ver más ",
	"sv":      " visa mer ",
	"th":      " ดูเพิ่มเติม",
	"zh-Hant": "查看",
	"zh-Hans": "查看",
	"tr":      "görünüm",
	"vi":      "lượt xem",
}
