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

// PathValue is the combination of path, separator and value.
type PathValue struct {
	Path      string
	Separator string
	Value     Value
}

// IsUndefined returns true if this value is undefined.
func (pv *PathValue) IsUndefined() bool {
	return pv.Value == nil || pv.IsError()
}

// IsError returns true if this value is an error.
func (pv *PathValue) IsError() bool {
	_, ok := pv.Value.(error)
	return ok
}

// AsError returns the error value in case of an error.
func (pv *PathValue) AsError() error {
	if pv.IsError() {
		return pv.Value.(error)
	}
	return nil
}

// AsString returns the value as string.
func (pv *PathValue) AsString(dv string) string {
	if pv.IsUndefined() {
		return dv
	}
	switch tv := pv.Value.(type) {
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
	switch tv := pv.Value.(type) {
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
	switch tv := pv.Value.(type) {
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
	switch tv := pv.Value.(type) {
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
		return reflect.DeepEqual(pv.Value, to.Value)
	}
}

// Process processes the value with the passed processor.
func (pv *PathValue) Process(process ValueProcessor) error {
	return process(pv)
}

// SplitPath splits the path into its parts.
func (pv *PathValue) SplitPath() []string {
	return splitPath(pv.Path, pv.Separator)
}

// String implements fmt.Stringer.
func (pv *PathValue) String() string {
	if pv.IsUndefined() {
		return "null"
	}
	return fmt.Sprintf("%v", pv.Value)
}

// PathValues contains a list of paths and their values.
type PathValues []PathValue

// EOF
