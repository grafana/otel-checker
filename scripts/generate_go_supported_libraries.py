# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "pyyaml",
# ]
# ///

#!/usr/bin/env python3

import os
import re
import sys
import yaml
import argparse
from pathlib import Path
from typing import Dict, Any, List, Optional
from common import Signals, Instrumentation, signals_match_file

TRACE_PATTERNS = [
    r'otel\.Tracer',
    r'TracerProvider',
    r'MeterProvider',
    r'SpanContextConfig',
    r'TraceState',
    r'TraceID',
    r'TraceFlags',
    r'SpanID',
    r'WithSpan',
    r'StartSpan'
]
METRIC_PATTERNS = [
    r'otel\.Meter',
    r'WithMeterProvider',
    r'metric\.Float64(Counter|UpDownCounter|Histogram)',
    r'metric\.Int64(Counter|UpDownCounter|Histogram)',
    r'NewInt64(Counter|UpDownCounter|Histogram)',
    r'NewFloat64(Counter|UpDownCounter|Histogram)',
    r'override _updateMetricInstruments'
]

def parse_go_mod_file(file_path: Path) -> Dict[str, Any]:
    """Parse a go.mod file and extract the module name and dependencies."""
    dependencies = {}
    module_name = None
    go_version = None
    
    with open(file_path, 'r') as f:
        content = f.read()

    # Extract the module name
    module_match = re.search(r'^module\s+(.+)$', content, re.MULTILINE)
    if module_match:
        module_name = module_match.group(1).strip()
    
    # Extract Go version
    go_match = re.search(r'^go\s+([0-9]+\.[0-9]+(?:\.[0-9]+)?)', content, re.MULTILINE)
    if go_match:
        go_version = go_match.group(1).strip()
    
    # Extract dependencies
    require_pattern = re.compile(r'^\t(.+?)\s+(v[0-9]+\.[0-9]+\.[0-9]+(?:-[a-zA-Z0-9.]+)?)$', re.MULTILINE)
    for match in require_pattern.finditer(content):
        dep_name = match.group(1).strip()
        dep_version = match.group(2).strip()
        dependencies[dep_name] = dep_version

    return {
        "module_name": module_name,
        "dependencies": dependencies,
        "go_version": go_version
    }

def calculate_version_range(version: str) -> str:
    """
    Calculate version range for Go dependencies where upper bound is next major version.
    Example: for v1.2.3, returns [1.2.3,2.0.0)
    """
    # Strip 'v' prefix if present
    clean_version = version.lstrip('v')

    # Parse major version
    major_version = int(clean_version.split('.')[0])

    # Calculate next major version
    next_major = major_version + 1

    # Create version range with upper bound as next major version
    return f"[{clean_version},{next_major}.0.0)"

def find_matching_dependency(file_path: Path, go_mod_data: Dict[str, Any]) -> Optional[Dict[str, Any]]:
    """
    Find the dependency that matches the directory structure of the go.mod file.

    Example:
    - Path: instrumentation/google.golang.org/grpc/otelgrpc/go.mod
    - Look for: google.golang.org/grpc
    """
    if not go_mod_data["module_name"]:
        return None

    # Get the relative path parts starting from the instrumentation directory
    rel_path_parts = []
    curr_dir = file_path.parent
    while curr_dir.name and curr_dir.name != "instrumentation":
        rel_path_parts.insert(0, curr_dir.name)
        curr_dir = curr_dir.parent

    if not rel_path_parts:
        return None

    # Build potential dependency prefixes based on directory structure
    dependencies = go_mod_data["dependencies"]

    # Look for a matching dependency
    for i in range(len(rel_path_parts) - 1, -1, -1):
        # Skip the last part if it looks like an instrumentation name (e.g., otelgrpc)
        if i == len(rel_path_parts) - 1 and rel_path_parts[i].startswith("otel"):
            continue

        # Try to build a path from components
        potential_path = "/".join(rel_path_parts[:i+1])

        # Assume this path is a built-in module if it doesn't contain a dot
        if not "." in potential_path:
            print("Potential path is a built-in module:", potential_path)
            return {
                "library": potential_path,
                "name": go_mod_data["module_name"],
                "version": go_mod_data.get("go_version", "unknown"),
                "module": go_mod_data["module_name"]
            }
        
        # Check if this path or any dependency starts with this path
        for dep_name in dependencies:
            if potential_path == dep_name or dep_name.startswith(f"{potential_path}/"):
                return {
                    "name": dep_name,
                    "version": dependencies[dep_name],
                    "module": go_mod_data["module_name"]
                }

    # If nothing matched directly, look at all parts
    for dep_name in dependencies:
        for part in rel_path_parts:
            if part in dep_name and not part.startswith("otel"):
                return {
                    "name": dep_name,
                    "version": dependencies[dep_name],
                    "module": go_mod_data["module_name"]
                }

    return None

