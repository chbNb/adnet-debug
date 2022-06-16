package process_pipeline

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
)

type MAXReqparamTransFilter struct {
}

func (mrtf *MAXReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()

	r := mvutil.RequestParams{}
	var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderMAXReqParam(body, &r)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("request params error:[%s]", err)
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderMAXReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var maxParams protobuf.BidRequest
	err := proto.Unmarshal(body, &maxParams)
	if err != nil {
		return nil, err
	}
	// 获取Adslot内容
	var deals []*protobuf.BidRequest_AdSlot_Deal
	var native *protobuf.BidRequest_AdSlot_Native
	var nativeTem []*protobuf.BidRequest_AdSlot_Native_NativeTemplate
	var minPrice *uint32
	var adslotId *uint32
	//var unitWidth *int32
	//var unitHeight *int32
	for _, v := range maxParams.Adslot {
		deals = v.Deals
		native = v.Native
		minPrice = v.MinCpmPrice
		adslotId = v.Id
		//unitWidth = v.Width
		//unitHeight = v.Height
		nativeTem = v.Native.NativeTemplate
		break
	}
	// 价格比较，若比设置值低则过滤
	onlinePriceFloor, ifFind := extractor.GetONLINE_PRICE_FLOOR()
	var resPrice int32
	if ifFind {
		price, ok := onlinePriceFloor["max360"]
		if ok {
			resPrice = int32(10000 * price)
		} else {
			resPrice = int32(200000)
		}
	}
	if minPrice == nil {
		return nil, errors.New("minPrice error")
	}
	if int32(*minPrice) > resPrice {
		pkg := ""
		if maxParams.Mobile.PackageName != nil {
			pkg = *maxParams.Mobile.PackageName
		}
		mvutil.Logger.Runtime.Warnf("pkg=[%s], price=[%d]", pkg, *minPrice)
		return nil, errors.New("minPrice is too high")
	}
	// 获取deals内容
	var dealId *int64
	for _, v := range deals {
		dealId = v.DealId
		break
	}
	query := make(mvutil.RequestQueryMap)
	if maxParams.UserAgent != nil {
		query["useragent"] = mvutil.RenderQueryMapV(*maxParams.UserAgent)
	}
	if maxParams.Ip != nil {
		query["client_ip"] = mvutil.RenderQueryMapV(*maxParams.Ip)
	}
	if maxParams.Mobile.Device.Model != nil {
		query["model"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.Model)
	}
	if maxParams.Mobile.Device.Os == nil {
		return nil, errors.New("platform empty")
	}
	if *maxParams.Mobile.Device.Os == mvconst.PlatformNameAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		if maxParams.Mobile.Device.AndroidId != nil {
			query["android_id"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.AndroidId)
		}
	} else if *maxParams.Mobile.Device.Os == mvconst.PlatformNameIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		if maxParams.Mobile.Device.Idfa != nil {
			query["idfa"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.Idfa)
		}
		query["http_req"] = mvutil.RenderQueryMapV("2")
	}
	if maxParams.Mobile.Device.OsVersion != nil {
		query["os_version"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.OsVersion)
	}
	if maxParams.Mobile.Device.Imei != nil {
		query["imei"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.Imei)
	}
	if maxParams.Mobile.Device.Mac != nil {
		query["mac"] = mvutil.RenderQueryMapV(*maxParams.Mobile.Device.Mac)
	}
	if maxParams.Mobile.PackageName != nil {
		query["package_name"] = mvutil.RenderQueryMapV(*maxParams.Mobile.PackageName)
	}
	// network映射值
	if maxParams.Mobile.Device.Network != nil {
		if *maxParams.Mobile.Device.Network == 1 {
			query["network_type"] = mvutil.RenderQueryMapV("9")
		} else if *maxParams.Mobile.Device.Network == 5 {
			query["network_type"] = mvutil.RenderQueryMapV("0")
		} else {
			network := strconv.FormatUint(uint64(*maxParams.Mobile.Device.Network), 10)
			query["network_type"] = mvutil.RenderQueryMapV(network)
		}
	}
	if maxParams.Mobile.Device.ScreenOrientation != nil {
		orientation := protobuf.BidRequest_Mobile_Device_ScreenOrientation_value[maxParams.Mobile.Device.ScreenOrientation.String()]
		query["orientation"] = mvutil.RenderQueryMapV(strconv.FormatInt(int64(orientation), 10))
	}
	if maxParams.Mobile.Device.ScreenWidth != nil && maxParams.Mobile.Device.ScreenHeight != nil {
		screenWidthStr := strconv.Itoa(int(*maxParams.Mobile.Device.ScreenWidth))
		screenHeightStr := strconv.Itoa(int(*maxParams.Mobile.Device.ScreenHeight))
		query["screen_size"] = mvutil.RenderQueryMapV(screenWidthStr + "x" + screenHeightStr)
	}
	//if unitWidth != nil && unitHeight != nil {
	//	width := strconv.Itoa(int(*unitWidth))
	//	height := strconv.Itoa(int(*unitHeight))
	//	query["unit_size"] = mvutil.RenderQueryMapV(width + "x" + height)
	//}

	if native.AdNum != nil {
		adNum := strconv.FormatUint(uint64(*native.AdNum), 10)
		query["ad_num"] = mvutil.RenderQueryMapV(adNum)
	}
	r.Param.MaxBid = *maxParams.Bid
	r.Param.AdslotId = *adslotId
	// 获取底价
	r.Param.MaxPrice = resPrice
	if dealId != nil {
		r.Param.MaxDealId = *dealId
	}
	// 获取templateid
	if len(nativeTem) > 0 {
		for _, v := range nativeTem {
			if v.TemplateId != nil && *v.TemplateId == 9 {
				r.Param.TemplateId = *v.TemplateId
			}
			for _, val := range v.ImageSize {
				r.Param.MaxImgW = *(val.Width)
				r.Param.MaxImgH = *(val.Height)
				break
			}
		}
	}
	if r.Param.TemplateId != 9 {
		return nil, errors.New("not allow template type")
	}
	MaxImgWStr := strconv.Itoa(int(r.Param.MaxImgW))
	MaxImgHStr := strconv.Itoa(int(r.Param.MaxImgH))
	query["unit_size"] = mvutil.RenderQueryMapV(MaxImgWStr + "x" + MaxImgHStr)
	AppAndUnitConf, ifFind := extractor.GetMAX_APPID_AND_UNITID()
	if ifFind {
		conf, ok := AppAndUnitConf[*maxParams.Mobile.Device.Os]
		if ok {
			query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(conf.AppId, 10))
			query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(conf.UnitId, 10))
		} else {
			return nil, errors.New("can not get appid & unitid")
		}
	}
	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")
	return query, nil
}
