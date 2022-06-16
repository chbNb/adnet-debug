package mvutil

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
)

type AppInfo struct {
	AppId              int64             `bson:"appId,omitempty" json:"appId"`
	App                App               `bson:"app,omitempty" json:"app"`
	Publisher          Publisher         `bson:"publisher,omitempty" json:"publisher"`
	Updated            int64             `bson:"updated,omitempty" json:"updated"`
	LandingPageVersion []string          `bson:"landingPageVersion,omitempty" json:"landingPageVersion"`
	RealPackageName    string            `bson:"realPackageName,omitempty" json:"realPackageName"`
	JumptypeConfig     map[string]int32  `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	JumptypeConfigV2   map[string]int32  `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`
	Rewards            []Reward          `bson:"rewards,omitempty" json:"rewards,omitempty"`
	BlackCategoryList  *[]int64          `bson:"blackCategoryList,omitempty" json:"blackCategoryList,omitempty"`
	BlackPackageList   *[]string         `bson:"blackPackageList,omitempty" json:"blackPackageList,omitempty"`
	DirectOfferConfig  DirectOfferConfig `bson:"directOfferConfig,omitempty" json:"directOfferConfig,omitempty"`
}

type App struct {
	AppId           int64               `bson:"appId,omitempty" json:"appId"`
	Grade           int                 `bson:"grade,omitempty" json:"grade"`
	Status          int                 `bson:"status,omitempty" json:"status"`
	Platform        int                 `bson:"platform,omitempty" json:"platform"`
	IsIncent        int                 `bson:"isIncent,omitempty" json:"isIncent"`
	Name            string              `bson:"name,omitempty" json:"name"`
	OfficialName    string              `bson:"officialName,omitempty" json:"officialName"`
	DevinfoEncrypt  int                 `bson:"devinfoEncrypt,omitempty" json:"devinfoEncrypt,omitempty"`
	BtClass         int                 `bson:"btClass,omitempty" json:"btClass,omitempty"`
	FrequencyCap    int                 `bson:"frequencyCap,omitempty" json:"frequencyCap,omitempty"`
	StorekitLoading int32               `bson:"storekitLoading,omitempty" json:"storekitLoading,omitempty"`
	DevIdAllowNull  int                 `bson:"devIdAllowNull,omitempty" json:"devIdAllowNull,omitempty"`
	OfferPreference []int               `bson:"offerPreference,omitempty" json:"offerPreference,omitempty"`
	ImpressionCap   int                 `bson:"impressionCap,omitempty" json:"impressionCap"`
	Domain          string              `bson:"domain,omitempty" json:"domain,omitempty"`
	DomainVerify    int                 `bson:"domain_verify,omitempty" json:"domain_verify,omitempty"`
	Plct            int                 `bson:"plct,omitempty" json:"plct"`
	Coppa           int                 `bson:"coppa,omitempty" json:"coppa"`
	IabCategoryV2   map[string][]string `bson:"iabCategoryV2" json:"iabCategoryV2"`
	StoreUrl        string              `bson:"storeUrl,omitempty" json:"storeUrl,omitempty"`
	M1Num           int                 `bson:"m1Num,omitempty" json:"m1Num,omitempty"`
	BundleId        string              `bson:"bundleId,omitempty" json:"bundleId,omitempty"`
	Ccpa            int                 `bson:"ccpa,omitempty" json:"ccpa,omitempty"`
	LivePlatform    int                 `bson:"livePlatform,omitempty" json:"livePlatform,omitempty"`
}

type RudeceRule struct {
	Priority int `bson:"priority,omitempty" json:"priority"`
	Install  int `bson:"install,omitempty" json:"install"`
	Status   int `bson:"status,omitempty" json:"status"`
	Start    int `bson:"start,omitempty" json:"start"`
	Price    int `bson:"price,omitempty" json:"price"`
}

type VTAConf struct {
	Status int `bson:"status,omitempty" json:"status"`
	Rate   int `bson:"rate,omitempty" json:"rate"`
	Rule   int `bson:"rule,omitempty" json:"rule"`
}

type FillRate struct {
	Rate   int `bson:"rate,omitempty" json:"rate"`
	Status int `bson:"status,omitempty" json:"status"`
}

func IsDevinfoEncrypt(appInfo *AppInfo) bool {
	if appInfo == nil {
		return true
	}
	if appInfo.App.DevinfoEncrypt == mvconst.DEVINFO_ENCRYPT_DONT {
		return false
	}
	return true
}

func AppFcaDefault(fca int) bool {
	if fca == 0 || fca == 2 {
		return true
	}
	return false
}
