# XAP `xap-1.0.0` — Frozen Schema

**Status: FROZEN at end of Session 1.** This document is the wire-schema contract
for protocol version `xap-1.0.0`. The engine and server build against exactly
these structures. Changes to a frozen field (rename, retype, reorder-semantics,
removal) require a new protocol version string; new *optional* fields may be added
within a major version only if they are `omitempty` and absent from all existing
canonical vectors' digests. This mirrors the frozen-schema discipline used for
prior protocol version strings.

Canonicalization: [RFC 8949](https://www.rfc-editor.org/rfc/rfc8949) §4.2 Core
Deterministic Encoding. Signature envelope: COSE_Sign1
([RFC 9052](https://www.rfc-editor.org/rfc/rfc9052)). Signature alg: Ed25519 reference, registry-agile — the registry also
admits `ecdsa-p256` and the post-quantum hybrid `hybrid-ecdsa-p384-ml-dsa-65`
(ECDSA P-384 + ML-DSA-65, both-must-pass, one COSE_Sign1 slot; see SPEC.md §11).
A registry addition does not change this version: the envelope and canonical
encoding are invariant across the algorithm choice. Digest: SHA-256
(registry-agile). CBOR keys are the `cbor:"..."` tags below; field declaration
order is **not** significant (deterministic encoding sorts map keys).

## Machine Authority Token (MAT) — ¶0041

```
v                string        protocol version, "xap-1.0.0"
id               string        artifact instance id
parent_id        string?       parent artifact id (derived MATs)
machine_identity MachineIdentity
scope            ExecutionScope
boundary         PermissionBoundary
trust_vector     TrustVector
proof_obligations []ProofObligation
constraints      []Constraint
delegation       DelegationRights
issuer           IssuerIdentity
replay           ReplayProtection
```

Sub-structures: `MachineIdentity{kind, public_key?, cert_ref?, attestation?,
composite?}`; `ExecutionScope{actions?, resources?, policy?, unconstrained?}`;
`PermissionBoundary{max_impact, max_privilege_delta, resource_quotas?,
exclusions?}`; `TrustVector{score?, level?}`; `ProofObligation{category,
max_age_seconds}`; `DelegationRights{allowed, max_depth}`; `IssuerIdentity{id,
kid?}`; `ReplayProtection{not_before, not_after, nonce, instance_id}`.

`ExecutionScope.unconstrained` names the dimensions a scope deliberately does
not restrict — `"actions"`, `"resources"`, or both. **An absent scope list
permits nothing; permitting a whole dimension requires naming it here.**

Absence is not a statement. Without this field an empty list has to mean either
"nothing is permitted" or "everything is permitted", and whichever is chosen the
other is what some issuer meant. Reading absence as "everything" makes the most
permissive grant in the protocol the one requiring the least typing, and leaves
an artifact that says nothing about a dimension indistinguishable from one that
deliberately opened it. Delegation carries the same rule: a child may declare a
dimension unconstrained only where its parent already did, and a child that
enumerates a dimension its parent left unconstrained is narrowing, which is
always permitted.

Optional and `omitempty`, so it is absent from every canonical vector digest
issued before it existed — what this document's own frozen-schema rule requires
of a within-version addition.

## Constraint (portable constraint representation) — ¶0087

```
id         string
type       string   temporal | network_zone | rate_limit | param_bound | resource_state | latency_bound
not_before,not_after string?   (temporal)
zones      []string?           (network_zone)
param      string?, min,max *float?  (param_bound)
key        string?, equals string?, in []string?  (resource_state)
max_rate   *int?               (rate_limit)
max_ms     int?                (latency_bound; evaluation budget, ¶0077)
```

Evaluation is deterministic and reproducible (¶0016): the same `(constraint,
context)` yields the same outcome on any implementation. Unknown types fail
closed.

## Runtime Context — ¶0047

```
time           string (RFC3339 UTC)
network_zone   string?
params         map[string]float64?
resource_state map[string]string?
risk_score     int?
rate           map[string]int64?   (keyed by rate_limit constraint id)
```

Context digest = SHA-256 over the canonical context (¶0018).

## Verifiable Execution Receipt — ¶0050/0088/0097

```
v                    string
id                   string
artifact_id          string
action               string?  the operation decided (¶0097 "reference to an
                              execution request or operation")
resource             string?  the operation's resource target
decision             string   permit | deny | permit_with_controls
controls             []string?
context_digest       bytes    (32, SHA-256)
rationale_codes      []string?
constraint_outcomes  []ConstraintOutcome?
evidence_refs        []EvidenceRef?
policy_digest        bytes?   (constraint-compilation digest, ¶0076)
timing               Timing
prior_hash           bytes?   (SHA-256 over prior receipt envelope, ¶0063)
resource_state_digest bytes?  (optimistic concurrency, ¶0054)
speculative          bool?    (¶0078)
enforcement_point    string?
commitment_digest    bytes?   (when a commitment governs, ¶0084A)
commitment_compliance ActionCompliance?  (¶0084A)
provenance           ProvenanceRef?      (multi-agent, ¶0084A)
```

Sub-structures: `Timing{start, complete, elapsed_ms, max_ms?}`;
`ConstraintOutcome{id, satisfied, code?}`; `EvidenceRef{category, digest,
fresh}`; `ActionCompliance{action, commitment_check, scope_check, boundary_check,
constraint_outcome, code?}`; `ProvenanceRef{parent_artifact_id,
parent_commitment_digest}`.

The enforcement point signature is the COSE_Sign1 envelope over the canonical
receipt payload.

`action` and `resource` are optional so that receipts predating them still decode
and so selective disclosure (¶0071, ¶0079) can withhold them. A verifier reports
the scope check as *not performed* rather than passed whenever the receipt does
not disclose enough to re-evaluate a dimension the MAT enumerates — **partial
disclosure is not a pass**. A receipt naming an action while withholding its
resource has had its resource evaluated against nothing, and calling that a
passing check asserts a gate never applied; pipeline step 2 (¶0046), whose
exceedance denies unconditionally regardless of every constraint outcome, would
otherwise be the one gate no independent party could reproduce.

`commitment_compliance` records what the enforcement point checked when a
commitment governs. Three of its booleans — `commitment_check`, `scope_check`,
`boundary_check` — are facts an independent party recomputes from the commitment
and the governing MAT, and **a verifier MUST recompute them rather than read them
back**. A structure whose whole purpose is to let someone check the enforcement
point's work cannot itself be taken on that point's word; otherwise a receipt may
claim `scope_check = true` for an action plainly outside its MAT's scope and
nothing contradicts it.

The comparison is **symmetric**, unlike the scope check above. A denial excuses
*acting* on an out-of-scope operation; it does not excuse *claiming* to have run
a check that would have returned the opposite answer. An assertion about a
reproducible fact is either correct or false, and the decision it accompanies
does not change which. `constraint_outcome` is the exception: it needs the
runtime context, and is reproduced by the constraint-outcome comparison of §9
when that context is supplied.

## Commitment Object — ¶0095B

```
v                 string
id                string
agent_identity    MachineIdentity
session_id        string
declared_actions  DeclaredActionSet
temporal_validity TemporalValidity
binding           CommitmentBinding
resource_targets  []string?
provenance        CommitmentProvenance?
action_window     TemporalValidity?
```

Sub-structures: `DeclaredActionSet{action_types, resources?, param_ranges?, unconstrained?}`;
`TemporalValidity{not_before, not_after}`; `CommitmentBinding{artifact_id,
constraint_digest}` — `constraint_digest` MUST equal the governing MAT's
constraint digest (Commitment Binding Verification, ¶0084A);
`CommitmentProvenance{parent_artifact_id, parent_commitment_digest}`.

The agent signature is the COSE_Sign1 envelope over the canonical commitment
payload.

`DeclaredActionSet.unconstrained` carries the same vocabulary and the same rule
as `ExecutionScope.unconstrained`: an absent list declares nothing, not
everything. A structure whose whole purpose is to be a *bounded* enumeration
must not become unbounded by omission.

**A commitment narrows the authority it binds to.** Every action and resource it
declares MUST be one the governing MAT authorizes, and it may not declare a
dimension unconstrained where that MAT enumerates it — otherwise the marker
would launder an escalation through the commitment layer instead of the
delegation layer. The same requirement binds `resource_targets` — a second resource list, separate
from `declared_actions.resources`, and every bit as much an over-claim surface —
and `param_ranges`: a declared range looser than the governing MAT's constraint
of the same id is an over-claim wearing the costume of a self-restriction, since
the commitment appears to bind the agent while relaxing what the MAT required. A
range with no counterpart in the MAT is an additional restriction the agent
chose and is always permitted.

A verifier reproduces this (`commitment_scope`), with the same
asymmetry as the scope check: an over-claiming commitment invalidates a receipt
that **permitted** under it, while a denial is the enforcement point doing
exactly what ¶0095A requires and remains consistent. Without this the declared
set is signed and carried but never evaluated, leaving
`COMMITMENT_SCOPE_VIOLATION` a code an enforcement point asserts and no
independent party can reach.

## Codes

The rationale/error/rejection code registry is frozen in
`constants` (see SPEC.md §10). Codes are append-only within a major version.
