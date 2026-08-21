# xap-spec

[![ci](https://img.shields.io/github/actions/workflow/status/Vidimuslabs/xap-spec/ci.yml?branch=main&style=flat-square&logo=github&logoColor=C9A961&labelColor=475569&label=ci)](https://github.com/Vidimuslabs/xap-spec/actions/workflows/ci.yml)
[![protocol](https://img.shields.io/badge/protocol-xap--1.0.0-2D5F4F?style=flat-square&labelColor=475569)](docs/SPEC.md)
[![wire schema](https://img.shields.io/badge/wire_schema-frozen-2D5F4F?style=flat-square&labelColor=475569)](docs/SCHEMA.md)
[![signatures](https://img.shields.io/badge/signatures-hybrid_post--quantum-2D5F4F?style=flat-square&labelColor=475569)](docs/SPEC.md)
[![DOI](https://img.shields.io/badge/DOI-10.5281%2Fzenodo.21144476-2D5F4F?style=flat-square&labelColor=475569)](https://doi.org/10.5281/zenodo.21144476)
[![spec license](https://img.shields.io/badge/spec-CC--BY--4.0-2D5F4F?style=flat-square&labelColor=475569)](LICENSE.md)

The normative specification, protocol constants, and canonical conformance
vectors for the **Execution Authority Protocol (XAP)**, protocol version
`xap-1.0.0`.

XAP governs execution authority: an enforcement point verifies a signed authority
artifact, evaluates constraints against runtime context *at execution time*, makes
a ternary decision (permit / deny / permit-with-controls), and emits a signed,
verifiable execution receipt that an independent third party can check by
recomputing digests — with no access to the enforcement point.

## Contents

- [`docs/SPEC.md`](docs/SPEC.md) — protocol-layer reference specification.
- [`docs/SCHEMA.md`](docs/SCHEMA.md) — the frozen `xap-1.0.0` wire schema.
- [`docs/PERFORMANCE.md`](docs/PERFORMANCE.md) — non-normative performance reference.
- [`constants/`](constants/) — protocol version, the ternary decision vocabulary,
  the rationale/error/rejection code registry, and the registered-algorithm
  tables. Consumed by the SDK and the engine.
- [`vectors/`](vectors/) — golden conformance vectors (CBOR hex + JSON contexts)
  with an expected-outcome [`manifest.json`](vectors/manifest.json). A conforming
  implementation reproduces every expected outcome.
- [`openapi/`](openapi/) — OpenAPI definition of the verification and issuance
  API surface.
- [`proto/`](proto/) — the transport-agnostic handshake contract (gRPC).
- [`clients/`](clients/) — multi-language **HTTP API clients** (Python, TypeScript)
  covering every OpenAPI `operationId`. Thin node drivers — cryptographic
  verification remains [xap-go](https://github.com/Vidimuslabs/xap-go).
- [`provenance/`](provenance/) — the image-signing public key and verification guide.

## Verify against the spec

The reference verifier is **[xap-go](https://github.com/Vidimuslabs/xap-go)**. It
replays every vector in this repository:

```
go run github.com/Vidimuslabs/xap-go/cmd/xap vectors run
```

## Cite

Paper: DOI [10.5281/zenodo.21144476](https://doi.org/10.5281/zenodo.21144476).
See [`CITATION.cff`](CITATION.cff).

## Status

Reference specification · protocol `xap-1.0.0` (frozen) · hybrid post-quantum ·
patent-pending.

## License

Specification content is licensed **CC-BY-4.0**. XAP is patent-pending. Vidimus
Labs **intends** to make implementation of XAP available under a **royalty-free
patent covenant** (final terms pending counsel review) — see
[`LICENSE.md`](LICENSE.md) and the public statement at
[vidimuslabs.com/ip](https://vidimuslabs.com/ip). Covered XAP patent rights are
as defined in the final covenant. The covenant does not extend to AGIV, CVEAR,
or AIRAP (separate provisional patent filings). It does not make Vidimus Labs
software open source. Until the covenant is published, website and LICENSE
language describe intent rather than an executed grant. No freedom-to-operate
clearance is given or implied.

---

Protected by U.S. Patent No. [PENDING ISSUANCE — 19/570,167] and pending applications. [www.vidimuslabs.com/patents](https://www.vidimuslabs.com/patents)
