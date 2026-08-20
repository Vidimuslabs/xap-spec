/** HTTP client for XAP OpenAPI operationIds (xap-spec/openapi/xap.yaml). */

export class XAPError extends Error {
  constructor(
    readonly status: number,
    readonly body: string,
    readonly path: string,
  ) {
    super(`XAP ${path}: HTTP ${status}: ${body.slice(0, 200)}`);
    this.name = "XAPError";
  }
}

export type ClientOptions = {
  token?: string;
  timeoutMs?: number;
  fetchImpl?: typeof fetch;
  userAgent?: string;
};

export type RuntimeContext = {
  time?: string;
  network_zone?: string;
  [key: string]: unknown;
};

export type ExecuteInput = {
  mat: string;
  action: string;
  resource: string;
  context: RuntimeContext;
  impact?: number;
  evidence?: Array<Record<string, unknown>>;
  receipt_id?: string;
  commitment?: string;
  resource_keys?: string[];
};

export class Client {
  readonly baseUrl: string;
  readonly token?: string;
  readonly timeoutMs: number;
  private readonly fetchImpl: typeof fetch;
  private readonly userAgent: string;

  constructor(baseUrl: string, opts: ClientOptions = {}) {
    this.baseUrl = baseUrl.replace(/\/$/, "");
    this.token = opts.token;
    this.timeoutMs = opts.timeoutMs ?? 30_000;
    this.fetchImpl = opts.fetchImpl ?? globalThis.fetch.bind(globalThis);
    this.userAgent = opts.userAgent ?? "xap-client-ts/0.1.0";
  }

  private async request(
    method: string,
    path: string,
    opts: {
      body?: unknown;
      query?: Record<string, string>;
      auth?: boolean;
      accept?: string;
    } = {},
  ): Promise<unknown> {
    let url = this.baseUrl + path;
    if (opts.query) {
      const q = new URLSearchParams(opts.query);
      url += `?${q.toString()}`;
    }
    const headers: Record<string, string> = {
      Accept: opts.accept ?? "application/json",
      "User-Agent": this.userAgent,
    };
    let body: string | undefined;
    if (opts.body !== undefined) {
      body = JSON.stringify(opts.body);
      headers["Content-Type"] = "application/json";
    }
    if (opts.auth || this.token) {
      if (!this.token) {
        throw new XAPError(401, "admin token required", path);
      }
      headers.Authorization = `Bearer ${this.token}`;
    }
    const ctrl = new AbortController();
    const t = setTimeout(() => ctrl.abort(), this.timeoutMs);
    try {
      const res = await this.fetchImpl(url, {
        method,
        headers,
        body,
        signal: ctrl.signal,
      });
      const text = await res.text();
      if (!res.ok) {
        throw new XAPError(res.status, text, path);
      }
      if (!text) return undefined;
      const ct = res.headers.get("content-type") ?? "";
      if (ct.includes("json") || text.startsWith("{") || text.startsWith("[")) {
        return JSON.parse(text) as unknown;
      }
      return text;
    } finally {
      clearTimeout(t);
    }
  }

  /** operationId: verifyReceipt */
  verifyReceipt(input: {
    receipt: string;
    mat?: string;
    context?: RuntimeContext;
    prior_receipt?: string;
    commitment?: string;
  }): Promise<Record<string, unknown>> {
    return this.request("POST", "/verify", { body: input }) as Promise<
      Record<string, unknown>
    >;
  }

  /** operationId: getAnchors */
  getAnchors(): Promise<Array<Record<string, unknown>>> {
    return this.request("GET", "/anchors") as Promise<
      Array<Record<string, unknown>>
    >;
  }

  /** operationId: computeDigest */
  computeDigest(context: RuntimeContext): Promise<{ digest_hex?: string }> {
    return this.request("POST", "/digest", { body: context }) as Promise<{
      digest_hex?: string;
    }>;
  }

  /** operationId: listRevocations */
  listRevocations(): Promise<Array<Record<string, unknown>>> {
    return this.request("GET", "/revocations") as Promise<
      Array<Record<string, unknown>>
    >;
  }

  /** operationId: getMetrics */
  getMetrics(): Promise<string> {
    return this.request("GET", "/metrics", {
      accept: "text/plain",
    }) as Promise<string>;
  }

  /** operationId: issueArtifact */
  issueArtifact(matHex: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/artifacts", {
      body: { mat_hex: matHex },
      auth: true,
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: deriveArtifact */
  deriveArtifact(
    parentHex: string,
    childHex: string,
  ): Promise<Record<string, unknown>> {
    return this.request("POST", "/artifacts/derive", {
      body: { parent_hex: parentHex, child_hex: childHex },
      auth: true,
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: revokeArtifact */
  revokeArtifact(
    artifactId: string,
    reason?: string,
  ): Promise<Record<string, unknown>> {
    return this.request("POST", `/artifacts/${encodeURIComponent(artifactId)}/revoke`, {
      query: reason ? { reason } : undefined,
      auth: true,
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: executeRequest */
  executeRequest(input: ExecuteInput): Promise<Record<string, unknown>> {
    return this.request("POST", "/execute", {
      body: input,
      auth: true,
    }) as Promise<Record<string, unknown>>;
  }

  /** Alias for application code. */
  execute(input: ExecuteInput): Promise<Record<string, unknown>> {
    return this.executeRequest(input);
  }

  /** operationId: listReceipts */
  listReceipts(): Promise<Array<Record<string, unknown>>> {
    return this.request("GET", "/receipts", { auth: true }) as Promise<
      Array<Record<string, unknown>>
    >;
  }

  /** operationId: verifyChain */
  verifyChain(): Promise<Record<string, unknown>> {
    return this.request("GET", "/chain/verify", { auth: true }) as Promise<
      Record<string, unknown>
    >;
  }

  /** operationId: handshakeCapability */
  handshakeCapability(versions: string[]): Promise<Record<string, unknown>> {
    return this.request("POST", "/handshake/capability", {
      body: { versions },
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeChallenge */
  handshakeChallenge(sessionId: string): Promise<Record<string, unknown>> {
    return this.request("POST", `/handshake/${encodeURIComponent(sessionId)}/challenge`, {
      body: {},
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeProof */
  handshakeProof(
    sessionId: string,
    nonce: string,
  ): Promise<Record<string, unknown>> {
    return this.request("POST", `/handshake/${encodeURIComponent(sessionId)}/proof`, {
      body: { nonce },
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeAuthority */
  handshakeAuthority(
    sessionId: string,
    mat: string,
  ): Promise<Record<string, unknown>> {
    return this.request(
      "POST",
      `/handshake/${encodeURIComponent(sessionId)}/authority`,
      { body: { mat } },
    ) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeNegotiate — constraints is canonical CBOR hex. */
  handshakeNegotiate(
    sessionId: string,
    constraints: string,
  ): Promise<Record<string, unknown>> {
    return this.request(
      "POST",
      `/handshake/${encodeURIComponent(sessionId)}/negotiate`,
      { body: { constraints } },
    ) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeBind */
  handshakeBind(sessionId: string): Promise<Record<string, unknown>> {
    return this.request("POST", `/handshake/${encodeURIComponent(sessionId)}/bind`, {
      body: {},
    }) as Promise<Record<string, unknown>>;
  }

  /** operationId: handshakeStatus */
  handshakeStatus(sessionId: string): Promise<Record<string, unknown>> {
    return this.request(
      "GET",
      `/handshake/${encodeURIComponent(sessionId)}`,
    ) as Promise<Record<string, unknown>>;
  }
}
