// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package httpentity

import (
	"encoding/json"
	"io"
	"locale"
)

//////////////////////////////////////////////////////////////////////

type EncoderJSON struct {
	Data interface{}
}

func (e EncoderJSON) Locales() []locale.Spec {
	return []locale.Spec{}
}

func (e EncoderJSON) Write(w io.Writer, l locale.Spec) locale.Error {
	bytes, err := json.Marshal(e.Data)
	if err != nil {
		return locale.UntranslatedError(err)
	}
	_, err = w.Write(bytes)
	return locale.UntranslatedError(err)
}

func (e EncoderJSON) IsText() bool {
	return true
}

//////////////////////////////////////////////////////////////////////

type EncoderJSONStr struct {
	Data locale.Stringer
}

func (e EncoderJSONStr) Locales() []locale.Spec {
	return e.Data.Locales()
}

func (e EncoderJSONStr) Write(w io.Writer, l locale.Spec) locale.Error {
	bytes, err := json.Marshal(e.Data.L10NString(l))
	if err != nil {
		return locale.UntranslatedError(err)
	}
	_, err = w.Write(bytes)
	return locale.UntranslatedError(err)
}

func (e EncoderJSONStr) IsText() bool {
	return true
}

//////////////////////////////////////////////////////////////////////

type EncoderTXT struct {
	Data locale.Stringer
}

func (e EncoderTXT) Locales() []locale.Spec {
	return e.Data.Locales()
}

func (e EncoderTXT) Write(w io.Writer, l locale.Spec) locale.Error {
	_, err := w.Write([]byte(e.Data.L10NString(l)))
	return locale.UntranslatedError(err)
}

func (e EncoderTXT) IsText() bool {
	return true
}
