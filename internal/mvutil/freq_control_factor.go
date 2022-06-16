package mvutil

type FreqControlFactor struct {
	FactorKey  string  `bson:"factorKey,omitempty" json:"factorKey,omitempty"`
	FactorRate float64 `bson:"factorRate,omitempty" json:"factorRate,omitempty"`
	Status     int64   `bson:"status,omitempty" json:"status,omitempty"`
	Updated    int64   `bson:"updated,omitempty" json:"updated,omitempty"`
}
