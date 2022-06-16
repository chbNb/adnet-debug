package output

import (
	"bytes"
	"html/template"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func RenderJssdkRes(r *mvutil.RequestParams) *Setting {
	setting := renderSetting(r)
	return &setting
}

func renderSetting(r *mvutil.RequestParams) Setting {
	var jsSetting Setting
	if mvutil.IsJsVideo(r.Param.AdType) {
		if &r.UnitInfo.AdSourceTime == nil {
			adSourceTime := map[string]int64{}
			jsSetting.AdSourceTime = &adSourceTime
		} else {
			jsSetting.AdSourceTime = &r.UnitInfo.AdSourceTime
		}
		jsSetting.Offset = mvutil.IntGetData(r.UnitInfo.Setting.Offset, 3)
		jsSetting.Autoplay = mvutil.IntGetData(r.UnitInfo.Setting.Autoplay, 2)
		jsSetting.Clicktype = mvutil.IntGetData(r.UnitInfo.Setting.ClickType, 1)
		jsSetting.DLNet = mvutil.IntGetData(r.UnitInfo.Setting.DLNet, 1)
		jsSetting.ShowImage = mvutil.IntGetData(r.UnitInfo.Setting.ShowImage, 1)
	} else if r.Param.AdType == mvconst.ADTypeWXNative || r.Param.AdType == mvconst.ADTypeWXBanner {
		jsSetting.Offset = mvutil.IntGetData(r.UnitInfo.Setting.Offset, 5)
		jsSetting.ApiCacheNum = mvutil.IntGetData(r.UnitInfo.Setting.ApiCacheNum, 3)
		jsSetting.RefreshFq = mvutil.IntGetData(r.UnitInfo.Setting.RefreshFq, 30)
	} else if r.Param.AdType == mvconst.ADTypeWXAppwall {
		jsSetting.IconImg = &r.UnitInfo.Unit.EntraImage
		jsSetting.IconTitle = &r.UnitInfo.Unit.EntraTitle
	} else if r.Param.AdType == mvconst.ADTypeWXRewardImg {
		jsSetting.DailyPlayCap = mvutil.IntGetData(r.UnitInfo.Setting.DailyPlayCap, 0)
		jsSetting.GifCtd = mvutil.IntGetData(r.UnitInfo.Setting.GifCtd, 5)
	} else {
		jsSetting.Hang = mvutil.IntGetData(r.UnitInfo.Unit.Hang, 0)
		jsSetting.IsIncent = mvutil.IntGetData(r.UnitInfo.Unit.IsIncent, 0)
		jsSetting.IsServerCall = mvutil.IntGetData(r.UnitInfo.Unit.IsServerCall, 0)
		jsSetting.VideoSkipTime = mvutil.IntGetData(r.UnitInfo.Setting.VideoSkipTime, 0)
		jsSetting.DailyPlayCap = mvutil.IntGetData(r.UnitInfo.Setting.DailyPlayCap, 0)
		jsSetting.Orientation = mvutil.IntGetData(r.Param.FormatOrientation, 0)
		jsSetting.CloseButtonDelay = mvutil.IntGetData(r.UnitInfo.Setting.CloseButtonDelay, 0)
		if &r.AppInfo.Rewards == nil {
			jsSetting.Reward = &[]smodel.Reward{}
		} else {
			jsSetting.Reward = &r.AppInfo.Rewards
		}
		if &r.UnitInfo.Unit.RecallNet == nil {
			recallNet := ""
			jsSetting.RecallNet = &recallNet
		} else {
			jsSetting.RecallNet = &r.UnitInfo.Unit.RecallNet
		}
		if &r.UnitInfo.Unit.EndcardTemplate == nil {
			endcardTemplate := ""
			jsSetting.EndcardTemplate = &endcardTemplate
		} else {
			jsSetting.EndcardTemplate = &r.UnitInfo.Unit.EndcardTemplate
		}
	}
	if !mvutil.IsWxAdType(r.Param.AdType) {
		jsSetting.CookieAchieve = mvutil.IntGetData(r.UnitInfo.Unit.CookieAchieve, 0)
		jsSetting.OffsetMax = mvutil.IntGetData(r.UnitInfo.Setting.Offset, 3)
		jsSetting.VideoInteractiveType = mvutil.IntGetData(r.UnitInfo.Setting.VideoInteractiveType, 3)
		jsSetting.MuteMode = mvutil.IntGetData(r.UnitInfo.Setting.MuteMode, 1)
		jsSetting.IsReady = mvutil.IntGetData(r.UnitInfo.Setting.ReadyRate, 100)
	}
	jsSetting.Plct = mvutil.IntGetData(r.AppInfo.App.Plct, 900)

	return jsSetting
}

func RenderJssdkResInHtml(res string, r mvutil.Params, data *MobvistaResult) (string, error) {
	t, err := template.ParseFiles("./conf/vijs_template_banner.html")
	if err != nil {
		return "", err
	}
	htmlData := HtmlData{
		Data:   data,
		AppId:  strconv.FormatInt(r.AppID, 10),
		UnitId: strconv.FormatInt(r.UnitID, 10),
		AdType: strconv.FormatInt(int64(r.AdType), 10),
		Sign:   r.Sign,
		Idfa:   r.IDFA,
		Imei:   r.IMEI,
		Gaid:   r.GAID,
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, htmlData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

type HtmlData struct {
	Data   *MobvistaResult
	AppId  string
	UnitId string
	AdType string
	Sign   string
	Idfa   string
	Imei   string
	Gaid   string
}
