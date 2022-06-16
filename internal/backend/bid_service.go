package backend

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"

	jsoniter "github.com/json-iterator/go"

	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	adxmodel "gitlab.mobvista.com/ADN/adx_common/model"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/vast"
)

// bid server 的相关返回常量
const (
	// mtg response ext abTest 字段中 bid server
	BidServerAdxABTestKey = "BidServer"
)

// bid server error code
const (
	// 成功出价
	BidServerErrorCodeSuccess = 1
	// 放弃出价: 低价值流量被过滤
	BidServerErrorCodeFilterLowValueReq = 2
	// 请求回退: 已过滤, 认为信息不够回退到pioneer处理
	BidServerErrorCodeRedirectPioneer = 3
	// 请求回退: 未过滤, 内部错误回退到pioneer处理
	BidServerErrorCodeInnerError = 4
	// 放弃出价: 认为流量有亏损
	BidServerErrorCodeLowProfit = 5
)

var (
	LoadHttpError                  = errors.New("load_http") // 请求 bid server 返回的 Load Url 错误
	BidHttpError                   = errors.New("bid_http")  // 请求 bid 的 http 错误
	BidHttpTimeoutError            = errors.New("bid_http_timeout")
	LoadUnmarshallError            = errors.New("load_unmarshall") // 请求 bid server 返回的 load url 反序列化失败
	RenderAdsDspExtNil             = errors.New("dspext_nil")      // 渲染广告的时候 dsp ext 是空
	RenderAdsNotMatchMas           = errors.New("not_match_mas")   // 通过 Load url 渲染广告的时候发现此时 dspId 不是 mas
	RenderAdxTrackerEmpty          = errors.New("ads_render_adx_tk_empty")
	RenderInParamEmpty             = errors.New("ads_render_in_param_empty")
	RenderAdsNoAds                 = errors.New("ads_render_no_ads")
	RenderAdsUnmarshallDspExtError = errors.New("ads_render_unmarshall_dsp_ext_error")
	RenderAdsNoImp                 = errors.New("ads_render_no_imp")
	RenderAdsAdxOriginRespEmpty    = errors.New("ads_render_adx_origin_imp_empty")
	RenderParseVast                = errors.New("ads_render_vast_parse_error")
	RenderVastNoInline             = errors.New("ads_render_vast_no_line")
	RenderDecodeVast               = errors.New("ads_render_decode_vast")
	RenderFillAd                   = errors.New("ads_render_fill_ad")
	RenderAdNotVast                = errors.New("ads_render_ad_not_vast")
	BidBidServerWinButNoAds        = errors.New("bid_bs_win_no_ads")
)

// Load 请求 pioneer
// 二阶段 bid server 成功返回后的 load 阶段
// 此时会请求 pioneer 通过 load_url 获取素材信息
func RequestLoadUrl(LoadUrl string) (*mtgrtb.BidResponse, error) {

	if LoadUrl == "" {
		return nil, errors.New("bid server load url is empty")
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(LoadUrl)
	req.Header.SetMethod("GET")

	// 记录Load时延
	now := time.Now()
	err := utility.HttpClientApp().Do(getLoadTimeout(), req, resp)
	costTime := (float64)(time.Since(now) / time.Millisecond)
	// metrics.SetGaugeWithLabelValues(costTime, 35, "load")
	metrics.AddSummaryWithLabelValues(costTime, 37, "load") // 统计响应分位

	//
	if err != nil {
		mvutil.Logger.Runtime.Warnf("[bid-server] request bid server load url:[%s] error: [%s]", LoadUrl, err.Error())
		return nil, LoadHttpError
	}

	// 不是 err 才拿状态码
	metrics.IncCounterWithLabelValues(34, strconv.Itoa(resp.StatusCode()))

	bidResp := new(mtgrtb.BidResponse)
	err = proto.Unmarshal(resp.Body(), bidResp)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("[bid-server] request bid server load url:[%s] unmarshall error:[%s]", LoadUrl, err.Error())
		return nil, LoadUnmarshallError
	}

	// success
	return bidResp, nil
}

