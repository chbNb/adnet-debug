package mvutil

type ConfigAlgorithmFillRate struct {
	Platform    int     `bson:"platform,omitempty" json:"platform,omitempty"`
	PublisherId int     `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	AppId       int     `bson:"appId,omitempty" json:"appId,omitempty"`
	UnitId      int     `bson:"unitId,omitempty" json:"unitId,omitempty"`
	Area        string  `bson:"area,omitempty" json:"area,omitempty"`
	UpdatedDate string  `bson:"updatedDate,omitempty" json:"updatedDate,omitempty"`
	UniqueKey   string  `bson:"uniqueKey,omitempty" json:"uniqueKey"`
	Rate        int     `bson:"rate,omitempty" json:"rate"`
	Updated     int64   `bson:"updated,omitempty" json:"updated"`
	Status      int     `bson:"status,omitempty" json:"status"`
	EcpmFloor   float64 `bson:"ecpmFloor,omitempty" json:"ecpmFloor"`
	ControlMode int     `bson:"controlMode,omitempty" json:"controlMode"`
}
