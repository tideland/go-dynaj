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
)

//--------------------
// TREE FUNCTIONS
//--------------------

// insertValue recursively inserts a value at the end of the keys list.
func insertValue(element Element, keys Keys, value Value) (Element, error) {
	if len(keys) == 0 {
		return value, nil
	}

	switch tnode := element.(type) {
	case nil:
		return createValue(keys, value)
	case Object:
		return insertValueInObject(tnode, keys, value)
	case Array:
		return insertValueInArray(tnode, keys, value)
	default:
		return nil, fmt.Errorf("document is not a valid JSON structure")
	}
}

// createValue creates a value at the end of the keys list.
func createValue(keys Keys, value Value) (Element, error) {
	// Check if we are at the end of the keys list.
	if len(keys) == 0 {
		return value, nil
	}
	h, t := headTail(keys)
	// Check for array index first.
	index, ok := asIndex(h)
	if ok {
		// It's an array index.
		arr := make(Array, index+1)
		element, err := createValue(t, value)
		if err != nil {
			return nil, err
		}
		arr[index] = element
		return arr, nil
	}
	// It's an object key.
	obj := Object{h: nil}
	element, err := createValue(t, value)
	if err != nil {
		return nil, err
	}
	obj[h] = element
	return obj, nil
}

// insertValueInObject inserts a value in a JSON object at the end of the keys list.
func insertValueInObject(obj Object, keys Keys, value Value) (Element, error) {
	h, t := headTail(keys)
	// Create object if keys list has only one element.
	if len(t) == 0 {
		if isObjectOrArray(obj[h]) {
			return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
		}
		_, ok := asIndex(h)
		if ok {
			return nil, fmt.Errorf("cannot insert value at %v: index %q in object", keys, h)
		}
		obj[h] = value
		return obj, nil
	}
	// Insert value in element.
	element := obj[h]
	if isValue(element) {
		return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
	}
	newElement, err := insertValue(element, t, value)
	if err != nil {
		return nil, err
	}

	obj[h] = newElement
	return obj, nil
}

// insertValueInArray inserts a value in an array at a given path.
func insertValueInArray(arr Array, keys Keys, value Value) (Element, error) {
	h, t := headTail(keys)
	// Convert path head into index.
	index, ok := asIndex(h)
	switch {
	case !ok:
		return nil, fmt.Errorf("cannot insert value at %v: invalid index %q", keys, index)
	case index < 0:
		return nil, fmt.Errorf("cannot insert value at %v: negative index %d", keys, index)
	case index >= len(arr):
		tmp := make(Array, index+1)
		copy(tmp, arr)
		arr = tmp
	}
	// Insert value if last element in path.
	if len(t) == 0 {
		if isObjectOrArray(arr[index]) {
			return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
		}
		arr[index] = value
		return arr, nil
	}
	// Insert value in element.
	element := arr[index]
	if isValue(element) {
		return nil, fmt.Errorf("cannot insert value at %v: would corrupt document", keys)
	}
	newElement, err := insertValue(element, t, value)
	if err != nil {
		return nil, err
	}

	arr[index] = newElement
	return arr, nil
}

// deleteElement recursively deletes a element at the end of the keys list.
func deleteElement(element Element, keys Keys, deep bool) (Element, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	switch tnode := element.(type) {
	case nil:
		return nil, fmt.Errorf("cannot delete value at %v: invalid path", keys)
	case Object:
		return deleteElementInObject(tnode, keys, deep)
	case Array:
		return deleteElementInArray(tnode, keys, deep)
	default:
		return nil, fmt.Errorf("cannot delete value at %v: path too long", keys)
	}
}

// deleteElementInObject deletes a value in a JSON object at the end of the keys list.
func deleteElementInObject(obj Object, keys Keys, deep bool) (Element, error) {
	h, t := headTail(keys)
	// Delete object if keys list has only one element.
	if len(t) == 0 {
		if deep {
			// Delete all.
			delete(obj, h)
			return obj, nil
		}
		// Not deep, so delete only if value.
		if isObjectOrArray(obj[h]) {
			return nil, fmt.Errorf("cannot delete value at %v: is no value", keys)
		}
		delete(obj, h)
		return obj, nil
	}
	// Delete element.
	element := obj[h]
	newElement, err := deleteElement(element, t, deep)
	if err != nil {
		return nil, err
	}
	obj[h] = newElement
	return obj, nil
}

// deleteElementInArray deletes a value in a JSON array at a given path.
func deleteElementInArray(arr Array, keys Keys, deep bool) (Element, error) {
	h, t := headTail(keys)
	// Convert head in index.
	index, ok := asIndex(h)
	switch {
	case !ok:
		return nil, fmt.Errorf("cannot delete value at %v: invalid index %q", keys, h)
	case index < 0:
		return nil, fmt.Errorf("cannot delete value at %v: negative index %d", keys, index)
	case index >= len(arr):
		return nil, fmt.Errorf("cannot delete value at %v: index %d out of range", keys, index)
	}
	// Delete object if keys list has only one element.
	if len(t) == 0 {
		if deep {
			// Delete all.
			copy(arr[index:], arr[index+1:])
			arr = arr[:len(arr)-1]
		}
		// Not deep, so delete only if value.
		if isObjectOrArray(arr[index]) {
			return nil, fmt.Errorf("cannot delete value at %v: is no value", keys)
		}
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
		return arr, nil
	}
	// Delete value in element.
	element := arr[index]
	newElement, err := deleteElement(element, t, deep)
	if err != nil {
		return nil, err
	}
	arr[index] = newElement
	return arr, nil
}

// EOF
