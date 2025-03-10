# Contributing to otel-checker

Thank you for your interest in contributing to otel-checker! This document provides guidelines and instructions for contributing to this project.

## Development Environment Setup

1. Ensure you have Go installed (1.24 or higher)
2. Clone the repository:
   ```
   git clone https://github.com/yourusername/otel-checker.git
   cd otel-checker
   ```
3. Install dependencies:
   ```
   go mod download
   ```
                   

## Running locally

1. Find your Go path:
```
❯ go env GOPATH
/Users/maryliag/go
```
2. Clone this repo in the go path folder, so you will have:
```
/Users/maryliag/go/src/otel-checker
```
3. Run
```
go run main.go
```

## Create binary and run from different directory

1. Build binary
```
go build
```
2. Install
```
go install
```
3. You can confirm it was installed with:
```
❯ ls $GOPATH/bin
otel-checker
```
4. Use from any other directory
```
otel-checker \
	-language=js \
	-components=sdk
```

Or start directly from the source code:
```
go run otel-checker \
	-language=js \
	-components=sdk
```

## Using the Makefile

We provide a Makefile with several useful commands to simplify common development tasks. The Makefile helps ensure consistent code quality and streamlines the development workflow.

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Builds the application using `go install` |
| `make test` | Runs all tests in the project |
| `make clean` | Removes build artifacts and cleans the Go cache |
| `make fmt` | Formats all Go code using `gofmt` |
| `make lint` | Lints the code using `golangci-lint` |
| `make deps` | Updates dependencies using `go mod tidy` |
| `make help` | Displays help information about available commands |

### Dependency: golangci-lint

For linting, we use golangci-lint. If not already installed, you can install it with:

```
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Contribution Workflow

1. Create a fork of the repository
2. Create a new branch for your feature or bug fix
3. Make your changes
4. Run `make fmt` to format your code
5. Run `make lint` to ensure code quality
6. Run `make test` to make sure all tests pass
7. Commit your changes with a descriptive message
8. Submit a pull request to the main repository

## Before Submitting Pull Requests

Please ensure:

1. Your code follows the project's style and conventions
2. All tests pass (`make test`)
3. Code is properly formatted (`make fmt`)
4. Linting passes without issues (`make lint`)

## Code Review Process

Once you submit a pull request:

1. Maintainers will review your code
2. They may request changes or improvements
3. Once approved, your PR will be merged into the main branch
