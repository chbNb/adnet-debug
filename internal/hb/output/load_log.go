package output

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func FormatLoadLog(in *mvutil.ReqCtx, filterCode int, filterMsg string) string {
	var dspID string
	dspExt, err := in.ReqParams.GetDspExt()
	if err == nil {
		dspID = strconv.FormatInt(dspExt.DspId, 10)
	}
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatInt(in.ReqParams.Param.PublisherID, 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatInt(in.ReqParams.Param.AppID, 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatInt(in.ReqParams.Param.UnitID, 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extra)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(constant.OpenApi)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(mvutil.GetAdTypeStr(in.ReqParams.Param.AdType))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(helpers.GetOs(in.ReqParams.Param.Platform))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strings.ToLower(in.ReqParams.Param.SDKVersion))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.AppVersionName)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.CountryCode)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatInt(in.ReqParams.Param.CityCode, 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(constant.APIOffer))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extchannel)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(dspID)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Token)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.RequestID)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.PriceBigDecimal)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(in.ReqParams.Nbr))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(filterCode))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extra2)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("0")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ClientIP)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(in.ReqParams.Param.Extrvtemplate))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extendcard)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("0")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ExtBigTemId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ExtBigTplOfferData)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ExtPlacementId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.DspExt)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatFloat(in.ReqParams.Param.Dmt, 'f', 2, 64))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatFloat(in.ReqParams.Param.Dmf, 'f', 2, 64))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Ct)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(int(in.ReqParams.Param.HBBidTestMode)))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(in.ReqParams.Param.PowerRate))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(in.ReqParams.Param.Charging))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.TotalMemory)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ResidualMemory)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.BidFloorBigDecimal)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.MediationName)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ChannelInfo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extra3)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.Extalgo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ExtAdxAlgo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.HBExtPfData)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(mvutil.RawUrlEncode(in.ReqParams.Param.UserAgent))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.FormatInt(in.ReqParams.Param.BidUnitID, 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(in.ReqParams.Param.BidAdType))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(filterMsg)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.OSVersion)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.IDFA)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.IDFV)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.IMEI)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.AndroidID)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.GAID)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.AppSettingId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.UnitSettingId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.RewardSettingId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.ExtData)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.PioneerExtdataInfo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(in.ReqParams.Param.PioneerOfferExtdataInfo)
	return buf.String()
}
