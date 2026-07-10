# XAP — Execution Authority Protocol

**Reference specification, protocol layer.** Version `xap-1.0.0`.

This document transcribes the protocol layer of the Autonomous Machine Identity
and Authority Protocol (AMIAP), also denominated the Execution Authority
Protocol (XAP). It is a normative-intent description of message structures,
processing sequences, and the verification algorithm. It is derived from and
subordinate to the patent specification (`AMIAP_Specification.docx` as amended by
`AMIAP_PRELIMINARY_AMENDMENT.docx`); paragraph anchors (¶NNNN) throughout point
to the governing text. Where this document and the amended specification differ,
the amended specification governs.

The protocol is independent of wire format and transport (¶0082, ¶0086). This
reference build uses deterministic CBOR (RFC 8949 §4.2 Core Deterministic
Encoding) for canonicalization, COSE_Sign1 for signature envelopes, Ed25519 as
the reference signature algorithm, and SHA-256 for digests. Signature and digest
algorithms are drawn from registries (§11, ¶0018/0066): the registries admit
alternatives — including a post-quantum hybrid — without changing the protocol
version, because the envelope shape (COSE_Sign1 over the canonical payload) is
invariant across the algorithm choice.

---

## 1. Protocol category

XAP defines the **Execution Authority Protocol** category (¶0020): the machine
governance of whether a *specific execution operation* may proceed *at a specific
instant*, under verified integrity, within scope and boundary, under reproducible
constraint evaluation, producing a cryptographic proof structure that an
independent party can verify by recomputing digests — with no access to
enforcement point state (¶0017). This category is disjoint from authentication,
authorization delegation, and channel security.

## 2. Authority artifact — Machine Authority Token (MAT)

A MAT (FIG. 2, ¶0041) is a signed structure with nine semantic fields:

| # | Field | Meaning |
|---|-------|---------|
| 122 | machine identity | public key / cert ref / attestation-bound / composite |
| 124 | execution scope | bounded enumeration of permitted operations + resource patterns |
| 126 | permission boundary | max impact, max privilege delta, resource quotas, exclusion lists — a hard ceiling |
| 128 | trust vector | quantitative/qualitative trust assessment |
| 130 | integrity proof obligations | evidence categories + freshness windows |
| 132 | execution constraints | portable, runtime-evaluable conditions (¶0087) |
| 134 | delegation rights | permission + depth |
| 136 | issuer identity + signature | issuer id; signature over the canonical serialization |
| 138 | replay protection + expiry | validity interval, nonce, instance id |

The issuer signature is carried by the COSE_Sign1 envelope over the canonical
MAT payload. The **constraint digest** — SHA-256 over the canonical constraint
set — is the value a commitment object binds to (§6, ¶0084A).

## 3. Message taxonomy (¶0083, as amended)

Core: Capability (now carrying a **protocol version identifier** for negotiation,
¶0083 addition), Challenge (freshness element), Proof (integrity evidence bound
to the freshness element), Authority (signed MAT), Constraint Negotiation
(subject to the monotonic rule), Execution Request, Execution Decision (permit /
deny / permit-with-controls), Receipt (the cryptographic proof structure).

Commitment layer (¶0083 addition): Commitment Object, Commitment Validation,
Commitment Violation, Commitment Revocation.

## 4. Handshake (FIG. 3, ¶0043)

capability (with version negotiation) → challenge (session-binding nonce) →
proof (bound to the challenge) → authority conveyance → constraint negotiation
under the **monotonic rule** (no message may relax a constraint or expand a
boundary) → receipt binding.

## 5. Enforcement pipeline (FIG. 4, ¶0044–0050)

The order is normative, and the unconditional-denial semantics are load-bearing:

1. **Signature verification** (¶0045). Failure → **unconditional deny**.
2. **Scope & boundary check** (¶0046). Exceedance → **unconditional deny,
   regardless of constraint outcome**.
3. **Constraint evaluation** (¶0047). Runtime context obtained **at execution
   time** (not session time). One rationale code per constraint with its binary
   outcome. Identical inputs → identical outcome (¶0016).
