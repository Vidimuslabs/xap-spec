# Security Policy

Vidimus Labs takes the security of the Execution Authority Protocol (XAP) and its
reference materials seriously.

## Reporting a vulnerability

Please report suspected security vulnerabilities **privately** to
**security@vidimuslabs.com**. Do **not** open a public issue, pull request, or
discussion for a security report.

Where possible, include:

- the affected file(s) and the commit or `xap-1.0.0` schema element,
- a description of the issue and its security impact, and
- steps, inputs, or a proof-of-concept to reproduce it.

We aim to acknowledge reports within a few business days and will keep you
updated as we investigate. Please allow us reasonable time to remediate before
any public disclosure. Reporters who wish to be credited will be acknowledged.

## Scope

This repository holds the **normative specification, protocol constants, and
conformance vectors** — the verification-side reference. It contains no signing
keys and no enforcement code. Reports about canonicalization or signature
malleability, ambiguities that could weaken verification, incorrect or
under-constrained conformance vectors, or spec/schema issues with security
consequences are especially welcome. The reference verifier lives in
[xap-go](https://github.com/Vidimuslabs/xap-go).
