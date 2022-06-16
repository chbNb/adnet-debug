package mvutil

type MasAbtestConfig struct {
	Key     string          `bson:"key" json:"key,omitempty"`
	Updated int64           `bson:"updated" json:"updated,omitempty"`
	Value   []MasUnitAbtest `bson:"value" json:"value,omitempty"`
}

type MasUnitAbtest struct {
	UnitID    string   `bson:"unit_id" json:"unit_id,omitempty"`
	Key       string   `bson:"key" json:"key,omitempty"`
	Blacklist []string `bson:"blacklist" json:"blacklist,omitempty"`
	Rate      float64  `bson:"rate" json:"rate,omitempty"`
}

type SupportVideoBanner struct {
	Key     string `bson:"key" json:"key,omitempty"`
	Updated int64  `bson:"updated" json:"updated,omitempty"`
	Value   string `bson:"value" json:"value,omitempty"`
}
