package clients

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mae-pax/consul-loadbalancer/balancer"
	"github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type MAdxClient struct {
	logger         *logger.Log
	backendCfg     *mvutil.ServiceDetail
	consulResolver *balancer.ConsulResolver
	//Client         *utility.HttpClient
}

const (
	keepalive int = 200
	idle      int = 2000
	dns       int = 30
)

func init() {
	utility.InitFastHttpClient(keepalive, keepalive, idle)
}

func NewMAdxClient(cfg *mvutil.ServiceDetail, logger *logger.Log, consulAdxLog *logger.Log, consulWatchLog *logger.Log) (*MAdxClient, error) {
	adxClient := &MAdxClient{
		logger:     logger,
		backendCfg: cfg,
	}
	if cfg.UseConsul {
		if cfg.ConsulCfg == nil {
			return nil, errors.New("consul config is nil")
		}
		consulResolver, err := balancer.NewConsulResolver(
			cfg.ConsulCfg.Cloud,
			cfg.ConsulCfg.Address,
			cfg.ConsulCfg.Service,
			"clb/adx/cpu_threshold.json",
			"clb/adx/zone_cpu.json",
			"clb/adx/instance_factor.json",
			"clb/adx/onlinelab_factor.json",
			time.Duration(cfg.ConsulCfg.Internal)*time.Second,
			time.Duration(cfg.ConsulCfg.Timeout)*time.Millisecond,
			"clb/adx/services.json",
		)

		if err != nil {
			logger.Errorf("InitAdxConsul error:" + err.Error())
			return nil, err
		}
		if consulAdxLog != nil {
			consulResolver.SetLogger(consulAdxLog)
		}
		if consulWatchLog != nil {
			consulResolver.SetWatcher(consulWatchLog)
		}
		err = consulResolver.Start()
		if err != nil {
			logger.Errorf("Start AdxConsul node update error:" + err.Error())
			return nil, err
		}
		adxClient.consulResolver = consulResolver
	}

	// adxClient.Client = utility.NewHttpClient(cfg.Timeout, keepalive, idle, dns)
	return adxClient, nil
}

func (mc *MAdxClient) GetNode() string {
	if os.Getenv("FORCE_ADX_ENDPIONT_FROM_ENV") == "1" && len(os.Getenv("ADX_SERVICE_NAME")) > 0 {
		return os.Getenv("ADX_SERVICE_NAME")
	}
	if mc.backendCfg.UseConsul && mc.consulResolver != nil {
		ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), "adx")
		if ratio > 0 && ratio > rand.Float64() {
			node := mc.consulResolver.SelectNode()
			if node == nil || len(node.Host) <= 0 {
				mc.logger.Warnf("not find node by consul, use default address")
				return mc.backendCfg.HttpURL
			} else {
				metrics.IncCounterWithLabelValues(22, "adx", mvutil.Zone(), node.Zone)
				var addr = node.Host
				if strings.Index(addr, ":") == -1 && node.Port != 0 {
					addr += ":" + strconv.Itoa(node.Port)
				}
				return "http://" + addr + "/open_rtb"
			}
		}
	}
	return mc.backendCfg.HttpURL
}

func (mc *MAdxClient) Stop() {
	if mc.backendCfg.UseConsul && mc.consulResolver != nil {
		mc.consulResolver.Stop()
	}
}
