package mvutil

type AdvertiserInfo struct {
	AdvertiserId int64      `bson:"advertiserId,omitempty" json:"advertiserId"`
	Advertiser   Advertiser `bson:"advertiser,omitempty" json:"advertiser,omitempty"`
	Updated      int64      `bson:"updated,omitempty" json:"updated"`
}

type Advertiser struct {
	AdvertiserId   int64  `bson:"advertiserId,omitempty" json:"advertiserId"`
	Status         int    `bson:"status,omitempty" json:"status"`
	PublisherId    int64  `bson:"publisherId,omitempty" json:"publisherId"`
	Name           string `bson:"name,omitempty" json:"name"`
	IsShowAdChoice bool   `bson:"isShowAdChoice,omitempty" json:"isShowAdChoice"`
	AdLogolink     string `bson:"adLogolink,omitempty" json:"adLogolink"`
	AdchoiceIcon   string `bson:"adchoiceIcon,omitempty" json:"adchoiceIcon"`
	AdchoiceLink   string `bson:"adchoiceLink,omitempty" json:"adchoiceLink"`
	AdchoiceSize   string `bson:"adchoiceSize,omitempty" json:"adchoiceSize"`
	AdvLogo        string `bson:"advLogo,omitempty" json:"advLogo"`
	AdvName        string `bson:"advName,omitempty" json:"advName"`
	PlatformLogo   string `bson:"platformLogo,omitempty" json:"platformLogo"`
	PlatformName   string `bson:"platformName,omitempty" json:"platformName"`
}
