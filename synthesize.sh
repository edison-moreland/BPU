#!/usr/bin/env bash

PROJECT_DIR="$(dirname "$(realpath "$0")")"
OUTPUT_DIR="${PROJECT_DIR}/.out"
MODULES_DIR="${PROJECT_DIR}/modules"
TECHMAPS_DIR="${PROJECT_DIR}/techmaps"

NETLIST_SKIN="${PROJECT_DIR}/skin.svg"
LW_LIBERTY_FILE="${PROJECT_DIR}/logicworld.lib"

MODULE="$1"

# Dependencies
# Yosys 0.46 (git sha1 e97731b9d, clang++ 16.0.0 -fPIC -O3)
# netlistsvg
# sv2v v0.0.12

cat "${LW_LIBERTY_FILE}"

mkdir -p "${OUTPUT_DIR}"

sv2v -w "${OUTPUT_DIR}" "${MODULES_DIR}/${MODULE}.sv"
yosys -o "${OUTPUT_DIR}/${MODULE}.json" -S "${OUTPUT_DIR}/${MODULE}.v" \
    -p 'prep -auto-top -flatten' \
    -p 'abc -g OR,NOR,NAND,XOR,XNOR; opt -full' \
    -p 'dfflegalize -cell $_DLATCH_P_ 0 -cell $_DFF_P_ 0' \
    -p "techmap -autoproc -map ${TECHMAPS_DIR}/ff2latch.v; opt_merge" \
    -p 'abc -g OR,NOR,NAND,XOR,XNOR' \
    # -p 'freduce -inv; opt -full' \
    # -p 'abc -liberty logicworld.lib'\
    # -p 'read_liberty -lib logicworld.lib'\
    # -p "dfflibmap -liberty logicworld.lib"\

# TODO: Elk layout file to increase readibility?
netlistsvg "${OUTPUT_DIR}/${MODULE}.json" -o "${OUTPUT_DIR}/${MODULE}.svg" --skin "${NETLIST_SKIN}"

