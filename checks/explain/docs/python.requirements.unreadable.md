---
id: python.requirements.unreadable
title: 'requirements.txt could not be read'
severity: error
---

## Why this matters

`otel-checker` looks for `requirements.txt` in the current directory and
parses it to build the list of your project's dependencies, then matches
each against the OpenTelemetry Python instrumentation catalog. If the
file exists but the read failed (permission denied, symlink to a
nonexistent target, mid-write), the supported-libraries check is
skipped — the main output of `check sdk --language=python`.

Note: if `requirements.txt` simply doesn't exist, the checker silently
falls through without producing any finding. This error is specifically
for a *read failure* on a file the checker found on disk.

## How to fix

1. Confirm the file is readable from the current directory:

   ```bash
   ls -l requirements.txt
   cat requirements.txt >/dev/null
   ```

2. If your project lives in a sub-directory, `cd` into it first:

   ```bash
   cd services/api
   otel-checker check sdk --language=python
   ```

3. If the file exists but the read failed with `permission denied`, fix
   the permission bits:

   ```bash
   chmod +r requirements.txt
   ```

4. If your project uses `pyproject.toml` / `poetry.lock` /
   `pipfile.lock` instead, generate a `requirements.txt` for the
   checker to consume:

   ```bash
   pip freeze > requirements.txt
   # or:
   poetry export -f requirements.txt --output requirements.txt
   ```

## Example

Minimal `requirements.txt`:

```text
opentelemetry-api==1.28.0
opentelemetry-sdk==1.28.0
opentelemetry-exporter-otlp==1.28.0
```

## Related

- `python.requirements.empty` — related warning when the file is
  readable but contains no parseable dependencies.
- `python.requirements.parse-error` — related warning when a specific
  line couldn't be parsed.
