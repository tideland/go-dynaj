# Tideland Go Dynamic JSON

[![GitHub release](https://img.shields.io/github/release/tideland/go-dynaj.svg)](https://github.com/tideland/go-dynaj)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-dynaj/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-dynaj)](https://github.com/tideland/go-dynaj/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/dynaj?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/dynaj?tab=packages)
[![Workflow](https://img.shields.io/github/workflow/status/tideland/go-dynaj/Go)](https://github.com/tideland/go-dynaj/actions/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-dynaj)](https://goreportcard.com/report/tideland.dev/go/dynaj)

## Description

**Tideland Go Dynamic JSON** provides a simple dynamic handling of JSON 
documents. Values can be retrieved, set and added by paths like "foo/bar/3".
Methods provide typesafe access to the values as well as flat and deep 
processing.

```go
doc, err := dynaj.Unmarshal(myDoc)
if err != nil {
    ...
}
name := doc.ValueAt("/name").AsString("")
street := doc.ValueAt("/address/street").AsString("unknown")
```

Another way is to create an empty document with

```go
doc := dynaj.NewDocument()
```

Here as well as in parsed documents values can be set with

```go
err := doc.SetValueAt("/a/b/3/c", 4711)
```

Additionally values of the document can be processed recursively using

```go
err := doc.Root().Process(func(pv *dynaj.PathValue) error {
    ...
})
```

or from deeper nodes with `doc.ValueAt("/a/b/3").Process(...)`.
Additionally flat processing is possible with

```go
err := doc.ValueAt("/x/y/z").Range(func(pv *dynaj.PathValue) error {
    ...
})
````

To retrieve the differences between two documents the function
`dynaj.Compare()` can be used:

```go
diff, err := dynaj.Compare(firstDoc, secondDoc)
````

privides a `dynaj.Diff` instance which helps to compare individual
paths of the two document.

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://themue.dev)
