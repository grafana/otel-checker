---
id: ruby.gemfile.missing
title: 'Gemfile is missing'
severity: error
---

## Why this matters

A `Gemfile` in the current directory is where you declare which gems
your Ruby project depends on — including the OpenTelemetry gems the
checker needs to verify are present. Without one, Bundler can't resolve
dependencies, `Gemfile.lock` won't exist, and the checker has nothing
to inspect.

## How to fix

1. If you're running `otel-checker` from the wrong directory (a monorepo
   root, a `bin/` folder), `cd` into the Ruby project directory:

   ```bash
   cd services/api
   otel-checker check sdk --language=ruby
   ```

2. If the project has no `Gemfile` yet, initialize one:

   ```bash
   bundle init
   ```

   Then add the OpenTelemetry gems and any application dependencies,
   and run `bundle install` to produce `Gemfile.lock`.

## Example

Minimal `Gemfile` for a service that ships OpenTelemetry via OTLP:

```ruby
source "https://rubygems.org"

gem "opentelemetry-api"
gem "opentelemetry-sdk"
gem "opentelemetry-exporter-otlp"
gem "opentelemetry-instrumentation-all"
```

Then:

```bash
bundle install
otel-checker check sdk --language=ruby
```

## Related

- `ruby.gemfile-lock.missing` — related failure once a `Gemfile` exists
  but hasn't been resolved yet.
- `ruby.gem.missing-required` — companion check that verifies the
  OpenTelemetry base gems are listed once the lock file is generated.
