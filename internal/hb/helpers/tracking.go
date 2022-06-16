package helpers

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func GenBackendContent(backendId, contentType int, requestKey string) string {
	dataBuf := bytes.NewBufferString("")
	dataBuf.WriteString(strconv.Itoa(backendId))
	dataBuf.WriteString(constant.SplitterComma)
	dataBuf.WriteString(requestKey)
	dataBuf.WriteString(constant.SplitterComma)
	dataBuf.WriteString(strconv.Itoa(contentType))
	return dataBuf.String()
}

func ConfigCenterKey(region string) string {
	switch region {
	case constant.SG:
		return constant.CommonSG
	case constant.VG:
		return constant.CommonVG
	case constant.FK:
		return constant.CommonFK
	case constant.OH:
		return constant.CommonOH
	default:
		return constant.CommonSG
	}
}

func GenSchema(httpReq int) string {
	if httpReq == 2 {
		return "https"
	}
	return "http"
}

func GenFullUrl(httpReq int, query string) string {
	if len(query) == 0 {
		return ""
	}
	if strings.HasPrefix(query, "http") {
		return query
	}
	return GenSchema(httpReq) + "://" + query
}

func GenPlayableUrl(playableUrl string, protocal int, httpReq int) string {
	switch protocal {
	case 1:
		playableUrl = "http://" + playableUrl
	case 2:
		playableUrl = "https://" + playableUrl
	default:
		playableUrl = GenSchema(httpReq) + "://" + playableUrl
	}
	return playableUrl
}

func GenBackendData(backendId int, campaignId int64) string {
	buf := bytes.NewBufferString("")
	buf.WriteString(strconv.Itoa(backendId))
	buf.WriteString(constant.SplitterComma)
	buf.WriteString(strconv.FormatInt(campaignId, 10))
	buf.WriteString(constant.SplitterComma)
	buf.WriteString("1")
	return buf.String()
}

func GenParamOZ(loadReq *params.LoadReqData, ad *params.Ad) string {
	var ExtBigTemplateId string
	if len(ad.BigTemlate) > 0 {
		ExtBigTemplateId = ad.BigTemlate
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{
		loadReq.PublisherIdStr,
		loadReq.AppIdStr,
		loadReq.UnitIdStr,
		strconv.FormatInt(loadReq.AdvertiserId, 10),
		strconv.FormatInt(loadReq.CampaignId, 10),
		"",
		loadReq.Scenario,
		"",
		loadReq.ImageSize,
		strconv.Itoa(loadReq.RequestType),
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.CountryCode,
		"",
		"",
		loadReq.Mcc + loadReq.Mnc,
		loadReq.Extra,
		"",
		loadReq.Extra3,
		loadReq.Extra4,
		loadReq.Extra5,
		"",
		strconv.Itoa(loadReq.Extra7),
		strconv.Itoa(loadReq.Extra8),
		loadReq.Extra9,
		loadReq.Extra10,
		loadReq.BidId,
		"",
		"",
		"",
		"",
		loadReq.ServerIp,
		"",
		"",
		"",
		"",
		"",
		loadReq.AppVersionName,
		"",
		loadReq.RemoteIp,
		"",
		"",
		"",
		loadReq.CityCode,
		strconv.FormatInt(int64(loadReq.Extra13), 10),
		strconv.FormatInt(loadReq.Extra14, 10),
		loadReq.Extra15,
		strconv.Itoa(loadReq.Extra16),
		"",
		"",
		"",
		loadReq.Extra20,
		"",
		strconv.FormatInt(loadReq.Extfinalsubid, 10),
		"",
		"",
		"",
		loadReq.ExtpackageName,
		"",
		strconv.Itoa(loadReq.ExtflowTagId),
		"",
		strconv.FormatInt(loadReq.Extendcard, 10),
		strconv.Itoa(loadReq.ExtrushNoPre),
		"",
		loadReq.ExtfinalPackageName,
		strconv.Itoa(loadReq.Extnativex),
		"",
		"",
		"",
		strconv.FormatInt(int64(loadReq.Extctype), 10),
		strconv.Itoa(loadReq.Extrvtemplate),
		"",
		"",
		"",
		"",
		"",
		"",
		strconv.Itoa(loadReq.Extb2t),
		loadReq.Extchannel,
		"",
		"",
		"",
		loadReq.Extbp,
		strconv.FormatInt(int64(loadReq.Extsource), 10),
		"",
		loadReq.Extalgo,
		loadReq.ExtthirdCid,
		loadReq.ExtifLowerImp,
		"",
		"",
		loadReq.ExtsystemUseragent,
		"",
		"",
		"",
		"",
		"",
		loadReq.ExtMpNormalMap,
		"",
		"",
		loadReq.ReplaceBrand,
		loadReq.ReplaceModel,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.Oaid,
		"",
		"",
		"",
		strings.Join(loadReq.ExtBigTplOfferDataList, ";"),
		ExtBigTemplateId,
		"",
		"",
		strconv.FormatInt(loadReq.PlacementId, 10),
	}, "|")))
}