4. **Integrity proof validation** (¶0048). Freshness evaluated at execution time.
   Failure → deny **or degraded mode** (never fail-open, ¶0064).
5. **Ternary decision + controls** (¶0049): throttle, sandbox, canary, step-up
   proof (incl. risk-triggered, ¶0072), parameter redaction, post-action
   verification, rollback conditioning, co-signature.
6. **Receipt generation** (¶0050). Every request — including denials — produces a
   receipt.

Latency-bounded evaluation (FIG. 4A, ¶0051–0053): a maximum evaluation latency
bound; a per-constraint-type timeout path (degraded / deny-with-timeout /
suspend); elapsed time recorded in every receipt. Race handling (¶0054–0055):
optimistic concurrency by default (resource state digest in the receipt), with
pessimistic-lock and speculative interfaces; conflict resolution serializes,
backs off, or denies both — all with rationale codes.

## 6. Commitment layer (§VIII, ¶0060 replaced; ¶0084A; ¶0095A/B)

An autonomous agent issues a **commitment object** before a session (¶0095B): an
agent identity, a session id, a **declared action set** (bounded action types,
resource targets, parameter ranges), a temporal validity, and a **commitment
binding** naming the governing MAT and carrying the **constraint digest** of that
MAT's constraint set. Optional: parameter constraints, resource target set,
**multi-agent provenance** (parent artifact id + parent commitment digest), and a
temporal action window.

Resource-controller flow (¶0095A): MAT signature verify (deny on fail) →
commitment signature verify (deny on fail) → **commitment scope validation**
(declared set within MAT scope and boundary; exceedance = unconditional rejection
with `COMMITMENT_SCOPE_VIOLATION`) → per proposed action: commitment evaluation
(action vs. declared set) + scope/boundary check + runtime context **at proposal
time** + constraint evaluation + integrity validation → execution control →
**commitment proof structure** (a receipt carrying the commitment digest and a
per-action **commitment compliance field**).

**Commitment Binding Verification** (¶0084A): the constraint digest in the
commitment must equal a fresh digest of the governing MAT's constraint set;
mismatch is rejected. **Commitment revocation** by commitment digest
unconditionally blocks all further actions under that commitment regardless of
declared-set membership (¶0083 Commitment Revocation Message).

## 7. Receipt — verifiable execution receipt (¶0050, ¶0088, ¶0097 as amended)

Fields: authority artifact id; ternary decision (+ applied controls); runtime
context digest; rationale code set; constraint outcomes; evidence reference set;
policy/constraint-compilation digest (¶0076); evaluation timing (start /
complete / elapsed / bound, ¶0053); prior-receipt hash (chaining, FIG. 11);
resource state digest (¶0054); speculative flag (¶0078); enforcement point
signature (in the envelope). **When a commitment governs:** commitment digest and
a commitment compliance field (¶0084A). Compact profile (¶0080) and selective
disclosure (¶0071/¶0079) are supported.

The chain link is SHA-256 over the prior receipt's signed envelope; an
append-only log cannot delete or reorder receipts without breaking the chain.

## 8. Delegation invariants (FIG. 5, ¶0057)

A child MAT is a valid derivation of a parent iff all four hold:

- (i) child scope ⊆ parent scope;
- (ii) child boundary ≤ parent boundary (more restrictive);
- (iii) each child constraint ≥ as strict as the corresponding parent constraint;
- (iv) child proof obligations ⊇ parent obligations.

Plus acyclicity (¶0073) and delegation depth (¶0041 field 134). A derived
artifact failing derivation proof validation is unconditionally rejected.

## 9. Verification algorithm (¶0095)

An independent verifier, with the receipt, optionally the governing MAT and
reproduced context, the prior receipt, and the governing commitment — and the
public trust anchors — performs:

1. Validate the receipt signature.
2. Check version, ternary decision validity, and that all rationale/error codes
   are registered (§10).
3. Check evaluation timing is within the recorded bound.
4. If the MAT is supplied: validate its signature and structure, and confirm the
   receipt's artifact id binds to it.
5. If the reproduced context is supplied: recompute the context digest and
   compare; recompute each recorded constraint outcome and compare; check
   decision consistency.
