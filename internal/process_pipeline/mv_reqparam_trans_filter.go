package process_pipeline

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MVReqparamTransFilter struct {
}

const (
	OpenRTBV23      = "2.3"
	OpenRTBV25      = "2.5"
	HBS2SHeaderFlag = "openrtb"
)

// 从req提取出ip、useragent以及body等参数，封装到RequestParams中
func (mrtf *MVReqparamTransFilter) Process(data interface{}) (res interface{}, err error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		// return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}
	// 将query和post的请求参数放到rawQuery这个map中
	rawQuery := mvutil.RequestQueryMap(in.URL.Query())
	if in.Method == "POST" { // ADD post data
		in.ParseForm()
		for k, v := range in.PostForm {
			rawQuery[k] = v
		}
	}
	// in.body放到r中
	r := mvutil.RequestParams{}
	// 根据请求路径判断是否是RTB，设置到header中
	if in.URL.Path == mvconst.PATHMopubBid {
		in.Header.Set(HBS2SHeaderFlag, OpenRTBV23)
	}
	r.Header = in.Header
	// TODO refactor hb and adnet
	// s2s,c2s，跟第三方监控平台有关。c2s是由客户端直接发送到第三方服务的；s2s是先发送到开发者服务端，再从服务端发送到第三方服务端
	if in.Header.Get(HBS2SHeaderFlag) == OpenRTBV23 || in.Header.Get(HBS2SHeaderFlag) == OpenRTBV25 {
		body, err := ioutil.ReadAll(in.Body)
		if err != nil {
			return nil, errors.Wrap(err, filter.S2SBidDataError.Error())
		}
		r.PostData = body
	} else if len(rawQuery) <= 0 {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		// return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}
	RenderReqParam(in, &r, rawQuery)

	return &r, nil
}

func RenderReqParam(in *http.Request, r *mvutil.RequestParams, rawQuery mvutil.RequestQueryMap) {
	r.Method = in.Method
	r.Param.RequestPath = in.URL.Path
	// TODO refactor hb and adnet
	// 判断头部竞价
	if r.Param.RequestPath == mvconst.PATHBidAds || r.Param.RequestPath == mvconst.PATHBid || r.Param.RequestPath == mvconst.PATHMopubBid || r.Param.RequestPath == mvconst.PATHLoad {
		r.IsHBRequest = true
	}
	// 记录请求url
	r.Param.RequestURI = in.RequestURI
	// 解析请求参数
	if astMode := in.URL.Query().Get("ast_mode"); astMode != "" {
		rawQuery["ast_mode"] = in.URL.Query()["ast_mode"]
	}
	// IP和ua,请求参数有就从参数取，没有就从header中取
	clientIp, err := rawQuery.GetString("client_ip", true, "")
	r.Param.ParamCIP = clientIp
	r.Param.ClientIP = clientIp
	if len(r.Param.ClientIP) == 0 || err != nil {
		// 请求没有传client_ip，使用http头部解析的ip
		r.Param.ClientIP = GetClientIP(in)
	} else {
		ip := net.ParseIP(r.Param.ClientIP)
		if ip == nil {
			// 请求client_ip字段为非法ip
			r.Param.ClientIP = GetClientIP(in)
		}
	}

	userAgent, err := rawQuery.GetString("useragent", true, "")
	if err != nil || len(userAgent) <= 0 {
		userAgent = in.Header.Get("User-Agent")
	}
	r.Param.UserAgent = userAgent
	r.Param.Extra9 = r.Param.UserAgent
	r.Param.ExtsystemUseragent = mvutil.RawUrlEncode(in.Header.Get("User-Agent"))
	r.Param.MvLine = in.Header.Get("Mv-Line")

	r.QueryMap = rawQuery
}
