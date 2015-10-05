// Copyright 2015 Luke Shumaker

package httpentity

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/url"
	"net/http"
)

func (r Response) WriteEntity(w io.Writer) error {
	if r.entity == nil {
		return nil
	}
	mimetype, _, _ := mime.ParseMediaType(r.Headers.Get("Content-Type"))
	encoders := r.entity.Encoders()
	return encoders[mimetype](w)
}

func ReadEntity(r io.Reader, contenttype string) (interface{}, error) {
	mimetype, params, err := mime.ParseMediaType(contenttype)
	if err != nil {
		return nil, err
	}
	switch mimetype {
	case "application/x-www-form-urlencoded":
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		entity, err := url.ParseQuery(string(bytes))
		if err != nil {
			return nil, err
		}
		return entity, nil
	case "multipart/form-data":
		boundary, ok := params["boundary"]
		if !ok {
			return nil, http.ErrMissingBoundary
		}
		reader := multipart.NewReader(r, boundary)
		form, err := reader.ReadForm(0)
		if err != nil {
			return nil, err
		}
		entity := make(map[string]interface{}, len(form.Value)+len(form.File))
		for k, v := range form.Value {
			entity[k] = v
		}
		for k, v := range form.File {
			if _, exists := entity[k]; exists {
				values := entity[k].([]string)
				list := make([]interface{}, len(values)+len(v))
				i := uint(0)
				for _, value := range values {
					list[i] = value
					i++
				}
				for _, value := range v {
					list[i] = value
					i++
				}
				entity[k] = list
			} else {
				entity[k] = v
			}
		}
		return entity, nil
	case "application/json":
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		var entity interface{}
		err = json.Unmarshal(bytes, &entity)
		if err != nil {
			return nil, err
		}
		return entity, nil
	case "application/octet-stream":
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
	return nil, nil
}
