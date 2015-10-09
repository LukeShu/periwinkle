// Copyright 2015 Luke Shumaker

package web

import (
	he "httpentity"
	"net/http"
	"periwinkle/cfg"
	"periwinkle/store"
	"time"
)

func Main() error {
	mux := http.NewServeMux()
	mux.Handle("/v1/", &he.Router{
		Prefix:      "/v1/",
		Root:        store.DirRoot,
		Middlewares: []he.Middleware{ postHack{}, database{}, session{} },
		Stacktrace:  cfg.Debug,
	})
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUiDir)))
	server := &http.Server{
		Addr:           cfg.WebAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	panic("not reached")
}
