---
id: collector.config.unreadable
title: 'The collector config.yaml could not be read'
severity: error
---

## Why this matters

`otel-checker check collector` looks for `config.yaml` inside the directory
passed via `--collector-config-path` (or the current directory when the
flag is omitted). Every downstream collector check — endpoint format,
receiver protocols, per-pipeline exporter presence — reads from that
file. If it can't be opened, the checker has nothing to inspect and skips
the whole collector-side analysis.

Common causes:

- Running the checker from a directory that doesn't contain the
  Collector's config.
- Passing `--collector-config-path` to a folder that doesn't hold a file
  literally named `config.yaml` (a different filename, or a symlink to
  a file the current user can't read).
- Permissions: the config lives on disk but the invoking user doesn't
  have read access.

## How to fix

1. Confirm the file exists at the expected path:

   ```bash
   ls -l ./config.yaml
   ```

2. If your config lives elsewhere, point the checker at its directory:

   ```bash
   otel-checker check collector --collector-config-path=./otel/
   ```

   The flag takes a directory; the file inside must be named
   `config.yaml`. Rename or symlink if your setup uses a different
   filename.

3. If the file exists but the read failed with `permission denied`, fix
   the permission bits so the checker's user can read it:

   ```bash
   chmod +r ./config.yaml
   ```

## Example

Typical layout:

```text
my-service/
├── otel/
│   └── config.yaml       ← collector config
└── ...
```

Invocation from `my-service/`:

```bash
otel-checker check collector --collector-config-path=./otel/
```

## Related

- `collector.config.parse-error` — related failure when the file is
  readable but not valid YAML.
