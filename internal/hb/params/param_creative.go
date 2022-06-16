package params

type ParamCreative struct {
	Gid     string  `json:"gid"`
	Tpid    int     `json:"tpid"`
	Crat    int     `json:"crat"`
	AdvCrid int     `json:"adv_crid"`
	Icc     int     `json:"icc"`
	Glist   string  `json:"glist"`
	Pi      float64 `json:"pi"`
	Po      float64 `json:"po"`
	Dco     int     `json:"dco"`
	Cid     int64   `json:"cid,omitempty"`
	Cname   string  `json:"cr_name,omitempty"`
	CpdIds  string  `json:"cpd_ids,omitempty"`
}
