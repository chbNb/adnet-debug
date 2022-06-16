package mvutil

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adx_common/constant"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"
)

func GetDspAdType(formatAdType int32) int {
	switch formatAdType {
	case mvconst.ADTypeRewardVideo:
		return constant.AD_TYPE_RV
	case mvconst.ADTypeInterstitialVideo:
		return constant.AD_TYPE_IV
	case mvconst.ADTypeNativeVideo, mvconst.ADTypeOnlineVideo:
		return constant.AD_TYPE_NV
	case mvconst.ADTypeNativePic:
		return constant.AD_TYPE_NI
	case mvconst.ADTypeSdkBanner:
		return constant.AD_TYPE_BANNER
	case mvconst.ADTypeInteractive:
		return constant.AD_TYPE_IA
	case mvconst.ADTypeSplash:
		return constant.AD_TYPE_SPLASH
	case mvconst.ADTypeNativeH5:
		return constant.AD_TYPE_NH5
	case mvconst.ADTypeInterstitialSdk:
		return constant.AD_TYPE_DI
	default:
		return 0
	}
}

func GetOpenrtbAdType(imp *openrtb.BidRequest_Imp) (int, error) {
	if imp.GetVideo() != nil && imp.GetVideo().GetExt() != nil {
		ext := imp.GetVideo().GetMvext()
		if ext.GetVideotype() == openrtb.BidRequest_Imp_Video_Ext_REWARDED_VIDEO {
			return constant.AD_TYPE_RV, nil
		} else if ext.GetVideotype() == openrtb.BidRequest_Imp_Video_Ext_INTERSTITIAL_VIDEO {
			return constant.AD_TYPE_IV, nil
		} else if ext.GetVideotype() == openrtb.BidRequest_Imp_Video_Ext_INSTREAM_VIDEO {
			return constant.AD_TYPE_ISV, nil
		}
	}
	// banner
	if imp.GetVideo() == nil && imp.GetNative() == nil && imp.GetBanner() != nil && imp.GetBanner().GetExt() != nil {
		ext := imp.GetBanner().GetMvext()
		if ext.GetAdtype() == openrtb.BidRequest_Imp_Banner_Ext_ADTYPE_SDK_BANNER {
			return constant.AD_TYPE_BANNER, nil
		} else if ext.GetAdtype() == openrtb.BidRequest_Imp_Banner_Ext_ADTYPE_SPLASH {
			return constant.AD_TYPE_SPLASH, nil
		} else if ext.GetAdtype() == openrtb.BidRequest_Imp_Banner_Ext_ADTYPE_INTERACTIVE_ADS {
			return constant.AD_TYPE_IA, nil
		} else if ext.GetAdtype() == openrtb.BidRequest_Imp_Banner_Ext_ADTYPE_INTERSTITIAL {
			return constant.AD_TYPE_DI, nil
		}
	}
	// adtype native
	if imp.GetVideo() == nil && imp.GetBanner() == nil && imp.GetNative() != nil {
		adType := constant.AD_TYPE_NI
		native := imp.GetNative()
		if native.GetMvext().GetNativetype() == openrtb.BidRequest_Imp_Native_Ext_NATIVE_H5 {
			adType = constant.AD_TYPE_NH5
		}
		var request mtgrtb.NativeRequest
		err := jsoniter.Unmarshal([]byte(native.GetRequest()), &request)
		if err != nil {
			return 0, errors.New("ssp params render decode Native.Request error:" + err.Error())
		}
		for _, as := range request.Assets {
			if as.Video != nil { // nativeType = NativeTypeVideo 或 NativeTypeDisplayVideo 都设置为NV, 竞价完成后会根据结果修正一次
				adType = constant.AD_TYPE_NV
				break
			}
		}
		return adType, nil
	}

	return 0, errors.New("adType error, unknown adType")
}

// adx adtype -> adnet adtype
func GetAdnetAdType(adtype int) int {
	switch adtype {
	case constant.AD_TYPE_RV:
		return mvconst.ADTypeRewardVideo
	case constant.AD_TYPE_NV:
		return mvconst.ADTypeNativeVideo
	case constant.AD_TYPE_IV:
		return mvconst.ADTypeInterstitialVideo
	case constant.AD_TYPE_NI:
		return mvconst.ADTypeNative
	case constant.AD_TYPE_IA:
		return mvconst.ADTypeInteractive
	case constant.AD_TYPE_BANNER:
		return mvconst.ADTypeSdkBanner
	case constant.AD_TYPE_ISV:
		return 0
	case constant.AD_TYPE_SPLASH:
		return mvconst.ADTypeSplash
	case constant.AD_TYPE_DI:
		return mvconst.ADTypeInterstitial
	case constant.AD_TYPE_NH5:
		return mvconst.ADTypeNativeH5
	}
	return 0
}
