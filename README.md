# xap-spec

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

See [`LICENSE.md`](LICENSE.md). License terms pending; no rights granted.

---

Protected by U.S. Patent No. [PENDING ISSUANCE — 19/570,167] and pending applications. [www.vidimuslabs.com/patents](https://www.vidimuslabs.com/patents)
