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

// isValue checks if the raw is a value and returns it
// type-safe. Otherwise nil and false are returned.
func isValue(raw any) (any, bool) {
	if raw == nil {
		return raw, true
	}
	_, ook := isObject(raw)
	_, aok := isArray(raw)
	if ook || aok {
		return nil, false
	}
	return raw, true
}

// isObject checks if the raw is an object and returns it
// type-safe. Otherwise nil and false are returned.
func isObject(raw any) (map[string]any, bool) {
	o, ok := raw.(map[string]any)
	return o, ok
}

// isArray checks if the raw is an array and returns it
// type-safe. Otherwise nil and false are returned.
func isArray(raw any) ([]any, bool) {
	a, ok := raw.([]any)
	return a, ok
}

// valueAt returns the value at the path parts.
func valueAt(node any, parts []string) (any, error) {
	length := len(parts)
	if length == 0 {
		// End of the parts.
		return node, nil
	}
	// Further access depends on part content node and type.
	head, tail := ht(parts)
	if head == "" {
		return node, nil
	}
	if o, ok := isObject(node); ok {
		// JSON object.
		field, ok := o[head]
		if !ok {
			return nil, fmt.Errorf("invalid path part '%s'", head)
		}
		return valueAt(field, tail)
	}
	if a, ok := isArray(node); ok {
		// JSON array.
		index, err := strconv.Atoi(head)
		if err != nil || index >= len(a) {
			return nil, fmt.Errorf("invalid path part '%s': %v", head, err)
		}
		return valueAt(a[index], tail)
	}
	// Parts left but field value.
	return nil, fmt.Errorf("path is too long")
}

// setValueAt sets the value at the path parts.
func setValueAt(root, value any, parts []string) (any, error) {
	h, t := ht(parts)
	return setNodeValueAt(root, value, h, t)
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

// setNodeValueAt is used recursively by setValueAt().
func setNodeValueAt(node, value any, head string, tail []string) (any, error) {
	// Check for nil node first.
	if node == nil {
		return addNodeValueAt(value, head, tail)
	}
	// Otherwise it should be an object or an array.
	if o, ok := isObject(node); ok {
		// JSON object.
		_, ok := isValue(o[head])
		switch {
		case !ok && len(tail) == 0:
			return nil, fmt.Errorf("setting value corrupts document")
		case ok && o[head] != nil && len(tail) > 0:
			return nil, fmt.Errorf("setting value corrupts document")
		case ok && len(tail) == 0:
			o[head] = value
		default:
			h, t := ht(tail)
			subnode, err := setNodeValueAt(o[head], value, h, t)
			if err != nil {
				return nil, err
			}
			o[head] = subnode
		}
		return o, nil
	}
	if a, ok := isArray(node); ok {
		// JSON array.
		index, err := strconv.Atoi(head)
		if err != nil {
			return nil, fmt.Errorf("invalid path part '%s'", head)
		}
		a = ensureArray(a, index+1)
		_, ok := isValue(a[index])
		switch {
		case !ok && len(tail) == 0:
			return nil, fmt.Errorf("setting value corrupts document")
		case ok && a[index] != nil && len(tail) > 0:
			return nil, fmt.Errorf("setting value corrupts document")
		case ok && len(tail) == 0:
			a[index] = value
		default:
			h, t := ht(tail)
			subnode, err := setNodeValueAt(a[index], value, h, t)
			if err != nil {
				return nil, err
			}
			a[index] = subnode
		}
		return a, nil
	}
	return nil, fmt.Errorf("invalid path part '%s'", head)
}

// addNodeValueAt is used recursively by setValueAt().
func addNodeValueAt(value any, head string, tail []string) (any, error) {
	// JSON value.
	if head == "" {
		return value, nil
	}
	index, err := strconv.Atoi(head)
	if err != nil {
		// JSON object.
		o := map[string]any{}
		if len(tail) == 0 {
			o[head] = value
			return o, nil
		}
		h, t := ht(tail)
		subnode, err := addNodeValueAt(value, h, t)
		if err != nil {
			return nil, err
		}
		o[head] = subnode
		return o, nil
	}
	// JSON array.
	a := ensureArray([]any{}, index+1)
	if len(tail) == 0 {
		a[index] = value
		return a, nil
	}
	h, t := ht(tail)
	subnode, err := addNodeValueAt(value, h, t)
	if err != nil {
		return nil, err
	}
	a[index] = subnode
	return a, nil
}

// ensureArray ensures the right len of an array.
func ensureArray(a []any, l int) []any {
	if len(a) >= l {
		return a
	}
	b := make([]any, l)
	copy(b, a)
	return b
}

// process processes node recursively.
func process(node any, parts []string, separator string, processor ValueProcessor) error {
	mkerr := func(err error, ps []string) error {
		return fmt.Errorf("cannot process '%s': %v", pathify(ps, separator), err)
	}
	// First check objects and arrays.
	if o, ok := isObject(node); ok {
		if len(o) == 0 {
			// Empty object.
			return processor(pathify(parts, separator), &Value{o, nil})
		}
		for field, subnode := range o {
			fieldparts := append(parts, field)
			if err := process(subnode, fieldparts, separator, processor); err != nil {
				return mkerr(err, parts)
			}
		}
		return nil
	}
	if a, ok := isArray(node); ok {
		if len(a) == 0 {
			// Empty array.
			return processor(pathify(parts, separator), &Value{a, nil})
		}
		for index, subnode := range a {
			indexparts := append(parts, strconv.Itoa(index))
			if err := process(subnode, indexparts, separator, processor); err != nil {
				return mkerr(err, parts)
			}
		}
		return nil
	}
	// Reached a value at the end.
	return processor(pathify(parts, separator), &Value{node, nil})
}

// pathify creates a path out of parts and separator.
func pathify(parts []string, separator string) string {
	return separator + strings.Join(parts, separator)
}

// EOF
