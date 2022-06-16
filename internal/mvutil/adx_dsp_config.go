package mvutil

type AdxDspConfig struct {
	DspID   int64   `bson:"dspId,omitempty" json:"dspId,omitempty"`
	Status  int     `bson:"status,omitempty" json:"status,omitempty"`
	Updated int64   `bson:"updated,omitempty" json:"updated,omitempty"`
	Target  *Target `bson:"target,omitempty" json:"target,omitempty"`
}

type Target struct {
	IsAdchoiceRequired int            `bson:"isAdchoiceRequired,omitempty" json:"isAdchoiceRequired,omitempty"`
	bundleType         int            `bson:"bundleType,omitempty" json:"bundleType,omitempty"`
	bundlePackage      []string       `bson:"bundlePackage,omitempty" json:"bundlePackage,omitempty"`
	SDKEcTemplate      map[string]int `bson:"sdkEcTemplate,omitempty" json:"sdkEcTemplate,omitempty"`
}
