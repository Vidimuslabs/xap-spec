// Package docs embeds the normative specification prose and exposes the parts of
// it that an implementation can be held to mechanically.
//
// Most of SPEC.md is argument, and argument is for people. The check name
// registry of §9.2 is not: those names are on the wire, in the `name` field of
// every verification result and in the list a relying party consults to require
// a minimum before acting. A registry that only a reader compares to the code is
// a registry that drifts from it, and the drift is invisible precisely because
// both halves look right on their own.
//
// So the table is parsed rather than admired. Like the rest of this module the
// package is stdlib-only, and it scans the one row shape a GitHub-flavoured
// markdown table has rather than taking a markdown dependency.
package docs

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed SPEC.md
var spec []byte

// Spec returns the raw specification document.
func Spec() []byte { return spec }

// VerificationCheck is one row of the §9.2 check name registry.
type VerificationCheck struct {
	// Name is the identifier carried in a verification result.
	Name string
	// Subject is the artifact whose claim the check tests ("receipt", "MAT",
	// "commitment"). It decides whether the check is reported at all when the
	// artifact was not presented.
	Subject string
	// MayReportNotPerformed marks a check whose inputs may not permit reaching a
	// verdict. Every one must be pinned as not_performed by a conformance
	// vector.
	MayReportNotPerformed bool
}

const checkRegistryHeading = "### 9.2 Check name registry"

// VerificationChecks returns the §9.2 registry in document order.
//
// It errors rather than returning nothing when the heading or the table is
// absent. A conformance gate built on a scanner that quietly returns an empty
// registry does not fail — it passes vacuously, certifying agreement with a
// table it never found.
func VerificationChecks() ([]VerificationCheck, error) {
	body := string(spec)
	i := strings.Index(body, checkRegistryHeading)
	if i < 0 {
		return nil, fmt.Errorf("docs: SPEC.md has no %q heading", checkRegistryHeading)
	}

	var out []VerificationCheck
	seen := map[string]bool{}
	for _, ln := range strings.Split(body[i:], "\n") {
		trimmed := strings.TrimSpace(ln)
		// A new section ends the registry.
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cells := splitRow(trimmed)
		// name | subject | step | may-report-not-performed | tests
		if len(cells) != 5 {
			continue
		}
		name := strings.Trim(cells[0], "`")
		// Skip the header and the alignment row.
		if name == "Check" || strings.HasPrefix(name, "---") {
			continue
		}
		if name == "" || name != strings.Trim(cells[0], "`") || strings.ContainsAny(name, " |") {
			return nil, fmt.Errorf("docs: §9.2 row has an unusable check name %q", cells[0])
		}
		if seen[name] {
			return nil, fmt.Errorf("docs: §9.2 lists %q twice", name)
		}
		seen[name] = true
		out = append(out, VerificationCheck{
			Name:                  name,
			Subject:               cells[1],
			MayReportNotPerformed: strings.Contains(cells[3], "✓"),
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("docs: §9.2 declared no checks")
	}
	return out, nil
}

// splitRow splits a markdown table row into its cells, dropping the leading and
// trailing pipes.
func splitRow(row string) []string {
	parts := strings.Split(strings.Trim(row, "|"), "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
