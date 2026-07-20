# OTel Me If It's Right

Checker for if the implementation of OpenTelemetry instrumentation is correct
by scanning the code in your repository, checking environment variables,
validating your Grafana token and more.

## Usage

For source installs, you need Go 1.25 or higher. Prebuilt binaries are also
available from the [GitHub Releases](https://github.com/grafana/otel-checker/releases)
page.

## Installation

### Prebuilt binaries

Download the archive for your operating system and CPU architecture from the
[latest release](https://github.com/grafana/otel-checker/releases/latest),
extract it, and put the `otel-checker` binary on your `PATH`. Release archives
are provided for Linux, macOS, and Windows on amd64 and arm64. The statically
linked Linux archives work on both glibc- and musl-based distributions, such as
Alpine Linux. `checksums.txt` contains SHA-256 checksums for verification.

### Go install

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
otel-checker explain                    # show explanations for every finding from a saved results file
otel-checker explain <id>               # show the explanation for a single ID
otel-checker explain list               # list every available explain ID
otel-checker version                    # print the binary version
otel-checker completion <shell>         # generate shell completion script
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

## Explanations

Each actionable finding (errors and warnings) is tagged with a stable explain ID
shown in square brackets at the end of the line:

```text
✖ SDK: package.json missing on path /src/inst [js.package-json.unreadable]
```

Look up the guidance for any finding with:

```bash
otel-checker explain js.package-json.unreadable
```

Or enumerate every available explain ID:

```bash
otel-checker explain list
```

To see explanations for **every** finding from a previous run, save the JSON output
once and call `explain` with no ID:

```bash
otel-checker check --language=js --format=json > results.json
otel-checker explain
```

With no ID, `explain` reads `./results.json` by default (also `./results.yaml` /
`./results.yml`) — the same file `serve` watches — and prints each explain doc
in turn, deduplicating repeated IDs and skipping findings that have none.
Pass `--data=<path>` to point at a different file.

The same IDs also appear as a `explain_id` field on every entry when using
`--format=json` or `--format=yaml`, so downstream tooling can match findings
against the explain catalog programmatically.

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

- Signal exporters (`OTEL_TRACES_EXPORTER`, `OTEL_METRICS_EXPORTER`,
  `OTEL_LOGS_EXPORTER`): must be `otlp`, `console`, or unset; `none` is rejected.
- Service name: `OTEL_SERVICE_NAME` (or `service.name` in
  `OTEL_RESOURCE_ATTRIBUTES`).
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

- `OTEL_EXPORTER_OTLP_ENDPOINT` matches
  `https://otlp-gateway-<zone>.grafana.net/otlp`.
- `OTEL_EXPORTER_OTLP_PROTOCOL` is `http/protobuf`.
- `OTEL_EXPORTER_OTLP_HEADERS` contains an `Authorization=Basic <token>` entry.
- Credential validation: tests the endpoint with the provided token.

### SDK

#### JavaScript

Run `otel-checker check sdk --language=js`:

- Node version — must be an Active LTS release (even-major, currently 22,
  24, or 26). Odd-major Current releases (23, 25) are flagged as well.
- `@opentelemetry/api` dependency in `package.json`.
- Auto-instrumentation mode (default):
  - `@opentelemetry/auto-instrumentations-node` dependency.
  - `NODE_OPTIONS` includes
    `--require @opentelemetry/auto-instrumentations-node/register`.
  - `OTEL_NODE_RESOURCE_DETECTORS` set to `all` or covers
    `env,host,os,serviceinstance`.
- Manual-instrumentation mode (`--manual-instrumentation`):
  - No concurrent auto-instrumentation via `NODE_OPTIONS`.
  - Uses `@opentelemetry/exporter-trace-otlp-http` (not `-otlp-proto`).
  - Warns when `ConsoleSpanExporter` / `ConsoleMetricExporter` is used
    outside debugging.
- Supported libraries: `package.json` dependencies are matched against the
  OpenTelemetry JS contrib registry.

#### Python

Run `otel-checker check sdk --language=python`:

- Supported libraries: `requirements.txt` dependencies are matched against the
  OpenTelemetry Python contrib registry.

#### .NET

Run `otel-checker check sdk --language=dotnet`:

- .NET version (>= 8.0).
- Auto-instrumentation environment variables:
  - `CORECLR_ENABLE_PROFILING` = `1`
  - `CORECLR_PROFILER` = `{918728DD-259F-4A6A-AC2B-B85E1B658318}`
  - `CORECLR_PROFILER_PATH` is set
  - `OTEL_DOTNET_AUTO_HOME` is set
- Supported instrumentations: `.csproj` NuGet dependencies are matched against
  the OpenTelemetry .NET auto-instrumentation library list.

> [!NOTE]
> Only .NET 8.0 and higher are supported

#### Java

Run `otel-checker check sdk --language=java`:

- Java version (>= 8).
- Supported libraries (discovered from a locally running Maven or Gradle,
  including a wrapper if found in the current directory or a parent):
  - Without `--manual-instrumentation`: libraries supported by the
    [Java Agent](https://github.com/open-telemetry/opentelemetry-java-instrumentation/).
  - With `--manual-instrumentation`: libraries with manual instrumentation
    support.

#### Go

Run `otel-checker check sdk --language=go`:

- Supported libraries for manual instrumentation, based on `go.mod` in the
  current directory.

#### Ruby

Run `otel-checker check sdk --language=ruby`:

- Ruby version (CRuby >= 3.3, JRuby >= 9.4, or TruffleRuby >= 22.1).
- Bundler installed.
- `Gemfile` and `Gemfile.lock` exist.
- Required gems: `opentelemetry-api`, `opentelemetry-sdk`,
  `opentelemetry-exporter-otlp`.
- Auto-instrumentation (default): `opentelemetry-instrumentation-all` or at
  least one specific instrumentation gem (e.g. `-rack`, `-rails`).

#### PHP

Run `otel-checker check sdk --language=php`:

- PHP version (>= 8.0).
- Composer installed.
- `composer.json` and `composer.lock` exist.
- Required packages: `open-telemetry/api`, `open-telemetry/sem-conv`,
  `open-telemetry/sdk`, `open-telemetry/exporter-otlp`.
- Auto-instrumentation (default): at least one instrumentation package (e.g.
  `symfony`, `pdo`, `laravel`, `wordpress`, `guzzle`).

### Collector

Run `otel-checker check collector`:

- `config.yaml` exists and is readable, and parses as valid YAML.
- At least one `otlp` receiver has the `http` protocol configured.
- At least one `otlphttp` / `otlp_http` exporter has an endpoint matching the
  Grafana Cloud format (`https://*.grafana.net/otlp`); warns when set to
  `localhost`.
- For each of the `traces`, `metrics`, and `logs` pipelines: the exporter list
  contains an `otlphttp` / `otlp_http` exporter, and the receiver list
  contains an `otlp` receiver.

Named components (e.g. `otlphttp/grafana_cloud`, `otlp/app`) are supported.

### Beyla

> [!NOTE]
> TBD

### Alloy

> [!NOTE]
> TBD
