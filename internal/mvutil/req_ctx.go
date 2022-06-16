package mvutil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtg_hb_rtb"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adx_common/model"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/native"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

// BackendCtx a context with a request to a backend
type BackendCtx struct {
	AdReqKeyName  string
	AdReqKeyValue string
	AdReqPkgName  string
	PriceFactor   float64
	// BackendItem   *BackendItem
	// Templates     map[string]BackendTemplate
	Content  int
	Elapsed  int // request backend cost in ms
	Region   []string
	RespData []byte
	Ads      *corsair_proto.BackendAds
	// FakeKeys      map[int64]*FakeKeyData
	BlackList     map[int64]*PkgBlackList
	IsBidAdServer bool
	IsBidMAS      bool // adserver dsp
	AsAds         *corsair_proto.BackendAds
	ReqPath       string
	Method        string
	Tmax          int
}

type FakeKeyData struct {
	ReqKeyName  string
	ReqKey      string
	BtPkg       string
	PriceFactor float64
}

type PkgBlackList struct {
	PkgNames []string
}

// MasParams MAS 将 p q r z al 参数自己拼接后返回，存到这里， 在拼url时优先使用这里的
type MasParams struct {
	Al       string
	P        string
	Q        string
	R        string
	Z        string
	K        string
	Oz       string
	SdkParam *mtgrtb.BidResponse_SdkParam
}

type RequestQueryMap map[string][]string

type RequestParams struct {
	Method         string // GET ,POST, PUT, DELETE
	Param          Params
	QueryMap       RequestQueryMap
	UnitInfo       *smodel.UnitInfo
	AppInfo        *smodel.AppInfo
	PublisherInfo  *smodel.PublisherInfo
	PlacementInfo  *smodel.PlacementInfo
	DebugInfo      string
	FlowTagID      int
	RandValue      int
	DspExt         string
	DspExtData     *model.DspExt
	IsBidRequest   bool
	IsFakeAs       bool
	DPrice         sync.Map
	ReqDPrice      []string
	PriceFactor    sync.Map
	ReqPriceFactor []string
	Adchoice       *native.Adchoice           // 第三方返回的 Adchoice
	AsResp         *mtgrtb.BidResponse_AsResp // Mas 返回的 bidresponse.as_resp 字段， 只当有mas的量返回时才会有值
	OmSDK          []OmSDK
	// ImageCreativeId int64
	// VideaCreativeId int64
	FreqDataFromAerospike map[string][]byte
	IsMoreAsAds           bool // 是否是“请求了多个广告，三方dsp胜出，用as/mas补充剩下几个广告位”的情况
	IsHBRequest           bool
	Header                http.Header
	PostData              []byte
	BidCur                string
	BidRespID             string
	Nbr                   int
	BidRejectCode         int
	Token                 string
	Price                 float64
	PriceBigDecimal       string
	BidFloorBigDecimal    string
	BidIsNotUSDCur        int
	Currency              string
	CurrencyPrice         float64
	BidWinUrl             string
	BidPrice              float64
	LoadRejectCode        int
	BidSkAdNetwork        *mtg_hb_rtb.BidResponse_SeatBid_Bid_Ext_Skadn
	IsTopon               bool
	ToponRequest          *openrtb.BidRequest
	ToponResponse         *openrtb.BidResponse
	Cloud                 string // 服务器云商
	Region                string // 服务器云商地区
}

func (rp *RequestParams) GetDspExt() (*model.DspExt, error) {
	// 直接走pioneer的情况下，还是生成dspExt，避免panic的情况。
	if IsRequestPioneerDirectly(&rp.Param) {
		rp.DspExtData = &model.DspExt{}
		return rp.DspExtData, nil
	}

	if rp.DspExt == "" {
		return nil, errors.New("getDspExt: no DspExt")
	}
	if rp.DspExtData != nil {
		return rp.DspExtData, nil
	}
	var dspExt model.DspExt
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(rp.DspExt), &dspExt)
	if err != nil {
		return nil, err
	}
	rp.DspExtData = &dspExt
	return &dspExt, nil
}

// ReqCtx a context with a request
type ReqCtx struct {
	MaxTimeout int
	FlowTagID  int
	RandValue  int
	AdsTest    bool
	Elapsed    int // request cost in ms
	Finished   int
	// CompleteCh      chan bool
	Backends        map[int]*BackendCtx
	OrderedBackends []int
	// ReqBackends     []string
	// ParamList       *corsair_proto.QueryParam
	ReqParams *RequestParams
	Result    *corsair_proto.QueryResult_
	// CtxMutex      *sync.Mutex
	DebugModeInfo []interface{}
	IsWhiteGdt    bool
	// IsNativeVIdeo 这个native请求的返回是否是native video
	// adnet-madx-dsp 的native请求中有部分同时支持视频和图片，如果返回的是video, 标记为true
	IsNativeVideo bool
	RespData      []byte // madx的原始返回
}

