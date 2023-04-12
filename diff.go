// Tideland Go Generic JSON Processor
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp // import "tideland.dev/go/gjp"

//--------------------
// DIFFERENCE
//--------------------

// Diff manages the two parsed documents and their differences.
type Diff struct {
	first  *Document
	second *Document
	paths  []string
}

// Compare parses and compares the documents and returns their differences.
func Compare(first, second []byte) (*Diff, error) {
	fd, err := Parse(first)
	if err != nil {
		return nil, err
	}
	sd, err := Parse(second)
	if err != nil {
		return nil, err
	}
	d := &Diff{
		first:  fd,
		second: sd,
	}
	err = d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CompareDocuments compares the documents and returns their differences.
func CompareDocuments(first, second *Document) (*Diff, error) {
	d := &Diff{
		first:  first,
		second: second,
	}
	err := d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// FirstDocument returns the first document passed to Diff().
func (d *Diff) FirstDocument() *Document {
	return d.first
}

// SecondDocument returns the second document passed to Diff().
func (d *Diff) SecondDocument() *Document {
	return d.second
}

// Differences returns a list of paths where the documents
// have different content.
func (d *Diff) Differences() []string {
	return d.paths
}

// DifferenceAt returns the differences at the given path by
// returning the first and the second value.
func (d *Diff) DifferenceAt(path string) (*PathValue, *PathValue) {
	fstPV := d.first.ValueAt(path)
	sndPV := d.second.ValueAt(path)
	return fstPV, sndPV
}

// compare iterates over the both documents looking for different
// values or even paths.
func (d *Diff) compare() error {
	firstPaths := map[string]struct{}{}
	firstProcessor := func(pv *PathValue) error {
		firstPaths[pv.path] = struct{}{}
		if !pv.Equals(d.second.ValueAt(pv.path)) {
			d.paths = append(d.paths, pv.path)
		}
		return nil
	}
	err := d.first.Process(firstProcessor)
	if err != nil {
		return err
	}
	secondProcessor := func(pv *PathValue) error {
		_, ok := firstPaths[pv.path]
		if ok {
			// Been there, done that.
			return nil
		}
		d.paths = append(d.paths, pv.path)
		return nil
	}
	return d.second.Process(secondProcessor)
}

// EOF
