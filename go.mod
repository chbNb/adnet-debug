module gitlab.mobvista.com/ADN/adnet

go 1.12

require (
	bou.ke/monkey v1.0.2
	git.apache.org/thrift.git v0.0.0-20151001171628-53dd39833a08
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/aerospike/aerospike-client-go v2.2.0+incompatible // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/aws/aws-sdk-go v1.20.3 // indirect
	github.com/bouk/monkey v1.0.1
	github.com/chaocai2001/micro_service v0.0.0-20180831105051-26edbc2148c5
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/easierway/concurrent_map v0.0.0-20190103024436-7073b0dd7e95
	github.com/easierway/go-kit v1.0.5
	github.com/easierway/pipefiter_framework v0.0.0-20180731055939-3e4803bd9b8e
	github.com/easierway/service_decorators v1.0.1 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/flatbuffers v1.12.0
	github.com/gorilla/mux v1.7.4
	github.com/json-iterator/go v1.1.12
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/mae-pax/consul-loadbalancer v0.1.9
	github.com/mae-pax/logger v0.3.2
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/ua-parser/uap-go v0.0.0-20210121150957-347a3497cc39
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/valyala/fasthttp v1.17.0
	github.com/xdg/stringprep v1.0.0 // indirect
	gitlab.mobvista.com/ADN/adx_common v1.6.5
	gitlab.mobvista.com/ADN/chasm v0.15.12
	gitlab.mobvista.com/ADN/exporter v0.1.6
	gitlab.mobvista.com/ADN/lego v1.2.6
	gitlab.mobvista.com/ADN/mtg_openrtb v0.23.10
	gitlab.mobvista.com/ADN/structs v0.10.6
	gitlab.mobvista.com/ADN/treasure_box_sdk v1.1.13
	gitlab.mobvista.com/adserver/recommend_protocols v1.10.10
	gitlab.mobvista.com/algo-engineering/abtest-sdk-go v1.0.9
	gitlab.mobvista.com/mae/geo-lib v0.1.9
	gitlab.mobvista.com/mae/go-kit v0.1.13
	gitlab.mobvista.com/mtech/mkv v1.9.4
	gitlab.mobvista.com/voyager/clickmode v0.1.5
	gitlab.mobvista.com/voyager/common v0.6.20
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	google.golang.org/grpc v1.42.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace google.golang.org/grpc v1.42.0 => google.golang.org/grpc v1.29.1
