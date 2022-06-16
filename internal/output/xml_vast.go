package output

import (
	"encoding/xml"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

// 判断是否遵循vast协议
func Isvast(params *mvutil.RequestParams) bool {
	// if params.Param.AdType == mvconst.ADTypeOnlineVideo && params.Param.IsVast {
	// 	return true
	// }
	// return false
	return params.Param.IsVast
}

type NOADS struct {
	value string `xml:",omitempty"`
}

type VAST struct {
	Ad      []*AD  `xml:"Ad,omitempty"`
	Version string `xml:"version,attr"`
}

type AD struct {
	Inline *Inline `xml:"InLine,omitempty"`
	Id     int64   `xml:"id,attr"`
}

type Inline struct {
	AdSystem    *AdSystem     `xml:"AdSystem,omitempty"`
	AdTitle     string        `xml:"AdTitle,omitempty"`
	Description string        `xml:"Description,omitempty"`
	Impression  []*Impression `xml:"Impression,omitempty"`
	Creatives   *Creatives    `xml:"Creatives,omitempty"`
	Extensions  *Extensions   `xml:"Extensions,omitempty"`
}

type Impression struct {
	Value string `xml:",cdata"`
}

type PopupMod struct {
	Value string `xml:",chardata"`
}

type Extensions struct {
	Extension *Extension `xml:"Extension,omitempty"`
}

type Extension struct {
	Type         string       `xml:"type,attr"`
	CTA          *CTA         `xml:"CTA,omitempty"`
	StarRating   float32      `xml:"StarRating,omitempty"`
	CommentNum   int          `xml:"CommentNum,omitempty"`
	SystemIcon   *SystemIcon  `xml:"SystemIcon,omitempty"`
	SkipMessage  *SkipMessage `xml:"SkipMessage,omitempty"`
	Click        *Click       `xml:"Click,omitempty"`
	VerticalView string       `xml:"VerticalView,omitempty"`
	PauseMod     string       `xml:"PauseMod,omitempty"`
	PopupMod     *PopupMod    `xml:"PopupMod,omitempty"`
	FullStar     *FullStar    `xml:"FullStar,omitempty"`
	HalfStar     *HalfStar    `xml:"HalfStar,omitempty"`
	EmptyStar    *EmptyStar   `xml:"EmptyStar,omitempty"`
	CTAButton    *CTAButton   `xml:"CTAButton,omitempty"`
	Review       *Review      `xml:"Review,omitempty"`
	Package      string       `xml:"Package,omitempty"`
}

type Review struct {
	Value string `xml:",cdata"`
}

type CTA struct {
	Value string `xml:",cdata"`
}

type FullStar struct {
	Value string `xml:",cdata"`
}

type HalfStar struct {
	Value string `xml:",cdata"`
}

type EmptyStar struct {
	Value string `xml:",cdata"`
}

type CTAButton struct {
	Value string `xml:",cdata"`
}

type VerticalView struct {
	Value string `xml:",chardata"`
}

type Click struct {
	Buttonexposure string `xml:"buttonexposure,attr"`
	Clickway       string `xml:"clickway,attr"`
	Value          string `xml:",chardata"`
}

type SkipMessage struct {
	Whenshown string `xml:"whenshown,attr"`
	Count     string `xml:"count,attr"`
	Value     string `xml:",chardata"`
}

type SystemIcon struct {
	Height       string `xml:"height,attr"`
	Width        string `xml:"width,attr"`
	CreativeType string `xml:"creativeType,attr"`
	Value        string `xml:",cdata"`
}

type Creatives struct {
	Creative []*Creative `xml:"Creative,omitempty"`
}

type Creative struct {
	Sequence     string        `xml:"sequence,attr"` // 每个广告素材应显示的数字顺序
	Id           string        `xml:"id,attr"`       // 广告素材的广告服务器定义的标识符
	Linear       *Linear1      `xml:"Linear,omitempty"`
	CompanionAds *CompanionAds `xml:"CompanionAds,omitempty"`
}

type CompanionAds struct {
	Companion []*Companion `xml:"Companion,omitempty"`
}

type Companion struct {
	Height         string           `xml:"height,attr"`
	Width          string           `xml:"width,attr"`
	ClickThrough   *ClickThrough    `xml:"CompanionClickThrough,omitempty"`
	StaticResource *StaticResource  `xml:"StaticResource,omitempty"`
	ClickTracking  []*ClickTracking `xml:"CompanionClickTracking,omitempty"`
}

type StaticResource struct {
	CreativeType *string `xml:"creativeType,attr,omitempty"`
	Value        *string `xml:",cdata"`
}

type Linear1 struct {
	SkipOffset     string          `xml:"skipoffset,attr"`
	Kipoffset      string          `xml:"Kipoffset,omitempty"`
	Duration       string          `xml:"Duration,omitempty"`
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty"`
	VideoClicks    *VideoClicks    `xml:"VideoClicks,omitempty"`
	MediaFiles     *MediaFiles     `xml:"MediaFiles,omitempty"`
}

type MediaFiles struct {
	MediaFile *MediaFile `xml:"MediaFile,omitempty"`
}

type MediaFile struct {
	Bitrate  int32  `xml:"bitrate,attr"`
	Delivery string `xml:"delivery,attr"`
	Height   int32  `xml:"height,attr"`
	Width    int32  `xml:"width,attr"`
	Type     string `xml:"type,attr"`
	Value    string `xml:",cdata"`
}

type TrackingEvents struct {
	Tracking []*Tracking `xml:"Tracking,omitempty"`
}

type Tracking struct {
	Event string `xml:"event,attr,omitempty"`
	Value string `xml:",cdata"`
}

type VideoClicks struct {
	ClickThrough  *ClickThrough    `xml:"ClickThrough,omitempty"`  // 点击时视频播放器在web浏览器窗口中打开此 URI
	ClickTracking []*ClickTracking `xml:"ClickTracking,omitempty"` // 用于在广告素材文件处理点击时跟踪点击
	// CustomClick CustomClick `xml:"CustomClick"`//元素用于跟踪线性广告素材中的其他非点击点击次数
}

type ClickTracking struct {
	Value string `xml:",cdata"`
}

type ClickThrough struct {
	ClickAvailDelay string `xml:"clickAvaildelay,attr,omitempty"`
	Value           string `xml:",cdata"`
}

// type Impression struct {
//	//Id    string `xml:"id,attr"`
//	Value string `xml:",cdata"`
// }

type AdSystem struct {
	Version string `xml:"version,attr"`
	Value   string `xml:",chardata"`
}

func RanderVastData(params *mvutil.RequestParams, resOnline OnlineResult) ([]byte, error) {
	unitId := params.Param.UnitID
	hasExtensions := false
	var extensionsType string
	var whenShown string
	var buttonexposure string
	var clickway string

	if IsAfreecatvUnit(unitId) {
		hasExtensions = true
		if !IsNewAfreecatvUnit(unitId) {
			whenShown = "00:00:00"
			buttonexposure = "1"
			clickway = "button"
		}
	}

	adSystemName := "Mintegral"
	adSystemVersion := "1.0"
	sequence := "1"
	kipoffset := 10
	delivery := "progressive"
	vast := VAST{Version: "2.0"}
	for _, value := range resOnline.Data.Ads {
		var ad AD
		ad.Id = value.CampaignID
		var inline Inline

		var adSystem AdSystem
		adSystem.Version = adSystemVersion
		adSystem.Value = adSystemName
		inline.AdSystem = &adSystem

		inline.AdTitle = value.AppName
		inline.Description = value.AppDesc

		var impression Impression
		impression.Value = value.ImpressionURL
		inline.Impression = append(inline.Impression, &impression)

		if value.AdTrackingPoint != nil && len(value.AdTrackingPoint.Impression) > 0 {
			for _, v := range value.AdTrackingPoint.Impression {
				var adTrackingImpression Impression
				adTrackingImpression.Value = v
				inline.Impression = append(inline.Impression, &adTrackingImpression)
			}
		}
		// 若为头条的广告，则将头条的impression存入Impression标签中
		// for _, impressionVal := range value.AdTrackingPoint.Impression {
		// 	inline.Impression = append(inline.Impression, impressionVal)
		// }

		var creative0, creative1 Creative
		creative0.Sequence = sequence

		if params.Param.ImageCreativeid != 0 { //原来的代码，不知道为什么这样写，判断的ImageCreativeid，而使用的是VideoCreativeid
			creative0.Id = strconv.FormatInt(params.Param.VideoCreativeid, 10)
		} else {
			creative0.Id = ""
		}
		// 凤凰视频召回的广告需遵循传来的MaxDuration，MinDuration
		if params.Param.RequestPath == mvconst.PATHIFENGADX {
			// 给adid赋值
			params.Param.IFENGAdId = value.CampaignID
			if value.VideoLength < params.Param.MinDuration ||
				value.VideoLength > params.Param.MaxDuration {
				return nil, errorcode.EXCEPTION_PARAMS_ERROR
			}
		}

		var linear Linear1
		if IsNewAfreecatvUnit(unitId) && value.VideoLength > 15 {
			linear.SkipOffset = "00:00:15"
		}

		// 判断广告返回类型
		if params.Param.AdFormat == mvconst.SKIPPABLE_LINEAR_ADS {
			linear.Kipoffset = time.Unix(int64(kipoffset), 0).Format("00:04:05")
		}
		linear.Duration = time.Unix(int64(value.VideoLength), 0).Format("00:04:05")

		// 组合ad_tracking
		var iconData, imageData Companion
		var TrackingEvents TrackingEvents
		for _, percentageVal := range value.AdTrackingPoint.Play_percentage {
			// 给joox下毒，不返回video strat
			prisonUnits, _ := extractor.GetUNIT_WITHOUT_VIDEO_START()
			if len(prisonUnits) > 0 && percentageVal.Rate == 0 {
				if mvutil.InInt64Arr(params.Param.UnitID, prisonUnits) {
					continue
				}
			}
			// percentageVal := playPercentage
			evenName := mvconst.GetVastPercentage(percentageVal.Rate)

			if len(evenName) > 0 {
				var tracking Tracking
				tracking.Event = evenName
				tracking.Value = percentageVal.Url
				TrackingEvents.Tracking = append(TrackingEvents.Tracking, &tracking)
			}
		}

		// 整理ad_tracking 的mute
		if len(value.AdTrackingPoint.Mute) > 0 {
			adTrackingMute := value.AdTrackingPoint.Mute
			for _, muteSub := range adTrackingMute {
				var tracking Tracking
				tracking.Event = "mute"
				tracking.Value = muteSub
				TrackingEvents.Tracking = append(TrackingEvents.Tracking, &tracking)
			}
		}
		linear.TrackingEvents = &TrackingEvents

		// video clicks
		var ClickThrough ClickThrough
		var VideoClicks VideoClicks
		ClickThrough.Value = value.ClickURL
		if IsNewAfreecatvUnit(unitId) {
			ClickThrough.ClickAvailDelay = "00:00:00"
		}
		VideoClicks.ClickThrough = &ClickThrough
		linear.VideoClicks = &VideoClicks
		if len(value.NoticeURL) > 0 {
			var clickTracking ClickTracking
			clickTracking.Value = value.NoticeURL
			linear.VideoClicks.ClickTracking = append(linear.VideoClicks.ClickTracking, &clickTracking)
		}
		if len(value.AdTrackingPoint.Click) > 0 {
			adTrackingClick := value.AdTrackingPoint.Click
			for _, ctVal := range adTrackingClick {
				var clickTracking ClickTracking
				clickTracking.Value = ctVal
				linear.VideoClicks.ClickTracking = append(linear.VideoClicks.ClickTracking, &clickTracking)
			}
		}

		// 给video 的mime默认值
		if len(value.VideoMime) <= 0 {
			value.VideoMime = "video/mp4"
		}
		var MediaFile MediaFile
		var MediaFiles MediaFiles
		MediaFile.Bitrate = value.Bitrate
		MediaFile.Delivery = delivery
		MediaFile.Height = value.VideoHeight
		MediaFile.Width = value.VideoWidth
		MediaFile.Type = value.VideoMime
		MediaFile.Value = value.VideoURL
		MediaFiles.MediaFile = &MediaFile
		linear.MediaFiles = &MediaFiles

		creative0.Linear = &linear

		if !hasExtensions {
			// companionads部分
			// icon mongo没有resolution，在此提供一个默认值
			if len(value.IconResolution) <= 0 {
				value.IconResolution = "128x128"
			}
			if len(value.IconMime) <= 0 {
				value.IconMime = "image/jpeg"
			}
			width, height := getWidthAndHeight(value.IconResolution)
			iconData.Width = width
			iconData.Height = height
			var StaticResource1, StaticResource2 StaticResource
			StaticResource1.CreativeType = &value.IconMime
			StaticResource1.Value = &value.IconURL
			iconData.StaticResource = &StaticResource1

			// creative1.Sequence = sequence
			// creative1.Id = string(params.Param.ImageCreativeid)

			// 大图
			width, height = getWidthAndHeight(value.ImageResolution)
			if len(width) > 0 && len(height) > 0 {
				imageData.Height = height
				imageData.Width = width
			} else {
				imageData.Height = "627"
				imageData.Width = "1200"
			}
			StaticResource2.CreativeType = &value.ImageMime
			StaticResource2.Value = &value.ImageURL
			imageData.StaticResource = &StaticResource2
			// 大图新增返回click_url
			imageData.ClickThrough = linear.VideoClicks.ClickThrough
			if len(linear.VideoClicks.ClickTracking) > 0 {
				imageData.ClickTracking = linear.VideoClicks.ClickTracking
			}
			creative1.Sequence = sequence
			if params.Param.ImageCreativeid != 0 {
				creative1.Id = strconv.FormatInt(params.Param.ImageCreativeid, 10)
			} else {
				creative1.Id = ""
			}
		}

		// 针对韩国开发者返回extension
		var count string
		var skipValue string
		var creativeUrls map[string]string
		if hasExtensions {
			// 判断是否支持pausemode（）
			// canReturnPausemod := canReturnPauseMod(params)
			// 若为cpm单子，则设置时长为15s，否则设为5
			if !IsNewAfreecatvUnit(unitId) {
				if value.CType == 3 {
					count = "00:00:15"
					skipValue = "Y"
					if value.VideoLength < 15 {
						skipValue = "N"
					}

				} else {
					count = "00:00:05"
					skipValue = "Y"
					if value.VideoLength < 5 {
						skipValue = "N"
					}
				}
			}
			var Extensions Extensions
			var Extension Extension
			extensionsType = "VideoADExtra"
			Extension.Type = extensionsType
			if !IsNewAfreecatvUnit(unitId) {
				var SkipMessage SkipMessage
				SkipMessage.Whenshown = whenShown
				SkipMessage.Count = count
				SkipMessage.Value = skipValue
				Extension.SkipMessage = &SkipMessage
				var Click Click
				Click.Buttonexposure = buttonexposure
				Click.Clickway = clickway
				Click.Value = "Y"
				Extension.Click = &Click
				var PopupMod PopupMod
				PopupMod.Value = "N"
				Extension.PopupMod = &PopupMod
				Extension.VerticalView = "N"
			}
			Extension.PauseMod = "Y"
			Extensions.Extension = &Extension
			inline.Extensions = &Extensions
		} else {
			// 因为针对所有开发者返回mintegral的logo，logo存放在配置中，所以需先获取配置
			gameloftConf, ok := extractor.GetGAMELOFT_CREATIVE_URLS()
			if !ok {
				mvutil.Logger.Runtime.Warnf("RanderVastData get GAMELOFT_CREATIVE_URLS error")
			}
			if params.Param.Platform == mvconst.PlatformIOS {
				creativeUrls = gameloftConf["2"]
			} else {
				creativeUrls = gameloftConf["1"]
			}
			var Extension Extension
			var Extensions Extensions
			Extension.Type = "Extra"
			var cta CTA
			cta.Value = value.CtaText
			Extension.CTA = &cta
			Extension.StarRating = value.Rating
			Extension.CommentNum = value.NumberRating
			// 应豆瓣开发者要求，vast返回时下发package，用于拼接storekit
			Extension.Package = value.PackageName

			var SystemIcon SystemIcon
			SystemIcon.Width = "94"
			SystemIcon.Height = "22"
			SystemIcon.CreativeType = "image/png"
			SystemIcon.Value = creativeUrls["systemIcon"]
			Extension.SystemIcon = &SystemIcon
			Extensions.Extension = &Extension
			inline.Extensions = &Extensions

			// 为gameloft开发者下毒
			if params.Param.PublisherID == 13026 {
				if len(creativeUrls) > 0 {
					var FullStar FullStar
					var EmptyStar EmptyStar
					var HalfStar HalfStar
					var CTAButton CTAButton
					var Review Review
					FullStar.Value = creativeUrls["fullStar"]
					EmptyStar.Value = creativeUrls["emptyStar"]
					HalfStar.Value = creativeUrls["halfStar"]
					CTAButton.Value = creativeUrls["ctaButton"]
					Review.Value = creativeUrls["review"]
					inline.Extensions.Extension.FullStar = &FullStar
					inline.Extensions.Extension.EmptyStar = &EmptyStar
					inline.Extensions.Extension.HalfStar = &HalfStar
					inline.Extensions.Extension.CTAButton = &CTAButton
					inline.Extensions.Extension.Review = &Review

				}
			}
		}

		var Creatives Creatives
		Creatives.Creative = append(Creatives.Creative, &creative0)

		if !hasExtensions {
			var CompanionAds CompanionAds
			CompanionAds.Companion = append(CompanionAds.Companion, &iconData)
			CompanionAds.Companion = append(CompanionAds.Companion, &imageData)
			creative1.CompanionAds = &CompanionAds

			Creatives.Creative = append(Creatives.Creative, &creative1)
		}
		inline.Creatives = &Creatives
		ad.Inline = &inline
		vast.Ad = append(vast.Ad, &ad)
	}

	output, err := xml.MarshalIndent(vast, ""+"\n", " ")
	myString := []byte(xml.Header + string(output))

	return myString, err

}

func getWidthAndHeight(resolution string) (width, height string) {
	arr := strings.Split(resolution, "x")
	if len(arr) >= 2 && len(arr[0]) > 0 && len(arr[1]) > 0 {
		width = arr[0]
		height = arr[1]
	}
	return
}

func VastReturnEmpty() ([]byte, error) {
	var noads NOADS
	output, err := xml.MarshalIndent(noads, ""+"\n", " ")
	myString := []byte(xml.Header + string(output))
	return myString, err
}

func canReturnPauseMod(params *mvutil.RequestParams) bool {
	returnPausemodConf := extractor.GetIS_RETURN_PAUSEMOD()
	if mvutil.InInt64Arr(params.UnitInfo.UnitId, returnPausemodConf.UnitIds) {
		return true
	}
	if returnPausemodConf.TotalRate == 100 {
		return true
	}
	rateRand := rand.Intn(100)
	if returnPausemodConf.TotalRate > rateRand {
		return true
	}
	return false
}

// 判断是否为afreecatv的unit
func IsAfreecatvUnit(unitId int64) bool {
	extensionsConf, _ := extractor.GetHAS_EXTENSIONS_UNIT()
	newAfreeConf := extractor.GetNEW_AFREECATV_UNIT()
	return mvutil.InInt64Arr(unitId, extensionsConf) || mvutil.InInt64Arr(unitId, newAfreeConf)
}

func IsNewAfreecatvUnit(unitId int64) bool {
	newAfreeConf := extractor.GetNEW_AFREECATV_UNIT()
	return mvutil.InInt64Arr(unitId, newAfreeConf)
}

func IsVastReturnInJson(param *mvutil.Params) bool {
	if !param.IsVast {
		return false
	}
	adnetConfList := extractor.GetADNET_CONF_LIST()
	if publisherList, ok := adnetConfList["vastReturnInJsonPub"]; ok {
		if mvutil.InInt64Arr(param.PublisherID, publisherList) {
			return true
		}
	}
	return false
}

func RenderVastReturnInJson(vastRes []byte, resOnline OnlineResult, req *mvutil.RequestParams) []byte {
	var or OnlineResult
	or.Status = 1
	or.Msg = "success"
	var oad OnlineAd
	for _, ad := range resOnline.Data.Ads {
		oad.DeepLink = ad.DeepLink
		oad.Vast = string(vastRes)
		oad.BidPrice = ad.BidPrice
		oad.CreativeId = ad.CreativeId
		oad.WinNoticeURL = req.BidWinUrl
	}
	or.Data.Ads = append(or.Data.Ads, oad)
	res, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(or)
	return res
}
