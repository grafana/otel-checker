---
id: php.composer.not-found
title: 'Composer is not installed'
severity: error
---

## Why this matters

`otel-checker` runs `composer -v` to confirm Composer is available.
Composer is the standard dependency manager for PHP — the OpenTelemetry
PHP packages ship as Composer libraries, and `composer install`
generates the `composer.lock` file the checker reads to verify which
packages are actually pinned.

## How to fix

Install Composer:

- **Official installer**:

  ```bash
  curl -sS https://getcomposer.org/installer | php
  sudo mv composer.phar /usr/local/bin/composer
  ```

- **macOS**:

  ```bash
  brew install composer
  ```

- **Debian / Ubuntu**:

  ```bash
  sudo apt-get install composer
  ```

- **Docker**: copy from the official Composer image, or install in your
  own image:

  ```dockerfile
  COPY --from=composer:latest /usr/bin/composer /usr/local/bin/composer
  ```

Verify:

```bash
composer -V
# Composer version 2.x ...
```

## Example

Bootstrap flow after a fresh install:

```bash
curl -sS https://getcomposer.org/installer | php -- \
  --install-dir=/usr/local/bin --filename=composer
composer install
otel-checker check sdk --language=php
```

## Related

- [Composer installation](https://getcomposer.org/download/)
- `php.runtime.not-found` — related failure when PHP itself isn't
  installed.
- `php.composer-lock.missing` — related failure when Composer is
  present but hasn't produced a lock file yet.
