// Tideland Go Generic JSON Processor
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
	"encoding/json"
	"fmt"
	"strconv"
)

//--------------------
// DOCUMENT
//--------------------

// Document represents one JSON document.
type Document struct {
	root Node
}

// Unmarshal parses the JSON-encoded data and stores the result
// as new document.
func Unmarshal(data []byte) (*Document, error) {
	var root any
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal document: %v", err)
	}
	return &Document{
		root: root,
	}, nil
}

// NewDocument creates a new empty document.
func NewDocument() *Document {
	return &Document{}
}

// Length returns the number of elements for the given path.
func (d *Document) Length(path string) int {
	node, err := valueAt(d.root, splitPath(path))
	if err != nil {
		return -1
	}
	// Return len based on type.
	switch n := node.(type) {
	case Object:
		return len(n)
	case Array:
		return len(n)
	default:
		return 1
	}
}

// SetValueAt sets the value at the given path.
func (d *Document) SetValueAt(path string, value Value) error {
	keys := splitPath(path)
	root, err := insertValueInNode(d.root, keys, value)
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

// ValueAt returns the addressed value.
func (d *Document) ValueAt(path string) *PathValue {
	pv := &PathValue{
		path: path,
	}
	node, err := valueAt(d.root, splitPath(path))
	if err != nil {
		pv.err = fmt.Errorf("invalid path %q: %v", path, err)
	} else {
		pv.node = node
	}
	return pv
}

// Root returns the root path value.
func (d *Document) Root() *PathValue {
	return &PathValue{
		path: Separator,
		node: d.root,
	}
}

// Clear removes the document data.
func (d *Document) Clear() {
	d.root = nil
}

// MarshalJSON implements json.Marshaler.
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.root)
}

// String implements fmt.Stringer.
func (d *Document) String() string {
	data, err := json.Marshal(d.root)
	if err != nil {
		return fmt.Sprintf("cannot marshal document: %v", err)
	}
	return string(data)
}

//--------------------
// DOCUMENT HELPERS
//--------------------

// insertValue recursively inserts a value at the end of the keys list.
func insertValueInNode(node Node, keys []string, value Value) (Node, error) {
	if len(keys) == 0 {
		return value, nil
	}

	switch tnode := node.(type) {
	case nil:
		return createValue(keys, value)
	case Object:
		return insertValueInObject(tnode, keys, value)
	case Array:
		return insertValueInArray(tnode, keys, value)
	default:
		return nil, fmt.Errorf("document is not a valid JSON structure")
	}
}

// createValue creates a value at the end of the keys list.
func createValue(keys []string, value Value) (Node, error) {
	// Check if we are at the end of the keys list.
	if len(keys) == 0 {
		return value, nil
	}
	h, t := ht(keys)
	// Check for array index first.
	index, err := strconv.Atoi(h)
	if err == nil {
		// It's an array index.
		arr := make([]any, index+1)
		arr[index], err = createValue(t, value)
		if err != nil {
			return nil, err
		}
		return arr, nil
	}
	// It's an object key.
	obj := Object{h: nil}
	obj[h], err = createValue(t, value)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// insertValueInObject inserts a value in a JSON object at the end of the keys list.
func insertValueInObject(obj Object, keys []string, value Value) (Node, error) {
	h, t := ht(keys)
	// Create object if keys list has only one element.
	if len(t) == 0 {
		if isObjectOrArray(obj[h]) {
			return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
		}
		obj[h] = value
		return obj, nil
	}
	// Insert value in node.
	node := obj[h]
	if isValue(node) {
		return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
	}
	newNode, err := insertValueInNode(node, t, value)
	if err != nil {
		return nil, err
	}

	obj[h] = newNode
	return obj, nil
}

// insertValueInArray inserts a value in an array at a given path.
func insertValueInArray(arr []any, path []string, value Value) (Node, error) {
	h, t := ht(path)
	// Convert path head into index.
	index, err := strconv.Atoi(h)
	switch {
	case err != nil:
		return nil, fmt.Errorf("invalid index %q in array", h)
	case index < 0:
		return nil, fmt.Errorf("negative index %d for array", index)
	case index >= len(arr):
		tmp := make(Array, index+1)
		copy(tmp, arr)
		arr = tmp
	}
	// Insert value if last element in path.
	if len(t) == 0 {
		if isObjectOrArray(arr[index]) {
			return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", path)
		}
		arr[index] = value
		return arr, nil
	}
	// Insert value in node.
	node := arr[index]
	if isValue(node) {
		return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", path)
	}
	newNode, err := insertValueInNode(node, t, value)
	if err != nil {
		return nil, err
	}

	arr[index] = newNode
	return arr, nil
}

// EOF
