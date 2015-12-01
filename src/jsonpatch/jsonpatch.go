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
	"locale"
	"encoding/json"

	evan "github.com/evanphx/json-patch"
)

var _ Patch = JSONPatch{}
var _ Patch = JSONMergePatch{}

// JSONPatch is an RFC 6902 JSON Patch document.
type JSONPatch evan.Patch

// JSONMergePatch is an RFC 7286 JSON Merge Patch document.
type JSONMergePatch json.RawMessage

// A Patch document.
type Patch interface {
	Apply(in interface{}, out interface{}) locale.Error
}

// Apply the patch to an object; as if the object were
// marshalled/unmarshalled JSON.
func (patch JSONPatch) Apply(in interface{}, out interface{}) locale.Error {
	inBytes, uerr := json.Marshal(in)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	outBytes, uerr := (evan.Patch(patch)).Apply(inBytes)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	return locale.UntranslatedError(json.Unmarshal(outBytes, out))
}

// Apply the patch to an object; as if the object were
// marshalled/unmarshalled JSON.
func (patch JSONMergePatch) Apply(in interface{}, out interface{}) locale.Error {
	inBytes, uerr := json.Marshal(in)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	outBytes, uerr := evan.MergePatch(inBytes, patch)
	if uerr != nil {
		return locale.UntranslatedError(uerr)
	}
	return locale.UntranslatedError(json.Unmarshal(outBytes, out))
}
