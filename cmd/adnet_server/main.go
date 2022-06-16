package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	rprof "runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/chaocai2001/micro_service/microservice_helper"
	"gitlab.mobvista.com/ADN/adnet/internal/backend"
	"gitlab.mobvista.com/ADN/adnet/internal/consuls"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/geo"
	"gitlab.mobvista.com/ADN/adnet/internal/hot_data"
	"gitlab.mobvista.com/ADN/adnet/internal/mkv"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/service"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_config"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	algoab "gitlab.mobvista.com/algo-engineering/abtest-sdk-go"
	"gitlab.mobvista.com/algo-engineering/abtest-sdk-go/pkg/types"
)

var AppVersion = "unknown"

func echoStep(dFlag *int) {
	if *dFlag != 0 {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println("file:" + file + " : line=" + strconv.Itoa(line))
	}
}

func main() {
	startTimeMs := time.Now().UnixNano() / 1e6
	//  设置并发度
	CORE_NUM := runtime.NumCPU() // number of core
	runtime.GOMAXPROCS(CORE_NUM * 4)
	// debug.SetGCPercent(200
	rand.Seed(time.Now().UnixNano())

	//  处理命令行参数
	confPath := flag.String("c", "./conf/adnet_server.conf", "configure file path")
	showVersion := flag.Bool("version", false, "show version")
	dFlag := flag.Int("d", 0, "debug")
	pFile := flag.String("pfile", "/tmp/cpu.prof", "prof file")
	startMode := flag.String("mode", "", "start mode")
	flag.Parse()

	// 记录灰度标记
	process_pipeline.StartMode = *startMode
	// if *dFlag == 2 {
	// 	rprof.WriteHeapProfile(f)
	// }
	if *showVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	echoStep(dFlag)

	process_pipeline.AppVersion = AppVersion

	// 读取配置文件
	config := mvutil.Config
	if err := config.LoadConfigFile(*confPath); err != nil {
		fmt.Println("load config failed: ", err.Error())
		os.Exit(1)
	}

	// 初始化metrics,用于监控上报
	metricsRegion := config.AreaConfig.HttpConfig.Region
	if len(os.Getenv("REGION")) > 0 {
		metricsRegion = os.Getenv("REGION")
	}
	if err := metrics.InitMetrics(metricsRegion, getGitTag(AppVersion),
		config.CommonConfig.MetricsConf.MetricsConfPath, ""); err != nil {
		fmt.Println("initMetrics faild: ", err.Error())
		os.Exit(1)
	}

	//b, _ := json.Marshal(config.TBConfig)
	//fmt.Println(string(b))
	//os.Exit(0)
	// parse Version
	if !mvutil.ParseSDKVersion(config.AreaConfig.HttpConfig.SDkVersions) {
		fmt.Println("parse sdkVersion error")
		os.Exit(1)
	}
	echoStep(dFlag)
	req_context.GetInstance().SetServerIpUrl(config.AreaConfig.HttpConfig.ServerIpUrl)
	req_context.GetInstance().UpdateServerIpInfo()
	echoStep(dFlag)
	// 初始化 ua-parser
	err := mvutil.InitUaParser()
	if err != nil {
		fmt.Printf("init ua-parse error=[%s]\n", err.Error())
		os.Exit(1)
	}
	echoStep(dFlag)
	exit := make(chan os.Signal)
	// 初始化log
	mvutil.InitLogHandler(
		config.CommonConfig.LogConfig.ReqConf,
		config.CommonConfig.LogConfig.RunConf,
		config.CommonConfig.LogConfig.WatchConf,
		config.CommonConfig.LogConfig.CreativeConf,
		config.CommonConfig.LogConfig.ReduceFillConf,
		config.CommonConfig.LogConfig.LossRequestConf,
		config.CommonConfig.LogConfig.DspCreativeDataConf,
		config.CommonConfig.LogConfig.TreasureBoxConf,
		config.CommonConfig.LogConfig.ConsulAdxConf,
		config.CommonConfig.LogConfig.ConsulWatchConf,
		config.CommonConfig.LogConfig.ConsulAerospikeConf,
		config.CommonConfig.LogConfig.AerospikeConf,
		config.CommonConfig.LogConfig.MappingServerConf,
		config.CommonConfig.LogConfig.ConsulMappingServerConf,
	)
	logger := mvutil.Logger
	// 初始化backend     估计是需要与backend交互
	for _, service := range config.AreaConfig.Service.ServiceDetail {
		err := backend.BackendManager.AddBackend(service, logger.Runtime, logger.ConsulAdxLog, logger.ConsulWatchLog)
		if err != nil {
			fmt.Println("NewBackend error:", err.Error())
			os.Exit(1)
		}
	}
	echoStep(dFlag)
	watcher.Init(logger.Watch)
	logger.Runtime.Infof("ExtraConfig.ActiveDataCollecter: %#v", config.AreaConfig.ExtraConfig.ActiveDataCollecter)
	hot_data.InitActiveDataCollecter(config.AreaConfig.ExtraConfig.ActiveDataCollecter, logger.Runtime)
	echoStep(dFlag)
	// watch启动
	if err := geo.SetGeoClient(config.AreaConfig.HttpConfig.GeoConfig); err != nil {
		fmt.Printf("init ip cluster redis failed:%s\n", err.Error())
		os.Exit(1)
	}
	echoStep(dFlag)
	if config.AreaConfig.HttpConfig.UseCtRedisConsul {
		if err := redis.InitPoolFromConsul(config.AreaConfig.CtRedisConsulConfig, config.AreaConfig.RedisLocalConfig.ConnectTimeout, config.AreaConfig.RedisLocalConfig.ReadTimeout,
			config.AreaConfig.RedisLocalConfig.WriteTimeout, config.AreaConfig.RedisLocalConfig.PoolSize, logger.Runtime); err != nil {
			// fmt.Println("init local redis failed: ", err.Error())
			fmt.Printf("init creative redis failed: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		// redis local cluster
		if err := redis.InitLocalRedis(config.AreaConfig.RedisLocalConfig.HostPort, config.AreaConfig.RedisLocalConfig.ConnectTimeout, config.AreaConfig.RedisLocalConfig.ReadTimeout,
			config.AreaConfig.RedisLocalConfig.WriteTimeout, config.AreaConfig.RedisLocalConfig.PoolSize); err != nil {
			// fmt.Println("init local redis failed: ", err.Error())
			fmt.Printf("init local redis failed: %s\n", err.Error())
			os.Exit(1)
		}
		// redis local algo init
		if err := redis.InitLocalAlgoRedis(config.AreaConfig.RedisAlgoConfig.HostPort, config.AreaConfig.RedisAlgoConfig.ConnectTimeout, config.AreaConfig.RedisAlgoConfig.ReadTimeout,
			config.AreaConfig.RedisAlgoConfig.WriteTimeout, config.AreaConfig.RedisAlgoConfig.PoolSize); err != nil {
			fmt.Printf("init algo redis failed: %s\n", err.Error())
			os.Exit(1)
		}
	}
	echoStep(dFlag)
	// mkv init
	if err := mkv.InitClient(); err != nil {
		fmt.Printf("init mkv failed: %s\n", err.Error())
		os.Exit(1)
	}
	if config.AreaConfig.LBConsulConfig != nil && config.AreaConfig.LBConsulConfig.AdnetAerospike != nil && config.AreaConfig.LBConsulConfig.AdnetAerospike.Enable {
		consulResolver, err := mvutil.NewAdnetConsulResolver(config.AreaConfig.LBConsulConfig)
		if err != nil {
			fmt.Println("NewAdnetConsulResolver failure: ", err.Error())
			os.Exit(1)
		}
		consulResolver.SetLogger(logger.ConsulAerospikeLog)
		if err := consulResolver.Start(); err != nil {
			fmt.Println("AdnetConsulResolver Start failure: ", err.Error())
			os.Exit(1)
		}
		mkv.SetConsulBuild(consulResolver, config.AreaConfig.LBConsulConfig.AdnetAerospike)
	}

	// mapping server consul
	if config.AreaConfig.LBConsulConfig != nil && config.AreaConfig.LBConsulConfig.MappingServer != nil && config.AreaConfig.LBConsulConfig.MappingServer.Enable {
		mappingServerResolver, err := mvutil.NewMappingServerConsulResolver(config.AreaConfig.LBConsulConfig)
		if err != nil {
			fmt.Println("NewMappingServerConsulResolver failure: ", err.Error())
			os.Exit(1)
		}
		// 记录consul log
		mappingServerResolver.SetLogger(logger.ConsulMappingServerLog)
		if err := mappingServerResolver.Start(); err != nil {
			fmt.Println("NewMappingServerConsulResolver Start failure: ", err.Error())
			os.Exit(1)
		}
		consuls.SetMappingServerResolver(mappingServerResolver, config.AreaConfig.LBConsulConfig.MappingServer)
	}

	if len(config.AreaConfig.HttpConfig.TreasureBoxConfigPath) > 0 { //配置必须开启
		logger.Runtime.Info("use treasure_box_tools in server")
		extractor.InitConfig(logger.Runtime)                               // 给extractor增加日志记录
		tbLogger := &tb_config.TBLogger{MaePaxLogger: mvutil.Logger.TBLog} //赋值Long实例
		//设置treasure box 的配置
		if err := tb_tools.Init(config.AreaConfig.HttpConfig.TreasureBoxConfigPath, tbLogger, map[string]func(interface{}) error{
			"ConfigPreProc":       extractor.ConfigPreProc,
			"ConfigCenterPreProc": extractor.ConfigCenterPreProc,
		}); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		config.AreaConfig.ExtraConfig.DbConfig = mvutil.RemoveExtractorCollections(config.AreaConfig.ExtraConfig.DbConfig,
			tb_tools.GetAllRegisteredTables()) //需要去除的表
	}

	if config.AreaConfig.LBConsulConfig != nil {
		algoab.Init(types.Settings{
			App:                    "MTG",
			ConsulAddr:             config.AreaConfig.LBConsulConfig.Address,
			ConsulTransportTimeout: 3 * time.Second,
			Monitor:                utility.NewAlgoABTestMetric(),
		})
	}

	//if len(config.AreaConfig.IpConfig.GrpcAddress) > 0 {
	//	c, err := geoclient.NewCorsairClient(config.AreaConfig.IpConfig.GrpcAddress)
	//	if err == nil {
	//		geo.SetGeoClient(c)
	//	}
	//}

	// init bucket
	decorator, err := utility.CreateRateLimitDecorator(1*time.Second, config.AreaConfig.HttpConfig.RateLimit, config.AreaConfig.HttpConfig.MaxRateLimit)
	if err != nil {
		fmt.Printf("CreateRateLimitDecorator error: %s", err.Error())
		os.Exit(1)
	}
	process_pipeline.RateLimitDecorator = decorator
	process_pipeline.UpdateBucket(config.AreaConfig.HttpConfig.MaxRateLimit, config.AreaConfig.HttpConfig.RateLimit)
	// process_pipeline.Bucket = microservice_helper.CreateTokenBucket(config.AreaConfig.HttpConfig.MaxRateLimit, config.AreaConfig.HttpConfig.RateLimit, 1*time.Second)
	process_pipeline.RedisBucket = microservice_helper.CreateTokenBucket(config.AreaConfig.HttpConfig.IpRedisLimit, config.AreaConfig.HttpConfig.IpRedisLimit, 1*time.Second)
	echoStep(dFlag)
	//extractor.NewMDbLoaderRegistry()
	//if err := extractor.InitExtractor(config.AreaConfig.ExtraConfig, config.AreaConfig.CpMongoConsulConfig, logger.Runtime, config.AreaConfig.HttpConfig.UseMongoConsul); err != nil {
	//	fmt.Printf("init extractor failed: %s\n", err.Error())
	//	os.Exit(1)
	//}
	echoStep(dFlag)
	//extractor.DecodeConfig()
	echoStep(dFlag)
	// extractor 启动
	//extractor.Working()
	//extractor.DecodeConfigWorking()
	echoStep(dFlag)

	// watch启动
	watcher.RunWatch()
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	echoStep(dFlag)
	if *dFlag != 0 {
		f, err := os.OpenFile(*pFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()
		if *dFlag == 1 {
			rprof.StartCPUProfile(f)
			defer rprof.StopCPUProfile()
		} else if *dFlag == 2 {
			err = trace.Start(f)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer trace.Stop()
		}
	}
	echoStep(dFlag)

	echoStep(dFlag)
	logger.Runtime.Infof("ready to start server at port: %d", config.AreaConfig.HttpConfig.Port)
	httpConnector := service.CreateHTTPConnector(fmt.Sprintf(":%d", config.AreaConfig.HttpConfig.Port))
	// 注册路由
	httpConnector.InitRouter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// start server
		err := httpConnector.Start(ctx)
		for tryTime := 0; err != nil && tryTime < 300; tryTime++ {
			printTonohug(fmt.Sprintf("listen on port:%d err:%s", config.AreaConfig.HttpConfig.Port, err), -1)
			time.Sleep(1 * time.Second)
			err = httpConnector.Start(ctx)
		}
		if err != nil {
			printTonohug(fmt.Sprintf("[exit]listen on port:%d err:%s", config.AreaConfig.HttpConfig.Port, err), 2)
		}
	}()
	echoStep(dFlag)
	//启动完成，统计耗时
	printTonohug(fmt.Sprintf("Http connector is started [port:%d][use_time_ms:%d].",
		config.AreaConfig.HttpConfig.Port, time.Now().UnixNano()/1e6-startTimeMs), -1)
	<-exit // waiting for SIGINT
	//printTonohug("extractor.Stop", -1)
	//extractor.Stop()
	if tb_tools.Enable() {
		printTonohug("tb_tools.Stop", -1)
		tb_tools.Stop()
	}
	printTonohug("watcher.Stop", -1)
	watcher.Stop()
	process_pipeline.StopBucket()
	printTonohug("watcher.WriteDisk", -1)
	watcher.WriteDisk()
	printTonohug("logger.Flush", -1)
	printTonohug("rprof.StopCPUProfile", -1)
	rprof.StopCPUProfile()
	// f.Close()
	printTonohug("httpConnector.Stop", -1)
	httpConnector.Stop(context.Background())
	printTonohug("metrics.Stop", -1)
	metrics.Stop()
	printTonohug("return", -1)
	return
}

func printTonohug(log string, exitCode int) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %s\n", timeStr, log)
	if exitCode != -1 {
		os.Exit(exitCode)
	}
}

func getGitTag(version string) string {
	var tag string
	tagPos := strings.Index(version, "appversion")
	if tagPos == -1 {
		return ""
	}

	tagPos = strings.Index(version[tagPos:], ":")
	if tagPos == -1 {
		return ""
	}
	tagPos++
	tagPosEnd := strings.IndexAny(version[tagPos:], "\n[")
	if tagPosEnd == -1 {
		tag = version[tagPos:]
	} else {
		tag = version[tagPos : tagPosEnd+tagPos]
	}

	return strings.TrimSpace(tag)
}
