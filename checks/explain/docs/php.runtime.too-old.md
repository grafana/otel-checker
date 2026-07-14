---
id: php.runtime.too-old
title: 'PHP version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` requires PHP 8.0 or newer. Everything below has reached
end-of-life
([endoflife.date/php](https://endoflife.date/php)) and the OpenTelemetry
PHP packages target 8.0+: they rely on typed properties, constructor
promotion, and other features that older PHP releases don't provide.
Installing them on PHP 7.x either fails at `composer install` or throws
`ParseError` at request time.

## How to fix

Upgrade to a supported PHP release:

- **macOS**: `brew install php` (installs the current stable release).
- **Debian / Ubuntu**: use Ondrej's PPA if the distro package is old:

  ```bash
  sudo add-apt-repository ppa:ondrej/php
  sudo apt-get update
  sudo apt-get install php8.3 php8.3-cli
  ```

- **Docker**: bump the base image, e.g. `FROM php:8.3-cli` or
  `FROM php:8.3-fpm`.
- **CI**: pin the PHP version explicitly (e.g. `shivammathur/setup-php`
  on GitHub Actions).

Verify:

```bash
php -v
# PHP 8.3.x ...
```

## Example

GitHub Actions workflow that pins PHP 8.3 before invoking the checker:

```yaml
- uses: shivammathur/setup-php@v2
  with:
    php-version: '8.3'
- run: otel-checker check sdk --language=php
```

## Related

- [PHP release schedule](https://www.php.net/supported-versions.php)
- [endoflife.date/php](https://endoflife.date/php) — current status of
  every PHP major.
- `php.runtime.not-found` — related failure when PHP isn't installed at
  all.
