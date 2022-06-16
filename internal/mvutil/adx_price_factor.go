package mvutil

type AdxPriceFactor struct {
	AdSourceId int     `bson:"adSourceId,omitempty" json:"adSourceId"`
	Key        string  `bson:"key,omitempty" json:"key"`
	Factor     float64 `bson:"factor,omitempty" json:"factor"`
	Updated    int64   `bson:"updated,omitempty" json:"updated"`
}
