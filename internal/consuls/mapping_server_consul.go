package consuls

import (
	"strconv"

	"github.com/mae-pax/consul-loadbalancer/balancer"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

var MappingServerClient struct {
	MappingServerConsulResolver *balancer.ConsulResolver
	Config                      *mvutil.MappingServer
}

func SetMappingServerResolver(r *balancer.ConsulResolver, cfg *mvutil.MappingServer) {
	MappingServerClient.MappingServerConsulResolver = r
	MappingServerClient.Config = cfg
}

func GetNode() string {
	if MappingServerClient.MappingServerConsulResolver != nil {
		node := MappingServerClient.MappingServerConsulResolver.SelectNode()
		metrics.IncCounterWithLabelValues(22, "mapping_server", mvutil.Zone(), node.Zone)
		if node == nil || len(node.Host) == 0 && MappingServerClient.Config != nil {
			return MappingServerClient.Config.HttpUrl
		} else {
			return node.Host + ":" + strconv.Itoa(node.Port)
		}
	}
	if MappingServerClient.Config != nil {
		return MappingServerClient.Config.HttpUrl
	}
	return ""
}
