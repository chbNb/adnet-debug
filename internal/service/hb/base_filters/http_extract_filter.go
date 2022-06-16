package base_filters

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
)

type HttpExtractFilter struct {
}

const (
	S2SHeaderFlag = "openrtb"
)

// 提取httpReq中的参数，如path、host等，关键是提取clientIp
func (hef *HttpExtractFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, filter.HttpExtractFilterParamError
	}
	httpReqData := &params.HttpReqData{}
	httpReqData.Path = in.URL.Path
	httpReqData.Host = in.Host
	httpReqData.QueryData = params.HttpQueryMap(in.URL.Query())
	headerFlag := in.Header.Get(S2SHeaderFlag)
	if len(headerFlag) > 0 && headerFlag == "2.5" {
		body, err := ioutil.ReadAll(in.Body)
		if err != nil {
			return nil, errors.Wrap(err, filter.S2SBidDataError.Error())
		}
		httpReqData.PostData = body
	}

	err := hef.renderCommonData(in, httpReqData)
	if err != nil {
		return nil, err
	}

	return httpReqData, nil
}

func (hef *HttpExtractFilter) renderCommonData(req *http.Request, data *params.HttpReqData) error {
	if req == nil || data == nil {
		return filter.RenderCommonDataError
	}
	clientIp := data.QueryData.GetString("client_ip", true, "")
	if len(clientIp) == 0 {
		clientIp = hef.getClientIp(req)
	} else {
		if !helpers.IsCorrectIp(clientIp) {
			clientIp = hef.getClientIp(req)
		}
	}
	data.ClientIp = clientIp
	// userAgent := data.QueryData.GetString("useragent", true, "")
	// if len(userAgent) == 0 {
	// userAgent =
	// }
	data.UserAgent = req.Header.Get("User-Agent")
	return nil
}

func (hef *HttpExtractFilter) requestHeader(req *http.Request, key string) string {
	if values, ok := req.Header[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (hef *HttpExtractFilter) getClientIp(req *http.Request) string {
	forwardedByClientIP := true
	if forwardedByClientIP {
		clientIP := strings.TrimSpace(hef.requestHeader(req, "X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}
		clientIP = hef.requestHeader(req, "X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
