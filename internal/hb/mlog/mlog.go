package mlog

import "github.com/mae-pax/logger"

type MLogger struct {
	Bid                  *logger.Log
	Load                 *logger.Log
	Event                *logger.Log
	Watch                *logger.Log
	Creative             *logger.Log
	DspCreative          *logger.Log
	Runtime              *logger.Log
	Request              *logger.Log
	LossRequest          *logger.Log
	ReqMonitor           *logger.SampleLog
	TreasureBox          *logger.Log
	ConsulAerospike      *logger.Log
	ConsulAdx            *logger.Log
	ConsulWatch          *logger.Log
	ConsulAdnetAerospike *logger.Log
	MappingServer        *logger.Log
	ConsulMappingServer  *logger.Log
	DeviceAerospike      *logger.Log
}

var Logger *MLogger

func InitLogHandler(bid, load, event, watch, creative, dspCreative, runtime, request, lossRequest, reqMonitor, treasureBox,
	consulAerospike, consuleAdx, consulWatch, consulAdnetAerospike, mappingServer, consulMappingServer, deviceAerospike string) {
	ml := MLogger{}

	bidLogOptions := logger.NewFromYaml(bid)
	bidLogOptions.CloseConsoleDisplay()
	ml.Bid = bidLogOptions.InitLogger("time", "", true, true)

	loadLogOptions := logger.NewFromYaml(load)
	loadLogOptions.CloseConsoleDisplay()
	ml.Load = loadLogOptions.InitLogger("time", "", true, true)

	eventLogOptions := logger.NewFromYaml(event)
	eventLogOptions.CloseConsoleDisplay()
	ml.Event = eventLogOptions.InitLogger("time", "", true, true)

	watchLogOptions := logger.NewFromYaml(watch)
	watchLogOptions.CloseConsoleDisplay()
	ml.Watch = watchLogOptions.InitLogger("time", "", true, true)

	creativeLogOptions := logger.NewFromYaml(creative)
	creativeLogOptions.CloseConsoleDisplay()
	ml.Creative = creativeLogOptions.InitLogger("", "", false, true)

	dspCreativeLogOptions := logger.NewFromYaml(dspCreative)
	dspCreativeLogOptions.CloseConsoleDisplay()
	ml.DspCreative = dspCreativeLogOptions.InitLogger("", "", false, true)

	runtimeLogOptions := logger.NewFromYaml(runtime)
	runtimeLogOptions.CloseConsoleDisplay()
	ml.Runtime = runtimeLogOptions.InitLogger("time", "level", false, true)

	requestLogOptions := logger.NewFromYaml(request)
	requestLogOptions.CloseConsoleDisplay()
	ml.Request = requestLogOptions.InitLogger("", "", false, true)

	lossRequestLogOptions := logger.NewFromYaml(lossRequest)
	lossRequestLogOptions.CloseConsoleDisplay()
	ml.LossRequest = lossRequestLogOptions.InitLogger("", "", false, true)

	reqMonitorLogOptions := logger.NewFromYaml(reqMonitor)
	reqMonitorLogOptions.CloseConsoleDisplay()
	ml.ReqMonitor = reqMonitorLogOptions.InitSampleLogger("time", "level", false, true)

	treasureBoxLogOptions := logger.NewFromYaml(treasureBox)
	treasureBoxLogOptions.CloseConsoleDisplay()
	ml.TreasureBox = treasureBoxLogOptions.InitLogger("time", "level", false, true)

	consulAerospikeOptions := logger.NewFromYaml(consulAerospike)
	consulAerospikeOptions.CloseConsoleDisplay()
	ml.ConsulAerospike = consulAerospikeOptions.InitLogger("time", "level", false, true)

	consulAdxOptions := logger.NewFromYaml(consuleAdx)
	consulAdxOptions.CloseConsoleDisplay()
	ml.ConsulAdx = consulAdxOptions.InitLogger("time", "level", false, true)

	consulWatchLogOptions := logger.NewFromYaml(consulWatch)
	consulWatchLogOptions.CloseConsoleDisplay()
	ml.ConsulWatch = consulWatchLogOptions.InitLogger("time", "", true, true)

	consulAdnetAerospikeOptions := logger.NewFromYaml(consulAdnetAerospike)
	consulAdnetAerospikeOptions.CloseConsoleDisplay()
	ml.ConsulAdnetAerospike = consulAdnetAerospikeOptions.InitLogger("time", "level", false, true)

	mappingServerLogOptions := logger.NewFromYaml(mappingServer)
	mappingServerLogOptions.CloseConsoleDisplay()
	ml.MappingServer = mappingServerLogOptions.InitLogger("", "", false, true)

	consulMappingServerLogOptions := logger.NewFromYaml(consulMappingServer)
	consulMappingServerLogOptions.CloseConsoleDisplay()
	ml.ConsulMappingServer = consulMappingServerLogOptions.InitLogger("time", "level", false, true)

	aerospikeLogOptions := logger.NewFromYaml(deviceAerospike)
	aerospikeLogOptions.CloseConsoleDisplay()
	ml.DeviceAerospike = aerospikeLogOptions.InitLogger("time", "level", false, true)

	Logger = &ml
}
