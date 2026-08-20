#!/usr/bin/env bash
# Run multi-lang client gates (no pip required for Python).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "== openapi coverage =="
"$ROOT/scripts/check-openapi-coverage.sh"

echo "== python unittest =="
(cd "$ROOT/python" && PYTHONPATH=. python3 -m unittest discover -s tests -v)

echo "== typescript =="
(cd "$ROOT/typescript" && npm test)

echo "== all client gates green =="
