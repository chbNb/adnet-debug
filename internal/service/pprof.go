package service

import (
	"net/http"
	"net/http/pprof"
	"path/filepath"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func (conn *HTTPConnector) InitPprof() {
	var prefix = mvutil.Config.AreaConfig.HttpConfig.PprofPath
	var names = []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"}
	var handles = map[string]http.HandlerFunc{
		"profile": pprof.Profile,
		"symbol":  pprof.Symbol,
		"trace":   pprof.Trace,
		"cmdline": pprof.Cmdline,
	}
	conn.router.HandleFunc(prefix+"/", pprof.Index)
	for _, name := range names {
		conn.router.Handle(filepath.Join(prefix, name), pprof.Handler(name))
	}
	for path, handle := range handles {
		conn.router.HandleFunc(filepath.Join(prefix, path), handle)
	}
}
