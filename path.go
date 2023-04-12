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

// splitPath splits and cleans the path into parts.
func splitPath(path string) []string {
	parts := strings.Split(path, Separator)
	out := []string{}
	for _, part := range parts {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// ht retrieves head and tail from a list of keys.
func ht(keys []string) (string, []string) {
	switch len(keys) {
	case 0:
		return "", []string{}
	case 1:
		return keys[0], []string{}
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
func valueAt(node Node, keys []string) (Node, error) {
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
func pathify(parts []string) string {
	return Separator + strings.Join(parts, Separator)
}

// process processes node recursively.
func process(node any, keys []string, processor ValueProcessor) error {
	mkerr := func(err error, ps []string) error {
		return fmt.Errorf("cannot process '%s': %v", pathify(ps), err)
	}

	switch tnode := node.(type) {
	case Object:
		// A JSON object.
		if len(tnode) == 0 {
			return (&PathValue{
				path: pathify(keys),
			}).Process(processor)
		}
		for key, subnode := range tnode {
			objectKeys := append(keys, key)
			if err := process(subnode, objectKeys, processor); err != nil {
				return mkerr(err, keys)
			}
		}
	case Array:
		// A JSON array.
		if len(tnode) == 0 {
			return (&PathValue{
				path: pathify(keys),
			}).Process(processor)
		}
		for index, subnode := range tnode {
			arrayKeys := append(keys, strconv.Itoa(index))
			if err := process(subnode, arrayKeys, processor); err != nil {
				return mkerr(err, keys)
			}
		}
	default:
		// A single value at the end.
		return (&PathValue{
			path: pathify(keys),
			node: tnode,
		}).Process(processor)
	}
	return nil
}

// EOF
