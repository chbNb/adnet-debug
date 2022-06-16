package output

import (
	"math/rand"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/mae/go-kit/decimal"
)

const (
	Add = "add"
	Sub = "sub"
	Mul = "mul"
	Div = "div"
)

func GetBidByImpId(seatbids []*mtgrtb.BidResponse_SeatBid, impId string) (*mtgrtb.BidResponse_SeatBid_Bid, error) {
	if len(seatbids) == 0 {
		return nil, errors.New("get BidByImpId seatbids is empty")
	}
	for _, seatbid := range seatbids {
		for _, bid := range seatbid.GetBid() {
			if bid.GetImpid() == impId {
				return bid, nil
			}
		}
	}
	return nil, errors.New("get BidByImpId has no bid match impId:" + impId)
}

func GetHBDomainPrefix(cc string) string {
	regionPrefix := req_context.GetInstance().GetRegionPrefix()
	if req_context.GetInstance().Region == req_context.Seoul && cc == "CN" {
		regionPrefix = "cn"
	}

	if req_context.GetInstance().Region == "virginia" && (req_context.GetInstance().Cloud == "aws" || req_context.GetInstance().Cloud == "aws-k8s") {
		regionPrefix = regionPrefix + "-aws"
	} else {
		regionPrefix = regionPrefix + "-new"
	}

	return regionPrefix
}

func GetDomain(cc string, multiZone bool) string {
	var domain, regionPrefix string
	rootDomain := "hb.rayjump.com"
	regionPrefix = GetHBDomainPrefix(cc)
	// 根据cc选择load域名, 目前并没有使用场景, sdk 是使用 token 按 '_' 符合分割后的后缀拼接 -hb.rayjump.com 作为 load 请求的域名
	hbLoadDomainConf := extractor.GetHBLoadDomainByCountryCodeConf()
	if hbLoadDomain, ok := hbLoadDomainConf[cc]; ok && len(hbLoadDomain) > 0 {
		rootDomain = hbLoadDomain
	}
	//
	domain = regionPrefix + "-" + rootDomain
	if multiZone {
		domain = regionPrefix + "-" + mvutil.Zone() + "-" + rootDomain
	}
	// 全链路灰度返回的灰度访问域名
	if os.Getenv("FORCE_LOAD_ENDPIONT_PREFIX_FROM_ENV") == "1" && len("LOAD_ENDPIONT_PREFIX") > 0 {
		domain = os.Getenv("LOAD_ENDPIONT_PREFIX") + "-" + rootDomain
	}
	return domain
}

func PriceConversion(rawPrice, bit float64, operation string) (float64, error) {
	dec1 := decimal.NewMDecimal()
	err := dec1.FromFloat64(rawPrice)
	if err != nil {
		return 0, err
	}

	dec2 := decimal.NewMDecimal()
	err = dec2.FromFloat64(bit)
	if err != nil {
		return 0, err
	}
	dec3 := decimal.NewMDecimal()
	switch operation {
	case Add:
		err = decimal.Add(dec1, dec2, dec3)
	case Sub:
		err = decimal.Sub(dec1, dec2, dec3)
	case Mul:
		err = decimal.Mul(dec1, dec2, dec3)
	case Div:
		err = decimal.Div(dec1, dec2, dec3, 3)
	}
	if err != nil {
		return 0, err
	}
	price, err := dec3.ToFloat64()
	if err != nil {
		return 0, err
	}
	return price, nil
}

func PriceConversionStringBit(rawPrice float64, bit, operation string) (float64, error) {
	dec1 := decimal.NewMDecimal()
	err := dec1.FromFloat64(rawPrice)
	if err != nil {
		return 0, err
	}

	dec2 := decimal.NewMDecimal()
	err = dec2.FromString([]byte(bit))
	if err != nil {
		return 0, err
	}
	dec3 := decimal.NewMDecimal()
	switch operation {
	case Add:
		err = decimal.Add(dec1, dec2, dec3)
	case Sub:
		err = decimal.Sub(dec1, dec2, dec3)
	case Mul:
		err = decimal.Mul(dec1, dec2, dec3)
	case Div:
		err = decimal.Div(dec1, dec2, dec3, 3)
	}
	if err != nil {
		return 0, err
	}
	price, err := dec3.ToFloat64()
	if err != nil {
		return 0, err
	}
	return price, nil
}

