// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package web

import (
	he "httpentity"
	"httpentity/util"
	"io"
	"net"
	"net/http"
	"periwinkle/cfg"
	"periwinkle/senders"
	"periwinkle/store"
	"stoppable"
)

var server stoppable.HTTPServer

func Main(socket net.Listener) error {
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

	// The static web UI
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUiDir)))

	// External API callbacks
	mux.Handle("/callbacks/twilio-sms", http.HandlerFunc(senders.Url_handler))

	// Now actually run.
	server = stoppable.HTTPServer{
		Server: http.Server{Handler: mux},
		Socket: socket,
	}
	return server.Serve()
}

func Stop() {
	server.Stop()
}
