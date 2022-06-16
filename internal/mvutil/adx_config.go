package mvutil

type AdxMediaConfig struct {
	UnitId           int64               `bson:"unitId,omitempty" json:"unitId"`
	EcpmFloorSwitch  int                 `bson:"ecpmFloorSwitch,omitempty" json:"ecpmFloorSwitch"`
	EcpmFloor        float64             `bson:"ecpmFloor,omitempty" json:"ecpmFloor"`
	DspWhiteList     []int               `bson:"dspWhiteList,omitempty" json:"dspWhiteList"`
	BlackIABCategory map[string][]string `bson:"blackIABCategory,omitempty" json:"blackIABCategory"`
	BlackDomain      []string            `bson:"blackDomain,omitempty" json:"blackDomain"`
	BlackBundle      []string            `bson:"blackBundle,omitempty" json:"blackBundle"`
	Status           int                 `bson:"status,omitempty" json:"status"`
	Updated          int64               `bson:"updated,omitempty" json:"updated"`
	EcpmFloors       map[string]float64  `bson:"ecpmFloors,omitempty" json:"ecpmFloors"`
}

type AdxTrafficMediaConfig struct {
	UnitId       int64    `bson:"unitId,omitempty" json:"unitId,omitempty"`
	TrafficType  int      `bson:"trafficType,omitempty" json:"trafficType,omitempty"`
	AdType       int64    `bson:"adType,omitempty" json:"adType,omitempty"`
	Mode         int      `bson:"mode,omitempty" json:"mode,omitempty"`
	DeviceId     []string `bson:"deviceId,omitempty" json:"deviceId,omitempty"`
	Area         string   `bson:"area,omitempty" json:"area,omitempty"`
	DspWhiteList []int64  `bson:"dspWhiteList,omitempty" json:"dspWhiteList,omitempty"`
	Status       int      `bson:"status,omitempty" json:"status,omitempty"`
	Updated      int64    `bson:"updated,omitempty" json:"updated,omitempty"`
}

type AdxConfigMediaEcpmFloor struct {
	UnitId  int64   `bson:"unitId,omitempty" json:"unitId"`
	Country string  `bson:"country,omitempty" json:"country"`
	Time    int     `bson:"time,omitempty" json:"time"`
	Ecpm    float64 `bson:"ecpm,omitempty" json:"ecpm"`
	Status  int     `bson:"status,omitempty" json:"status"`
	Updated int64   `bson:"updated,omitempty" json:"updated"`
}
