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
	"encoding/json"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/dynaj"
)

//--------------------
// TESTS
//--------------------

// TestBuilding tests the creation of documents.
func TestBuilding(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Most simple document.
	doc := dynaj.NewDocument()
	err := doc.SetValueAt("", "foo")
	assert.NoError(err)

	sv := doc.NodeAt("").AsString("bar")
	assert.Equal(sv, "foo")

	// Positive cases.
	doc = dynaj.NewDocument()
	err = doc.SetValueAt("/a/b/x", 1)
	assert.NoError(err)
	err = doc.SetValueAt("/a/b/y", true)
	assert.NoError(err)
	err = doc.SetValueAt("/a/c", "quick brown fox")
	assert.NoError(err)
	err = doc.SetValueAt("/a/d/0/z", 47.11)
	assert.NoError(err)
	err = doc.SetValueAt("/a/d/1/z", nil)
	assert.NoError(err)
	err = doc.SetValueAt("/a/d/2", 2)
	assert.NoError(err)

	iv := doc.NodeAt("a/b/x").AsInt(0)
	assert.Equal(iv, 1)
	bv := doc.NodeAt("a/b/y").AsBool(false)
	assert.Equal(bv, true)
	sv = doc.NodeAt("a/c").AsString("")
	assert.Equal(sv, "quick brown fox")
	fv := doc.NodeAt("a/d/0/z").AsFloat64(8.15)
	assert.Equal(fv, 47.11)
	nv := doc.NodeAt("a/d/1/z").IsUndefined()
	assert.True(nv)

	nodes, err := doc.Root().Query("*x")
	assert.NoError(err)
	assert.Length(nodes, 1)

	// Now provoke errors.
	err = doc.SetValueAt("/a/d", "stupid")
	assert.ErrorContains(err, "cannot insert value")
	err = doc.SetValueAt("/a/d/0", "stupid")
	assert.ErrorContains(err, "cannot insert value")
	err = doc.SetValueAt("/a/d/2/z", "stupid")
	assert.ErrorContains(err, "cannot insert value")
	err = doc.SetValueAt("/a/b/y/z", "stupid")
	assert.ErrorContains(err, "cannot insert value")
	err = doc.SetValueAt("a", "stupid")
	assert.ErrorMatch(err, ".*corrupt.*")
	err = doc.SetValueAt("a/b/x/y", "stupid")
	assert.ErrorMatch(err, ".*corrupt.*")
	err = doc.SetValueAt("/a/d/x", "stupid")
	assert.ErrorMatch(err, ".*invalid index.*")
	err = doc.SetValueAt("/a/d/-1", "stupid")
	assert.ErrorMatch(err, ".*negative index.*")

	// Legal change of values.
	err = doc.SetValueAt("/a/b/x", 2)
	assert.NoError(err)
	iv = doc.NodeAt("a/b/x").AsInt(0)
	assert.Equal(iv, 2)
}

// TestParseError tests the returned error in case of
// an invalid document.
func TestParseError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs := []byte(`abc{def`)

	doc, err := dynaj.Unmarshal(bs)
	assert.Nil(doc)
	assert.ErrorContains(err, "cannot unmarshal document")
}

// TestClear tests to clear a document.
func TestClear(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	doc.Clear()
	err = doc.SetValueAt("/", "foo")
	assert.NoError(err)
	foo := doc.NodeAt("/").AsString("<undefined>")
	assert.Equal(foo, "foo")
}

// TestLength tests retrieving values as strings.
func TestLength(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	l := doc.Length("X")
	assert.Equal(l, -1)
	l = doc.Length("")
	assert.Equal(l, 4)
	l = doc.Length("B")
	assert.Equal(l, 3)
	l = doc.Length("B/2")
	assert.Equal(l, 5)
	l = doc.Length("/B/2/D")
	assert.Equal(l, 2)
	l = doc.Length("/B/1/S")
	assert.Equal(l, 3)
	l = doc.Length("/B/1/S/0")
	assert.Equal(l, 1)
}

// TestNotFound tests the handling of not found values.
func TestNotFound(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)

	// Check if is undefined.
	node := doc.NodeAt("you-wont-find-me")
	assert.False(node.IsUndefined())
	assert.True(node.IsError())
	assert.ErrorContains(node.Err(), "invalid path")
}

// TestString verifies the string representation of a document.
func TestString(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	s := doc.String()
	assert.Equal(s, string(bs))
}

// TestAsString tests retrieving values as strings.
func TestAsString(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	sv := doc.NodeAt("A").AsString("default")
	assert.Equal(sv, "Level One")
	sv = doc.NodeAt("B/0/B").AsString("default")
	assert.Equal(sv, "100")
	sv = doc.NodeAt("B/0/C").AsString("default")
	assert.Equal(sv, "true")
	sv = doc.NodeAt("B/0/D/B").AsString("default")
	assert.Equal(sv, "10.1")
	sv = doc.NodeAt("Z/Z/Z").AsString("default")
	assert.Equal(sv, "default")

	sv = doc.NodeAt("A").String()
	assert.Equal(sv, "Level One")
	sv = doc.NodeAt("Z/Z/Z").String()

	// Difference between invalid path and nil value.
	assert.Contains("invalid path", sv)
	doc.SetValueAt("Z/Z/Z", nil)
	sv = doc.NodeAt("Z/Z/Z").String()
	assert.Equal(sv, "null")
}

