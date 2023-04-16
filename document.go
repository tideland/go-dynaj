// Tideland Go Dynamic JSON
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
	"encoding/json"
	"fmt"
)

//--------------------
// DOCUMENT
//--------------------

// Document represents one JSON document.
type Document struct {
	root Element
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
func (d *Document) Length(path Path) int {
	node, err := elementAt(d.root, splitPath(path))
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
func (d *Document) SetValueAt(path Path, value Value) error {
	keys := splitPath(path)
	root, err := insertValue(d.root, keys, value)
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

// DeleteValueAt deletes the value at the given path. If it is inside
// an object the key is deleted, if it is inside an array the elements
// are shifted.
func (d *Document) DeleteValueAt(path Path) error {
	keys := splitPath(path)
	root, err := deleteElement(d.root, keys, false)
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

// NodeAt returns the addressed value.
func (d *Document) NodeAt(path Path) *Node {
	node := &Node{
		path: path,
	}
	element, err := elementAt(d.root, splitPath(path))
	if err != nil {
		node.err = fmt.Errorf("invalid path %q: %v", path, err)
	} else {
		node.element = element
	}
	return node
}

// Root returns the root path value.
func (d *Document) Root() *Node {
	return &Node{
		path:    Separator,
		element: d.root,
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

// EOF
