package storage

import (
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
)

// BiddingKey struct mkv storage key
type BiddingKey struct {
	BiddingID string `json:"bid,omitempty"`
}

// BiddingCompressor of mkv bidding storage compressor
type BiddingCompressor struct{}

func NewBiddingCompressor() *BiddingCompressor {
	return &BiddingCompressor{}
}

func (c *BiddingCompressor) Compress(key interface{}) string {
	return key.(*BiddingKey).BiddingID
}

// BiddingVal struct of mkv storage value
type BiddingVal struct {
	BiddingRequest  *params.BiddingRequest `json:"bidRequest,omitempty"`
	BiddingResponse *mtgrtb.BidResponse    `json:"bidResponse,omitempty"`
}

// BiddingSerializer of mkv bidding storage serailizer
type BiddingSerializer struct{}

func NewBiddingSerializer() *BiddingSerializer {
	return &BiddingSerializer{}
}

// Marshal key
func (s *BiddingSerializer) Marshal(val interface{}) (buf []byte, err error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&val)
}

// Unmarshal key
func (s *BiddingSerializer) Unmarshal(buf []byte, val interface{}) (err error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(buf, &val)
}
