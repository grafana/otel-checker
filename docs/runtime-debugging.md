# Still missing telemetry?

A static instrumentation finding can point you toward a problem, but it can't
tell you what happened at runtime. Passing dependency or configuration checks
doesn't show that the application ran, the agent attached, telemetry was
exported, or a destination stored it. Even a successful connectivity check
doesn't confirm that the application sent data.

For runtime investigation, use the existing gcx skills and its
[diagnostic guide][diagnostics]. That guide covers installation, skill setup,
connecting to a destination, evidence limits, repair approval, and cleanup;
this page focuses on what the checker can tell you. The gcx guide is pinned to
v1.4.0, so check `gcx version` and the relevant command's `--help` for your
installed release. When reproducing a problem, note the versions of both the
checker and gcx that you actually used.

## Preserve the checker result

Keep the command you ran, its relevant output, and the finding's `explain_id`
locally. Before looking into runtime behavior, use
[Explanations](../README.md#explanations) to understand the finding. Some
results don't have an explain ID; that's fine—keep the actual result rather
than trying to invent one.

If you ask an agent to investigate, give it enough context to distinguish what
you know from what still needs checking. Include:

- The selected finding and its explanation, or that the relevant checks passed.
- Language/runtime and the configuration being checked.
- The missing signal, application/service identity, and expected operation.
- The actual deployment/destination and time window, if known.
- The application/Collector scope you own and what access is unavailable.

Before sharing, check the output for credentials and sensitive paths. Keep full
result files local, and don't put findings, service names, paths, or secrets in
guide URLs. Remember that the checker reports on the configuration it
inspected; the running process may be using something else.

## Establish runtime access

You can start with the shared guide even if gcx isn't installed; it explains
installation and skill setup. To query a backend, you'll need a reachable
Grafana destination and the right permissions.

If you don't have access to a backend, you can [start a local docker-lgtm
instance][lgtm] and connect an approved example or application. That lets you
make a new local proof; it won't show that an existing production exporter or
an inaccessible destination works. Don't change an application's exporter
without permission. If you can't access a system, say that it wasn't observed
rather than concluding that data was lost.

## Ask for an evidence-based investigation

Once you've followed the shared setup, you can start an investigation in the
relevant checkout with a request like this:

> otel-checker reported `<finding/explain ID>`, and I'm missing `<signal>` from
> `<service>`. Use the existing gcx skills with `<config path and context>`.
> Treat the checker result as a starting point, not proof of the runtime cause.
> Inspect only `<application/deployment scope>`. State what you could not observe;
> do not infer data loss from missing access. Diagnose first and ask before
> making changes, restarting resources, or enabling payload logging.

If the checks passed, replace the first sentence with the checks you ran and
describe the remaining symptom. You don't need a finding to investigate missing
telemetry.

After a repair is approved, exercise a fresh operation and use the shared guide
to check whether the original symptom is gone. Re-running static checks alone
can't confirm runtime recovery. Note anything that remains unknown, and remove
any temporary diagnostics.

[diagnostics]: https://github.com/grafana/gcx/blob/v1.4.0/docs/guides/diagnose-missing-telemetry.md
[lgtm]: https://github.com/grafana/docker-otel-lgtm/blob/main/docs/gcx-integration.md
