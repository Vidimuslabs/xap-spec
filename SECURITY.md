# Security Policy

Vidimus Labs takes the security of the Execution Authority Protocol (XAP) and its
reference materials seriously.

## We are asking you to break this

XAP's entire claim is that a receipt can be verified by someone with no access to
the system that issued it. A claim like that is worth exactly what independent
scrutiny says it is worth, so we would rather you find a flaw than take our word
for it. Engineers, cryptographers, and academic researchers are explicitly
invited to attack this specification.

This repository is the design-level target. A protocol flaw is worth more to us
than an implementation bug, and it is far cheaper to fix now than after receipts
signed under the frozen `xap-1.0.0` schema are in the world — those are meant to
verify years later, so the wire format cannot be quietly revised. Concretely, we
would like to know if you can find:

- an **ambiguity in the canonicalization rules** that lets two distinct encodings
  present as one value, or one value canonicalize two ways;
- a **malleability** in the receipt, MAT, or commitment schema — any field an
  attacker can vary without changing what the signature covers;
- a **weakness in the hybrid composition** (ECDSA P-384 + ML-DSA-65, both-must-
  pass, private-use COSE alg `-65537`) — separability, cross-protocol reuse, or
  any way one half's failure does not deny;
- a **delegation invariant that does not hold** — a derivation the four monotonic
  rules (¶0057) permit but that widens authority;
- a **conformance vector that is wrong or under-constrained**, such that a
  non-conforming implementation still reproduces every expected outcome;
- a **gap between the specification text and what a verifier can actually
  reproduce** — anything an independent party is told to check but cannot.

The reference verifier is [xap-go](https://github.com/Vidimuslabs/xap-go);
you need **both** repos. From xap-go: `go test ./...` and
`go run ./cmd/xap vectors run` replay every vector here against it. Read
xap-go's `SECURITY.md` for the implementation attack surface and the full list
of adversarial tests already in tree.

**Do not probe live Vidimus hosts** (including `vidimuslabs.com/verify`). Safe
harbor does not cover production infrastructure — attack a local checkout.

## Reporting a vulnerability

Please report suspected security vulnerabilities **privately** to
**security@vidimuslabs.com**. Do **not** open a public issue, pull request, or
discussion for a security report.

Where possible, include:

- the affected file(s) and the commit or `xap-1.0.0` schema element,
- a description of the issue and its security impact, and
- steps, inputs, or a proof-of-concept to reproduce it.

We aim to acknowledge reports within a few business days and will keep you
updated as we investigate.

## Publication

**You may publish your findings.** We ask that you give us 90 days from the
acknowledgement of your report before doing so, or until a fix ships if that
comes sooner, and we will tell you as soon as a fix is out. If we have not fixed
an issue in 90 days, publish anyway — a deadline that moves whenever the vendor
finds it inconvenient is not a deadline, and research that cannot be published
is not research.

We will not ask you to sign a non-disclosure agreement as a condition of
reporting, and we will not treat a good-faith disclosure as a breach of any
term.

## Safe harbor

We consider good-faith security research on this specification and its reference
materials to be authorized conduct. We will not initiate legal action against
you, or refer you for prosecution, for research that stays within this policy —
specifically research that respects the scope below, avoids privacy violations
and service disruption, and does not access or modify data that is not yours.

If a third party brings action against you for research conducted in good faith
under this policy, we will make that authorization clear.

The specification content is CC-BY-4.0, so analysing, quoting, and republishing
it — including in a paper that demonstrates a flaw — is already permitted. XAP
is patent-pending; see [`LICENSE`](LICENSE) and [`NOTICE.md`](NOTICE.md) for the
licensing position.

## Scope

**In scope:** this repository — the **normative specification, protocol
constants, and conformance vectors**, the OpenAPI definition, and the protobuf
handshake contract. It contains no signing keys and no enforcement code. Reports
about canonicalization or signature malleability, ambiguities that could weaken
verification, incorrect or under-constrained conformance vectors, or spec and
schema issues with security consequences are especially welcome. The reference
verifier lives in [xap-go](https://github.com/Vidimuslabs/xap-go).

**Out of scope — do not test these (safe harbor does not cover them):**

- Any live Vidimus host: `vidimuslabs.com`, `api.vidimuslabs.com`,
  `*.vidimuslabs.com`, and DNS/edge configuration behind them (including the
  public `/verify` page — use a local xap-go / WASM build instead).
- The private enforcement engine and server (not published; not this invite).

Also out of scope: volumetric denial of service, social engineering of Vidimus
Labs staff or contractors, and physical attacks.

**Resource exhaustion — where the line sits, because the two sides of it look
alike.** A single input that makes a parse, decode or verify path allocate without
bound, recurse without bound, or hang **is in scope**, and we would like to hear
about it: that is a defect in the artifact, or an under-specification in this
document that permits one, and either is ours to fix. Sending enough well-formed
traffic to exhaust a node **is not** — availability under load is a property of how
a deployment is provisioned, and no bound inside a process answers a distributed
source.

That distinction has an owner, and the specification names it: the `servers` block
states that *"XAP is self-hosted: each operator runs their own instance and sets
the host. There is no Vidimus-hosted endpoint."* So availability belongs to the
operator's ingress — rate limiting, a WAF, L3/L4 filtering, redundant replicas.
The node's own in-process bounds, and what they deliberately do not cover, are
documented for operators alongside the server distribution.

## Recognition

Reporters who wish to be credited are acknowledged by name in the release notes
and, for findings that change the protocol or the schema, in the specification
itself. Where a finding warrants a CVE we will credit you in it. For substantive
protocol-level findings we are glad to offer formal acknowledgement in the next
revision of this specification and its Zenodo deposit, which is citable — tell us
how you would like to be named and whether you have an ORCID.
