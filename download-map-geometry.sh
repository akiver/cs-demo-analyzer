#!/bin/bash

# Downloads the CS2 map geometry used by the tests to compute the time to damage/kill from
# https://gitlab.com/akiver/cs2-map-geometry
#
# The version is pinned because regenerating the geometry after a map update changes the snapshots, to use a newer
# version bump MAP_GEOMETRY_REF and update the snapshots.
# MAPS must contain the map of every CS2 demo used by the tests.

set -e

MAP_GEOMETRY_REF="v1"
MAPS="de_ancient de_anubis de_mirage de_nuke de_train de_vertigo"
FOLDER="cs2-map-geometry"
BASE_URL="https://gitlab.com/akiver/cs2-map-geometry/-/raw/$MAP_GEOMETRY_REF"

# Files downloaded for another version are stale.
if [ "$(cat "$FOLDER/.ref" 2>/dev/null)" != "$MAP_GEOMETRY_REF" ]; then
  rm -rf "$FOLDER"
fi
mkdir -p "$FOLDER"

for map in $MAPS; do
  path="$FOLDER/$map.tri"
  if [ -f "$path" ]; then
    continue
  fi

  echo "Downloading map geometry $path"
  # Download to a temporary file so that an interrupted download is not considered as done.
  curl -s -f -L -o "$path.tmp" "$BASE_URL/$map.tri"
  mv "$path.tmp" "$path"
done

echo "$MAP_GEOMETRY_REF" > "$FOLDER/.ref"
