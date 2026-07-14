---
id: php.composer-lock.missing
title: 'composer.lock is missing'
severity: error
---

## Why this matters

`composer.lock` records the exact resolved version of every package your
project uses — including transitive dependencies. `otel-checker` reads
it to see which OpenTelemetry packages are actually installed and at
what versions. Without a lock file, Composer hasn't resolved the graph
yet, and the checker can't tell whether `open-telemetry/api`,
`open-telemetry/sdk`, `open-telemetry/sem-conv`, and
`open-telemetry/exporter-otlp` are present.

## How to fix

Run `composer install` to resolve dependencies and generate the lock
file:

```bash
composer install
```

Commit `composer.lock` to the repo — it's not a build artifact, it's a
guarantee that every environment resolves the same versions. Running
`composer install` on a fresh checkout should be idempotent when the
lock file is present.

## Example

Clean install flow:

```bash
composer install
ls composer.lock         # exists now
otel-checker check sdk --language=php
```

## Related

- `php.composer-json.missing` — related failure when `composer.json`
  itself doesn't exist yet.
- `php.composer.not-found` — related failure when Composer isn't
  installed.
- [Composer basic usage: composer.lock](https://getcomposer.org/doc/01-basic-usage.md#commit-your-composer-lock-file-to-version-control)
