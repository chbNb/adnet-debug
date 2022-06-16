package mvutil

type SspProfitDistributionRule struct {
	Type            int8           `bson:"type,omitempty" json:"type,omitempty"`
	PlusMax         int            `bson:"plusMax,omitempty" json:"plusMax,omitempty"`
	FixedEcpm       float64        `bson:"fixedEcpm,omitempty" json:"fixedEcpm,omitempty"`
	UnitId          int64          `bson:"unitId,omitempty" json:"unitId,omitempty"`
	Area            string         `bson:"area,omitempty" json:"area,omitempty"`
	DailyDeficitCap map[string]int `bson:"dailyDeficitCap,omitempty" json:"dailyDeficitCap"`
	Updated         int64          `bson:"updated,omitempty" json:"updated,omitempty"`
	UpdatedDate     string         `bson:"updatedDate,omitempty" json:"updatedDate,omitempty"`
}
