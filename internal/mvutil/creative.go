package mvutil

import (
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
	"gopkg.in/mgo.v2/bson"
)

// CreativeInfo struct.
type CreativeInfo struct {
	ID               bson.ObjectId          `bson:"_id,omitempty" json:"_id,omitempty"`
	CampaignID       int64                  `bson:"campaignId,omitempty" json:"campaignId"`
	Status           int32                  `bson:"status,omitempty" json:"status"`
	Source           int32                  `bson:"source,omitempty" json:"source"`
	PackageName      string                 `bson:"packageName,omitempty" json:"packageName"`
	CountryCode      string                 `bson:"countryCode,omitempty" json:"countryCode"`
	SourceCreativeID int32                  `bson:"sourceCreativeId,omitempty" json:"sourceCreativeId"`
	Updated          int64                  `bson:"updated,omitempty" json:"updated,omitempty"`
	UpdatedDate      string                 `bson:"updatedDate,omitempty" json:"updatedDate"`
	Content          map[string]interface{} `bson:"content,omitempty" json:"content"`
}

type Content struct {
	Url             string      `bson:"url,omitempty" json:"url,omitempty"`
	VideoLength     int32       `bson:"videoLength,omitempty" json:"videoLength,omitempty"`
	VideoSize       int32       `bson:"videoSize,omitempty" json:"videoSize,omitempty"`
	VideoResolution string      `bson:"videoResolution,omitempty" json:"videoResolution,omitempty"`
	Width           int32       `bson:"width,omitempty" json:"width,omitempty"`
	Height          int32       `bson:"height,omitempty" json:"height,omitempty"`
	VideoTruncation int32       `bson:"videoTruncation,omitempty" json:"videoTruncation,omitempty"`
	WatchMile       int32       `bson:"watchMile,omitempty" json:"watchMile,omitempty"`
	BitRate         int32       `bson:"bitrate,omitempty" json:"bitrate,omitempty"`
	ScreenShot      string      `bson:"screenShot,omitempty" json:"screenShot,omitempty"`
	Value           interface{} `bson:"value,omitempty" json:"value,omitempty"`
	Resolution      string      `bson:"resolution,omitempty" json:"resolution,omitempty"`
	Mime            string      `bson:"mime,omitempty" json:"mime,omitempty"`
	Attribute       string      `bson:"attribute,omitempty" json:"attribute,omitempty"`
	AdvCreativeId   string      `bson:"advCreativeId,omitempty" json:"advCreativeId,omitempty"`
	CreativeId      int64       `bson:"creativeId,omitempty" json:"creativeId,omitempty"`
	FMd5            string      `bson:"fMd5,omitempty" json:"fMd5,omitempty"`
	Source          int32       `bson:"source,omitempty" json:"source,omitempty"`
	Orientation     int         `bson:"orientation,omitempty" json:"orientation,omitempty"`
	Protocal        int         `bson:"protocal,omitempty" json:"protocal,omitempty"`
}

var CreativeImageV2 []string = []string{
	ad_server.CreativeType_SIZE_320x50.String(),
	ad_server.CreativeType_SIZE_300x250.String(),
	ad_server.CreativeType_SIZE_480x320.String(),
	ad_server.CreativeType_SIZE_320x480.String(),
	ad_server.CreativeType_SIZE_300x300.String(),
	ad_server.CreativeType_SIZE_1200x627.String(),
	ad_server.CreativeType_JS_TAG_320x50.String(),
	ad_server.CreativeType_JS_TAG_300x250.String(),
	ad_server.CreativeType_JS_TAG_480x320.String(),
	ad_server.CreativeType_JS_TAG_320x480.String(),
	ad_server.CreativeType_JS_TAG_300x300.String(),
	ad_server.CreativeType_JS_TAG_1200x627.String(),
}

var CreativeVideo []string = []string{
	ad_server.CreativeType_VIDEO.String(),
	ad_server.CreativeType_JS_TAG.String(),
}

var CreativeEndcard []string = []string{
	ad_server.CreativeType_ENDCARD.String(),
}

var CreativeIcon []string = []string{
	ad_server.CreativeType_ICON.String(),
}

var CreativeAppName []string = []string{
	ad_server.CreativeType_APP_NAME.String(),
}

var CreativeAppDesc []string = []string{
	ad_server.CreativeType_APP_DESC.String(),
}

var CreativeAppRATE []string = []string{
	ad_server.CreativeType_APP_RATE.String(),
}

var CreativeCtaButton []string = []string{
	ad_server.CreativeType_CTA_BUTTON.String(),
}

var CreativeRating []string = []string{
	ad_server.CreativeType_COMMENT.String(),
}

var AppRating []string = []string{
	ad_server.CreativeType_APP_RATE.String(),
}

var Playable []string = []string{
	ad_server.CreativeType_PLAYABLE_URL.String(),
	ad_server.CreativeType_PLAYABLE_ZIP.String(),
}
