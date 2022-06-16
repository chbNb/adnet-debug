package corsair_proto

import "gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"

type BackendAds struct {
	BackendId           int32
	RequestKey          string
	CampaignList        []*Campaign
	Strategy            string
	RunTimeVariables    *RunTimeVariable
	AdTemplate          *ad_server.ADTemplate
	NoneResultReason    *NoneResultReason
	FilterReason        []*FilterReason
	FrameList           []*Frame
	AlgoFeatInfo        *string
	IfLowerImp          *int32
	IsAdServerTest      *int32
	ResourceType        *ad_server.InteractiveResourceType
	EndScreenTemplateId *int32
	BannerHtml          *string
	BannerUrl           *string
	BigTemplateInfo     *BigTemplate
	ExtAdxAlgo          string
	EcpmFloor           float64
	DspId               int64
	RKS                 map[string]string
}

func NewBackendAds() *BackendAds {
	return &BackendAds{}
}
func (p *BackendAds) GetCampaignList() []*Campaign {
	return p.CampaignList
}

func (p *BackendAds) GetBannerHtml() string {
	if p.BannerHtml == nil {
		return ""
	}
	return *p.BannerHtml
}
func (p *BackendAds) GetBannerUrl() string {
	if p.BannerUrl == nil {
		return ""
	}
	return *p.BannerUrl
}
