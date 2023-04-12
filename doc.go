// Tideland Go Generic JSON Processor
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package gjp provides the generic parsing and processing of JSON
// documents by paths. Values can be retrieved by paths like "foo/bar/3".
// Functions with given default value help to retrieve the values in
// a type safe way.
//
//	doc, err := gjp.Parse(myDoc)
//	if err != nil {
//	    ...
//	}
//	name := doc.ValueAt("/name").AsString("")
//	street := doc.ValueAt("/address/street").AsString("unknown")
//
// Another way is to create an empty document with
//
//	doc := gjp.NewDocument()
//
// Here as well as in parsed documents values can be set with
//
//	err := doc.SetValueAt("/a/b/3/c", 4711)
//
// Additionally values of the document can be processed recursively
// using
//
//	 err := doc.Process(func(pv *gjp.PathValue) error {
//		    ...
//	 })
//
// To retrieve the distance between two documents the function
// gjp.Compare() can be used:
//
//	diff, err := gjp.Compare(firstDoc, secondDoc)
//
// privides a gjp.Diff instance which helps to compare individual
// paths of the two document.
package gjp // import "tideland.dev/go/text/gjp"

// EOF
