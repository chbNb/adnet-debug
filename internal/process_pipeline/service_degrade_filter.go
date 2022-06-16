package process_pipeline

import (
	"errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"math/rand"
)

var (
	ServiceDegradeFilterInputError = errors.New("service degrade filter input error")
	ServiceDegradeHitError         = errors.New("service degrade hit this request")
)

// 服务降级
type ServiceDegradeFilter struct{}

func (sdf *ServiceDegradeFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, ServiceDegradeFilterInputError
	}
	/**
	服务请求级别的降级配置格式
	{
	    "aliyun":{
	        "virginia":{
	            "/load": 0,
	        }
	    },
	    "aliyun-k8s":{
	        "virginia":{
	            "/load":0
	        }
	    },
	    "aws":{
	        "frankfurt":{
	            "/load":0
	        },
	        "seoul":{
	            "/load":0
	        },
	        "singapore":{
	            "/load":0
	        },
	        "virginia":{
	            "/load":0
	        }
	    },
	    "aws-k8s":{
	        "frankfurt":{
	            "/load":0
	        },
	        "seoul":{
	            "/load":0
	        },
	        "singapore":{
	            "/load":0
	        },
	        "virginia":{
	            "/load":0
	        }
	    }
	}
	*/
	rate := extractor.GetServiceDegradeRate(mvutil.Cloud(), mvutil.Region(), in.Param.RequestPath)
	if rate > 0 && rate > rand.Float64() {

		// 直接返回
		metrics.IncCounterWithLabelValues(30, in.Param.RequestPath)
		return in, ServiceDegradeHitError
	}
	return in, nil
}
