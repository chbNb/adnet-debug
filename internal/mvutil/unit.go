package mvutil

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

const FixedEcpm = 5

type UnitInfo struct {
	AdSourceLen     int
	UnitId          int64              `bson:"unitId,omitempty" json:"unitId,omitempty"`
	Unit            Unit               `bson:"unit,omitempty" json:"unit,omitempty"`
	AppId           int64              `bson:"appId,omitempty" json:"appId,omitempty"`
	Setting         Setting            `bson:"setting,omitempty" json:"setting,omitempty"`
	AdSourceCountry map[string]int     `bson:"adSourceCountry,omitempty" json:"adSourceCountry,omitempty"`
	AdSourceData    []AdSourceDataInfo `bson:"adSourceData,omitempty" json:"adSourceData,omitempty"`
	Updated         int64              `bson:"updated,omitempty" json:"updated,omitempty"`
	// VirtualRewardRaw bson.Raw           `bson:"virtualReward,omitempty" json:"virtualReward,omitempty"`
	VirtualReward    VirtualReward    `bson:"virtualReward,omitempty" json:"virtualReward,omitempty"`
	Endcard          *smodel.EndCard  `bson:"endcard,omitempty" json:"endcard,omitempty"`
	JumptypeConfig   map[string]int32 `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	MPToMV           *MPMVObj         `bson:"mp2mv,omitempty" json:"mp2mv,omitempty"`
	MVToMP           *MPMVObj         `bson:"mv2mp,omitempty" json:"mv2mp,omitempty"`
	JumptypeConfigV2 map[string]int32 `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`
	// AdSourceTimeRaw   bson.Raw         `bson:"adSourceTime,omitempty" json:"adSourceTime,omitempty"`
	AdSourceTime      *map[string]int64   `bson:"adSourceTime,omitempty" json:"adSourceTime,omitempty"`
	CdnSetting        []*CdnSetting       `bson:"CDNSetting,omitempty" json:"CDNSetting,omitempty"`
	EcpmFloors        map[string]float64  `bson:"ecpmFloors,omitempty" json:"ecpmFloors,omitempty"`
	BlackCategoryList *[]int64            `bson:"blackCategoryList,omitempty" json:"blackCategoryList,omitempty"`
	BlackPackageList  *[]string           `bson:"blackPackageList,omitempty" json:"blackPackageList,omitempty"`
	DirectOfferConfig DirectOfferConfig   `bson:"directOfferConfig,omitempty" json:"directOfferConfig,omitempty"`
	FakeRuleV2        map[string]FakeRule `bson:"fakeRuleV2,omitempty" json:"fakeRuleV2,omitempty"`
}

type CdnSetting struct {
	Id     int    `bson:"id,omitempty" json:"id,omitempty"`
	Url    string `bson:"url,omitempty" json:"url,omitempty"`
	Weight int    `bson:"weight,omitempty" json:"weight,omitempty"`
}

type Reward struct {
	ID     int64  `bson:"id,omitempty" json:"id,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Amount int64  `bson:"amount,omitempty" json:"amount,omitempty"`
}

