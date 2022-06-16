package process_pipeline

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	//"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type OpenApiV3Filter struct {
}

type ErrorInfo struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func (oaf *OpenApiV3Filter) Process(data interface{}) (interface{}, error) {
	return "open_api_v3", nil
}

func genAdBackendData(id int, campaignId string, used int) string {
	var databuf bytes.Buffer
	databuf.WriteString(strconv.Itoa(id))
	databuf.WriteString(":")
	databuf.WriteString(campaignId)
	// databuf.WriteString(strconv.FormatInt(ad.CampaignId, 10))
	databuf.WriteString(":")
	databuf.WriteString(strconv.Itoa(used))
	databuf.WriteString(";")
	return databuf.String()
}

func GetBackendConfig(reqCtx *mvutil.ReqCtx) string {
	// 获取adBackend使用distributionInfo.adBackend，如果默认就使用mobvista
	adBackendConfig := "1:0:0"

	textConfig := []string{}
	for _, backendId := range reqCtx.OrderedBackends {
		backendCtx, ok := reqCtx.Backends[backendId]
		if !ok {
			continue
		}

		if backendId == mvconst.Mobvista || backendId == mvconst.Pioneer {

			if len(backendCtx.AdReqKeyName) == 0 {
				// 有点恶心
				textConfig = append(textConfig, "1:0:0")
				continue
			}
		}

		textConfig = append(textConfig,
			fmt.Sprintf("%d:%s:%d", backendId, backendCtx.AdReqKeyName, backendCtx.Content))
	}

	if len(textConfig) > 0 {
		adBackendConfig = strings.Join(textConfig, ";")
	}
	return adBackendConfig
}

func FilterAds(reqCtx *mvutil.ReqCtx, id int, ads *corsair_proto.BackendAds, adCount *int64, backendData *[]string) {
	// mvutil.Logger.Runtime.Infof("======== FilterAds ====== %#v", ads)
	var filtered_ads []*corsair_proto.Campaign
	for _, ad := range ads.GetCampaignList() {

		*backendData = append(*backendData, genAdBackendData(id, ad.CampaignId, 0))

		if reqCtx.ReqParams.Param.TNum > 0 && *adCount >= int64(reqCtx.ReqParams.Param.TNum) {
			// 大模板不做限制
			if reqCtx.ReqParams.Param.AdType != mvconst.ADTypeInteractive && !reqCtx.ReqParams.Param.BigTemplateFlag {
				continue
			}
		}

		if dspExt, _ := reqCtx.ReqParams.GetDspExt(); id != mvconst.Mobvista && !reqCtx.ReqParams.IsFakeAs && !(dspExt != nil && dspExt.DspId == mvconst.MAS) { // 不走as和mas
			// mobvista ads has no creatives, adnnet will fill creatives
			if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeNative {
				if len(reqCtx.ReqParams.Param.VideoVersion) > 0 && reqCtx.ReqParams.Param.VideoAdType != mvconst.VideoAdTypeNOLimit {
					if (reqCtx.ReqParams.Param.VideoAdType == mvconst.VideoAdTypeNOVideo && len(ad.GetVideoURL()) > 0) ||
						(reqCtx.ReqParams.Param.VideoAdType == mvconst.VideoAdTypeOnlyVideo && len(ad.GetVideoURL()) == 0) {
						continue
					}

				}
				if len(reqCtx.ReqParams.Param.VideoVersion) == 0 {
					// 没有传video_version的native广告，默认不支持视频
					if len(ad.GetVideoURL()) > 0 {
						continue
					}
				}
			} else if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeRewardVideo ||
				reqCtx.ReqParams.Param.AdType == mvconst.ADTypeInterstitialVideo {
				if len(ad.GetVideoURL()) == 0 && len(ads.GetBannerHtml()) == 0 && len(ads.GetBannerUrl()) == 0 &&
					len(ad.GetUrlTemplate()) == 0 && len(ad.GetHtmlTemplate()) == 0 {
					// 对banner的判断是支持只返回mraid的情况
					continue
				}
			}

		}

		filtered_ads = append(filtered_ads, ad)
		(*backendData)[len(*backendData)-1] = genAdBackendData(id, ad.CampaignId, 1)
		*adCount = *adCount + 1
	}

	ads.CampaignList = filtered_ads
}
