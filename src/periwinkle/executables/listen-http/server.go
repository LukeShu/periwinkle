// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package main

import (
	he "httpentity"
	"httpentity/util"
	"io"
	"log"
	"net"
	"net/http"
	"periwinkle"
	"periwinkle/email_handlers"
	"periwinkle/store"
	"stoppable"
)

func makeServer(socket net.Listener, cfg *periwinkle.Cfg) *stoppable.HTTPServer {
	std_decoders := map[string]func(io.Reader, map[string]string) (interface{}, error){
		"application/x-www-form-urlencoded": heutil.DecoderFormUrlEncoded,
		"multipart/form-data":               heutil.DecoderFormData,
		"application/json":                  heutil.DecoderJSON,
		"application/json-patch+json":       heutil.DecoderJSONPatch,
		"application/merge-patch+json":      heutil.DecoderJSONMergePatch,
	}
	std_middlewares := []he.Middleware{
		MiddlewarePostHack,
		MiddlewareDatabase(cfg),
		MiddlewareSession,
	}
	mux := http.NewServeMux()
	// The main REST API
	mux.Handle("/v1/", he.Router{
		Prefix:         "/v1/",
		Root:           store.DirRoot,
		Decoders:       std_decoders,
		Middlewares:    std_middlewares,
		Stacktrace:     cfg.Debug,
		LogRequest:     cfg.Debug,
		TrustForwarded: cfg.TrustForwarded,
	}.Init())
	// URL shortener service
	mux.Handle("/s/", he.Router{
		Prefix:         "/s/",
		Root:           store.DirShortUrls,
		Decoders:       std_decoders,
		Middlewares:    std_middlewares,
		Stacktrace:     cfg.Debug,
		LogRequest:     cfg.Debug,
		TrustForwarded: cfg.TrustForwarded,
	}.Init())

	// The static web UI
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUiDir)))

	smsCallbackServer := handlers.SmsCallbackServer{}
	go func() {
		err := smsCallbackServer.Serve()
		if err != nil {
			log.Printf("Could not serve SmsCallbackServer: %v\n", err)
		}
	}()

	// External API callbacks
	mux.Handle("/callbacks/twilio-sms", http.HandlerFunc(smsCallbackServer.ServeHTTP))

	// Make the server
	return &stoppable.HTTPServer{
		Server: http.Server{Handler: mux},
		Socket: socket,
	}
}
