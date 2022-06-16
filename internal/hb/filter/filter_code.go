package filter

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type FilterCode int

// Load Error Code
const (
	QueryBidError FilterCode = 20301
	//GetBidByImpIdError       = 20302
	//TokenInvalidateError     = 20303
	UnitIDInvalidateError      FilterCode = 20304
	AdTypeInvalidateError      FilterCode = 20305
	MissingRequiredParamsError FilterCode = 20306
	// 若ironsource 的流量，根据堆栈等信息过滤
	IllegalParamError FilterCode = 20307
	//LoadDecodeDspExtError    = 20401
	//AdmEmptyError            = 20402
	//UnmarshalAsRespError     = 20403
	//DecodeAsQueryResultError = 20404
	//RenderAsCampaignError    = 20405
	//DecodeNativeError        = 20501
	//FillNativeError          = 20502
	//DecodeVastError          = 20503
	//NoVastAdsError           = 20504
	//FillVastError            = 20505
	//NoAdsInTrackingLink      = 20601
	//InjectVideoLinkError     = 20602
	//RenderRKSError           = 20701
	//MarshalJsonError         = 20702

	// bid server load 时出错
	BidServerLoadError FilterCode = 20320
)

// Bid Error Code
const (
	BidNoError FilterCode = 10001

	HttpExtractFilterParamError FilterCode = 10101
	S2SBidDataError             FilterCode = 10102
	RenderCommonDataError       FilterCode = 10103

	ReqParamFilterInputError FilterCode = 10201
	RenderBidReqDateError    FilterCode = 10202
	BidRequestUnitInValidate FilterCode = 10203
	BidRequestAppInValidate  FilterCode = 10204
	BidRequestImpEmpty       FilterCode = 10205
	BuyerUidEmpty            FilterCode = 10206
	BuyerUidDataInvalidate   FilterCode = 10207
	DecodeBidDataError       FilterCode = 10208
	BidRequestInvalidate     FilterCode = 10209
	BidMSDKVersionTooLow     FilterCode = 10210

	AreaTargetFilterInputError    FilterCode = 10301
	ClientIpInvalidate            FilterCode = 10302
	QueryNetServiceError          FilterCode = 10303
	BidInUnitCountryCodeBlackList FilterCode = 10304

	UserAgentDataFilterInputError FilterCode = 10401

	ReplaceBrandModelFilterInputError FilterCode = 10501
	ReplaceModelParamsError           FilterCode = 10502

	RenderCoreDataFilterInputError FilterCode = 10601
	AppNotFoundError               FilterCode = 10602
	PublisherNotFoundError         FilterCode = 10603
	UnitNotFoundError              FilterCode = 10604
	UnitAdNumSetNoneError          FilterCode = 10605
	IVOrientationInValidateError   FilterCode = 10606
	IVRecallNetInValidateError     FilterCode = 10607
	RenderScreenSizeInputError     FilterCode = 10608
	UnitIsNotActiveError           FilterCode = 10609
	UnitNotFoundAppError           FilterCode = 10610
	AppPlatformError               FilterCode = 10611
	AdTypeNotSupport               FilterCode = 10612
	AppDisableHb                   FilterCode = 10613
	IllegalSdkVersionForGP         FilterCode = 10614
	UnitBiddingTypeInvalidate      FilterCode = 10615
	PlacementIdInvalidateOfUnit    FilterCode = 10616
	LowFlowFilter                  FilterCode = 10617

	BuildAsRequestFilterInputError FilterCode = 10701
	ComposeAdServerRequestError    FilterCode = 10702
	ComposeAdServerJsonError       FilterCode = 10703

	TrafficSampleFilterInputError FilterCode = 10801

	BidAdxFilterInputError        FilterCode = 10901
	BidAdxComposeHttpRequestError FilterCode = 10902
	BidRequestNot200              FilterCode = 10903 // request adx after
	BidRespNoAd                   FilterCode = 10904
	BidAdxDoBidError              FilterCode = 10905
	BidAdxReadRespError           FilterCode = 10906
	BidAdxDecodeRespError         FilterCode = 10907
	DecodeDspExtError             FilterCode = 10908
	BiddingPriceError             FilterCode = 10909

	BidCacheFilterInputError FilterCode = 11001
	BidCacheError            FilterCode = 11002
	BidImpIdError            FilterCode = 11003

	FormatOutputFilterInputError FilterCode = 11101

	PriceFactorFilterInputError FilterCode = 12211

	RankerInfoFilterInputError FilterCode = 13212
)

