package backend

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	rtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
)

type PioneerBackend struct {
	Backend
}

func (backend PioneerBackend) composeHttpRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, req *fasthttp.Request) error {
	if reqCtx == nil || backendCtx == nil || req == nil {
		return errors.New("pioneerBackend composeHttpRequest params invalidate")
	}
	bidRequest, err := renderBidRequest(reqCtx, backendCtx)
	if err != nil {
		return err
	}
	protoData, err := proto.Marshal(bidRequest)
	if err != nil {
		return err
	}
	req.Header.SetMethod(backendCtx.Method)
	req.SetRequestURI(backendCtx.ReqPath)
	req.SetBody(protoData)
	return nil
}

func (backend PioneerBackend) filterBackend(reqCtx *mvutil.ReqCtx) int {
	return mvconst.BackendOK
}

func (backend *PioneerBackend) getRequestNode() string {
	return backend.PioneerClient.GetNode() + "/bid"
}

func (backend *PioneerBackend) parseHttpResponse(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (int, error) {
	if reqCtx == nil || backendCtx == nil {
		return ERR_Param, errors.New("MAdxBackend parseHttpResponse params invalidate")
	}
	var rdata rtb.BidResponse
	err := proto.Unmarshal(backendCtx.RespData, &rdata)
	if err != nil {
		watcher.AddWatchValue("mAdx_unmarshal_error", float64(1))
		return ERR_ParseRsp, err
	}
	// 解析返回的协议字段
	GetInfoFromRData(rdata, reqCtx)

	if len(rdata.Seatbid) == 0 {
		watcher.AddWatchValue("pioneer_no_ads", float64(1))
		return ERR_NoAds, errors.New("return no ads, result.id=" + *rdata.Id)
	}
	reqCtx.RespData = backendCtx.RespData
	err = fillAd(&rdata, reqCtx, backendCtx)
	if err != nil {
		return ERR_ParseRsp, err
	}

	// b, _ := json.Marshal(backendCtx.Ads)
	// mvutil.Logger.Runtime.Info("ads:" + string(b))
	return ERR_OK, nil
}
