package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type HupuResult struct {
	UCResponseId string        `json:"id"`
	Bidid        string        `json:"bidid,omitempty"`
	Cur          string        `json:"cur,omitempty"`
	Seatbid      []HupuSeatbid `json:"seatbid"`
}

type HupuSeatbid struct {
	Bid []HupuBid `json:"bid"`
}

type HupuBid struct {
	Id         string `json:"id"`
	UnitId     int    `json:"impid"`
	Price      int    `json:"price"`
	CampaignId string `json:"adid,omitempty"`
	//Adm           string  `json:"adm,omitempty"`
	FunCampaignId int64      `json:"crid,omitempty"`
	DealId        string     `json:"dealid,omitempty"`
	Ext           HupuExt    `json:"ext,omitempty"`
	Native        HupuNative `json:"native,omitempty"`
}

type HupuExt struct {
	//PreviewUrl string   `json:"ldp,omitempty"`
	Pm      []string `json:"pm,omitempty"`
	Cm      []string `json:"cm,omitempty"`
	Em      []string `json:"em,omitempty"`
	IconUrl string   `json:"icon_url,omitempty"`
}

type HupuNative struct {
	Link  HupuLink  `json:"link,omitempty"`
	Title string    `json:"title,omitempty"`
	Image HupuImage `json:"image,omitempty"`
	Video HupuVideo `json:"video,omitempty"`
}

