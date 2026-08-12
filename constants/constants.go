// Package constants defines the versioned protocol constants for the Execution
// Authority Protocol (XAP). It is the single source of truth for the protocol
// version string, the ternary decision vocabulary, the canonical rationale/
// error/rejection code registry, and the registered cryptographic algorithm
// tables. Both the verification SDK and the enforcement engine import this
// package so that codes appearing in a signed proof structure mean exactly one
// thing across every implementation.
//
// Spec anchors: protocol version negotiation (¶0083 addition, Capability
// Message); rationale/error/rejection code taxonomy (¶0084 addition); ternary
// decision (¶0049); algorithm agility for digests and signatures (¶0018,
// technology decision — registered-algorithms table).
package constants

// ProtocolVersion is the protocol version identifier embedded in every
// authority artifact, commitment object, and receipt. It is negotiated in the
// Capability Message (¶0083 addition — Protocol Version Field).
const ProtocolVersion = "xap-1.0.0"

// Decision is the ternary execution control decision emitted by the enforcement
// point (¶0049, Execution Decision Message ¶0083).
type Decision string

const (
	// DecisionPermit authorizes the operation without added controls.
	DecisionPermit Decision = "permit"
	// DecisionDeny refuses the operation. All unconditional-denial paths
	// (signature failure ¶0045, boundary exceedance ¶0046, revocation ¶0065)
	// resolve to this value.
	DecisionDeny Decision = "deny"
	// DecisionPermitWithControls authorizes the operation subject to one or
	// more applied execution controls (¶0049).
	DecisionPermitWithControls Decision = "permit_with_controls"
)

// Valid reports whether d is a recognized ternary decision value.
func (d Decision) Valid() bool {
	switch d {
	case DecisionPermit, DecisionDeny, DecisionPermitWithControls:
		return true
	default:
		return false
	}
}

// Canonical rationale, error, and rejection codes (¶0084 addition). Every code
// that appears in a receipt is bound by the enforcement point's signature and
// forms part of the independently verifiable proof record. The registry is
// versioned with the protocol: codes are append-only across a major version.
const (
	// CodeArtifactSignatureFailure: the authority artifact signature failed
	// verification against the configured trust anchor set (¶0045).
	CodeArtifactSignatureFailure = "ARTIFACT_SIGNATURE_FAILURE"
	// CodeArtifactScopeExceedance: the requested operation falls outside the
	// execution scope or exceeds the permission boundary (¶0046).
	CodeArtifactScopeExceedance = "ARTIFACT_SCOPE_EXCEEDANCE"
	// CodeConstraintEvaluationFailure: a specific constraint evaluated to false;
	// carried with the constraint ID and evaluation rationale (¶0047).
	CodeConstraintEvaluationFailure = "CONSTRAINT_EVALUATION_FAILURE"
	// CodeIntegrityEvidenceFailure: machine integrity evidence failed validation
	// against one or more integrity proof obligations (¶0048).
	CodeIntegrityEvidenceFailure = "INTEGRITY_EVIDENCE_FAILURE"
	// CodeConstraintEvaluationTimeout: constraint evaluation did not complete
	// within the maximum evaluation latency bound; carried with a timeout
	// handling disposition code (¶0052).
	CodeConstraintEvaluationTimeout = "CONSTRAINT_EVALUATION_TIMEOUT"
	// CodeCommitmentObjectSignatureFailure: the commitment object signature
	// failed verification against the agent identity (¶0095A).
	CodeCommitmentObjectSignatureFailure = "COMMITMENT_OBJECT_SIGNATURE_FAILURE"
	// CodeCommitmentScopeViolation: the commitment object's declared action set
	// exceeds the governing artifact's execution scope or permission boundary
	// (¶0095A, ¶0083 addition — Commitment Validation Message).
	CodeCommitmentScopeViolation = "COMMITMENT_SCOPE_VIOLATION"
	// CodeCommitmentActionViolation: a proposed action falls outside the
	// commitment object's declared action set (¶0083 addition — Commitment
	// Violation Message).
	CodeCommitmentActionViolation = "COMMITMENT_ACTION_VIOLATION"
	// CodeCommitmentRevocation: the governing commitment object has been revoked;
	// all further actions are blocked regardless of declared-set membership
	// (¶0083 addition — Commitment Revocation Message).
	CodeCommitmentRevocation = "COMMITMENT_REVOCATION"
	// CodeSpeculativeConfirmationFailure: a speculative evaluation (¶0078) was
	// not confirmed, because the resource state it read had changed by the time
	// confirmation re-read it. The speculative permit is void and this code
	// carries the reason.
	//
	// Distinct from CONSTRAINT_EVALUATION_FAILURE on purpose: that says a
	// constraint did not hold when it was evaluated, this says the ground moved
	// between evaluation and confirmation. A relying party reconstructing what
	// happened needs to tell those apart, and an audit protocol that collapses
	// them loses the only evidence that a race occurred.
	CodeSpeculativeConfirmationFailure = "SPECULATIVE_CONFIRMATION_FAILURE"
)

