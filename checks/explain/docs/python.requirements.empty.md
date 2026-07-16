---
id: python.requirements.empty
title: 'No dependencies found in requirements.txt'
severity: warning
---

## Why this matters

`otel-checker` read `requirements.txt` successfully but couldn't extract
any dependencies from it — every line was either blank, a comment, or
in a form the parser doesn't recognize. That means the
supported-libraries report will be empty even though the check itself
ran to completion.

The parser accepts lines of the form `name==version` (an exact pin) and
handles `pip-compile --generate-hashes` output — trailing `\`
continuations are stripped and `--hash=…` lines are skipped silently, so
a hashed pin still parses as its `name==version` value.

Common shapes the parser does *not* recognize:

- Loose specifiers: `requests>=2.28`, `django~=4.2`, `flask` (bare name).
- Editable installs: `-e .`, `-e git+https://...`.
- File references: `-r requirements-dev.txt`.
- Comments: `# ...` (silently ignored, not an error).
- URL-only dependencies with no version pin.

## How to fix

1. If your `requirements.txt` uses loose specifiers, either pin exact
   versions (recommended for reproducible builds):

   ```bash
   pip freeze > requirements.txt
   ```

2. If your project uses another dependency manager (Poetry, PDM, Pipenv,
   uv), export a compatible `requirements.txt`:

   ```bash
   poetry export -f requirements.txt --output requirements.txt --without-hashes
   pipenv requirements > requirements.txt
   uv pip compile pyproject.toml -o requirements.txt
   ```

3. If your project genuinely has no dependencies yet, install at least
   the OpenTelemetry base packages and re-run:

   ```bash
   pip install opentelemetry-api opentelemetry-sdk opentelemetry-exporter-otlp
   pip freeze > requirements.txt
   ```

## Example

Frozen `requirements.txt` the checker parses cleanly:

```text
opentelemetry-api==1.28.0
opentelemetry-sdk==1.28.0
opentelemetry-exporter-otlp==1.28.0
django==4.2.11
requests==2.32.3
```

## Related

- `python.requirements.parse-error` — related warning when specific
  lines fail to parse.
- `python.requirements.unreadable` — related error when the file itself
  can't be read.
