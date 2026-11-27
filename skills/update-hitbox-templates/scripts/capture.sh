#!/usr/bin/env bash
# Captures the hitboxes of the players of a demo by playing it in CS2 with the dev build of the CS Demo Manager plugin,
# and stores the capture next to a copy of the demo in debug/hitboxes/.
#
# Usage: capture.sh <demo.dem> [start tick] [dump count]
#
# The hitboxes are dumped every 24 ticks from the start tick, 750 times by default, which covers 18000 ticks. The start
# tick defaults to 20% of the demo.
#
# Environment:
#   CSDM_REPO    cs-demo-manager repository, default: ../cs-demo-manager next to this repository
#   CS2_BIN_DIR  folder of the CS2 executable, where the plugin writes csdm.log, default: the default Steam library
set -euo pipefail

TICKS_BETWEEN_DUMPS=24
PLAYBACK_SPEED=3

if [ $# -lt 1 ]; then
    sed -n '2,12p' "$0" | sed 's/^# \{0,1\}//'
    exit 1
fi

script_dir=$(cd "$(dirname "$0")" && pwd)
root=$(cd "$script_dir/../../.." && pwd)
csdm_repo=${CSDM_REPO:-$root/../cs-demo-manager}

is_windows=0
case "$(uname)" in
MSYS* | MINGW* | CYGWIN*) is_windows=1 ;;
esac

if [ -n "${CS2_BIN_DIR:-}" ]; then
    bin_dir=$CS2_BIN_DIR
elif [ "$is_windows" = 1 ]; then
    bin_dir="C:/Program Files (x86)/Steam/steamapps/common/Counter-Strike Global Offensive/game/bin/win64"
else
    bin_dir="$HOME/.steam/steam/steamapps/common/Counter-Strike Global Offensive/game/bin/linuxsteamrt64"
fi

demo=$1
if [ ! -f "$demo" ]; then
    echo "Demo not found: $demo"
    exit 1
fi
if [ ! -f "$csdm_repo/out/cli.js" ]; then
    echo "$csdm_repo/out/cli.js not found: run 'vp run dev:cli' in the cs-demo-manager repository, or set CSDM_REPO."
    exit 1
fi
if [ ! -d "$bin_dir" ]; then
    echo "CS2 not found in $bin_dir: set CS2_BIN_DIR to the folder of the CS2 executable."
    exit 1
fi

is_game_running() {
    if [ "$is_windows" = 1 ]; then
        tasklist 2>/dev/null | grep -qi '^cs2\.exe'
    else
        pgrep -x cs2 >/dev/null 2>&1
    fi
}

# The watch command would send the demo to the running game and return right away.
if is_game_running; then
    echo "CS2 is already running: close it first."
    exit 1
fi

# An absolute path the game understands: on Windows, the drive letter form rather than the /c/ one of the shell.
demo_dir=$(cd "$(dirname "$demo")" && { pwd -W 2>/dev/null || pwd; })
name=$(basename "$demo" .dem)
demo_path="$demo_dir/$name.dem"

playback_ticks=$(node "$script_dir/demo-ticks.mjs" "$demo_path" | tr -d '\r')
dump_count=${3:-750}
if [ -n "${2:-}" ]; then
    start_tick=$2
elif [ "$playback_ticks" -gt 0 ]; then
    start_tick=$((playback_ticks / 5))
else
    echo "The demo doesn't say how long it is: pass the start tick."
    exit 1
fi

quit_tick=$((start_tick + TICKS_BETWEEN_DUMPS * (dump_count + 1) + 50))
# A command at a tick the demo never reaches never runs: the game would stay open without quitting.
if [ "$playback_ticks" -gt 0 ] && [ "$quit_tick" -ge "$playback_ticks" ]; then
    echo "The capture would end at tick $quit_tick, the demo ends at $playback_ticks: start earlier or dump less."
    exit 1
fi
# Windows limits a command line to 32767 characters, each dump adds about 35.
if [ "$is_windows" = 1 ] && [ "$dump_count" -gt 850 ]; then
    echo "At most 850 dumps fit in a Windows command line: run several captures with different start ticks."
    exit 1
fi

# The plugin appends to csdm_hitboxes.jsonl next to its log: an existing one is set aside and put back at the end.
dump_path="$bin_dir/csdm_hitboxes.jsonl"
kept_dump_path="$dump_path.kept"
if [ -f "$dump_path" ]; then
    mv "$dump_path" "$kept_dump_path"
fi
restore_dump() {
    # The running game may still append to the dump: deleting it would lose the partial capture, and putting the kept
    # one back would have later dumps appended to it.
    if is_game_running; then
        echo "CS2 is running, the dumps are left as they are. Once it's closed:"
        echo "- $dump_path holds this partial capture, delete it or move it away."
        if [ -f "$kept_dump_path" ]; then
            echo "- $kept_dump_path is the previous capture, rename it back to csdm_hitboxes.jsonl."
        fi
        return
    fi
    rm -f "$dump_path"
    if [ -f "$kept_dump_path" ]; then
        mv "$kept_dump_path" "$dump_path"
    fi
}
trap restore_dump EXIT

args=(--tick "$start_tick" --exec "96:demo_timescale $PLAYBACK_SPEED")
for ((i = 1; i <= dump_count; i++)); do
    args+=(--exec "$((start_tick + i * TICKS_BETWEEN_DUMPS)):csdm_dump_hitboxes")
done
args+=(--exec "$quit_tick:quit")

echo "Capturing $name from tick $start_tick to $quit_tick ($dump_count dumps)..."
# The playback lasts the ticks at 64 per second divided by the speed, the game needs up to a few minutes to start.
max_seconds=$(((quit_tick - start_tick) / 64 / PLAYBACK_SPEED + 300))
(cd "$csdm_repo" && timeout "$max_seconds" node out/cli.js watch "$demo_path" "${args[@]}") || true

if is_game_running; then
    echo "CS2 is still running, the quit command was never reached: close the game and check $bin_dir/csdm.log."
    exit 1
fi
if [ ! -s "$dump_path" ]; then
    echo "No capture written. Check $bin_dir/csdm.log:"
    echo "- 'Executing: csdm_dump_hitboxes' lines followed by no 'players written' line: the installed plugin isn't a dev"
    echo "  build, build its Debug configuration."
    echo "- 'Fields missing' or 'Hitboxes of ... not found': a CS2 update moved what the plugin reads, fix the plugin first."
    exit 1
fi

captures_dir="$root/debug/hitboxes"
mkdir -p "$captures_dir"
capture_path="$captures_dir/$name.jsonl"
mv "$dump_path" "$capture_path"

# The generator needs the demo of each capture: keep a copy next to it unless the repository already has the demo.
if [ ! -f "$captures_dir/$name.dem" ] && [ ! -f "$root/cs-demos/cs2/$name.dem" ]; then
    cp "$demo_path" "$captures_dir/$name.dem"
    echo "Demo copied to debug/hitboxes/$name.dem"
fi

player_count=$(grep -c '"type":"player"' "$capture_path" || true)
written_count=$(grep -c 'players written' "$bin_dir/csdm.log" || true)
skipped_count=$(grep -c 'skipped: no' "$bin_dir/csdm.log" || true)
echo "debug/hitboxes/$name.jsonl: $player_count poses from $written_count dumps, $skipped_count players skipped for a missing model, bones or hitboxes"
