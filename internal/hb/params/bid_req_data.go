package params

import (
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

type BidReqData struct {
	S2SBidId                  string
	BidId                     string
	ClientIp                  string
	Platform                  int
	OsVersion                 string
	PackageName               string
	AppPackageName            string
	AppVersionName            string
	AppVersionCode            string
	Orientation               int
	Brand                     string
	Model                     string
	Oaid                      string
	OaidMd5                   string
	AndroidId                 string
	AndroidIdMd5              string
	Imei                      string
	ImeiMd5                   string
	Gaid                      string
	GaidMd5                   string
	Idfa                      string
	IdfaMd5                   string
	Idfv                      string
	MNC                       string
	MCC                       string
	NetworkType               int
	Language                  string
	UserAgent                 string
	SdkVersion                string
	ScreenSize                string
	EnImei                    string
	EnAndroidId               string
	AppId                     int64
	AppIdStr                  string
	UnitId                    int64
	UnitIdStr                 string
	UnitSize                  string
	InstallIds                string
	OpenIdfa                  string
	HttpReq                   int
	Channel                   int
	ApiVersion                string
	BidFloor                  float64
	ImpId                     string
	DeviceType                int
	Os                        string
	ScreenWidth               int32
	ScreenHeight              int32
	DidSha1                   string
	DidMd5                    string
	DpidSha1                  string
	DpidMd5                   string
	CountryCode               string
	CityName                  string
	CityCode                  string
	RegionString              string // ????????? eg: ca????????????
	ReplaceBrand              string
	ReplaceModel              string
	AppName                   string
	RealPackageName           string
	RequestType               int
	PublisherId               int64
	PublisherType             int
	AdType                    int
	AdTypeStr                 string
	VideoAdType               int
	FormatOrientation         int
	IfSupportSeperateCreative int32
	Scenario                  string
	IsLowDevice               bool
	DebugCountryCode          string
	Offset                    int32
	SysId                     string
	BkupId                    string
	ExtSysId                  string
	AsRequestData             string
	SupportDownload           bool
	StoreUrl                  string
	BundleId                  string
	AppCat                    []string
	Coppa                     int
	BidResp                   *mtgrtb.BidResponse
	Algorithm                 string
	Price                     float64
	Nbr                       int
	Token                     string
	BidRawData                string
	DspId                     int64
	AppFrequencyCap           int
	JumpTypeConfig            map[string]int32
	AppDevinfoEncrypt         int
	UnitEndCard               *smodel.EndCard
	NVTemplate                int32
	VideoEndType              int
	AppStorekitLoading        int32
	ExtflowTagId              int
	RandValue                 int
	OsVersionCode             int32
	FakeKeys                  map[int64]*FakeKey
	AdBackend                 string
	AdBackendData             string
	RemoteIp                  string
	ServerIp                  string
	Mac                       string
	AdNum                     int
	BackendConfig             string
	ThirdTemplate             string
	DspExt                    string
	ReqKeyName                string
	ExtsystemUseragent        string
	ReqBackend                string
	BidRejectCode             int
	RejectData                string
	PriceFactor               string
	UnitNVTemplate            int32
	ExcludeIdS                string
	TokenTimeStamp            string
	BannerUnitWidth           int32
	BannerUnitHeight          int32
	CloseId                   string
	DisplayInfos              []DisplayInfo
	RefreshTime               int64
	BidIsNotUSDCur            int
	Currency                  string
	CurrencyPrice             float64
	VideoWidth                int32
	VideoHeight               int32
	UnitBtClass               int
	AppBtClass                int
	BidTestMode               int32
	ExcludePackageNames       map[string]bool
	ImpExcludePkg             string
	ABTestDeviceKey           string
	ABTestDeviceVal           map[string][]byte
	UnitAlac                  int
	UnitAlecfc                int
	UnitMof                   int
	MofUnitId                 int64
	RemovePubImp              bool
	BigTemplateFlag           bool
	RandNum                   int32
	ReqType                   string
	PlacementId               int64
	BlackDomain               []string
	BlackBundle               []string
	UnitBlackPackageList      *[]string
	AppBlackPackageList       *[]string
	BlackIABCategory          map[string][]string
	Dmt                       float64
	Dmf                       float64
	Ct                        string
	PowerRate                 int
	Charging                  int
	TotalMemory               string
	ResidualMemory            string
	RankerInfo                string
	MediationName             string
	ToponChannelInfo          string
	Extra3                    string
	ExtAlgo                   string
	ExtAdxAlgo                string
	ExtData2Log               string
	ExtData                   ExtData
	DcoTestFlag               int32
	Tmax                      int32
	AsTestMode                int32
	WebEnvData                WebEnv // web???????????????????????????h5??????????????????(?????????????????????)
}

type DisplayInfo struct {
	CampaignId string `json:"cid"`
	RequestId  string `json:"rid"`
}

type ExtData struct {
	PriceFactor          float64 `json:"pf,omitempty"`      // ????????????- ????????????
	PriceFactorGroupName string  `json:"pf_g,omitempty"`    // ????????????- ???????????????
	PriceFactorTag       int     `json:"pf_t,omitempty"`    // ????????????- ???????????????1=A, 2=B, 3=B'
	PriceFactorFreq      *int    `json:"pf_f,omitempty"`    // ????????????- ????????????????????????
	PriceFactorHit       int     `json:"pf_h,omitempty"`    // ????????????- ???????????????????????? 1=?????????2=?????????
	Send2RS              int     `json:"pf_s2rs,omitempty"` // ????????????- ???????????????RS??? 1=?????????Hb?????????????????????2????????????HB????????????????????????
}

func (ed *ExtData) SetPriceFactor(pf float64) {
	if pf > constant.PriceFactor_MAXValue || pf <= constant.PriceFactor_MINValue {
		ed.PriceFactor = constant.PriceFactor_DefaultValue
	} else {
		ed.PriceFactor = pf
	}
}

type WebEnv struct {
	Webgl *int `json:"webgl"` // ????????????webgl 0????????????1????????????2????????????
}
