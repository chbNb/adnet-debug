package output

import (
	"bytes"
	"encoding/json"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type UCResult struct {
	UCResponseId string    `json:"id"`
	Cur          string    `json:"cur,omitempty"`
	Seatbid      []Seatbid `json:"seatbid"`
	Bidid        string    `json:"bidid,omitempty"`
	NotEmptyAd   bool      `json:"-"`
}

type Seatbid struct {
	Bid []Bid `json:"bid"`
}

type Bid struct {
	Id            string  `json:"id"`
	Nurl          string  `json:"nurl,omitempty"`
	UnitId        string  `json:"impid"`
	Price         float32 `json:"price"`
	CampaignId    string  `json:"adid,omitempty"`
	Adm           string  `json:"adm,omitempty"`
	FunCampaignId string  `json:"crid,omitempty"`
	Ext           Ext     `json:"ext,omitempty"`
}

type UCAd struct {
}

type AdmData struct {
	Native Native `json:"native"`
}

type Native struct {
	Imptrackers []string `json:"imptrackers"`
	Link        Link     `json:"link"`
	Assets      []Assets `json:"assets"`
}

type Link struct {
	ClickUrl  string   `json:"url"`
	NoticeUrl []string `json:"clicktrackers"`
}

type Assets struct {
	Id    int    `json:"id"`
	Data  *Data  `json:"data,omitempty"`
	Title *Title `json:"title,omitempty"`
	Img   *Img   `json:"img,omitempty"`
}

type Data struct {
	Desc *string `json:"value,omitempty"`
}

type Title struct {
	Title *string `json:"text,omitempty"`
	Len   *int    `json:"len,omitempty"`
}

type Img struct {
	Url    *string `json:"url,omitempty"`
	Width  *int    `json:"w,omitempty"`
	Height *int    `json:"h,omitempty"`
}

type Ext struct {
	PreviewUrl string   `json:"ldp,omitempty"`
	Pm         []Pm     `json:"pm,omitempty"`
	Cm         []string `json:"cm,omitempty"`
}

type Pm struct {
	Point       int    `json:"point"`
	TrackingUrl string `json:"url"`
}

func RenderUCRes(mr MobvistaResult, r *mvutil.RequestParams) UCResult {
	var uc UCResult
	var bid Bid
	var seatbid Seatbid
	uc.UCResponseId = r.Param.UCResponseId
	uc.Cur = "USD"
	for _, v := range mr.Data.Ads {
		bid = renderUcAd(v, r)
		if r.Param.NewRTBFlag && bid.Price > 0.0 {
			uc.NotEmptyAd = true
		}
		seatbid.Bid = append(seatbid.Bid, bid)
		uc.Seatbid = append(uc.Seatbid, seatbid)
	}
	return uc
}

func renderUcAd(ad Ad, r *mvutil.RequestParams) Bid {
	var bid Bid
	var admData AdmData
	bid.Price = 0.0
	if r.Param.NewRTBFlag {
		// uc 切量到实时出价
		if IsOnlineEcpmUnit(&r.Param) {
			bid.Price = float32(ad.OnlineApiBidPrice)
		} else {
			ecppv := extractor.GetUnitFixedEcpm(r.Param.UnitID, r.Param.CountryCode)
			if ecppv > 0.0 {
				bid.Price = float32(ecppv)
			} else {
				mvutil.Logger.Runtime.Warnf("request_id: %s, unit_id: %s, cc: %s, has not config manage revenue.", r.Param.RequestID, bid.UnitId, r.Param.CountryCode)
				watcher.AddWatchValue("unit_no_config_manage_revenue", float64(1))
			}
		}
	}
	bid.Id = "372674938"
	bid.Nurl = ""
	bid.UnitId = strconv.FormatInt(r.Param.UnitID, 10)
	bid.CampaignId = strconv.FormatInt(ad.CampaignID, 10)
	admData = renderAdm(ad)
	//使得特殊字符不转义
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(admData)
	admJson := buffer.Bytes()
	admStr := string(admJson)
	bid.Adm = admStr
	return bid
}

func renderAdm(ad Ad) AdmData {
	var adm AdmData
	adm.Native.Imptrackers = append(adm.Native.Imptrackers, ad.ImpressionURL)
	if len(ad.DeepLink) > 0 {
		adm.Native.Link.ClickUrl = ad.DeepLink
		adm.Native.Link.NoticeUrl = append(adm.Native.Link.NoticeUrl, ad.ClickURL)
	} else {
		adm.Native.Link.ClickUrl = ad.ClickURL
		//统一逻辑到 requestType = online 处理
		//adm.Native.Link.NoticeUrl = append(adm.Native.Link.NoticeUrl, ad.NoticeURL)
	}

	if ad.AdTrackingPoint != nil {
		if len(ad.AdTrackingPoint.Impression) > 0 {
			adm.Native.Imptrackers = append(adm.Native.Imptrackers, ad.AdTrackingPoint.Impression...)
		}
		//ad.AdTrackingPoint.Click = []string{} //测试代码
		if len(ad.AdTrackingPoint.Click) > 0 {
			adm.Native.Link.NoticeUrl = append(adm.Native.Link.NoticeUrl, ad.AdTrackingPoint.Click...)
		}
	}
	if len(adm.Native.Link.NoticeUrl) == 0 {
		adm.Native.Link.NoticeUrl = make([]string, 0)
	}

	// 获取大图的宽高
	imaResulotionList := strings.Split(ad.ImageResolution, "x")
	if len(imaResulotionList) == 1 {
		imaResulotionList = []string{"0", "0"}
	}
	// assets
	ids := []int{1, 2, 3, 4}

	for _, id := range ids {
		var assetsList Assets
		var data Data
		var title Title
		var img Img
		if id == 1 {
			data.Desc = &ad.AppDesc
			assetsList.Data = &data
		} else if id == 2 {
			title.Title = &ad.AppName
			titleLen := len(ad.AppName)
			title.Len = &titleLen
			assetsList.Title = &title
		} else if id == 3 {
			img.Url = &ad.IconURL
			iconSize := 0
			img.Height = &iconSize
			img.Width = &iconSize
			assetsList.Img = &img
		} else if id == 4 {
			img.Url = &ad.ImageURL
			imgWidth, _ := strconv.Atoi(imaResulotionList[0])
			imgHeight, _ := strconv.Atoi(imaResulotionList[1])
			img.Width = &imgWidth
			img.Height = &imgHeight
			assetsList.Img = &img
		}
		assetsList.Id = id
		adm.Native.Assets = append(adm.Native.Assets, assetsList)
	}
	return adm
}
