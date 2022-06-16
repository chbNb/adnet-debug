package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/voyager/clickmode/clkmode_context"
)

func NewClickmodeContext(params *mvutil.Params, r *mvutil.RequestParams) *clkmode_context.ClickModeContext {
	return &clkmode_context.ClickModeContext{
		PublisherId:   params.PublisherID,
		AppId:         params.AppID,
		PublisherInfo: r.PublisherInfo,
		AppInfo:       r.AppInfo,
		UnitInfo:      r.UnitInfo,
		AdType:        mvutil.GetAdTypeStr(params.AdType),
		AdvertiserId:  int(params.AdvertiserID),
		CampaignId:    params.CampaignID,
		ClientIP:      params.ClientIP,
		PlatForm:      params.PlatformName,
		SDKVersion:    params.SDKVersion,
		CountryCode:   params.CountryCode,
		Model:         params.Model,
		IMEI:          params.IMEI,
		AndroidId:     params.AndroidID,
		GAID:          params.GAID,
		IDFA:          params.IDFA,
		OAID:          params.OAID,
		PingMode:      params.PingMode,
		IDFV:          params.IDFV,
		RequestType:   params.RequestType,
		UnitId:        params.UnitID,
		LinkType:      params.LinkType,
		SysId:         params.SysId,
		BkupId:        params.BkupId,
		Region:        params.ExtcdnType,
	}
}
