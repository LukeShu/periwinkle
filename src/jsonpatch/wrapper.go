// Copyright 2015 Luke Shumaker

package jsonpatch

import (
	"encoding/json"
	evan "github.com/evanphx/json-patch"
	matt "github.com/mattbaird/jsonpatch"
)

type Patch interface {
	Apply(in interface{}, out interface{}) error
}

type jsonPatch evan.Patch
var _ Patch = jsonPatch{}

func NewJSONPatch(str []byte) (Patch, error) {
	patch, err := evan.DecodePatch(str)
	return jsonPatch(patch), err
}

func (patch jsonPatch) Apply(in interface{}, out interface{}) error {
	in_bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	out_bytes, err := (evan.Patch(patch)).Apply(in_bytes)
	if err != nil {
		return err
	}
	return json.Unmarshal(out_bytes, out)
}

type jsonMergePatch json.RawMessage
var _ Patch = jsonMergePatch{}

func NewJSONMergePatch(str []byte) (Patch, error) {
	return jsonMergePatch(str), nil
}

func (patch jsonMergePatch) Apply(in interface{}, out interface{}) error {
	in_bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	out_bytes, err := evan.MergePatch(in_bytes, patch)
	if err != nil {
		return err
	}
	return json.Unmarshal(out_bytes, out)
}

func DiffPatch(a interface{}, b interface{}) (Patch, error) {
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	p_bytes, err := evan.CreateMergePatch(a_bytes, b_bytes)
	if err != nil {
		return nil, err
	}
	return jsonMergePatch(p_bytes), nil
}

func DiffMergePatch(a interface{}, b interface{}) (Patch, error) {
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	p, err := matt.CreatePatch(a_bytes, b_bytes)
	p_bytes, err := json.Marshal(p)
	return jsonMergePatch(p_bytes), err
}

func Equal(a interface{}, b interface{}) bool {
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return evan.Equal(a_bytes, b_bytes)
}
