# Running `otel-checker` against a Kubernetes-hosted application

This runbook is for engineers helping a customer whose OpenTelemetry-instrumented
application runs in a Kubernetes pod. It explains how to validate that
instrumentation with the current `otel-checker` binary.

## Scope

`otel-checker` is a local process that scans the current working directory
for source files (`.csproj`, `package.json`, `go.mod`, `pom.xml`,
`requirements.txt`, `pyproject.toml`, …) and reads environment variables
from its own process. Nothing in the binary talks to the Kubernetes API.

To cover an application deployed to Kubernetes, split the work in two:

| Part | Where it runs | What it validates |
|---|---|---|
| **A. Live-pod checks** | Inside the target container's namespace | `env` variables (`OTEL_*`), Grafana Cloud endpoint / auth / connectivity |
| **B. Source-repo checks** | Developer machine or CI, against the app's Git repo | SDK setup — installed OTel packages, versions, target framework, supported instrumentations |

Neither part alone is enough. Part A confirms the *running deployment* is
wired up; Part B confirms the *built artifact* has the right OTel packages.

Neither part inspects the Collector / Alloy / Beyla pods themselves in
this runbook — those are separate deployments; run `otel-checker check
collector` (or `alloy` / `beyla`) against their config files the same
way as Part B.

## Prerequisites

- `kubectl` configured for the customer's cluster with permission to `get`
  and `exec` pods in the target namespace.
- For Part A via `kubectl debug`: permission on the
  `pods/ephemeralcontainers` subresource, and a cluster on Kubernetes
  1.25 or later (Ephemeral Containers GA).
- A Linux `otel-checker` binary matching the pod's CPU architecture —
  download from
  [github.com/grafana/otel-checker/releases](https://github.com/grafana/otel-checker/releases).
  The statically-linked Linux archives work on both glibc and musl
  (Alpine) targets.
- For Part B: `git` and the language's toolchain (`dotnet`, `mvn`,
  `npm`, `python`, `go`, …) installed wherever the checker runs.

## Part A — Check the running pod

The env and Grafana Cloud components need three things: the target
container's env vars, the ability to reach the OTLP endpoint from the
pod's network, and the `otel-checker` binary itself. Two ways to deliver
the binary, pick whichever fits the customer's image:

### Option 1: `kubectl debug` with a custom image (works with distroless / read-only apps)

1. Build a minimal image that ships the checker. No official image is
   published — the Dockerfile is one layer:

   ```dockerfile
   FROM alpine:3
   RUN apk add --no-cache ca-certificates
   COPY otel-checker /otel-checker
   ENTRYPOINT ["/bin/sh"]
   ```

   Build and push to a registry the cluster can pull from:

   ```bash
   docker build --platform=linux/amd64 -t <your-registry>/otel-checker:latest .
   docker push <your-registry>/otel-checker:latest
   ```

2. Attach as an ephemeral debug container sharing the target's process
   namespace:

   ```bash
   POD=<pod-name>
   NS=<namespace>
   APP_CONTAINER=<container-name>

   kubectl debug -it $POD -n $NS \
     --image=<your-registry>/otel-checker:latest \
     --target=$APP_CONTAINER \
     --profile=general \
     -- sh
   ```

3. From the debug shell, adopt the target container's env and working
   directory via `/proc`:

   ```sh
   APP_PID=$(pgrep -n -f <app-binary-or-keyword>)   # e.g. dotnet, java, node, python
   cd /proc/$APP_PID/cwd                            # target's working directory
   export $(cat /proc/$APP_PID/environ | tr '\0' '\n' | xargs -d '\n')

   /otel-checker check grafana-cloud --language=<lang> --format=json > /tmp/results.json
   cat /tmp/results.json
   ```

   Replace `<lang>` with one of `dotnet`, `go`, `java`, `js`, `python`,
   `ruby`, `php`.

4. Copy the results file out for review or web-UI replay:

   ```bash
   kubectl cp $NS/$POD:/tmp/results.json ./results.json -c debugger
   otel-checker serve --data=./results.json          # opens the local web UI
   otel-checker explain                              # markdown docs for every flagged ID
   ```

### Option 2: `kubectl cp` + `kubectl exec` (works only if the app image has `sh` and is writable)

1. Copy the Linux binary into the running container:

   ```bash
   kubectl cp ./otel-checker $NS/$POD:/tmp/otel-checker -c $APP_CONTAINER
   kubectl exec $NS/$POD -c $APP_CONTAINER -- chmod +x /tmp/otel-checker
   ```

2. Run it — the exec'd process automatically inherits the container's
   env, so no `/proc` juggling:

   ```bash
   kubectl exec $NS/$POD -c $APP_CONTAINER -- \
     /tmp/otel-checker check grafana-cloud --language=<lang> --format=json > results.json
   ```

3. Same follow-up as Option 1 — `serve` and `explain` locally.

### What Part A catches

