---
id: ruby.cruby.too-old
title: 'CRuby version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` requires CRuby (MRI Ruby, the default C implementation)
version 3.3 or newer. Everything below 3.3 has reached end-of-life
([endoflife.date/ruby](https://endoflife.date/ruby)), so it no longer
receives security fixes and the OpenTelemetry Ruby gems no longer
guarantee compatibility with it. Installing the OpenTelemetry gems on
Ruby 3.2 or older either fails at `bundle install` or produces runtime
errors that don't obviously point at OpenTelemetry.

## How to fix

Upgrade to a supported CRuby release (3.3, 3.4, or newer):

- **Version manager** (recommended for local development):

  ```bash
  # rbenv
  rbenv install 3.3.5
  rbenv local 3.3.5

  # asdf
  asdf install ruby 3.3.5
  asdf local ruby 3.3.5
  ```

- **Docker**: bump the base image, e.g. `FROM ruby:3.3-alpine`.
- **CI**: pin the Ruby version explicitly, e.g. with `actions/setup-ruby`.

Pin the version in `.ruby-version` so the repo advertises the required
release:

```text
3.3.5
```

Verify:

```bash
ruby -v
# ruby 3.3.5 ...
```

## Example

`.ruby-version` + CI setup on GitHub Actions:

```yaml
- uses: ruby/setup-ruby@v1
  with:
    ruby-version: '.ruby-version'
    bundler-cache: true
```

## Related

- [Ruby release schedule](https://www.ruby-lang.org/en/downloads/branches/)
- [endoflife.date/ruby](https://endoflife.date/ruby) — current status of
  every Ruby major.
- `ruby.jruby.too-old` — related check for JRuby.
- `ruby.runtime.not-found` — related failure when Ruby isn't installed
  at all.