type MPMVObj struct {
	PublisherId int `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	AppId       int `bson:"appId,omitempty" json:"appId,omitempty"`
	UnitId      int `bson:"unitId,omitempty" json:"unitId,omitempty"`
}

type VirtualReward struct {
	Name         string `bson:"name,omitempty" json:"name,omitempty"`
	ExchangeRate int    `bson:"exchange_rate,omitempty" json:"exchange_rate,omitempty"`
	StaticReward int    `bson:"static_reward,omitempty" json:"static_reward,omitempty"`
}
type AdSourceDataEntry struct {
	AdSourceId int `bson:"adSourceId,omitempty" json:"adSourceId,omitempty"`
	Status     int `bson:"status,omitempty" json:"status,omitempty"`
	Priority   int `bson:"priority,omitempty" json:"priority,omitempty"`
}

type AdSourceDataInfo []AdSourceDataEntry

type Unit struct {
	UnitID               int64               `bson:"unitId,omitempty" json:"unitId,omitempty"`
	IsIncent             int                 `bson:"isIncent,omitempty" json:"isIncent,omitempty"`
	BtClass              int                 `bson:"btClass,omitempty" json:"btClass,omitempty"`
	Status               int                 `bson:"status,omitempty" json:"status,omitempty"`
	AdType               int32               `bson:"adType,omitempty" json:"adType,omitempty"`
	VideoAds             int                 `bson:"videoAds,omitempty" json:"videoAds,omitempty"`
	Orientation          int                 `bson:"orientation,omitempty" json:"orientation,omitempty"`
	Templates            []int               `bson:"templates,omitempty" json:"templates,omitempty"`
	NVTemplate           int32               `bson:"nvTemplate,omitempty" json:"nvTemplate,omitempty"`
	VideoEndType         int                 `bson:"videoEndType,omitempty" json:"videoEndType,omitempty"`
	RecallNet            *string             `bson:"recallNet,omitempty" json:"recallNet,omitempty"`
	DevIdAllowNull       int                 `bson:"devIdAllowNull,omitempty" json:"devIdAllowNull,omitempty"`
	ImpressionCap        int                 `bson:"impressionCap,omitempty" json:"impressionCap,omitempty"`
	EntranceImg          map[string]string   `bson:"entranceImg,omitempty" json:"entranceImg,omitempty"`
	CookieAchieve        int                 `bson:"cookieAchieve,omitempty" json:"cookieAchieve,omitempty"`
	Hang                 int                 `bson:"hang,omitempty" json:"hang,omitempty"`
	EndcardTemplate      string              `bson:"endcardTemplate,omitempty" json:"endcardTemplate,omitempty"`
	IsServerCall         int                 `bson:"isServerCall,omitempty" json:"isServerCall,omitempty"`
	EntraImage           *string             `bson:"entraImage,omitempty" json:"entraImage,omitempty"`
	RedPointShow         *bool               `bson:"redPointShow,omitempty" json:"redPointShow,omitempty"`
	RedPointShowInterval *int                `bson:"redPointShowInterval,omitempty" json:"redPointShowInterval,omitempty"`
	EntraTitle           string              `bson:"entraTitle,omitempty" json:"entraTitle,omitempty"`
	Alac                 *int                `bson:"alac,omitempty" json:"alac,omitempty"`
	Alecfc               *int                `bson:"alecfc,omitempty" json:"alecfc,omitempty"`
	Mof                  *int                `bson:"mof,omitempty" json:"mof,omitempty"`
	MofUnitId            int64               `bson:"mofUnitId,omitempty" json:"mofUnitId,omitempty"`
	BlackIABCategory     map[string][]string `bson:"blackIABCategory,omitempty" json:"blackIABCategory,omitempty"`
	BlackDomain          []string            `bson:"blackDomain,omitempty" json:"blackDomain,omitempty"`
	BlackBundle          []string            `bson:"blackBundle,omitempty" json:"blackBundle,omitempty"`
	PlacementId          int64               `bson:"plmtId,omitempty" json:"plmtId,omitempty"`
	BiddingType          int                 `bson:"biddingType,omitempty" json:"biddingType,omitempty"`
}

type OutputType struct {
	Template int `bson:"template,omitempty" json:"template,omitempty"`
	Type     int `bson:"type,omitempty" json:"type,omitempty"`
}

type Setting struct {
	ApiRequestNum        int32 `bson:"apiRequestNum,omitempty" json:"apiRequestNum,omitempty"`
	ApiCacheNum          int   `bson:"apiCacheNum,omitempty" json:"apiCacheNum,omitempty"`
	Offset               int   `bson:"mo,omitempty" json:"mo,omitempty"`
	Autoplay             int   `bson:"autoplay,omitempty" json:"autoplay,omitempty"`
	ClickType            int   `bson:"clickType,omitempty" json:"clickType,omitempty"`
	DLNet                int   `bson:"dlnet,omitempty" json:"dlnet,omitempty"`
	ShowImage            int   `bson:"showImage,omitempty" json:"showImage,omitempty"`
	VideoSkipTime        int   `bson:"videoSkipTime,omitempty" json:"videoSkipTime,omitempty"`
	DailyPlayCap         int   `bson:"dailyPlayCap,omitempty" json:"dailyPlayCap,omitempty"`
	CloseButtonDelay     int   `bson:"closeButtonDelay,omitempty" json:"closeButtonDelay,omitempty"`
	VideoInteractiveType int   `bson:"videoInteractiveType,omitempty" json:"videoInteractiveType,omitempty"`
	MuteMode             int   `bson:"muteMode,omitempty" json:"muteMode,omitempty"`
	ReadyRate            int   `bson:"readyRate,omitempty" json:"readyRate,omitempty"`
	AdPlatform           int   `bson:"adPlatform,omitempty" json:"adPlatform,omitempty"`
	GifCtd               int   `bson:"gifCtd,omitempty" json:"gifCtd,omitempty"`
	RefreshFq            int   `bson:"refreshFq,omitempty" json:"refreshFq,omitempty"`
}

func (u *UnitInfo) GetAdSourceID(countryCode string, adSourceID int) int {
	resultList := u.GetAdSourceIDList(countryCode)
	if adSourceID <= 0 {
		if u.UnitId <= int64(0) {
			return mvconst.ADSourceAPIOffer
		}
		adSourceID = mvconst.ADSourceAPIOffer
	}
	if InArray(adSourceID, resultList) {
		return adSourceID
	}
	return 0
}

func (u *UnitInfo) GetAdSourceIDList(countryCode string) []int {
	var adSourceList []int
	key, ok := u.AdSourceCountry[countryCode]
	if !ok {
		// return adSourceList
		if u.AdSourceLen != 0 {
			key = 0
		} else {
			return adSourceList
		}
	}
	if key >= len(u.AdSourceData) {
		return adSourceList
	}
	aList := u.AdSourceData[key]
	for _, info := range aList {
		if info.Status != 1 {
			continue
		}
		adSourceList = append(adSourceList, info.AdSourceId)
	}
	return adSourceList
}

type FakeRule struct {
	Install int     `bson:"install,omitempty" json:"install,omitempty"`
	Start   string  `bson:"start,omitempty" json:"start,omitempty"`
	Status  int     `bson:"status,omitempty" json:"status,omitempty"`
	Updated string  `bson:"updated,omitempty" json:"updated,omitempty"`
	Ecppv   float64 `bson:"ecppv,omitempty" json:"ecppv,omitempty"`
	Type    int     `bson:"type,omitempty" json:"type,omitempty"`
}

type DirectOfferConfig struct {
	Type     int     `bson:"type,omitempty" json:"type,omitempty"`
	Status   int     `bson:"status,omitempty" json:"status,omitempty"`
	OfferIds []int64 `bson:"offerId,omitempty" json:"offerId,omitempty"`
}
