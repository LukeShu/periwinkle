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
	"periwinkle/senders"
	"time"
)

func Main() error {
	std_decoders := map[string]func(io.Reader, map[string]string) (interface{}, error){
		"application/x-www-form-urlencoded": heutil.DecoderFormUrlEncoded,
		"multipart/form-data":               heutil.DecoderFormData,
		"application/json":                  heutil.DecoderJSON,
		"application/json-patch+json":       heutil.DecoderJSONPatch,
		"application/merge-patch+json":      heutil.DecoderJSONMergePatch,
	}
	std_middlewares := []he.Middleware{
		postHack{},
		database{},
		session{},
	}

	mux := http.NewServeMux()
	// The main REST API
	mux.Handle("/v1/", &he.Router{
		Prefix:      "/v1/",
		Root:        store.DirRoot,
		Decoders:    std_decoders,
		Middlewares: std_middlewares,
		Stacktrace:  cfg.Debug,
		LogRequest:  cfg.Debug,
	})
	// URL shortener service
	mux.Handle("/s/", &he.Router{
		Prefix:      "/s/",
		Root:        store.DirShortUrls,
		Decoders:    std_decoders,
		Middlewares: std_middlewares,
		Stacktrace:  cfg.Debug,
		LogRequest:  cfg.Debug,
	})

	mux.Handle("/webui/twilio/sms", http.HandlerFunc(senders.Url_handler))


	// The static web UI
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUiDir)))
	// External API callbacks
	//mux.Handle("/callbacks/MY_CALLBACK", pkg.MY_CALLBACK); // FOR FUTURE USE

	// Now actually run.
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
