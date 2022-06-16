package process_pipeline

import (
	"net/http"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type JssdkReqparamTransFilter struct {
}

func (jrtf *JssdkReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	r := mvutil.RequestParams{}
	r.QueryMap = mvutil.RequestQueryMap(in.URL.Query())
	RenderJssdkParam(&r)
	RenderReqParam(in, &r, r.QueryMap)
	// jssdk request type
	r.Param.RequestType = mvconst.REQUEST_TYPE_SITE
	return &r, nil
}

func RenderJssdkParam(r *mvutil.RequestParams) {
	r.QueryMap["ad_type"] = r.QueryMap["content_type"]
}
