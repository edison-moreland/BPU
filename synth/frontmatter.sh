#!/bin/bash
set -euo pipefail

# Parsing for the frontmatter
# Example:
# ```systemverilog
# //visualization: netlistsvg
# ```

# Read frontmatter into variable $frontmatter
function read_frontmatter() {
    FILE_PATH="$1"

    # This awk will give us every line that starts with #-, and it will split the fields
    # In:
    #   //visualization: netlistsvg
    # Out: ('-' is a placeholder for unit seperator '\x1f')
    #   visualization-netlistsvg
    mapfile -t frontmatter < <(awk -F': ' '/^\/\// {gsub(/^\/\//, "", $1); print tolower($1) "\x1F" $2} ' < "${FILE_PATH}" )
}

# Get the value of a key in the frontmatter
function get_frontmatter_key() {
    echo "${frontmatter}" | grep -e "^${1}" | cut -d$'\x1f' -f2
}