"""Unit tests with a fake HTTP handler (stdlib only — no pytest/pip)."""

from __future__ import annotations

import json
import sys
import threading
import unittest
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from xap_client import Client, XAPError


class _Handler(BaseHTTPRequestHandler):
    def log_message(self, *args):  # noqa: ARG002
        pass

    def _read(self):
        n = int(self.headers.get("Content-Length", "0"))
        return self.rfile.read(n) if n else b""

    def do_GET(self):  # noqa: N802
        if self.path.endswith("/anchors"):
            body = json.dumps([{"kid_hex": "ab", "alg": "ed25519"}]).encode()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return
        if "/chain/verify" in self.path:
            if self.headers.get("Authorization") != "Bearer secret":
                self.send_response(403)
                self.end_headers()
                self.wfile.write(b"forbidden")
                return
            body = b'{"intact":true,"count":1}'
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return
        self.send_response(404)
        self.end_headers()

    def do_POST(self):  # noqa: N802
        raw = self._read()
        if self.path.endswith("/verify"):
            body = b'{"valid":true}'
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return
        if self.path.endswith("/execute"):
            data = json.loads(raw.decode())
            assert data["action"] == "read"
            body = b'{"decision":"permit","receipt_id":"r1"}'
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return
        self.send_response(404)
        self.end_headers()


class ClientTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.httpd = HTTPServer(("127.0.0.1", 0), _Handler)
        t = threading.Thread(target=cls.httpd.serve_forever, daemon=True)
        t.start()
        host, port = cls.httpd.server_address
        cls.base = f"http://{host}:{port}/xap/v1"

    @classmethod
    def tearDownClass(cls):
        cls.httpd.shutdown()

    def test_get_anchors(self):
        anchors = Client(self.base).get_anchors()
        self.assertEqual(anchors[0]["alg"], "ed25519")

    def test_verify_receipt(self):
        r = Client(self.base).verify_receipt("aabb")
        self.assertTrue(r["valid"])

    def test_admin_required(self):
        with self.assertRaises(XAPError) as cm:
            Client(self.base).verify_chain()
        self.assertEqual(cm.exception.status, 401)

    def test_execute_with_token(self):
        c = Client(self.base, token="secret")
        r = c.execute_request(
            mat="ee",
            action="read",
            resource="/x",
            context={"time": "2026-01-01T00:00:00Z", "network_zone": "dmz"},
        )
        self.assertEqual(r["decision"], "permit")

    def test_verify_chain_auth(self):
        c = Client(self.base, token="secret")
        self.assertTrue(c.verify_chain()["intact"])


if __name__ == "__main__":
    unittest.main()