Running `check grafana-cloud` and `check env` (or the default `check`
without arguments, filtered to those components) validates:

- `OTEL_SERVICE_NAME`, `OTEL_RESOURCE_ATTRIBUTES` (including
  `service.instance.id`, `service.namespace`, `deployment.environment.name`)
- `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_PROTOCOL`,
  `OTEL_EXPORTER_OTLP_HEADERS` (presence, format, auth header shape)
- `OTEL_EXPORTER_OTLP_HEADERS` decoded and verified against the
  Grafana Cloud stack the token belongs to
- DNS + TCP reachability from the pod to the endpoint
- Language-specific runtime env (`NODE_OPTIONS`, `PYTHONPATH`,
  `CORECLR_*`, `OTEL_DOTNET_AUTO_HOME`, `JAVA_TOOL_OPTIONS`, …)

## Part B — Check the source repository

The SDK component needs the project files and the language toolchain,
which typically aren't present in a production container image (multi-stage
Docker builds strip them). Run it against the customer's Git checkout
instead.

### Locally

```bash
git clone <customer-repo>
cd <customer-repo>/<service-directory>          # cd into the directory that holds the project file
otel-checker check sdk --language=<lang>
```

For manual instrumentation, add `--manual-instrumentation` (and
`--instrumentation-file=<path>` for JS).

### In CI

Add a job to the customer's pipeline (GitHub Actions example):

```yaml
- name: Install otel-checker
  env:
    OTEL_CHECKER_VERSION: '0.3.3'   # pin to a release from https://github.com/grafana/otel-checker/releases
  run: |
    curl -sSfL "https://github.com/grafana/otel-checker/releases/download/v${OTEL_CHECKER_VERSION}/otel-checker_${OTEL_CHECKER_VERSION}_linux_amd64.tar.gz" \
      | sudo tar -xz -C /usr/local/bin otel-checker

- name: Validate OTel instrumentation
  working-directory: services/<service>
  run: otel-checker check sdk --language=<lang> --format=json > otel-checker.json

- uses: actions/upload-artifact@v4
  with:
    name: otel-checker-results
    path: services/<service>/otel-checker.json
```

### What Part B catches

- Required OTel SDK and instrumentation packages present
- Package versions above the supported minimum
- Language runtime / target framework at or above the OTel-supported minimum
- Recognised distro / SDK (`Microsoft.NET.Sdk` vs `Microsoft.NET.Sdk.Web`, etc.)
- Supported instrumentation libraries detected for each dependency
- Manual-instrumentation entry point wired up correctly (JS)

## Coverage matrix — which mode covers what

| Check family | Part A (live pod) | Part B (source repo) |
|---|---|---|
| `env` (OTEL_* variables) | ✅ | ⚠️ only if env is also set locally |
| `grafana-cloud` (endpoint / auth / reachability) | ✅ | ⚠️ reachability from CI, not from the pod |
| `sdk` (packages, versions, target framework) | ❌ files usually absent | ✅ |
| `collector`, `alloy`, `beyla` config | run separately against the config file | run separately against the config file |

## Troubleshooting

**`no .csproj files found in directory`** (or the equivalent for other
languages) — the current working directory doesn't contain a project
file. In Part A this usually means the app container ships only the
compiled artifact; that's expected — the SDK check is a Part B
responsibility. In Part B, `cd` into the directory that holds the
project file before running.

**`could not check .NET version: exec: "dotnet": executable file not found`**
— the container has the .NET runtime but not the SDK. Same story: skip
the SDK check in Part A, run it in Part B against the repo.

**`kubectl debug` returns "ephemeral containers are disabled"** — the
cluster is older than 1.25 or has the feature gate off. Fall back to
Option 2 (`kubectl cp` + `exec`).

**Env vars set at the deployment level but missing in the checker's
output** — the checker only sees the environment of the process that
runs it. If you're in a debug shell that doesn't inherit them, use the
`/proc/$APP_PID/environ` export from Option 1 or exec inside the app
container directly (Option 2).

**Grafana Cloud connectivity check fails from the pod but works from
your laptop** — the pod has restricted egress. Check `NetworkPolicy`,
egress firewall rules, and any proxy the pod should be routing through
(`HTTPS_PROXY` / `HTTP_PROXY` env vars).

**`--format=json` output is not valid JSON** — the checker writes
human-readable output to stderr by default; make sure you're redirecting
`stdout` only (`> results.json`, not `&> results.json`) and passing
`--format=json`.

## Interpreting findings

Every error and warning carries a stable ID (e.g.
`env.otel-service-name.unset`,
`grafana-cloud.headers.missing-auth`). Look it up with:

```bash
otel-checker explain <id>               # one finding
otel-checker explain                    # every ID in ./results.json
otel-checker explain list               # every ID the binary knows about
otel-checker serve --data=./results.json   # local web UI with the same docs
```

The explain docs describe the failure mode, the impact, and concrete
remediation — usually the fastest way to hand the customer a
specific fix.