func GenParamP(loadReq *params.LoadReqData) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		AdTypeStr(loadReq.AdType),
		loadReq.ImageSize,
		"",
		loadReq.Os,
		loadReq.OsVersion,
		loadReq.SdkVersion,
		loadReq.Model,
		loadReq.ScreenSize,
		strconv.Itoa(loadReq.Orientation),
		"",
		loadReq.Language,
		loadReq.NetworkTypeName,
		loadReq.Mcc + loadReq.Mnc,
		"",
		"",
		loadReq.Extra3,
		loadReq.Extra4,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.ClientIp,
		loadReq.Imei,
		loadReq.Mac,
		loadReq.AndroidId,
		"",
		"",
		"",
		"",
		loadReq.Gaid,
		loadReq.Idfa,
		"",
		loadReq.Brand,
		loadReq.RemoteIP,
		loadReq.SessionID,
		loadReq.ParentSessionID,
		"",
		"",
		"",
		"",
		"1",
		"",
		loadReq.Idfv + "," + loadReq.OpenIdfa,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.Extstats,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.Extstats2,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.ReplaceBrand,
		loadReq.ReplaceModel,
		"",
		"",
		loadReq.ExtCreativeNew,
		"",
		"",
		"",
		"",
		"",
		"",
		loadReq.ExtSysId,
		strconv.FormatFloat(loadReq.ApiVersion, 'f', 1, 64),
	}, "|")))
}

func GenParamQ(loadReq *params.LoadReqData, campaign *smodel.CampaignInfo, ad *params.Ad) string {
	oriPrice := loadReq.PriceIn
	price := loadReq.PriceOut

	var extSlotId string
	if ad.OfferExtData != nil {
		extSlotId = strconv.Itoa(int(ad.OfferExtData.SlotId))
	}

	return "a_" + Base64Encode(strings.Join([]string{
		"2.0",
		loadReq.Domain,
		loadReq.BidId,
		loadReq.Extra5,
		loadReq.PublisherIdStr,
		loadReq.AppIdStr,
		loadReq.UnitIdStr,
		strconv.FormatInt(loadReq.AdvertiserId, 10),
		strconv.FormatInt(loadReq.CampaignId, 10),
		strconv.FormatFloat(oriPrice, 'f', 2, 64),
		strconv.FormatFloat(price, 'f', 2, 64),
		loadReq.Extra,
		strconv.Itoa(constant.OPENAPI_V3),
		loadReq.CountryCode,
		strconv.Itoa(loadReq.Extra8),
		loadReq.Extra10,
		strconv.FormatInt(int64(loadReq.Extra13), 10),
		strconv.FormatInt(loadReq.Extra14, 10),
		loadReq.CityCode,
		loadReq.AppVersionName,
		constant.OpenApi,
		strconv.Itoa(loadReq.Extra16), "|-",
		strconv.Itoa(loadReq.Extra7),
		strconv.Itoa(loadReq.Extbtclass),
		strconv.FormatInt(loadReq.Extfinalsubid, 10),
		strconv.Itoa(loadReq.ExtdeleteDevid),
		strconv.Itoa(loadReq.ExtinstallFrom),
		strconv.Itoa(loadReq.ExtnativeVideo),
		strconv.Itoa(loadReq.ExtflowTagId), "",
		strconv.FormatInt(loadReq.Extendcard, 10), "",
		strconv.FormatInt(loadReq.ExtdspRealAppid, 10),
		loadReq.ExtfinalPackageName,
		strconv.Itoa(loadReq.Extnativex),
		strconv.FormatInt(loadReq.CreativeId, 10),
		strconv.Itoa(loadReq.Extattr), "",
		strconv.FormatInt(int64(loadReq.Extctype), 10),
		strconv.Itoa(loadReq.Extrvtemplate), "",
		strconv.Itoa(loadReq.Extplayable),
		strconv.Itoa(loadReq.Extb2t),
		loadReq.Extchannel,
		strconv.Itoa(loadReq.Extabtest2),
		loadReq.Extbp,
		strconv.FormatInt(int64(loadReq.Extsource), 10),
		strconv.Itoa(loadReq.ExtappearUa),
		strconv.Itoa(loadReq.ExtCDNAbTest),
		loadReq.ExtMpNormalMap,
		strconv.Itoa(loadReq.Extplayable2),
		loadReq.Extabtest3,
		loadReq.ExtData2,
		strconv.FormatFloat(loadReq.BidPrice, 'f', 2, 64), "",
		strconv.FormatInt(loadReq.ExtMtgId, 10),
		extSlotId,
		loadReq.ExtCreativeNew,
	}, "|"))
}

