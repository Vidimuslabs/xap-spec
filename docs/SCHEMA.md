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
composite?}`; `ExecutionScope{actions?, resources?, policy?}`;
`PermissionBoundary{max_impact, max_privilege_delta, resource_quotas?,
exclusions?}`; `TrustVector{score?, level?}`; `ProofObligation{category,
max_age_seconds}`; `DelegationRights{allowed, max_depth}`; `IssuerIdentity{id,
kid?}`; `ReplayProtection{not_before, not_after, nonce, instance_id}`.

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

Sub-structures: `DeclaredActionSet{action_types, resources?, param_ranges?}`;
`TemporalValidity{not_before, not_after}`; `CommitmentBinding{artifact_id,
constraint_digest}` — `constraint_digest` MUST equal the governing MAT's
constraint digest (Commitment Binding Verification, ¶0084A);
`CommitmentProvenance{parent_artifact_id, parent_commitment_digest}`.

The agent signature is the COSE_Sign1 envelope over the canonical commitment
payload.

## Codes

The rationale/error/rejection code registry is frozen in
`constants` (see SPEC.md §10). Codes are append-only within a major version.
