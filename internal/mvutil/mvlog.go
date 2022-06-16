package mvutil

import "github.com/mae-pax/logger"

var Logger *MidwayLog

type MidwayLog struct {
	Watch                  *logger.Log
	Creative               *logger.Log
	Request                *logger.Log
	LossRequest            *logger.Log
	Runtime                *logger.Log
	ReduceFill             *logger.Log
	DspCreativeData        *logger.Log
	TBLog                  *logger.Log
	ConsulAdxLog           *logger.Log
	ConsulWatchLog         *logger.Log
	ConsulAerospikeLog     *logger.Log
	AerospikeLog           *logger.Log
	MappingServerLog       *logger.Log
	ConsulMappingServerLog *logger.Log
}

func InitLogHandler(requestLog, runtimeLog, watchLog, creativeLog, reduceFillLog,
	lossRequestLog, dspCreativeDataLog, tbLog, consulAdxLog, consulWatchLog, consulAerospike, aerospikeLog, mappingServerLog, consulMappingServerLog string) {
	ml := MidwayLog{}

	requestLogOptions := logger.NewFromYaml(requestLog)
	requestLogOptions.CloseConsoleDisplay()
	ml.Request = requestLogOptions.InitLogger("", "", false, true)

	runtimeLogOptions := logger.NewFromYaml(runtimeLog)
	runtimeLogOptions.CloseConsoleDisplay()
	ml.Runtime = runtimeLogOptions.InitLogger("time", "level", false, true)

	watchLogOptions := logger.NewFromYaml(watchLog)
	watchLogOptions.CloseConsoleDisplay()
	ml.Watch = watchLogOptions.InitLogger("time", "", true, true)

	creativeLogOptions := logger.NewFromYaml(creativeLog)
	creativeLogOptions.CloseConsoleDisplay()
	ml.Creative = creativeLogOptions.InitLogger("", "", false, true)

	reduceFillLogOptions := logger.NewFromYaml(reduceFillLog)
	reduceFillLogOptions.CloseConsoleDisplay()
	ml.ReduceFill = reduceFillLogOptions.InitLogger("", "", false, true)

	lossRequestLogOptions := logger.NewFromYaml(lossRequestLog)
	lossRequestLogOptions.CloseConsoleDisplay()
	ml.LossRequest = lossRequestLogOptions.InitLogger("", "", false, true)

	dspCreativeDataLogOptions := logger.NewFromYaml(dspCreativeDataLog)
	dspCreativeDataLogOptions.CloseConsoleDisplay()
	ml.DspCreativeData = dspCreativeDataLogOptions.InitLogger("", "", false, true)

	tbLogOptions := logger.NewFromYaml(tbLog)
	tbLogOptions.CloseConsoleDisplay()
	ml.TBLog = tbLogOptions.InitLogger("time", "level", false, true)

	consulAdxLogOptions := logger.NewFromYaml(consulAdxLog)
	consulAdxLogOptions.CloseConsoleDisplay()
	ml.ConsulAdxLog = consulAdxLogOptions.InitLogger("time", "level", false, true)

	consulWatchLogOptions := logger.NewFromYaml(consulWatchLog)
	consulWatchLogOptions.CloseConsoleDisplay()
	ml.ConsulWatchLog = consulWatchLogOptions.InitLogger("time", "", true, true)

	consulAerospikeLogOptions := logger.NewFromYaml(consulAerospike)
	consulAerospikeLogOptions.CloseConsoleDisplay()
	ml.ConsulAerospikeLog = consulAerospikeLogOptions.InitLogger("time", "level", false, true)

	aerospikeLogOptions := logger.NewFromYaml(aerospikeLog)
	aerospikeLogOptions.CloseConsoleDisplay()
	ml.AerospikeLog = aerospikeLogOptions.InitLogger("time", "level", false, true)

	mappingServerLogOptions := logger.NewFromYaml(mappingServerLog)
	mappingServerLogOptions.CloseConsoleDisplay()
	ml.MappingServerLog = mappingServerLogOptions.InitLogger("", "", false, true)

	consulMappingServerLogOptions := logger.NewFromYaml(consulMappingServerLog)
	consulMappingServerLogOptions.CloseConsoleDisplay()
	ml.ConsulMappingServerLog = consulMappingServerLogOptions.InitLogger("time", "level", false, true)
	Logger = &ml
}
