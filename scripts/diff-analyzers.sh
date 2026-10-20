#!/bin/bash
# Shows how the CS:GO analyzer (internal/analyzer/csgo, demoinfocs-golang v4) differs from the CS2 analyzer
# (internal/analyzer/cs2, demoinfocs-golang v6). Both are mostly the same code: a change to one usually belongs in
# the other too. Package names and demoinfocs import paths are normalised so only real differences show.
set -e
cd "$(dirname "$0")/.."

normalise() {
	local dir=$1 pkg=$2 version=$3
	for f in "$dir"/*.go; do
		echo "### $(basename "$f")"
		sed -e "s|demoinfocs-golang/$version/|demoinfocs-golang/vX/|" -e "1s/^package $pkg\$/package analyzer/" "$f"
	done
}

diff -u --label cs2 --label csgo \
	<(normalise internal/analyzer/cs2 cs2 v6) \
	<(normalise internal/analyzer/csgo csgo v4) \
	"$@" || true
