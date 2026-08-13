# XAP — Execution Authority Protocol

**Reference specification, protocol layer.** Version `xap-1.0.0`.

This document transcribes the protocol layer of the Autonomous Machine Identity
and Authority Protocol (AMIAP), also denominated the Execution Authority
Protocol (XAP). It is a normative-intent description of message structures,
processing sequences, and the verification algorithm. It is derived from and
subordinate to the patent specification of U.S. patent application 19/570,167,
as amended by the preliminary amendment filed in that application; paragraph
anchors (¶NNNN) throughout point to the governing text. Where this document and
the amended specification differ, the amended specification governs.

The protocol is independent of wire format and transport (¶0082, ¶0086). This
reference build uses deterministic CBOR ([RFC 8949](https://www.rfc-editor.org/rfc/rfc8949)
§4.2 Core Deterministic Encoding) for canonicalization, COSE_Sign1
([RFC 9052](https://www.rfc-editor.org/rfc/rfc9052)) for signature envelopes,
Ed25519 as the reference signature algorithm, and SHA-256 for digests. Signature and digest
algorithms are drawn from registries (§11, ¶0018/0066): the registries admit
alternatives — including a post-quantum hybrid — without changing the protocol
version, because the envelope shape (COSE_Sign1 over the canonical payload) is
invariant across the algorithm choice.

Canonicalization is a **verification** requirement and not only an encoding
convention. An implementation MUST reject input that decodes successfully but is
not in canonical form — map keys out of the sorted order RFC 8949 §4.2 requires,
or integers at non-minimal width — as well as the encodings that change meaning
(duplicate keys, indefinite-length items, invalid UTF-8). Ordering and width
change only the bytes, which is precisely why they matter: the protocol
identifies artifacts by digests over those bytes, and a digest is an identity
only when one value has exactly one encoding. Re-encoding a decoded value and
requiring byte equality is sufficient, since canonical encoding is a function.

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

The instance id in field 138 and the artifact instance identifier (¶0084
Authority Identifier Field) are **the same value**, and an implementation
**MUST** reject a MAT in which they differ. They were separately described and
never related, which left every consumer to choose which one identifies the
artifact — and both answers were defensible from the text. A MAT **MUST** also
carry a nonce: freshness needs state a verifier may not have, but an absent
nonce means there is no replay protection at all, which is a different and worse
thing.

Constraint ids **MUST** be present and unique within an artifact. They are how
outcomes are matched to constraints at verification and how parent and child
constraints are paired during derivation; two constraints sharing an id collapse
in both lookups, so an outcome recorded for that id marks both as covered while
only one was ever evaluated — and the issuer's ordering decides which becomes
invisible.

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
suspend); elapsed time recorded in every receipt.

**A timeout receipt records `elapsed_ms` equal to `max_ms`, never greater.**
Evaluation is abandoned *at* the bound, so the bound is reached and not
exceeded; the timeout is signalled by the `CONSTRAINT_EVALUATION_TIMEOUT`
rationale code and the disposition, never by `elapsed_ms` running past `max_ms`.
This is normative because the alternative is unverifiable: a verifier checks
`elapsed_ms <= max_ms` (§9 step 3), so a receipt recording elapsed beyond the
bound could never verify, and the protocol would be unable to express a truthful
timeout at all. A receipt whose `elapsed_ms` exceeds `max_ms` is therefore
malformed regardless of which codes it carries. Conformance vector:
`receipt_constraint_timeout`.

Race handling (¶0054–0055):
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

The chain link is the prior receipt's **digest** — SHA-256 over its canonical
payload — so an append-only log cannot delete or reorder receipts without
breaking the chain.

It is deliberately **not** a hash of the signed envelope. A COSE_Sign1 envelope
has no unique encoding for a given receipt: ECDSA admits both `(r, s)` and
`(r, n−s)`, and the COSE unprotected header bucket is by construction outside
the signature, so a third party holding no key material can produce a second,
byte-distinct envelope that verifies as the same receipt. A link over envelope
bytes therefore lets a log holder manufacture "the chain broke" about receipts
nobody tampered with, and cannot serve as a dedup or replay key. Restricting
signatures to low-`s` would not be sufficient on its own; the header bucket is
malleable independently of the signature algorithm.

The payload is what the signature covers and what canonicalization pins to one
encoding, which is what makes it an identity. Implementations **MUST** reject a
payload that decodes but is not in canonical form (§ canonicalization); a digest
identifies a value only when that value has exactly one encoding.

## 8. Delegation invariants (FIG. 5, ¶0057)

A child MAT is a valid derivation of a parent iff all four hold:

- (i) child scope ⊆ parent scope;
- (ii) child boundary ≤ parent boundary (more restrictive);
- (iii) each child constraint ≥ as strict as the corresponding parent constraint;
- (iv) child proof obligations ⊇ parent obligations.

Plus acyclicity (¶0073) and delegation depth (¶0041 field 134). A derived
artifact failing derivation proof validation is unconditionally rejected.

Two rules make invariant (i) mean what it says, and an implementation that omits
either satisfies the letter of "subset" while widening authority:

- **An absent scope dimension is *unconstrained*, not empty.** A scope list that
  is absent imposes no restriction on that dimension, so a child that simply
  omits a dimension its parent constrained is **broader** than its parent, not
  narrower. Treating the child's list as a set and testing membership reads that
  omission as the empty set — vacuously a subset — and admits an escalation. If
  the parent constrains a dimension, the child **must** constrain it too.
- **The same rule binds every dimension, not only scope.** A child may not drop
  a resource quota its parent set (an absent quota is an unstated one, not a
  smaller one), may not empty a constraint's value set while keeping its ID (the
  empty set is trivially a subset, so this would neutralise the constraint while
  satisfying invariant iii), and a root permitting delegation **must state a
  depth** — an unstated `max_depth` is not unlimited, or the shallowest possible
  statement would produce the deepest possible grant.
- **Resource containment is not a prefix test on an un-normalized string.** A
  parent pattern `svc/*` textually covers `svc/../db/main`, which resolves
  outside `svc` entirely, so a traversal segment converts a narrowing pattern
  into an escape. A resource containing a `..` path segment is covered by
  nothing, in derivation and in the scope check of §9 step 5 alike. Whole
  segments are matched, so ordinary names containing dots stay usable.

## 9. Verification algorithm (¶0095)

An independent verifier, with the receipt, optionally the governing MAT and
reproduced context, the prior receipt, and the governing commitment — and the
public trust anchors — performs:

1. Validate the receipt signature. An artifact verifies only against a trust
   anchor **registered for that artifact kind** — issuer for a MAT, enforcement
   point for a receipt, agent for a commitment object. The three signing roles
   are distinct (¶0041 field 136, ¶0050, ¶0095B) and a raw public key states no
   purpose of its own, so an anchor set that does not record the role trusts
   every key for everything, and an agent's key can mint the authority that
   agent operates under. An anchor naming no role signs nothing.
2. Check version, ternary decision validity, and that all rationale/error codes
   are registered (§10). Check that the receipt's declared controls match its
   decision: `permit_with_controls` **MUST** name at least one control, and no
   other decision may name any. `permit_with_controls` is the only decision that
   may permit an operation whose constraints did not all hold, and the controls
   are what compensate (¶0049); a receipt claiming the decision and naming no
   control claims the exemption without the thing that earns it.
3. Check evaluation timing. Three separate questions, and a verifier answers as
   many as its inputs allow:
   (a) elapsed is within the bound the **receipt** declares — self-consistency,
   since both values come from the artifact under judgement;
   (b) the receipt's timing agrees with itself — `complete` does not precede
   `start`, `elapsed_ms` is not negative, and `elapsed_ms` matches the
   `start`..`complete` window within the one second RFC3339 second-granularity
   truncation can account for. `elapsed_ms` is the value every latency gate is
   applied to, so a receipt understating it escapes a bound it visibly exceeded
   while carrying the two timestamps that refute the claim;
   (c) if the MAT is supplied, elapsed is within the bound the **MAT
   authorizes** — the strictest `latency_bound` constraint it states. A receipt
   declaring `max_ms = 0` (unbounded) under a MAT that states a bound **fails**,
   as does one declaring a bound wider than authorized. Report **not performed**
   where the MAT states no `latency_bound`. Delegation already refuses a child
   that widens this bound (§8); verification refusing a receipt that ignores it
   is the same rule at the other end.
   Where `start` and `complete` are withheld, (b) is **not performed**.
4. If the MAT is supplied: validate its signature and structure, and confirm the
   receipt's artifact id binds to it. Confirm the MAT's signed issuer key id is
   the key id that verified it — an artifact naming an issuer it was not issued
   by is not a well-formed artifact, and a MAT that names no issuing key id
   fails, since absence cannot excuse a binding check. A MAT must also carry a
   replay nonce and an instance id (field 138); their presence is the only part
   of replay protection a stateless verifier can establish, and their absence
   means there is no replay protection rather than replay protection it cannot
   evaluate — see §9.1.
   Confirm the receipt's evaluation `start` falls inside the MAT's validity
   interval (¶0065). Both values are signed, so this asks no clock anything and
   every verifier reaches the same answer; it is distinct from checking an
   artifact against *now*, which depends on the verifier's clock and is
   therefore a separate, caller-supplied lifecycle check.
5. If the MAT is supplied and the receipt names the operation: **recompute the
   scope and boundary check** (¶0046, pipeline step 2). Confirm the receipt's
   `action` and `resource` fall within the MAT's execution scope and permission
   boundary. A receipt recording any decision other than `deny` for an operation
   outside scope **fails verification** — exceedance denies unconditionally,
   regardless of every constraint outcome, so a receipt that permits an
   out-of-scope operation is self-inconsistent. A `deny` of an out-of-scope
   operation is consistent and passes. Report this check as **not performed**
   rather than passed whenever the receipt does not disclose enough to
   re-evaluate every dimension the MAT constrains — that is, when the MAT
   constrains actions and `action` is absent, or constrains resources and
   `resource` is absent, or both are absent (selective disclosure, ¶0071,
   ¶0079, or issuers predating those fields). **Partial disclosure is not a
   pass.** A receipt naming an action while withholding its resource has had its
   resource evaluated against nothing, and reporting that as a passing scope
   check asserts a gate that was never applied — which is exactly what turns the
   one unconditional gate into an unverifiable assertion.
6. If the reproduced context is supplied: recompute the context digest and
   compare; where the receipt carries a `resource_state_digest` **and** the
   `resource_keys` it covers, recompute that digest over the reproduced values
   of exactly those keys (¶0054) — a digest carried without its keys, or naming
   a key the reproduced context does not carry, is unreproducible and reports
   **not performed** rather than being recomputed over the subset that happens
   to be available; recompute each recorded constraint outcome and compare; check
   decision consistency. Report the outcome check as **not performed** where the
   receipt records no outcome for some constraint the MAT states — withholding
   is legitimate (¶0071, ¶0079) and is not the same claim as having been
   checked. A receipt recording no outcomes at all agrees with everything it was
   shown, having been shown nothing, and reporting that as passed asserts a gate
   that was never applied. Decision consistency: a `permit` requires every
   constraint to hold; a `deny` is consistent with any constraint state; a
   `permit_with_controls` is consistent exactly when controls are named.
7. If a prior receipt is supplied: confirm the chain link — the receipt's
   `prior_hash` equals the prior receipt's **digest** (§7), not a hash of its
   envelope.
8. If a commitment is supplied: validate its signature, verify its binding to the
   MAT, and confirm the receipt's commitment digest. Confirm that every operation
   the commitment declares is one the governing MAT authorizes — a commitment
   narrows the authority it binds to and may declare less than was granted, never
   more (¶0095A) — under the same asymmetry as step 5: an over-claiming
   commitment fails a receipt that **permitted** under it, while a `deny` is the
   enforcement point doing what ¶0095A requires and stays consistent.
   Where the governing MAT is **not supplied**, every check above that reproduces
   a claim against it — the binding, the scope narrowing, and the MAT-dependent
   compliance recomputations — is reported **not performed**. It is **NOT**
   omitted: an absent check and a passed check are indistinguishable in a result,
   so omission lets whoever presents the receipt delete a gate by withholding an
   input rather than fail it. A check reproducible from the commitment alone
   still runs; needing *some* artifact is not the same as needing *this* one.
   Confirm the receipt's
   evaluation `start` falls inside the commitment's `temporal_validity` and,
   where declared, its `action_window` (¶0095B), by the same signed-values
   reasoning as step 4; report **not performed** where the commitment declares
   no such window. Confirm the commitment's `agent_identity` is the machine
   identity the governing MAT authorizes (field 122): the MAT names who the
   authority is FOR and the commitment names who it was generated BY, and a
   disagreement means the commitment is not the one that artifact governs.
   Report **not performed** where either side discloses no identity.
9. A receipt marked `speculative` records an evaluation **pending confirmation**
   (¶0078) and is not a final authorization. Confirmation **MUST** issue a
   receipt: an ordinary chained, signed, non-speculative receipt naming the
   speculative one it settles in `confirms`. Without that, "pending
   confirmation" has nothing to be pending on, and the mode emits only artifacts
   no verifier will accept.
   A receipt is issued in **both** outcomes. Where the resource state still
   holds, the confirming receipt carries the original decision; where it moved,
   it carries `deny` with `SPECULATIVE_CONFIRMATION_FAILURE` (§10) — a race that
   voided a permit is precisely the event an audit log exists to carry, and
   reporting it only to the caller would leave the one interesting case
   unrecorded. A receipt that settles a speculative evaluation is never itself
   speculative, or confirmation could not terminate.
   Where the confirmed receipt is supplied, a verifier confirms `confirms` names
   its digest, that the confirmed receipt was speculative, and that this one is
   not. A verifier **MUST** surface the speculative distinction rather than
   reporting such a receipt as simply valid; the reference
   verifier fails an explicit finality check and carries the flag in its result,
   so a relying party that accepts speculative receipts does so deliberately.

Verification succeeds using only the receipt, reproduced inputs, and public
anchors — **never** enforcement point internal state (¶0017).

**Integrity proof validation is not reproducible, and a verifier must not report
it as verified.** Pipeline step 4 (§5, ¶0048) validates integrity evidence
against the MAT's proof obligations, freshness included, at execution time. From
the signed structures a verifier can reproduce two things: that every category in
the MAT's `proof_obligations` appears among the receipt's `evidence_refs`, and
that each such reference records `fresh = true`. It cannot reproduce the
freshness determination. `EvidenceRef` carries `{category, digest, fresh}` and no
timestamp, so `max_age_seconds` has no input to evaluate against — `fresh` is the
enforcement point's assertion about a liveness observation made at execution time
and not re-observable from the receipt afterwards. Neither can the evidence
itself be checked, the receipt carrying only a digest of it.

Implementations **SHOULD** perform the two reproducible checks; the reference
verifier does. Both follow the disclosure and asymmetry rules already established
for the scope check in step 5. A receipt disclosing **no** evidence references at
all has exercised selective disclosure (¶0071, ¶0079), and both checks are
reported **not performed** rather than failed; once it discloses any, it is
making a readable claim about coverage, and partial disclosure is not a pass. A
receipt that **denies** on uncovered or stale evidence is the enforcement point
doing what ¶0048 requires and is consistent; only a receipt recording a decision
other than `deny` while its own references leave an obliged category unreferenced,
or record `fresh = false`, **fails verification**.

What no implementation may do is present step 4 as independently verified. These
two checks establish that the receipt's account of its evidence is internally
consistent with the artifact governing it — not that the evidence was good. Of
the six pipeline stages this is the one an outside party takes on trust, and
saying so is the difference between a bounded claim and an overstated one.



### 9.1 Signed fields this protocol does NOT verify

A field being signed means the issuer committed to it. It does not mean a
verifier can check it, and the two are routinely confused, because a schema that
carries a field implies someone evaluates it. An implementation **MUST NOT**
report a field in this section as verified.

**Once a check's subject artifact is supplied, that check is reported, never
omitted.** Not performed, passed and failed are the three answers; silence is not
a fourth. A result that simply lacks an entry cannot be distinguished from one
that passed, so dropping an unrunnable check hands the presenting party a way to
remove a gate by withholding an input rather than fail it.

**The subject of a check is whoever made the claim it tests, not every artifact
the computation touches.** The scope check reads the governing MAT, but the claim
under examination — that this action was within authority — is the *receipt's*.
So it is the receipt's presence that obliges the check to be reported, and the
MAT is merely the input needed to reach a verdict on it. A receipt is always
supplied, therefore **every check testing a receipt-carried claim is reported for
every verification**, reading `not performed` where an input it needed was
absent. Where the MAT is not supplied that covers, at minimum, the artifact
binding, the scope check, evidence coverage and freshness, the authorized latency
bound, the policy digest, evaluation-within-validity, the reproduced constraint
outcomes and decision consistency; where the reproduced context is not supplied,
the context digest and the resource-state digest; where the prior receipt is not
supplied, the chain link.

This does not oblige an implementation to emit checks about artifacts nobody
supplied. A receipt presented alone makes no commitment claims, so §9 step 8's
checks have no subject and are absent rather than not-performed; likewise the
MAT's own signature is the MAT's claim about itself, and there is none to report
when no MAT was given. The rule bites where a claim IS presented and the input
needed to reproduce it is the one thing left out.

**NOT PERFORMED belongs to a verifier lacking inputs, never to an artifact
withholding what it was obliged to provide.** The distinction is load-bearing,
because the party who benefits from a check not running is often the party who
controls whether it can run — an agent signs its own commitment, an enforcement
point its own receipt. Where an artifact declines to supply something the
protocol requires of it, that **MUST** fail rather than report not performed;
otherwise the weaker answer is a downgrade the signer selects. Concretely: a
commitment governed by a MAT that names a machine identity must name one
(¶0095B), and a commitment must declare its temporal validity. Selective
disclosure (¶0071, ¶0079) is the legitimate case and stays not-performed — the
distinction is whether the omitted value was owed.

A verification result **SHOULD** name the checks it did not perform, so a relying
party can require a minimum before acting. `valid` answers "was anything
refuted", not "how much was established", and those diverge.

The section is deliberately short. Most of what once sat here was unverifiable
only because the protocol had not said what the field meant, or had not
disclosed the input needed to reproduce it — both of which are fixable, and were
fixed rather than documented:

- `policy_digest` is now defined over the governing MAT's **portable** constraint
  set (¶0087), not an engine's internal compiled form. An internal form is an
  implementation's own business, so a digest over it could never be compared;
  over the portable form it binds a receipt to the exact constraints it was
  evaluated under. Verified in step 4.
- `resource_state_digest` is now accompanied by `resource_keys`, naming the
  variables it covers. The digest was unreproducible only because the subset was
  undisclosed. Verified in step 6.
- `enforcement_point` is bound to the identity the OPERATOR registered for the
  signing key. A receipt cannot bind its own name — the enforcement point picks
  both the name it writes and the key it signs with, so comparing the two
  compares its own choices. Binding needs a statement from outside, and the
  anchor set is the operator's statement about keys. Where an anchor records no
  subject the name is unverifiable and reports NOT PERFORMED: there the missing
  input is the operator's, which is the distinction §9.1 draws below. Verified
  in step 2.
- Replay is checked against a record of receipts the relying party has already
  **acted on**, keyed on the receipt digest. ¶0017 forbids relying on
  ENFORCEMENT-POINT state; it says nothing about the verifier's own, and a party
  that remembers what it accepted learns nothing from the issuer by doing so. A
  relying party keeping no record reports this NOT PERFORMED. See step 2.

  The unit is the **receipt**, and the record is consulted but never updated by
  verification. Both matter. A MAT's nonce identifies the ARTIFACT, and one
  artifact authorizes many operations, so keying replay on it makes every
  receipt after the first under that MAT look like a replay — which breaks
  walking an append-only log, the ordinary auditing case. A commitment's
  `session_id` has the same shape: a session has many actions. And a verifier
  that recorded what it inspected would make verification a mutation, so looking
  at one receipt twice would report the second look as a replay. A replay guard
  is updated when a receipt is acted upon, which only the relying party knows.

What genuinely cannot be verified is a judgement, because there is no correct
value to recompute against:

- `trust_vector` (MAT field 128) is the issuer's assessment of the subject. Two
  issuers may score the same subject differently and neither is wrong.
- `scope.policy` is an opaque policy expression evaluated by the engine; the
  protocol assigns it no structure, so there is nothing to check it against.

Both are the issuer's word, and a verifier should present them as the issuer's
word. The test for membership in this section is whether an independent party,
given the artifacts and the public anchors, can reach a conclusion — not whether
anything currently checks it. Every field that failed that test for a fixable
reason has been fixed.

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
| `SPECULATIVE_CONFIRMATION_FAILURE` | a speculative evaluation was not confirmed: the resource state it read had moved by the time confirmation re-read it (¶0078) |

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
bytes each; the remainder is the ML-DSA-65 ([FIPS 204](https://csrc.nist.gov/pubs/fips/204/final)) signature over the same
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
generator signs the hybrid vectors deterministically ([RFC 6979](https://www.rfc-editor.org/rfc/rfc6979) for the ECDSA
half, deterministic ML-DSA) so every vector — classical and hybrid — reproduces
byte-for-byte.

Measured latency and size characteristics of the reference implementation
(non-normative) are in [`PERFORMANCE.md`](PERFORMANCE.md): a full hybrid issuance
runs in ~1 % of a typical `latency_bound` budget.

---

Zenodo paper: DOI [10.5281/zenodo.21144476](https://doi.org/10.5281/zenodo.21144476).
