# xap-client (TypeScript)

```bash
npm install && npm run build
```

```ts
import { Client } from "xap-client";

const c = new Client("http://localhost:8080/xap/v1");
console.log(await c.getAnchors());

const admin = new Client("http://localhost:8080/xap/v1", {
  token: "eval-admin-token-change-me",
});
console.log(await admin.verifyChain());
```

Uses global `fetch` (Node 18+). See the parent [README](../README.md).
