// Tideland Go Dynamic JSON
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package dynaj provides a simple dynamic handling of JSON documents.
// Values can be retrieved, set and added by paths like "foo/bar/3".
// Methods provide typesafe access to the values as well as flat and
// deep processing.
//
//	doc, err := dynaj.Unmarshal(myDoc)
//	if err != nil {
//	    ...
//	}
//	name := doc.ValueAt("/name").AsString("")
//	street := doc.ValueAt("/address/street").AsString("unknown")
//
// Another way is to create an empty document with
//
//	doc := dynaj.NewDocument()
//
// Here as well as in parsed documents values can be set with
//
//	err := doc.SetValueAt("/a/b/3/c", 4711)
//
// Additionally values of the document can be processed recursively
// using
//
//	err := doc.Root().Process(func(pv *dynaj.PathValue) error {
//	    ...
//	})
//
// or from deeper nodes with doc.ValueAt("/a/b/3").Process(...).
// Additionally flat processing is possible with
//
//	err := doc.ValueAt("/x/y/z").Range(func(pv *dynaj.PathValue) error {
//	    ...
//	})
//
// To retrieve the differences between two documents the function
// dynaj.Compare() can be used:
//
//	diff, err := dynaj.Compare(firstDoc, secondDoc)
//
// privides a dynaj.Diff instance which helps to compare individual
// paths of the two document.
package dynaj // import "tideland.dev/go/text/dynaj"

// EOF
