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
	mux.Handle("/v1/", he.NetHttpHandler(cfg.Debug, "/v1/", store.DirRoot, postHack{}, database{}, session{}))
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
