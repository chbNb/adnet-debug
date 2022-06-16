package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MangoResult struct {
	Version int        `json:"version"`  //Ctrl + V
	Bid     string     `json:"bid"`      //Ctrl + V
	ErrCode int        `json:"err_code"` //芒果是不论有无广告填充都需要返回的。
	Ads     []AdsMango `json:"ads"`
}

type AdsMango struct {
	SpaceId    string      `json:"space_id"`    //request_id
	Price      int         `json:"price"`       //注意这里！
	CreativeId string      `json:"creative_id"` //campaign_id
	Iurl       []IurlMango `json:"iurl"`
	Curl       []string    `json:"curl"`     //点击检测地址
	Duration   int         `json:"duration"` //视频时长
	Ctype      int         `json:"ctype"`    //素材类型
}

type IurlMango struct {
	Event int    `json:"event"` //曝光监测事件上报时间点
	Url   string `json:"url"`   //impression_url
}

func RenderNoRes(r *mvutil.RequestParams) MangoResult {
	var mango MangoResult
	mango.Version = r.Param.MangoVersion
	mango.Bid = r.Param.MangoBid
	mango.ErrCode = 204
	return mango
}

func RenderMangoRes(mr MobvistaResult, r *mvutil.RequestParams, creative map[string]string) MangoResult {
	var mango MangoResult
	//先填充bidresponse的变量：
	mango.Version = r.Param.MangoVersion
	mango.Bid = r.Param.MangoBid
	for _, v := range mr.Data.Ads {
		ads, res := renderMangoAds(v, *r, creative)
		if res {
			mango.Ads = append(mango.Ads, ads)
		}
	}
	// 有填充返回200，无填充返回204
	if len(mango.Ads) > 0 {
		mango.ErrCode = 200
	} else {
		mango.ErrCode = 204
	}

	return mango
}

func renderMangoAds(ad Ad, r mvutil.RequestParams, creative map[string]string) (AdsMango, bool) {
	var adlist AdsMango
	var Iurl IurlMango
	var Bidprice int
	// 优先处理创意ID。
	creativeId := strconv.FormatInt(ad.CampaignID, 10) //基数指代进制。
	value, ok := creative[creativeId]
	if !ok {
		mvutil.Logger.Runtime.Warnf("MANGO RTB get Creative not in whitelist.campaignId=[%d],requestId=[%s]", ad.CampaignID, r.Param.RequestID)
		return adlist, false
	}
	adlist.CreativeId = value
	//adlist.Price = int(ad.Price) // 在mongo配置上线之前使用这个内容
	onlinePriceFloorUnits, _ := extractor.GetONLINE_PRICE_FLOOR_APPID()

	price, ok := onlinePriceFloorUnits[strconv.FormatInt(r.Param.AppID, 10)]
	if !ok {
		Bidprice = 0
	} else {
		Bidprice = int(price)
	}
	if Bidprice < r.Param.MangoMinPrice {
		mvutil.Logger.Runtime.Warnf("MANGO filter by price.price=[%d],requestId=[%s]", r.Param.MangoMinPrice, r.Param.RequestID)
		return adlist, false
	}

	adlist.Price = Bidprice
	//开始处理广告唯一ID
	adlist.SpaceId = r.Param.RequestID //与oppo处理方法一致。
	adlist.Duration = ad.VideoLength
	// 贴片广告不需传title和desc
	//adlist.Title = ad.AppName
	//adlist.Desc = ad.AppDesc
	// duration 默认15s
	adlist.Duration = 15
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if mangoDuration, ok := adnConf["mangoDurtion"]; ok {
		adlist.Duration = mangoDuration
	}
	adlist.Ctype = 2 //因为是视频，所以写死。
	// 放置impression_url
	Iurl.Event = 0
	Iurl.Url = ad.ImpressionURL
	adlist.Iurl = append(adlist.Iurl, Iurl)
	// 开始处理Irul
	for _, v := range ad.AdTracking.Play_percentage {
		switch v.Rate {
		case 0:
			Iurl.Event = 0
			Iurl.Url = v.Url
			adlist.Iurl = append(adlist.Iurl, Iurl)
		case 25:
			Iurl.Event = 1
			Iurl.Url = v.Url
			adlist.Iurl = append(adlist.Iurl, Iurl)
		case 50:
			Iurl.Event = 2
			Iurl.Url = v.Url
			adlist.Iurl = append(adlist.Iurl, Iurl)
		case 75:
			Iurl.Event = 3
			Iurl.Url = v.Url
			adlist.Iurl = append(adlist.Iurl, Iurl)
		case 100:
			Iurl.Event = 4
			Iurl.Url = v.Url
			adlist.Iurl = append(adlist.Iurl, Iurl)
		}
	}
	adlist.Curl = append(adlist.Curl, ad.ClickURL)
	return adlist, true
}
