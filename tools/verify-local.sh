#!/usr/bin/env bash
set -euo pipefail
cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.."
case "${1:-all}" in
  --help) echo 'Usage: tools/verify-local.sh [all] (requires Go from go.mod)'; exit 0 ;;
  all) ;;
  *) echo 'Unknown validation group' >&2; exit 2 ;;
esac
[[ $# -le 1 ]] || { echo 'Expected at most one group' >&2; exit 2; }
go test ./...
go vet ./...
