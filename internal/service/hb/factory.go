package hb

const (
	Bid          = "bid"
	Load         = "load"
	Event        = "event"
	Health       = "health"
	QueryMem     = "query_mem"
	ReloadOne    = "reload_one"
	ReloadAll    = "reload_all"
	HotTable     = "hot_table"
	HotDataAll   = "hot_data_all"
	HotDataRange = "hot_data_range"
)

// var NotFoundPipelineError = errors.New("can't find the pipeline definition")
//
// type HBHandler struct {
// 	Pipeline pf.Filter
// }
//
// func CreateHandler(name string) (http.Handler, error) {
// 	switch name {
// 	case Bid:
// 		return CreateBidHandler(), nil
// 	case Load:
// 		return CreateLoadHandler(), nil
// 	case Win, Loss:
// 		return CreateEventHandler(), nil
// 	}
// 	return nil, NotFoundPipelineError
// }

// func CreateService(name string) *WallTimePipeline {
// 	switch name {
// 	case Bid:
// 		filters := []pf.Filter{&base.HttpExtractFilter{}}
// 		return &WallTimePipeline{
// 			Name:        Bid,
// 			TimeElapsed: make([]AtomicInt, len(filters)),
// 		}
// 	case Load:
// 		filters := []pf.Filter{&base.HttpExtractFilter{}}
// 		return &WallTimePipeline{
// 			Name:        Load,
// 			TimeElapsed: make([]AtomicInt, len(filters)),
// 		}
// 	case Win:
// 		filters := []pf.Filter{&base.HttpExtractFilter{}}
// 		return &WallTimePipeline{
// 			Name:        Win,
// 			TimeElapsed: make([]AtomicInt, len(filters)),
// 		}
// 	case Loss:
// 		filters := []pf.Filter{&base.HttpExtractFilter{}}
// 		return &WallTimePipeline{
// 			Name:        Loss,
// 			TimeElapsed: make([]AtomicInt, len(filters)),
// 		}
// 	default:
// 		return &WallTimePipeline{}
// 	}
// }
