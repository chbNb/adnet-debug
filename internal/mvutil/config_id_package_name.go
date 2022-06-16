package mvutil

type ConfigIDPackagename struct {
	UrlIDMD5    string `bson:"urlIdMd5,omitempty" json:"urlIdMd5"`
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
	Updated     int64  `bson:"updated,omitempty" json:"updated"`
	Status      int    `bson:"status,omitempty" json:"status"`
}