// Codes returns the registered rationale/error/rejection codes in registry
// order. Implementations use it to validate that a code carried in a receipt is
// a recognized member of the taxonomy for this protocol version.
func Codes() []string {
	return []string{
		CodeArtifactSignatureFailure,
		CodeArtifactScopeExceedance,
		CodeConstraintEvaluationFailure,
		CodeIntegrityEvidenceFailure,
		CodeConstraintEvaluationTimeout,
		CodeCommitmentObjectSignatureFailure,
		CodeCommitmentScopeViolation,
		CodeCommitmentActionViolation,
		CodeCommitmentRevocation,
		CodeSpeculativeConfirmationFailure,
	}
}

// KnownCode reports whether code is a registered member of the taxonomy.
func KnownCode(code string) bool {
	for _, c := range Codes() {
		if c == code {
			return true
		}
	}
	return false
}

// TimeoutDisposition is the configured timeout handling path recorded alongside
// CodeConstraintEvaluationTimeout (¶0052).
type TimeoutDisposition string

const (
	// TimeoutDegraded enters degraded mode with scope restriction (¶0052, ¶0064).
	TimeoutDegraded TimeoutDisposition = "degraded"
	// TimeoutDeny issues a denial with a timeout rationale code (¶0052).
	TimeoutDeny TimeoutDisposition = "deny"
	// TimeoutSuspend suspends execution pending recovery (¶0052).
	TimeoutSuspend TimeoutDisposition = "suspend"
)

// DigestAlg names a registered digest algorithm. SHA-256 is the reference
// algorithm; the table exists for algorithm agility (¶0018).
type DigestAlg string

// DigestSHA256 is the reference digest algorithm used wherever a "digest"
// appears in the protocol.
const DigestSHA256 DigestAlg = "sha-256"

// SignatureAlg names a registered signature algorithm. Ed25519 is the reference
// algorithm; the signer/verifier interface admits ECDSA P-256, a post-quantum
// hybrid, and HSM-backed signers (¶0066, FIG. 14). Adding a value is the
// registry-agility extension point and does not change the protocol version:
// the envelope stays COSE_Sign1 over the canonical payload and existing
// artifacts verify unchanged.
type SignatureAlg string

const (
	// SigEd25519 is the reference signature algorithm.
	SigEd25519 SignatureAlg = "ed25519"
	// SigECDSAP256 is a registered alternative (interface-reserved).
	SigECDSAP256 SignatureAlg = "ecdsa-p256"
	// SigHybridECDSAP384MLDSA65 is the post-quantum hybrid: a composite of
	// ECDSA P-384 (over SHA-384) and ML-DSA-65 (FIPS 204), verified
	// both-must-pass. It carries the whole portfolio's transition posture
	// (matching the CVEAR/AGIV/AIRAP receipts). The composite occupies the single
	// COSE_Sign1 signature slot as the 96-byte ECDSA r‖s followed by the ML-DSA-65
	// signature, under the private-use COSE algorithm id -65537 (private-use is
	// deliberate: the IETF composite-ML-DSA COSE code points are still draft).
	SigHybridECDSAP384MLDSA65 SignatureAlg = "hybrid-ecdsa-p384-ml-dsa-65"
)

// COSEAlgHybridECDSAP384MLDSA65 is the COSE algorithm identifier carried in the
// protected header of a SigHybridECDSAP384MLDSA65 envelope. It is drawn from the
// COSE private-use range (< -65536) so it can never be invalidated by a future
// IANA registration for composite ML-DSA. See SPEC.md §11.
const COSEAlgHybridECDSAP384MLDSA65 = -65537

// HybridECDSAP384SigLen is the fixed length, in bytes, of the ECDSA P-384 half
// (raw r‖s, 48 bytes each) that prefixes the composite hybrid signature. The
// remainder of the composite is the ML-DSA-65 signature.
const HybridECDSAP384SigLen = 96

// LifecycleState is an authority artifact lifecycle state (FIG. 13, ¶0065).
// Revoked and Expired are unconditionally rejected.
type LifecycleState string

const (
	StateIssued   LifecycleState = "issued"
	StateActive   LifecycleState = "active"
	StateNarrowed LifecycleState = "narrowed"
	StateRevoked  LifecycleState = "revoked"
	StateExpired  LifecycleState = "expired"
)
