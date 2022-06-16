package mvutil

type PublisherInfo struct {
	PublisherId      int64            `bson:"publisherId,omitempty" json:"publisherId"`
	Publisher        Publisher        `bson:"publisher,omitempty" json:"publisher"`
	Updated          int64            `bson:"updated,omitempty" json:"updated"`
	JumptypeConfig   map[string]int32 `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	OffsetList       map[string]int32 `bson:"offsetList,omitempty" json:"offsetList,omitempty"`
	JumptypeConfigV2 map[string]int32 `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`
	// ExcludeAdvertiserIds []int            `bson:"excludeAdvertiserIds,omitempty" json:"excludeAdvertiserIds"`
}

type Publisher struct {
	PublisherId int64  `bson:"publisherId,omitempty" json:"publisherId"`
	Status      int    `bson:"status,omitempty" json:"status"`
	Apikey      string `bson:"apiKey,omitempty" json:"apiKey"`
	Type        int    `bson:"type,omitempty" json:"type"`
	// MvSourceStatus int    `bson:"mvSourceStatus,omitempty" json:"mvSourceStatus"`
	// Created        int    `bson:"created,omitempty" json:"created"`
	// Username       string `bson:"username,omitempty" json:"username"`
	// ForceDeviceId  int    `bson:"forceDeviceId,omitempty" json:"forceDeviceId"`
	// BlockCategory  []int  `bson:"blockCategory,omitempty" json:"blockCategory"`
}
