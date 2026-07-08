# xap-spec

Reference specification, protocol constants, and canonical conformance vectors
for the **Execution Authority Protocol (XAP)**, protocol version `xap-1.0.0`.

- [`docs/SPEC.md`](docs/SPEC.md) — protocol-layer reference specification.
- [`docs/SCHEMA.md`](docs/SCHEMA.md) — the frozen `xap-1.0.0` wire schema.
- [`constants/`](constants/) — protocol version, ternary decision vocabulary,
  the rationale/error/rejection code registry, and the registered algorithm
  tables. Consumed by the SDK and the engine.
- [`vectors/`](vectors/) — golden conformance vectors (CBOR hex + JSON contexts)
  with an expected-outcome [`manifest.json`](vectors/manifest.json). A conforming
  implementation reproduces every expected outcome.
- [`openapi/`](openapi/) — OpenAPI definition of the verification and issuance
  API surface.

XAP governs execution authority: an enforcement point verifies a signed authority
artifact, evaluates constraints against runtime context at execution time, makes
a ternary decision, and emits a signed verifiable execution receipt that an
independent third party can verify by recomputing digests — with no access to the
enforcement point.

Paper: DOI [10.5281/zenodo.21144476](https://doi.org/10.5281/zenodo.21144476).

## License

See [`LICENSE.md`](LICENSE.md). License terms pending; no rights granted.

---

Protected by U.S. Patent No. [PENDING ISSUANCE — 19/570,167] and pending applications. www.vidimuslabs.com/patents
