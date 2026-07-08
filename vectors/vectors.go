// Package vectors embeds the canonical conformance test vectors for the
// Execution Authority Protocol and exposes them to any conformance runner.
// The vectors are static signed artifacts (COSE_Sign1 CBOR) plus a manifest
// describing, for each vector, the scenario kind, the input files, and the
// expected verification outcome. A conforming verifier — the reference SDK and,
// later, the engine — must reproduce every expected outcome.
//
// This package is verification-side and stdlib-only: it embeds bytes and parses
// the manifest. It contains no signing keys and no issuance logic. The vectors
// themselves are minted by a private generator using throwaway test keys; only
// the public halves of those keys appear here, as trust anchors.
package vectors

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed manifest.json data
var files embed.FS

// Manifest is the top-level conformance manifest.
type Manifest struct {
	// Version is the protocol version the vectors target.
	Version string `json:"version"`
	// Anchors are the trust anchors (public keys) needed to verify the vectors.
	Anchors []Anchor `json:"anchors"`
	// Vectors is the ordered set of conformance scenarios.
	Vectors []Vector `json:"vectors"`
}

// Anchor is a trust anchor: a key identifier, its signature algorithm, and the
// public key, all hex-encoded where binary.
type Anchor struct {
	KIDHex string `json:"kid_hex"`
	Alg    string `json:"alg"`
	PubHex string `json:"pub_hex"`
}

// Vector describes one conformance scenario. Kind selects how a runner
// interprets the remaining fields.
type Vector struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	// Expect is "valid" when the scenario must verify/accept, "invalid" when it
	// must be rejected.
	Expect string `json:"expect"`
	// ExpectCode, when set, is the rationale/error code the receipt must carry
	// (e.g. COMMITMENT_ACTION_VIOLATION).
	ExpectCode string `json:"expect_code,omitempty"`
	// File references and scenario parameters (relative to the data directory).
	ReceiptFile      string   `json:"receipt_file,omitempty"`
	MATFile          string   `json:"mat_file,omitempty"`
	ParentMATFile    string   `json:"parent_mat_file,omitempty"`
	CommitmentFile   string   `json:"commitment_file,omitempty"`
	ContextFile      string   `json:"context_file,omitempty"`
	PriorReceiptFile string   `json:"prior_receipt_file,omitempty"`
	ReceiptFiles     []string `json:"receipt_files,omitempty"`
	// ExpectDigestHex is the expected canonical digest for "canon" vectors.
	ExpectDigestHex string `json:"expect_digest_hex,omitempty"`
	// At is an RFC3339 evaluation instant for lifecycle checks.
	At string `json:"at,omitempty"`
}

// Load parses and returns the embedded manifest.
func Load() (*Manifest, error) {
	b, err := files.ReadFile("manifest.json")
	if err != nil {
		return nil, fmt.Errorf("vectors: read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("vectors: parse manifest: %w", err)
	}
	return &m, nil
}

// File returns the raw bytes of a named data file (a COSE envelope, a context
// JSON, etc.), relative to the vectors data directory.
func File(name string) ([]byte, error) {
	b, err := files.ReadFile("data/" + name)
	if err != nil {
		return nil, fmt.Errorf("vectors: read %q: %w", name, err)
	}
	return b, nil
}

// DataFS returns the embedded data directory as an fs.FS.
func DataFS() (fs.FS, error) {
	return fs.Sub(files, "data")
}
