#!/usr/bin/env bash
# Verification harness for the kabuka CLI. See ../SKILL.md for usage.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)"
ROOT="${KABUKA_VERIFY_ROOT:-${TMPDIR:-/tmp}/kabuka-verify}"
ROOT="${ROOT%/}"

usage() {
  cat <<EOF
usage:
  $0 build   [run-id]                 build ./kabuka for a new run, print run id
  $0 doctor  <run-id>                 read-only health check of a run
  $0 run     <run-id> <name> -- ARGS  run the built binary with ARGS, record evidence
  $0 cleanup <run-id>                 remove bin/ and work/, keep evidence/
root: $ROOT
EOF
  exit 2
}

run_dir() {
  [[ -n "${1:-}" ]] || usage
  echo "$ROOT/$1"
}

# HEAD plus a fingerprint of uncommitted Go sources, so a binary built from
# a dirty tree still matches until the sources change again.
repo_rev() {
  local rev untracked
  rev="$(git -C "$REPO" rev-parse --short HEAD)"
  untracked="$(git -C "$REPO" ls-files --others --exclude-standard -- '*.go' go.mod go.sum)"
  if [[ -z "$(git -C "$REPO" status --porcelain -- '*.go' go.mod go.sum)" ]]; then
    echo "$rev"
    return
  fi
  echo "$rev-dirty-$( { git -C "$REPO" diff HEAD -- '*.go' go.mod go.sum
    for f in $untracked; do echo "$f"; cat "$REPO/$f"; done; } | shasum | cut -c1-8)"
}

cmd_build() {
  local id="${1:-$(date +%Y%m%d-%H%M%S)-$$}"
  local dir="$ROOT/$id"
  [[ -e "$dir" ]] && { echo "run $id already exists at $dir" >&2; exit 1; }
  mkdir -p "$dir/bin" "$dir/work" "$dir/evidence"
  (cd "$REPO" && go build -o "$dir/bin/kabuka" .)
  repo_rev > "$dir/bin/REVISION"
  echo "$id"
}

cmd_doctor() {
  local dir; dir="$(run_dir "${1:-}")"
  local ok=1
  check() { if eval "$2"; then echo "ok    $1"; else echo "FAIL  $1"; ok=0; fi; }
  check "run dir $dir"                  "[[ -d '$dir' ]]"
  check "binary built"                  "[[ -x '$dir/bin/kabuka' ]]"
  check "binary answers --help"         "'$dir/bin/kabuka' --help >/dev/null 2>&1"
  local built now
  built="$(cat "$dir/bin/REVISION" 2>/dev/null || echo none)"
  now="$(repo_rev)"
  check "binary rev $built matches checkout $now" "[[ '$built' == '$now' ]]"
  check "work/ is empty"                "[[ -z \"\$(ls -A '$dir/work' 2>/dev/null)\" ]]"
  local code
  code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 10 'https://finance.yahoo.co.jp/' || echo 000)"
  check "finance.yahoo.co.jp reachable (HTTP $code)" "[[ '$code' == 200 ]]"
  [[ $ok == 1 ]] || { echo "doctor: NOT healthy (rev mismatch means sources changed: build a new run)"; exit 1; }
  echo "doctor: healthy"
}

cmd_run() {
  local dir; dir="$(run_dir "${1:-}")"; shift
  local name="${1:-}"; shift || true
  [[ -n "$name" && "${1:-}" == "--" ]] || usage
  shift
  [[ -x "$dir/bin/kabuka" ]] || { echo "no binary in $dir; run build first" >&2; exit 1; }
  local ev="$dir/evidence/$name"
  [[ -e "$ev.cmd" ]] && { echo "evidence $name already exists; pick a new name" >&2; exit 1; }
  printf 'kabuka' > "$ev.cmd"; [[ $# -gt 0 ]] && printf ' %q' "$@" >> "$ev.cmd"; echo >> "$ev.cmd"
  local rc=0
  (cd "$dir/work" && "$dir/bin/kabuka" "$@") > "$ev.stdout" 2> "$ev.stderr" || rc=$?
  echo "$rc" > "$ev.exit"
  # The CLI must not write files: work/ is its cwd and must stay empty.
  local leaked; leaked="$(ls -A "$dir/work")"
  if [[ -n "$leaked" ]]; then echo "files written to cwd: $leaked" > "$ev.sideeffects"; else echo "none" > "$ev.sideeffects"; fi
  echo "\$ $(cat "$ev.cmd")"
  echo "--- stdout";  cat "$ev.stdout"
  echo "--- stderr";  cat "$ev.stderr"
  echo "--- exit $rc"
  echo "--- files written: $(cat "$ev.sideeffects")"
  echo "--- evidence: $ev.{cmd,stdout,stderr,exit,sideeffects}"
}

cmd_cleanup() {
  local dir; dir="$(run_dir "${1:-}")"
  [[ -d "$dir" ]] || { echo "no run at $dir" >&2; exit 1; }
  rm -rf "$dir/bin" "$dir/work"
  echo "removed $dir/bin $dir/work"
  echo "evidence kept: $dir/evidence ($(ls "$dir/evidence" 2>/dev/null | wc -l | tr -d ' ') files)"
}

case "${1:-}" in
  build)   shift; cmd_build "$@" ;;
  doctor)  shift; cmd_doctor "$@" ;;
  run)     shift; cmd_run "$@" ;;
  cleanup) shift; cmd_cleanup "$@" ;;
  *)       usage ;;
esac
