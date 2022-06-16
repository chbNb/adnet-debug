package mvutil

type BackendMetric struct {
	FilterCode   int
	IsReqBackend bool
	BackendId    int
}

// func NewBackendMetric() *BackendMetric {
// 	return &BackendMetric{FilterCode: make(map[int]int32, 0), ReqBackend: make([]string, 0)}
// }
