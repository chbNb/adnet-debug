package mvutil

type PlacementInfo struct {
	PlacementId         int64 `bson:"placementId,omitempty" json:"placementId"`
	ImpressionCap       int   `bson:"impressionCap,omitempty" json:"impressionCap"`
	ImpressionCapPeriod int   `bson:"impressionCapPeriod,omitempty" json:"impressionCapPeriod"`
	Status              int   `bson:"status,omitempty" json:"status,omitempty"`
	Updated             int64 `bson:"updated,omitempty" json:"updated"`
}
