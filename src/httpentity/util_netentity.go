// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"encoding/json"
	"fmt"
)

type netString string

func (s netString) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain": s.text,
		"application/json": s.json,
	}
}

func (s netString) text(w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

func (s netString) json(w io.Writer) (err error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}

type netList []interface{}

func (l netList) Encoders() map[string]Encoder {
	return map[string]Encoder{
		"text/plain": l.text,
		"application/json": l.json,
	}
}

func (l netList) text(w io.Writer) error {
	for _, line := range l {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l netList) json(w io.Writer) (err error) {
	bytes, err := json.Marshal(l)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}
