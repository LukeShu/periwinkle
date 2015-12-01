// Copyright 2015 Luke Shumaker

package httpentity

import (
	"encoding/json"
	"fmt"
	"io"
	"locale"
	"net/http"
	"net/url"
)

func ErrorToNetEntity(status int16, err locale.Error) NetEntity {
	return NetPrintf("%d %s %v", status, http.StatusText(int(status)), err)
}

//////////////////////////////////////////////////////////////////////

// NetString is a string that implements httpentity.NetEntity.
type NetStringer struct {
	locale.Stringer
}

var _ NetEntity = NetStringer{}

// NetPrintf is fmt.Sprintf as a NetString.
func NetPrintf(format string, a ...interface{}) NetStringer {
	return NetStringer{locale.Sprintf(format, a...)}
}

// Encoders fulfills the httpentity.NetEntity interface.
func (s NetStringer) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain":       EncoderTXT{s},
		"application/json": EncoderJSONStr{s},
	}
}

//////////////////////////////////////////////////////////////////////

type NetLocations []*url.URL

type netLocationsTXT struct{ NetLocations }
type netLocationsJSON struct{ NetLocations }

// Encoders fulfills the httpentity.NetEntity interface.
func (l NetLocations) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain":       netLocationsTXT{l},
		"application/json": netLocationsJSON{l},
	}
}

func (l NetLocations) Locales() []locale.Spec {
	return []locale.Spec{}
}

func (l NetLocations) IsText() bool {
	return true
}

func (l netLocationsTXT) Write(w io.Writer, loc locale.Spec) locale.Error {
	for _, line := range l.NetLocations {
		_, uerr := fmt.Fprintln(w, line)
		if uerr != nil {
			return locale.UntranslatedError(uerr)
		}
	}
	return nil
}

func (l netLocationsJSON) Write(w io.Writer, loc locale.Spec) locale.Error {
	strs := make([]string, len(l.NetLocations))
	for i, u := range l.NetLocations {
		strs[i] = u.String()
	}
	bytes, uerr := json.Marshal(strs)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	_, uerr = w.Write(bytes)
	return locale.UntranslatedError(uerr)
}

//////////////////////////////////////////////////////////////////////

type NetJSON struct {
	Data interface{}
}

// Encoders fulfills the httpentity.NetEntity interface.
func (l NetJSON) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"application/json": l,
	}
}

func (l NetJSON) Locales() []locale.Spec {
	return []locale.Spec{}
}

func (l NetJSON) IsText() bool {
	return true
}

func (l NetJSON) Write(w io.Writer, loc locale.Spec) locale.Error {
	bytes, uerr := json.Marshal(l.Data)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	_, uerr = w.Write(bytes)
	return locale.UntranslatedError(uerr)
}
