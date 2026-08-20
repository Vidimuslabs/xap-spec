# XAP multi-language clients (Phase 5)

Thin **HTTP API clients** for a self-hosted XAP node, generated-from-contract in
spirit: every `operationId` in
[`openapi/xap.yaml`](../openapi/xap.yaml) is wrapped once per language. The gRPC
handshake contract is [`proto/handshake.proto`](../proto/handshake.proto) —
generate stubs with your language’s `protoc` plugin (see below).

| Language | Package | Install |
|---|---|---|
| **Go** | [xap-go](https://github.com/Vidimuslabs/xap-go) | full **verify-only** SDK (crypto) — not a thin HTTP client |
| **Python** | `python/xap_client` | `pip install -e python/` |
| **TypeScript** | `typescript/` | `cd typescript && npm install && npm run build` |

## Scope (be honest)

These clients talk to the **node HTTP API** (`{scheme}://{host}/xap/v1`). They:

- call public routes (verify, anchors, digest, revocations, metrics, handshake)
- call admin routes when you pass a bearer token (issue, execute, chain, …)
- **do not** implement COSE / hybrid PQC verification locally

Independent cryptographic verification remains **xap-go** (and the WASM build).
Use these clients to *drive* a node; use xap-go to *trust* a receipt offline.

There is **no Vidimus-hosted endpoint**. You set the base URL of your node.

## Quick start

### Python

```bash
cd python && pip install -e .
python -c "
from xap_client import Client
c = Client('http://localhost:8080/xap/v1')
print(c.get_anchors())
"
```

With admin token:

```python
c = Client("http://localhost:8080/xap/v1", token="eval-admin-token-change-me")
c.execute(mat=envelope_hex, action="read", resource="/gw/payments/1",
          context={"time": "2026-08-20T12:00:00Z", "network_zone": "dmz"})
```

### TypeScript

```bash
cd typescript && npm install && npm run build
```

```ts
import { Client } from "xap-client";
const c = new Client("http://localhost:8080/xap/v1");
console.log(await c.getAnchors());
```

## Handshake

**HTTP** (under `/xap/v1/handshake/…`) is covered by both clients  
(`handshakeCapability`, `handshakeChallenge`, …).

**gRPC** — generate from the public proto (any language):

```bash
# Example: Go
protoc -I ../../proto \
  --go_out=. --go-grpc_out=. \
  ../../proto/handshake.proto

# Example: Python
python -m grpc_tools.protoc -I ../../proto \
  --python_out=. --grpc_python_out=. \
  ../../proto/handshake.proto
```

Server implementation stays private (`xap-engine`); only the contract is public.

## Contract coverage

```bash
./scripts/check-openapi-coverage.sh
```

Fails if an `operationId` in `xap.yaml` is missing from Python or TypeScript.

## Layout

```
xap-clients/
  python/xap_client/     # Client + types
  typescript/src/        # Client + types
  scripts/               # coverage gate
  examples/              # minimal scripts
```

## Status

Phase 5 multi-lang slice · protocol `xap-1.0.0` · self-hosted only · patent-pending.

---

Protected by U.S. Patent No. [PENDING ISSUANCE — 19/570,167] and pending applications.
www.vidimuslabs.com/patents
