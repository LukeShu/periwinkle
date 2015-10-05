// Copyright 2015 Luke Shumaker

package httpentity

import (
	"io"
	"encoding/json"
)

type netString string

func (s netString) Encoders() map[string]Encoder {
	return map[string]Encoder{"text/plain": s.write}
}

func (s netString) write(w io.Writer) error {
	_, err := w.Write([]byte(s))
	return err
}

type netJson struct{
	body interface{}
}

func (j netJson) Encoders() map[string]Encoder {
	return map[string]Encoder{"application/json": j.write}
}

func (j netJson) write(w io.Writer) (err error) {
	bytes, err := json.Marshal(j.body)
	if err != nil {
		return
	}
	_, err = w.Write(bytes)
	return
}
