// Copyright 2015 Luke Shumaker

// Package jsondiff provides utilites for creating RFC 6902 (JSON Patch)
// and RFC 7386 (JSON Merge Patch) patches.
//
// It does this by wrapping github.com/evanphx/json-patch and
// github.com/mattbaird/jsonpatch ; there's very little actual code
// here; it is just a set of wrappers to provide a nicer interface.
package jsondiff

import (
	"encoding/json"
	"jsonpatch"
	"locale"

	evan "github.com/evanphx/json-patch"
	matt "github.com/mattbaird/jsonpatch"
)

// Diff two objects, and produce an RFC 7386 JSON Merge Patch.
func NewJSONMergePatch(a interface{}, b interface{}) (jsonpatch.JSONMergePatch, locale.Error) {
	// convert a to json
	aBytes, uerr := json.Marshal(a)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// convert b to json
	bBytes, uerr := json.Marshal(b)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// diff them
	pBytes, uerr := evan.CreateMergePatch(aBytes, bBytes)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// return
	var ret jsonpatch.JSONMergePatch
	uerr = json.Unmarshal(pBytes, &ret)
	return ret, locale.UntranslatedError(uerr)
}

// Diff two objects, and produce an RFC 6902 JSON Patch.
func NewJSONPatch(a interface{}, b interface{}) (jsonpatch.JSONPatch, locale.Error) {
	// convert a to json
	aBytes, uerr := json.Marshal(a)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// convert b to json
	bBytes, uerr := json.Marshal(b)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// diff them
	p, uerr := matt.CreatePatch(aBytes, bBytes)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	pBytes, uerr := json.Marshal(p)
	if uerr != nil {
		return nil, locale.UntranslatedError(uerr)
	}
	// return
	var ret jsonpatch.JSONPatch
	uerr = json.Unmarshal(pBytes, &ret)
	return ret, locale.UntranslatedError(uerr)
}

// Test whether two objects have equivalent JSON structures.
func Equal(a interface{}, b interface{}) bool {
	// convert a to json
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	// convert b to json
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	// diff them and return
	return evan.Equal(aBytes, bBytes)
}
