package mvutil

import (
	"errors"
	"strings"
	"time"

	"github.com/koding/multiconfig"
	"github.com/mae-pax/consul-loadbalancer/util"
	"github.com/spf13/viper"
	"github.com/ua-parser/uap-go/uaparser"
)

var localZone string
var localCloud string
var localRegion string

func InitRegion(region string) {
	localRegion = region
}

func InitZone(cloud string) {
	cloud = strings.Split(cloud, "-")[0]
	localCloud = cloud
	localZone = util.Zone(cloud)
}

func Region() string {
	return localRegion
}

func Cloud() string {
	if localCloud == "" {
		return "aws"
	}
	return localCloud
}

func Zone() string {
	if localZone == "" {
		// 默认aws
		InitZone("aws")
	}
	return localZone
}

// func UseConsulServices(region, service string) bool {
// 	ratio := extractor.GetUseConsulServicesV2Ratio(Cloud(), region, service)
// 		if ratio > 0 && ratio > rand.Float64() {
// }

type AdnetConfig struct {
	AreaConfig   *AreaConfig
	CommonConfig *CommonConfig
	Region       string
	Cloud        string
	MetricsConf  string
}

type AreaConfig struct {
	HttpConfig          HttpConfig
	RedisLocalConfig    RedisLocalConfig
	RedisClusterConfig  RedisClusterConfig
	RedisAlgoConfig     RedisAlgoConfig
	IpInfoClusterConfig IpInfoClusterConfig
	Service             ServiceConfig
	ExtraConfig         *ExtraConfig
	CpMongoConsulConfig *ConsulConfig
	AsConsulConfig      ConsulConfig
	AdxConsulConfig     ConsulConfig
	AsDsConsulConfig    ConsulConfig
	CtRedisConsulConfig ConsulConfig
	IpConfig            IpConfig
	LBConsulConfig      *Consul
}

type BackendConfig struct {
	Service ServiceConfig
}

type CommonConfig struct {
	LogConfig          LogConfig
	TrackConfig        TrackConfig
	DVIConfig          DVIConfig
	IpConfig           IpConfig
	DefaultValueConfig DefaultValueConfig
	ChetConfig         ChetConfig
	SampleConfig       SampleConfig
	MetricsConf        MetricsConf
}

type CorsairConfig struct {
	ExtraConfig ExtraConfig
	Service     ServiceConfig
}

type HttpConfig struct {
	Port                  int      `required`
	Region                string   `required`
	ServerIpUrl           string   `required`
	DefaultIconUrls       []string `required`
	PprofPath             string   `required`
	SDkVersions           []string `required`
	AndroidNot302Ver      []int64  `required`
	RateLimit             int      `required`
	MaxRateLimit          int      `required`
	DNSInterval           int      `required`
	ConfigCenterKey       string   `required`
	ToMongo               bool     `required`
	CommonPath            string   `required`
	MgoExtractorPath      string   `required`
	BackendPath           string   `required`
	RedisDecode           bool     `required`
	RuntimeLogRate        int      `required`
	IpRedisLimit          int      `required`
	MKVConf               string   `required`
	MKVConfSE             string
	MKVPKGConf            string
	UseCtRedisConsul      bool
	UseMongoConsul        bool   `required`
	TreasureBoxConfigPath string `required`
	Cloud                 string
	RegionName            string
	ConsuleConfig         string
	GeoConfig             string
}

type DVIConfig struct {
	DVIKeys           []string
	DVIMaps           []string
	AndroidExpVersion []int32
}

type DefaultValueConfig struct {
	TrueNumRewardVideo    int
	TrueNumFeedsVideo     int
	TrueNumOfferwall      int
	TrueNumInterstitalSDK int
	TrueNumAppwall        int
	TrueNumDefault        int
}

type LogConfig struct {
	ReqConf                 string `required`
	RunConf                 string `required`
	WatchConf               string `required`
	OutputFullReqRes        bool   `required`
	CreativeConf            string `required`
	AdserverConf            string `required`
	ReduceFillConf          string `required`
	LossRequestConf         string `required`
	DspCreativeDataConf     string `required`
	TreasureBoxConf         string `required`
	ConsulAdxConf           string `required`
	ConsulWatchConf         string `required`
	ConsulAerospikeConf     string `required`
	AerospikeConf           string `required`
	MappingServerConf       string `required`
	ConsulMappingServerConf string `required`
}

