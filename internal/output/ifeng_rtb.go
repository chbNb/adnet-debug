package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type IFENGResult struct {
	Id  string    `json:"id"`  // ifengid
	CId string    `json:"cid"` // requestid
	Ad  []IFENGAd `json:"ad"`
}

type IFENGAd struct {
	Id       string        `json:"id"`    // ifengtagid
	ImpId    string        `json:"impid"` // ifengimpid
	AdId     string        `json:"adid"`  // campaign_id
	Creative IFENGCreative `json:"creative"`
}

type IFENGCreative struct {
	Type int    `json:"type"` // 写死3
	Vast string `json:"vast"` // using cheetah solution...
}

func RenderIFENGRes(mr MobvistaResult, r *mvutil.RequestParams, vastxml []byte) IFENGResult {
	var IFENG IFENGResult
	// start from Id
	IFENG.Id = r.Param.IFENGId
	IFENG.CId = r.Param.RequestID
	for _, v := range mr.Data.Ads {
		adlist := RenderIFENGCreative(v, *r, vastxml)
		IFENG.Ad = append(IFENG.Ad, adlist)
	}
	return IFENG
}

func RenderIFENGCreative(ad Ad, r mvutil.RequestParams, vastxml []byte) IFENGAd {
	var adlist IFENGAd
	// start from Id
	adlist.Id = r.Param.IFENGTagId
	adlist.ImpId = r.Param.IFENGImpId
	adlist.AdId = strconv.FormatInt(r.Param.IFENGAdId, 10)
	// filter vast
	adlist.Creative.Type = 3
	adlist.Creative.Vast = string(vastxml)
	return adlist
}
