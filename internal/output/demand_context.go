package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
)

func NewDemandContext(params *mvutil.Params, r *mvutil.RequestParams) *demand.Context {
	uaInfo := mvutil.UaParser.Parse(params.UserAgent)
	uaOs := mvutil.UrlEncode(uaInfo.Os.Family + " " + uaInfo.Os.ToVersionString())
	var skVersion string
	var skAdNetworkIds []string
	if params.Skadnetwork != nil {
		skVersion = params.Skadnetwork.Ver
		skAdNetworkIds = params.Skadnetwork.Adnetids
	}
	return &demand.Context{
		// ClickId:         params.RequestID,
		PublisherId:     params.PublisherID,
		PublisherType:   params.PublisherType,
		AppId:           params.AppID,
		PublisherInfo:   r.PublisherInfo,
		AppInfo:         r.AppInfo,
		UnitInfo:        r.UnitInfo,
		PackageName:     params.ExtpackageName,
		AdType:          mvutil.GetAdTypeStr(params.AdType),
		AdvertiserId:    int64(params.AdvertiserID),
		CampaignId:      params.CampaignID,
		CreativeId:      params.CreativeId,
		ClientIP:        params.ClientIP,
		ImageSize:       params.ImageSize,
		PlatformName:    params.PlatformName,
		OSVersion:       params.OSVersion,
		SDKVersion:      params.SDKVersion,
		ScreenSize:      params.ScreenSize,
		Orientation:     params.Orientation,
		NetworkTypeName: params.NetworkTypeName,
		CountryCode:     params.CountryCode,
		CityCode:        params.CityCode,
		CityString:      params.CityString,
		MNC:             params.MNC,
		MCC:             params.MCC,
		Language:        params.Language,
		UserAgent:       params.UserAgent,
		Model:           params.Model,
		Brand:           params.Brand,
		IMEI:            params.IMEI,
		Md5IMEI:         params.ImeiMd5,
		//Sha1IMEI:            params.Imei,
		MAC:          params.MAC,
		AndroidId:    params.AndroidID,
		Md5AndroidId: params.AndroidIDMd5,
		GAID:         params.GAID,
		Md5GAID:      params.GAIDMd5,
		//Sha1GAID:            params.GAID,
		IDFA:    params.IDFA,
		Md5IDFA: params.IDFAMd5,
		OAID:    params.OAID,
		//ChannelPriceIn:  0, // 分渠道价格
		//ChannelPriceOut: 0, // 分渠道价格
		//ChannelBp:           params.Extchannel,
		UAOs:                uaOs,
		AdvCreativeId:       int64(params.AdvCreativeID),
		AdvCreativeName:     params.CreativeName,
		CsetName:            params.CsetName,
		AdvAdType:           params.CreativeAdType,
		MtgId:               params.ExtMtgId,
		ExtFinalSubId:       strconv.FormatInt(getSubId(params), 10),
		ExtFinalPackageName: params.ExtfinalPackageName,
		SkVersion:           skVersion,
		SkAdNetworkIds:      skAdNetworkIds,
		AppBundle:           params.RealPackageName,
		PingMode:            params.PingMode,
		IDFV:                params.IDFV,
		RequestType:         params.RequestType,
		UnitId:              params.UnitID,
		LinkType:            params.LinkType,
		SysId:               params.SysId,
		BkupId:              params.BkupId,
		Region:              params.ExtcdnType,
		MappingIDFA:         params.MappingIdfa,
		PassThroughData:     "",
		MiSkSpt:             params.MiSkSpt,
		MiSkSptDet:          mvutil.RenderIntslice(params.MiSkSptDet),
		SupplyPackageName:   params.PackageName,
	}
}

func ParseDemandContext4Params(ctx *demand.Context, params *mvutil.Params) *mvutil.Params {
	if ctx == nil {
		return params
	}

	// 上报了啥日志落啥
	params.ExtMtgId = ctx.MtgId
	params.Extfinalsubid, _ = strconv.ParseInt(ctx.ExtFinalSubId, 10, 64)
	params.ExtfinalPackageName = ctx.ExtFinalPackageName // 确认算法没有使用(也不应该使用，因为这个是BT之后的结果)
	params.PriceIn = ctx.ChannelPriceIn
	params.PriceOut = ctx.ChannelPriceOut
	params.LocalCurrency = ctx.LocalCurrency
	params.LocalChannelPriceIn = ctx.LocalChannelPriceIn
	params.SkadnetworkDataStr = ctx.SkAdNetworkDataStr
	params.SkAdNetworkId = ctx.SkAdNetwork.SkNetworkId
	params.SkCid = ctx.SkAdNetwork.SkCid
	params.SkTargetId = ctx.SkAdNetwork.SkTargetId
	params.SkNonce = ctx.SkAdNetwork.SkNonce
	params.SkSourceId = ctx.SkAdNetwork.SkSourceId
	params.SkTmp = ctx.SkAdNetwork.SkTimestamp
	params.SkSign = ctx.SkAdNetwork.SkSign
	params.SkNeed = ctx.SkAdNetwork.SkNeed
	params.SkViewSign = ctx.SkAdNetwork.SkViewSign
	//params.ExtChannelBp = ctx.ChannelBp
	params.RequestID = ctx.ClickId
	return params
}
