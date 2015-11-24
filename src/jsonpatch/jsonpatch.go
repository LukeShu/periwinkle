// Copyright 2015 Luke Shumaker

// Package jsonpatch provides sane bindings to RFC 6902 (JSON Patch)
// and RFC 7386 (JSON Merge Patch).
//
// To read in a Patch or Merge Patch, declare the variable, then use
// json.Unmarshal(); like you would any other JSON document:
//
//     bytes, err := ioutil.ReadAll(inputstream)
//     var patch jsonpatch.JSONPatch
//     err := json.Unmarshal(bytes, &patch)
//     err := patch.Apply(oldObj, &newObj)
//
// or
//
//     bytes, err := ioutil.ReadAll(inputstream)
//     var patch jsonpatch.JSONMergePatch
//     err := json.Unmarshal(bytes, &patch)
//     err := patch.Apply(oldObj, &newObj)
//
// It does this by wrapping github.com/evanphx/json-patch ; there's
// very little actual code here; it is just a set of wrappers to
// provide a nicer interface.
package jsonpatch

import (
	"encoding/json"

	evan "github.com/evanphx/json-patch"
)

var _ Patch = JSONPatch{}
var _ Patch = JSONMergePatch{}

type JSONPatch evan.Patch
type JSONMergePatch json.RawMessage

type Patch interface {
	Apply(in interface{}, out interface{}) error
}

func (patch JSONPatch) Apply(in interface{}, out interface{}) error {
	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	outBytes, err := (evan.Patch(patch)).Apply(inBytes)
	if err != nil {
		return err
	}
	return json.Unmarshal(outBytes, out)
}

func (patch JSONMergePatch) Apply(in interface{}, out interface{}) error {
	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	outBytes, err := evan.MergePatch(inBytes, patch)
	if err != nil {
		return err
	}
	return json.Unmarshal(outBytes, out)
}
