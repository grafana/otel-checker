---
id: ruby.runtime.not-found
title: 'No supported Ruby runtime found'
severity: error
---

## Why this matters

The Ruby checker probes three runtimes — CRuby (`ruby -v`), JRuby
(`jruby --version`), and TruffleRuby — and requires at least one to
report a supported version. None of them responded here, which means no
Ruby interpreter is on the current `PATH` (or all three returned
versions below the minimums). Without a runtime the checker can't
validate that your Gemfile is compatible with the OpenTelemetry Ruby
gems, and your service can't run at all.

Supported minimums:

- CRuby ≥ 3.3
- JRuby ≥ 9.4
- TruffleRuby ≥ 22.1 (best-effort)

## How to fix

Install a Ruby runtime that meets one of the minimums:

- **CRuby via a version manager** (recommended for local development):

  ```bash
  # rbenv
  rbenv install 3.3.5
  rbenv local 3.3.5

  # asdf
  asdf install ruby 3.3.5
  asdf local ruby 3.3.5
  ```

- **Docker**: use an official image, e.g. `FROM ruby:3.3-alpine`.
- **System package**: install from your distro's package manager or from
  [ruby-lang.org](https://www.ruby-lang.org/en/downloads/).

Verify the checker's environment sees Ruby:

```bash
which ruby
ruby -v
```

## Example

`.ruby-version` file pinning a supported release for the repo:

```text
3.3.5
```

## Related

- `ruby.cruby.too-old` / `ruby.jruby.too-old` — related checks when a
  runtime is found but its version is below the minimum.
- [Ruby downloads](https://www.ruby-lang.org/en/downloads/)
