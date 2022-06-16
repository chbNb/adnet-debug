package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type pptvResponse struct {
	ID      string        `json:"id"`
	Seatbid []pptvSeatbid `json:"seatbid"`
}

type pptvSeatbid struct {
	PPTVBid []pptvBid `json:"bid"`
	Cur     string    `json:"cur"`
}

type pptvBid struct {
	ID     string   `json:"id"`
	ImpId  string   `json:"impid"`
	Price  float64  `json:"price"`
	Crid   string   `json:"crid"`
	Vcurl  []string `json:"vcurl"`
	Vsurl  []string `json:"vsurl"`
	Event  []Event  `json:"event"`
	DealId string   `json:"dealid"`
}

type Event struct {
	Name     string `json:"name"`
	Trackurl string `json:"trackurl"`
}

func RenderPPTVRes(mr MobvistaResult, r *mvutil.RequestParams) (pptvResponse, error) {
	var pptvRsp pptvResponse
	pptvRsp.ID = r.Param.OnlineReqId
	pptvSeatbid, err := renderSeatbid(mr, r)
	pptvRsp.Seatbid = append(pptvRsp.Seatbid, pptvSeatbid)
	return pptvRsp, err
}

func renderSeatbid(mr MobvistaResult, r *mvutil.RequestParams) (pptvSeatbid, error) {
	var seatbid pptvSeatbid
	seatbid.Cur = "CNY"
	bid, err := renderBid(mr, r)
	seatbid.PPTVBid = append(seatbid.PPTVBid, bid)
	return seatbid, err
}

func renderBid(mr MobvistaResult, r *mvutil.RequestParams) (pptvBid, error) {
	var bid pptvBid
	crIdConf, ok := extractor.GetCREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS()
	if !ok {
		return bid, errorcode.EXCEPTION_RETURN_EMPTY
	}
	bid.ID = r.Param.RequestID
	bid.ImpId = "1"
	bid.Price = float64(r.Param.MaxPrice)
	for _, v := range mr.Data.Ads {
		camId := strconv.FormatInt(v.CampaignID, 10)
		if crId, ok := crIdConf[camId]; ok {
			bid.Crid = crId
		} else {
			mvutil.Logger.Runtime.Warnf("PPTV got unverified ad, requestId:[%s], campaignId:[%s].",
				r.Param.RequestID, camId)
			return bid, errorcode.EXCEPTION_RETURN_EMPTY
		}
		dealId, ifFind := extractor.GetPPTV_DEAL_ID()
		if ifFind {
			bid.DealId = dealId
		}
		bid.Vcurl = append(bid.Vcurl, v.ClickURL)
		bid.Vsurl = append(bid.Vsurl, v.ImpressionURL)
		bid.Event = renderEvent(v)
	}
	return bid, nil
}

func renderEvent(ad Ad) []Event {
	var event Event
	var eventArr []Event
	for _, v := range ad.AdTracking.Play_percentage {
		switch v.Rate {
		case 0:
			event.Name = "start"
		case 25:
			event.Name = "firstQuartile"
		case 50:
			event.Name = "midpoint"
		case 75:
			event.Name = "thirdQuartile"
		case 100:
			event.Name = "complete"
		}
		event.Trackurl = v.Url
		eventArr = append(eventArr, event)
	}
	return eventArr
}
