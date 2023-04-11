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
func splitPath(path, separator string) []string {
	// Split the path by the separator.
	parts := strings.Split(path, separator)
	out := []string{}
	for _, part := range parts {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// ht retrieves head and tail from parts.
func ht(parts []string) (string, []string) {
	switch len(parts) {
	case 0:
		return "", []string{}
	case 1:
		return parts[0], []string{}
	default:
		return parts[0], parts[1:]
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
func valueAt(node Node, path []string) (Node, error) {
	if len(path) == 0 {
		// End of the path.
		return node, nil
	}
	// Further access depends on part content node and type.
	h, t := ht(path)
	if h == "" {
		return node, nil
	}
	switch n := node.(type) {
	case Object:
		// JSON object.
		field, ok := n[h]
		if !ok {
			return nil, fmt.Errorf("invalid path %q", path)
		}
		return valueAt(field, t)
	case Array:
		// JSON array.
		index, err := strconv.Atoi(h)
		if err != nil || index >= len(n) {
			return nil, fmt.Errorf("invalid path %q: %v", h, err)
		}
		return valueAt(n[index], t)
	}
	// Path is longer than existing node structure.
	return nil, fmt.Errorf("path is too long")
}

// pathify creates a path out of parts and separator.
func pathify(parts []string, separator string) string {
	return separator + strings.Join(parts, separator)
}

// process processes node recursively.
func process(node any, parts []string, separator string, processor ValueProcessor) error {
	mkerr := func(err error, ps []string) error {
		return fmt.Errorf("cannot process '%s': %v", pathify(ps, separator), err)
	}

	switch n := node.(type) {
	case map[string]any:
		// A JSON object.
		if len(n) == 0 {
			return (&PathValue{
				Path:      pathify(parts, separator),
				Separator: separator,
			}).Process(processor)
		}
		for field, subnode := range n {
			fieldparts := append(parts, field)
			if err := process(subnode, fieldparts, separator, processor); err != nil {
				return mkerr(err, parts)
			}
		}
	case []any:
		// A JSON array.
		if len(n) == 0 {
			return (&PathValue{
				Path:      pathify(parts, separator),
				Separator: separator,
			}).Process(processor)
		}
		for index, subnode := range n {
			indexparts := append(parts, strconv.Itoa(index))
			if err := process(subnode, indexparts, separator, processor); err != nil {
				return mkerr(err, parts)
			}
		}
	default:
		// A single value at the end.
		pv := &PathValue{
			Path:      pathify(parts, separator),
			Separator: separator,
			Value:     n,
		}
		return pv.Process(processor)
	}

	return nil
}

// EOF
