package mvutil

type AppPackageMtgID struct {
	MtgID      int64  `bson:"mtgId,omitempty" json:"mtgId,omitempty"`
	AppPackage string `bson:"appPackage,omitempty" json:"appPackage,omitempty"`
	Updated    int64  `bson:"updated,omitempty" json:"updated,omitempty"`
}
