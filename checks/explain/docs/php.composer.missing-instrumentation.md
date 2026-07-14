---
id: php.composer.missing-instrumentation
title: 'No OpenTelemetry instrumentation package is installed'
severity: error
---

## Why this matters

Auto-instrumentation for a PHP application is delivered as a set of
per-framework Composer packages —
`open-telemetry/opentelemetry-auto-symfony`,
`open-telemetry/opentelemetry-auto-laravel`,
`open-telemetry/opentelemetry-auto-guzzle`, and so on. `otel-checker`'s
auto-instrumentation mode expects at least one of these to be pinned in
`composer.lock`.

Without any instrumentation package, the SDK from `open-telemetry/sdk`
still initializes but nothing hooks into the frameworks and libraries
in your app — no request spans, no ORM spans, no HTTP-client spans —
even though the tracer provider itself is configured correctly.

Note that PHP auto-instrumentation also requires the
`open-telemetry/opentelemetry` PHP extension to be loaded in the SDK's
`.ini` — the Composer package alone isn't enough. Installing that
extension is outside this check's scope but part of any working
auto-instrumentation setup.

## How to fix

Install the instrumentation package(s) that match your stack. Common
picks:

```bash
# Framework
composer require open-telemetry/opentelemetry-auto-symfony
composer require open-telemetry/opentelemetry-auto-laravel
composer require open-telemetry/opentelemetry-auto-slim
composer require open-telemetry/opentelemetry-auto-yii

# CMS
composer require open-telemetry/opentelemetry-auto-wordpress

# HTTP / DB clients
composer require open-telemetry/opentelemetry-auto-guzzle
composer require open-telemetry/opentelemetry-auto-pdo
composer require open-telemetry/opentelemetry-auto-mongodb

# PSR standards (any that apply)
composer require open-telemetry/opentelemetry-auto-psr3   # logs
composer require open-telemetry/opentelemetry-auto-psr15  # middleware
composer require open-telemetry/opentelemetry-auto-psr18  # http-client
```

You also need the PHP extension loaded:

```bash
pecl install opentelemetry
echo "extension=opentelemetry.so" >> /path/to/php.ini
```

Verify:

```bash
php -m | grep opentelemetry
composer show | grep opentelemetry-auto
```

## Example

`composer.json` fragment for a Symfony service with HTTP + database
instrumentation:

```json
{
  "require": {
    "open-telemetry/api": "^1.0",
    "open-telemetry/sdk": "^1.0",
    "open-telemetry/exporter-otlp": "^1.0",
    "open-telemetry/opentelemetry-auto-symfony": "^1.0",
    "open-telemetry/opentelemetry-auto-pdo": "^1.0",
    "open-telemetry/opentelemetry-auto-guzzle": "^1.0"
  }
}
```

## Related

- [OpenTelemetry PHP: zero-code instrumentation](https://opentelemetry.io/docs/zero-code/php/)
- [OpenTelemetry PHP contrib repositories](https://github.com/open-telemetry?q=opentelemetry-php&type=all)
- `php.composer.missing-required` — related error for the four core
  packages every setup needs.
