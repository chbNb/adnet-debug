package mvutil

import (
	"encoding/json"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type Params struct {
	RequestPath         string
	RequestID           string
	ClientIP            string
	ParamCIP            string
	RemoteIP            string
	ServerIP            string
	Platform            int
	PlatformName        string
	OSVersion           string
	OSVersionCode       int32
	PackageName         string
	AppPackageName      string // 实际传给as的packagename
	AppVersionName      string
	AppVersionCode      string
	Orientation         int
	Brand               string
	Model               string
	ReplaceBrand        string
	ReplaceModel        string
	RModel              string
	AndroidID           string
	AndroidIDMd5        string
	IMEI                string
	ImeiMd5             string
	MAC                 string
	GAID                string
	GAIDMd5             string
	IDFA                string
	IDFAMd5             string
	MNC                 string
	MCC                 string
	MCCMNC              string
	Carrier             string
	NetworkType         int
	NetworkTypeName     string
	Language            string
	TimeZone            string
	UserAgent           string
	SDKVersion          string
	GPVersion           string
	GPSV                string
	ScreenSize          string
	LAT                 string
	LNG                 string
	GPST                string
	GPSAccuracy         string
	GPSType             int
	PkgSource           string // mp传来的宿主应用包名，暂时无用
	D1                  string
	D2                  string
	D3                  string
	AppID               int64
	UnitID              int64
	PublisherID         int64
	PublisherType       int
	Sign                string
	Category            int
	AdNum               int32
	PingMode            int
	UnitSize            string
	ExcludeIDS          string
	Offset              int32
	SessionID           string
	IsNewSession        bool
	ParentSessionID     string
	OnlyImpression      int
	NetWork             string
	ImpressionImage     int
	AdSourceID          int
	AdType              int32
	AdTypeStr           string
	NativeInfo          string
	NativeInfoList      []NativeInfoEntry
	Template            int
	RealAppID           int64
	IsOffline           int
	Scenario            string
	TNum                int
	ImageSize           string
	ImageSizeID         int
	InstallIDS          string
	DisplayCIDS         string
	FrameNum            int
	RequireNum          int // v4接口迁移，表示一帧需多少条广告
	IDFV                string
	OpenIDFA            string
	IDFVOpenIDFA        string
	PriceFloor          float64
	HTTPReq             int32
	DVI                 string
	BaseIDS             string
	FlowTagID           int
	CountryCode         string
	CityCode            int64
	CityString          string
	RandValue           int
	ApiRequestNum       int32
	ApiCacheNum         int
	VersionFlag         int32
	VideoVersion        string
	VideoAdType         int
	VideoW              int32
	VideoH              int32
	AppName             string
	ScreenWidth         int
	ScreenHeigh         int
	HeaderUA            string
	UaFamily            string
	RealPackageName     string
	Extra               string
	Extra2              string
	Extra3              string
	Extra4              string
	Extra5              string // origin requestid
	Extra6              int
	Extra7              int
	Extra8              int
	Extra9              string
	Extra10             string
	Extra11             string
	Extra13             int32
	Extra14             int64
	Extra15             int
	Extra16             int
	Extra20             string
	Extra20List         []string
	PowerRate           int
	Charging            int
	TotalMemory         string
	ResidualMemory      string
	CID                 string
	RankerInfo          string
	ApiVersion          float64
	ApiVersionCode      int32
	DisplayCamIds       []int64
	UnSupportSdkTrueNum int32
	Network             int64
	RequestType         int
	ExtPlayableList     string
	IARst               int
	Gotk                bool
	MpReqDomain         string // mp 请求的域名，主要用于替换adywind的tracking域名
	RegionString        string // 州

	// 中间参数
	CreativeId         int64
	IsRVBack           bool
	EndcardCreativeID  int64
	CreativeDescSource int32
	AlgoMap            map[int64]string
	AdvertiserID       int32
	CampaignID         int64
	Domain             string
	TrackingCdnDomain  string
	JumpType           string
	LinkType           int
	OfferType          int32
	CUA                int
	BackendID          int32
	AdBackendConfig    string
	QueryR             string
	QueryP             string
	QueryZ             string
	QueryQ             string
	QueryQ2            string // 广告主trackingUrl用，不strreplace
	QueryAL            string
	QueryCSP           string
	QueryMP            string
	QueryCRs           QueryR
	QueryRsList        []QueryR
	QueryRs            string
	CC                 string
	BackendList        []int32
	EndcardUrl         string
	CreativeAdType     int
	AdvCreativeID      int
	CreativeName       string
	CsetName           string
	ThirdCidList       []string
	ExtPlayableArr     string
	IAADRst            int
	IAPlayableUrl      string
	IAOrientation      int
	IsNewUrl           bool
	FBFlag             int
	Fallback           bool
	PlayableFlag       bool            // playable 是否使用adserver
	NewCreativeFlag    bool            // 是否走素材三期逻辑
	NeedWebviewPrison  bool            // android webview >= 72 下毒
	NewMoreOfferFlag   bool            // 标记是否新more offer 逻辑
	NewMofImpFlag      bool            // 标记more offer是否走展示批量上报逻辑
	NewRTBFlag         bool            // 标记是否新 UC RTB
	MofAbFlag          bool            // 标记是否走新逻辑。1则走新逻辑，0或者空值则走旧逻辑。（目前仅提供给endcard用于决定是否传model参数值使用）
	TmpIds             map[string]bool //SDK上传本地保存的模板ID
	RTBSourceBytes     []byte

	Dmt               float64 // 设备内存总空间
	Dmf               float64 // 设备内存剩余空间
	Ct                string  // cpu type
	ChannelInfo       string  // topon私有API透传算法参数
	PlacementId       int64   // sdk 传递的placement id
	FinalPlacementId  int64   // 最终的placementid
	WebEnvData        WebEnv  // web做的环境检查，存放h5透传来的信息(目前仅安卓会传)
	Open              int     // 是否为开源版本
	TS                string  // 客户端时间戳
	SignWithTimeStamp string  // new sign
	Skadnetwork       *Skadnetwork
	//TrackingInfo        *TrackingInfo //设备动态track信息 -- 通过这些信息以及固定的设备信息，组成临时的有效track id
	BandWidth int64  //客户端的带宽（单位：bps）
	TKI       string // 设备动态track信息 通过这些信息以及固定的设备信息，组成临时的有效track id
	NewTKI    string // 新参数名及结构的tki
	OsvUpTime string // 操作系统更新时间
	Ram       string // 手机内存
	UpTime    string // ios 开机时间
	NewId     string // tki内记录的new_id
	OldId     string // tki内记录的old_id

	// ext参数
	Extbtclass           int
	Extfinalsubid        int64
	ExtMtgId             int64
	ExtdeleteDevid       int
	ExtinstallFrom       int
	Extstats             string
	ExtpackageName       string
	ExtnativeVideo       int
	ExtflowTagId         int
	ExtcdnType           string
	Extendcard           string
	ExtrushNoPre         int
	ExtdspRealAppid      int64
	ExtfinalPackageName  string
	Extnativex           int
	Extattr              int
	Extstats2            string
	Extadstacking        string
	Extctype             int32
	Extrvtemplate        int
	Extabtest1           int
	ExtcreativeCdnFlag   string
	ExtcreativeSizeFlag  string
	Extplayable          int
	PriceIn              float64 // give
	PriceOut             float64 // receive
	LocalCurrency        int     // 币种
	LocalChannelPriceIn  float64 // local币种出价
	ExtpriceIn           string
	ExtpriceOut          string
	UseAlgoPrice         bool
	AlgoPriceIn          float64
	AlgoPriceOut         float64
	Extb2t               int
	Extchannel           string
	ExtadvInstallTime    string
	ExtthreesInstallTime string
	CNDomainTest         int
	Extabtest2           int
	Extbp                string
	Extsource            int32
	Extreject            string
	Extalgo              string
	ExtthirdCid          string
	ExtifLowerImp        int32
	Extnvt2              int32
	ExttrueNvt2          string
	ExtsystemUseragent   string
	ExtappearUa          int
	ExtCDNAbTest         int
	ExtMpNormalMap       string
	ExtAdxRequestID      string
	ExtBrand             string
	ExtModel             string
	ExtCreativeNew       string
	ExtData              string
	ExtDataInit          ExtData
	ExtTagArrStr         string
	Extplayable2         int
	Extabtest3           string // creative 从mongo 还是redis中读取
	ExtSysId             string
	ExtApiVersion        string
	ExtData2             string                 // 因之前ext_data不支持放到q字段，所以新建ext_data2存放必须放到q参数里的信息。
	extData2             map[string]interface{} // 插入切量标记
	ExtAdxAlgo           string
	ExtData2MAS          string // 透传给MAS的extdata2字段
	ExtDeviceId          string // 设备ID扩展标识 json结构
	UseDynamicTmax       string
	DynamicTmax          int32
	MasResponseTime      int64
	BidDeviceGEOCC       string

	// midway 额外参数
	MWadBackend     string
	MWadBackendData string
	MWbackendConfig string
	MWRandValue     int
	MWFlowTagID     int
	Debug           int
	RequestKey      string
	MWplayInfo      string
	MwCreatvie      []int32
	// 算法现在debug模式
	DebugMode               bool
	ImpCap                  int
	IsBlockByImpCap         bool
	OnlyRequestThirdDsp     bool // 仅请求三方DSP,仅ADX支持
	RespOnlyRequestThirdDsp bool // (ADX)返回是否仅调用三方DSP
	// 11-unknown 4-phone 5-tablet
	DeviceType int
	// Req-Device-Type 0--Unknown 1-phone 2-tablet 3-smartTv
	ReqDeviceType int
	// sdkversioncode
	FormatSDKVersion supply_mvutil.SDKVersionItem
	FormatAdType     int32
	// UC接口参数
	UCResponseId string
	// 风行所需字段
	FunRequestId string
	FunImpId     string
	// 虎扑所需字段
	HupuRequestId string
	HupuImpId     string
	DealId        string
	// 标记是否为低配设备
	LowDevice bool
	// 是否使用vast协议
	IsVast bool
	// vast协议参数：当传入的ad_format为1时，则返回linear ads，当传入的ad_format为2时，则返回skippable linear ads。若没有传，则默认返回linear ads
	AdFormat int
	AID      string
	NCP      string // 小程序请求ads接口不做参数强校验

	// filter condition
	ExcludePackageNames map[string]bool
	ImpExcPkgNames      []string
	// aerospike device key
	DeviceKey string
	// aerospike device installed package
	DeviceInstalledPackages map[string][]byte

	VideoCreativeid int64
	ImageCreativeid int64
	// 咪咕所需回传字段
	MIGUId    string
	MIGUImpId string
	MIGUCur   string
	// 咪咕所需文字长度字段
	MIGUTitleLen int
	MIGUDescLen  int
	// 猎豹所需回传字段
	CMId    string
	CMImpId string

	RequestTime int64
	// 凤凰视频回传字段
	IFENGId    string
	IFENGImpId string
	IFENGTagId string
	IFENGAdId  int64
	// 360max所需字段
	MaxBid     string
	MaxDealId  int64
	TemplateId int32
	AdslotId   uint32
	MaxPrice   int32
	MaxImgW    int32
	MaxImgH    int32
	// 金山云所需回传字段
	KSYUNAdSlotID  string
	KSYUNSKey      string
	KSYUNRequestID string
	KSYUNMaxCPM    int64
	ResInHtml      bool
	MaxDuration    int
	MinDuration    int
	AllowType      []int
	OnlineReqId    string // online api对接回传的id，可共用

	// request接口校验sign使用
	TimeStamp string
	FcaSwitch bool // 频次控制整体开关
	DspExt    string
	DPrice    int64

	// format orientation
	FormatOrientation        int
	HasWX                    bool // 标记是否安装了微信
	UseDs                    bool
	SysId                    string // 自有id
	EncryptedSysId           string // adnet加密后的sysId
	NewEncryptedSysId        string
	EncryptedSysIdTimestamp  int64
	BkupId                   string // 备用自有id
	EncryptedBkupId          string // adnet加密后的bkupId
	EncryptedBkupIdTimestamp int64
	NewEncryptedBkupId       string
	RuId                     string // mapping server生成的ruid
	EncryptedRuid            string // 加密后的ruid

	Ndm            int // 标记jssdk是否需要jssdk专属域名
	ChetDomain     string
	MTrackDomain   string
	JssdkCdnDomain string

	AdSense     int // 0或不传：为普通mtg广告请求 1：返回广告数据，同时需要下发广告位点击exa_click
	BidFloor    float64
	BidPrice    float64
	BidsPrice   string
	PriceFactor string
	IsFakeAs    bool
	AdxBidFloor float64

	ExtReduceFillReq   int     // 是否为有效降填充请求 0.否 1.是
	ExtReduceFillResp  int     // 降填充请求是否为有效返回 0.否 1.是
	PlPkg              string  // 小米替包广告位传来的包名list
	FillBackendId      []int32 // mtg 被广告填充所用的backendid
	Mof                int     // 标记是否为more offer请求
	MofData            string  // 算法对于more offer需要的上下文信息。
	ParentUnitId       int64   // more offer的父unit_id（rv, iv unit_id）
	UcParentUnitId     int64   // 固定unit的more offer的父unit_id（rv, iv unit_id）
	H5Type             int     // 1表示endcard；2表示playable
	TplGroupId         int     // template group id,more offer请求需要
	ImgFMD5            string  // 算法需要的图片filemd5值
	NewMoreOfferParams string  // more offer 算法和h5需要的信息
	RequestURI         string  // 请求url，用于抓取请求url
	MofVersion         int     // more offer version >=2则表示走新展示上报逻辑
	// MofUnitId    int64 // more offer unit id
	ReqType                     string   // 记录sdk上报的req_type
	ReqTypeAABTest              int      // 记录req_type=3的aabtest结果
	CleanDeviceTest             int      // 清除空设备ID测试， 1 不清除 2 清除
	DisplayCampaignABTest       int      // SDK展示过单子对应包名是否继续召回广告abtest结果
	VcnABTest                   int      // vcn（缓存条数） abtest结果
	CtnSizeTag                  string   // ctnSize标记
	DisplayCampaignPackageNames []string // SDK展示过单子对应包名
	ThirdPartyInjectParams      string   // 三方注入参数
	MofType                     int      // 用于标记mof类型，关闭场景广告位2，more offer 为1

	ReduceFillList                []string
	ReqBackends                   []string
	BackendReject                 []string
	ExtReduceFillValue            string // redis fillrate value
	ExtCampaignTagList            map[int64]*CampaignTagInfo
	AdClickReplaceclickByClickTag bool // 是否走clickurl替换adtracking.click逻辑标记
	TemplateGroupId               int
	DisplayInfoData               []DisplayInfo
	NewDisplayInfoData            []NewDisplayInfo
	SdkBannerUnitWidth            int64 // sdk banner的unit_size的width
	SdkBannerUnitHeight           int64 // sdk banner的unit_size的height

	// algo experiment
	PubFlowExpectPrice      float64 // 开发者的交付价格
	FillEcpmFloor           float64 // 降填充底价（算法计算生成的 unit+cc 降填充控制的底价），不是在 portal 上配置的 eCPM Floor
	FillEcpmFloorKey        string  // redis fillrate key
	FillEcpmFloorVer        string  // redis fillrate update version
	ABTestTags              map[string]int
	ABTestTagStr            string
	ExtABTestList           string
	AsTestMode              int32                // as测试模式号
	DebugModeTimeout        int                  // debugmode超时时间
	ThirdPartyABTestRes     ThirdPartyABTestData // 记录三方广告主情况下abtest结果
	ThirdPartyABTestStr     string               // 三方abtest结果
	MvLine                  string               // 中国专线abtest标记
	MangoBid                string               // 芒果tv的bid
	MangoVersion            int                  // 芒果tv的version
	MangoMinPrice           int                  // 芒果tv的最低竞价价格
	OAID                    string               // 匿名设备标识符
	IMSI                    string               // 国际移动用户识别码
	H5Data                  string               // 用于h5传递业务标记，便于数据分析
	SupportAdChoice         bool
	WaterfallFallback       bool            // 是否waterfall兜底
	SupportMoattag          bool            // 是否waterfall兜底
	BigTemplateFlag         bool            // 标记是否切到大模板
	BigTemplateId           int64           // 大模板id
	BigTemplateSlotMap      map[int32]int64 // 大模板slot map
	ExtBigTplOfferData      string          // unit维度记录offerid及大图素材id，视频素材id，endcard素材id，视频模版id，endcard模版id,slotid
	ExtBigTplOfferDataList  string
	ExtSlotId               string   // slot id
	ExtBigTemId             string   // 大模板id
	BigTempalteAdxPvUrl     []string // 大模板adx的pv上报url
	BigTemplateUrl          string   // mas返回的大模板url
	ExtPlacementId          string   // placementId
	RandNum                 int32    // 内部测试使用,201能固定召回201模板
	ReduceFillConfig        *smodel.ConfigAlgorithmFillRate
	ReqFillEcpmFloor        float64
	PolarisFlag             bool
	PolarisTplData          string // unit维度的new_tpl_data。记录component template id，template type，template group id，apiframework
	PolarisCreativeData     string // 记录模版对应的素材信息。ComponentTypeId为99表示offer维度额外召回的特殊素材。
	TplGrayTag              string // 记录大模版在灰度时对应的id标记，用于进行新模版url上线的abtest效果分析。
	RespFillEcpmFloor       float64
	BigTplABTestParamsStr   string // 大模板abtest框架下发参数
	CrLangABTestTag         bool   // 素材lang abtest标记
	VideoDspTplABTestId     int    // 视频模版三方广告源abtest结果。
	EndcardDspTplABTestId   int    // endcard模版三方广告源abtest结果。
	EndcardTplId            string
	IosStorekitPoisonFlag   bool // ios sdk问题版本storekit下毒标记
	SupportTrackingTemplate bool // 是否支持tracking参数模板化

	StartMode             string            // 判断是否为灰度机器(启动模式)
	AdxStartMode          map[string]string // adx返回的启动模式
	StartModeTags         map[string]string // 启动模式标记
	StartModeTagsStr      string            // 启动模式标记
	IfSupDco              int32             // dco切量标记
	IsThirdPartyMoreoffer int               // 是否为三方dsp带来的more_offer

	// HB param
	HBS2SBidID    string
	HBBidTestMode int32
	// HBBidFloor     float64
	HBTmax         int32
	Os             string
	ImpID          string
	MediationName  string
	TokenTimeStamp string
	RegionName     string
	HBExtPfData    string
	LossReqFlag    bool // 是否记录了 loss request
	Algorithm      string
	Ccpa           int32
	BidUnitID      int64
	BidAdType      int
	// HB param end

	ReduceFill                   string
	IsLowFlowUnitReq             bool
	CreativeCompressData         map[ad_server.CreativeType]CreativeCompressDataMap
	ThirdPartyDspVideoCreativeId int64 // adnet生成的三方dsp 视频素材id
	// demand side context
	DemandContext           *demand.Context
	AsDebugParam            *ad_server.QueryParam
	MasDebugParam           *mtgrtb.BidRequest_AsInfo
	AsABTestResTag          string   // as abtest返回的实验结果
	JunoCommonLogInfoJson   string   // juno 实验记录的切量标记
	PioneerExtdataInfo      string   // pioneer透传过来的信息，为json string（本期记录人群包信息）
	PioneerOfferExtdataInfo string   // pioneer返回的信息，为json string，支持拓展透传,和pioneer_extdata_info的区别是它会分campaignid记录标记
	CountryBlackPackageList []string // 针对country的黑名单包名，这些包名不做召回
	RwPlus                  int32    // Reward plus开关
	YLHHit                  int      // hit click mode test
	VideoFmd5               string
	AdspaceType             int32 // 开发者广告位设置为插屏全屏or插屏半屏。1＝全屏 2＝半屏
	MaterialType            int32 // 开发者广告类型的 素材设置，0= 图片＋视频 1=图片 2=视频

	ReplacedImpTrackDomain   string // 替换后的展示trakcing域名
	ReplacedClickTrackDomain string // 替换后的点击trakcing域名

	OnlineApiBidPrice string // 返回给online api 开发者的bid_price

	ToponThirdPartyImpUrl    string // topon请求三方dsp的imp tracking url
	ToponThirdPartyNoticeUrl string // topon请求三方dsp的click 上报 tracking url
	ToponThirdPartyParamMP   string // topon 请求三方dsp的mp参数

	IsTowardAdx bool //debug参数

	SkadnetworkDataStr string // 带到trcking中的skadnetwork信息
	SkAdNetworkId      string
	SkCid              string
	SkTargetId         string
	SkNonce            string
	SkSourceId         string
	SkTmp              string
	SkSign             string
	SkNeed             int
	SkViewSign         string

	OnlineFilterByBidFloor bool   // 判断online api 是否被底价过滤掉
	MainDomain             string // jssdk校验的域名

	DeeplinkType int //1表示仅支持通过click_url跳转到应用形式的deeplink；2表示支持双链deeplink，即支持优先跳转deeplink，若deeplink无法打开，回退preivew_url，同时无论是否能打开均异步发送click_url。
	HtmlSupport  int
	IfSupH265    int32 // h265 abtest切量标记

	LimitTrk      string // 是否限制广告追踪
	Att           string // app 级别 idfa 授权状态
	Brt           string // 屏幕亮度
	Vol           string // 音量
	Lpm           string // 是否为低电量模式
	Font          string // 设备默认字体大小
	IsReturnWtick int    // 判断是否返回wtick

	NeedCreativeDataCIds string // 需要as记录传给rs的素材信息的单子

	PcmReportendpoint string // 服务端接受归因信息的端点地址
	FwType            int    // 1/2 1表示CN包，2表示海外包。MTG CN的请求为1
	HardwareModel     string // 硬件型号，如D22AP
	TargetIds         string // 指定召回的offerid

	DebugDspId        int      //用于debug onlineAPI切量, -1 = adserver,    6=adx->adserver, 13=adx->mas
	BtPriceOutPercent *float64 // bt 打折系数。打折：不管做不做bt，有配置，就要打折

	Ntbarpt    int // 通知栏常驻设置。枚举值 0和1。ntbarpt=1表示不常驻 为0或不下发此字段表示常驻。
	Ntbarpasbl int // 设置通知栏是否可针对apk下载执行暂停。枚举值 0和1。ntbarpasbl=1表示可暂停，为0或不下发此字段表示不可暂停
	AtatType   int // 表示控制anpk安装完成后的激活控制逻辑。 2或不下发此字段表示安装完成后无额外处理 1控制表示需要在检测到用户已经安装完成后，弹窗提示激活。弹窗提示仅在当前广告任务所属广告处于展示阶段下进行，广告内容移除后不进行提示。 0控制表示检测到已安装完成后，自动触发激活

	// mapping server
	MappingServerDebug     string // mapping server debug模式开关
	MappingServerResCode   string // mapping server返回码
	MappingServerMessage   string // 计算的中间结果
	MappingServerDebugInfo string // debug模式下的信息
	MappingServerFrom      string // adnet/hb

	StackList             string // 调用堆栈信息
	ClassNameList         string // delegate 遵守的类名继承链路信息
	ProtocolList          string // delegate 遵守的协议信息
	TrafficInfo           string // 聚合平台信息
	DecryptTrafficInfoStr string // 解密后的聚合平台信息
	FilterRequestReason   string // 请求被过滤原因

	ToponTemplateSupportVersion int32  // topon支持素材模版的版本，1表示支持
	MappingIdfa                 string // mapping 到的idfa
	DspMoreOfferInfo            string // dsp more_offer传来支持召回单子需要的相关信息
	GdprConsent                 string // 表示客户端设置的GDPR consent信息状态
	DspMof                      int    // 1则为dsp流量触发的more_offer，其他值为sdk流量触发的more_offer

	AppSettingId    string // app setting abtest 标记
	UnitSettingId   string // unit setting abtest 标记
	RewardSettingId string // reward setting abtest 标记

	MiSkSpt                    string            // 是否支持小米storekit。1表示支持，0表示不支持。-1表示未安装小米商店
	MiSkSptDet                 string            // 支持的detailstyle
	AdnLibABTestTags           map[string]string // adn lib库abtest标记
	ParentId                   string            // more_offer cache的parent_id
	MofRequestDomain           string            // more_offer的请求域名
	MoreOfferRequestId         string            // more_offer的requestid
	IsHBRequest                bool              // 是否为hb请求
	CachedCampaignIds          string            // 需要as记录传给rs的已经缓存在sdk的单子
	OnlineApiNeedOfferBidPrice string            // 表示算法是否需要针对hb online api请求的每个offer单独出价 1表示需要
	FixedEcpm                  float64           // unit+cc配置的fixed_ecpm
	EncryptHarmonyInfo         string            // 加密后的鸿蒙info
	DecryptHarmonyInfo         string            // 加密后的鸿蒙info
	TokenRule                  int32             // headerbidding流量，如果1表示该广告可以被应用在其他token中，比如其他token如果load，本地已有标记为1的有效缓存广告可直接回调load成功。不下发或2表示token各自独立。默认不下发。
	Dnt                        string            // 值为1时，用户自行退出个性化广告推荐的能力
	CNTrackingDomainTag        bool              // 有没有切量到cn tracking
	ParentAdType               string            // 对于dsp的流量，传递b,vre,vin 格式的值 对于sdk，传递sdk_banner，rewarded_video,interstitial_video
	ParentExchange             string            // 主unit 的adx名字

	GetAdsErrbackendList   []int  // 请求广告失败时的backend list
	GetAdsErr              string // 请求广告失败原因
	NewIVClearEndScreenUrl bool
	UseCdnTrackingDomain   int

	// ===== bid server 相关 =====
	IsHitRequestBidServer int32                // 是否命中请求BidServer切量
	BidServerCtx          *smodel.BidServerCtx // bid server 的返回结果
	BidServerAdxResponse  *mtgrtb.BidResponse  // bid server 命中时, adx 放回的原始响应, 用于在Load阶段重新拼接某些参数, 如果不是bid server 实验则为nil
}

type CreativeCompressDataMap struct {
	Url       string
	Fmd5      string
	VideoSize int32
}

type SDKVersionItem struct {
	SDKType        string
	SDKNumber      string
	SDKVersionCode int32
}

type QueryR struct {
	Gid      string  `json:"gid"`
	Tpid     int     `json:"tpid"`
	Crat     int     `json:"crat"`
	AdvCrid  int     `json:"adv_crid"`
	Icc      int     `json:"icc"`
	Glist    string  `json:"glist"`
	Pi       float64 `json:"pi"`
	Po       float64 `json:"po"`
	Dco      int     `json:"dco"`
	Cid      int64   `json:"cid,omitempty"`
	Cname    string  `json:"cr_name,omitempty"`
	CpdIds   string  `json:"cpd_ids,omitempty"`
	CsetName string  `json:"cset_name,omitempty"`
}

type NativeInfoEntrys []NativeInfoEntry

type NativeInfoEntry struct {
	AdTemplate int `json:"id"`
	RequireNum int `json:"ad_num"`
}

type NewDisplayInfos []NewDisplayInfo

type NewDisplayInfo struct {
	RequestId  string `json:"1"`
	CampaignId string `json:"2"`
}

type DisplayInfos []DisplayInfo

type DisplayInfo struct {
	CampaignId string `json:"cid"`
	RequestId  string `json:"rid"`
}

type ExtData struct {
	HasWX                             bool      `json:"has_wx,omitempty"`
	CctAbTest                         *int      `json:"cct,omitempty"`
	Alac                              int       `json:"alac,omitempty"`
	Alecfc                            int       `json:"alecfc,omitempty"`
	Mof                               int       `json:"mof,omitempty"`
	MofUnitId                         int64     `json:"mof_uid,omitempty"`
	ParentUnitId                      int64     `json:"parent_id,omitempty"`
	H5Type                            int       `json:"h5_t,omitempty"`
	IsThirdParty                      int       `json:"is_tp,omitempty"`
	UseAlgoPrice                      bool      `json:"useAlgoPrice,omitempty"` // 是否使用算法价格
	AlgoPriceIn                       float64   `json:"algoPriceIn,omitempty"`  // 算法自定义PriceIn
	AlgoPriceOut                      float64   `json:"algoPriceOut,omitempty"` // 算法自定义PriceOut
	PriceIn                           float64   `json:"priceIn,omitempty"`      // 原业务PriceIn
	PriceOut                          float64   `json:"priceOut,omitempty"`     // 原业务PriceOut
	ReqTypeTest                       int       `json:"rt_test,omitempty"`      // req_type aabtest结果
	CleanDeviceTest                   int       `json:"cldev_test,omitempty"`   // 清除 device id abtest
	IsMoreOffer                       int       `json:"is_mof,omitempty"`       // 判断是否为more offer请求
	DisplayCampaignABTest             int       `json:"display_test,omitempty"` // SDK展示过单子 abtest
	VcnABTest                         int       `json:"vcn_test,omitempty"`     // vcn(缓存条数) abtest
	CtnSizeTest                       string    `json:"ctn_size,omitempty"`     // ctnSize的值
	IsReplaceAdClick                  bool      `json:"is_rp_click,omitempty"`  // 标记是否有走替换adclick的逻辑
	AdjustS2S                         string    `json:"adjust_s2s,omitempty"`
	ClickWithUaLangTag                int       `json:"cwua_tag,omitempty"`    // 标记点击header里是否带ua和lang的abtest的结果
	H5Handle                          int       `json:"h5_handle,omitempty"`   // 标记是否由h5处理点击上报
	CloseAdTag                        string    `json:"clsad,omitempty"`       // 标记是否要出关闭场景广告
	MofType                           int       `json:"mof_type,omitempty"`    // 区分是more offer 还是 close button ad
	CrtRid                            string    `json:"crt_rid,omitempty"`     // 主offer的request_id
	ClickInAdTracking                 int       `json:"ciat,omitempty"`        // 标记是否使用adtracking跳转。1为使用adtracking跳，2为使用click_url跳
	ClickInServer                     int       `json:"cis,omitempty"`         // 标记是否走af 白名单通道（处理逻辑：click_url使用market地址，notice_url去除notice参数）。1为是，2为否
	OldClickMode                      string    `json:"old_cm,omitempty"`      // 修改前的clickmode值
	ECCDNTag                          int       `json:"ec_cdn,omitempty"`      // endcard CND 实验结果
	MvLine                            string    `json:"mv_l,omitempty"`        // 中国专线abtest标记
	ReturnAfTokenParam                int       `json:"rt_af_tk,omitempty"`    // 标记是否在click_url有返回appsflyer的token给sdk
	ExcludePkg                        string    `json:"exc_pkg,omitempty"`     // mtg 点击过不召回mtg广告
	H5Data                            string    `json:"h5_d,omitempty"`        // 记录h5传来的业务标记
	ExcludePsbPkg                     string    `json:"exc_psb_pkg,omitempty"` // third postback 不召回实验标记
	ExcludeAopPkg                     string    `json:"exc_aop_pkg,omitempty"` // analysis offline package 不召回实验标记
	CNTrackDomain                     int       `json:"cntd,omitempty"`        // 中国专线tracking域名切量标记
	IdfaTag                           int       `json:"idfa_t,omitempty"`      // idfa 切量tag
	GaidTag                           int       `json:"gaid_t,omitempty"`      // gaid 切量tag
	ImeiTag                           int       `json:"imei_t,omitempty"`      // imei 切量tag
	AndroidIdTag                      int       `json:"aid_t,omitempty"`       // android_id 切量tag
	ImeiMd5Tag                        int       `json:"imei_md5_t,omitempty"`  // imeiMd5 切量tag
	ServerUniqClickTime               int       `json:"suct,omitempty"`        // 服务端点击设备去重时间窗
	ImpExcludePkg                     string    `json:"imp_exc_pkg,omitempty"` // mtg 展示过不召回mtg广告
	ReduceFillMode                    int       `json:"rdfl_mode,omitempty"`   // 降填充控制模式 1 填充率 2 底价
	Dmt                               float64   `json:"dmt,omitempty"`         // 设备内存总空间
	Dmf                               float64   `json:"dmf,omitempty"`         // 设备内存剩余空间
	CpuType                           string    `json:"ct,omitempty"`          // cpu type
	AdnetStartMode                    string    `json:"netg,omitempty"`        // adnet start mode
	AdxStartMode                      string    `json:"adxg,omitempty"`        // adx start mode
	PriceFactor                       float64   `json:"pf,omitempty"`          // 频次控制- 价格系数
	PriceFactorGroupName              string    `json:"pf_g,omitempty"`        // 频次控制- 实验组名称
	PriceFactorTag                    int       `json:"pf_t,omitempty"`        // 频次控制- 实验标签，1=A, 2=B, 3=B'
	PriceFactorFreq                   *int      `json:"pf_f,omitempty"`        // 频次控制- 获取到当前的频次
	PriceFactorHit                    int       `json:"pf_h,omitempty"`        // 频次控制- 是否能命中概率， 1=命中，2=不命中
	Send2RS                           int       `json:"pf_s2rs,omitempty"`     // 频次控制- 是否发送给RS， 1=发送（Hb不处理价格），2不发送（HB需要处理价格）
	VideoDspTplABTest                 int       `json:"vdt,omitempty"`         // 视频模版三方广告源abtest结果。仅request日志记录
	EndcardDspTplABTest               int       `json:"ecdt,omitempty"`        // endcard模版三方广告源abtest结果。仅request日志记录
	IfSupDco                          int       `json:"dco,omitempty"`         // 标记dco切量结果
	ImpressionCap                     int       `json:"imp_c,omitempty"`       // placement的impressionCap
	ImpressionCapTime                 int64     `json:"imp_t,omitempty"`       // placement的impressionCap对应的时间点（TS)
	VideoCompressAbtestTag            int       `json:"v_cp,omitempty"`        // 视频压缩abtest 切量标记
	ImageCompressAbtestTag            int       `json:"i_cp,omitempty"`        // 大图压缩abtest 切量标记
	DemandLibABTest                   int       `json:"dla,omitempty"`         // demand lib abtest 1-> demand lib 2-> inside
	IconCompressAbtestTag             int       `json:"ic_cp,omitempty"`       // icon压缩abtest 切量标记
	TreasureBoxAbtestTag              int       `json:"tb_t,omitempty"`        // treasure box sdk方式读取数据 1-> treasuer box, 2->自管理方式
	RwPlus                            int32     `json:"rw_plus,omitempty"`     // 是否开启大模板召回开关
	V5AbtestTag                       string    `json:"v5_t,omitempty"`        // V5的实验标记， 5_5, 5_3, 或者控
	SDKOpen                           int       `json:"sdk_open,omitempty"`    // sdk-opensource
	ReplacedImpTrackDomainId          int       `json:"ritd,omitempty"`        // 替换后的展示tracking域名对应的id
	ReplacedClickTrackDomainId        int       `json:"rctd,omitempty"`        // 替换后的展示tracking域名对应的id
	BandWidth                         int64     `json:"bw,omitempty"`          // 带宽
	AdxAbTest                         AdxAbTest `json:"adx_t,omitempty"`       //  ADX 返回的abtest字段，用于统计数据使用
	SqsCollect                        int       `json:"sqs_c,omitempty"`       // 值为1表示设备信息需要发到sqs做展示，点击频次控制。值为0或无此key表示不发送到sqs。值为2表示按照原有的逻辑处理
	BigTemplateTag                    int       `json:"big_tt,omitempty"`      // 大模版切量标记。1为切量，2为没切量。
	FreExcludePkgList                 []string  `json:"fepl,omitempty"`        // 记录频次控制传递给as的包名
	ToponRequestId                    string    `json:"tprid,omitempty"`       // topon request id
	TKSysTag                          string    `json:"tkst,omitempty"`        // tracking 集群切量
	OsvUpTime                         string    `json:"osvut,omitempty"`       // 操作系统更新时间
	UpTime                            string    `json:"upt,omitempty"`         // iOS 开机时间（格式为时间戳）
	NewId                             string    `json:"tkinid,omitempty"`      // tki内记录的new_id
	OldId                             string    `json:"tkioid,omitempty"`      // tki内记录的oldid
	ImeiABTest                        int       `json:"imei_abt,omitempty"`    // imei abtest
	DeeplinkType                      int       `json:"dlt,omitempty"`         //1表示仅支持通过click_url跳转到应用形式的deeplink；2表示支持双链deeplink
	HtmlSupport                       int       `json:"hs,omitempty"`          //是否支持html输出，1=yes，2=no
	LimitTrk                          string    `json:"l_trk,omitempty"`       // 是否限制广告追踪
	Att                               string    `json:"att,omitempty"`         // app 级别 idfa 授权状态
	Brt                               string    `json:"brt,omitempty"`         // 屏幕亮度
	Vol                               string    `json:"vol,omitempty"`         // 音量
	Lpm                               string    `json:"lpm,omitempty"`         // 是否为低电量模式
	Font                              string    `json:"font,omitempty"`        // 设备默认字体大小
	IsReturnWtickTag                  int       `json:"irwt,omitempty"`        // 是否给sdk返回wtick=1
	FwType                            int       `json:"fw_t,omitempty"`        // 1/2 1表示CN包，2表示海外包。MTG CN的请求为1
	HardwareModel                     string    `json:"h_mod,omitempty"`       // 硬件型号
	Ntbarpt                           int       `json:"ntbarpt,omitempty"`
	Ntbarpasbl                        int       `json:"ntbarpasbl,omitempty"`
	AtatType                          int       `json:"atatType,omitempty"`
	AppsflyerUaABTestTag              int       `json:"af_uat,omitempty"` // af的uaabtest tag
	PackageReplace                    int       `json:"pkg_rp,omitempty"` // 1 替换 2 没替换 空或0无意义
	WTick                             int       `json:"wtick,omitempty"`  // 1 替换 2 没替换 空或0无意义
	NewClickmodeTag                   int       `json:"ncm,omitempty"`
	ClickmodeGroupTag                 string    `json:"cmgt,omitempty"`     // clickmode 切量维度标记
	ClickMode6NotInGpAndAppstore      int       `json:"cm6n,omitempty"`     // clickmode 6情况下，link_type非gp，appstore单子
	MappingIdfaTag                    string    `json:"mp_idfa,omitempty"`  // mapping idfa abtest标记
	GdprConsent                       string    `json:"gdpr_c,omitempty"`   // 表示客户端设置的GDPR consent信息状态 0(没有传值,旧版本)，1（true）2（false）3（unknown）
	TKCNABTestTag                     int       `json:"tkcn"`               // tracking cn 集群切量标记
	TKCNABTestAATag                   int       `json:"tkcnaa"`             // tracking cn 集群切量AA标记
	MappingIdfaCoverIdfaTag           string    `json:"mp_ici,omitempty"`   // mapping idfa 替换idfa的abtest标记
	AppSettingId                      string    `json:"a_stid,omitempty"`   // app setting abtest 标记
	UnitSettingId                     string    `json:"u_stid,omitempty"`   // unit setting abtest 标记
	RewardSettingId                   string    `json:"r_stid,omitempty"`   // reward setting abtest 标记
	MiskSpt                           string    `json:"misk_spt,omitempty"` // 是否支持小米storekit。1表示支持，0表示不支持。-1表示未安装小米商店
	MoreofferAndAppwallMvToPioneerTag string    `json:"maamtp,omitempty"`   // more_offer/appwall迁移aabtest标记
	PioneerHttpCode                   int       `json:"p_hcode,omitempty"`  // more_offer和appwall直接请求pioneer，非200的httpcode
	OnlineApiNeedOfferBidPrice        string    `json:"olnobp,omitempty"`   // 表示算法是否需要针对hb online api请求的每个offer单独出价 1表示需要
	TrackDomainByCountryCode          int       `json:"tdbcc,omitempty"`    // 根据country code选择的tracking域名
	ClickmodeGlobalConfigTag          string    `json:"cmgc_t,omitempty"`
	DecryptHarmonyInfo                string    `json:"dhm_info,omitempty"`   // 解密后的鸿蒙info
	LoadCDNTag                        int       `json:"lcdnt,omitempty"`      // load cdn tag
	TmaxABTestTag                     int       `json:"a_tmax,omitempty"`     // tmax abtest tag
	UseDynamicTmax                    string    `json:"dy_tmax,omitempty"`    // use dynamic tmax
	MasResponseTime                   int64     `json:"mas_rt,omitempty"`     // mas bid response timestamp
	DecryptTrafficInfoStr             string    `json:"dtraf_info,omitempty"` // 解密后的聚合信息
	Dnt                               string    `json:"sdk_dnt,omitempty"`
	AerospikeGzipEnable               int       `json:"aerospike_gzip_enable,omitempty"` // 是否开启Aerospike的Gzip压缩 0-关闭, 1-开启
	AerospikeRemoveRedundancyEnable   int       `json:"as_rm_red,omitempty"`             // Aerospike 开启冗余去除 abtest
	GoVersion                         string    `json:"go_v,omitempty"`                  // go升级ab
	MultiVcn                          int       `json:"m_vcn,omitempty"`                 // hb 聚合是否开启缓存大于1
	VcnCampaigns                      string    `json:"vcn_cids,omitempty"`              // hb 聚合缓存未展示的 campaign ids
	ReqTimeout                        int       `json:"timeout,omitempty"`               // backend timeout tag
	ExpIds                            string    `json:"expIds,omitempty"`                // algo abtest 实验id
	MpToPioneerTag                    string    `json:"mptpi,omitempty"`                 // mp 流量迁移aabtest标记
	GetAdsErrbackendList              []int     `json:"gaebl,omitempty"`                 // 请求广告失败时的backend list
	GetAdsErr                         string    `json:"gaerr,omitempty"`                 // 请求广告失败原因
	DeviceGEOCCMatch                  int       `json:"d_geo_cc,omitempty"`              // hb bid request device.geo.country 和 IP 信息的 country 是否一致, 0: device.geo.country 空, 1: 一致, 2: 不一致
	ThreeLetterCountry                string    `json:"t_cc,omitempty"`                  // 三位国家码
	HBSubsidyType                     int       `json:"subsidy_t,omitempty"`             // 扶持类型。普通流量：0 扶持流量+垂类保量 ：101 扶持流量+非垂类保量 ：102 冷启动流量：200
	BidServerTag                      int32     `json:"bs_tag,omitempty"`                // 0:未请求bid-server 1: bid-server出价竞价成功 2: bid-server回退pioneer竞价成功 3:命中了bid-server, 但是竞价失败了
	BidServerRsPrice                  float64   `json:"bs_rs_price,omitempty"`           // 二阶段load rs的精准出价
	LoadRejectCode                    int       `json:"lrjc,omitempty"`
	CdnTrackingDomainABTestTag        string    `json:"cdn_tdabt,omitempty"` // imp，click,only_imp 的cdn切量标记
}
type AdxAbTest map[string]map[string]string

type ThirdPartyABTestData struct {
	ClickmodeRes               int     `json:"cm_res,omitempty"`    // mv dsp clickmode abtest 标记。1则为clickmode6，2则返回默认值0
	CNTrackDomain              int     `json:"cntd,omitempty"`      // 中国专线tracking域名切量标记
	VideoDspTplABTest          int     `json:"vdt,omitempty"`       // 视频模版三方广告源abtest结果。
	EndcardDspTplABTest        int     `json:"ecdt,omitempty"`      // endcard模版三方广告源abtest结果。
	PriceFactor                float64 `json:"pf,omitempty"`        // 频次控制- 价格系数
	PriceFactorGroupName       string  `json:"pf_g,omitempty"`      // 频次控制- 实验组名称
	PriceFactorTag             int     `json:"pf_t,omitempty"`      // 频次控制- 实验标签，1=A, 2=B, 3=B'
	PriceFactorFreq            *int    `json:"pf_f,omitempty"`      // 频次控制- 获取到当前的频次
	PriceFactorHit             int     `json:"pf_h,omitempty"`      // 频次控制- 是否能命中概率， 1=命中，2=不命中
	ImpressionCap              int     `json:"imp_c,omitempty"`     // placement的impressionCap
	ImpressionCapTime          int64   `json:"imp_t,omitempty"`     // placement的impressionCap对应的时间点（TS)
	TPDspVideoCrId             int64   `json:"dsp_crid,omitempty"`  // adnet生成的三方dsp 视频素材id
	YLHHit                     int     `json:"ylh_hit,omitempty"`   // 命中clickmode实验
	V5AbtestTag                string  `json:"v5_t,omitempty"`      // V5的实验标记， 5_5, 5_3, 或者控
	BandWidth                  int64   `json:"bw,omitempty"`        // 带宽
	TKSysTag                   string  `json:"tkst,omitempty"`      // tracking 集群切量
	EndcardTplId               int     `json:"ec_tpid,omitempty"`   // ec 模版abtest标记（三方dsp模版）
	RvTplId                    int     `json:"rv_tpid,omitempty"`   // 视频 模版abtest标记（三方dsp模版）
	TplCreativeDomainTag       string  `json:"tcdt,omitempty"`      // 记录模版url上需要替换域名的宏的标记
	TKCNABTestTag              int     `json:"tkcn"`                // tracking cn 集群切量标记
	SDKVideoFullClick          int     `json:"vid_fc"`              // 三方dsp全屏可点切量标记
	AdspaceType                int32   `json:"adspace_type"`        // 1＝全屏 2＝半屏
	MaterialType               int32   `json:"material_type"`       // 开发者广告类型的 素材设置，0= 图片＋视频 1=图片 2=视频
	TemplateType               int32   `json:"iv_tpl_t"`            // 新插屏素材组合类型
	ThirdPartyDspTplTag        int     `json:"tpdtt"`               // 三方dsp 模版获取abtest标记
	CdnTrackingDomainABTestTag string  `json:"cdn_tdabt,omitempty"` // imp，click,only_imp 的cdn切量标记
}

type ExtPriceFactorData struct {
	PriceFactor          float64 `json:"pf,omitempty"`                    // 频次控制- 价格系数
	PriceFactorGroupName string  `json:"pf_g,omitempty"`                  // 频次控制- 实验组名称
	PriceFactorTag       int     `json:"pf_t,omitempty"`                  // 频次控制- 实验标签，1=A, 2=B, 3=B'
	PriceFactorFreq      *int    `json:"pf_f,omitempty"`                  // 频次控制- 获取到当前的频次
	PriceFactorHit       int     `json:"pf_h,omitempty"`                  // 频次控制- 是否能命中概率， 1=命中，2=不命中
	Send2RS              int     `json:"pf_s2rs,omitempty"`               // 频次控制- 是否发送给RS， 1=发送（Hb不处理价格），2不发送（HB需要处理价格），
	AerospikeGzipEnable  int     `json:"aerospike_gzip_enable,omitempty"` // 是否开启Aerospike的Gzip压缩 0-关闭, 1-开启
}

type ReduceFillLog struct {
	CampaignID       string
	BackendID        int32
	IsReduceFill     bool    // 是否被降填充
	IsWhiteListDev   bool    // 是否在设备白名单列表
	IsAlgoExperiment bool    // 是否是算法实验
	ReduceEcpmFloor  float64 // 降填充取到的 ecpm floor 美元
	FillPrice        float64 // 返回广告的价格 美元
	ReduceFillKey    string  // 降填充取词表的key
	Version          string  // key 版本
}

type CampaignTagInfo struct {
	CampaignID        int64
	CDNAbTest         int
	TemplateGroup     int64
	VideoTemplateId   int64
	EndCardTemplateId int64
	VideoCreativeid   int64
	IsReduceFill      int64
}

type MofData struct {
	Vfmd5  string `json:"v_fmd5,omitempty"`  // 视频素材fmd5值
	Ifmd5  string `json:"i_fmd5,omitempty"`  // 图片素材fmd5值
	CrtRid string `json:"crt_rid,omitempty"` // 主offer的request_id
}

type WebEnv struct {
	Webgl *int `json:"webgl"` // 是否支持webgl 0是默认，1是支持，2是不支持
}

type Skadnetwork struct {
	Ver      string   `json:"ver,omitempty"`
	Adnetids []string `json:"adnetids,omitempty"`
	Tag      string   `json:"tag,omitempty"` // 0表示未配置，1表示配置了大写形式的 id，2表示配置了小写形式的id
}

type TrackingInfo struct {
	Zone            string `json:"zone,omitempty"`
	OsVersionUpTime int64  `json:"os_version_up_time,omitempty"`
	Uptime          int64  `json:"uptime,omitempty"`
	Abstract        string `json:"abstract,omitempty"`
	NewId           string `json:"new_id,omitempty"`
	OldId           string `json:"old_id,omitempty"`
}

type ExtDeviceId struct {
	Ruid        string `json:"ruid,omitempty"`
	MappingIdfa string `json:"mpidfa,omitempty"`
}

func RenderSDKVersion(sdkversion string) SDKVersionItem {
	var item SDKVersionItem
	if !strings.Contains(sdkversion, "_") {
		item.SDKNumber = sdkversion
	} else {
		arr := strings.Split(sdkversion, "_")
		if len(arr[0]) > 0 {
			item.SDKType = arr[0]
		}
		if len(arr[1]) > 0 {
			item.SDKNumber = arr[1]
		}
	}
	typeArr := []string{"mi", "ma", "mal", "sa", "nxa", "nxi", "js", "mp"}
	if InStrArray(item.SDKType, typeArr) && strings.Contains(item.SDKNumber, ".") {

	} else {
		item.SDKNumber = sdkversion
	}
	item.SDKVersionCode = GetVersionCode(item.SDKNumber)
	return item
}

func IsMP(path string) bool {
	pathList := []string{
		mvconst.PATHMPAD,
		mvconst.PATHMPNewAD,
		mvconst.PATHMPADV2,
	}
	return InStrArray(path, pathList)
}

func IsMpPingmode0(r *RequestParams) bool {
	return IsMP(r.Param.RequestPath) && r.Param.PingMode == 0
}

func IsBigTemplate(bigTempalteId int64) bool {
	return bigTempalteId != mvconst.UnsupportBigTpl && bigTempalteId != mvconst.SupportButNotBigTpl
}

type AdvCreativeMap struct {
	AdvCreativeId      string
	AdvCreativeName    string
	AdvCreativeGroupId int64
}

func (p *Params) UnmarshalExtData2() {
	if len(p.ExtData2) == 0 {
		return
	}
	if p.extData2 != nil {
		return
	}
	extData2 := make(map[string]interface{})
	if err := json.Unmarshal([]byte(p.ExtData2), &extData2); err != nil {
		return
	}
	p.extData2 = extData2
	return
}

func (p *Params) MarshalExtData2() {
	if p.extData2 == nil {
		return
	}
	data, _ := json.Marshal(p.extData2)
	p.ExtData2 = string(data)
}

func (p *Params) SetExtData2(key string, value interface{}) {
	if p.extData2 == nil || len(p.extData2) == 0 {
		p.extData2 = make(map[string]interface{})
	}
	p.extData2[key] = value
}
