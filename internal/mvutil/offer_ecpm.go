package mvutil

type OfferECPM struct {
	UnitCountry string               `bson:"unitCountry,omitempty" json:"unitCountry"`
	Value       []map[string]float32 `bson:"value,omitempty" json:"value"`
	Updated     int64                `bson:"updated,omitempty" json:"updated"`
}