//
func RenderLoadResult(val *storage.ReqCtxVal, reqCtx *mvutil.ReqCtx) (*mvutil.BackendCtx, error) {
	dspExt, err := reqCtx.ReqParams.GetDspExt()
	if err != nil {
		return nil, RenderAdsDspExtNil
	}
	if dspExt.DspId != mvconst.MAS {
		// 如果走有这里, 但是胜出的dsp又不是mas, 说明是有问题的
		return nil, RenderAdsNotMatchMas
	}

	adxOrigonResp := val.ReqParams.Param.BidServerAdxResponse
	loadUrl := val.ReqParams.Param.BidServerCtx.LoadUrl
	bidResp, err := RequestLoadUrl(loadUrl)
	if err != nil {
		// 请求 Load url 失败
		return nil, err
	}

	// 构造一个 backendCtx
	backendCtx := mvutil.NewBackendCtx("", "", "", 1.00, 1, []string{"ALL"})
	backendCtx.Ads.BackendId = mvconst.MAdx

	err = ParseLoadResponse(bidResp, adxOrigonResp, reqCtx, backendCtx)
	if err != nil {
		// 解析 pioneer 放回的素材失败
		mvutil.Logger.Runtime.Warnf("[bid-server] parse pioneer resp ads error:[%s]", err.Error())
		return nil, err
	}

	return backendCtx, nil
}