func (code FilterCode) String() string {
	switch code {
	case HttpExtractFilterParamError:
		return "http extract filter param is not http.Request"
	case S2SBidDataError:
		return "S2S bid request data empty"
	case RenderCommonDataError:
		return "Render Common Data error"
	case ReqParamFilterInputError:
		return "req_param_filter input error"
	case RenderBidReqDateError:
		return "render bid request data error"
	case BidRequestUnitInValidate:
		return "Bid Request unit InValidate"
	case BidRequestAppInValidate:
		return "Bid Request App InValidate"
	case BidRequestImpEmpty:
		return "Bid Request Imp Is Empty"
	case BuyerUidEmpty:
		return "Buyeruid is Empty"
	case BuyerUidDataInvalidate:
		return "Buyeruid data is invalidate"
	case AreaTargetFilterInputError:
		return "area target filter input is error"
	case ClientIpInvalidate:
		return "client_ip is invalidate"
	case QueryNetServiceError:
		return "query netacuity server error"
	case UserAgentDataFilterInputError:
		return "user agent data filter input error"
	case BidMSDKVersionTooLow:
		return "Bid Request mtg sdk version too low"
	case ReplaceBrandModelFilterInputError:
		return "replace brand model filter input error"
	case ReplaceModelParamsError:
		return "replaceBrand params error"
	case RenderCoreDataFilterInputError:
		return "render core data filter input error"
	case AppNotFoundError:
		return "App not found error"
	case PublisherNotFoundError:
		return "Publisher not found"
	case UnitNotFoundError:
		return "Unit not found error"
	case UnitAdNumSetNoneError:
		return "Unit ad num set none error"
	case IVOrientationInValidateError:
		return "iv orientation invalidate error"
	case IVRecallNetInValidateError:
		return "iv recallnet invalidate error"
	case RenderScreenSizeInputError:
		return "renderScreenSize input is error"
	case UnitIsNotActiveError:
		return "unit is not active error"
	case UnitNotFoundAppError:
		return "unit not found app"
	case AppPlatformError:
		return "app platform error"
	case AdTypeNotSupport:
		return "ad type not support"
	case AppDisableHb:
		return "current app is disable header bidding"
	case IllegalSdkVersionForGP:
		return "illegal sdk version for Google play"
	case UnitBiddingTypeInvalidate:
		return "unit bidding type invalidate"
	case PlacementIdInvalidateOfUnit:
		return "placement_id invalidate of the unit_id"
	case BuildAsRequestFilterInputError:
		return "build as request filter input is invalidate"
	case ComposeAdServerRequestError:
		return "compose ad server request error"
	case BidAdxFilterInputError:
		return "bid adx filter input is error"
	case BidRequestNot200:
		return "bid request is not 200"
	case BidRespNoAd:
		return "bid response has no ad"
	case BidCacheFilterInputError:
		return "bid cache filter input error"
	case BidCacheError:
		return "bid cache error"
	case BidInUnitCountryCodeBlackList:
		return "bid in unit and country code blacklist"
	case FormatOutputFilterInputError:
		return "format output filter input error"
	case BidImpIdError:
		return "bid imp id error"
	case DecodeDspExtError:
		return "decode Dsp ext error"
	case BiddingPriceError:
		return "real time bidding price error"
	case DecodeBidDataError:
		return "decode bid data error"
	case ComposeAdServerJsonError:
		return "adserver data json error"
	case BidAdxComposeHttpRequestError:
		return "compose adx http request error"
	case BidAdxDoBidError:
		return "bid adx http do error"
	case BidAdxReadRespError:
		return "bid adx read resp error"
	case BidAdxDecodeRespError:
		return "bid adx decode resp error"
	case TrafficSampleFilterInputError:
		return "traffic sample filter input error"
	case QueryBidError:
		return "get cache miss of the load request"
	case UnitIDInvalidateError:
		return "the load request unit_id param not the same as bid the request"
	case AdTypeInvalidateError:
		return "the load request ad_type param not the same as bid the request"
	case MissingRequiredParamsError:
		return "the load request missing required params"
	case IllegalParamError:
		return "load params is illegal"
	case BidServerLoadError:
		return "load by bid server load url error"
	default:
		return "unknown"
	}
}

func (code FilterCode) Int() int {
	return int(code)
}

type MsgError struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func (a FilterCode) Message() string {
	var msg MsgError
	msg.Status = int(a)
	msg.Msg = a.String()
	rJson, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&msg)
	if err != nil {
		return "json format error"
	}
	return string(rJson)
}

func (a FilterCode) Error() string {
	return a.Message()
}

func UnmarshalMessage(message string) (code FilterCode, err error) {
	var msg MsgError
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(message), &msg)
	if err != nil {
		return code, err
	}
	code = FilterCode(msg.Status)
	return code, nil
}

func New(code FilterCode) error {
	return &PipeLineHandlerError{code}

}

type PipeLineHandlerError struct {
	code FilterCode
}

func (e *PipeLineHandlerError) Error() string {
	return e.code.Error()
}

func (e *PipeLineHandlerError) GetCode() FilterCode {
	return e.code
}

func FormatErrorMessage(errMsg string) (string, string) {
	var wrapErrMsg, rawErrMsg string
	if strings.Contains(errMsg, ": ") {
		errMsgs := strings.SplitN(errMsg, ": ", -1)
		wrapErrMsgs := errMsgs[:len(errMsgs)-1]
		wrapErrMsg = strings.Join(wrapErrMsgs, ": ")
		rawErrMsgs := errMsgs[len(errMsgs)-1:]
		rawErrMsg = rawErrMsgs[0]
	} else {
		rawErrMsg = errMsg
	}
	return wrapErrMsg, rawErrMsg
}
