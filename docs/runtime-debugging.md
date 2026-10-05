# Still missing telemetry?

A static instrumentation finding is a starting point, not a runtime diagnosis.
Passing dependency or configuration checks do not prove that the application
executed, the agent attached, telemetry was exported, or a destination stored it.
Even a successful destination connectivity check does not prove app ingestion.

Use the existing gcx skills for runtime investigation. Follow the gcx-owned
[diagnostic guide][diagnostics] for installation, skill setup, connection,
evidence limits, repair approval, and cleanup; this page covers checker evidence.
The linked gcx guide is pinned to v1.4.0. Check `gcx version` and the relevant
command `--help` for your installed release, and record your actual checker and
gcx versions when reproducing a problem.

## Preserve the checker result

Keep the original command, relevant output, and finding's `explain_id` locally.
Use [Explanations](../README.md#explanations) to understand the finding before
investigating runtime behavior. A missing explain ID is not a blocker: preserve
the actual result instead of inventing an ID.

When sharing evidence with your agent, include:

- The selected finding and its explanation, or that the relevant checks passed.
- Language/runtime and the configuration being checked.
- The missing signal, application/service identity, and expected operation.
- The actual deployment/destination and time window, if known.
- The application/Collector scope you own and what access is unavailable.

Review output for credentials and sensitive paths before sharing it. Keep full
result files local; don't embed findings, service names, paths, or secrets in
guide URLs. A checker result describes the inspected configuration, which may
not be the configuration used by the running process.

## Establish runtime access

You do not need gcx installed to start reading the shared guide: it begins with
installation and skill setup. Backend queries require a reachable Grafana
destination and appropriate permissions.

If no backend exists, [start a local docker-lgtm instance][lgtm] and connect an
approved example or application. This is a new local proof, not proof that an
existing production exporter or inaccessible destination works. Do not change
an application's exporter without permission. If access is missing, report that
boundary as unobserved rather than claiming data loss.

## Ask for an evidence-based investigation

Start your agent in the relevant checkout after following the shared setup:

> otel-checker reported `<finding/explain ID>`, and I'm missing `<signal>` from
> `<service>`. Use the existing gcx skills with `<config path and context>`.
> Treat the checker result as a starting point, not proof of the runtime cause.
> Inspect only `<application/deployment scope>`. State what you could not observe;
> do not infer data loss from missing access. Diagnose first and ask before
> making changes, restarting resources, or enabling payload logging.

If checks passed, replace the first sentence with the checks that passed and
the remaining symptom. No finding is required to investigate missing telemetry.

After an approved repair, exercise a fresh operation and verify the original
symptom using the shared guide. Re-running static checks alone cannot establish
runtime recovery. Record remaining unknowns and restore temporary diagnostics.

[diagnostics]: https://github.com/grafana/gcx/blob/v1.4.0/docs/guides/diagnose-missing-telemetry.md
[lgtm]: https://github.com/grafana/docker-otel-lgtm/blob/main/docs/gcx-integration.md
