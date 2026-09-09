---
id: config.file.unreadable
title: 'The declarative config file could not be found'
severity: error
---

## Why this matters

`otel-checker check config` reads the OpenTelemetry declarative
configuration file — the language-agnostic YAML defined by the OTel
Configuration specification
([spec](https://opentelemetry.io/docs/specs/otel/configuration/)) — from the
full path passed via `--config-path`. When the flag is omitted, the
checker falls back to `otel-config.yaml` (then `otel-config.yml`) in
the current working directory. Every downstream config check — file
format, provider presence, exporter endpoints — reads from that file.
If it can't be opened, the checker has nothing to inspect and skips
the whole config-side analysis.

## How to fix

- Pass the exact path with `--config-path`:

  ```bash
  otel-checker check config --config-path=./deploy/otel-config.yaml
  ```

- Or run the checker from the directory that already contains
  `otel-config.yaml` / `otel-config.yml`.

- Verify the file exists and is readable by the current user:

  ```bash
  ls -l ./otel-config.yaml
  ```

- If the customer's SDK reads its configuration from a different
  filename (some SDKs use `OTEL_EXPERIMENTAL_CONFIG_FILE` to point at
  an arbitrary path), point `--config-path` at the same file so the
  checker validates what the application actually loads.

## Related

- `config.file.parse-error` — related failure when the file exists but
  the YAML parser rejects it.
- `collector.config.unreadable` — sibling check for the Collector's
  own config file.