//
func ParseLoadResponse(bid *mtgrtb.BidResponse, adxOriginResp *mtgrtb.BidResponse, reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) error {
	if bid == nil || reqCtx == nil || backendCtx == nil || adxOriginResp == nil {
		return RenderInParamEmpty
	}

	if len(bid.Seatbid) == 0 {
		return RenderAdsNoAds
	}

	// 解析返回的协议字段
	GetInfoFromRData(*bid, reqCtx)

	// 通过 adxOriginResp 补齐 bidResp 中缺失的字段
	target, err := getBidByImpId(bid.GetSeatbid(), ImpId)
	if err != nil {
		// load url 返回值没有 bid
		return RenderAdsNoImp
	}
	//
	origin, err := getBidByImpId(adxOriginResp.GetSeatbid(), ImpId)
	if err != nil || origin.GetExt() == nil {
		return RenderAdsAdxOriginRespEmpty
	}

	//
	if target.Ext == nil {
		target.Ext = &mtgrtb.BidResponse_SeatBid_Bid_Ext{}
	}
	// 补齐 DataExt 中的 bd & cid, andtk用的csp参数会用到Dataext
	dspExt := adxmodel.DspExt{}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(origin.GetExt().GetDataext()), &dspExt)
	if err != nil {
		return RenderAdsUnmarshallDspExtError
	}
	dspExt.Bundle = target.GetBundle()
	dspExt.CampaignId = target.GetCid()
	dspExtBytes, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(dspExt)
	target.Ext.Dataext = proto.String(string(dspExtBytes))
	target.Adid = proto.String(fmt.Sprintf("%d-%s", mvconst.MAS, target.GetAdid()))
	target.Cid = proto.String(fmt.Sprintf("%d-%s", mvconst.MAS, target.GetCid()))
	target.Crid = proto.String(fmt.Sprintf("%d-%s", mvconst.MAS, target.GetCrid()))
	// adx 连接模版参数化逻辑补齐
	if origin.GetExt() != nil && origin.GetExt().GetAks() != nil {
		target.Ext.Aks = origin.GetExt().GetAks()
	}

	// adx aks 参数化逻辑补齐
	if adxOriginResp.GetExt() != nil && adxOriginResp.GetExt().GetRks() != nil {
		if bid.Ext == nil {
			bid.Ext = &mtgrtb.BidResponse_Ext{}
		}
		bid.Ext.Rks = adxOriginResp.Ext.Rks
	}

	// rs 精准出价落日志
	// 这里需要美元转下美分
	reqCtx.ReqParams.Param.ExtDataInit.BidServerRsPrice = target.GetPrice() * 100

	// 将price替换为bid的出价
	// 聚合平台会用 bid 阶段的出价结算, adntk则用Load阶段链接中的价格结算
	// 两边的价格需要一样, 不然账对不齐
	target.Price = proto.Float64(reqCtx.ReqParams.Price)

	//
	if reqCtx.ReqParams.Param.BidServerCtx.AdxTrackers == nil {
		return RenderAdxTrackerEmpty
	}

	//
	impTk := reqCtx.ReqParams.Param.BidServerCtx.AdxTrackers.ImpressionTracker
	clkTk := reqCtx.ReqParams.Param.BidServerCtx.AdxTrackers.ClickTracker
	trackEvents := reqCtx.ReqParams.Param.BidServerCtx.AdxTrackers.TrackingEvents

	// 解析 bid adm vast (pioneer放回的都是vast), 回填 adxtk 连接
	adm := target.GetAdm()
	if mvutil.IsAdmVast(adm) {
		bidVast, err := decodeVast(adm)
		if err != nil {
			return RenderParseVast
		}
		if len(bidVast.Ads) == 0 || bidVast.Ads[0].InLine == nil {
			return RenderVastNoInline
		}
		impressions := bidVast.Ads[0].InLine.Impressions
		// 放入 imp tk
		impressions = append(impressions, vast.Impression{URI: impTk})
		bidVast.Ads[0].InLine.Impressions = impressions

		// 放入 click ck
		for i, ad := range bidVast.Ads {
			onlyCompanion := true // 当只有companion没有linear时， 将adx自己的点击上报放到companion里
			companionIndex := -1
			for j, creative := range ad.InLine.Creatives {
				if creative.Linear != nil {
					onlyCompanion = false
					if creative.Linear.VideoClicks != nil {
						clkTracking := creative.Linear.VideoClicks.ClickTrackings
						clkTracking = append(clkTracking, vast.VideoClick{URI: clkTk})
						bidVast.Ads[i].InLine.Creatives[j].Linear.VideoClicks.ClickTrackings = clkTracking
					}

					// 放入 ck event
					trackingEvents := creative.Linear.TrackingEvents
					for _, e := range trackEvents {
						trackingEvents = append(trackingEvents, vast.Tracking{
							Event: e.Event,
							URI:   e.URI,
						})
					}
					bidVast.Ads[i].InLine.Creatives[j].Linear.TrackingEvents = trackingEvents
				}
				if creative.CompanionAds != nil {
					companionIndex = j
				}
			}
			if onlyCompanion && companionIndex >= 0 {
				for j, cp := range ad.InLine.Creatives[companionIndex].CompanionAds.Companions {
					clkTracking := cp.CompanionClickTracking
					clkTracking = append(clkTracking, vast.CDATAString{CDATA: clkTk})
					bidVast.Ads[i].InLine.Creatives[companionIndex].CompanionAds.Companions[j].CompanionClickTracking = clkTracking
				}
			}
		}

		// 回填
		admVast, err := xml.Marshal(bidVast)
		if err != nil {
			return RenderDecodeVast
		}
		header := strings.Replace(xml.Header, "\n", "", 1)
		admStr := header + string(admVast) // 增加 xml header
		target.Adm = proto.String(admStr)
	} else {
		// pioneer 返回的广告格式都是 vast
		return RenderAdNotVast
	}

	reqCtx.RespData = backendCtx.RespData
	err = fillAd(bid, reqCtx, backendCtx)
	if err != nil {
		mvutil.Logger.Runtime.Errorf("[bid-server] load fillAd error:[%s]", err.Error())
		return RenderFillAd
	}

	return nil
}

//
func getLoadTimeout() int {

	//
	adnConf, _ := extractor.GetADNET_SWITCHS()

	//
	if timeout, ok := adnConf["bidServerLoadTimeout"]; ok && timeout > 0 {
		return timeout
	}

	// 2s 兜底
	return 2000
}
