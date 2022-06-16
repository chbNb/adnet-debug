package clients

import (
	"errors"
	"os"
	"time"

	"github.com/easierway/go-kit/balancer"
	"github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type PioneerClient struct {
	logger         *logger.Log
	backendCfg     *mvutil.ServiceDetail
	consulResolver *balancer.ConsulResolver
}

func NewPioneerClient(cfg *mvutil.ServiceDetail, logger *logger.Log) (*PioneerClient, error) {
	pc := &PioneerClient{
		logger:     logger,
		backendCfg: cfg,
	}
	if cfg.UseConsul {
		if cfg.ConsulCfg == nil {
			return nil, errors.New("consul config is nil")
		}
		Resolver, err := balancer.NewConsulResolver(
			cfg.ConsulCfg.Address,
			cfg.ConsulCfg.Service,
			cfg.ConsulCfg.MyService,
			time.Duration(cfg.ConsulCfg.Internal)*time.Microsecond,
			cfg.ConsulCfg.ServiceRatio,
			cfg.ConsulCfg.CpuThreshold,
		)

		if err != nil {
			logger.Errorf("InitPioneerConsul error:" + err.Error())
			return nil, err
		}
		Resolver.SetLogger(logger)
		pc.consulResolver = Resolver
	}
	return pc, nil
}

func (pc *PioneerClient) GetNode() string {
	if os.Getenv("FORCE_PIONEER_ENDPIONT_FROM_ENV") == "1" && len(os.Getenv("PIONEER_SERVICE_NAME")) > 0 {
		return os.Getenv("PIONEER_SERVICE_NAME")
	}
	if pc.backendCfg.UseConsul && pc.consulResolver != nil {
		node := pc.consulResolver.DiscoverNode()
		if node == nil || len(node.Address) <= 0 {
			pc.logger.Warnf("[pioneer]not find node by consul, use default address: %s", pc.backendCfg.HttpURL)
			return pc.backendCfg.HttpURL
		} else {
			return "http://" + node.Address
		}
	} else {
		return pc.backendCfg.HttpURL
	}
}
