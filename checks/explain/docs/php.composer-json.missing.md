---
id: php.composer-json.missing
title: 'composer.json is missing'
severity: error
---

## Why this matters

`composer.json` is where you declare the packages your PHP project
depends on — including the OpenTelemetry packages the checker verifies.
Without one, Composer can't resolve dependencies, `composer.lock`
won't exist, and the checker has nothing to inspect.

## How to fix

1. If you're running `otel-checker` from the wrong directory (a
   monorepo root, a `bin/` folder), `cd` into the PHP project
   directory:

   ```bash
   cd services/api
   otel-checker check sdk --language=php
   ```

2. If the project has no `composer.json` yet, create one interactively:

   ```bash
   composer init
   ```

   Or write a minimal file by hand — see the Example below — then run
   `composer install` to generate `composer.lock`.

## Example

Minimal `composer.json` for a service that ships OpenTelemetry via
OTLP:

```json
{
  "name": "acme/api",
  "require": {
    "php": ">=8.0",
    "open-telemetry/api": "^1.0",
    "open-telemetry/sdk": "^1.0",
    "open-telemetry/exporter-otlp": "^1.0",
    "open-telemetry/sem-conv": "^1.0"
  }
}
```

Then:

```bash
composer install
otel-checker check sdk --language=php
```

## Related

- `php.composer-lock.missing` — related failure once `composer.json`
  exists but hasn't been resolved yet.
- `php.composer.missing-required` — companion check that verifies the
  OpenTelemetry base packages are pinned once the lock file is
  generated.
