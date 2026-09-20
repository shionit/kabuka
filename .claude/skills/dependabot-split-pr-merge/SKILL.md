---
name: dependabot-split-pr-merge
description: In this repo, detect "split Dependabot PRs" — cases where Dependabot opens separate PRs for dependencies that are tightly coupled within the same file (e.g. github/codeql-action/init and github/codeql-action/analyze pinned to the same version in the same job in .github/workflows/codeql-analysis.yml), where merging only one of them breaks CI due to a version mismatch — and use local git/gh CLI to create a combined branch and combined PR. The cloud GitHub App integration has no write access to this repo, so this skill relies on local git credentials (SSH key / gh CLI login) instead. Trigger on requests like "merge the dependabot PRs together", "consolidate the split PRs", "combine the codeql-action PRs", or "prepare the related PRs for merging together".
---

# Dependabot Split-PR Merge Skill

## Background

Dependabot opens one PR per dependency. But when multiple actions/versions are tightly coupled within the same file, merging only one of the resulting PRs leaves the others on an old version, and the version mismatch breaks CI.

Concrete example: in `.github/workflows/codeql-analysis.yml`, `github/codeql-action/init` and `github/codeql-action/analyze` are pinned to the same commit SHA within the same job. Dependabot creates two separate PRs for these, but merging only one leaves init and analyze on mismatched versions, and the "Analyze (go)" job fails.

## Prerequisites

- Runs against the local clone of kabuka.
- `gh` CLI is installed and authenticated (verify with `gh auth status`).
- Push access to `origin` is available.

## Procedure

### 1. Identify the candidate PRs

```bash
gh pr list --repo shionit/kabuka --search "author:app/dependabot" --state open \
  --json number,title,headRefName,files
```

### 2. Determine whether this is actually a split-PR situation

Check the diff of each PR:

```bash
gh pr diff <PR number> --repo shionit/kabuka
```

Treat it as a likely split PR if any of the following hold:

- Multiple open Dependabot PRs modify nearby lines in the same file (the same job / the same group of steps).
- Opening the target file shows that several places are designed to reference the same version/SHA (the setup breaks unless they stay in sync).
- The PR's CI log (`gh pr checks <PR number> --repo shionit/kabuka`) is actually failing, and the failure reason suggests a version mismatch, type mismatch, or interface incompatibility.

If it's unclear, don't force the "split PR" label — fall back to the normal flow and report it to the user individually as held/needs-review instead.

### 3. Do the merge work locally

```bash
git fetch origin
git checkout main
git pull origin main
git checkout -b dependabot-merge/<short-name-for-what-youre-combining>
```

Look at the diff of each Dependabot PR (`gh pr diff <PR number>`) and apply the changes to the target file by hand. In most cases this is a simple version string / commit SHA substitution.

If a lockfile needs to be regenerated (e.g. `go.sum`, `package-lock.json`), use the normal local command to bring it back into sync.

```bash
# Go example
go mod tidy

# Node.js example
npm install
```

If lockfile regeneration doesn't go cleanly, or the tooling can't handle it, don't force it — follow the "when this is hard to do" step below instead.

### 4. Commit and push

```bash
git add <changed files>
git commit -m "chore(deps): bump <target> together to keep versions in sync

Combines:
- #<PR number 1>
- #<PR number 2>"
git push origin dependabot-merge/<short-name-for-what-youre-combining>
```

### 5. Open the combined PR

```bash
gh pr create --repo shionit/kabuka \
  --title "chore(deps): bump <target> to <new version> (merged split Dependabot PRs)" \
  --body "Combines the following Dependabot PRs:
- #<PR number 1>
- #<PR number 2>

Reason: these are tightly coupled within the same workflow/config file, and merging them individually breaks CI due to a version mismatch." \
  --base main --head dependabot-merge/<short-name-for-what-youre-combining>
```

### 6. Check CI and report back

```bash
gh pr checks <new PR number> --repo shionit/kabuka --watch
```

- **If CI passes**: report this to the user. The user merges the PR themselves (`gh pr merge`) — this skill never clicks merge.
- **If CI fails**: stop the merge work for this repo here, and report the failure and the likely cause to the user. Do not attempt any further automated fix.
- **If something (e.g. lockfile regeneration) can't reasonably be done with CLI operations alone**: don't force it — report the situation and ask the user to confirm how to proceed.

## Things this skill must never do

- Never actually merge a PR (`gh pr merge`). Merging is always done by the user.
- Never close or edit the original Dependabot PRs on your own.
- If anything unexpected comes up (it's unclear whether this really is a split-PR case, the CI failure looks unrelated to a version mismatch, etc.), don't decide mergeability on your own — report the situation and ask the user to confirm.
