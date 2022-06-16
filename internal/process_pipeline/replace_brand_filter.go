package process_pipeline

import (
	"errors"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type ReplaceBrandFilter struct {
}

// 品牌、型号规范化   品牌型号这些数据需要传递给算法做分析的，所以为防止SDK传递过来的数据不规范，需要进行规范化
func (rbf *ReplaceBrandFilter) Process(data interface{}) (interface{}, error) {

	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("IpInfoFilter input type should be *params.RequestParams")
	}

	brand := strings.ToLower(in.Param.Brand)
	model := strings.ToLower(in.Param.Model)

	if in.PublisherInfo.PublisherId == mvconst.PUB_XIAOMI {
		brand = "xiaomi"
	}

	//var rep map[string]map[string]string
	rep, _ := extractor.GetREPLACE_BRAND_MODEL()
	filter := false
	if repBrand, ok := rep["brand"]; !ok || len(repBrand) == 0 {
		filter = true
	}
	if !filter {
		if repModel, ok := rep["model"]; !ok || len(repModel) == 0 {
			filter = true
		}
	}
	if filter {
		mvutil.Logger.Runtime.Warnf("replaceBrandFilter get replaceDevice error")
		in.Param.ReplaceBrand = brand
		in.Param.ReplaceModel = model
		in.Param.ExtBrand = in.Param.ReplaceBrand
		in.Param.ExtModel = in.Param.ReplaceModel
		return in, nil
	}

	repBrand := rep["brand"]
	repModel := rep["model"]

	if in.Param.Platform == mvutil.IOSPLATFORM && len(model) > 0 {
		if reModel, ok := repModel[model]; ok {
			model = reModel
		}
	}

	if in.Param.Platform == mvutil.ANDROIDPLATFORM {
		if newBrand, newModel, ok := trimModel(brand, model, rep["model_trim"]); ok {
			brand = newBrand
			model = newModel
		}
		if len(brand) > 0 && brand != "0" {
			new := brand

			if strings.Contains(new, ".") {
				new = strings.Replace(new, ".", "", -1)
			}
			if reBrand, ok := repBrand[new]; ok {
				new = reBrand
			}

			reModel, err := replaceBrand(brand, new, model)
			if err != nil {
				mvutil.Logger.Runtime.Warnf("request_id=[%s] replaceBrand old=[%s] new=[%s] model=[%s] error:%s", in.Param.RequestID, brand, new, model, err)
			} else {
				brand = new
				model = reModel
			}
		}
	}

	in.Param.ReplaceBrand = brand
	in.Param.ReplaceModel = model
	in.Param.ExtBrand = in.Param.ReplaceBrand
	in.Param.ExtModel = in.Param.ReplaceModel

	return in, nil
}

// replaceBrand 将型号前面的品牌前缀去掉
func replaceBrand(old, new, model string) (string, error) {
	if len(old) == 0 || len(new) == 0 || len(model) == 0 {
		return model, errors.New("replaceBrand params error")
	}

	model = strings.Replace(model, old, new, -1)
	for {
		has := false
		for _, linker := range linkers {
			brandStr := new + linker
			if strings.HasPrefix(model, brandStr) {
				model = strings.Replace(model, brandStr, "", 1) // brand+连接符替换成""
				has = true
				break
			}
		}
		if !has {
			break
		}
	}
	return model, nil
}

var linkers = []string{" ", "-", "_", ""} // 空字符串一定要放在最后位置

func trimModel(brand, model string, trimModelMap map[string]string) (newBrand string, newModel string, ok bool) {
	if len(brand) != 0 && brand != "0" || len(model) == 0 || len(trimModelMap) == 0 {
		return
	}
	for baseBrand := range trimModelMap {
		has := false
		brandStr := baseBrand + " "
		for {
			if strings.HasPrefix(model, brandStr) {
				model = strings.Replace(model, brandStr, "", 1) // brand+连接符替换成""
				has = true
			} else {
				break
			}
		}
		if has {
			return baseBrand, model, true
		}
	}
	return
}
