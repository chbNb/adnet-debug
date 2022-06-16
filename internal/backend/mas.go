package backend

import (
	"errors"
	"strconv"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtg_hb_rtb"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

func ConstructMasInfo(reqCtx *mvutil.ReqCtx) (asInfo *mtgrtb.BidRequest_AsInfo, err error) {
	// 对照 composeAdServerRequest()
	if reqCtx == nil || reqCtx.ReqParams == nil {
		return nil, errors.New("constructAsInfo params is invalidate")
	}
	reqParam := reqCtx.ReqParams.Param
	asInfo = new(mtgrtb.BidRequest_AsInfo)
	reqId := mvutil.GetRequestID()
	asInfo.SessionId = &reqId
	adTypeStr := mvutil.GetAdTypeStr(reqParam.AdType)
	if reqParam.RequestPath == mvconst.PATHREQUEST && len(adTypeStr) <= 0 {
		adTypeStr = reqParam.AdTypeStr
	}
	asInfo.AdTypeStr = &adTypeStr
	ImageSizeId := mtgrtb.BidRequest_AsInfo_ImageSizeEnum(ad_server.ImageSizeEnum(reqParam.ImageSizeID))
	asInfo.ImageSizeId = &ImageSizeId
	asInfo.UnitSize = &reqParam.UnitSize
	asInfo.PriceFloor = &reqParam.PriceFloor
	asInfo.Scenario = &reqParam.Scenario
	asInfo.AdNum = &reqParam.AdNum
	tnum := int64(reqCtx.ReqParams.Param.TNum)
	// 达到展示控cap的条件，trueNum传0给AdServer
	if reqCtx.ReqParams.Param.IsBlockByImpCap {
		tnum = int64(0)
	}
	asInfo.TrueNum = &tnum
	asInfo.UnSupportSdkTruenum = &reqParam.UnSupportSdkTrueNum
	asInfo.Offset = &reqParam.Offset
	asInfo.NativeInfoList = make([]*mtgrtb.BidRequest_AsInfo_NativeInfo, len(reqParam.NativeInfoList))
	if len(reqParam.NativeInfoList) > 0 {
		for i := 0; i < len(reqParam.NativeInfoList); i++ {
			info := reqParam.NativeInfoList[i]
			adTemplate := mtgrtb.BidRequest_AsInfo_ADTemplate(info.AdTemplate)
			requiredNum := int32(reqParam.AdNum)
			nativeInfo := mtgrtb.BidRequest_AsInfo_NativeInfo{
				AdTemplate: &adTemplate,
				RequireNum: &requiredNum,
			}
			asInfo.NativeInfoList[i] = &nativeInfo
		}
	}
	asInfo.ShowedCampaignIdList = reqParam.DisplayCamIds
	ResourceType := mtgrtb.BidRequest_AsInfo_InteractiveResourceType(reqParam.IARst)
	asInfo.ResourceType = &ResourceType
	asInfo.AppVersionCode = &reqParam.AppVersionCode
	asInfo.SdkVersion = &reqParam.SDKVersion
	asInfo.OsVersionCodeV2 = &reqParam.OSVersionCode
	asInfo.ExcludeIdSet = renderExcludeIDslice(reqParam.ExcludeIDS)
	Timestamp := time.Now().UTC().Unix()
	asInfo.Timestamp = &Timestamp
	// no appid, unitid
	RequestType := mtgrtb.BidRequest_AsInfo_RequestType(ad_server.RequestType(reqParam.RequestType))
	asInfo.RequestType = &RequestType
	adSourceList := renderAdSource(reqCtx.ReqParams)
	adSources := make([]mtgrtb.BidRequest_AsInfo_ADSource, len(adSourceList))
	for i := 0; i < len(adSourceList); i++ {
		adSources[i] = mtgrtb.BidRequest_AsInfo_ADSource(adSourceList[i])
	}
	asInfo.AdSourceList = adSources
	asInfo.ApiVersion = &reqParam.ApiVersionCode
	asInfo.DebugMode = &reqParam.DebugMode
	asInfo.Imei = &reqParam.IMEI
	asInfo.Mac = &reqParam.MAC
	asInfo.DevId = &reqParam.AndroidID
	asInfo.Oaid = &reqParam.OAID
	asInfo.ScreenSize = &reqParam.ScreenSize
	orientation := 0
	if mvutil.IsHBS2SUseVideoOrientation(reqCtx.ReqParams) && reqParam.Orientation == mvconst.ORIENTATION_BOTH {
		orientation = reqParam.FormatOrientation
	} else {
		orientation = reqParam.Orientation
	}
	Orientation := mtgrtb.BidRequest_AsInfo_Orientation(orientation)
	asInfo.Orientation = &Orientation
	// no mcc mnc networktype, language, ip, adnserverIp, countryCode
	// no pSessionId
	asInfo.Timezone = &reqParam.TimeZone
	Category := mtgrtb.BidRequest_AsInfo_Category(reqCtx.ReqParams.Param.Category)
	asInfo.Category = &Category
	asInfo.RealAppId = &reqParam.RealAppID
	asInfo.DeviceModel = &reqParam.ReplaceModel // UnifiedDeviceModel
	asInfo.DeviceBrand = &reqParam.ReplaceBrand // UnifiedDeviceBrand
	NetworkType := mtgrtb.BidRequest_AsInfo_NetworkType(reqParam.NetworkType)
	asInfo.NetworkType = &NetworkType
	IfSupportSeperateCreative := int32(2)
	// 是否走素材三期逻辑
	if reqCtx.ReqParams.Param.NewCreativeFlag {
		IfSupportSeperateCreative = int32(3)
	}
	asInfo.IfSupportSeperateCreative = &IfSupportSeperateCreative
	asInfo.NetworkId = &reqParam.Network
	asInfo.InstallIdSet = renderExcludeIDslice(reqCtx.ReqParams.Param.InstallIDS)
	RecallADNOffer := int32(1)
	asInfo.RecallADNOffer = &RecallADNOffer
	asInfo.RankerInfo = &reqParam.RankerInfo
	asInfo.ExtData2 = &reqParam.ExtData2MAS // 发送给pion
	asInfo.VideoVersion = &reqParam.VideoVersion
	asInfo.LowDevice = &reqParam.LowDevice
	asInfo.TargetPackageNameSet = renderPackageListSlice(reqParam.PlPkg)
	if reqParam.Mof == 1 {
		if reqCtx.ReqParams.Param.UcParentUnitId > 0 {
			asInfo.ParentUnitId = &reqCtx.ReqParams.Param.UcParentUnitId
		} else {
			asInfo.ParentUnitId = &reqCtx.ReqParams.Param.ParentUnitId
		}
		asInfo.MofData = &reqParam.MofData
	}
	asInfo.GPVersion = &reqParam.GPVersion
	excludePkgNames := make([]string, 0)
	for pkg := range reqParam.ExcludePackageNames {
		excludePkgNames = append(excludePkgNames, pkg)
	}
	asInfo.DEPRECATEDExcludePackageNameSet = excludePkgNames
	flowTagIdStr := strconv.Itoa(reqCtx.FlowTagID)
	asInfo.FlowTagId = &flowTagIdStr
	asInfo.Channel = &reqParam.Extchannel
	CDNAbTastStr := strconv.Itoa(reqParam.ExtCDNAbTest)
	asInfo.CdnAbtest = &CDNAbTastStr
	asInfo.ReqType = &reqParam.ReqType
	asInfo.ServerIp = &reqParam.ServerIP
	asInfo.RemoteIp = &reqParam.RemoteIP
	asInfo.CloseAdTag = &reqParam.ExtDataInit.CloseAdTag
	asInfo.CtnSizeTag = &reqParam.CtnSizeTag
	pingMode := int32(reqParam.PingMode)
	asInfo.PingMode = &pingMode
	asInfo.HttpReq = &reqParam.HTTPReq
	asInfo.RequestId = &reqParam.RequestID
	asInfo.SystemId = &reqParam.SysId
	asInfo.SysbkupId = &reqParam.BkupId
	asInfo.Ruid = &reqParam.RuId
	asInfo.FillEcpmFloor = &reqParam.FillEcpmFloor
	cityCode := int32(reqParam.CityCode)
	asInfo.CityCode = &cityCode
	if reqParam.PubFlowExpectPrice > 0.0 {
		asInfo.PubFlowExpectPrice = &reqParam.PubFlowExpectPrice
	}
	asInfo.MvLine = &reqParam.MvLine
	extra := mvutil.GetCloudExtra(extractor.GetCLOUD_NAME()) + "_pioneer"
	// 对于关闭场景广告，加上后缀方便在报表区分数据
	if mvutil.IsAppwallOrMoreOffer(reqParam.AdType) && (reqParam.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD ||
		reqParam.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD_MORE_OFFER) {
		extra += "_clsad"
	}
	if reqCtx.ReqParams.IsHBRequest {
		extra = reqCtx.ReqParams.Param.Algorithm + "-pioneer"
		reqCtx.ReqParams.Param.Algorithm = extra
		reqCtx.ReqParams.Param.Extra = extra
	}
	if reqCtx.ReqParams.IsTopon {
		extra = "topon-nothb"
		reqCtx.ReqParams.Param.Extra = extra
		reqCtx.ReqParams.Param.Algorithm = extra
	}
	asInfo.Extra = &extra
	asInfo.ExcludePkg = &reqParam.ExtDataInit.ExcludePkg
	asInfo.ExcludePsbPkg = &reqParam.ExtDataInit.ExcludePsbPkg
	asInfo.Region = &reqParam.ExtcdnType
	if reqParam.BigTemplateFlag {
		IfSupportBigTemplate := int32(1)
		asInfo.IfSupportBigTemplate = &IfSupportBigTemplate
	}

	if reqCtx.ReqParams.Param.Skadnetwork != nil {
		if demand.IsSupMtgSKAdnetwork(reqCtx.ReqParams.AppInfo.RealPackageName, reqCtx.ReqParams.Param.Skadnetwork.Adnetids, reqCtx.ReqParams.Param.Skadnetwork.Ver) {
			supportSkadn := true
			asInfo.IfSupportSkadn = &supportSkadn
		}
	}

	asInfo.ImpExcludePkg = &reqParam.ExtDataInit.ImpExcludePkg
	if reqParam.PolarisFlag {
		supportPolaris := int32(1)
		asInfo.IfSupportPolaris = &supportPolaris
	}

	asInfo.PlacementId = &reqCtx.ReqParams.Param.FinalPlacementId
	asInfo.PassthroughData = RenderPassthroughData(&reqCtx.ReqParams.Param)
	asInfo.IfSupportDco = &reqCtx.ReqParams.Param.IfSupDco
	asInfo.TestMode = &reqCtx.ReqParams.Param.AsTestMode
	asInfo.RandNum = &reqCtx.ReqParams.Param.RandNum
	asInfo.IdfaMd5 = &reqCtx.ReqParams.Param.IDFAMd5
	asInfo.GaidMd5 = &reqCtx.ReqParams.Param.GAIDMd5
	asInfo.WebEnv = RenderWebEnv(&reqCtx.ReqParams.Param)
	asInfo.Idfv = &reqCtx.ReqParams.Param.IDFV
	asInfo.Openidfa = &reqCtx.ReqParams.Param.OpenIDFA
	asInfo.RwPlus = &reqCtx.ReqParams.Param.RwPlus
	asInfo.IfSupH265 = &reqCtx.ReqParams.Param.IfSupH265
	asInfo.NeedCreativeDataCIds = renderExcludeIDslice(reqCtx.ReqParams.Param.NeedCreativeDataCIds)
	asInfo.DebugCampaignIdList = renderExcludeIDslice(reqCtx.ReqParams.Param.TargetIds)
	asInfo.MappingIdfa = &reqCtx.ReqParams.Param.MappingIdfa
	asInfo.MiSkSpt = &reqCtx.ReqParams.Param.MiSkSpt
	asInfo.MiSkSptDet = mvutil.RenderIntslice(reqCtx.ReqParams.Param.MiSkSptDet)
	asInfo.SupplyPackageName = &reqCtx.ReqParams.Param.PackageName
	asInfo.DspJunoRes = &reqCtx.ReqParams.Param.DspMoreOfferInfo
	mofTypeInt32 := int32(reqCtx.ReqParams.Param.MofType)
	asInfo.MoreOfferType = &mofTypeInt32 // pioneer暂无使用的地方
	asInfo.ParentId = &reqCtx.ReqParams.Param.ParentId
	asInfo.MoreofferUnitId = &reqCtx.ReqParams.UnitInfo.Unit.MofUnitId
	asInfo.SystemUseragent = &reqCtx.ReqParams.Param.ExtsystemUseragent
	asInfo.MoreOfferRequestId = &reqCtx.ReqParams.Param.MoreOfferRequestId
	var moreofferTrafficType int32
	if reqCtx.ReqParams.Param.DspMof == 1 {
		moreofferTrafficType = 1
	}
	asInfo.MoreofferTrafficType = &moreofferTrafficType
	asInfo.CachedCampaignIds = renderExcludeIDslice(reqCtx.ReqParams.Param.CachedCampaignIds)
	asInfo.FixedEcpm = &reqCtx.ReqParams.Param.FixedEcpm
	asInfo.Dnt = &reqCtx.ReqParams.Param.Dnt
	isMp := mvutil.IsMpad(reqCtx.ReqParams.Param.RequestPath)
	asInfo.IsMp = &isMp
	// 透传字段
	// asInfo.AdnLibPassthroughData
	return
}

func fillMasAd(rdata *mtgrtb.BidResponse, reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, ad *corsair_proto.Campaign, k int) {

	asResp := rdata.GetAsResp()
	if asResp == nil {
		return
	}
	reqCtx.ReqParams.Param.RespFillEcpmFloor = asResp.GetEcpmFloor()
	reqCtx.ReqParams.AsResp = asResp

	// more_offer的requestid
	moreOfferRequestId := asResp.GetMoreOfferRequestid()
	if len(moreOfferRequestId) > 0 {
		reqCtx.ReqParams.Param.RequestID = moreOfferRequestId
	}

	sdkParams := asResp.GetSdkParam()
	// 防止越界
	if k+1 > len(sdkParams) {
		return
	}
	if len(sdkParams) > 0 {
		sdkParam := sdkParams[k]
		ad.AdvertiserId = sdkParam.AdvId
		ad.CType = sdkParam.Ctype
		ad.AppSize = sdkParam.AppSize
		ad.ClickMode = sdkParam.ClickMode
		ad.WatchMile = sdkParam.WatchMile
		ad.FCA = sdkParam.Fca
		if sdkParam.ImageSizeId != nil {
			imageSizeId := *sdkParam.ImageSizeId
			ad.ImageSizeId = ad_server.ImageSizeEnum(imageSizeId)
		}
		if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && //都要加流量判断
			sdkParam.ImageSize != nil && len(*sdkParam.ImageSize) > 0 {
			ad.ImageSize = sdkParam.ImageSize
		}
		masAdvImp := sdkParam.GetAdvImp()
		advImps := make([]*corsair_proto.AdvImp, len(masAdvImp))
		for i := 0; i < len(masAdvImp); i++ {
			advImp := &corsair_proto.AdvImp{
				URL:    masAdvImp[i].GetUrl(),
				Second: masAdvImp[i].GetSec(),
			}
			advImps[i] = advImp
		}
		ad.AdvImpList = advImps
		ad.AdURLList = sdkParam.GetAdUrlList()
		// ad_render: sub_category_name, md5_file, storekit, storekit_time, endcard_click_result, imp_ua, c_ua, deep_link
		ad.RetargetOffer = sdkParam.RetargetOffer
		ad.LandingType = sdkParam.LandingType
		ad.LinkType = sdkParam.LinkType
		for _, lb := range sdkParam.GetLoopback() {
			ad.LoopBack[lb.GetKey()] = lb.GetValue()
		}
		ad.OfferType = sdkParam.OfferType
		if len(sdkParam.GetPackageName()) > 0 {
			ad.PackageName = sdkParam.PackageName
		}
		ad.AdTemplate = (*ad_server.ADTemplate)(sdkParam.AdTemplate)
		ad.WTick = sdkParam.Wtick
		ad.Floatball = sdkParam.Flb
		ad.FloatballSkipTime = sdkParam.FlbSkiptime
		if reqCtx.ReqParams.IsBidRequest && sdkParam.GetSkAdnetwork() != nil && len(sdkParam.GetSkAdnetwork().GetVersion()) > 0 {
			version := sdkParam.GetSkAdnetwork().GetVersion()
			network := sdkParam.GetSkAdnetwork().GetNetwork()
			campaign := sdkParam.GetSkAdnetwork().GetCampaign()
			itunesItem := sdkParam.GetSkAdnetwork().GetItunesitem()
			nonce := sdkParam.GetSkAdnetwork().GetNonce()
			sourceApp := sdkParam.GetSkAdnetwork().GetSourceapp()
			timestamp := sdkParam.GetSkAdnetwork().GetTimestamp()
			signature := sdkParam.GetSkAdnetwork().GetSignature()

			reqCtx.ReqParams.BidSkAdNetwork = &mtg_hb_rtb.BidResponse_SeatBid_Bid_Ext_Skadn{
				Version:    &version,
				Network:    &network,
				Campaign:   &campaign,
				Itunesitem: &itunesItem,
				Nonce:      &nonce,
				Sourceapp:  &sourceApp,
				Timestamp:  &timestamp,
				Signature:  &signature,
			}
		}
	}
}