// TODO
// 此函数的逻辑在后面架构统一之后需要清理这种兼容逻辑
//
func SetBidCache(token string, val *storage.ReqCtxVal) (err error) {
	var (
		getMapValueOK bool
		aerospikeConf *mvutil.HBAerospikeConf
	)
	zone := mvutil.Zone()        // eu-central-1a or eu-central-1b
	zoneID := zone[:len(zone)-1] // eu-central-1
	// 切换到 AerospikeMultiZone 的时候，兼容发版过程不同版本的 token 写入 aerospike 集群
	if strings.Contains(token, zoneID) {
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region+"-"+zone]
	} else {
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region]
	}
	if req_context.GetInstance().Cfg.ConsulCfg.Aerospike.Enable {
		serviceName := req_context.GetInstance().Cfg.ConsulCfg.Aerospike.ServiceName
		if req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone {
			serviceName = serviceName + "-" + mvutil.Zone()
		}
		ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), serviceName)
		if ratio > 0 && ratio > rand.Float64() {
			key := &storage.ReqCtxKey{Token: token}
			client, e := req_context.GetInstance().GetMkvAerospikeClient()
			if e != nil {
				err = aerospikeClientSetValue(aerospikeConf, getMapValueOK, token, val)
			} else {
				err = client.Set(key, val)
			}
		} else {
			err = aerospikeClientSetValue(aerospikeConf, getMapValueOK, token, val)
		}
	} else {
		err = aerospikeClientSetValue(aerospikeConf, getMapValueOK, token, val)
	}
	return err
}

func aerospikeClientSetValue(conf *mvutil.HBAerospikeConf, getConfOK bool, key string, val *storage.ReqCtxVal) (err error) {
	// BidCacheClient 如果开启了 AerospikeMultiZone 则是对应 az 的 client, 否则是不区分 az 的旧集群 client
	if getConfOK {
		queryMigrate := strings.Contains(key, "migrate")
		client, e := req_context.GetInstance().GetMkvAerospikeClientWithConf(conf, queryMigrate)
		if e == nil {
			key := &storage.ReqCtxKey{Token: key}
			err = client.Set(key, val)
		} else if !queryMigrate { // 对于包含 migrate 的 token, 这里有错误兜底没用, 因为 BidCacheClient 是非 migrate 集群的 client
			err = req_context.GetInstance().BidCacheClient.SetReqCtx(key, val)
		}
	} else {
		err = req_context.GetInstance().BidCacheClient.SetReqCtx(key, val)
	}
	return err
}

// TODO
// 此函数的逻辑在后面架构统一之后需要清理这种兼容逻辑
//
func GetBidCache(token string) (val *storage.ReqCtxVal, err error) {
	var (
		getMapValueOK bool
		aerospikeConf *mvutil.HBAerospikeConf
	)
	zone := mvutil.Zone()        // eu-central-1a or eu-central-1b
	zoneID := zone[:len(zone)-1] // eu-central-1
	zone = zoneID + token[len(token)-1:]
	// 切换到 AerospikeMultiZone 的时候，兼容旧格式的 token 的 load 请求
	if strings.Contains(token, zoneID) {
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region+"-"+zone]
	} else {
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region]
	}
	if req_context.GetInstance().Cfg.ConsulCfg.Aerospike.Enable {
		serviceName := req_context.GetInstance().Cfg.ConsulCfg.Aerospike.ServiceName
		if req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone {
			serviceName = serviceName + "-" + zone
		}
		ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), serviceName)
		if ratio > 0 && ratio > rand.Float64() {
			key := &storage.ReqCtxKey{Token: token}
			client, e := req_context.GetInstance().GetMkvAerospikeClient()
			if e != nil {
				val, err = aerospikeClientQuery(aerospikeConf, getMapValueOK, token)
			} else {
				val, err = client.GetReqCtx(key)
			}
		} else {
			val, err = aerospikeClientQuery(aerospikeConf, getMapValueOK, token)
		}
	} else {
		val, err = aerospikeClientQuery(aerospikeConf, getMapValueOK, token)
	}
	return val, err
}

func aerospikeClientQuery(conf *mvutil.HBAerospikeConf, getConfOK bool, key string) (val *storage.ReqCtxVal, err error) {
	// BidCacheClient 如果开启了 AerospikeMultiZone 则是对应 az 的 client, 否则是不区分 az 的旧集群 client
	if getConfOK {
		queryMigrate := strings.Contains(key, "migrate")
		client, e := req_context.GetInstance().GetMkvAerospikeClientWithConf(conf, queryMigrate)
		if e == nil {
			k := &storage.ReqCtxKey{Token: key}
			val, err = client.GetReqCtx(k)
		} else if !queryMigrate { // 对于包含 migrate 的 token, 这里有错误兜底没用, 因为 BidCacheClient 是非 migrate 集群的 client
			val, err = req_context.GetInstance().BidCacheClient.GetReqCtx(key)
		}
	} else {
		val, err = req_context.GetInstance().BidCacheClient.GetReqCtx(key)
	}
	return val, err
}
