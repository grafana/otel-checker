# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Project Overview

otel-checker ("OTel Me If It's Right") is a Go CLI tool that validates
OpenTelemetry instrumentation implementations. It scans code repositories,
checks environment variables, validates Grafana Cloud tokens, and verifies
correct SDK/collector/Beyla/Alloy configuration across 7 languages (Go, JS,
Java, Python, .NET, Ruby, PHP).

## Build Commands

```bash
# Build and install
mise run build

# Run all tests
mise run test

# Run all checks (lint + test)
mise run check

# Update dependencies
mise run deps

# Regenerate supported library lists from upstream OTel contrib repos.
# Expects sibling clones of opentelemetry-{go,js}-contrib by default;
# pass --clone-folders to point at a different parent directory.
mise run generate
mise run generate --clone-folders=/path/to/parent-of-clones
```

## Linting

```bash
# Auto-fix and verify (recommended dev workflow)
mise run lint:fix

# Verify only (same command used in CI)
mise run lint

```

Linting is powered by [grafana/flint](https://github.com/grafana/flint).

Run `mise run lint:fix` before committing changes.
If output includes `fixed`, keep those changes.
If output includes `partial` or `review`, address the remaining issues and
run `mise run lint:fix` again.

Example output:
flint: fixed: gofmt — commit before pushing | partial: cargo-clippy

## Architecture

### Package Organization

- **`cmd/otel-checker/`** — CLI entry point and cobra subcommand wiring
  (`check`, `serve`, `explain`, `version`, etc.)
- **`checks/checks.go`** — Orchestrator: always runs env checks first, then routes to component-specific checkers
- **`checks/env/`** — Common OTel environment variable validation
- **`checks/sdk/`** — Language-specific SDK checkers, each in its own
  subpackage (`go/`, `js/`, `java/`, `python/`, `dotnet/`,
  `rubyChecker.go`, `phpChecker.go`)
- **`checks/sdk/supported/`** — Shared library support checking logic
- **`checks/collector/`** — OTel Collector YAML config validation
- **`checks/beyla/`** — Beyla-specific checks (stub — no checks yet)
- **`checks/alloy/`** — Grafana Alloy checks (stub — no checks yet)
- **`checks/grafana/`** — Grafana Cloud connectivity/auth validation
- **`checks/utils/`** — `Commands` struct, flag validation, typed errors,
  `Reporter`/`ComponentReporter` pattern (aggregates checks/warnings/errors
  plus per-finding explain IDs)
- **`checks/output/`** — Pluggable result renderers (text/JSON/YAML) via a
  `Renderer` interface
- **`checks/explain/`** — Registry of explanation docs keyed by stable
  kebab-namespaced IDs (`docs/*.md`, embedded). `Lookup`/`All` are the
  public API
- **`checks/webserver/`** — Embedded web UI (`tmpl/`, `static/`) for
  `--web-server` and the `serve` subcommand
- **`scripts/`** — Python scripts to generate `supported-libraries.yaml` from upstream OTel contrib repos

### Key Patterns

- **Reporter pattern**: `ComponentReporter` accumulates checks/warnings/errors;
  `Reporter` aggregates multiple component reporters. Each finding carries
  an optional explain ID via the `AddXxxWithExplain` family of methods.
- **Explain IDs**: stable kebab-namespaced strings (e.g.
  `env.otel-service-name.unset`) that point at a markdown doc in
  `checks/explain/docs/`. The CLI's `explain <id>` subcommand and the web
  UI's `/explain/{id}` route both look them up via `explain.Lookup`. A
  coverage test asserts every literal ID in the source resolves to a
  registered doc.
- **Renderer interface**: `output.Renderer` lets callers plug in custom
  output. Text output orders findings as Errors → Warnings → Successful
  Checks, and appends `[explain.id]` to each line that has one.
- **Generated files**: `supported-libraries.yaml` files in `checks/sdk/go/`
  and `checks/sdk/js/` — regenerate via `mise run generate`, don't edit
  manually
- **Embedded resources**: Static files, templates, and the explain doc
  registry use `//go:embed`

## CLI Usage

```bash
# Per-component verbs
otel-checker check sdk           --language=<lang> [--manual-instrumentation ...]
otel-checker check collector     [--collector-config-path=<file>]
otel-checker check beyla         --language=<lang>
otel-checker check alloy         --language=<lang>
otel-checker check grafana-cloud --language=<lang>

# Multi-component in one invocation (positional, comma-separated, no spaces)
otel-checker check <comp1,comp2> --language=<lang>

# All components at once (no positional argument)
otel-checker check --language=<lang>

# Web UI replay of saved JSON results
otel-checker serve --data=results.json

# Explain a finding (looked up against checks/explain/docs/*.md)
otel-checker explain <id>            # single ID
otel-checker explain                 # every flagged ID in ./results.json
otel-checker explain list            # every registered ID

# Languages: dotnet, go, java, js, python, ruby, php
# Components: sdk, collector, beyla, alloy, grafana-cloud
# Output formats (--format): text (default), json, yaml

# Examples
otel-checker check sdk --language=js
otel-checker check sdk --language=java --manual-instrumentation
otel-checker check sdk,collector --language=js
otel-checker check --language=js                                 # every component
otel-checker check sdk --language=python --web-server --listen=127.0.0.1:9000
otel-checker check --language=js --format=json > results.json    # capture for explain/serve
```

## Code Conventions

- Go 1.24+
- Testing: stretchr/testify for assertions
- Color output: `fatih/color` (green/yellow/red)
- Each checker component in its own package
- Language checkers follow pattern `Check<Lang>Setup()`
- Errors reported via ComponentReporter, not panics

## CI

- `mise run check` (lint + test) on PRs
- Linting via flint (shellcheck, shfmt, rumdl, ryl, taplo, actionlint,
  typos, editorconfig-checker, golangci-lint, gofmt, ruff, ruff-format,
  biome, biome-format, lychee, renovate-deps)
- Python scripts use uv for dependencies
- Security issues should be reported via [Grafana's security issue reporting page](https://grafana.com/legal/report-a-security-issue/) and not directly in this repository.
