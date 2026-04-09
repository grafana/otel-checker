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

# Regenerate supported library lists from upstream OTel contrib repos
mise run generate
```

## Linting

```bash
# Auto-fix and verify (recommended dev workflow)
mise run lint:fix

# Verify only (same command used in CI)
mise run lint

# Go linting only
mise run lint:go
```

Linting is powered by [grafana/flint](https://github.com/grafana/flint) v2.

## Architecture

### Package Organization

- **`main.go`** — Entry point. Parses CLI args, runs checks, optional web server on `:8080`
- **`checks/checks.go`** — Orchestrator: always runs env checks first, then routes to component-specific checkers
- **`checks/env/`** — Common OTel environment variable validation
- **`checks/sdk/`** — Language-specific SDK checkers, each in its own subpackage (`go/`, `js/`, `java/`, `python/`, `dotnet/`, `rubyChecker.go`, `phpChecker.go`)
- **`checks/sdk/supported/`** — Shared library support checking logic
- **`checks/collector/`** — OTel Collector YAML config validation
- **`checks/beyla/`** — Beyla-specific checks
- **`checks/alloy/`** — Grafana Alloy checks
- **`checks/grafana/`** — Grafana Cloud connectivity/auth validation
- **`checks/utils/`** — CLI parsing, Reporter pattern (aggregates checks/warnings/errors)
- **`scripts/`** — Python scripts to generate `supported-libraries.yaml` from upstream OTel contrib repos
- **`static/`, `tmpl/`** — Embedded web UI assets

### Key Patterns

- **Reporter pattern**: `ComponentReporter` accumulates checks/warnings/errors, `Reporter` aggregates multiple component reporters
- **Generated files**: `supported-libraries.yaml` files in `checks/sdk/go/` and `checks/sdk/js/` — regenerate via `mise run generate`, don't edit manually
- **Embedded resources**: Static files and templates use `//go:embed`

## CLI Usage

```bash
# Required flags
otel-checker -language=<lang> -components=<comp1,comp2>

# Languages: dotnet, go, java, js, python, ruby, php
# Components: sdk, collector, beyla, alloy, grafana-cloud

# Examples
otel-checker -language=js -components=sdk
otel-checker -language=java -components=sdk,collector -manual-instrumentation
otel-checker -language=python -components=sdk,grafana-cloud -web-server
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
- Linting via flint v2 (shellcheck, shfmt, prettier, markdownlint, codespell, actionlint, editorconfig, lychee, renovate-deps, gofmt) + golangci-lint
- Python scripts use uv for dependencies