def find_go_mod_files(repo_path: Path) -> List[Path]:
    """Find all go.mod files in the instrumentation directory."""
    instrumentation_dir = repo_path / "instrumentation"
    if not instrumentation_dir.exists():
        print(f"Error: {instrumentation_dir} does not exist", file=sys.stderr)
        sys.exit(1)

    return list(instrumentation_dir.glob("**/go.mod"))

def check_instrumentation_signals(src_dir: Path) -> Signals:
    """Check for traces and metrics signals in the instrumentation code."""
    signals = Signals()
    for file_path in src_dir.rglob('*.go'):
        signals.update(signals_match_file(file_path, 
                                               metric_patterns=METRIC_PATTERNS, 
                                               trace_patterns=TRACE_PATTERNS))

        if signals.traces and signals.metrics:
            break
    return signals

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from OpenTelemetry Go Contrib repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry Go Contrib repository')
    parser.add_argument('--output', '-o', default='checks/sdk/go/supported-libraries.yaml',
                        help='Output path for the YAML file (default: checks/sdk/go/supported-libraries.yaml)')
    args = parser.parse_args()

    repo_path = Path(args.repo_dir)
    go_mod_files = find_go_mod_files(repo_path)

    supported_libraries = {}

    print(f"Found {len(go_mod_files)} go.mod files")

    for go_mod_file in go_mod_files:
        if go_mod_file.parent.name in ["example", "test"]:
            continue

        try:
            go_mod_data = parse_go_mod_file(go_mod_file)
            rel_path = go_mod_file.relative_to(repo_path)

            matching_dep = find_matching_dependency(go_mod_file, go_mod_data)

            if matching_dep:
                library_name = matching_dep.get("library", False) or matching_dep["name"]
                module_name = matching_dep["module"]
                version = matching_dep["version"]

                # Calculate version range with proper upper bound
                version_range = calculate_version_range(version)

                # Check instrumentation signals
                signals = check_instrumentation_signals(go_mod_file.parent)

                print(f"Found match for {rel_path}:")
                print(f"  Library: {library_name}")
                print(f"  Version: {version}")
                print(f"  Version Range: {version_range}")
                print(f"  Module: {module_name}")
                print(f"  Signals: {signals}")

                entry = Instrumentation(
                    name=library_name,
                    source_path=str(rel_path.parent.as_posix()),
                    signals=signals,
                    link=module_name,
                    target_versions_library=[version_range]
                )

                if library_name not in supported_libraries:
                    supported_libraries[library_name] = [entry]
                else:
                    supported_libraries[library_name].append(entry)
            else:
                print(f"No matching dependency found for {rel_path}")
        except Exception as e:
            print(f"Error processing {go_mod_file}: {e}", file=sys.stderr)

    output_path = Path(args.output)

    if not output_path:
        print("Error: Output path is not specified or invalid.", file=sys.stderr)
        sys.exit(1)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    with open(output_path, 'w') as f:
        yaml.dump(supported_libraries, f, sort_keys=True)

    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main()

