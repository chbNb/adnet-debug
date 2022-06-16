package consuls

import (
	"time"

	"github.com/easierway/go-kit/balancer"
	mlogger "github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

var CreativeRedisResolver *balancer.ConsulResolver

func InitCreativeRedisConsul(consulConfig mvutil.ConsulConfig, log *mlogger.Log) error {
	Resolver, err := balancer.NewConsulResolver(consulConfig.Address, consulConfig.Service, consulConfig.MyService,
		time.Duration(consulConfig.Internal)*time.Millisecond, consulConfig.ServiceRatio, consulConfig.CpuThreshold)
	if err != nil {
		log.Error("InitCreativeRedisConsul error:" + err.Error())
		return err
	}
	Resolver.SetLogger(log)
	CreativeRedisResolver = Resolver
	return nil
}
