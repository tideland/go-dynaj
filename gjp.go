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

//--------------------
// TYPES
//--------------------

// Node may be a JSON object, array or value.
type Node = any

// Object represents a JSON object.
type Object = map[string]any

// Array represents a JSON array.
type Array = []any

// Value contains one JSON value.
type Value = any

// EOF
