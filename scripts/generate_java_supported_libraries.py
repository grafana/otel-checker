# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "pyyaml",
#     "requests",
# ]
# ///

#!/usr/bin/env python3

import argparse
import sys
import yaml
import urllib.request
from pathlib import Path
from typing import Dict, Any, List, Optional
from common import Signals, Instrumentation, signals_match_file

TRACE_PATTERNS: List[str] = [
    r"import\s+io\.opentelemetry\.(?:api|sdk)\.trace\.(?:[^;]+);",
    r"import\s+io\.opentelemetry\.extension\.annotations\.WithSpan;",
    r"@WithSpan"
]
METRIC_PATTERNS: List[str] = [
    r"import\s+io\.opentelemetry\.(?:api|sdk)\.metrics\.(?:[^;]+);",
    # Consider adding common metric instrument class usages if imports are not always explicit
    # r"Meter\.counterBuilder\(", 
    # r"Meter\.histogramBuilder\(",
    # r"Meter\.gaugeBuilder\(" 
]

DEFAULT_YAML_URL = "https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml"

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from existing file in OpenTelemetry Java repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry Java repository')
    parser.add_argument('--input-yaml-url', default=DEFAULT_YAML_URL, help='URL to the existing Java supported-libraries.yaml file.')
    parser.add_argument('--output', '-o', default='checks/sdk/java/supported-libraries.yaml',
                      help='Output path for the YAML file (default: checks/sdk/java/supported-libraries.yaml)')
    args = parser.parse_args()

    repo_scan_path = Path(args.repo_dir).resolve()
    output_file_path = Path(args.output)

    if not repo_scan_path.is_dir():
        print(f"Error: Provided repository path '{repo_scan_path}' is not a valid directory.", file=sys.stderr)
        print(f"Directory exists: {repo_scan_path.exists()}", file=sys.stderr)
        # is file?
        print(f"Is a file: {repo_scan_path.is_file()}", file=sys.stderr)
        sys.exit(1)

    print(f"Input YAML URL: {args.input_yaml_url}")
    print(f"Java Instrumentation Repo Path: {repo_scan_path}")
    print(f"Output file: {output_file_path}")

    raw_yaml_content = None
    try:
        with urllib.request.urlopen(args.input_yaml_url) as response:
            raw_yaml_content = response.read().decode('utf-8')
            print(f"Successfully downloaded YAML from {args.input_yaml_url}")
    except Exception as e:
        print(f"Error downloading YAML from {args.input_yaml_url}: {e}", file=sys.stderr)
        sys.exit(1)

    if not raw_yaml_content:
        print("Error: No content downloaded from YAML URL.", file=sys.stderr)
        sys.exit(1)
        
    parsed_yaml_data: Dict[str, Any] = {}
    try:
        parsed_yaml_data = yaml.safe_load(raw_yaml_content)
        if not isinstance(parsed_yaml_data, dict):
            print(f"Error: Downloaded YAML does not seem to be a valid dictionary structure. Top level type is {type(parsed_yaml_data)}", file=sys.stderr)
            sys.exit(1)
        if "libraries" not in parsed_yaml_data:
            print("Error: Downloaded YAML does not contain 'libraries' key.", file=sys.stderr)
            sys.exit(1)
        print(f"Successfully parsed YAML. Found {len(parsed_yaml_data['libraries'])} library entries.")
    except yaml.YAMLError as e:
        print(f"Error parsing YAML content: {e}", file=sys.stderr)
        sys.exit(1)

    processed_instrumentations_count = 0
    libraries = parsed_yaml_data.get("libraries", {})
    for lib_name, lib_data in libraries.items():
        for inst_entry in lib_data:
            if not isinstance(inst_entry, dict):
                print(f"Warning: Skipping an instrumentation entry for '{lib_name}' as it's not a dictionary.", file=sys.stderr)
                continue

            source_path_str = inst_entry.get("source_path")
            if not source_path_str or not isinstance(source_path_str, str):
                # print(f"Warning: No 'source_path' found or invalid for an instrumentation under '{lib_name}'. Signals cannot be determined.", file=sys.stderr)
                inst_entry['signals'] = {'traces': False, 'metrics': False, 'note': 'Source path missing or invalid'}
                continue

            # Construct the full path. Assume source_path_str is relative to the repo_scan_path.
            # In java-instrumentation, source_path often looks like:
            # "instrumentation/apache-httpasyncclient-4.1/javaagent/src/main/java/io/opentelemetry/javaagent/instrumentation/apachehttpasyncclient"
            # So, repo_scan_path should be the root of the 'opentelemetry-java-instrumentation' checkout.
            source_code_full_path = repo_scan_path / source_path_str
            print(f"Processing instrumentation '{lib_name}' with source path '{source_path_str}'")

            if not source_code_full_path.is_dir():
                inst_entry['signals'] = {'traces': False, 'metrics': False, 'note': f'Source path {source_path_str} not found in repo.'}
                continue

            signals = Signals()

            for file in source_code_full_path.rglob('*.java'):
                signals.update(signals_match_file(file, TRACE_PATTERNS, METRIC_PATTERNS))
                if signals.traces and signals.metrics:
                    break

            inst_entry['signals'] = signals.to_dict()
            processed_instrumentations_count += 1

    print(f"Finished processing. Added signal data to {processed_instrumentations_count} instrumentations.")

    try:
        output_file_path.parent.mkdir(parents=True, exist_ok=True)
        with open(output_file_path, 'w') as f:
            yaml.dump(parsed_yaml_data, f, sort_keys=False, indent=2) # sort_keys=False to preserve original order from input YAML
        print(f"Successfully wrote augmented YAML to {output_file_path}")
    except IOError as e:
        print(f"Error writing YAML to {output_file_path}: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e: # Catch any other unexpected errors during writing
        print(f"An unexpected error occurred while writing the output YAML: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
