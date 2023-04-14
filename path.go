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
	"fmt"
	"strconv"
	"strings"
)

//--------------------
// PROCESSING FUNCTIONS
//--------------------

// splitPath splits and cleans the path into keys.
func splitPath(path Path) Keys {
	keys := strings.Split(path, Separator)
	out := []string{}
	for _, key := range keys {
		if key != "" {
			out = append(out, key)
		}
	}
	return out
}

// ht retrieves head and tail from a list of keys.
func ht(keys Keys) (Key, Keys) {
	switch len(keys) {
	case 0:
		return "", Keys{}
	case 1:
		return keys[0], Keys{}
	default:
		return keys[0], keys[1:]
	}
}

// isObjectOrArray checks if the node is an object or an array.
func isObjectOrArray(node Node) bool {
	switch node.(type) {
	case Object, Array:
		return true
	default:
		return false
	}
}

// isValue checks if the node is a value.
func isValue(node Node) bool {
	switch node.(type) {
	case Object, Array, nil:
		return false
	default:
		return true
	}
}

// valueAt returns the value at the given path.
func valueAt(node Node, keys Keys) (Node, error) {
	if len(keys) == 0 {
		// End of the path.
		return node, nil
	}
	// Further access depends on part content node and type.
	h, t := ht(keys)
	if h == "" {
		return node, nil
	}
	switch n := node.(type) {
	case Object:
		// JSON object.
		field, ok := n[h]
		if !ok {
			return nil, fmt.Errorf("invalid path %q", pathify(keys))
		}
		return valueAt(field, t)
	case Array:
		// JSON array.
		index, err := strconv.Atoi(h)
		if err != nil || index >= len(n) {
			return nil, fmt.Errorf("invalid path %q: %v", pathify(keys), err)
		}
		return valueAt(n[index], t)
	}
	// Path is longer than existing node structure.
	return nil, fmt.Errorf("path is too long")
}

// pathify creates a path out of keys.
func pathify(keys Keys) Path {
	return Separator + strings.Join(keys, Separator)
}

// appendKey appends a key to a path.
func appendKey(path Path, key Key) Path {
	if len(path) == 1 {
		// Root path.
		return path + key
	}
	return path + Separator + key
}

// EOF
