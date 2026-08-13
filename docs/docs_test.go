package docs

import (
	"testing"
)

// TestVerificationChecksRecoversTheRegistry guards the scanner the check-name
// gate is built on.
//
// The gate in xap-go compares its verifier's vocabulary to what this function
// returns. A scanner that stopped finding the table would therefore report
// agreement rather than a problem, which is the failure mode the registry exists
// to remove — so the properties that make the result meaningful are asserted
// here, in the module that owns the document.
func TestVerificationChecksRecoversTheRegistry(t *testing.T) {
	checks, err := VerificationChecks()
	if err != nil {
		t.Fatalf("VerificationChecks: %v", err)
	}

	// Three anchors chosen because they sit at different points of the table and
	// carry different columns: one that always reaches a verdict, one that may
	// not, and one whose subject is an artifact other than the receipt.
	want := map[string]VerificationCheck{
		"receipt_signature": {Subject: "receipt", MayReportNotPerformed: false},
		"scope_check":       {Subject: "receipt", MayReportNotPerformed: true},
		"commitment_scope":  {Subject: "commitment", MayReportNotPerformed: true},
	}
	got := map[string]VerificationCheck{}
	for _, c := range checks {
		got[c.Name] = c
	}
	for name, w := range want {
		g, ok := got[name]
		if !ok {
			t.Errorf("the scanner did not recover %q from §9.2", name)
			continue
		}
		if g.Subject != w.Subject {
			t.Errorf("%s: subject %q, want %q", name, g.Subject, w.Subject)
		}
		if g.MayReportNotPerformed != w.MayReportNotPerformed {
			t.Errorf("%s: may-report-not-performed %v, want %v", name, g.MayReportNotPerformed, w.MayReportNotPerformed)
		}
	}

	// Every row must carry a subject: it is the column that decides whether a
	// check is reported at all when its artifact was not presented, so a blank
	// one is not a cosmetic gap.
	for _, c := range checks {
		if c.Subject == "" {
			t.Errorf("check %q declares no subject artifact", c.Name)
		}
	}
}
