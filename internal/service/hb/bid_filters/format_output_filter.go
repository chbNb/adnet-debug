package bid_filters

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output/c2s"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output/s2s"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	adn_output "gitlab.mobvista.com/ADN/adnet/internal/output"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

type FormatOutputFilter struct {
}

func (fof *FormatOutputFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, filter.FormatOutputFilterInputError
	}

	fof.renderToken(in)
	// 是否需要Aerospik Gzip压缩
	gzipRate := extractor.GetHBAerospikeGzipRate(mvutil.Cloud(), mvutil.Region())
	if gzipRate > 0 && gzipRate > rand.Float64() {
		// 开启
		in.ReqParams.Param.ExtDataInit.AerospikeGzipEnable = 1
	}

	// 是否移除 AppInfo, UnitInfo, PublisherInfo
	removeRedundancyRate := extractor.GetHBAerospikeRemoveRedundancyRate(mvutil.Cloud(), mvutil.Region())
	if removeRedundancyRate > 0 && removeRedundancyRate > rand.Float64() {
		// 开启
		in.ReqParams.Param.ExtDataInit.AerospikeRemoveRedundancyEnable = 1
	}

	// 记录需要记录到hb request日志的字段
	adn_output.RenderHBRequestLogParam(in.ReqParams)

	// 记录了 loss request 就不记录 request
	if !in.ReqParams.Param.LossReqFlag {
		mvutil.StatHBRequestLog(in)
	}

	if in.Result == nil || in.ReqParams.Nbr != constant.OK {
		return in, errors.New(strings.Join(in.ReqParams.Param.BackendReject, ";"))
	}

	var (
		respByte      []byte
		err           error
		getMapValueOK bool
		aerospikeConf *mvutil.HBAerospikeConf
	)

	in.ReqParams.Param.HBExtPfData = in.ReqParams.Param.ExtData
	tokenStrs := strings.Split(in.ReqParams.Token, "_") // fb6a79a7-440a-4593-a344-ff8c2421b18a-1637839753_fk-new
	tokenStrs1 := ""
	if len(tokenStrs) >= 2 {
		tokenStrs1 = tokenStrs[1]
	}
	randVal := rand.Intn(100)
	// 修改 token 格式
	if req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone {
		tokenZoneTag := mvutil.Zone()
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region+"-"+tokenZoneTag]
		if getMapValueOK {
			if aerospikeConf.MigrateEnable && randVal < aerospikeConf.MigrateRate {
				// fb6a79a7-440a-4593-a344-ff8c2421b18a-1637839753-migrate_fk-new-eu-central-1a
				// fk-new-eu-central-1a-hb.rayjump.com
				in.ReqParams.Token = tokenStrs[0] + "-migrate_" + tokenStrs1 + "-" + tokenZoneTag
			} else {
				// fb6a79a7-440a-4593-a344-ff8c2421b18a-1637839753_fk-new-eu-central-1a
				// fk-new-eu-central-1a-hb.rayjump.com
				in.ReqParams.Token = tokenStrs[0] + "_" + tokenStrs1 + "-" + tokenZoneTag
			}
		} else {
			// fb6a79a7-440a-4593-a344-ff8c2421b18a-1637839753_fk-new-eu-central-1a
			// fk-new-eu-central-1a-hb.rayjump.com
			in.ReqParams.Token = tokenStrs[0] + "_" + tokenStrs1 + "-" + tokenZoneTag
		}
	} else {
		aerospikeConf, getMapValueOK = extractor.GetHBAerospikeConf().ConfMap[req_context.GetInstance().Cloud+"-"+req_context.GetInstance().Region]
		if getMapValueOK {
			if aerospikeConf.MigrateEnable && randVal < aerospikeConf.MigrateRate {
				// fb6a79a7-440a-4593-a344-ff8c2421b18a-1637839753-migrate_fk-new
				// fk-new-hb.rayjump.com
				in.ReqParams.Token = tokenStrs[0] + "-migrate_" + tokenStrs1
			}
		}
	}

	if len(in.ReqParams.Param.HBS2SBidID) > 0 {
		resp, _ := s2s.RenderOutput(in)
		respByte, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resp)
		if err != nil {
			return in, errors.Wrap(err, "json marshal")
		}
	} else {
		resp, _ := c2s.RenderOutput(in)
		respByte, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resp)
		if err != nil {
			return in, errors.Wrap(err, "json marshal")
		}
	}

	// 放到这一步才写入，为了减少不必要的缓存
	// 清空不必要的数据
	in.ReqParams.PostData = make([]byte, 0)
	in.ReqParams.Param.DeviceInstalledPackages = make(map[string][]byte, 0)
	in.ReqParams.FreqDataFromAerospike = make(map[string][]byte, 0)
	in.ReqParams.Param.ExcludePackageNames = make(map[string]bool, 0)
	//in.ReqParams.UnitInfo.AdSourceCountry = make(map[string]int, 0)
	in.ReqParams.AppInfo.DefaultSkIds = ""
	in.ReqParams.AppInfo.AppSkIds = make([]smodel.AppSkId, 0)
	// 清除 AppInfo, UnitInfo, PublishInfo. load 时分别再从 $req.param.AppID, $req.param.UnitID, $req.param.PublisherID 渲染回去
	// 全局变量控制, 防止发版滚动更新出现新旧pod并行导致解析bug
	if in.ReqParams.Param.ExtDataInit.AerospikeRemoveRedundancyEnable == 1 {
		in.ReqParams.AppInfo = nil
		in.ReqParams.UnitInfo = nil
		in.ReqParams.PublisherInfo = nil
	}
	for _, v := range in.Backends {
		v.RespData = make([]byte, 0)
	}
	reqCtx := &storage.ReqCtxVal{
		ReqParams:       in.ReqParams,
		Result:          in.Result,
		MaxTimeout:      in.MaxTimeout,
		FlowTagID:       in.FlowTagID,
		RandValue:       in.RandValue,
		AdsTest:         in.AdsTest,
		Elapsed:         in.Elapsed,
		Finished:        in.Finished,
		OrderedBackends: in.OrderedBackends,
		DebugModeInfo:   in.DebugModeInfo,
		IsWhiteGdt:      in.IsWhiteGdt,
		IsNativeVideo:   in.IsNativeVideo,
	}
	// reqCtxByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(reqCtx)
	// req_context.GetInstance().MLogs.ReqMonitor.Debugf("AdType: %s, Token: %s, ReqCtxVal: %s, ReqCtxVal Size: %d", mvutil.GetAdTypeStr(in.ReqParams.Param.AdType), in.ReqParams.Token, reqCtxByte, len(reqCtxByte))
	err = output.SetBidCache(in.ReqParams.Token, reqCtx)
	// err = req_context.GetInstance().BidCacheClient.SetReqCtx(in.ReqParams.Token, reqCtx)
	if err != nil {
		// if strings.Contains(err.Error(), "Record too big") {
		// 	b, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(reqCtx)
		// 	req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidCacheClient SetReqCtx failure: %s, ReqCtxVal size: %d, error: %s", b, len(b), err.Error())
		// }

		watcher.AddWatchValue(constant.SetBidCacheNotSuccess, 1)
		metrics.IncCounterWithLabelValues(6)
		return in, errors.Wrap(filter.New(filter.BidCacheError), err.Error())
	}

	req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(in, 0, ""))

	return &respByte, nil
}

func (fof *FormatOutputFilter) renderToken(in *mvutil.ReqCtx) {
	regionPrefix := output.GetHBDomainPrefix(in.ReqParams.Param.CountryCode)
	if in.ReqParams != nil && in.ReqParams.Param.ExtDataInit.LoadCDNTag == 1 && extractor.GetLoadDomainABTest() != nil {
		regionPrefix += "-" + extractor.GetLoadDomainABTest().CDNPrefix
	}
	// 全链路灰度返回的 load 访问域名前缀
	if os.Getenv("FORCE_LOAD_ENDPIONT_PREFIX_FROM_ENV") == "1" && len(os.Getenv("LOAD_ENDPIONT_PREFIX")) > 1 {
		regionPrefix = os.Getenv("LOAD_ENDPIONT_PREFIX")
	}
	in.ReqParams.Token = in.ReqParams.Token + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + regionPrefix
}
