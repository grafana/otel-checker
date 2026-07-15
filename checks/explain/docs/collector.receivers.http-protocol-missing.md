---
id: collector.receivers.http-protocol-missing
title: 'The OTLP receiver does not declare the http protocol'
severity: warning
---

## Why this matters

The OTLP receiver in the OpenTelemetry Collector can accept data over
gRPC, HTTP, or both — you opt in per protocol under
`receivers.otlp.protocols`. The Grafana Cloud SDK setup this checker
covers uses `http/protobuf` on the SDK side, so if the Collector's OTLP
receiver only enables `grpc` (or leaves `http` unset), incoming SDK
exports have nowhere to land: gRPC clients on the wrong port fail to
connect, HTTP clients on the wrong port get 404s.

## How to fix

Add an `http:` block under `receivers.otlp.protocols`. The block can be
empty (which uses the default endpoint `0.0.0.0:4318`) or override the
bind address / port:

```yaml
receivers:
  otlp:
    protocols:
      grpc:                 # keep gRPC if you also need it
        endpoint: 0.0.0.0:4317
      http:                 # add HTTP
        endpoint: 0.0.0.0:4318
```

Then reference the receiver in whichever pipelines should accept data
over OTLP:

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
```

## Example

Minimal receivers block with both protocols enabled:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
```

## Related

- [OTLP receiver documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/receiver/otlpreceiver)
- [OpenTelemetry Collector configuration](https://opentelemetry.io/docs/collector/configuration/)
