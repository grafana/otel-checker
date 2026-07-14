---
id: js.node-version.not-lts
title: 'Node.js is running an odd-major (non-LTS) release'
severity: error
---

## Why this matters

Node.js releases follow an even/odd major-version cadence: even-numbered
majors (22, 24, 26, …) transition to Active LTS about six months after
release and receive roughly 30 months of maintenance, while odd-numbered
majors (23, 25, …) are Current releases that ship in April, drop to
maintenance in October of the same year, and go end-of-life quickly after
that.

The OpenTelemetry JavaScript project supports the Active LTS lines. Running
your service on a Current release (Node 23, 25, …) may work today, but the
runtime is out of support before you finish your next quarter and there's
no long-term promise that OTel and its instrumentation packages will keep
their patch coverage on that line. Pinning to LTS keeps the production
runtime supported for years without another Node bump.

## How to fix

Switch to the nearest Active LTS release — at the time of writing, Node
22, 24, or 26.

- Local development with a version manager:

  ```bash
  # nvm
  nvm install 24
  nvm use 24

  # fnm
  fnm install 24
  fnm use 24
  ```

- Docker: change the base image, e.g. `FROM node:24-alpine`.
- CI: set the `node-version` input on `actions/setup-node` (or equivalent)
  to an even-numbered major.

## Example

`.nvmrc` pinning an LTS version for the repo:

```text
24
```

GitHub Actions:

```yaml
- uses: actions/setup-node@v7
  with:
    node-version-file: '.nvmrc'
```

## Related

- [Node.js release schedule](https://github.com/nodejs/release#release-schedule)
- [endoflife.date/nodejs](https://endoflife.date/nodejs) — current status
  of every Node major.
- `js.node-version.too-old` — related check that fires when Node is below
  the minimum supported LTS.