func NewReqCtx() *ReqCtx {
	return &ReqCtx{
		MaxTimeout: 0,
		FlowTagID:  0,
		RandValue:  0,
		Elapsed:    0,
		Finished:   0,
		// ReqBackends:   make([]string, 0),
		// CompleteCh:    make(chan bool, 1),
		Backends: make(map[int]*BackendCtx),
		// CtxMutex:      new(sync.Mutex),
		DebugModeInfo: make([]interface{}, 0),
	}
}

// func (ctx *ReqCtx) OnBackendFinish() {
// 	ctx.CtxMutex.Lock()
// 	defer ctx.CtxMutex.Unlock()
//
// 	ctx.Finished++
//
// 	if ctx.Finished >= len(ctx.Backends) {
// 		ctx.CompleteCh <- true
// 	}
// }

func NewBackendCtx(keyName, reqKey, keyPag string, priceFactor float64, content int, region []string) *BackendCtx {
	return &BackendCtx{
		AdReqKeyName:  keyName,
		AdReqKeyValue: reqKey,
		AdReqPkgName:  keyPag,
		PriceFactor:   priceFactor,
		Content:       content,
		Elapsed:       -1,
		Region:        region,
		RespData:      make([]byte, 0),
		// FakeKeys:      make(map[int64]*FakeKeyData),
		BlackList: make(map[int64]*PkgBlackList),
		Ads:       corsair_proto.NewBackendAds(),
	}
}

func NewMobvistaCtx() *BackendCtx {
	return &BackendCtx{
		AdReqKeyName:  "",
		AdReqKeyValue: "",
		Content:       0,
		Elapsed:       -1,
		Region:        []string{},
		RespData:      []byte{},
		Ads:           corsair_proto.NewBackendAds(),
	}
}

// GetString 获取收入参数key代表的string值
func (v RequestQueryMap) GetString(key string, filterSpace bool, def ...string) (string, error) {
	if val, ok := v[key]; ok {
		rawString := strings.Join(val, "")
		rawString = strings.TrimSpace(rawString)
		if filterSpace {
			rawString = strings.Replace(rawString, "\t", "", -1)
			rawString = strings.Replace(rawString, "\n", "", -1)
			// rawString = strings.Trim(rawString, "\\n\\t")
		}
		return rawString, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return "", fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetUint8 获取收入参数key代表的uint8值
func (v RequestQueryMap) GetUint8(key string, def ...uint8) (uint8, error) {
	if val, ok := v[key]; ok {
		u64, err := strconv.ParseUint(strings.TrimSpace(strings.Join(val, "")), 10, 8)
		return uint8(u64), err
	}

	if len(def) > 0 {
		return def[0], nil
	}
	return 0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetInt 获取收入参数key代表的Int值
func (v RequestQueryMap) GetInt(key string, def ...int) (int, error) {
	if val, ok := v[key]; ok {
		return strconv.Atoi(strings.TrimSpace(strings.Join(val, "")))
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return 0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetUint32 获取输入参数key代表的uint32值
func (v RequestQueryMap) GetUint32(key string, def ...uint32) (uint32, error) {
	if val, ok := v[key]; ok {
		u64, err := strconv.ParseUint(strings.TrimSpace(strings.Join(val, "")), 10, 8)
		return uint32(u64), err
	}

	if len(def) > 0 {
		return def[0], nil
	}
	return 0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetInt64 获取输入参数key对应的Int64值
func (v RequestQueryMap) GetInt64(key string, def ...int64) (int64, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseInt(strings.TrimSpace(strings.Join(val, "")), 10, 64)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return -1, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetUint64 获取输入参数key对应的UInt64值
func (v RequestQueryMap) GetUint64(key string, def ...uint64) (uint64, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseUint(strings.TrimSpace(strings.Join(val, "")), 10, 64)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return 0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

// GetFloat 获取输入参数key对应的float64值
func (v RequestQueryMap) GetFloat(key string, def ...float64) (float64, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseFloat(strings.TrimSpace(strings.Join(val, "")), 64)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return 0.0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

func (v RequestQueryMap) GetBool(key string, def ...bool) (bool, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseBool(strings.TrimSpace(strings.Join(val, "")))
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return false, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

func RenderQueryMapV(str string) []string {
	var sList []string
	sList = append(sList, str)
	return sList
}

type OmSDK struct {
	VerificationParameters string `json:"verification_p"`
	EventtrackerUrl        string `json:"et_url"`
	VendorKey              string `json:"vkey"`
}
