// Copyright 2015 Luke Shumaker

package httpentity

import (
	"encoding/json"
	"fmt"
	"io"
)

type NetString string

func (s NetString) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain":       s.text,
		"application/json": s.json,
	}
}

func (s NetString) text(w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

func (s NetString) json(w io.Writer) (err error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}

type NetList []interface{}

func (l NetList) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain":       l.text,
		"application/json": l.json,
	}
}

func (l NetList) text(w io.Writer) error {
	for _, line := range l {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l NetList) json(w io.Writer) (err error) {
	bytes, err := json.Marshal(l)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}
