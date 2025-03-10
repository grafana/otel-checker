# Go Supported Libraries

This directory contains the configuration and checker for Go libraries that can be instrumented with OpenTelemetry.

## Supported Libraries File

The `supported-libraries.yaml` file contains information about which Go libraries are supported by OpenTelemetry instrumentation, along with their version ranges and source paths.

### File Format

```yaml
import_path:
  instrumentations:
    - name: import_path
      srcPath: path/to/instrumentation
      link: go.opentelemetry.io/contrib/path/to/instrumentation
      target_versions:
        library:
          - [min_version,max_version)
```

Example:
```yaml
go.mongodb.org/mongo-driver:
  instrumentations:
  - name: go.mongodb.org/mongo-driver
    srcPath: instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo
    link: go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo
    target_versions:
      library:
      - '[1.17.3,2.0.0)'
```

## Generating the File

The `supported-libraries.yaml` file is generated from the OpenTelemetry Go Contrib repository using the `generate_go_supported_libraries.py` script.

### Prerequisites

- Python 3.x
- PyYAML package (`pip install pyyaml`)
- A local clone of the [OpenTelemetry Go Contrib repository](https://github.com/open-telemetry/opentelemetry-go-contrib)

### Usage

```bash
# Using default output path
python3 scripts/generate_go_supported_libraries.py /path/to/opentelemetry-go-contrib

# Specifying custom output path
python3 scripts/generate_go_supported_libraries.py /path/to/opentelemetry-go-contrib -o custom/path/supported-libraries.yaml
```

### How It Works

The script:
1. Finds all instrumentation packages in the Go Contrib repository
2. Extracts supported version information from go.mod files
3. Converts version ranges to a consistent format
4. Generates a YAML file with the supported libraries and their version ranges

## Maintenance

The `supported-libraries.yaml` file should be updated when:
1. New instrumentations are added to the OpenTelemetry Go Contrib repository
2. Version ranges for existing instrumentations change
3. Instrumentation paths change

To update the file:
1. Update your local clone of the OpenTelemetry Go Contrib repository
2. Run the generation script
3. Review the changes in the generated YAML file
4. Commit the changes if they look correct

## Implementation Details

The Go checker scans `go.mod` files in your project to identify dependencies and compares them against the supported libraries list. It checks whether your dependencies are supported by OpenTelemetry instrumentation and reports any unsupported or out-of-range versions.

## Version Range Format

Version ranges are specified in the following format:
- `[1.0.0,2.0.0)` means version 1.0.0 (inclusive) to 2.0.0 (exclusive)
- `[1.0.0,)` means version 1.0.0 or higher
- `[,2.0.0)` means any version below 2.0.0 