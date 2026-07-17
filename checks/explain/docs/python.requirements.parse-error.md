---
id: python.requirements.parse-error
title: 'A line in requirements.txt could not be parsed'
severity: warning
---

## Why this matters

`otel-checker`'s Python parser understands lines of the form
`package==version` — an exact pin. It also handles `pip-compile
--generate-hashes` output: trailing `\` continuations are stripped and
`--hash=…` lines are skipped silently, so a hashed pin still parses as
its `name==version` value. Every other line in `requirements.txt`
(loose specifier, editable install, file reference, URL, environment
marker) is skipped, and this warning fires so you know which lines were
ignored.

The affected packages don't show up in the supported-libraries report,
even if OpenTelemetry has instrumentation for them.

## How to fix

1. If the flagged line is a genuine dependency you care about, pin it
   to an exact version:

   ```text
   # Before (skipped by the parser)
   requests>=2.28

   # After
   requests==2.32.3
   ```

   `pip freeze` produces the pinned form for every currently-installed
   package:

   ```bash
   pip freeze > requirements.txt
   ```

2. If the flagged line is one the parser will never handle (editable
   install, `-r` reference, URL), either leave it and accept that its
   packages won't appear in the report, or export a separate frozen
   file for the checker to consume. Hashed `pip-compile` output is
   handled — you don't need to strip `--hash=` lines yourself.

3. If your project uses another dependency manager, export a
   `requirements.txt` in the pinned form:

   ```bash
   poetry export -f requirements.txt --output requirements.txt --without-hashes
   uv pip compile pyproject.toml -o requirements.txt
   ```

## Example

Lines the parser handles vs. skips:

```text
opentelemetry-api==1.28.0         # ✓ parsed
django==4.2.11                    # ✓ parsed
click==8.1.7 \                    # ✓ parsed (pip-compile hashed pin)
    --hash=sha256:abc...          # ✓ silently skipped (hash line)

requests>=2.28                    # ✗ skipped (loose specifier)
-e .                              # ✗ skipped (editable install)
-r extra.txt                      # ✗ skipped (file reference)
git+https://github.com/x/y.git    # ✗ skipped (URL)
```

## Related

- `python.requirements.empty` — related warning when *no* lines could
  be parsed.
- `python.requirements.unreadable` — related error when the file itself
  can't be read.
