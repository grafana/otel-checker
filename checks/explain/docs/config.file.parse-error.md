---
id: config.file.parse-error
title: 'The declarative config file could not be parsed'
severity: error
---

## Why this matters

The checker read the declarative config file but the YAML parser
refused it. Every downstream config check — file format, provider
presence, exporter endpoints — reads from the parsed document, so a
parse failure aborts config-side analysis entirely.

Common causes:

- Bad indentation. YAML is whitespace-significant and mixing tabs with
  spaces (or using different indent widths under the same key) is the
  most frequent trip-up.
- An unquoted value that YAML interprets as another type — a version
  string like `1.10` becomes a float, a `Yes`/`No` gets coerced to a
  boolean.
- An unterminated string, missing colon, or a stray character.
- Duplicate keys the parser rejects.

## How to fix

1. The finding message quotes the underlying parser error, which
   points at a specific line and column. Open the file at that
   position and inspect the surrounding indentation.

2. Validate the file with an external YAML parser:

   ```bash
   yq eval otel-config.yaml
   # or:
   python -c 'import yaml; yaml.safe_load(open("otel-config.yaml"))'
   ```

3. Quote version-like or boolean-like values:

   ```yaml
   file_format: "1.1"     # not: 1.1
   ```

## Related

- [OpenTelemetry Configuration spec](https://opentelemetry.io/docs/specs/otel/configuration/)
- `config.file.unreadable` — related failure when the file can't be
  opened at all.
