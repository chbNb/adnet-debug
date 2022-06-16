package geo

import "gitlab.mobvista.com/mae/geo-lib/pkg/client"

func SetGeoClient(cfg string) error {
	return client.NewGeoLibClient(cfg)
}

func GetGeo(ip string) (client.Geo, bool, error) {
	return client.GetIpGeoData(ip)
}
