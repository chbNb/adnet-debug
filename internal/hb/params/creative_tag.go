package params

type ExtCreativeNew struct {
	PlayWithoutVideo int    `json:"pwv,omitempty"` // playable_ads_without_video
	VideoEndType     int32  `json:"vet,omitempty"` // VideoEndType
	TemplateGroupId  *int   `json:"t_group,omitempty"`
	EndScreenId      string `json:"es_id,omitempty"`
	IsCreativeNew    bool   `json:"is_new,omitempty"`
}