func GenContent(videoUrl string) int {
	content := constant.Image
	if len(videoUrl) > 0 {
		content = constant.Video
	}
	return content
}

func GenParamZ(loadReq *params.LoadReqData) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{
		loadReq.ImageSize,
		loadReq.Extra3,
		loadReq.Extra4,
		strconv.Itoa(loadReq.AdNum),
		loadReq.ReplaceBrand,
		loadReq.ReplaceModel,
		"",
		loadReq.ExtApiVersion,
		loadReq.Mcc + loadReq.Mnc,
		loadReq.Oaid,
		loadReq.ExtBigTemId,
		strconv.FormatInt(loadReq.PlacementId, 10),
	}, "|")))
}

func GenParamCSP(loadReq *params.LoadReqData) string {
	return Base64Encode(strings.Join([]string{
		loadReq.AdBackend,
		loadReq.AdBackendData,
		strconv.Itoa(loadReq.ExtflowTagId),
		strconv.Itoa(loadReq.RandValue),
		loadReq.BackendConfig,
		strconv.Itoa(loadReq.AdNum),
		strconv.Itoa(loadReq.AdNum),
		loadReq.DspExt,
	}, "|"))
}

func GenMidwayParam(loadReq *params.LoadReqData, adBackend, adBackendData, backendConfig, dspExt, fakePrice, priceFactor string) string {
	rawMidwayParam := strings.Join([]string{
		// request_id
		loadReq.BidId,
		// publisher_id
		loadReq.PublisherIdStr,
		// app_id
		loadReq.AppIdStr,
		// unit_id
		loadReq.UnitIdStr,
		// ad_backend
		loadReq.AdBackend,
		// ad_backend_data
		loadReq.AdBackendData,
		// country_code
		loadReq.CountryCode,
		// city_code
		loadReq.CityCode,
		// platform
		strconv.Itoa(loadReq.Platform),
		// adType
		strconv.Itoa(loadReq.AdType),
		// os_version
		loadReq.OsVersion,
		// sdk_version
		loadReq.SdkVersion,
		// app_version
		loadReq.AppVersionName,
		// device_brand
		loadReq.Brand,
		// device_model
		loadReq.Model,
		// screen_size
		loadReq.ScreenSize,
		// orientation
		strconv.Itoa(loadReq.Orientation),
		// language
		loadReq.Language,
		// network_type
		strconv.Itoa(loadReq.NetworkType),
		// mcc_mnc
		loadReq.Mcc + loadReq.Mnc,
		// client_ip
		loadReq.ClientIp,
		// remote_ip
		"0",
		// server_ip
		loadReq.ServerIp,
		// mei
		loadReq.Imei,
		// mac
		"0",
		// android_id
		loadReq.AndroidId,
		// gaid
		loadReq.Gaid,
		// idfa
		loadReq.Idfa,
		// flow_tag_id
		strconv.Itoa(loadReq.ExtflowTagId),
		// request_num
		strconv.Itoa(loadReq.AdNum),
		// t_num
		strconv.Itoa(loadReq.AdNum),
		// rand_value
		strconv.Itoa(loadReq.RandValue),
		// backend_config
		loadReq.BackendConfig,
		// tracking_json
		loadReq.PlayInfo,
		// scenairo
		constant.OpenApi,
		// is_video(native)
		"0",
		// extra
		loadReq.Extra,
		// third_template
		loadReq.MidwayCreativeData,
		// dsp_ext
		loadReq.DspExt,
		// fake_dsp_price
		fakePrice,
		// req_backend
		"0",
		// reject_code
		"",
		// ext_is_imp_timeout
		"0",
		// price_factor
		loadReq.PriceFactor,
		// mgt_channel
		loadReq.Extchannel,
		// 对应ThirdPartyABTestStr参数，不过没有ThirdPartyABTest内容
		loadReq.ExtData2Log,
	}, "|")

	enMidwayParam := Base64Encode(rawMidwayParam)
	enMidwayParam = strings.Replace(enMidwayParam, "=", "%3D", -1)
	enMidwayParam = strings.Replace(enMidwayParam, "+", "%2B", -1)
	enMidwayParam = strings.Replace(enMidwayParam, "/", "%2F", -1)
	return enMidwayParam
}
