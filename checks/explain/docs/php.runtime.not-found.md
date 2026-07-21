---
id: php.runtime.not-found
title: 'PHP is not installed'
severity: error
---

## Why this matters

`otel-checker` runs `php -v` to confirm a PHP interpreter is available.
This error means no `php` binary is on the current `PATH`. Without a
PHP interpreter, none of the downstream checks (version, Composer,
composer.lock) can run, and your service can't execute at all.

## How to fix

Install PHP 8.0 or newer:

- **macOS**:

  ```bash
  brew install php
  ```

- **Debian / Ubuntu**:

  ```bash
  sudo apt-get install php php-cli
  ```

- **RHEL / Fedora**:

  ```bash
  sudo dnf install php-cli
  ```

- **Docker**: use an official image, e.g. `FROM php:8.3-cli`.
- **Windows**: install from [windows.php.net/download](https://windows.php.net/download).

Verify:

```bash
which php
php -v
# PHP 8.3.x ...
```

## Example

Base image for a PHP service with the OpenTelemetry extension available:

```dockerfile
FROM php:8.3-cli
RUN apt-get update && apt-get install -y --no-install-recommends \
      unzip curl && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
```

## Related

- `php.runtime.too-old` — related failure when PHP is found but its
  version is below the minimum.
- `php.composer.not-found` — related failure for the companion tool.
- [PHP downloads](https://www.php.net/downloads)
