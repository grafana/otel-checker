---
id: php.composer.missing-required
title: 'A required OpenTelemetry package is missing from composer.lock'
severity: error
---

## Why this matters

The checker looks in `composer.lock` for four packages that any
OpenTelemetry PHP setup needs regardless of whether you use auto- or
manual-instrumentation:

- `open-telemetry/api` — the tracer / meter interfaces your code (and
  instrumentation packages) call into.
- `open-telemetry/sem-conv` — semantic-convention attribute keys used
  by the API and every instrumentation.
- `open-telemetry/sdk` — the SDK that actually assembles spans, metrics,
  and logs.
- `open-telemetry/exporter-otlp` — the exporter that ships those
  signals to Grafana Cloud via OTLP.

If any of them are missing, either the SDK never registers a tracer
provider (so instrumentation packages produce no-op spans), or spans
are produced but there's no exporter to ship them anywhere.

## How to fix

Add the missing package to `composer.json` and run `composer install`.
The finding names the specific package in its message.

```json
{
  "require": {
    "php": ">=8.0",
    "open-telemetry/api": "^1.0",
    "open-telemetry/sem-conv": "^1.0",
    "open-telemetry/sdk": "^1.0",
    "open-telemetry/exporter-otlp": "^1.0"
  }
}
```

Then:

```bash
composer install
otel-checker check sdk --language=php
```

Or install in one step:

```bash
composer require \
  open-telemetry/api \
  open-telemetry/sem-conv \
  open-telemetry/sdk \
  open-telemetry/exporter-otlp
```

## Example

Minimal instrumentation bootstrap code:

```php
use OpenTelemetry\SDK\Sdk;
use OpenTelemetry\SDK\Trace\TracerProvider;

$tracerProvider = TracerProvider::builder()
    ->addSpanProcessor(/* OTLP span processor */)
    ->build();

Sdk::builder()
    ->setTracerProvider($tracerProvider)
    ->setAutoShutdown(true)
    ->buildAndRegisterGlobal();
```

## Related

- `php.composer.missing-instrumentation` — related error when
  instrumentation packages (e.g. `opentelemetry-auto-symfony`) are
  missing.
- [OpenTelemetry PHP: getting started](https://opentelemetry.io/docs/languages/php/getting-started/)
