package output

import "gitlab.mobvista.com/ADN/adnet/internal/mvutil"

type V2Result struct {
	Stat              int      `json:"status"`
	Msg               string   `json:"msg"`
	OnlyImpressionUrl string   `json:"only_impression_url"`
	Data              []V2Data `json:"data"`
}

type V2Data struct {
	Id              int64    `json:"id"` // Campaign_Id
	Title           string   `json:"title"`
	Desc            string   `json:"desc"`
	PackageName     string   `json:"package_name"`
	IconUrl         string   `json:"icon_url"`
	ImageUrl        string   `json:"image_url"`
	ImpressionUrl   string   `json:"impression_url"`
	ClickUrl        string   `json:"click_url"`
	NoticeUrl       string   `json:"notice_url"`
	AppSize         string   `json:"app_size"`
	ClickMode       int      `json:"click_mode"`
	CCT             int      `json:"c_ct"`
	LinkType        int      `json:"link_type"`
	SubCategoryName []string `json:"sub_category_name"`
}

func RenderV2Res(mr MobvistaResult, r *mvutil.RequestParams) V2Result {
	var result V2Result
	result.Stat = mr.Status
	result.Msg = mr.Msg
	result.OnlyImpressionUrl = mr.Data.OnlyImpressionURL
	for _, v := range mr.Data.Ads {
		data := renderV2Data(v, *r)
		result.Data = append(result.Data, data)
	}
	return result
}

func renderV2Data(ad Ad, r mvutil.RequestParams) V2Data {
	var Data V2Data
	Data.Id = ad.CampaignID
	Data.Title = ad.AppName
	Data.Desc = ad.AppDesc
	Data.PackageName = ad.PackageName
	Data.IconUrl = ad.IconURL
	Data.ImageUrl = ad.ImageURL
	Data.ImpressionUrl = ad.ImpressionURL
	Data.ClickUrl = ad.ClickURL
	Data.NoticeUrl = ad.NoticeURL
	Data.AppSize = ad.AppSize
	Data.ClickMode = ad.ClickMode
	Data.CCT = ad.ClickCacheTime
	Data.LinkType = ad.CampaignType
	Data.SubCategoryName = append(Data.SubCategoryName, ad.SubCategoryName...)
	return Data
}
