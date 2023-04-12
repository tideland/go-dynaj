// Tideland Go Generic JSON Processing - Value
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp // import "tideland.dev/go/gjp"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"reflect"
	"strconv"
)

//--------------------
// PATH VALUE
//--------------------

// ValueProcessor describes a function for the processing of
// values while iterating over a document.
type ValueProcessor func(pv *PathValue) error

// PathValue is the combination of path and its node value.
type PathValue struct {
	path string
	node Node
	err  error
}

// IsUndefined returns true if this value is undefined.
func (pv *PathValue) IsUndefined() bool {
	return pv.node == nil && pv.err == nil
}

// IsError returns true if this value is an error.
func (pv *PathValue) IsError() bool {
	return pv.err != nil
}

// Err returns the error if there is one.
func (pv *PathValue) Err() error {
	return pv.err
}

// AsString returns the value as string.
func (pv *PathValue) AsString(dv string) string {
	if pv.IsUndefined() {
		return dv
	}
	switch tv := pv.node.(type) {
	case string:
		return tv
	case int:
		return strconv.Itoa(tv)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(tv)
	}
	return dv
}

// AsInt returns the value as int.
func (pv *PathValue) AsInt(dv int) int {
	if pv.IsUndefined() {
		return dv
	}
	switch tv := pv.node.(type) {
	case string:
		i, err := strconv.Atoi(tv)
		if err != nil {
			return dv
		}
		return i
	case int:
		return tv
	case float64:
		return int(tv)
	case bool:
		if tv {
			return 1
		}
		return 0
	}
	return dv
}

// AsFloat64 returns the value as float64.
func (pv *PathValue) AsFloat64(dv float64) float64 {
	if pv.IsUndefined() {
		return dv
	}
	switch tv := pv.node.(type) {
	case string:
		f, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			return dv
		}
		return f
	case int:
		return float64(tv)
	case float64:
		return tv
	case bool:
		if tv {
			return 1.0
		}
		return 0.0
	}
	return dv
}

// AsBool returns the value as bool.
func (pv *PathValue) AsBool(dv bool) bool {
	if pv.IsUndefined() {
		return dv
	}
	switch tv := pv.node.(type) {
	case string:
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return dv
		}
		return b
	case int:
		return tv == 1
	case float64:
		return tv == 1.0
	case bool:
		return tv
	}
	return dv
}

// Equals compares a value with the passed one.
func (pv *PathValue) Equals(to *PathValue) bool {
	switch {
	case pv.IsUndefined() && to.IsUndefined():
		return true
	case pv.IsUndefined() || to.IsUndefined():
		return false
	default:
		return reflect.DeepEqual(pv.node, to.node)
	}
}

// Process processes the value with the passed processor.
func (pv *PathValue) Process(process ValueProcessor) error {
	return process(pv)
}

// Path returns the path of the value.
func (pv *PathValue) Path() string {
	return pv.path
}

// SplitPath splits the path into its keys.
func (pv *PathValue) SplitPath() []string {
	return splitPath(pv.path)
}

// String implements fmt.Stringer.
func (pv *PathValue) String() string {
	if pv.IsUndefined() {
		return "null"
	}
	if pv.IsError() {
		return fmt.Sprintf("error: %v", pv.err)
	}
	return fmt.Sprintf("%v", pv.node)
}

// PathValues contains a list of paths and their values.
type PathValues []*PathValue

// EOF
