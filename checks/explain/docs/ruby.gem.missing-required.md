---
id: ruby.gem.missing-required
title: 'A required OpenTelemetry gem is missing from Gemfile.lock'
severity: error
---

## Why this matters

The checker looks in `Gemfile.lock` for three gems that any
OpenTelemetry Ruby setup needs regardless of whether you use auto- or
manual-instrumentation:

- `opentelemetry-api` — the tracer / meter interfaces your code (and
  instrumentation gems) call into.
- `opentelemetry-sdk` — the SDK that actually assembles spans, metrics,
  and logs.
- `opentelemetry-exporter-otlp` — the exporter that ships those signals
  to Grafana Cloud via OTLP.

If any of them is missing, either the SDK never registers a tracer
provider (so instrumentation gems produce no-op spans), or spans are
produced but there's no exporter to ship them anywhere.

## How to fix

Add the missing gem to your `Gemfile` and run `bundle install`. The
finding names the specific gem in its message.

```ruby
# Gemfile
source "https://rubygems.org"

gem "opentelemetry-api"
gem "opentelemetry-sdk"
gem "opentelemetry-exporter-otlp"
```

Then:

```bash
bundle install
otel-checker check sdk --language=ruby
```

## Example

Minimal instrumentation bootstrap (e.g. in `config/initializers/otel.rb`
for a Rails app):

```ruby
require "opentelemetry/sdk"
require "opentelemetry/exporter/otlp"

OpenTelemetry::SDK.configure do |c|
  c.service_name = "checkout"
  c.use_all       # requires opentelemetry-instrumentation-all
end
```

## Related

- `ruby.gem.missing-instrumentation` — related error when instrumentation
  gems (e.g. `opentelemetry-instrumentation-rails`) are missing.
- [OpenTelemetry Ruby: getting started](https://opentelemetry.io/docs/languages/ruby/getting-started/)
