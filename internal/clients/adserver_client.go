package clients

import (
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/easierway/go-kit/balancer"
	"github.com/mae-pax/logger"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type AdServerClient struct {
	logger         *logger.Log
	backendCfg     *mvutil.ServiceDetail
	consulResolver *balancer.ConsulResolver
}

func NewAdServerClient(cfg *mvutil.ServiceDetail, logger *logger.Log) (*AdServerClient, error) {
	ac := &AdServerClient{
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
			logger.Errorf("InitAsConsul error:" + err.Error())
			return nil, err
		}
		Resolver.SetLogger(logger)
		ac.consulResolver = Resolver
	}
	return ac, nil
}

func (ac *AdServerClient) GetCampaigns(req *ad_server.QueryParam, timeout int) (*ad_server.QueryResult_, error) {
	node := ac.GetNode()
	if len(node) <= 0 {
		return nil, errors.New("no as backend")
	}

	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	// node跟consul有关，估计是从配置中心获得rpc服务器的地址
	transport, err := thrift.NewTSocketTimeout(node, time.Millisecond*time.Duration(timeout))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("connect[%s] as timeout[%d]", node, timeout))
	}
	defer transport.Close()
	useTransport := transportFactory.GetTransport(transport)

	if err := useTransport.Open(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("AdServerNode[%s], open transport error", node))
	}
	defer useTransport.Close()

	client := ad_server.NewRecommendSrvClientFactory(useTransport, protocolFactory)
	// 开始rpc调用？
	res, err := client.GetCampaigns(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("AdServerNode[%s], GetCampaigns error", node))
	}

	return res, nil
}

func (ac *AdServerClient) GetNode() string {
	if ac.backendCfg.UseConsul && ac.consulResolver != nil {
		node := ac.consulResolver.DiscoverNode()
		if node == nil || len(node.Address) <= 0 {
			ac.logger.Warnf("[adserver]not find node by consul, use default address: %s", ac.backendCfg.HttpURL)
			return ac.backendCfg.HttpURL
		} else {
			return node.Address
		}
	} else {
		return ac.backendCfg.HttpURL
	}
}
