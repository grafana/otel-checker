---
id: js.node-version.too-old
title: 'Node.js version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` reads the local Node.js version via `node -v` and expects
at least Node 22. Current versions of `@opentelemetry/api`,
`@opentelemetry/sdk-node`, and `@opentelemetry/auto-instrumentations-node`
target Node 22+: they use modern syntax and runtime APIs that older
versions don't provide. Installing them on Node 20 or older either fails
outright at `npm install` or throws at startup with unhelpful
`SyntaxError` / `ReferenceError` messages that don't point at OpenTelemetry
as the culprit.

## How to fix

Upgrade Node to a supported LTS release (22 or newer). The mechanism
depends on your setup:

- Local development with a version manager:

  ```bash
  # nvm
  nvm install --lts
  nvm use --lts

  # fnm
  fnm install --lts
  fnm use lts-latest
  ```

- Docker: bump the base image, e.g. `FROM node:22-alpine`.
- CI: update the `node-version` input on `actions/setup-node`, GitLab
  `image:`, etc.
- System package: install a current release from
  [nodejs.org](https://nodejs.org/) or your distro's up-to-date repo.

After upgrading, verify:

```bash
node -v
# v22.11.0
```

## Example

`.nvmrc` pinning a supported version for the repo:

```text
22
```

GitHub Actions:

```yaml
- uses: actions/setup-node@v7
  with:
    node-version-file: '.nvmrc'
```

## Related

- [Node.js release schedule](https://github.com/nodejs/release#release-schedule)
- `js.node-version.unknown` — related failure when `node` isn't callable
  at all.
