package params

import (
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type AsAdData struct {
	AdTemplate          int64
	AlgoFeatInfo        string
	IfLowerImp          int32
	ResourceType        int64
	EndScreenTemplateId int32
	CampaignId          int64
	AdSource            int64
	AdAdTemplate        int64
	ImageSizeId         int64
	OfferType           int32
	BtType              int64
	CreativeId          string
	Downloadtest        int32
	AdElementTemplate   int64
	CreativeId2         string
	CreativeMap         map[ad_server.CreativeType]int64
	DcoCrMap            map[ad_server.CreativeType]map[string]int64
	CreativeTypeIdMap   map[ad_server.CreativeType]int64
	CreativeTypeIdMap2  map[ad_server.CreativeType]int64
	DynamicCreative     map[ad_server.CreativeType]string
	Playable            bool
	EndcardUrl          string
	ExtPlayable         int32
	VideoEndType        int32
	Orientation         int64
	UsageVideo          bool
	TemplateGroup       int64
	VideoTemplateId     int64
	EndCardTemplateId   int64
	MiniCardTemplateId  int64
	BidPrice            float64
	RankType            int64
	Strategy            string
	CreativeDataMap     map[string]map[ad_server.CreativeType]int64
	BigTemplateId       int64
	BigTemplateSlotMap  map[int32]int64
	CpdIds              []string
}