type HupuImage struct {
	Url    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type HupuVideo struct {
	Url      string `json:"url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Size     int    `json:"size,omitempty"`
}

type HupuLink struct {
	Ldp      string `json:"ldp,omitempty"`
	LdpType  int    `json:"ldptype"`
	DeepLink string `json:"dp,omitempty"`
}

func RenderHupuRes(mr MobvistaResult, r *mvutil.RequestParams, creative []int64) (HupuResult, bool) {
	var hupu HupuResult
	var seatbid HupuSeatbid

	hupu.UCResponseId = r.Param.HupuRequestId
	hupu.Bidid = r.Param.RequestID
	for _, v := range mr.Data.Ads {
		bid, ok := renderHupuAd(v, r, creative)
		if !ok {
			return hupu, false
		}
		seatbid.Bid = append(seatbid.Bid, bid)
	}
	hupu.Seatbid = append(hupu.Seatbid, seatbid)
	return hupu, true
}

func renderHupuAd(ad Ad, r *mvutil.RequestParams, creative []int64) (HupuBid, bool) {
	var bid HupuBid
	// 先定位requestid
	bid.Id = r.Param.RequestID
	bid.Price = getHupuPrice()
	// 这里指的是曝光ID
	bid.UnitId, _ = strconv.Atoi(r.Param.HupuImpId)
	//bid.Adm = mvutil.Base64Decode(ad.VideoURL)

	// 针对unit切量，切成免审模式。
	if !mvutil.InInt64Arr(ad.CampaignID, creative) && IsNeedReviewCreativeUnit(r.Param.UnitID) {
		mvutil.Logger.Runtime.Warnf("HUPU got unverified ad, requestId:[%s], campaignId:[%d].",
			r.Param.RequestID, ad.CampaignID)
		return bid, false
	}
	if !IsNeedReviewCreativeUnit(r.Param.UnitID) && len(ad.ImageURL) == 0 && IsHupuSplash(r.Param.DealId) {
		mvutil.Logger.Runtime.Warnf("HUPU image url empty, requestId:[%s], campaignId:[%d].",
			r.Param.RequestID, ad.CampaignID)
		return bid, false
	}
	bid.FunCampaignId = ad.CampaignID
	bid.DealId = r.Param.DealId
	bid.Ext = renderHupuExt(ad, r)
	if !IsNeedReviewCreativeUnit(r.Param.UnitID) {
		bid.Native = renderNative(ad, r)
	}
	return bid, true
}

func IsNeedReviewCreativeUnit(unitId int64) bool {
	adnConfList := extractor.GetADNET_CONF_LIST()
	if hupuNeedReviewCreativeUnitList, ok := adnConfList["hupuNeedReviewCreativeUnitList"]; ok && mvutil.InInt64Arr(unitId, hupuNeedReviewCreativeUnitList) {
		return true
	}
	return false
}

func renderNative(ad Ad, r *mvutil.RequestParams) HupuNative {
	var native HupuNative
	var hupuLink HupuLink
	var hupuImage HupuImage
	var hupuVideo HupuVideo
	native.Title = ad.AppName
	hupuImage.Url = ad.ImageURL
	imaResulotionList := strings.Split(ad.ImageResolution, "x")
	if len(imaResulotionList) == 2 {
		hupuImage.Width, _ = strconv.Atoi(imaResulotionList[0])
		hupuImage.Height, _ = strconv.Atoi(imaResulotionList[1])
	}
	// 兼容第三方无下发videoHeight和videoWidth的情况
	compatibleVideoHW(&ad)
	hupuVideo.Url = mvutil.Base64Decode(ad.VideoURL)
	hupuVideo.Duration = ad.VideoLength
	hupuVideo.Size = ad.VideoSize
	hupuVideo.Height = int(ad.VideoHeight)
	hupuVideo.Width = int(ad.VideoWidth)
	native.Image = hupuImage
	native.Video = hupuVideo

	hupuLink.Ldp = ad.ClickURL
	hupuLink.LdpType = 0
	if len(ad.DeepLink) > 0 {
		hupuLink.DeepLink = ad.DeepLink
		hupuLink.Ldp = ad.PreviewUrl
		// 对于切rtdsp的情况下，没有
		if len(hupuLink.Ldp) == 0 && r.DspExtData != nil && mvutil.InInt64Arr(r.DspExtData.DspId, []int64{mvconst.MVDSP_Retarget, mvconst.MAS}) {
			hupuLink.Ldp = ad.ClickURL
		}
	}

	native.Link = hupuLink
	return native
}

func renderHupuExt(ad Ad, r *mvutil.RequestParams) HupuExt {
	var ext HupuExt
	ext.Pm = append(ext.Pm, ad.ImpressionURL)
	// 虎扑支持切量adx
	if ad.AdTrackingPoint != nil && len(ad.AdTrackingPoint.Impression) > 0 {
		ext.Pm = append(ext.Pm, ad.AdTrackingPoint.Impression...)
	}

	//complateUrl := ""
	for _, v := range ad.AdTracking.Play_percentage {
		if v.Rate == 100 {
			//complateUrl = v.Url
			ext.Em = append(ext.Em, v.Url) // 支持adnTracking  & adxTracking
		}
	}
	//ext.Em = append(ext.Em, complateUrl)
	if IsNeedReviewCreativeUnit(r.Param.UnitID) {
		ext.Cm = append(ext.Cm, ad.ClickURL)
		ext.IconUrl = ad.IconURL
	}

	// 虎扑支持切量adx
	if ad.AdTrackingPoint != nil && len(ad.AdTrackingPoint.Click) > 0 {
		ext.Cm = append(ext.Cm, ad.AdTrackingPoint.Click...)
	}
	// 免审模式下
	if !IsNeedReviewCreativeUnit(r.Param.UnitID) {
		// adnet->as,adnet->adx->as召回deeplink单子的情况下，需要把click_url放到cm中，不然无法获取点击上报
		if len(ad.DeepLink) > 0 && (r.DspExtData == nil || mvutil.InInt64Arr(r.DspExtData.DspId, []int64{mvconst.FakeAdserverDsp, 0})) {
			ext.Cm = append(ext.Cm, ad.ClickURL)
		}
	}

	return ext
}

func getHupuPrice() int {
	conf, ok := extractor.GetHUPU_DEFAULT_PRICE()
	if !ok {
		mvutil.Logger.Runtime.Warnf("Hupu get HUPU_DEFAULT_PRICE error")
	}
	price := 0
	if conf > 0 {
		price = conf * 100
	}
	return price
}

func IsHupuSplash(dealid string) bool {
	listConf := extractor.GetADNET_CONF_LIST()

	if hupuSplashDealIdList, ok := listConf["hupuSplashDealIdList"]; ok && len(hupuSplashDealIdList) > 0 {
		dealidInt, err := strconv.ParseInt(dealid, 10, 64)
		if err != nil {
			return false
		}
		return mvutil.InInt64Arr(dealidInt, hupuSplashDealIdList)
	}
	return dealid == "10059" || dealid == "10058" || dealid == "10260" || dealid == "10261"
}
