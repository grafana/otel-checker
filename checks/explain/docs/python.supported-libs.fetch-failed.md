---
id: python.supported-libs.fetch-failed
title: 'Could not fetch the supported Python libraries list'
severity: error
---

## Why this matters

The Python checker downloads the OpenTelemetry Python contrib README
from GitHub at check time
(`https://raw.githubusercontent.com/open-telemetry/opentelemetry-python-contrib/refs/heads/main/instrumentation/README.md`)
and parses the supported-libraries table out of it. If that fetch fails
or the parse errors, the checker has no catalog to compare your
`requirements.txt` against and skips the supported-libraries report —
the main output of `check sdk --language=python`.

Common causes:

- **No network access.** The environment running `otel-checker` can't
  reach `raw.githubusercontent.com`. Common in air-gapped CI, networks
  with strict egress rules, or an unreachable proxy.
- **Upstream format drift.** The Python contrib README structure
  changed in a way the parser doesn't handle. The fetch itself
  succeeded but the extraction failed.

## How to fix

1. Verify network reachability from the environment running the checker:

   ```bash
   curl -sI https://raw.githubusercontent.com/open-telemetry/opentelemetry-python-contrib/refs/heads/main/instrumentation/README.md
   ```

   Anything other than `HTTP/2 200` is the underlying failure.

2. If you're behind a proxy, export `HTTPS_PROXY` for the checker:

   ```bash
   export HTTPS_PROXY=http://proxy.internal:3128
   ```

3. If the URL itself returns 200 but the checker still fails, the
   upstream format may have drifted. Upgrade `otel-checker`:

   ```bash
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

4. If it still fails after the upgrade, open an issue on
   `grafana/otel-checker` with the error message.

## Example

Quick reachability check:

```bash
curl -o /dev/null -s -w "%{http_code}\n" \
  https://raw.githubusercontent.com/open-telemetry/opentelemetry-python-contrib/refs/heads/main/instrumentation/README.md
# 200
```

## Related

- [opentelemetry-python-contrib](https://github.com/open-telemetry/opentelemetry-python-contrib) —
  upstream source of the supported-libraries list.
- `python.requirements.unreadable` / `python.requirements.empty` —
  companion checks on your project's dependency file.