6. If a prior receipt is supplied: confirm the chain link.
7. If a commitment is supplied: validate its signature, verify its binding to the
   MAT, and confirm the receipt's commitment digest.

Verification succeeds using only the receipt, reproduced inputs, and public
anchors — **never** enforcement point internal state (¶0017).

## 10. Rationale / error / rejection code registry (¶0084 addition)

All codes appearing in a receipt are signature-bound and part of the verifiable
record.

| Code | Meaning |
|------|---------|
| `ARTIFACT_SIGNATURE_FAILURE` | authority artifact signature failed (¶0045) |
| `ARTIFACT_SCOPE_EXCEEDANCE` | outside scope or over boundary (¶0046) |
| `CONSTRAINT_EVALUATION_FAILURE` | a constraint evaluated false (+ id + rationale) (¶0047) |
| `INTEGRITY_EVIDENCE_FAILURE` | integrity evidence failed obligations (¶0048) |
| `CONSTRAINT_EVALUATION_TIMEOUT` | evaluation exceeded the latency bound (+ disposition) (¶0052) |
| `COMMITMENT_OBJECT_SIGNATURE_FAILURE` | commitment signature failed (¶0095A) |
| `COMMITMENT_SCOPE_VIOLATION` | declared set exceeds MAT scope/boundary (¶0095A) |
| `COMMITMENT_ACTION_VIOLATION` | proposed action outside declared set (¶0083) |
| `COMMITMENT_REVOCATION` | commitment revoked; all actions blocked (¶0083) |

## 11. Registered cryptographic algorithms (¶0018, ¶0066)

Digest and signature algorithms are named by registry so an implementation can
recompute digests and verify signatures deterministically. Registries are
append-only within a protocol version; adding an algorithm is not a version
change, because the canonical encoding and the COSE_Sign1 envelope shape are
invariant across the choice.

**Digest algorithms.** `sha-256` (reference), used wherever a digest appears.

**Signature algorithms.**

| Name | COSE alg | Signature slot | Notes |
|------|---------|----------------|-------|
| `ed25519` | `-8` (EdDSA) | 64-byte Ed25519 | reference algorithm (¶0066) |
| `ecdsa-p256` | `-7` (ES256) | ASN.1 ECDSA | registered alternative; HSM/KMS-friendly |
| `hybrid-ecdsa-p384-ml-dsa-65` | `-65537` (private use) | `ECDSA-P384(SHA-384) r‖s` (96 B) `‖ ML-DSA-65` | post-quantum hybrid, **both-must-pass** |

The **hybrid** carries the classical and post-quantum signatures in the single
COSE_Sign1 signature slot as a fixed-layout composite: the first 96 bytes are the
ECDSA P-384 signature (over SHA-384 of the canonical payload) as raw `r‖s`, 48
bytes each; the remainder is the ML-DSA-65 (FIPS 204) signature over the same
payload. A verifier splits at 96 bytes and accepts only if **both** halves verify
against the anchor's ECDSA P-384 and ML-DSA-65 public keys — an attacker must
forge both schemes. The COSE algorithm id is taken from the private-use range
(`< -65536`) so it cannot collide with, or be invalidated by, a future IANA
registration for composite ML-DSA (still an IETF draft). The hybrid gives XAP the
same post-quantum posture as the CVEAR/AGIV/AIRAP authority artifacts. Classical
and post-quantum private keys may be held under independent custody (e.g. the
ECDSA half in an HSM, the ML-DSA half in software).

## 12. Conformance

The `vectors/` directory holds golden vectors with an expected-outcome manifest.
A conforming implementation reproduces every expected outcome. The reference SDK
(`xap-go`) does so via `xap vectors run`; engine-issued receipts verify against
the same SDK — the two-implementation cross-check for independent verifiability.
Hybrid anchors carry two public keys (ECDSA P-384 SPKI + ML-DSA-65); the
generator signs the hybrid vectors deterministically (RFC 6979 for the ECDSA
half, deterministic ML-DSA) so every vector — classical and hybrid — reproduces
byte-for-byte.

---

Zenodo paper: DOI [10.5281/zenodo.21144476](https://doi.org/10.5281/zenodo.21144476).