type TrackConfig struct {
	TrackHost     string `required`
	PlayTrackPath string `required`
}

type MetricsConf struct {
	MetricsConfPath string `required`
}

type IpConfig struct {
	FeatureCode  int    `required`
	APIID        int    `required`
	NetServerIP  string `required`
	TimeoutDelay int    `required`
	Expire       int    `required`
	GrpcAddress  string
}

type ChetConfig struct {
	ChetHost         string `required`
	VgChetHost       string `required`
	ImpressionPath   string `required`
	ClickPath        string `required`
	ToutiaoImpPath   string `required`
	ToutiaoClickPath string `required`
}

type ExtraConfig struct {
	ModifyInterval                int
	IntervalFactor                int
	UseExpiredMap                 bool
	EMBatchDeleteTime             time.Duration
	EMRetryAgrainSleepMicrosecond int64
	EMExpiredDefaultTime          int64
	ActiveDataCollecter           []string
	Mongo                         string
	TimeOut                       int
	ReadTimeOut                   int
	MaxPoolSize                   int
	UpdateOffset                  int
	Db                            string
	DbConfig                      []DbConfig `required`
}

type DbConfig struct {
	Collection string // 表名
	Index      int
}

type ServiceConfig struct {
	ServiceDetail []*ServiceDetail `required`
}

type ServiceDetail struct {
	Name      string // 服务名称
	ID        int    // 服务编号
	Timeout   int    // 服务请求超时时间
	Workers   int    // 工作协程数
	HttpURL   string // 服务请求httpURL
	HttpsURL  string // 服务请求httpsURL
	Path      string // 服务请求path
	Method    string // 请求服务的Http Method
	UseConsul bool
	ConsulCfg *ConsulConfig
}

type ConsulConfig struct {
	Cloud        string
	Address      string
	Service      string
	MyService    string
	Internal     int
	Timeout      int
	ServiceRatio float64
	CpuThreshold float64
}

type IpInfoClusterConfig struct {
	HostPort       string
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	PoolSize       int
}

type RedisAlgoConfig struct {
	HostPort       string
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	PoolSize       int
}

type RedisClusterConfig struct {
	HostPort       string
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	PoolSize       int
}

type RedisLocalConfig struct {
	HostPort       string
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	PoolSize       int
}

type SampleConfig struct {
	RandType   int
	RandFactor int
}

var Config *AdnetConfig

var UaParser *uaparser.Parser

var SDKVersions map[string]int32

var IphoneModels map[string]string

var ServerIP, RemoteIP string

// var UaParser *uaparser.Parser

// ----config loaded from mongo
type RVConfig struct {
	Key     string     `bson:"key" json:"key"`
	Value   RVTemplate `bson:"value,omitempty" json:"value"`
	Updated int64      `bson:"updated,omitempty" json:"updated"`
}

type EndScreenConfig struct {
	Key     string        `bson:"key" json:"key"`
	Value   OfferwallUrls `bson:"value,omitempty" json:"value"`
	Updated int64         `bson:"updated,omitempty" json:"updated"`
}

type RVTemplate struct {
	Id         int    `bson:"id,omitempty" json:"id"`
	Url        string `bson:"url,omitempty" json:"url"`
	Paused_url string `bson:"paused_url,omitempty" json:"paused_url"`
}

type OfferwallUrls struct {
	Http  OfferwallUrls_ `bson:"http,omitempty" json:"http"`
	Https OfferwallUrls_ `bson:"https,omitempty" json:"https"`
}

type OfferwallUrls_ struct {
	Rewardvideo_end_screen string `bson:"rewardvideo_end_screen,omitempty" json:"rewardvideo_end_screen"`
}

var Rewardvideo_end_screen EndScreenConfig
var DefRVTemplate RVConfig

// adserver test group
type AdServerTestConfig struct {
	Rate int `json:"rate"`
}

// var AdServerTestConfig AdServerTestConfig

