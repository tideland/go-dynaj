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

// joinPaths joins the given paths into one.
func joinPaths(paths ...Path) Path {
	out := Keys{}
	for _, path := range paths {
		out = append(out, splitPath(path)...)
	}
	return pathify(out)
}

// headTail retrieves the head and the tail key from a list of keys.
func headTail(keys Keys) (Key, Keys) {
	switch len(keys) {
	case 0:
		return "", Keys{}
	case 1:
		return keys[0], Keys{}
	default:
		return keys[0], keys[1:]
	}
}

// asIndex converts the given key into an index.
func asIndex(key Key) (int, bool) {
	index, err := strconv.Atoi(key)
	if err != nil {
		return 0, false
	}
	return index, true
}

// elementAt returns the element at the given path recursively
// starting at the given element.
func elementAt(element Element, keys Keys) (Element, error) {
	if len(keys) == 0 {
		// End of the path.
		return element, nil
	}
	// Further access depends on part content node and type.
	h, t := headTail(keys)
	if h == "" {
		return element, nil
	}
	switch typed := element.(type) {
	case Object:
		// JSON object.
		field, ok := typed[h]
		if !ok {
			return nil, fmt.Errorf("invalid path %q", pathify(keys))
		}
		return elementAt(field, t)
	case Array:
		// JSON array.
		index, ok := asIndex(h)
		if !ok {
			return nil, fmt.Errorf("invalid path %q: no index", pathify(keys))
		}
		if index < 0 || index >= len(typed) {
			return nil, fmt.Errorf("invalid path %q: index out of range", pathify(keys))
		}
		return elementAt(typed[index], t)
	}
	// Path is longer than existing node structure.
	return nil, fmt.Errorf("key or index not found")
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

// isObjectOrArray checks if the element is an object or an array.
func isObjectOrArray(element Element) bool {
	switch element.(type) {
	case Object, Array:
		return true
	default:
		return false
	}
}

// isValue checks if the element is a single value.
func isValue(element Element) bool {
	switch element.(type) {
	case Object, Array, nil:
		return false
	default:
		return true
	}
}

// EOF
