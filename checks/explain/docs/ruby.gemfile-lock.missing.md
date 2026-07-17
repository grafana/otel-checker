---
id: ruby.gemfile-lock.missing
title: 'Gemfile.lock is missing'
severity: error
---

## Why this matters

`Gemfile.lock` records the exact resolved version of every gem your
project uses — including transitive dependencies. `otel-checker` reads
it to see which OpenTelemetry gems are actually installed and at what
versions. Without a lock file, Bundler hasn't resolved the graph yet,
and the checker can't tell whether `opentelemetry-api`,
`opentelemetry-sdk`, and the exporter gem are present.

## How to fix

Run `bundle install` to resolve dependencies and generate the lock file:

```bash
bundle install
```

Commit `Gemfile.lock` to the repo — it's not a build artifact, it's a
guarantee that every environment resolves the same versions. Running
`bundle install` on a fresh checkout should be idempotent when the lock
file is present.

## Example

Clean install flow:

```bash
bundle install
ls Gemfile.lock         # exists now
otel-checker check sdk --language=ruby
```

## Related

- `ruby.gemfile.missing` — related failure when `Gemfile` itself
  doesn't exist yet.
- `ruby.bundler.not-found` — related failure when `bundle` isn't
  installed.
- [Bundler: understanding Gemfile.lock](https://bundler.io/guides/gemfile.html)
