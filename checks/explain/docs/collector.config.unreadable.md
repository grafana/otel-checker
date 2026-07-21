---
id: collector.config.unreadable
title: 'The collector config.yaml could not be read'
severity: error
---

## Why this matters

`otel-checker check collector` reads the Collector's config file from the
full path passed via `--collector-config-path`. When the flag is omitted,
it falls back to `config.yaml` (then `config.yml`) in the current working
directory. Every downstream collector check — endpoint format, receiver
protocols, per-pipeline exporter presence — reads from that file. If it
can't be opened, the checker has nothing to inspect and skips the whole
collector-side analysis.

Common causes:

- Running the checker from a directory that doesn't contain a
  `config.yaml` / `config.yml` and not passing `--collector-config-path`.
- Passing `--collector-config-path` with a wrong path (typo, relative
  path resolved from a different directory, or pointing at a directory
  instead of the file itself).
- Permissions: the config lives on disk but the invoking user doesn't
  have read access.

## How to fix

1. Confirm the file exists at the expected path:

   ```bash
   ls -l ./otel/my-collector.yaml
   ```

2. If your config lives elsewhere or uses a non-standard filename, pass
   the full file path via `--collector-config-path` (the flag accepts
   any filename):

   ```bash
   otel-checker check collector --collector-config-path=./otel/my-collector.yaml
   ```

3. If the file exists but the read failed with `permission denied`, fix
   the permission bits so the checker's user can read it:

   ```bash
   chmod +r ./otel/my-collector.yaml
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
otel-checker check collector --collector-config-path=./otel/config.yaml
```

## Related

- `collector.config.parse-error` — related failure when the file is
  readable but not valid YAML.
