#!/usr/bin/env python3
"""Fetch trust anchors from a self-hosted node (public, no token)."""
import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "python"))
from xap_client import Client

base = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080/xap/v1"
for a in Client(base).get_anchors():
    print(a.get("kid_hex"), a.get("alg"))
