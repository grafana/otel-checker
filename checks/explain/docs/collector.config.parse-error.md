---
id: collector.config.parse-error
title: 'The collector config.yaml could not be parsed'
severity: error
---

## Why this matters

The checker read the Collector's `config.yaml` but the YAML parser
refused it. Every downstream collector check — endpoint format,
receiver protocols, per-pipeline exporter presence — reads from the
parsed document, so a parse failure aborts collector analysis
entirely.

Common causes:

- Bad indentation. YAML is whitespace-significant and mixing tabs with
  spaces (or using different indent widths under the same key) is the
  most frequent trip-up.
- An unquoted value that YAML interprets as another type — a version
  string like `1.10` becomes a float and drops the trailing zero, an
  IP address like `10:00` becomes a sexagesimal number, a `Yes`/`No`
  gets coerced to a boolean.
- An unterminated string, missing colon, or a stray character.
- Comments or duplicate keys the parser rejects.

## How to fix

1. The finding message quotes the underlying parser error, which points
   at a specific line and column. Open the file at that position and
   inspect the surrounding indentation.

2. Validate the file with an external YAML parser to double-check:

   ```bash
   yq eval config.yaml
   # or, if yq isn't available:
   python -c 'import yaml,sys; yaml.safe_load(open("config.yaml"))'
   ```

3. Ask the Collector itself to validate the config (it does a stricter
   schema check than a generic YAML parser):

   ```bash
   otelcol validate --config=./config.yaml
   ```

4. Common quick fixes:

   - Convert tabs to spaces.
   - Quote version-like values: `"1.10"` instead of `1.10`.
   - Escape colons inside unquoted values, or quote the whole string.

## Example

A version string coerced silently — parses but the Collector then
misbehaves on startup:

```yaml
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    version: 1.10          # becomes 1.1 as a float
```

Fix:

```yaml
    version: "1.10"        # stays as a string
```

## Related

- [OpenTelemetry Collector config reference](https://opentelemetry.io/docs/collector/configuration/)
- `collector.config.unreadable` — related failure when the file can't
  be opened at all.
