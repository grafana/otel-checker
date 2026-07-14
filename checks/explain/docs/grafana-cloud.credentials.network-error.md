---
id: grafana-cloud.credentials.network-error
title: 'Credential test could not reach the Grafana Cloud endpoint'
severity: error
---

## Why this matters

To validate credentials, `otel-checker` sends a real `POST` to
`<OTEL_EXPORTER_OTLP_ENDPOINT>/v1/metrics` with your `Authorization`
header. This error fires when that request never made it to a response —
DNS didn't resolve, TCP couldn't connect, TLS failed, the connection
timed out (10 second budget), or a proxy in the path refused the request.

Because the request never reached Grafana Cloud, the checker cannot say
whether the token itself is valid. What it *can* say is that your app
running in the same environment will hit the same failure and produce no
telemetry.

Common causes:

- The machine (or CI runner) can't reach `*.grafana.net` — a corporate
  proxy, firewall, or NetworkPolicy is blocking egress.
- A typo in `OTEL_EXPORTER_OTLP_ENDPOINT` — a stale region, wrong
  cluster name, or extra path segment produces a DNS or 404 close to
  the network layer.
- The endpoint uses `http://` instead of `https://`, and something in
  the path is redirecting or dropping the plaintext request.

## How to fix

1. Reproduce the request from the same environment to see the underlying
   error:

   ```bash
   curl -v -X POST "$OTEL_EXPORTER_OTLP_ENDPOINT/v1/metrics" \
     -H "Content-Type: application/x-protobuf" \
     -H "$(echo "$OTEL_EXPORTER_OTLP_HEADERS" | tr ',' '\n')"
   ```

   The verbose output distinguishes DNS from TCP from TLS from
   application-level failures.

2. If DNS or connect fails, check egress rules. Grafana Cloud endpoints
   are `*.grafana.net`; open outbound `443/tcp` to that hostname (or
   route it through your proxy).

3. If you're behind a proxy, export `HTTPS_PROXY` / `HTTP_PROXY` for both
   the checker and your app.

4. If the endpoint URL contains a typo, correct it (see
   `grafana-cloud.endpoint.invalid-format`).

## Example

Sanity check without the checker:

```bash
curl -o /dev/null -s -w "%{http_code}\n" \
  -X POST "$OTEL_EXPORTER_OTLP_ENDPOINT/v1/metrics" \
  -H "Authorization: Basic <base64-token>" \
  -H "Content-Type: application/x-protobuf"
# A 4xx / 5xx here still means the network worked. Anything else (curl error)
# is the same class of failure this ID reports.
```

## Related

- `grafana-cloud.endpoint.invalid-format` — related check when the URL
  itself is malformed (a common cause of network-layer failures).
- `grafana-cloud.credentials.unauthorized` — related error when the
  request *does* reach the gateway but is rejected as unauthenticated.
