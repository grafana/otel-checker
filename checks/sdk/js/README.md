# JavaScript Supported Libraries

This directory contains the configuration for supported JavaScript libraries that can be instrumented with OpenTelemetry.

## Supported Libraries File

The `supported-libraries.yaml` file contains information about which JavaScript libraries are supported by OpenTelemetry instrumentation, along with their version ranges and source paths.

### File Format

```yaml
library_name:
  instrumentations:
    - name: library_name
      srcPath: plugins/node/instrumentation-name
      target_versions:
        LIBRARY:
          - [min_version,max_version)
```

Example:
```yaml
amqplib:
  instrumentations:
    - name: amqplib
      srcPath: plugins/node/instrumentation-amqplib
      target_versions:
        LIBRARY:
          - [0.5.5,1)
```

## Generating the File

The `supported-libraries.yaml` file is generated from the OpenTelemetry JS Contrib repository using the `generate_js_supported_libraries.py` script.

### Prerequisites

- Python 3.x
- PyYAML package (`pip install pyyaml`)
- A local clone of the [OpenTelemetry JS Contrib repository](https://github.com/open-telemetry/opentelemetry-js-contrib)

### Usage

```bash
# Using default output path
python3 scripts/generate_js_supported_libraries.py /path/to/opentelemetry-js-contrib

# Specifying custom output path
python3 scripts/generate_js_supported_libraries.py /path/to/opentelemetry-js-contrib -o custom/path/supported-libraries.yaml
```

### How It Works

The script:
1. Reads README.md files from each instrumentation directory in the JS Contrib repository
2. Extracts supported version information from the "Supported Versions" section
3. Converts version ranges to a consistent format (e.g., `>=0.5.5 <1` becomes `[0.5.5,1)`)
4. Generates a YAML file with the supported libraries and their version ranges

## Maintenance

The `supported-libraries.yaml` file should be updated when:
1. New instrumentations are added to the OpenTelemetry JS Contrib repository
2. Version ranges for existing instrumentations change
3. Instrumentation paths change

To update the file:
1. Update your local clone of the OpenTelemetry JS Contrib repository
2. Run the generation script
3. Review the changes in the generated YAML file
4. Commit the changes if they look correct

## Version Range Format

The script converts version ranges from the format in README files to a consistent format:
- `>=0.5.5 <1` becomes `[0.5.5,1)`
- `>=1.0.0` becomes `[1.0.0,)`
- `<2.0.0` becomes `[,2.0.0)`

This format is compatible with the version range parsing in the JavaScript checker. 