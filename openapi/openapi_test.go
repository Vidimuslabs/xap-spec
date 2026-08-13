package openapi

import (
	"strings"
	"testing"
)

// TestOperationsRecoversTheDocument guards the scanner itself.
//
// Everything downstream of this package compares an implementation to what
// Operations returns, so a scanner that stopped recognising the document would
// not report a problem — it would report agreement, having found nothing to
// disagree with. The checks below are the properties that make the result
// trustworthy rather than merely non-empty.
func TestOperationsRecoversTheDocument(t *testing.T) {
	ops, err := Operations()
	if err != nil {
		t.Fatalf("Operations: %v", err)
	}

	// The verification surface is the part of this document that exists in
	// order to be reachable by parties holding no authority. If the scanner
	// cannot find it, it has not read the paths block.
	for _, want := range []string{"POST /verify", "GET /anchors", "POST /digest"} {
		if !contains(ops, want) {
			t.Errorf("the scanner did not recover %q from the document", want)
		}
	}

	seen := map[string]bool{}
	for _, op := range ops {
		if op.Method != strings.ToUpper(op.Method) || op.Method == "" {
			t.Errorf("operation %v has a malformed method", op)
		}
		if !strings.HasPrefix(op.Path, "/") {
			t.Errorf("operation %v has a path that is not rooted", op)
		}
		if seen[op.String()] {
			t.Errorf("%v recovered twice", op)
		}
		seen[op.String()] = true
	}
}

// TestBasePathIsRecovered guards the other half of the address.
//
// The paths block says which surfaces exist; the server URL says where they are.
// A parity check that compares only the first passes while every route sits at
// the wrong prefix, because each individual path still matches.
func TestBasePathIsRecovered(t *testing.T) {
	base, err := BasePath()
	if err != nil {
		t.Fatalf("BasePath: %v", err)
	}
	if base != "/xap/v1" {
		t.Errorf("base path = %q, want %q — if the document's server URL changed "+
			"deliberately, every implementation's mount point changes with it", base, "/xap/v1")
	}
	if strings.HasSuffix(base, "/") {
		t.Errorf("base path %q has a trailing slash; joining it to a path would double the separator", base)
	}
}

func contains(ops []Operation, want string) bool {
	for _, o := range ops {
		if o.String() == want {
			return true
		}
	}
	return false
}
