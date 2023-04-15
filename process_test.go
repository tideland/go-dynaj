// Tideland Go Dynamic JSON - Unit Tests
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dynaj_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestProcess tests the processing of documents.
func TestProcess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	values := []string{}
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)

	// Verify iteration of all nodes.
	err = doc.Root().Process(processor)
	assert.NoError(err)
	assert.Length(values, 27)
	assert.Contains(`"/B/0/B" = "100"`, values)
	assert.Contains(`"/B/0/C" = "true"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verifiy processing error.
	processor = func(node *dynaj.Node) error {
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
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)

	// Verify iteration of all nodes.
	err = doc.NodeAt("/B/0/D").Process(processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	values = []string{}
	err = doc.NodeAt("/B/1").Process(processor)
	assert.NoError(err)
	assert.Length(values, 8)
	assert.Contains(`"/B/1/S/2" = "white"`, values)
	assert.Contains(`"/B/1/B" = "200"`, values)

	// Verifiy iteration of non-existing path.
	err = doc.NodeAt("/B/3").Process(processor)
	assert.ErrorContains(err, "invalid path")

	// Verify procesing error.
	processor = func(node *dynaj.Node) error {
		return errors.New("ouch")
	}
	err = doc.NodeAt("/A").Process(processor)
	assert.ErrorContains(err, "ouch")
}

// TestRange tests the range processing of documents.
func TestRange(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	values := []string{}
	processor := func(node *dynaj.Node) error {
		value := fmt.Sprintf("%q = %q", node.Path(), node.AsString("<undefined>"))
		values = append(values, value)
		return nil
	}
	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)

	// Verify range of object.
	values = []string{}
	err = doc.NodeAt("/B/0/D").Range(processor)
	assert.NoError(err)
	assert.Length(values, 2)
	assert.Contains(`"/B/0/D/A" = "Level Three - 0"`, values)
	assert.Contains(`"/B/0/D/B" = "10.1"`, values)

	// Verify range of array.
	values = []string{}
	err = doc.NodeAt("/B/1/S").Range(processor)
	assert.NoError(err)
	assert.Length(values, 3)
	assert.Contains(`"/B/1/S/0" = "orange"`, values)
	assert.Contains(`"/B/1/S/1" = "blue"`, values)
	assert.Contains(`"/B/1/S/2" = "white"`, values)

	// Verify range of value.
	values = []string{}
	err = doc.NodeAt("/A").Range(processor)
	assert.NoError(err)
	assert.Length(values, 1)
	assert.Contains(`"/A" = "Level One"`, values)

	// Verify range of non-existing path.
	err = doc.NodeAt("/B/0/D/X").Range(processor)
	assert.ErrorContains(err, "invalid path")

	// Verify range of mixed types.
	err = doc.NodeAt("/B/0").Range(processor)
	assert.ErrorContains(err, "is object or array")

	// Verify procesing error.
	processor = func(node *dynaj.Node) error {
		return errors.New("ouch")
	}
	err = doc.NodeAt("/A").Range(processor)
	assert.ErrorContains(err, "ouch")
}

// TestRootQuery tests querying a document.
func TestRootQuery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	nodes, err := doc.Root().Query("Z/*")
	assert.NoError(err)
	assert.Length(nodes, 0)
	nodes, err = doc.Root().Query("*")
	assert.NoError(err)
	assert.Length(nodes, 27)
	nodes, err = doc.Root().Query("/A")
	assert.NoError(err)
	assert.Length(nodes, 1)
	nodes, err = doc.Root().Query("/B/*")
	assert.NoError(err)
	assert.Length(nodes, 24)
	nodes, err = doc.Root().Query("/B/[01]/*")
	assert.NoError(err)
	assert.Length(nodes, 18)
	nodes, err = doc.Root().Query("/B/[01]/*A")
	assert.NoError(err)
	assert.Length(nodes, 4)
	nodes, err = doc.Root().Query("*/S/*")
	assert.NoError(err)
	assert.Length(nodes, 8)
	nodes, err = doc.Root().Query("*/S/3")
	assert.NoError(err)
	assert.Length(nodes, 1)

	// Verify the content
	nodes, err = doc.Root().Query("/A")
	assert.NoError(err)
	assert.Equal(nodes[0].Path(), "/A")
	assert.Equal(nodes[0].AsString(""), "Level One")
}

// TestValueAtQuery tests querying a document starting at a deeper node.
func TestValueAtQuery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	nodes, err := doc.NodeAt("/B/0/D").Query("Z/*")
	assert.NoError(err)
	assert.Length(nodes, 0)
	nodes, err = doc.NodeAt("/B/0/D").Query("*")
	assert.NoError(err)
	assert.Length(nodes, 2)
	nodes, err = doc.NodeAt("/B/0/D").Query("A")
	assert.NoError(err)
	assert.Length(nodes, 1)
	nodes, err = doc.NodeAt("/B/0/D").Query("B")
	assert.NoError(err)
	assert.Length(nodes, 1)
	nodes, err = doc.NodeAt("/B/0/D").Query("C")
	assert.NoError(err)
	assert.Length(nodes, 0)
	nodes, err = doc.NodeAt("/B/1").Query("S/*")
	assert.NoError(err)
	assert.Length(nodes, 3)
	nodes, err = doc.NodeAt("/B/1").Query("S/2")
	assert.NoError(err)
	assert.Length(nodes, 1)

	// Verify non-existing path.
	nodes, err = doc.NodeAt("Z/Z/Z").Query("/A")
	assert.ErrorContains(err, "invalid path")
	assert.Length(nodes, 0)
}

// EOF
