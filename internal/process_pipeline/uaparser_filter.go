package process_pipeline

import (
	"errors"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type UaParserFilter struct {
}

func getPlatformFromFamily(family string) int {
	if strings.ToLower(family) == "android" {
		return mvconst.PlatformAndroid
	} else if strings.ToLower(family) == "ios" {
		return mvconst.PlatformIOS
	} else {
		return mvconst.PlatformOther
	}
}

// 根据ua解析出OSversion、platform和model
func (upf *UaParserFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("UaParserFilter input type should be *params.RequestParams")
	}

	if (in.Param.Platform == 0 || len(in.Param.OSVersion) == 0) && len(in.Param.UserAgent) > 0 {
		//info := mvutil.UaParser.Parse(in.Param.UserAgent)
		if len(in.Param.OSVersion) == 0 {
			os := mvutil.UaParser.ParseOs(in.Param.UserAgent)
			in.Param.OSVersion = strings.ToLower(os.ToVersionString())
		}
		//从ua中解析platform
		if in.Param.Platform == 0 {
			pl := mvutil.UaParser.ParseOs(in.Param.UserAgent)
			in.Param.Platform = getPlatformFromFamily(pl.Family)
		}
		//从ua中解析model
		if in.Param.Model == "0" {
			model := mvutil.UaParser.ParseDevice(in.Param.UserAgent)
			if len(model.Family) > 0 {
				in.Param.Model = model.Family
			}
		}
	}
	if len(in.Param.OSVersion) == 0 {
		return nil, errorcode.EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED
	}

	return in, nil
}
