package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MIGUResult struct {
	Id      string        `json:"id"`    //Ctrl + V
	Bidid   string        `json:"bidid"` //我们的request_id
	Seatbid []SeatbidMIGU `json:"seatbid"`
	Cur     string        `json:"cur"` //Ctrl + V
}

type SeatbidMIGU struct {
	Bid []BidMIGU `json:"bid"`
}

type BidMIGU struct {
	Impid       string     `json:"impid"`
	Price       float64    `json:"price"`
	Nurl        string     `json:"nurl,omitempty"`
	Creative_id string     `json:"creative_id"` //对应的是我们的campaign_id
	Native_ad   *Native_ad `json:"native_ad,omitempty"`
	Video_ad    *Video_ad  `json:"video_ad,omitempty"`
}

type Native_ad struct {
	Imptrackers   []string `json:"imptrackers,omitempty"`
	Clicktrackers []string `json:"clicktrackers,omitempty"`
}

type Video_ad struct {
	Starttrackers    []string `json:"starttrackers,omitempty"`
	Middletrackers   []string `json:"middletrackers,omitempty"`
	Completetrackers []string `json:"completetrackers,omitempty"`
	Clicktrackers    []string `json:"clicktrackers,omitempty"`
}

func RenderMIGURes(mr MobvistaResult, r *mvutil.RequestParams, creative map[string]string) (MIGUResult, bool) {
	var migu MIGUResult
	var seatbid SeatbidMIGU
	//先填充需要回填的值。
	//bidid为Requestid，貌似没有被传递到ad.go?
	migu.Bidid = r.Param.RequestID
	migu.Id = r.Param.MIGUId
	migu.Cur = "CNY" //根据文档，写死。
	if len(creative) > 0 {
		for _, v := range mr.Data.Ads {
			adlist, ok := renderMIGUAdm(v, *r, creative)
			if !ok {
				return migu, false
			} else {
				seatbid.Bid = append(seatbid.Bid, adlist)
			}
		}
		migu.Seatbid = append(migu.Seatbid, seatbid)
		return migu, true
	} else {
		return migu, false
	}

}

func renderMIGUAdm(ad Ad, r mvutil.RequestParams, creative map[string]string) (BidMIGU, bool) {
	var adlist BidMIGU
	//先从creative_id开始
	Creative_id := strconv.FormatInt(ad.CampaignID, 10) //基数指代进制。
	value, ok := creative[Creative_id]
	if !ok {
		mvutil.Logger.Runtime.Warnf("MIGU got unverified ad, requestId:[%s], campaignId:[%d].",
			r.Param.RequestID, ad.CampaignID)
		return adlist, false
	} else {
		adlist.Creative_id = value
	}
	//adlist.Price = 500 //这个价格需要在上线前修改，去MongoDB加一条属性。
	onlinePriceFloorUnits, ifFind := extractor.GetONLINE_PRICE_FLOOR_APPID()
	if ifFind {
		price, ok := onlinePriceFloorUnits[strconv.FormatInt(r.Param.AppID, 10)]
		if !ok {
			adlist.Price = 0
		} else {
			adlist.Price = price
		}
	}
	adlist.Impid = r.Param.MIGUImpId
	//填充native_ad的情况
	if r.Param.AdType == mvconst.ADTypeNative {
		var nativead Native_ad
		nativead.Imptrackers = append(nativead.Imptrackers, ad.ImpressionURL)
		nativead.Clicktrackers = append(nativead.Clicktrackers, ad.ClickURL)
		adlist.Native_ad = &nativead
	} else if r.Param.AdType == mvconst.ADTypeOnlineVideo {
		var videoad Video_ad
		for _, v := range ad.AdTracking.Play_percentage {
			switch v.Rate {
			//视频开始播放，将imrpssion和starttrackers放一起。
			case 0:
				videoad.Starttrackers = append(videoad.Starttrackers, ad.ImpressionURL, v.Url)
			case 50:
				videoad.Middletrackers = append(videoad.Middletrackers, v.Url)
			case 100:
				videoad.Completetrackers = append(videoad.Completetrackers, v.Url)
			}
		}
		videoad.Clicktrackers = append(videoad.Clicktrackers, ad.ClickURL)
		adlist.Video_ad = &videoad
	}
	return adlist, true
}
