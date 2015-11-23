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

	evan "github.com/evanphx/json-patch"
	matt "github.com/mattbaird/jsonpatch"
)

// Diff two objects, and produce an RFC 7386 JSON Merge Patch.
func NewJSONMergePatch(a interface{}, b interface{}) (jsonpatch.JSONMergePatch, error) {
	// convert a to json
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	// convert b to json
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	// diff them
	p_bytes, err := evan.CreateMergePatch(a_bytes, b_bytes)
	if err != nil {
		return nil, err
	}
	// return
	var ret jsonpatch.JSONMergePatch
	err = json.Unmarshal(p_bytes, &ret)
	return ret, err
}

// Diff two objects, and produce an RFC 6902 JSON Patch.
func NewJSONPatch(a interface{}, b interface{}) (jsonpatch.JSONPatch, error) {
	// convert a to json
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	// convert b to json
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	// diff them
	p, err := matt.CreatePatch(a_bytes, b_bytes)
	p_bytes, err := json.Marshal(p)
	// return
	var ret jsonpatch.JSONPatch
	err = json.Unmarshal(p_bytes, &ret)
	return ret, err
}

// Test whether two objects have equivalent JSON structures.
func Equal(a interface{}, b interface{}) bool {
	// convert a to json
	a_bytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	// convert b to json
	b_bytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	// diff them and return
	return evan.Equal(a_bytes, b_bytes)
}
