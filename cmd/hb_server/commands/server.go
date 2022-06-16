package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ua-parser/uap-go/uaparser"
	"gitlab.mobvista.com/ADN/adnet/internal/backend"
	"gitlab.mobvista.com/ADN/adnet/internal/consuls"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/geo"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/config"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/mlog"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mkv"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	creativeredis "gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_config"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	algoab "gitlab.mobvista.com/algo-engineering/abtest-sdk-go"
	"gitlab.mobvista.com/algo-engineering/abtest-sdk-go/pkg/types"
)

const (
	DevMode     string = "dev"
	TestMode    string = "test"
	ProductMode string = "product"
)

var Appversion string

func NewServe() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "serve start",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			// 逐个读取配置文件，将k-v存放到cfg中。ConfigDirFlag这些信息应该是在命令行输入的
			cfg, err := config.ParseConfig(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag))
			if err != nil {
				fmt.Println(err)
			}
			config.SetConfig(cfg)
			// 设置请求上下文
			req_context.GetInstance().SetRegion(viper.GetString(config.RegionFlag))
			req_context.GetInstance().SetCloud(viper.GetString(config.CloudFlag))
			req_context.GetInstance().SetCfg(cfg)
			parser, err := uaparser.NewWithOptions(config.UaParserConfig(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag)), uaparser.EOsLookUpMode|uaparser.EUserAgentLookUpMode|uaparser.EDeviceLookUpMode, 10, 20, true, false)
			if err != nil {
				fmt.Printf("init parser err:%s\n", err.Error())
				os.Exit(1)
			}
			req_context.GetInstance().SetUaParser(parser)
			process_pipeline.StartMode = viper.GetString(config.StartModeFlag) //？
			// 如果是productMode，调用阿里云提供的接口获取ecs实例的ip地址
			if strings.TrimSpace(viper.GetString(config.ModeFlag)) == ProductMode {
				serverIP, err := helpers.QueryServerIp(cfg.ServerCfg.ServerMetaUrl)
				if err != nil {
					fmt.Printf("query server ip error:%s\n", err.Error())
					os.Exit(1)
				}
				req_context.GetInstance().SetServerIp(serverIP)
			}
			// 配置aerospike缓存客户端
			if req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone {
				err = req_context.GetInstance().SetBidCacheClient(storage.Aerospike, config.AerospikeConfigWithZone(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag), mvutil.Zone()))
				if err != nil {
					fmt.Printf("init bid cache client error:%s\n", err.Error())
					os.Exit(1)
				}
			} else {
				err = req_context.GetInstance().SetBidCacheClient(storage.Aerospike, config.AerospikeConfig(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag)))
				if err != nil {
					fmt.Printf("init bid cache client error:%s\n", err.Error())
					os.Exit(1)
				}
			}
			// test mode ad cache
			err = req_context.GetInstance().SetAdCacheClient(storage.Aerospike, config.AerospikeConfig(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag)))
			if err != nil {
				fmt.Printf("init ad cache client error:%s\n", err.Error())
				os.Exit(1)
			}
			// geo客户端
			err = geo.SetGeoClient(cfg.GeoConfigPath)
			if err != nil {
				fmt.Printf("init ipinfo cache client error:%s\n", err.Error())
				os.Exit(1)
			}
			// err = req_context.GetInstance().SetAsCacheClient(storage.Aerospike, config.AerospikeConfigAs(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag)))
			// if err != nil {
			// 	fmt.Printf("init ad cache client error:%s\n", err.Error())
			// 	os.Exit(1)
			// }
			// 素材redis客户端
			err = req_context.GetInstance().SetCreativeCacheClient(config.CreativeCacheConfig(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag)), "creative")
			if err != nil {
				fmt.Printf("init creative cache client error:%s\n", err.Error())
				os.Exit(1)
			}
			// 根据配置文件，为每个变量配置一个log
			mlog.InitLogHandler(cfg.LogCfg.Bid,
				cfg.LogCfg.Load,
				cfg.LogCfg.Event,
				cfg.LogCfg.Watch,
				cfg.LogCfg.Creative,
				cfg.LogCfg.DspCreative,
				cfg.LogCfg.Runtime,
				cfg.LogCfg.Request,
				cfg.LogCfg.LRequest,
				cfg.LogCfg.ReqMonitor,
				cfg.LogCfg.TreasureBox,
				cfg.LogCfg.ConsulAerospike,
				cfg.LogCfg.ConsulAdx,
				cfg.LogCfg.ConsulWatch,
				cfg.LogCfg.ConsulAdnetAerospike,
				cfg.LogCfg.MappingServer,
				cfg.LogCfg.ConsulMappingServer,
				cfg.LogCfg.DeviceAerospike,
			)
			req_context.GetInstance().SetLoggers(mlog.Logger)

			watcher.Init(req_context.GetInstance().MLogs.Watch)

			// 初始化metrics
			region := viper.GetString(config.CloudFlag) + "-" + viper.GetString(config.RegionFlag)
			fileConfig := viper.GetString(config.ConfigDirFlag)
			if fileConfig == "" {
				fileConfig = "./config"
			}
			metricsPath := fmt.Sprintf("%s/%s/%s/metrics.yaml", fileConfig, viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag))
			if err := metrics.InitMetrics(region, Appversion, metricsPath, ""); err != nil {
				fmt.Println("initMetrics failure: ", err.Error())
				os.Exit(1)
			}
			// bid cache aerospike consul           NewConsulResolver
			if cfg.ConsulCfg != nil && cfg.ConsulCfg.Aerospike != nil && cfg.ConsulCfg.Aerospike.Enable {
				if err := req_context.GetInstance().SetAerospikeConsulBuild(); err != nil {
					fmt.Println("SetAerospikeConsulBuild failure: ", err.Error())
					os.Exit(1)
				}
				req_context.GetInstance().Cfg.AerospikeConsulBuild.SetLogger(mlog.Logger.ConsulAerospike)
				// 定期从配置中心更新本地配置？
				if err := req_context.GetInstance().Cfg.AerospikeConsulBuild.Start(); err != nil {
					fmt.Println("AerospikeConsulBuild Start failure: ", err.Error())
					os.Exit(1)
				}
			}
			// adx server consul				逻辑跟上面相似
			if cfg.ConsulCfg != nil && cfg.ConsulCfg.Adx != nil && cfg.ConsulCfg.Adx.Enable {
				if err := req_context.GetInstance().SetAdxConsulBuild(); err != nil {
					fmt.Println("SetAdxConsulBuild failure: ", err.Error())
					os.Exit(1)
				}
				req_context.GetInstance().Cfg.AdxConsulBuild.SetLogger(mlog.Logger.ConsulAdx)
				req_context.GetInstance().Cfg.AdxConsulBuild.SetWatcher(mlog.Logger.ConsulWatch)
				if err := req_context.GetInstance().Cfg.AdxConsulBuild.Start(); err != nil {
					fmt.Println("AdxConsulBuild Start failure: ", err.Error())
					os.Exit(1)
				}
			}
			// device tag aerospike consul		同上
			if cfg.ConsulCfg != nil && cfg.ConsulCfg.AdnetAerospike != nil && cfg.ConsulCfg.AdnetAerospike.Enable {
				consulResolver, err := mvutil.NewAdnetConsulResolver(cfg.ConsulCfg)
				if err != nil {
					fmt.Println("NewAdnetConsulResolver failure: ", err.Error())
					os.Exit(1)
				}
				consulResolver.SetLogger(mlog.Logger.ConsulAdnetAerospike)
				if err := consulResolver.Start(); err != nil {
					fmt.Println("AdnetConsulResolver Start failure: ", err.Error())
					os.Exit(1)
				}
				req_context.GetInstance().SetAdnetAerospikeConsulBuild(consulResolver)
				mkv.SetConsulBuild(consulResolver, cfg.ConsulCfg.AdnetAerospike)
			}

			// mapping server consul		同上
			if cfg.ConsulCfg != nil && cfg.ConsulCfg.MappingServer != nil && cfg.ConsulCfg.MappingServer.Enable {
				mappingServerResolver, err := mvutil.NewMappingServerConsulResolver(cfg.ConsulCfg)
				if err != nil {
					fmt.Println("NewMappingServerConsulResolver failure: ", err.Error())
					os.Exit(1)
				}
				// 记录consul log
				mappingServerResolver.SetLogger(mlog.Logger.ConsulMappingServer)
				if err := mappingServerResolver.Start(); err != nil {
					fmt.Println("NewMappingServerConsulResolver Start failure: ", err.Error())
					os.Exit(1)
				}
				consuls.SetMappingServerResolver(mappingServerResolver, cfg.ConsulCfg.MappingServer)
			}
			// abtest相关
			if cfg.ConsulCfg != nil {
				algoab.Init(types.Settings{
					App:                    "MTG",
					ConsulAddr:             cfg.ConsulCfg.Address,
					ConsulTransportTimeout: 3 * time.Second,
					Monitor:                utility.NewAlgoABTestMetric(),
				})
			}

			// TODO refactor hb and adnet
			if mvutil.Config.AreaConfig == (*mvutil.AreaConfig)(nil) {
				req_context.GetInstance().MLogs.Runtime.Info("init metadata config")
				// for hb get configcenter tracking domain
				configCenterArea := helpers.ConfigCenterKey(req_context.GetInstance().Region)
				adnCfg := mvutil.AdnetConfig{}
				areaCfg := mvutil.AreaConfig{}
				httpCfg := mvutil.HttpConfig{}
				tkCfg := mvutil.TrackConfig{}

				httpCfg.ConfigCenterKey = configCenterArea
				httpCfg.MKVConf = config.AerospikeConfigAs(viper.GetString(config.ConfigDirFlag), viper.GetString(config.CloudFlag), viper.GetString(config.RegionFlag))
				httpCfg.UseCtRedisConsul = false
				creativeredis.SetSubjectl(req_context.GetInstance().CreativeCacheClient.Conn)

				areaCfg.HttpConfig = httpCfg
				serviceDetail := mvutil.ServiceDetail{
					Name:      "MAdx",
					ID:        17,
					Workers:   100,
					HttpURL:   "",
					HttpsURL:  "",
					Path:      "/hbrtb",
					Method:    "POST",
					Timeout:   2000,
					UseConsul: false,
				}
				areaCfg.Service.ServiceDetail = append(areaCfg.Service.ServiceDetail, &serviceDetail)

				adnCfg.AreaConfig = &areaCfg
				adnCfg.Region = req_context.GetInstance().Region
				adnCfg.Cloud = req_context.GetInstance().Cloud
				tkCfg.TrackHost = req_context.GetInstance().Cfg.TkCfg.TrackHost
				tkCfg.PlayTrackPath = req_context.GetInstance().Cfg.TkCfg.PlayTrackPath
				adnCfg.CommonConfig = &mvutil.CommonConfig{LogConfig: mvutil.LogConfig{}, TrackConfig: tkCfg}
				// 上面各层包装，最终到这里
				mvutil.Config = &adnCfg
				mvutil.UaParser = req_context.GetInstance().UaParser
				// for pipeline filter logger
				// kv storage device tag in aerospike
				// k-v客户端
				err := mkv.InitClient()
				if err != nil {
					req_context.GetInstance().MLogs.Runtime.Errorf("init aladdin aerospike error: %s", err.Error())
					os.Exit(1)
				}

				mvutil.Logger = &mvutil.MidwayLog{
					Request:          mlog.Logger.Request,
					LossRequest:      mlog.Logger.LossRequest,
					Creative:         mlog.Logger.Creative,
					DspCreativeData:  mlog.Logger.DspCreative,
					Runtime:          mlog.Logger.Runtime,
					MappingServerLog: mlog.Logger.MappingServer,
					AerospikeLog:     mlog.Logger.DeviceAerospike,
				}
				// 添加Mintegral、MAdx、Pioneer广告后端client到manager.Backends中管理起来
				for _, service := range mvutil.Config.AreaConfig.Service.ServiceDetail {
					err := backend.BackendManager.AddBackend(service, mlog.Logger.Runtime, nil, nil)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}
			// todo 未看
			if len(cfg.TBConfigPath) > 0 {
				req_context.GetInstance().MLogs.Runtime.Info("use treasure_box_tools in server")
				extractor.InitConfig(mlog.Logger.Runtime)
				tbLogger := &tb_config.TBLogger{MaePaxLogger: mlog.Logger.TreasureBox}
				tbConfigPath := config.GetTBRegisterConfigPath(viper.GetString(config.ConfigDirFlag)) + "/"
				tb_tools.SetRegisterConfigPath(tbConfigPath)
				if err := tb_tools.Init(cfg.TBConfigPath, tbLogger, map[string]func(interface{}) error{
					"ConfigPreProc":       extractor.ConfigPreProc,
					"ConfigCenterPreProc": extractor.ConfigCenterPreProc,
				}); err != nil {
					fmt.Println("TB Init error: " + err.Error())
					os.Exit(1)
				}
				// 删除旧的 Extrator 内存加载，减少内存使用
				cfg.ExtraCfg.DbConfig = mvutil.RemoveExtractorCollections(cfg.ExtraCfg.DbConfig,
					tb_tools.GetAllRegisteredTables())
			}

			// todo
			watcher.RunWatch()

			elapsed := time.Since(start)
			req_context.GetInstance().MLogs.Runtime.Infof("service start time consuming: %s", elapsed)
		},
		Run: func(cmd *cobra.Command, args []string) {
			signChan := make(chan os.Signal, 1)
			signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
			httpConnector := hb.CreateHTTPConnector(viper.GetString(config.HTTPAddrFlag))
			httpConnector.InitRouter()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			req_context.GetInstance().MLogs.Runtime.Info("Http connector is started.")
			go httpConnector.Start(ctx)
			<-signChan
			httpConnector.Stop(context.Background())
			metrics.Stop()
		},
	}
	serveCmd.Flags().String(config.ConfigDirFlag, config.DefaultConfigDir(), "config dir path")
	serveCmd.Flags().String(config.HTTPAddrFlag, config.DefaultSrvAddr(), "http server serve addr")
	serveCmd.Flags().String(config.CloudFlag, helpers.GetCloudName(), "cloud name")
	serveCmd.Flags().String(config.RegionFlag, "singapore", "region name")
	serveCmd.Flags().String(config.ModeFlag, "test", "mode for serve options: test/dev/product")
	serveCmd.Flags().String(config.StartModeFlag, "0", "start gray mode: 0/1")
	viper.BindPFlags(serveCmd.Flags())
	return serveCmd
}