func (ac *AdnetConfig) LoadConfigFile(confpath string) error {
	m := &multiconfig.TOMLLoader{Path: confpath}
	AreaConfig := new(AreaConfig)
	if err := m.Load(AreaConfig); err != nil {
		return errors.New("load config file:" + confpath + " " + err.Error())
	}
	ac.AreaConfig = AreaConfig
	InitZone(AreaConfig.HttpConfig.Cloud)
	InitRegion(AreaConfig.HttpConfig.RegionName)

	extraConfig, err := ParseExtraConfig(AreaConfig.HttpConfig.MgoExtractorPath)
	if err != nil {
		return errors.New("load config file:" + AreaConfig.HttpConfig.MgoExtractorPath + " " + err.Error())
	}
	ac.AreaConfig.ExtraConfig = extraConfig

	lbConsulConfig, err := ParseConsulCfg(ac.AreaConfig.HttpConfig.ConsuleConfig)
	if err != nil {
		return errors.New("load config file:" + ac.AreaConfig.HttpConfig.ConsuleConfig + " " + err.Error())
	}
	ac.AreaConfig.LBConsulConfig = lbConsulConfig

	m = &multiconfig.TOMLLoader{Path: AreaConfig.HttpConfig.CommonPath}
	CommonConfig := new(CommonConfig)
	if err := m.Load(CommonConfig); err != nil {
		return errors.New("load config file:" + AreaConfig.HttpConfig.CommonPath + " " + err.Error())
	}
	ac.CommonConfig = CommonConfig

	// m = &multiconfig.TOMLLoader{Path: AreaConfig.HttpConfig.BackendPath}
	// BackendConfig := new(BackendConfig)
	// if err := m.Load(BackendConfig); err != nil {
	// 	return errors.New("load config file:" + AreaConfig.HttpConfig.BackendPath + " " + err.Error())
	// }
	// ac.AreaConfig.Service.ServiceDetail = append(ac.AreaConfig.Service.ServiceDetail, BackendConfig.Service.ServiceDetail...)
	return nil
}

