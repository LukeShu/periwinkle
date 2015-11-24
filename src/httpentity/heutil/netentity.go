// Copyright 2015 Luke Shumaker

package heutil

import (
	"encoding/json"
	"fmt"
	"io"
)

// NetString is a string that implements httpentity.NetEntity.
type NetString string

// NetPrintf is fmt.Sprintf as a NetString.
func NetPrintf(format string, a ...interface{}) NetString {
	return NetString(fmt.Sprintf(format, a...))
}

// Encoders fulfills the httpentity.NetEntity interface.
func (s NetString) Encoders() map[string]func(out io.Writer) error {
	return map[string]func(out io.Writer) error{
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

// NetList is an array that implements httpentity.NetEntity.
type NetList []interface{}

// Encoders fulfills the httpentity.NetEntity interface.
func (l NetList) Encoders() map[string]func(out io.Writer) error {
	return map[string]func(out io.Writer) error{
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

// NetMap is a map[string]interface{} that implements
// httpentity.NetEntity.
type NetMap map[string]interface{}

// Encoders fulfills the httpentity.NetEntity interface.
func (l NetMap) Encoders() map[string]func(out io.Writer) error {
	return map[string]func(out io.Writer) error{
		"text/plain":       l.text,
		"application/json": l.json,
	}
}

func (l NetMap) text(w io.Writer) error {
	for key, val := range l {
		_, err := fmt.Fprintf(w, "%s=%s", key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l NetMap) json(w io.Writer) (err error) {
	bytes, err := json.Marshal(l)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}
