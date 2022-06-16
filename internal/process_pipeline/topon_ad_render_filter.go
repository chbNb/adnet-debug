package process_pipeline

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"
	"strings"
)

type ToponAdRenderFilter struct {
}

func (tf *ToponAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("ToponAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	result, err := output.RenderOutput(in.ReqParams, in.Result)
	in.ReqParams.Param.Extra = "topon-nothb"
	in.ReqParams.Param.Algorithm = "topon-nothb"
	mvutil.StatRequestLog(&in.ReqParams.Param)
	//降填充日志
	mvutil.StatReduceFillLog(&in.ReqParams.Param)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] topon RenderOutput err: %s", in.ReqParams.Param.RequestID, err.Error())
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}
	res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)

	if err != nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] result=[%+v] to json error : %s", in.ReqParams.Param.RequestID, *result, err.Error())
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}

	if in.ReqParams.ToponResponse == nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] ToponResponse nil", in.ReqParams.Param.RequestID)
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}
	if len(in.ReqParams.Param.ReqBackends) == 1 && len(in.ReqParams.Param.BackendReject) == 1 {
		parts := strings.Split(in.ReqParams.Param.BackendReject[0], ":")
		if len(parts) == 2 && parts[1] != "2000" { // 2000: BackendOK
			mvutil.Logger.Runtime.Warnf("request_id=[%s] Topon Rejected:[%s]", in.ReqParams.Param.RequestID, in.ReqParams.Param.BackendReject[0])
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
	}
	if in.ReqParams.Param.ToponTemplateSupportVersion != 0 && in.ReqParams.ToponResponse.Mvext != nil {
		in.ReqParams.ToponResponse.Mvext.SdkResp = string(res)
	}
	csp := mvutil.SerializeCSP(&in.ReqParams.Param, in.ReqParams.Param.DspExt)
	if len(in.ReqParams.ToponResponse.Seatbid) > 0 && len(in.ReqParams.ToponResponse.Seatbid[0].Bid) > 0 {
		m := map[string]string{
			"csp": csp,
			"do":  in.ReqParams.Param.Domain,
			"c":   mvutil.SerializeC(&in.ReqParams.Param),
		}
		b, _ := jsoniter.ConfigFastest.Marshal(&m)
		if in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext == nil {
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext = &openrtb.BidResponse_SeatBid_Bid_Ext{}
		}
		in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Dataext = string(b)
		in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Linktype = openrtb.LinkType(in.ReqParams.Param.LinkType)
		if in.ReqParams.AsResp != nil {
			impressionUrl := output.CreateImpressionUrl(&in.ReqParams.Param)
			clickUrl := output.CreateClickUrl(&in.ReqParams.Param, true)
			// topon 针对pioneer情况，记录bid_price
			if len(in.ReqParams.Param.OnlineApiBidPrice) > 0 && in.ReqParams.Param.OnlineApiBidPrice != "0.00000000" {
				bidPrice := "&bid_p=" + in.ReqParams.Param.OnlineApiBidPrice
				impressionUrl += bidPrice
				clickUrl += bidPrice
			}
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers = append(
				in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers, impressionUrl)
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers = append(
				in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers, output.CreateOnlyImpressionUrl(in.ReqParams.Param, in.ReqParams))
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Clicktrackers = append(
				in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Clicktrackers, clickUrl)
		} else {
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers = append(
				in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Imptrackers, in.ReqParams.Param.ToponThirdPartyImpUrl)
			in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Clicktrackers = append(
				in.ReqParams.ToponResponse.Seatbid[0].Bid[0].Mvext.Clicktrackers, in.ReqParams.Param.ToponThirdPartyNoticeUrl)
		}
	}

	b, _ := proto.Marshal(in.ReqParams.ToponResponse)
	in.RespData = b
	resp := string(in.RespData)
	return &resp, nil
}
