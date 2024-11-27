#!/usr/bin/env bash
set -euo pipefail


ASSEMBLER_DIR="$(dirname "$(realpath "$0")")"

go run -C "${ASSEMBLER_DIR}/" . <(m4 "${@}")