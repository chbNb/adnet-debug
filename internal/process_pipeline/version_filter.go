package process_pipeline

import (
	"errors"
	"net/http"
)

type VersionFilter struct {
}

var AppVersion string

func (qmdf *VersionFilter) Process(data interface{}) (interface{}, error) {
	_, ok := data.(*http.Request)
	if !ok {
		return nil, errors.New("VersionFilter input type should be *http.Request")
	}

	return &AppVersion, nil
}
