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

	"tideland.dev/go/matcher"
)

//--------------------
// DOCUMENT
//--------------------

// Document represents one JSON document.
type Document struct {
	separator string
	root      any
}

// Parse reads a raw document and returns it as
// accessible document.
func Parse(data []byte, separator string) (*Document, error) {
	var root any
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal document: %v", err)
	}
	return &Document{
		separator: separator,
		root:      root,
	}, nil
}

// NewDocument creates a new empty document.
func NewDocument(separator string) *Document {
	return &Document{
		separator: separator,
	}
}

// Length returns the number of elements for the given path.
func (d *Document) Length(path string) int {
	node, err := valueAt(d.root, splitPath(path, d.separator))
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
	pathParts := splitPath(path, d.separator)
	root, err := insertValueInNode(d.root, pathParts, value)
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

// ValueAt returns the addressed value.
func (d *Document) ValueAt(path string) *PathValue {
	pv := &PathValue{
		Path:      path,
		Separator: d.separator,
	}
	n, err := valueAt(d.root, splitPath(path, d.separator))
	if err != nil {
		pv.Value = fmt.Errorf("cannot get value at %q: %v", path, err)
	} else {
		pv.Value = n
	}
	return pv
}

// Clear removes the document data.
func (d *Document) Clear() {
	d.root = nil
}

// Query allows to find pathes matching a given pattern.
func (d *Document) Query(pattern string) (PathValues, error) {
	pvs := PathValues{}
	err := d.Process(func(pv *PathValue) error {
		if matcher.Matches(pattern, pv.Path, false) {
			pvs = append(pvs, PathValue{
				Path:      pv.Path,
				Separator: pv.Separator,
				Value:     pv.Value,
			})
		}
		return nil
	})
	return pvs, err
}

// Process iterates over a document and processes its values.
// There's no order, so nesting into an embedded document or
// list may come earlier than higher level paths.
func (d *Document) Process(processor ValueProcessor) error {
	return process(d.root, []string{}, d.separator, processor)
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

// insertValue recursively inserts a value at a given path.
func insertValueInNode(node Node, path []string, value Value) (Node, error) {
	if len(path) == 0 {
		// return nil, fmt.Errorf("path cannot be empty")
		return value, nil
	}

	switch tnode := node.(type) {
	case nil:
		return createValue(path, value)
	case Object:
		return insertValueInObject(tnode, path, value)
	case Array:
		return insertValueInArray(tnode, path, value)
	default:
		return nil, fmt.Errorf("document is not a valid JSON structure")
	}
}

// createValue creates a value at a given path.
func createValue(path []string, value Value) (Node, error) {
	h, t := ht(path)
	// Create object if last element in path.
	if len(t) == 0 {
		return Object{h: value}, nil
	}
	// Create array if path head is an index.
	index, err := strconv.Atoi(h)
	switch {
	case err == nil:
		arr := make([]any, index+1)
		arr[index], err = createValue(t, value)
		if err != nil {
			return nil, err
		}
		return arr, nil
	case err != nil && len(t) == 0:
		return Object{h: value}, nil
	}
	// Create object.
	obj := Object{h: nil}
	obj[h], err = createValue(t, value)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// insertValueInObject inserts a value in a JSON object at a given path.
func insertValueInObject(obj Object, path []string, value Value) (Node, error) {
	h, t := ht(path)
	// Insert value if last element in path.
	if len(t) == 0 {
		if isObjectOrArray(obj[h]) {
			return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", path)
		}
		obj[h] = value
		return obj, nil
	}
	// Insert value in node.
	node := obj[h]
	if isValue(node) {
		return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", path)
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
