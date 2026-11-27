---
name: update-hitbox-templates
description: Regenerate the hitbox templates of the visibility check (internal/visibility/hitbox_templates.go) from hitboxes captured in CS2. Captures the hitbox capsules of the players of a demo with the dev build of the CS Demo Manager plugin, keeps the captures in debug/hitboxes/, runs cmd/hitbox-templates on all of them and reads its error report. Use when asked to update, regenerate, restore or improve the hitbox templates, to add a demo to them, after a CS2 update that changes the player models or animations, or when hitbox_templates.go is missing and the visibility package no longer compiles.
user-invocable: true
---

## Why

The analyzer has no access to the skeleton of the players: a demo only gives their position, view angles and state.
To know which parts of a body are visible, the visibility check places a **template** on the player: a set of points
that stand for the hitboxes, in the frame of the player (forward, right, up from the feet).

There are 36 templates, one per combination of stance (standing, crouching, airborne), movement (still, moving) and
held weapon (other, pistol, rifle, scoped, SMG, grenade). Each is the average of real poses captured in the game, so
the file is generated, never edited:

1. The `csdm_dump_hitboxes` command of the CS Demo Manager CS2 plugin writes the world-space hitbox capsules of every
   alive player at a tick of a demo playback. Run at many ticks, it gives a **capture**: a JSON lines file of poses.
2. `cmd/hitbox-templates` joins each pose with what the analyzer knows about the player at that tick (the demo is
   needed for that), sorts the poses by template and fits the points.

Two things follow from this and explain most of the steps below:

- The generator rebuilds the file from the captures it is given and nothing else. Passing one new capture alone
  **replaces** the templates with ones fitted on that capture only. Every run takes every capture.
- A template needs 30 poses of its own, otherwise it falls back on a coarser one (its stance). A capture of a few
  ticks produces a file where the 12 templates of a stance are identical. A useful capture has thousands of poses.

The hitboxes themselves (19 capsules, their bones, sizes and radii) are read from the game and are the same for
every agent model, so more demos never make _them_ more accurate. Only the poses vary.

## What you need

- This repository next to the [cs-demo-manager](https://github.com/akiver/cs-demo-manager) repository
  (`../cs-demo-manager`, or set `CSDM_REPO`), with its CLI built: `vp run dev:cli` there creates `out/cli.js`.
- A **dev build** of the CS2 plugin in `static/cs2/` of that repository, the only one with `csdm_dump_hitboxes`:
  the `Debug` configuration on Windows, `make DEV_COMMANDS=1 build` on Linux. Its
  `cs2-server-plugin/AGENTS.md` explains how to build and how a playback is started from the CLI.
- CS2 installed, Steam running, and CS2 **not** running. A capture takes over the screen for two to three minutes:
  tell the user before starting one, and ask first if the game is already open, it may be their session.

`debug/hitboxes/` is local and ignored by git: it holds `<name>.jsonl` for each capture and `<name>.dem` for its
demo, unless the demo is already in `cs-demos/cs2/` under the same name. If the folder is empty, nothing was
captured on this machine yet and every demo has to be captured again before generating.

## Choosing the demos

Ask the user which demos to capture, unless they already named them. Don't pick demos yourself.

Tell them the demos must be recent as the game only plays demos recorded since its last update that broke the older
ones. Most of `cs-demos/cs2/` is too old for that reason.

To only restore the file, no demo is needed, see [Restoring the file](#restoring-the-file).

## Steps

### 1. Capture each chosen demo

```sh
skills/update-hitbox-templates/scripts/capture.sh <demo.dem> [start tick] [dump count]
```

The script starts CS2 on the demo at 3 times the speed, dumps the hitboxes every 24 ticks, 750 times by default
(18000 ticks, about 95 seconds of playback), quits the game, moves the capture to `debug/hitboxes/<name>.jsonl` and
copies the demo next to it. It refuses to run when the game is open or when the capture would run past the end of
the demo, since a command at a tick the playback never reaches never runs and the game would stay open. A
`csdm_hitboxes.jsonl` already in the game folder is set aside and put back.

Choosing the window matters more than the number of demos:

- The start tick defaults to 20% of the demo. Early ticks are warmup and pistol rounds, a window in the middle of a
  half gives rifles, scoped weapons and SMGs.
- Look at the poses per template printed by the last generation (the `// standing, still, scoped, 435 poses`
  comments of the generated file) and pick demos and windows that feed the thin ones.
- Run one capture at a time: each one needs the game.

A healthy capture prints around 5000 to 6500 poses for 750 dumps and no skipped player. When it reports skipped
players, or writes nothing, read `csdm.log` next to the CS2 executable:

- `Executing: csdm_dump_hitboxes` lines followed by no `players written` line: the game doesn't know the command.
  The installed plugin isn't a dev build, or the playback settings of the dev app select an older plugin version.
- `Fields missing`, `Hitboxes of ... not found` or `skipped: no ...`: a CS2 update moved something the plugin reads
  in the game's memory. That is fixed in the plugin (`commands/csdm_dump_models.cpp`, whose `csdm_dump_pawn_model`
  command shows the current layout), not here.

### 2. Generate

```sh
skills/update-hitbox-templates/scripts/generate.sh
```

It runs `go run ./cmd/hitbox-templates -report` with every capture of `debug/hitboxes/` and its demo, prints the
poses per template, then runs `gofmt`, `go vet` and the tests of the packages that use the templates. When
`hitbox_templates.go` is missing, it first writes an empty set of templates so that the generator, which imports the
package, compiles.

### 3. Read the output

For each demo:

```
debug/hitboxes/<name>.dem: 5538/5543 captures joined, mean origin distance 0.69
```

- Nearly every pose must join. A handful missing is normal (a player the analyzer has no state for at that tick).
- The mean origin distance is how far the game's player origin is from the analyzer's at the same tick. It is under
  1 for a good capture. Around 2 or more means the playback and the ticks don't line up, which was seen with dumps
  taken right after a seek: capture again from a continuous playback.

Then the report, which fits the templates on the other demos and measures them on each demo (it needs two demos):

```
Signed distance of the points to the body: p10 -5.3 p50 -2.3 p90 3.5, 14% outside
Distance from each capsule axis to the closest point, minus the radius (p50 p90 p99):
  head_0        -2.2  -0.5   3.8
  ...
  ankle_l        2.0  12.7  19.5
```

- **Outside** is the share of template points more than 2 units away from the real body. With four demos it is 14%,
  and between 13% and 16% per demo.
- The table says how well each hitbox is covered: negative means a point lies inside the capsule. The torso and the
  head are well covered, the lower legs, ankles and hands are the loosest (p90 around 7 to 13), because one averaged
  pose can't follow a walk cycle.

Compare with these numbers and with the previous run. Expect small moves: going from two to four demos lowered the
outside share by 0.3 to 0.4 points on every demo. The error left comes from averaging poses, not from a lack of
them, so a large change in either direction points at a bad capture or at a change in the generator, not at the new
data.

### 4. Report

Tell the user:

- which demos were captured, with their window and number of poses;
- the join and origin distance lines, the outside share before and after, and any template that still falls back
  (fewer than 30 poses of its own);
- that `internal/visibility/hitbox_templates.go` changed and that the checks pass. The captures and demos stay
  local, the generated file is what gets committed.

## Restoring the file

`hitbox_templates.go` deleted or overwritten by a run on too few captures: run step 2 alone. With the same captures
it gives the same file.
