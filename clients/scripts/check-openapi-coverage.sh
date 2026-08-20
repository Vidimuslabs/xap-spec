#!/usr/bin/env bash
# Fail if any OpenAPI operationId is missing from Python or TypeScript clients.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SPEC="${ROOT}/../openapi/xap.yaml"
PY="${ROOT}/python/xap_client/client.py"
TS="${ROOT}/typescript/src/client.ts"

if [[ ! -f "$SPEC" ]]; then
  echo "FAIL missing OpenAPI at $SPEC" >&2
  exit 1
fi

mapfile -t OPS < <(grep -E '^\s+operationId:' "$SPEC" | awk '{print $2}' | sort -u)
fail=0

# OpenAPI operationId -> expected marker in source (snake or camel).
snake() {
  # verifyReceipt -> verify_receipt
  echo "$1" | sed -E 's/([a-z0-9])([A-Z])/\1_\2/g' | tr '[:upper:]' '[:lower:]'
}

echo "== OpenAPI operation coverage (${#OPS[@]} ops) =="
for op in "${OPS[@]}"; do
  sn=$(snake "$op")
  # Python: method def + docstring operationId
  if ! grep -qE "operationId: ${op}|def ${sn}\(" "$PY"; then
    echo "FAIL python missing $op (expect def ${sn} or docstring)"
    fail=$((fail + 1))
  else
    echo "OK   python $op"
  fi
  # TypeScript: method name is same camelCase as operationId, or docstring
  if ! grep -qE "operationId: ${op}|${op}\(|${op}:" "$TS"; then
    # method names: verifyReceipt stays verifyReceipt in TS
    if ! grep -qE "${op}\(" "$TS"; then
      echo "FAIL typescript missing $op"
      fail=$((fail + 1))
    else
      echo "OK   typescript $op"
    fi
  else
    echo "OK   typescript $op"
  fi
done

echo "== coverage: $fail failure(s) =="
exit "$fail"
