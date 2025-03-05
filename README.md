# OTel Me If It's Right

Checker if the implementation of OpenTelemetry instrumentation is correct.

## Usage

Requirement: Golang

### Installation
1. Install the `otel-checker` binary
```
go install github.com/grafana/otel-checker@latest
```
2. You can confirm it was installed with:
```
❯ ls $GOPATH/bin
otel-checker
```

### Flags

The available flags:
```
❯ otel-checker -h
Usage of otel-checker:
  -manual-instrumentation
    	Provide if your application is using manual instrumentation (auto instrumentation as default)
  -collector-config-path string
    	Path to collector's config.yaml file. Required if using Collector and the config file is not in the same location as the otel-checker is being executed from. E.g. "-collector-config-path=src/inst/"
  -components string
    	Instrumentation components to test, separated by ',' (required). Possible values: sdk, collector, beyla, alloy
  -instrumentation-file string
    	Name (including path) to instrumentation file. Required if using manual-instrumentation. E.g."-instrumentation-file=src/inst/instrumentation.js"
  -language string
    	Language used for instrumentation (required). Possible values: dotnet, go, java, js, python
  -package-json-path string
    	Path to package.json file. Required if instrumentation is in JavaScript and the file is not in the same location as the otel-checker is being executed from. E.g. "-package-json-path=src/inst/"
```

### Checks

#### Grafana Cloud
- Endpoints
- Authentication
- Service name
- Exporter protocol

#### SDK

##### JavaScript
- Node version
- Required dependencies on package.json
- Required environment variables
- Resource detectors
- Dependencies compatible with Grafana Cloud
- Usage of Console Exporter

#### Python
TBD

#### .NET
TBD

#### Java

- Java version
- Prints which libraries (as discovered from a locally running maven or gradle) are supported:
  - With `-manual-instrumentation`, the libraries for manual instrumentation are printed.
  - Without `-manual-instrumentation`, it will print the libraries supported by the [Java Agent](https://github.com/open-telemetry/opentelemetry-java-instrumentation/).
  - A maven or gradle wrapper will be used if found in the current directory or a parent directory.

#### Go
TBD

#### Ruby

- Ruby version
- Bundler installation
- Required dependencies in `Gemfile.lock`:

#### Collector
- Config receivers and exporters

#### Beyla
- Environment variables

#### Alloy
TBD

### Examples

Application with auto-instrumentation
![auto instrumentation exemple](./assets/auto.png)

Application with custom instrumentation using SDKs and Collector
![sdk and collector example](./assets/sdk.png)

## Development

Requirement: Golang

### Running locally

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

### Create binary and run from different directory

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
