// Copyright 2015 Luke Shumaker

package web

import (
	"fmt"
	he "httpentity"
	"httpentity/util"
	"io"
	"net/http"
	"periwinkle/cfg"
	"periwinkle/store"
	"time"
)

func Main() error {
	mux := http.NewServeMux()
	mux.Handle("/v1/", &he.Router{
		Prefix: "/v1/",
		Root:   store.DirRoot,
		Decoders: map[string]func(io.Reader, map[string]string) (interface{}, error){
			"application/x-www-form-urlencoded": heutil.DecoderFormUrlEncoded,
			"multipart/form-data":               heutil.DecoderFormData,
			"application/json":                  heutil.DecoderJSON,
			"application/json-patch+json":       heutil.DecoderJSONPatch,
			"application/merge-patch+json":      heutil.DecoderJSONMergePatch,
		},
		Middlewares: []he.Middleware{postHack{}, database{}, session{}},
		Stacktrace:  cfg.Debug,
		LogRequest:  cfg.Debug,
	})
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUiDir)))
	server := &http.Server{
		Addr:           cfg.WebAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	panic(fmt.Sprintf("Could not start HTTP server: %v", err))
}
