#!/usr/bin/env bash

PROJECT_DIR="$(dirname "$(realpath "$0")")"
OUTPUT_DIR="${PROJECT_DIR}/.out"
MODULES_DIR="${PROJECT_DIR}/modules"

MODULE="$1"

mkdir -p "${OUTPUT_DIR}"
yosys -o "${OUTPUT_DIR}/${MODULE}.json" -S "${MODULES_DIR}/${MODULE}.v" -p 'prep -auto-top -flatten; abc -g AND,XOR,OR'
netlistsvg "${OUTPUT_DIR}/${MODULE}.json" -o "${OUTPUT_DIR}/${MODULE}.svg"