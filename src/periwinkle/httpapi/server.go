// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package httpapi

import (
	he "httpentity"
	"httpentity/heutil"
	"httpentity/rfc7231"
	"io"
	"locale"
	"net"
	"net/http"
	"periwinkle"
	"periwinkle/domain_handlers"
	"stoppable"
)

func MakeServer(socket net.Listener, cfg *periwinkle.Cfg) *stoppable.HTTPServer {
	stdDecoders := map[string]func(io.Reader, map[string]string) (interface{}, locale.Error){
		"application/x-www-form-urlencoded": heutil.DecoderFormURLEncoded,
		"multipart/form-data":               heutil.DecoderFormData,
		"application/json":                  heutil.DecoderJSON,
		"application/json-patch+json":       heutil.DecoderJSONPatch,
		"application/merge-patch+json":      heutil.DecoderJSONMergePatch,
	}
	stdMiddlewares := []he.Middleware{
		MiddlewarePostHack,
		MiddlewareDatabase(cfg),
		MiddlewareSession,
	}
	mux := http.NewServeMux()
	// The main REST API
	mux.Handle("/v1/", he.Router{
		Prefix:           "/v1/",
		Root:             NewDirRoot(),
		Decoders:         stdDecoders,
		Middlewares:      stdMiddlewares,
		Stacktrace:       cfg.Debug,
		Log:              heutil.StderrLog,
		TrustForwarded:   cfg.TrustForwarded,
		MethodNotAllowed: rfc7231.StatusMethodNotAllowed,
	}.Init())
	// URL shortener service
	mux.Handle("/s/", he.Router{
		Prefix:           "/s/",
		Root:             NewDirShortURLs(),
		Decoders:         stdDecoders,
		Middlewares:      stdMiddlewares,
		Stacktrace:       cfg.Debug,
		Log:              heutil.StderrLog,
		TrustForwarded:   cfg.TrustForwarded,
		MethodNotAllowed: rfc7231.StatusMethodNotAllowed,
	}.Init())

	// The static web UI
	mux.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(cfg.WebUIDir)))

	// External API callbacks
	mux.Handle("/callbacks/twilio-sms", http.HandlerFunc(domain_handlers.TwilioSMSCallbackServer{DB: cfg.DB}.ServeHTTP))

	// Make the server
	return &stoppable.HTTPServer{
		Server: http.Server{Handler: mux},
		Socket: socket,
	}
}
