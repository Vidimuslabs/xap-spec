"""HTTP client for XAP OpenAPI operationIds (xap-spec/openapi/xap.yaml)."""

from __future__ import annotations

import json
import urllib.error
import urllib.request
from typing import Any, Mapping, MutableMapping, Optional, Sequence
from urllib.parse import urlencode


class XAPError(Exception):
    """Non-2xx response from the node."""

    def __init__(self, status: int, body: str, path: str) -> None:
        self.status = status
        self.body = body
        self.path = path
        super().__init__(f"XAP {path}: HTTP {status}: {body[:200]}")


class Client:
    """
    Self-hosted XAP node client.

    base_url must include the /xap/v1 prefix, e.g. http://localhost:8080/xap/v1.
    token is the admin bearer for privileged routes; omit for public-only use.
    """

    def __init__(
        self,
        base_url: str,
        token: Optional[str] = None,
        timeout: float = 30.0,
        user_agent: str = "xap-client-python/0.1.0",
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.token = token
        self.timeout = timeout
        self.user_agent = user_agent

    # ---- transport ----

    def _request(
        self,
        method: str,
        path: str,
        *,
        body: Any = None,
        query: Optional[Mapping[str, str]] = None,
        auth: bool = False,
        accept: str = "application/json",
    ) -> Any:
        url = self.base_url + path
        if query:
            url += "?" + urlencode(query)
        data = None
        headers: MutableMapping[str, str] = {
            "Accept": accept,
            "User-Agent": self.user_agent,
        }
        if body is not None:
            data = json.dumps(body).encode("utf-8")
            headers["Content-Type"] = "application/json"
        if auth or self.token:
            if not self.token:
                raise XAPError(401, "admin token required", path)
            headers["Authorization"] = f"Bearer {self.token}"
        req = urllib.request.Request(url, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read()
                if not raw:
                    return None
                ctype = resp.headers.get("Content-Type", "")
                if "json" in ctype or raw[:1] in (b"{", b"["):
                    return json.loads(raw.decode("utf-8"))
                return raw.decode("utf-8")
        except urllib.error.HTTPError as e:
            err_body = e.read().decode("utf-8", errors="replace")
            raise XAPError(e.code, err_body, path) from None

    # ---- public (no auth) ----

    def verify_receipt(
        self,
        receipt: str,
        *,
        mat: Optional[str] = None,
        context: Optional[Mapping[str, Any]] = None,
        prior_receipt: Optional[str] = None,
        commitment: Optional[str] = None,
    ) -> dict[str, Any]:
        """operationId: verifyReceipt"""
        body: dict[str, Any] = {"receipt": receipt}
        if mat is not None:
            body["mat"] = mat
        if context is not None:
            body["context"] = dict(context)
        if prior_receipt is not None:
            body["prior_receipt"] = prior_receipt
        if commitment is not None:
            body["commitment"] = commitment
        return self._request("POST", "/verify", body=body)

    def get_anchors(self) -> list[dict[str, Any]]:
        """operationId: getAnchors"""
        return self._request("GET", "/anchors")

    def compute_digest(self, context: Mapping[str, Any]) -> dict[str, Any]:
        """operationId: computeDigest"""
        return self._request("POST", "/digest", body=dict(context))

    def list_revocations(self) -> list[dict[str, Any]]:
        """operationId: listRevocations"""
        return self._request("GET", "/revocations")

    def get_metrics(self) -> str:
        """operationId: getMetrics"""
        return self._request("GET", "/metrics", accept="text/plain")

    # ---- admin / license-gated ----

    def issue_artifact(self, mat_hex: str) -> dict[str, Any]:
        """operationId: issueArtifact"""
        return self._request("POST", "/artifacts", body={"mat_hex": mat_hex}, auth=True)

    def derive_artifact(self, parent_hex: str, child_hex: str) -> dict[str, Any]:
        """operationId: deriveArtifact"""
        return self._request(
            "POST",
            "/artifacts/derive",
            body={"parent_hex": parent_hex, "child_hex": child_hex},
            auth=True,
        )

    def revoke_artifact(self, artifact_id: str, *, reason: Optional[str] = None) -> dict[str, Any]:
        """operationId: revokeArtifact"""
        q = {"reason": reason} if reason else None
        return self._request(
            "POST", f"/artifacts/{artifact_id}/revoke", query=q, auth=True
        )

    def execute_request(
        self,
        *,
        mat: str,
        action: str,
        resource: str,
        context: Mapping[str, Any],
        impact: Optional[int] = None,
        evidence: Optional[Sequence[Mapping[str, Any]]] = None,
        receipt_id: Optional[str] = None,
        commitment: Optional[str] = None,
        resource_keys: Optional[Sequence[str]] = None,
    ) -> dict[str, Any]:
        """operationId: executeRequest"""
        body: dict[str, Any] = {
            "mat": mat,
            "action": action,
            "resource": resource,
            "context": dict(context),
        }
        if impact is not None:
            body["impact"] = impact
        if evidence is not None:
            body["evidence"] = list(evidence)
        if receipt_id is not None:
            body["receipt_id"] = receipt_id
        if commitment is not None:
            body["commitment"] = commitment
        if resource_keys is not None:
            body["resource_keys"] = list(resource_keys)
        return self._request("POST", "/execute", body=body, auth=True)

    # Alias kept for readability in application code.
    execute = execute_request

    def list_receipts(self) -> list[dict[str, Any]]:
        """operationId: listReceipts"""
        return self._request("GET", "/receipts", auth=True)

    def verify_chain(self) -> dict[str, Any]:
        """operationId: verifyChain"""
        return self._request("GET", "/chain/verify", auth=True)

    # ---- handshake (HTTP binding) ----

    def handshake_capability(self, versions: Sequence[str]) -> dict[str, Any]:
        """operationId: handshakeCapability"""
        return self._request(
            "POST", "/handshake/capability", body={"versions": list(versions)}
        )

    def handshake_challenge(self, session_id: str) -> dict[str, Any]:
        """operationId: handshakeChallenge"""
        return self._request("POST", f"/handshake/{session_id}/challenge", body={})

    def handshake_proof(self, session_id: str, nonce: str) -> dict[str, Any]:
        """operationId: handshakeProof — nonce is the challenge nonce as hex."""
        return self._request(
            "POST",
            f"/handshake/{session_id}/proof",
            body={"nonce": nonce},
        )

    def handshake_authority(self, session_id: str, mat: str) -> dict[str, Any]:
        """operationId: handshakeAuthority — mat is COSE_Sign1 envelope CBOR hex."""
        return self._request(
            "POST",
            f"/handshake/{session_id}/authority",
            body={"mat": mat},
        )

    def handshake_negotiate(self, session_id: str, constraints: str) -> dict[str, Any]:
        """operationId: handshakeNegotiate — constraints is canonical CBOR hex."""
        return self._request(
            "POST",
            f"/handshake/{session_id}/negotiate",
            body={"constraints": constraints},
        )

    def handshake_bind(self, session_id: str) -> dict[str, Any]:
        """operationId: handshakeBind"""
        return self._request("POST", f"/handshake/{session_id}/bind", body={})

    def handshake_status(self, session_id: str) -> dict[str, Any]:
        """operationId: handshakeStatus"""
        return self._request("GET", f"/handshake/{session_id}")
