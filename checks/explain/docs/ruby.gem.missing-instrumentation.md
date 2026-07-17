---
id: ruby.gem.missing-instrumentation
title: 'No OpenTelemetry instrumentation gem is installed'
severity: error
---

## Why this matters

Auto-instrumentation for a Ruby application is delivered as a set of
per-library gems — `opentelemetry-instrumentation-rack`,
`opentelemetry-instrumentation-rails`, `opentelemetry-instrumentation-pg`,
and so on — or as the meta-gem `opentelemetry-instrumentation-all` that
bundles every stable instrumentation. `otel-checker`'s auto-instrumentation
mode expects at least one of these to be listed in `Gemfile.lock`.

Without any instrumentation gem, the SDK from `opentelemetry-sdk`
initializes but nothing hooks into the frameworks and libraries in your
app, so you get no spans, metrics, or logs from HTTP handlers, ORMs,
background jobs, or external calls — even though the tracer provider
itself is set up correctly.

## How to fix

The easiest option is to pull in the meta-gem, which enables every
stable instrumentation:

```ruby
# Gemfile
gem "opentelemetry-instrumentation-all"
```

Then run `bundle install` and enable them in your SDK bootstrap:

```ruby
require "opentelemetry/sdk"
require "opentelemetry/instrumentation/all"

OpenTelemetry::SDK.configure do |c|
  c.service_name = "checkout"
  c.use_all
end
```

If you only need instrumentation for a specific framework (a smaller
image footprint, or an unusual stack), install the per-library gems
individually instead — e.g. for Rails and Sidekiq:

```ruby
gem "opentelemetry-instrumentation-rails"
gem "opentelemetry-instrumentation-sidekiq"
```

Enable them in the bootstrap with `c.use "OpenTelemetry::Instrumentation::Rails"`.

## Example

`Gemfile` fragment for a Rack-based service:

```ruby
gem "opentelemetry-api"
gem "opentelemetry-sdk"
gem "opentelemetry-exporter-otlp"
gem "opentelemetry-instrumentation-rack"
gem "opentelemetry-instrumentation-net_http"
```

## Related

- [OpenTelemetry Ruby instrumentation registry](https://github.com/open-telemetry/opentelemetry-ruby-contrib)
- `ruby.gem.missing-required` — related error for the three core gems
  every setup needs.
