---
id: ruby.bundler.not-found
title: 'Bundler is not installed'
severity: error
---

## Why this matters

`otel-checker` runs `bundle -v` to confirm Bundler is available. Bundler
is the standard tool for resolving `Gemfile` dependencies and generating
`Gemfile.lock` — without it, the OpenTelemetry gems your service depends
on can't be installed reproducibly, and every teammate or CI runner ends
up with a slightly different set of versions.

## How to fix

Install Bundler with the standard `gem` command:

```bash
gem install bundler
```

Verify:

```bash
bundle -v
# Bundler version 2.6.x
```

If your project pins a specific Bundler version in `Gemfile.lock` (the
`BUNDLED WITH` section), install that version too:

```bash
gem install bundler:2.6.2
```

In Docker, install Bundler in your image:

```dockerfile
FROM ruby:4-alpine
RUN gem install bundler
```

## Example

Complete gem-install flow from a fresh clone:

```bash
gem install bundler
bundle install
otel-checker check sdk --language=ruby
```

## Related

- [Bundler installation guide](https://bundler.io/)
- `ruby.runtime.not-found` — related failure when Ruby itself isn't on
  PATH.
- `ruby.gemfile-lock.missing` — related failure when Bundler is
  installed but hasn't produced a lock file yet.
