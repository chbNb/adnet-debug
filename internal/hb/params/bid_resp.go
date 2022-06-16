package params

type BidResp struct {
	Status int          `json:"status"`
	Msg    string       `json:"msg"`
	Data   *BidRespData `json:"data,omitempty"`
}

type BidRespData struct {
	Bid     string            `json:"bid"`
	Token   string            `json:"token"`
	Price   float64           `json:"price"`
	LossUrl string            `json:"ln"`
	WinUrl  string            `json:"wn"`
	Cur     string            `json:"cur"`
	Macors  map[string]string `json:"macors"`
}