func ReadConfig(filePath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(filePath)
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseExtraConfig(fileName string) (*ExtraConfig, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	extraConfig := ExtraConfig{}
	extraConfig.ModifyInterval = viper.GetInt("extractor.modify_interval")
	extraConfig.IntervalFactor = viper.GetInt("extractor.interval_factor")
	extraConfig.Mongo = viper.GetString("extractor.mongo")
	extraConfig.Db = viper.GetString("extractor.db")
	extraConfig.TimeOut = viper.GetInt("extractor.conn_timeout")
	extraConfig.ReadTimeOut = viper.GetInt("extractor.read_timeout")
	extraConfig.MaxPoolSize = viper.GetInt("extractor.max_pool_size")
	extraConfig.UpdateOffset = viper.GetInt("extractor.offset")
	collections := viper.GetStringSlice("extractor.collections")
	for _, collection := range collections {
		dbConfig := DbConfig{Collection: collection}
		extraConfig.DbConfig = append(extraConfig.DbConfig, dbConfig)
	}
	extraConfig.UseExpiredMap = viper.GetBool("extractor.use_expire_map")
	extraConfig.EMBatchDeleteTime = viper.GetDuration("extractor.em_batch_delete_time")
	extraConfig.EMExpiredDefaultTime = viper.GetInt64("extractor.em_expired_default_time")
	extraConfig.EMRetryAgrainSleepMicrosecond = viper.GetInt64("extractor.em_retry_sleep_microsecond")
	extraConfig.ActiveDataCollecter = viper.GetStringSlice("extractor.active_em_collections")
	return &extraConfig, nil
}

func ParseSDKVersion(sdkVersions []string) bool {
	for _, ver := range sdkVersions {
		prefix := "empty"
		strVer := ""
		if strings.Contains(ver, "_") {
			verList := strings.Split(ver, "_")
			prefix = strings.ToLower(strings.TrimSpace(verList[0]))
			strVer = strings.TrimSpace(verList[1])
		} else {
			strVer = strings.TrimSpace(ver)
		}

		intVer, err := IntVer(strVer)
		if err != nil {
			return false
		}
		SDKVersions[prefix] = intVer
	}
	return true
}

func InitUaParser() error {
	parser, err := uaparser.NewWithOptions("./conf/adn_regexes.yaml", uaparser.EOsLookUpMode|uaparser.EUserAgentLookUpMode|uaparser.EDeviceLookUpMode, 10, 20, true, false)
	if err != nil {
		return err
	}
	UaParser = parser
	return nil
}

func init() {
	Config = new(AdnetConfig)
	SDKVersions = make(map[string]int32)
	IphoneModels = make(map[string]string, 55)
	// models
	IphoneModels["iPhone1,1"] = "iPhone 1"
	IphoneModels["iPhone1,2"] = "iPhone 3G"
	IphoneModels["iPhone2,1"] = "iPhone 3GS"
	IphoneModels["iPhone3,1"] = "iPhone 4"
	IphoneModels["iPhone3,2"] = "iPhone 4"
	IphoneModels["iPhone3,3"] = "iPhone 4"
	IphoneModels["iPhone4,1"] = "iPhone 4S"
	IphoneModels["iPhone5,1"] = "iPhone 5"
	IphoneModels["iPhone5,2"] = "iPhone 5"
	IphoneModels["iPhone5,3"] = "iPhone 5c"
	IphoneModels["iPhone5,4"] = "iPhone 5c"
	IphoneModels["iPhone6,1"] = "iPhone 5s"
	IphoneModels["iPhone6,2"] = "iPhone 5s"
	IphoneModels["iPhone7,1"] = "iPhone 6 Plus"
	IphoneModels["iPhone7,2"] = "iPhone 6"
	IphoneModels["iPhone8,1"] = "iPhone 6s"
	IphoneModels["iPhone8,2"] = "iPhone 6s Plus"
	IphoneModels["iPhone8,4"] = "iPhone SE"
	IphoneModels["iPhone9,3"] = "iPhone 7"
	IphoneModels["iPhone9,4"] = "iPhone 7 Plus"
	IphoneModels["iPad1,1"] = "iPad"
	IphoneModels["iPad2,1"] = "iPad 2"
	IphoneModels["iPad2,2"] = "iPad 2"
	IphoneModels["iPad2,3"] = "iPad 2"
	IphoneModels["iPad2,4"] = "iPad 2"
	IphoneModels["iPad3,1"] = "iPad 3"
	IphoneModels["iPad3,2"] = "iPad 3"
	IphoneModels["iPad3,3"] = "iPad 3"
	IphoneModels["iPad3,4"] = "iPad 4"
	IphoneModels["iPad3,5"] = "iPad 4"
	IphoneModels["iPad3,6"] = "iPad 4"
	IphoneModels["iPad2,5"] = "iPad Mini"
	IphoneModels["iPad2,6"] = "iPad Mini"
	IphoneModels["iPad2,7"] = "iPad Mini"
	IphoneModels["iPad4,1"] = "iPad Air"
	IphoneModels["iPad4,2"] = "iPad Air"
	IphoneModels["iPad4,3"] = "iPad Air"
	IphoneModels["iPad4,5"] = "iPad Mini 2"
	IphoneModels["iPad4,6"] = "iPad Mini 2"
	IphoneModels["iPad4,4"] = "iPad Mini 2"
	IphoneModels["iPad4,7"] = "iPad Mini 3"
	IphoneModels["iPad4,8"] = "iPad Mini 3"
	IphoneModels["iPad4,9"] = "iPad Mini 3"
	IphoneModels["iPad5,1"] = "iPad Mini 4"
	IphoneModels["iPad5,2"] = "iPad Mini 4"
	IphoneModels["iPad5,3"] = "iPad Air 2"
	IphoneModels["iPad5,4"] = "iPad Air 2"
	IphoneModels["iPod1,1"] = "iPod Touch 1"
	IphoneModels["iPod2,1"] = "iPod Touch 2"
	IphoneModels["iPod3,1"] = "iPod Touch 3"
	IphoneModels["iPod4,1"] = "iPod Touch 4"
	IphoneModels["iPod5,1"] = "iPod Touch 5"
	IphoneModels["iPod6,1"] = "iPod Touch 6"
	IphoneModels["i386"] = "32-bit Simulator"
	IphoneModels["x86_64"] = "64-bit Simulator"
}
