// Tideland Go Generic JSON Processing - Value
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj // import "tideland.dev/go/dynaj"

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
// NODE
//--------------------

// Processor defines the signature of function for processing
// a path value. This may be the iterating over the whole
// document or one object or array.
type Processor func(n *Node) error

// Node is the combination of path and its value.
type Node struct {
	path  Path
	value Value
	err   error
}

// IsUndefined returns true if this value is undefined.
func (node *Node) IsUndefined() bool {
	return node.value == nil && node.err == nil
}

// IsValue returns true if this node is a simple value.
func (node *Node) IsValue() bool {
	switch node.value.(type) {
	case Object, Array:
		return false
	default:
		return true
	}
}

// IsObject returns true if this node is an object.
func (node *Node) IsObject() bool {
	_, ok := node.value.(Object)
	return ok
}

// IsArray returns true if this node is an array.
func (node *Node) IsArray() bool {
	_, ok := node.value.(Array)
	return ok
}

// IsError returns true if this value is an error.
func (node *Node) IsError() bool {
	return node.err != nil
}

// Err returns the error if there is one.
func (node *Node) Err() error {
	return node.err
}

// AsString returns the value as string.
func (node *Node) AsString(dv string) string {
	if node.IsUndefined() {
		return dv
	}
	switch tv := node.value.(type) {
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
func (node *Node) AsInt(dv int) int {
	if node.IsUndefined() {
		return dv
	}
	switch tv := node.value.(type) {
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
func (node *Node) AsFloat64(dv float64) float64 {
	if node.IsUndefined() {
		return dv
	}
	switch tv := node.value.(type) {
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
func (node *Node) AsBool(dv bool) bool {
	if node.IsUndefined() {
		return dv
	}
	switch tv := node.value.(type) {
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
func (node *Node) Equals(other *Node) bool {
	switch {
	case node.IsUndefined() && other.IsUndefined():
		return true
	case node.IsUndefined() || other.IsUndefined():
		return false
	default:
		return reflect.DeepEqual(node.value, other.value)
	}
}

// Path returns the path of the value.
func (node *Node) Path() Path {
	return node.path
}

// SplitPath splits the path into its keys.
func (node *Node) SplitPath() Keys {
	return splitPath(node.path)
}

// NodeAt returns the node at the passed path.
func (node *Node) NodeAt(path Path) *Node {
	if node.IsUndefined() {
		return &Node{
			path:  path,
			value: nil,
		}
	}
	if node.IsValue() {
		return &Node{
			path:  path,
			value: node.value,
		}
	}
	// Navigate downstream.
	nodeAt := &Node{
		path: joinPaths(node.path, path),
	}
	value, err := elementAt(node.value, splitPath(path))
	if err != nil {
		nodeAt.err = fmt.Errorf("invalid path %q: %v", path, err)
	} else {
		nodeAt.value = value
	}
	return nodeAt
}

// Process iterates over the node and all its subnodes and
// processes them with the passed processor function.
func (node *Node) Process(process Processor) error {
	if node.err != nil {
		return node.err
	}
	switch typed := node.value.(type) {
	case Object:
		// A JSON object.
		if len(typed) == 0 {
			return process(&Node{
				path:  node.path,
				value: Object{},
			})
		}
		for key, subvalue := range typed {
			subpath := appendKey(node.path, key)
			subnode := &Node{
				path:  subpath,
				value: subvalue,
			}
			if err := subnode.Process(process); err != nil {
				return fmt.Errorf("cannot process %q: %v", subpath, err)
			}
		}
	case Array:
		// A JSON array.
		if len(typed) == 0 {
			return process(&Node{
				path:  node.path,
				value: Array{},
			})
		}
		for idx, subvalue := range typed {
			subpath := appendKey(node.path, strconv.Itoa(idx))
			subnode := &Node{
				path:  subpath,
				value: subvalue,
			}
			if err := subnode.Process(process); err != nil {
				return fmt.Errorf("cannot process %q: %v", subpath, err)
			}
		}
	default:
		// A single value at the end.
		err := process(&Node{
			path:  node.path,
			value: typed,
		})
		if err != nil {
			return fmt.Errorf("cannot process %q: %v", node.path, err)
		}
	}
	return nil
}

// Range takes  the node and processes it with the passed processor
// function. In case of an object all keys and in case of an array
// all indices will be processed. It is not working recursively.
func (node *Node) Range(process Processor) error {
	if node.err != nil {
		return node.err
	}
	switch typed := node.value.(type) {
	case Object:
		// A JSON object.
		for key := range typed {
			keypath := appendKey(node.path, key)
			if isObjectOrArray(typed[key]) {
				return fmt.Errorf("cannot process %q: is object or array", keypath)
			}
			err := process(&Node{
				path:  keypath,
				value: typed[key],
			})
			if err != nil {
				return fmt.Errorf("cannot process %q: %v", keypath, err)
			}
		}
	case Array:
		// A JSON array.
		for idx := range typed {
			idxpath := appendKey(node.path, strconv.Itoa(idx))
			if isObjectOrArray(typed[idx]) {
				return fmt.Errorf("cannot process %q: is object or array", idxpath)
			}
			err := process(&Node{
				path:  idxpath,
				value: typed[idx],
			})
			if err != nil {
				return fmt.Errorf("cannot process %q: %v", idxpath, err)
			}
		}
	default:
		// A single value at the end.
		err := process(&Node{
			path:  node.path,
			value: typed,
		})
		if err != nil {
			return fmt.Errorf("cannot process %q: %v", node.path, err)
		}
	}
	return nil
}

// Query iterates over the node and all its subnodes and returns
// all values with paths matching the passed pattern.
func (node *Node) Query(pattern string) (Nodes, error) {
	nodes := Nodes{}
	err := node.Process(func(pnode *Node) error {
		trimmedPath := strings.TrimPrefix(pnode.path, node.path+Separator)
		if matcher.Matches(pattern, trimmedPath, false) {
			nodes = append(nodes, &Node{
				path:  pnode.path,
				value: pnode.value,
			})
		}
		return nil
	})
	return nodes, err
}

// String implements fmt.Stringer.
func (node *Node) String() string {
	if node.IsUndefined() {
		return "null"
	}
	if node.IsError() {
		return fmt.Sprintf("error: %v", node.err)
	}
	return fmt.Sprintf("%v", node.value)
}

// Nodes contains a list of paths and their value.
type Nodes []*Node

// EOF
