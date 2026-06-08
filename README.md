# OTel Me If It's Right

Checker for if the implementation of OpenTelemetry instrumentation is correct
by scanning the code in your repository, checking environment variables,
validating your Grafana token and more.

## Usage

Requirement: Golang

## Installation

1. Install the `otel-checker` binary

   ```text
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

2. You can confirm it was installed with:

   ```terminal
   ❯ ls $GOPATH/bin
   otel-checker
   ```

## Commands

```terminal
otel-checker check                  # all components
otel-checker check sdk              # SDK only
otel-checker check collector        # Collector config only
otel-checker check beyla            # Beyla only
otel-checker check alloy            # Grafana Alloy only
otel-checker check grafana-cloud    # Grafana Cloud connectivity only
otel-checker serve                  # web UI for a previously-saved JSON result
otel-checker version                # print the binary version
otel-checker completion <shell>     # generate shell completion script
```

The `check` command takes an optional comma-separated list of components
(`check sdk,collector,beyla`). With no argument, every component is checked.

Run `otel-checker check --help` or `otel-checker check sdk --help` for the full
flag set on each subcommand.

## Examples

```bash
# Single-component checks
otel-checker check sdk --language=js
otel-checker check sdk --language=java --manual-instrumentation
otel-checker check collector --collector-config-path=./otel/
otel-checker check grafana-cloud --language=python

# Multi-component (positional, comma-separated, no spaces)
otel-checker check sdk,collector,beyla --language=js

# Every component at once
otel-checker check --language=js
```

## Output formats

By default results are printed as colored text. Use `--format=json` or
`--format=yaml` for machine-readable output suitable for CI pipelines:

```bash
otel-checker check sdk --language=go --format=json
```

## Web UI

Pass `--web-server` to any `check` invocation to also serve the results at
`http://127.0.0.1:8080`. Override the bind address with `--listen=host:port`;
the default binds to loopback only. Press `Ctrl-C` to shut the server down
cleanly.

### Serving a results file

`otel-checker serve` watches a JSON or YAML results file on disk and renders it
in the web UI. The page polls every few seconds, so the moment the file is
created or rewritten the browser picks up the new content automatically.

By default, `serve` looks for `./results.json`, then `./results.yaml`, then
`./results.yml` in the current directory. If none exist yet, the server still
starts and shows a placeholder pointing at the expected path — useful for
keeping the UI open while a long-running pipeline writes results.

```bash
# Start the server (looks for ./results.json by default)
otel-checker serve

# In another terminal, write the file; the UI updates on the next poll
otel-checker check sdk --language=go --format=json > results.json

# Point at a specific file or a different format
otel-checker serve --data=./out/results.yaml
```

## Checks

### Common Environment Variables

These checks are automatically performed for all languages and components.

- Best practices for setting common environment variables:
  - Service name
  - Exporter protocol

- Resource attributes checks:
  - Validates the presence of recommended OpenTelemetry resource attributes
  - Checks for the following attributes:
    - `service.name` (via `OTEL_SERVICE_NAME` or in `OTEL_RESOURCE_ATTRIBUTES`)
    - `service.namespace` (e.g., `shop`)
    - `deployment.environment.name` (e.g., `production`)
    - `service.instance.id` (e.g., `checkout-123`)
    - `service.version` (e.g., `1.2`)
  - For missing attributes, provides specific recommendations with example values
  - Follows the
    [OpenTelemetry specification](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/)
    for precedence (e.g., `OTEL_SERVICE_NAME` takes precedence over
    `service.name` in `OTEL_RESOURCE_ATTRIBUTES`)
  - Example warning: `Set OTEL_RESOURCE_ATTRIBUTES="service.namespace=shop": An optional namespace for service.name`

### Grafana Cloud

Run `otel-checker check grafana-cloud --language=<lang>` (or pass
`--components=grafana-cloud` to `check`):

- Endpoints
- Authentication

### SDK

#### JavaScript

Run `otel-checker check sdk --language=js`:

- Node version
- Required dependencies on package.json
- Required environment variables
- Resource detectors
- Dependencies compatible with Grafana Cloud
- Usage of Console Exporter
- Prints which libraries are supported based on the `package.json` in the current directory.

#### Python

Run `otel-checker check sdk --language=python`:

- Prints which libraries are supported:
  - The used libraries are discovered from `requirements.txt` in the current directory.

#### .NET

Run `otel-checker check sdk --language=dotnet`:

- .NET version
- Available instrumentation for .NET libraries and dependencies
- Auto-instrumentation environment variables

> [!NOTE]
> Only .NET 8.0 and higher are supported

#### Java

Run `otel-checker check sdk --language=java`:

- Java version
- Prints which libraries (as discovered from a locally running maven or gradle)
  are supported:
  - With `--manual-instrumentation`, the libraries for manual instrumentation are printed.
  - Without `--manual-instrumentation`, it will print the libraries supported by
    the [Java Agent](https://github.com/open-telemetry/opentelemetry-java-instrumentation/).
  - A maven or gradle wrapper will be used if found in the current directory or
    a parent directory.

#### Go

Run `otel-checker check sdk --language=go`:

- Prints which libraries are supported for manual instrumentation
  based on the `go.mod` in the current directory.

#### Ruby

Run `otel-checker check sdk --language=ruby`:

- Ruby version
- Bundler installation
- `Gemfile` and `Gemfile.lock` exist
- Required dependencies installed
- Optional auto-instrumentation dependencies installed

#### PHP

Run `otel-checker check sdk --language=php`:

- PHP version
- Composer installation
- `composer.json` and `composer.lock` exist
- Required dependencies in `composer.lock`
- Some auto-instrumentation dependencies installed

### Collector

Run `otel-checker check collector`:

- Config receivers and exporters

### Beyla

Run `otel-checker check beyla --language=<lang>`:

- Environment variables

### Alloy

> [!NOTE]
> TBD
