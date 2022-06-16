package protocol

import (
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"
)

func Openrtbv2ToMtgrtbResp(b []byte) (*mtgrtb.BidResponse, error) {
	oResp := new(openrtb.BidResponse)
	err := proto.Unmarshal(b, oResp)
	if err != nil {
		return nil, err
	}
	mtgResp := &mtgrtb.BidResponse{
		Id:    proto.String(oResp.Id),
		Bidid: proto.String(oResp.Bidid),
		Cur:   proto.String(oResp.Cur),
	}
	nbr := mtgrtb.BidResponse_NoBidReason(oResp.Nbr)
	mtgResp.Nbr = &nbr
	if oResp.Mvext != nil {
		mtgResp.RespAs = &oResp.Mvext.RespAs
		mtgResp.VastMas = &oResp.Mvext.VastMas
		b, _ := jsoniter.ConfigFastest.Marshal(oResp.Mvext.AsResp)
		asResp := new(mtgrtb.BidResponse_AsResp)
		jsoniter.ConfigFastest.Unmarshal(b, asResp)
		mtgResp.AsResp = asResp
	}
	mtgResp.Ext = &mtgrtb.BidResponse_Ext{
		RepGrayTags:        oResp.Mvext.RepGrayTags,
		OnlyRequstThirdDsp: &oResp.Mvext.OnlyRequstThirdDsp,
		Rks:                oResp.Mvext.Rks,
	}
	mtgResp.Seatbid = make([]*mtgrtb.BidResponse_SeatBid, len(oResp.Seatbid))
	for i, seatbid := range oResp.Seatbid {
		mtgResp.Seatbid[i] = &mtgrtb.BidResponse_SeatBid{
			Bid:  make([]*mtgrtb.BidResponse_SeatBid_Bid, len(seatbid.Bid)),
			Seat: &seatbid.Seat,
		}
		for j, bid := range seatbid.Bid {
			mtgResp.Seatbid[i].Bid[j] = &mtgrtb.BidResponse_SeatBid_Bid{
				Id:      proto.String(bid.Id),
				Impid:   proto.String(bid.Impid),
				Price:   proto.Float64(bid.Price),
				Adid:    proto.String(bid.Adid),
				Adm:     proto.String(bid.Adm),
				Adomain: bid.Adomain,
				Bundle:  proto.String(bid.Bundle),
				Nurl:    proto.String(bid.Nurl),
				Iurl:    proto.String(bid.Iurl),
				Cid:     proto.String(bid.Cid),
				Crid:    proto.String(bid.Crid),
				Cat:     bid.Cat,
				//Attr:                   bid.Attr,
				//Api:                    nil,
				//Protocol:               nil,
				//Qagmediarating:         nil,
				Imptracker:   proto.String(bid.Imptracker),
				Clicktracker: proto.String(bid.Clicktracker),
				Burl:         proto.String(bid.Burl),
			}
			qag := mtgrtb.QAGMediaRating(bid.Qag)
			mtgResp.Seatbid[i].Bid[j].Qagmediarating = &qag
			if bid.Mvext != nil {
				mtgResp.Seatbid[i].Bid[j].Ext = &mtgrtb.BidResponse_SeatBid_Bid_Ext{
					Dataext: &bid.Mvext.Dataext,
				}
			}
		}
	}
	//b, _ = json.Marshal(mtgResp)
	//fmt.Println("mtgreq:" ,string(b))

	return mtgResp, nil
}
