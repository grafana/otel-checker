---
id: js.package-lock-json.unreadable
title: 'package-lock.json could not be read'
severity: error
---

## Why this matters

`otel-checker`'s supported-libraries check prefers `package-lock.json`
because it captures the exact resolved version of every dependency,
including transitive ones. That precision lets the checker tell you which
libraries in your project have OpenTelemetry instrumentation available and
which don't. When the lock file exists but can't be read (permission
denied, corrupted, mid-write), the check either falls back to `package.json`
(losing version accuracy and transitive coverage) or fails outright.

## How to fix

1. Confirm the lock file exists and is readable from the current directory:

   ```bash
   ls -l package-lock.json
   ```

2. If the file is intentionally missing (e.g. the project uses `yarn.lock`
   or `pnpm-lock.yaml`), that's not the same as this error — the checker
   silently falls back to `package.json` in that case. If you're seeing
   this error, `package-lock.json` is present but the read failed.

3. Regenerate the lock file if it looks corrupt:

   ```bash
   rm package-lock.json
   npm install
   ```

4. Fix the permission bits so the process running `otel-checker` can read
   it:

   ```bash
   chmod +r package-lock.json
   ```

## Example

Typical clean state:

```bash
$ ls -l package.json package-lock.json
-rw-r--r-- 1 you staff  1204 Jul 14 10:00 package.json
-rw-r--r-- 1 you staff 82134 Jul 14 10:00 package-lock.json
```

## Related

- `js.package-json.unreadable` — related failure on the manifest file.
- `js.supported-libs.fetch-failed` — a check that consumes the parsed
  dependency list.
