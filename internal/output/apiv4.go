package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type V4Result struct {
	Status    int           `json:"status"`
	Msg       string        `json:"msg"`
	Data      V4Data        `json:"data"`
	DebugInfo []interface{} `json:"debuginfo,omitempty"`
}

type V4Data struct {
	SessionID         string            `json:"session_id"`
	ParentSessionID   string            `json:"parent_session_id"`
	AdType            int               `json:"ad_type"`
	AdSourceID        int               `json:"ad_source_id"`
	UnitSize          string            `json:"unit_size"`
	Frames            []V4Ads           `json:"frames"`
	OnlyImpressionURL string            `json:"only_impression_url,omitempty"`
	IAIcon            *string           `json:"ia_icon,omitempty"`
	IARst             *int              `json:"ia_rst,omitempty"`
	IAUrl             *string           `json:"ia_url,omitempty"`
	IAOri             *int              `json:"ia_ori,omitempty"`
	RKS               map[string]string `json:"rks,omitempty"`
	Setting           *Setting          `json:"setting,omitempty"`
}

type V4Ads struct {
	Ads      []Ad `json:"ads"`
	Template int  `json:"template,omitempty"`
}

func RenderV4Res(mr MobvistaResult, r *mvutil.RequestParams) V4Result {
	var result V4Result
	var ads V4Ads
	result.Status = 1
	result.Msg = "success"
	result.Data.SessionID = mr.Data.SessionID
	result.Data.ParentSessionID = mr.Data.ParentSessionID
	result.Data.AdType = mr.Data.AdType
	// 根据返回广告数量计算真实帧数
	frameNum := len(mr.Data.Ads) / r.Param.RequireNum
	// 计算最终返回的单子
	resNum := frameNum * r.Param.RequireNum
	// start render adsourceId
	for k, v := range mr.Data.Ads {
		result.Data.AdSourceID = v.AdSourceID
		ads.Ads = append(ads.Ads, v)
		// 广告填充满一帧则封装
		if (k+1)%r.Param.RequireNum == 0 {
			ads.Template = r.Param.Template
			result.Data.Frames = append(result.Data.Frames, ads)
			ads.Ads = []Ad{}
		}
		// 过滤多余的广告
		if k+1 >= resNum {
			break
		}
	}

	result.Data.UnitSize = mr.Data.UnitSize
	result.Data.OnlyImpressionURL = mr.Data.OnlyImpressionURL
	result.Data.IAIcon = mr.Data.IAIcon
	result.Data.IAOri = mr.Data.IAOri
	result.Data.IARst = mr.Data.IARst
	result.Data.IAUrl = mr.Data.IAUrl
	result.Data.RKS = mr.Data.RKS
	result.Data.Setting = mr.Data.Setting
	return result
}
