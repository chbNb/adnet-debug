package process_pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type ImageFilter struct {
}

type ImageResult struct {
	Status int       `json:"status"`
	Msg    string    `json:"msg"`
	Data   ImageDate `json:"data"`
}

type ImageDate struct {
	Image string `json:"image"`
	Rp    bool   `json:"rp"`
	Rpt   int    `json:"rpt"`
}

func (imgf *ImageFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("ParamRenderFilter input type should be *params.RequestParams")
	}

	unitId := in.Param.UnitID
	sign := in.Param.Sign

	if unitId <= 0 || len(sign) == 0 {
		return nil, errors.New("imageFilter unitId sign can not be empty")
	}

	rawSign := fmt.Sprintf("%d%s", in.Param.AppID, in.PublisherInfo.Publisher.Apikey)
	newSign := mvutil.Md5(rawSign)
	if newSign != sign && needCheckSign(in) {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by checkSign is validate appID=[%d],apikey=[%s],sign=[%s], param.sign=[%s]",
			in.Param.RequestID, in.Param.AppID, in.PublisherInfo.Publisher.Apikey, newSign, in.Param.Sign)
		return nil, errorcode.EXCEPTION_SIGN_ERROR
		//return mvconst.EXCEPTION_SIGN_ERROR, fmt.Errorf("EXCEPTION_SIGN_ERROR")
	}

	if int(in.Param.AdType) != int(ad_server.ADType_APPWALL) {
		return nil, errors.New("imageFilter unit not appwall")
	}

	//var data map[string]string
	var rs ImageDate
	rs.Image = in.UnitInfo.Unit.EntraImage
	if in.UnitInfo.Unit.RedPointShow == nil {
		rs.Rp = true
	} else {
		rs.Rp = *in.UnitInfo.Unit.RedPointShow
	}

	if in.UnitInfo.Unit.RedPointShowInterval == nil {
		rs.Rpt = 3600
	} else {
		rs.Rpt = *in.UnitInfo.Unit.RedPointShowInterval
	}

	var mr ImageResult
	mr.Status = mvconst.MESSAGE_SUCCESS
	mr.Msg = "success"
	mr.Data = rs

	result, err := json.Marshal(&mr)

	str := strings.TrimRight(string(result), "\n")
	if err != nil {
		return "", errors.New("imageFilter marshal failed")
	}
	return &str, nil

}
