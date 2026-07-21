# XAP — Performance (non-normative)

This document reports measured performance characteristics of the **reference
implementation** (the `xap-go` verify SDK and the `xap-engine` enforcement
pipeline). It is **non-normative**: performance is a property of an
implementation and its host, not of the protocol. The protocol's only timing
requirement is the per-decision `latency_bound` constraint (§ SPEC 5, ¶0077),
which an enforcement point evaluates against its own configured budget.

## How to reproduce

**Conformance correctness (public).** The reference verifier replays every
conformance vector in this repository:

```
go run github.com/Vidimuslabs/xap-go/cmd/xap vectors run
```

**Verify-side latency (public).** The in-browser verifier at
[vidimuslabs.com/verify](https://www.vidimuslabs.com/verify) times a real hybrid
receipt on the visitor's own hardware — a self-measured verify latency you can
reproduce yourself.

**Issuance / full-pipeline latency (internal).** The end-to-end issuance and
signer figures below are benchmarked in the enforcement engine (`go test ./engine/
-bench=. -benchmem`), which is the licensed, private component. The signature and
digest algorithms are registry-selectable (SPEC § 11) and the pipeline is
algorithm-agnostic, so the Ed25519 → hybrid delta is exactly the added cost of the
post-quantum hybrid.

## Indicative results

Environment: Go 1.26, linux/amd64, single core, in-process **software** keys.
These are indicative figures from a development machine — **regenerate on
representative target hardware before citing any number publicly** (see
"Publishing" below).

**End-to-end and signature latency**

| Operation | Ed25519 (reference) | Hybrid (ECDSA-P384 + ML-DSA-65) |
|-----------|--------------------:|--------------------------------:|
| Engine receipt issuance (`EvaluateRequest`, full pipeline) | ~0.13 ms | **~1.4 ms** |
| Issue a signed MAT/receipt envelope (`Signer.Sign`) | ~0.03 ms | ~0.5 ms |
| Verify an envelope (`ParseMAT`, both-must-pass) | ~0.08 ms | ~0.8 ms |

**Primitive costs** (why the hybrid costs what it does)

| Primitive | sign | verify |
|-----------|-----:|-------:|
| Ed25519 | ~0.025 ms | ~0.059 ms |
| ML-DSA-65 (hedged) | ~0.19 ms | **~0.024 ms** |
| ECDSA P-384 | ~0.27 ms | **~0.74 ms** |

**Envelope size**

| | bytes |
|---|------:|
| Canonical MAT payload | ~456 |
| Ed25519 COSE_Sign1 envelope | ~535 |
| Hybrid COSE_Sign1 envelope | ~3881 (≈ 7.3×) |

The hybrid signature is a 96-byte ECDSA P-384 half followed by the ~3309-byte
ML-DSA-65 signature (SPEC § 11).

## Interpretation

- **Well within the latency bound.** A full hybrid issuance is ~1.4 ms against a
  typical `latency_bound` budget of 100–200 ms — roughly **1 % of budget**, ~700
  issuances/second/core. Post-quantum protection is effectively free on the
  decision path.
- **The post-quantum half is the cheaper half.** ML-DSA-65 verify (~0.024 ms) is
  *faster* than Ed25519 verify, and its hedged sign beats ECDSA P-384. The added
  cost of "going hybrid" is dominated by the **classical ECDSA P-384** half — in
  Go, P-384 uses generic field arithmetic with no hardware acceleration (unlike
  P-256). If decision latency ever became a constraint, the lever is the classical
  curve, not the post-quantum scheme.
- **Size.** The ~3.4 KB hybrid signature makes envelopes ~7× larger than Ed25519.
  XAP nodes run locally and receipts are not typically carried over
  bandwidth-constrained links, so this is a deliberate, accepted trade for
  quantum resistance.

## Production custody caveat

These figures use in-process software keys. In production the **classical (ECDSA
P-384) half is held in an HSM/KMS** behind the signer interface (split custody);
the device round-trip (typically single-digit milliseconds) then dominates
issuance latency — still comfortably within a 100 ms+ bound. The **ML-DSA-65 half
signs in software** and stays sub-millisecond. Signing uses **hedged**
(randomized) ML-DSA for fault/side-channel resistance; verification is identical
either way.

## Publishing

Before any public performance claim (website, datasheet):
1. Regenerate on representative production-class hardware, not a dev machine.
2. Include **HSM/KMS-inclusive** issuance numbers, or state clearly that the
   figures are in-process and that HSM custody adds device latency.
3. Publish the methodology (the command above) and figures as **ranges**, so a
   customer reproduces them rather than being surprised.
