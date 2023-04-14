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
	"strings"

	"tideland.dev/go/matcher"
)

//--------------------
// PATH VALUE
//--------------------

// Processor defines the signature of function for processing
// a path value. This may be the iterating over the whole
// document or one object or array.
type Processor func(pv *PathValue) error

// PathValue is the combination of path and its node value.
type PathValue struct {
	path Path
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

// Path returns the path of the value.
func (pv *PathValue) Path() Path {
	return pv.path
}

// SplitPath splits the path into its keys.
func (pv *PathValue) SplitPath() Keys {
	return splitPath(pv.path)
}

// Process iterates over the node and all its subnodes and
// processes them with the passed processor function.
func (pv *PathValue) Process(process Processor) error {
	if pv.err != nil {
		return pv.err
	}
	switch tnode := pv.node.(type) {
	case Object:
		// A JSON object.
		if len(tnode) == 0 {
			return process(&PathValue{
				path: pv.path,
				node: Object{},
			})
		}
		for key, subnode := range tnode {
			subpath := appendKey(pv.path, key)
			subvalue := &PathValue{
				path: subpath,
				node: subnode,
			}
			if err := subvalue.Process(process); err != nil {
				return fmt.Errorf("cannot process %q: %v", subpath, err)
			}
		}
	case Array:
		// A JSON array.
		if len(tnode) == 0 {
			return process(&PathValue{
				path: pv.path,
				node: Array{},
			})
		}
		for idx, subnode := range tnode {
			subpath := appendKey(pv.path, strconv.Itoa(idx))
			subvalue := &PathValue{
				path: subpath,
				node: subnode,
			}
			if err := subvalue.Process(process); err != nil {
				return fmt.Errorf("cannot process %q: %v", subpath, err)
			}
		}
	default:
		// A single value at the end.
		err := process(&PathValue{
			path: pv.path,
			node: tnode,
		})
		if err != nil {
			return fmt.Errorf("cannot process %q: %v", pv.path, err)
		}
	}
	return nil
}

// Range takes  the node and processes it with the passed processor
// function. In case of an object all keys and in case of an array
// all indices will be processed. It is not working recursively.
func (pv *PathValue) Range(process Processor) error {
	if pv.err != nil {
		return pv.err
	}
	switch tnode := pv.node.(type) {
	case Object:
		// A JSON object.
		for key := range tnode {
			keypath := appendKey(pv.path, key)
			if isObjectOrArray(tnode[key]) {
				return fmt.Errorf("cannot process %q: is object or array", keypath)
			}
			err := process(&PathValue{
				path: keypath,
				node: tnode[key],
			})
			if err != nil {
				return fmt.Errorf("cannot process %q: %v", keypath, err)
			}
		}
	case Array:
		// A JSON array.
		for idx := range tnode {
			idxpath := appendKey(pv.path, strconv.Itoa(idx))
			if isObjectOrArray(tnode[idx]) {
				return fmt.Errorf("cannot process %q: is object or array", idxpath)
			}
			err := process(&PathValue{
				path: idxpath,
				node: tnode[idx],
			})
			if err != nil {
				return fmt.Errorf("cannot process %q: %v", idxpath, err)
			}
		}
	default:
		// A single value at the end.
		err := process(&PathValue{
			path: pv.path,
			node: tnode,
		})
		if err != nil {
			return fmt.Errorf("cannot process %q: %v", pv.path, err)
		}
	}
	return nil
}

// Query iterates over the node and all its subnodes and returns
// all values with paths matching the passed pattern.
func (pv *PathValue) Query(pattern string) (PathValues, error) {
	pvs := PathValues{}
	err := pv.Process(func(ppv *PathValue) error {
		ppvpath := strings.TrimPrefix(ppv.path, pv.path+Separator)
		println(pattern + "  =>  " + ppvpath)
		if matcher.Matches(pattern, ppvpath, false) {
			pvs = append(pvs, &PathValue{
				path: ppv.path,
				node: ppv.node,
			})
		}
		return nil
	})
	return pvs, err
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
