---
id: java.gradle.no-dependencies
title: 'Gradle reported no runtime dependencies'
severity: warning
---

## Why this matters

`otel-checker` runs `gradle dependencies --configuration=runtimeClasspath`
(via the wrapper if present) and parses the output to build the list of
libraries in your project — that's what it matches against the supported
Java instrumentation catalog. This warning fires when the command
completed but the output contained zero library entries, so the
supported-libraries check has nothing to report.

Common causes:

- **Empty project.** A brand-new Gradle project with no dependencies
  declared yet.
- **Wrong build file.** The checker read a Gradle file that isn't the
  application's — a settings-only root file in a multi-project build, or
  a script file the framework generated for something else.
- **Different configuration name.** Some setups (Android, custom source
  sets) put their runtime deps under a different configuration. The
  checker only looks at `runtimeClasspath`.
- **Build failed silently.** `gradle dependencies` produced output but
  the run itself errored earlier, leaving the parseable section empty.

## How to fix

1. Reproduce the exact command the checker runs so you can see what
   Gradle actually reports:

   ```bash
   ./gradlew --build-file=build.gradle dependencies --configuration=runtimeClasspath
   ```

   If the output is empty (or only shows a "No dependencies" note), your
   project genuinely has no runtime deps declared.

2. If your app lives inside a multi-project build, run `otel-checker`
   from the sub-project directory that contains the app's `build.gradle`
   (or `build.gradle.kts`), not the workspace root:

   ```bash
   cd services/api
   otel-checker check sdk --language=java
   ```

3. If your deps are on a custom configuration, add at least one entry
   under `implementation` / `runtimeOnly` so it ends up in
   `runtimeClasspath`, or extend that configuration from yours.

## Example

Minimal `build.gradle` that satisfies the check:

```groovy
plugins {
    id 'java'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web:3.5.0'
}
```

## Related

- [Gradle: viewing and debugging dependencies](https://docs.gradle.org/current/userguide/viewing_debugging_dependencies.html)
- `java.supported-libs.fetch-failed` — related error when the supported
  libraries list itself can't be fetched.
