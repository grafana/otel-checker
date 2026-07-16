---
id: php.runtime.too-old
title: 'PHP version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` requires PHP 8.0 or newer. Earlier versions have reached
end-of-life
([endoflife.date/php](https://endoflife.date/php)) and the OpenTelemetry
PHP packages target 8.0+: they rely on typed properties, constructor
promotion, and other features that older PHP releases don't provide.
Installing them on PHP 7.x either fails at `composer install` or throws
`ParseError` at runtime.

## How to fix

Upgrade to a supported PHP release:

- **macOS**: `brew install php` (installs the current stable release).
- **Debian / Ubuntu**: `sudo apt-get install php-cli`. This works on
  current distro releases (Ubuntu 22.04+, Debian 12+). On older
  releases the distro package predates PHP 8.0 — prefer Docker (see
  below) or install a supported version another way.
- **Docker**: use the [official PHP image](https://hub.docker.com/_/php)
  maintained by the PHP team, e.g. `FROM php:8.3-cli` or
  `FROM php:8.3-fpm`.
- **CI**: run the job inside the official PHP Docker image so the
  version is pinned by the image tag, not by the runner's default PHP.
  On GitHub Actions:

  ```yaml
  jobs:
    test:
      runs-on: ubuntu-latest
      container: php:8.3-cli
      steps:
        - uses: actions/checkout@v4
        - run: otel-checker check sdk --language=php
  ```

  GitLab / CircleCI have equivalent `image:` / `docker:` fields.

Verify:

```bash
php -v
# PHP 8.3.x ...
```

## Example

GitHub Actions workflow that pins PHP 8.3 via the official PHP Docker
image:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    container: php:8.3-cli
    steps:
      - uses: actions/checkout@v4
      - run: otel-checker check sdk --language=php
```

## Related

- [PHP release schedule](https://www.php.net/supported-versions.php)
- [endoflife.date/php](https://endoflife.date/php) — current status of
  every PHP major.
- `php.runtime.not-found` — related failure when PHP isn't installed at
  all.
