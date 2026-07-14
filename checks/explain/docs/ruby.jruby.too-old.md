---
id: ruby.jruby.too-old
title: 'JRuby version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` requires JRuby 9.4 or newer — the current stable line.
Support for JRuby in OpenTelemetry Ruby is best-effort, so older
releases lack fixes and stdlib behavior the OpenTelemetry SDK depends
on. Installing the OpenTelemetry gems on JRuby 9.3 or older either
fails at `bundle install` or produces subtle runtime issues.

## How to fix

Upgrade JRuby to 9.4 or newer:

- **Version manager**:

  ```bash
  # rbenv (with the jruby plugin) or asdf
  asdf install ruby jruby-9.4.9.0
  asdf local ruby jruby-9.4.9.0
  ```

- **Docker**: use an official JRuby image, e.g.
  `FROM jruby:9.4-alpine`.
- **System install**: download from
  [jruby.org/download](https://www.jruby.org/download).

Verify:

```bash
jruby --version
# jruby 9.4.x.x ...
```

Also confirm JVM compatibility — modern JRuby needs Java 8 or newer.

## Example

`.ruby-version` pointing at a supported JRuby release:

```text
jruby-9.4.9.0
```

## Related

- [JRuby downloads and compatibility](https://www.jruby.org/download)
- `ruby.cruby.too-old` — related check for CRuby (the default
  implementation).
- `ruby.runtime.not-found` — related failure when no Ruby is installed.
