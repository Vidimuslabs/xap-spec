// Package openapi embeds the OpenAPI definition of the XAP HTTP surface and
// exposes the operations it declares, so an implementation can be held to the
// document rather than asked to agree with it by hand.
//
// The document is normative and public: it is the contract a client is generated
// from. That makes an undeclared route and an undelivered path the same class of
// defect as a verification check that disagrees with the spec — and the same
// thing has to be true of both, that nobody has to remember to keep them level.
// A server serving a route this package does not list is serving something no
// generated client can reach; a path listed here that nothing serves is a 404
// with a schema.
//
// Like the rest of this module, it is stdlib-only. Operations are recovered by
// scanning the document's structure rather than by parsing YAML, because a
// dependency-free spec module is worth more than a general parser here: the
// scanner reads exactly the two line shapes OpenAPI requires of a paths block,
// and Operations reports an error rather than an empty set if the document ever
// stops having that shape.
package openapi

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed xap.yaml
var spec []byte

// Spec returns the raw OpenAPI document.
func Spec() []byte { return spec }

// Operation is one declared HTTP operation: an uppercase method and the path
// template it is declared under (e.g. "POST /artifacts/{id}/revoke").
type Operation struct {
	Method string
	Path   string
}

// String renders the operation as "METHOD /path", the same shape Go's
// http.ServeMux uses for a route pattern, so the two are directly comparable.
func (o Operation) String() string { return o.Method + " " + o.Path }

// httpMethods are the operation keys OpenAPI defines inside a path item.
// Anything else at that indent (parameters, summary, $ref) is not an operation.
var httpMethods = map[string]bool{
	"get": true, "put": true, "post": true, "delete": true,
	"options": true, "head": true, "patch": true, "trace": true,
}

// Operations returns every operation the document declares, sorted by path then
// method.
//
// It errors rather than returning nothing when the document has no paths block
// or the block yields no operations. A parity check built on a scanner that
// silently returns an empty set does not fail — it passes vacuously, certifying
// agreement with a document it never read.
func Operations() ([]Operation, error) {
	lines := strings.Split(string(spec), "\n")

	start := -1
	for i, ln := range lines {
		if ln == "paths:" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("openapi: no top-level paths block")
	}

	var (
		ops     []Operation
		current string
	)
	for _, ln := range lines[start:] {
		trimmed := strings.TrimSpace(ln)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(ln) - len(strings.TrimLeft(ln, " "))

		// A new top-level key ends the paths block.
		if indent == 0 {
			break
		}
		// Two-space indent under paths: a path item, keyed by the path template.
		if indent == 2 && strings.HasPrefix(trimmed, "/") {
			current = strings.TrimSuffix(trimmed, ":")
			continue
		}
		// Four-space indent under a path item: possibly a method.
		if indent == 4 && current != "" && strings.HasSuffix(trimmed, ":") {
			key := strings.TrimSuffix(trimmed, ":")
			if httpMethods[key] {
				ops = append(ops, Operation{Method: strings.ToUpper(key), Path: current})
			}
		}
	}
	if len(ops) == 0 {
		return nil, fmt.Errorf("openapi: paths block declared no operations")
	}

	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Path != ops[j].Path {
			return ops[i].Path < ops[j].Path
		}
		return ops[i].Method < ops[j].Method
	})
	return ops, nil
}
