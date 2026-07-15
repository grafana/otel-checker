---
id: java.supported-libs.fetch-failed
title: 'Could not fetch the supported Java libraries list'
severity: error
---

## Why this matters

The Java checker downloads the OpenTelemetry Java auto-instrumentation
catalog from GitHub at check time
(`https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml`)
so it can match your Maven / Gradle dependencies against the current
set of instrumented libraries. If that fetch fails, the checker has no
catalog to compare against and skips the supported-libraries report —
the most useful output of `check sdk --language=java`.

Common causes:

- **No network access.** The environment running `otel-checker` can't
  reach `raw.githubusercontent.com`. Common in air-gapped CI, restrictive
  corporate networks, or a broken proxy config.
- **GitHub-side outage.** Occasionally `raw.githubusercontent.com`
  returns 5xx errors or times out.
- **YAML shape drift.** If the upstream file schema changes in a way
  the checker doesn't handle, the fetch itself succeeds but the parse
  fails and gets reported under the same ID. Rare, but worth checking
  if the URL itself works.

## How to fix

1. Verify network reachability from the environment running the checker:

   ```bash
   curl -sI https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml
   ```

   Anything other than `HTTP/2 200` is the underlying failure. If DNS or
   TCP fails, open outbound `443/tcp` to `raw.githubusercontent.com`
   (or route it through your proxy).

2. If you're behind a proxy, export `HTTPS_PROXY` for the checker:

   ```bash
   export HTTPS_PROXY=http://proxy.internal:3128
   ```

3. If the URL itself returns 200 but the checker still fails, the schema
   likely drifted. Upgrade `otel-checker`:

   ```bash
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

## Example

Quick reachability check:

```bash
curl -o /dev/null -s -w "%{http_code}\n" \
  https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml
# 200
```

## Related

- [OpenTelemetry Java auto-instrumentation library list](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/supported-libraries.md)
- `java.gradle.no-dependencies` — related warning when the local
  dependency list itself is empty.
