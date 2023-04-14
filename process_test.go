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
	err = doc.Process(processor)
	assert.NoError(err)
	assert.Length(values, 27)
	assert.Contains(`"/B/0/B" = "100"`, values)
	assert.Contains(`"/B/0/C" = "true"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verifiy processing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.Process(processor)
	assert.ErrorContains(err, "ouch")
}

// TestProcessPath tests the processing of documents starting at a path.
func TestProcessPath(t *testing.T) {
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
	err = doc.ProcessPath("/B/0/D", processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	values = []string{}
	err = doc.ProcessPath("/B/1", processor)
	assert.NoError(err)
	assert.Length(values, 8)
	assert.Contains(`"/B/1/S/2" = "white"`, values)
	assert.Contains(`"/B/1/B" = "200"`, values)

	// Verifiy iteration of non-existing path.
	err = doc.ProcessPath("/B/3", processor)
	assert.ErrorContains(err, "invalid path")

	// Verify procesing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.ProcessPath("/A", processor)
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
	err = doc.Range("/B/0/D", processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	// Verify range of array.
	values = []string{}
	err = doc.Range("/B/1/S", processor)
	assert.NoError(err)
	assert.Length(values, 3)
	assert.Contains(`"/B/1/S/0" = "orange"`, values)
	assert.Contains(`"/B/1/S/1" = "blue"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verify range of value.
	values = []string{}
	err = doc.Range("/A", processor)
	assert.NoError(err)
	assert.Length(values, 1)
	assert.Contains(`"/A" = "Level One"`, values)

	// Verify range of non-existing path.
	err = doc.Range("/B/0/D/X", processor)
	assert.ErrorContains(err, "invalid path")

	// Verify range of mixed types.
	err = doc.Range("/B/0", processor)
	assert.ErrorContains(err, "is object or array")

	// Verify procesing error.
	processor = func(pv *gjp.PathValue) error {
		return errors.New("ouch")
	}
	err = doc.Range("/A", processor)
	assert.ErrorContains(err, "ouch")
}

// EOF
