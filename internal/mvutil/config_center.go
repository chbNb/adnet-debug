package mvutil

type ConfigCenter struct {
	Key     string                 `bson:"key,omitempty" json:"key"`
	Area    string                 `bson:"area,omitempty" json:"area"`
	Value   map[string]interface{} `bson:"value,omitempty" json:"value"`
	Updated int64                  `bson:"updated,omitempty" json:"updated"`
}

type TRACKING_DB struct {
	Write []string `bson:"write,omitempty" json:"write"`
}

type Fallback struct {
	Status            int `bson:"status,omitempty" json:"status"`                       // 总开关
	FillRateLimit     int `bson:"fillRateLimit,omitempty" json:"fillRateLimit"`         // 填充率阈值
	AvgCostLimit      int `bson:"avgCostLimit,omitempty" json:"avgCostLimit"`           // 平均响应时间阈值
	ActiveReqNumLimit int `bson:"activeReqNumLimit,omitempty" json:"activeReqNumLimit"` // adserver qpm超过多少才生效
	TestAdserverRate  int `bson:"testAdserverRate,omitempty" json:"testAdserverRate"`   // 降级时仍调用adserver概率
	AdserverLogRate   int `bson:"adserverLogRate,omitempty" json:"adserverLogRate"`     // adserver打日志概率
}

type MP_DOMAIN_CONF struct {
	SearchDomain  *string `bson:"searchDomain,omitempty" json:"searchDomain"`
	ReplaceDomain *string `bson:"replaceDomain,omitempty" json:"replaceDomain"`
}

type RateLimit struct {
	BucketSpeed int `bson:"bucketSpeed,omitempty" json:"bucketSpeed"`
	PreTokens   int `bson:"preTokens,omitempty" json:"preToken"`
}
