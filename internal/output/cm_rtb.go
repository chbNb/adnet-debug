package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type CMResult struct {
	Id      string      `json:"id,omitempty"`    // Paste from param
	Bidid   string      `json:"bidid,omitempty"` // request_id or extra5?
	Cur     string      `json:"cur"`             // CNY
	SeatBid []CMSeatbid `json:"seatbid"`
}

type CMSeatbid struct {
	Bid []CMBid `json:"bid"`
}

type CMBid struct {
	Id     string  `json:"id,omitempty"`    // request_id ?
	ImpId  string  `json:"impid,omitempty"` // Paste from param
	Price  float64 `json:"price"`           // map price on mongo
	Adid   string  `json:"adid"`            // map campaignid on mongo
	Bundle string  `json:"bundle"`          // pkgmame
	Adm    string  `json:"adm"`             // vast xml
	W      int32   `json:"w"`               // video_w
	H      int32   `json:"h"`               // video_h
	Nurl   string  `json:"nurl,omitempty"`  // win notice, still talking
}

func RenderCMRes(mr MobvistaResult, r *mvutil.RequestParams, vastxml []byte) (CMResult, bool) {
	var CM CMResult
	var seatbid CMSeatbid
	// start from Id
	CM.Id = r.Param.CMId
	CM.Bidid = r.Param.RequestID // 是否保持一致？
	CM.Cur = "CNY"
	for _, v := range mr.Data.Ads {
		adlist, ok := RenderCMBid(v, *r, vastxml)
		if !ok {
			return CM, false
		} else {
			seatbid.Bid = append(seatbid.Bid, adlist)
		}
	}
	CM.SeatBid = append(CM.SeatBid, seatbid)
	return CM, true
}

// mock a price, creativeid unsolved
func RenderCMBid(ad Ad, r mvutil.RequestParams, vastxml []byte) (CMBid, bool) {
	var adlist CMBid
	// start from ImpId
	adlist.ImpId = r.Param.CMImpId
	adlist.Id = r.Param.RequestID
	//adlist.Price = 5 // 元/CPM！
	// 单位元/CPM
	onlinePriceFloorUnits, ifFind := extractor.GetONLINE_PRICE_FLOOR_APPID()
	if !ifFind {
		return adlist, false
	}
	adlist.Price = (onlinePriceFloorUnits[strconv.FormatInt(r.Param.AppID, 10)] / 100) //单位转换
	adlist.Adid = strconv.FormatInt(ad.CampaignID, 10)
	// adm is using vast xml
	adlist.Adm = string(vastxml)
	if adlist.Adm != "" {
		adlist.Bundle = ad.PackageName
		// 兼容第三方无下发videoHeight和videoWidth的情况
		compatibleVideoHW(&ad)
		adlist.W = ad.VideoWidth
		adlist.H = ad.VideoHeight
		return adlist, true
	}
	return adlist, false
}
