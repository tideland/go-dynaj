// Tideland Go Generic JSON Processor - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/gjp"
)

//--------------------
// TESTS
//--------------------

// TestProcess tests the processing of documents.
func TestProcess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	values := []string{}
	processor := func(pv *gjp.PathValue) error {
		value := fmt.Sprintf("%q = %q", pv.Path(), pv.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := gjp.Unmarshal(bs)
	assert.NoError(err)

	// Verify iteration of all nodes.
	err = doc.Root().Process(processor)
	assert.NoError(err)
	assert.Length(values, 27)
	assert.Contains(`"/B/0/B" = "100"`, values)
	assert.Contains(`"/B/0/C" = "true"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verifiy processing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.Root().Process(processor)
	assert.ErrorContains(err, "ouch")
}

// TestValueAtProcess tests the processing of documents starting at a
// deeper node.
func TestValueAtProcess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	values := []string{}
	processor := func(pv *gjp.PathValue) error {
		value := fmt.Sprintf("%q = %q", pv.Path(), pv.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := gjp.Unmarshal(bs)
	assert.NoError(err)

	// Verify iteration of all nodes.
	err = doc.ValueAt("/B/0/D").Process(processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	values = []string{}
	err = doc.ValueAt("/B/1").Process(processor)
	assert.NoError(err)
	assert.Length(values, 8)
	assert.Contains(`"/B/1/S/2" = "white"`, values)
	assert.Contains(`"/B/1/B" = "200"`, values)

	// Verifiy iteration of non-existing path.
	err = doc.ValueAt("/B/3").Process(processor)
	assert.ErrorContains(err, "invalid path")

	// Verify procesing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.ValueAt("/A").Process(processor)
	assert.ErrorContains(err, "ouch")
}

// TestRange tests the range processing of documents.
func TestRange(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	values := []string{}
	processor := func(pv *gjp.PathValue) error {
		value := fmt.Sprintf("%q = %q", pv.Path(), pv.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := gjp.Unmarshal(bs)
	assert.NoError(err)

	// Verify range of object.
	values = []string{}
	err = doc.ValueAt("/B/0/D").Range(processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	// Verify range of array.
	values = []string{}
	err = doc.ValueAt("/B/1/S").Range(processor)
	assert.NoError(err)
	assert.Length(values, 3)
	assert.Contains(`"/B/1/S/0" = "orange"`, values)
	assert.Contains(`"/B/1/S/1" = "blue"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verify range of value.
	values = []string{}
	err = doc.ValueAt("/A").Range(processor)
	assert.NoError(err)
	assert.Length(values, 1)
	assert.Contains(`"/A" = "Level One"`, values)

	// Verify range of non-existing path.
	err = doc.ValueAt("/B/0/D/X").Range(processor)
	assert.ErrorContains(err, "invalid path")

	// Verify range of mixed types.
	err = doc.ValueAt("/B/0").Range(processor)
	assert.ErrorContains(err, "is object or array")

	// Verify procesing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.ValueAt("/A").Range(processor)
	assert.ErrorContains(err, "ouch")
}

// TestRootQuery tests querying a document.
func TestRootQuery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Unmarshal(bs)
	assert.NoError(err)
	pvs, err := doc.Root().Query("Z/*")
	assert.NoError(err)
	assert.Length(pvs, 0)
	pvs, err = doc.Root().Query("*")
	assert.NoError(err)
	assert.Length(pvs, 27)
	pvs, err = doc.Root().Query("/A")
	assert.NoError(err)
	assert.Length(pvs, 1)
	pvs, err = doc.Root().Query("/B/*")
	assert.NoError(err)
	assert.Length(pvs, 24)
	pvs, err = doc.Root().Query("/B/[01]/*")
	assert.NoError(err)
	assert.Length(pvs, 18)
	pvs, err = doc.Root().Query("/B/[01]/*A")
	assert.NoError(err)
	assert.Length(pvs, 4)
	pvs, err = doc.Root().Query("*/S/*")
	assert.NoError(err)
	assert.Length(pvs, 8)
	pvs, err = doc.Root().Query("*/S/3")
	assert.NoError(err)
	assert.Length(pvs, 1)

	// Verify the content
	pvs, err = doc.Root().Query("/A")
	assert.NoError(err)
	assert.Equal(pvs[0].Path(), "/A")
	assert.Equal(pvs[0].AsString(""), "Level One")
}

// TestValueAtQuery tests querying a document starting at a deeper node.
func TestValueAtQuery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Unmarshal(bs)
	assert.NoError(err)
	pvs, err := doc.ValueAt("/B/0/D").Query("Z/*")
	assert.NoError(err)
	assert.Length(pvs, 0)
	pvs, err = doc.ValueAt("/B/0/D").Query("*")
	assert.NoError(err)
	assert.Length(pvs, 2)
	pvs, err = doc.ValueAt("/B/0/D").Query("A")
	assert.NoError(err)
	assert.Length(pvs, 1)
	pvs, err = doc.ValueAt("/B/0/D").Query("B")
	assert.NoError(err)
	assert.Length(pvs, 1)
	pvs, err = doc.ValueAt("/B/0/D").Query("C")
	assert.NoError(err)
	assert.Length(pvs, 0)
	pvs, err = doc.ValueAt("/B/1").Query("S/*")
	assert.NoError(err)
	assert.Length(pvs, 3)
	pvs, err = doc.ValueAt("/B/1").Query("S/2")
	assert.NoError(err)
	assert.Length(pvs, 1)

	// Verify non-existing path.
	pvs, err = doc.ValueAt("Z/Z/Z").Query("/A")
	assert.ErrorContains(err, "invalid path")
	assert.Length(pvs, 0)
}

// EOF
