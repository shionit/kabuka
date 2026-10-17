---
name: verify-kabuka
description: Build and drive the kabuka stock-price CLI end to end against live Yahoo! Finance Japan, capturing stdout/stderr/exit-code evidence. Use when you need to prove a kabuka change works for a real user (symbol lookup, -d/--detail, -f/--format text|json|csv, error paths), not just that unit tests pass.
---

# Verify kabuka

kabuka is a short-lived Go CLI (`kabuka [symbol] [-d] [-f text|json|csv]`). It prints one stock quote scraped from `https://finance.yahoo.co.jp/` and exits. The only user surface is the terminal. There is no server, no config file and no local state, so "launch" means building the binary once per run and invoking it per check. There is no TTY behavior worth driving, so a plain subprocess is the harness (no tmux/PTY).

Everything goes through `scripts/kabuka-verify.sh` (paths below are relative to the repo root):

```bash
H=.claude/skills/verify-kabuka/scripts/kabuka-verify.sh
```

Shell state does not persist between tool calls. Repeat `H=...` and pass the run id explicitly on every call. Run under bash. zsh does not word-split `$args`, so `kabuka $args` turns into a single bogus argument.

## Launch

```bash
RUN=$($H build)          # optional: $H build my-run-id
echo "$RUN"
```

`build` runs `go build -o <root>/<run>/bin/kabuka .` from the repo root and records the git revision. It creates three dirs:

- `<root>/<run>/bin/`: the binary and `REVISION`.
- `<root>/<run>/work/`: an empty cwd for every invocation.
- `<root>/<run>/evidence/`: the proof artifacts.

`<root>` defaults to `${TMPDIR:-/tmp}/kabuka-verify`. Override it with `KABUKA_VERIFY_ROOT`, and pass the same value to every call.

Ready means `build` printed a run id and `doctor` reports `healthy`. Each run has its own binary and cwd, so concurrent runs never interfere. It never touches a globally installed `kabuka` on `PATH`. Do not verify with that one, because it may be an old release.

## Doctor

Run it first, and again whenever output looks wrong:

```bash
$H doctor "$RUN"
```

It is read-only. It checks four things:

- The binary exists and answers `--help`.
- The binary was built from the current Go sources. The rev is HEAD, plus a `-dirty-<hash>` fingerprint of uncommitted `.go`/`go.mod`/`go.sum` changes. Edit code after building and this fails, so build a new run.
- `work/` is still empty.
- `https://finance.yahoo.co.jp/` returns HTTP 200.

When the network check fails, every live drive will fail. Report that and stop, rather than reading the failures as regressions.

## Drive

```bash
$H run "$RUN" <evidence-name> -- <kabuka args...>
```

This runs `<root>/<run>/bin/kabuka <args>` with `work/` as cwd. It prints the command, stdout, stderr, exit code, and whether any file appeared in cwd. Evidence names must be unique within a run. Examples drawn from the feature map:

```bash
$H run "$RUN" quote-jp-code   -- 3994.T
$H run "$RUN" detail-json-us  -- -d AAPL -f=json
$H run "$RUN" error-notfound  -- ZZZZNOTREAL
```

Recipes per feature live in [`features/README.md`](features/README.md). Use them as the source of truth for what to drive and what output proves it.

Prices are live and change minute to minute. Assert shape: the symbol, numeric price matching `^[0-9]+(\.[0-9]+)?$`, field presence, and exit code. Never assert exact values. Outside market hours Yahoo may serve `---` in place of a number. That is expected, not a bug (see `CLAUDE.md`). Report it and retry during market hours if a numeric proof is required.

## Evidence

Each `run` writes the following to `<root>/<run>/evidence/`:

- `<name>.cmd`: the exact argv.
- `<name>.stdout`, `<name>.stderr`, `<name>.exit`: the program's output and exit code.
- `<name>.sideeffects`: `none`, or the files the CLI wrote to cwd.

Proof standards:

- Exercise the real binary built from this checkout against the live site. Don't call `Kabuka.fetch()` or use the unexported `baseURL`/`httpClient` test overrides. Those are unit-test seams, not user paths. Production hardcodes the Yahoo URL, so there is no mock boundary for a user-level proof.
- Record the action (`.cmd`) and the result (`.stdout`, `.exit`) together. A passing exit code alone is not proof. Read the output.
- kabuka has one legitimate side effect, the outbound HTTPS GET, and it must write no files. `.sideeffects` must read `none`.
- For a code change, capture a before/after pair. Build a run at the base commit and a run at your change, then drive the same names in both.
- Cite the evidence directory path in your report.

## Cleanup

```bash
$H cleanup "$RUN"
```

This removes `bin/` and `work/` for that run only. It keeps `evidence/` and prints how many files it holds. Confirm that count is nonzero before reporting. There are no background processes to kill, since every invocation exits on its own. Run cleanup after failed attempts too.

## Related

- `make test` runs offline unit tests. `make test-integration` (build tag `integration`) runs live-network fetch tests in Go. They complement this skill, which proves the built CLI behaves for a user.
- When every symbol suddenly fails with `Unknown market type.` or returns an empty price, Yahoo's DOM probably changed. See the selector gotcha in `CLAUDE.md`.
