# Contributing to otel-checker

Thank you for your interest in contributing to otel-checker! This document
provides guidelines and instructions for contributing to this project.

## Development Environment Setup

1. Ensure you have Go installed (1.24 or higher)
2. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/otel-checker.git
   cd otel-checker
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

## Running locally

1. Find your Go path:

```bash
❯ go env GOPATH
/Users/maryliag/go
```

1. Clone this repo in the go path folder, so you will have:

```text
/Users/maryliag/go/src/otel-checker
```

1. Run

```bash
go run ./cmd/otel-checker
```

## Create binary and run from different directory

1. Build binary

```bash
go build ./cmd/otel-checker
```

1. Install

```bash
go install ./cmd/otel-checker
```

1. You can confirm it was installed with:

```bash
❯ ls $GOPATH/bin
otel-checker
```

1. Use from any other directory

```bash
otel-checker check sdk --language=js
```

Or start directly from the source code:

```bash
go run ./cmd/otel-checker check sdk --language=js
```

## Using mise

We provide a `mise.toml` file with several useful commands to simplify common
development tasks. [mise](https://mise.jdx.dev/) helps ensure consistent code
quality and streamlines the development workflow.

### Available mise Commands

| Command             | Description                                     |
| ------------------- | ----------------------------------------------- |
| `mise run build`    | Builds the application using `go install ./cmd/otel-checker` |
| `mise run test`     | Runs all tests in the project                   |
| `mise run clean`    | Removes build artifacts and cleans the Go cache |
| `mise run lint:fix` | Auto-fix lint and formatting issues             |
| `mise run lint`     | Run all lints                                   |
| `mise run check`    | Run all checks (test and lint)                  |
| `mise run deps`     | Updates dependencies using `go mod tidy`        |

## Contribution Workflow

1. Create a fork of the repository
2. Create a new branch for your feature or bug fix
3. Make your changes
4. Run `mise run lint:fix` to format your code
5. Run `mise run lint` to ensure code quality
6. Run `mise run test` to make sure all tests pass
7. Commit your changes with a descriptive message
8. Submit a pull request to the main repository

## Before Submitting Pull Requests

Please ensure:

1. Your code follows the project's style and conventions
2. All tests pass (`mise run test`)
3. Code is properly formatted (`mise run lint:fix`)
4. Linting passes without issues (`mise run lint`)

## Code Review Process

Once you submit a pull request:

1. Maintainers will review your code
2. They may request changes or improvements
3. Once approved, your PR will be merged into the main branch

## Submit a pull request

Effective 2026-06-22, all Grafana Labs repositories [require signed commits][signed-commits].
To learn more about Git commit verification, refer to [About commit signature verification][signing-commits]
and [Checking your commit signature verification status][verifying-commits].

> [!NOTE]
> Pull requests containing any unsigned commits cannot be merged until all commits are signed.

[signed-commits]: https://docs.github.com/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches#require-signed-commits
[signing-commits]: https://docs.github.com/authentication/managing-commit-signature-verification/about-commit-signature-verification
[verifying-commits]: https://docs.github.com/authentication/troubleshooting-commit-signature-verification/checking-your-commit-and-tag-signature-verification-status
