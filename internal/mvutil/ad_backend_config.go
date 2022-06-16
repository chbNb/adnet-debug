package mvutil

type AdBackendConfigInfo struct {
	BackendId int64 `bson:"backendId,omitempty" json:"backendId"`
	//BackendName string                     `bson:"backendName,omitempty" json:"backendName"`
	Region  []string `bson:"region,omitempty" json:"region"`
	Content int      `bson:"content,omitempty" json:"content"`
	//Cooperation int                        `bson:"cooperation,omitempty" json:"cooperation"`
	Status        int                        `bson:"status,omitempty" json:"status"`
	Updated       int64                      `bson:"updated,omitempty" json:"updated"`
	AdReqKeys     map[string]string          `bson:"adReqKeys,omitempty" json:"adReqKeys"`
	AdReqPkgNames map[string]string          `bson:"adReqPkgNames,omitempty" json:"adReqPkgNames"`
	Templates     map[string]BackendTemplate `bson:"templates,omitempty" json:"templates"`
	OrigUnitIds   map[string][]int64         `bson:"origUnitIds,omitempty" json:"origUnitIds"`
}

type BackendTemplate struct {
	TemplateId    int32    `bson:"templateId,omitempty" json:"templateId"`
	EndCard       []string `bson:"endCard,omitempty" json:"endCard"`
	MiniCard      []string `bson:"miniCard,omitempty" json:"miniCard"`
	VideoTemplate []string `bson:"videoTemplate,omitempty" json:"videoTemplate"`
}
