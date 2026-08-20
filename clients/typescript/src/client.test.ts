import { createServer, type IncomingMessage, type ServerResponse } from "node:http";
import { once } from "node:events";
import test from "node:test";
import assert from "node:assert/strict";
import { Client, XAPError } from "./client.js";

async function withServer(
  handler: (req: IncomingMessage, res: ServerResponse) => void,
  fn: (base: string) => Promise<void>,
): Promise<void> {
  const server = createServer(handler);
  server.listen(0, "127.0.0.1");
  await once(server, "listening");
  const addr = server.address();
  if (!addr || typeof addr === "string") throw new Error("no port");
  const base = `http://127.0.0.1:${addr.port}/xap/v1`;
  try {
    await fn(base);
  } finally {
    server.close();
  }
}

test("getAnchors", async () => {
  await withServer((req, res) => {
    assert.equal(req.url, "/xap/v1/anchors");
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify([{ kid_hex: "ab", alg: "ed25519" }]));
  }, async (base) => {
    const c = new Client(base);
    const a = await c.getAnchors();
    assert.equal(a[0]?.alg, "ed25519");
  });
});

test("verifyChain requires token", async () => {
  await withServer((_req, res) => {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end("{}");
  }, async (base) => {
    const c = new Client(base);
    await assert.rejects(() => c.verifyChain(), (e: unknown) => {
      assert.ok(e instanceof XAPError);
      assert.equal(e.status, 401);
      return true;
    });
  });
});

test("execute with token", async () => {
  await withServer((req, res) => {
    assert.equal(req.headers.authorization, "Bearer tok");
    let body = "";
    req.on("data", (c) => (body += c));
    req.on("end", () => {
      const j = JSON.parse(body) as { action: string };
      assert.equal(j.action, "read");
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ decision: "permit" }));
    });
  }, async (base) => {
    const c = new Client(base, { token: "tok" });
    const r = await c.executeRequest({
      mat: "ee",
      action: "read",
      resource: "/x",
      context: { time: "2026-01-01T00:00:00Z", network_zone: "dmz" },
    });
    assert.equal(r.decision, "permit");
  });
});
