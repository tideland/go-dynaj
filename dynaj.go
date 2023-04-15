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

//--------------------
// CONSTANTS
//--------------------

const (
	// Separator is the default separator for paths.
	Separator = "/"
)

//--------------------
// TYPES
//--------------------

// Path represents a path in a JSON document. It is a string using
// the Separator as separator between the keys and indices.
type Path = string

// Key represents a key or string index in a JSON object.
type Key = string

// Keys represents a list of keys.
type Keys = []Key

// Element represents a JSON element, i.e. a simple value, an object or an array.
type Element = any

// Value represents a simple JSON value.
type Value = any

// Object represents a JSON object.
type Object = map[string]any

// Array represents a JSON array.
type Array = []any

// EOF
