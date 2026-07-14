---
id: beyla.open-port.unset
title: 'BEYLA_OPEN_PORT is not set'
severity: error
---

## Why this matters

`BEYLA_OPEN_PORT` tells Beyla which TCP port your application listens on.
Beyla uses eBPF to attach to the process bound to that port and observe
its traffic — that's how it produces spans and metrics without any code
change to your app. Without this variable, Beyla has no target: it
starts, but it never attaches to anything and no telemetry is produced.

Multiple ports can be specified as a comma-separated list, and port
ranges are supported (e.g. `8080-8089`).

## How to fix

Set `BEYLA_OPEN_PORT` to the port(s) your service actually binds. In
Docker:

```dockerfile
ENV BEYLA_OPEN_PORT=8080
```

In Kubernetes (typically as a sidecar container):

```yaml
env:
  - name: BEYLA_OPEN_PORT
    value: "8080"
```

Locally:

```bash
export BEYLA_OPEN_PORT=8080
sudo -E beyla
```

Beyla needs `CAP_SYS_ADMIN` (or root) to load its eBPF probes — that's
why the local invocation uses `sudo -E` to preserve the environment.

## Example

Comma-separated list, and a range:

```bash
export BEYLA_OPEN_PORT=8080,8081        # explicit list
export BEYLA_OPEN_PORT=8080-8089        # port range
```

## Related

- [Beyla configuration reference](https://grafana.com/docs/beyla/latest/configure/options/)
- `beyla.grafana-cloud-submit.unset` — companion check for the submit-type
  variable Beyla needs to know what to ship.