// TestAsInt tests retrieving values as ints.
func TestAsInt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	iv := doc.NodeAt("A").AsInt(-1)
	assert.Equal(iv, -1)
	iv = doc.NodeAt("B/0/B").AsInt(-1)
	assert.Equal(iv, 100)
	iv = doc.NodeAt("B/0/C").AsInt(-1)
	assert.Equal(iv, 1)
	iv = doc.NodeAt("B/0/S/2").AsInt(-1)
	assert.Equal(iv, 1)
	iv = doc.NodeAt("B/0/D/B").AsInt(-1)
	assert.Equal(iv, 10)
	iv = doc.NodeAt("Z/Z/Z").AsInt(-1)
	assert.Equal(iv, -1)
}

// TestAsFloat64 tests retrieving values as float64.
func TestAsFloat64(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	fv := doc.NodeAt("A").AsFloat64(-1.0)
	assert.Equal(fv, -1.0)
	fv = doc.NodeAt("B/0/B").AsFloat64(-1.0)
	assert.Equal(fv, 100.0)
	fv = doc.NodeAt("B/1/B").AsFloat64(-1.0)
	assert.Equal(fv, 200.0)
	fv = doc.NodeAt("B/0/C").AsFloat64(-99)
	assert.Equal(fv, 1.0)
	fv = doc.NodeAt("B/0/S/3").AsFloat64(-1.0)
	assert.Equal(fv, 2.2)
	fv = doc.NodeAt("B/1/D/B").AsFloat64(-1.0)
	assert.Equal(fv, 20.2)
	fv = doc.NodeAt("Z/Z/Z").AsFloat64(-1.0)
	assert.Equal(fv, -1.0)
}

// TestAsBool tests retrieving values as bool.
func TestAsBool(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := dynaj.Unmarshal(bs)
	assert.NoError(err)
	bv := doc.NodeAt("A").AsBool(false)
	assert.Equal(bv, false)
	bv = doc.NodeAt("B/0/C").AsBool(false)
	assert.Equal(bv, true)
	bv = doc.NodeAt("B/0/S/0").AsBool(false)
	assert.Equal(bv, false)
	bv = doc.NodeAt("B/0/S/2").AsBool(false)
	assert.Equal(bv, true)
	bv = doc.NodeAt("B/0/S/4").AsBool(false)
	assert.Equal(bv, true)
	bv = doc.NodeAt("Z/Z/Z").AsBool(false)
	assert.Equal(bv, false)
}

// TestMarshalJSON tests building a JSON document again.
func TestMarshalJSON(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Compare input and output.
	bsIn, _ := createDocument(assert)
	parsedDoc, err := dynaj.Unmarshal(bsIn)
	assert.NoError(err)
	bsOut, err := parsedDoc.MarshalJSON()
	assert.NoError(err)
	assert.Equal(bsOut, bsIn)

	// Now doc one.
	doc := dynaj.NewDocument()
	err = doc.SetValueAt("/a/2/x", 1)
	assert.NoError(err)
	err = doc.SetValueAt("/a/4/y", true)
	assert.NoError(err)
	bsIn = []byte(`{"a":[null,null,{"x":1},null,{"y":true}]}`)
	bsOut, err = doc.MarshalJSON()
	assert.NoError(err)
	assert.Equal(bsOut, bsIn)
}

//--------------------
// HELPERS
//--------------------

type levelThree struct {
	A string
	B float64
}

type levelTwo struct {
	A string
	B int
	C bool
	D *levelThree
	S []string
}

type levelOne struct {
	A string
	B []*levelTwo
	D time.Duration
	T time.Time
}

func createDocument(assert *asserts.Asserts) ([]byte, *levelOne) {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"1",
					"2.2",
					"true",
				},
			},
			{
				A: "Level Two - 1",
				B: 200,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 20.2,
				},
				S: []string{
					"orange",
					"blue",
					"white",
				},
			},
			{
				A: "Level Two - 2",
				B: 300,
				C: true,
				D: &levelThree{
					A: "Level Three - 2",
					B: 30.3,
				},
			},
		},
		D: 5 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 30, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	assert.NoError(err)
	return bs, lo
}

func createCompareDocument(assert *asserts.Asserts) []byte {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"0",
					"2.2",
					"false",
				},
			},
			{
				A: "Level Two - 1",
				B: 300,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 99.9,
				},
				S: []string{
					"orange",
					"blue",
					"white",
					"red",
				},
			},
		},
		D: 10 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 59, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	assert.NoError(err)
	return bs
}

// EOF
