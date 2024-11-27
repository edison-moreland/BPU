#!/usr/bin/env bash
set -euo pipefail

# Use: ./synth.sh <hdl dir>
# Runs the synth.tcl script for the given hdl dir

export PROJECT_DIR="$(dirname "$(realpath "$0")")"
export SYNTH_DIR="${PROJECT_DIR}/synth"

export TARGET="${1}"
export TARGET_DIR="${PROJECT_DIR}/hdl/${TARGET}"
export OUTPUT_DIR="${PROJECT_DIR}/.out/${TARGET}"

if [[ ! -d "${TARGET_DIR}" ]]; then
    printf "Can't find target ${TARGET}\n"
    exit 1
fi

if [[ ! -f "${TARGET_DIR}/synth.tcl" ]]; then
    printf "Can't find synth script for target ${TARGET}\n"
    exit 1
fi

mkdir -p "${OUTPUT_DIR}"
"${TARGET_DIR}/synth.tcl"
