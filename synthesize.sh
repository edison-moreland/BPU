#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(dirname "$(realpath "$0")")"
OUTPUT_DIR="./.out"
MODULES_DIR="./modules"

NETLIST_SKIN="./synth/netlistsvg_skin.svg"
LOGICWORLD_LIBERTY="./synth/logicworld.lib"
FF_INTO_LATCH_TECHMAP="./synth/flipflop2latch_techmap.v"

MODULE="$1"
MODULE_NAME="$(basename "${MODULE}")"

. "${PROJECT_DIR}/synth/frontmatter.sh"

# Dependencies
# Yosys 0.46 (git sha1 e97731b9d, clang++ 16.0.0 -fPIC -O3)
# netlistsvg
# sv2v v0.0.12
# probably gnu awk+grep

pushd "${PROJECT_DIR}"
mkdir -p "${OUTPUT_DIR}"

module_file="${MODULES_DIR}/${MODULE}.sv"
if [[ ! -f "${module_file}" ]]; then

    module_file="${MODULES_DIR}/${MODULE}.v"
    if [[ ! -f "${module_file}" ]]; then
        echo "Module not found"
        exit 1
    else
        cp "${module_file}" "${OUTPUT_DIR}"
    fi
else
    sv2v -w "${OUTPUT_DIR}" "${module_file}"
fi

case "$(get_frontmatter_key visualization)" in
    * | netlistsvg)
        yosys -o "${OUTPUT_DIR}/${MODULE_NAME}.json" -S "${OUTPUT_DIR}/"*.v \
            -p "prep -flatten -top ${MODULE_NAME}" \
            -p 'dfflegalize -cell $_DLATCH_P_ 0 -cell $_DFF_P_ 0 -cell $_DFF_PP0_ 0 -cell $_ALDFF_PP_ 0' \
            -p "techmap -autoproc -map ${FF_INTO_LATCH_TECHMAP}; opt_merge" \
            -p 'freduce -inv; opt -full' \
            -p "read_liberty -lib ${LOGICWORLD_LIBERTY}"\
            -p "abc -liberty ${LOGICWORLD_LIBERTY}; opt_clean"\

        # TODO: Elk layout file to increase readibility?
        netlistsvg "${OUTPUT_DIR}/${MODULE_NAME}.json" -o "${OUTPUT_DIR}/${MODULE_NAME}.svg" --skin "${NETLIST_SKIN}"
    ;;

    # netlistsvg-sop)
    #     #  Sum of products TODO - DOESNT WORK
    #     yosys -o "${OUTPUT_DIR}/${MODULE}.json" -S "${OUTPUT_DIR}/"*.v -p 'synth -flatten -top '"${MODULE}; abc -sop"
    #     # yosys-abc  -o "${OUTPUT_DIR}/${MODULE}.sop.blif" "${OUTPUT_DIR}/${MODULE}.blif"
    #     netlistsvg "${OUTPUT_DIR}/${MODULE}.json" -o "${OUTPUT_DIR}/${MODULE}.svg" --skin "${NETLIST_SKIN}"

    # ;;

esac


popd