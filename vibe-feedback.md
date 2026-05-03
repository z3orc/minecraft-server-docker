# Codebase Evaluation

## Overview

A Go-based runner for Minecraft (Fabric) servers in Docker. It manages the server lifecycle — downloading the JAR, writing `server.properties` from env vars, managing ops/whitelist JSON files, and forwarding OS signals for graceful shutdown. The architecture is clean and well-decomposed.

---

## Issues Found

### Bugs

**1. `parse()` is applied to values, not keys (`internal/minecraft/properties/properties.go:99`)**
The function lowercases everything and converts `_` to `-`. This is the right transform for property *keys*, but it's called on property *values*. Setting `RCON_PASSWORD=my_secret` writes `rcon.password=my-secret`. Setting `LEVEL_NAME=Cool_World` writes `level-name=cool-world`. This corrupts passwords, seeds, MOTDs, and any value with underscores.

**2. `parseProp` rejects values containing `=` (`internal/minecraft/properties/properties.go:88-89`)**
`strings.Split(prop, "=")` splits on every `=`, then checks `len(vals) != 2`. A MOTD or resource pack URL containing `=` (e.g., Base64 values in URLs) would fail. Should be `strings.SplitN(prop, "=", 2)`.

**3. `AddPlayerToList` in `internal/minecraft/manage/playerlist.go:61` is dead and broken**
It's never called (the typed functions `AddPlayerToOpsList`/`AddPlayerToWhitelist` are used instead). It also has a panic: it always asserts `playerList.(Whitelist)` at line 78 regardless of whether `listType` is `OPS_LIST`.

**4. `download.ServerExecutable` uses bare `http.Get` without a timeout (`internal/download/executable.go:17`)**
Every other HTTP call goes through `httpclient.New()` which sets a 10s timeout. A slow or hung download of the server JAR will stall the container forever.

---

### Inconsistencies / Dead Code

**5. CLI flags struct and comments are dead weight (`cmd/runner.go:13-47`)**
The `flags` struct and all the commented-out `flag.Parse()` code serve no purpose. Env vars are read directly. The Makefile's `run` target still passes `-jar`, `-timeout`, `-debug`, `-sigkill` which are silently ignored.

**6. `Dockerfile` references `golang:1.25-alpine` (`Dockerfile:1`)**
Go 1.25 doesn't exist. This will fail to pull. Should be `golang:1.24-alpine`.

**7. `entrypoint.sh` is copied into the image but never used (`Dockerfile:17, 36-38`)**
`ENTRYPOINT` calls `/app/runner` directly; the shell script is dead.

---

### Minor

**8. `properties.Write()` truncates the existing file (`internal/minecraft/properties/properties.go:63`)**
`os.Create` overwrites the whole file, so any properties not set via env vars (e.g., hand-edited entries) are erased on each restart. Probably intentional, but means the container is not idempotent for properties set by the Minecraft server itself on first run.

**9. `httpclient.New()` creates a new `*http.Client` per call**
Each API call allocates a fresh client, forgoing connection pooling. For the startup path (2 Fabric API calls + 1 Mojang call per player) this is harmless, but worth noting.

---

## Summary

The two most impactful bugs are **#1** (values silently corrupted) and **#2** (values with `=` rejected), both in `properties.go`. The Dockerfile version (`golang:1.25`) would prevent any Docker build from succeeding. `AddPlayerToList` should just be deleted.
