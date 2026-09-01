# Contributing

xap-spec is the normative specification, constants, and conformance vectors for
XAP protocol version `xap-1.0.0`. The wire schema is **frozen**.

## Develop

```
go test ./...
```

Replay the vectors with the reference verifier:

```
go run github.com/Vidimuslabs/xap-go/cmd/xap vectors run
```

## Scope

- Spec text, constants, vectors, OpenAPI, and the handshake proto belong here.
- Cryptographic verification belongs in [xap-go](https://github.com/Vidimuslabs/xap-go).
- Changes that alter the frozen `xap-1.0.0` wire schema need a new protocol
  version, not a silent revision.

## Security

See [`SECURITY.md`](SECURITY.md). Report vulnerabilities to
**security@vidimuslabs.com**. Do not open a public issue for a security report.

## License

Specification content is CC-BY-4.0. See [`LICENSE`](LICENSE) and
[`NOTICE.md`](NOTICE.md).
