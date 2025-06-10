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
from scripts.common import Signals, Instrumentation, signals_match_file

METRIC_PATTERNS = [
        r'from opentelemetry.metrics import',
        r'create_counter',
        r'create_up_down_counter',
        r'create_histogram',
        r'create_observable_gauge',
        r'Counter\(',
        r'UpDownCounter\(',
        r'Histogram\('
    ]

TRACE_PATTERNS = [
        r'from opentelemetry.trace import',
        r'SpanKind',
        r'start_as_current_span',
        r'start_span',
        r'set_span_in_context',
        r'set_attributes',
        r'add_event'
    ]

def check_instrumentation_signals(src_dir: Path) -> Signals:
    """Check if instrumentation supports traces and/or metrics."""
    signals = Signals()

    # First check package.py for explicit metrics support flag
    package_py = Path(src_dir) / "package.py"
    signals.update(signals_match_file(package_py, metric_patterns=[r'_supports_metrics\s*=\s*True']))

    # Walk through Python source files
    for file_path in src_dir.rglob('*.py(i)?'):
        if not file_path.is_file():
            continue
            
        signals.update(signals_match_file(file_path, 
                                           metric_patterns=METRIC_PATTERNS, 
                                           trace_patterns=TRACE_PATTERNS))

        if signals.traces and signals.metrics:
            break
    return signals

def find_instrumentation_dirs(repo_path: Path) -> List[Path]:
    """Find all Python instrumentation directories."""
    instrumentation_dir = repo_path / "instrumentation"
    if not instrumentation_dir.exists():
        print(f"Error: {instrumentation_dir} does not exist", file=sys.stderr)
        sys.exit(1)
    
    return [d for d in instrumentation_dir.iterdir() 
            if d.is_dir() and d.name.startswith('opentelemetry-instrumentation-')]

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from OpenTelemetry Python Contrib repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry Python Contrib repository')
    parser.add_argument('--output', '-o', default='checks/sdk/python/supported-libraries.yaml',
                        help='Output path for the YAML file')
    args = parser.parse_args()
    
    repo_path = Path(args.repo_dir)
    instrumentation_dirs = find_instrumentation_dirs(repo_path)
    
    supported_libraries = {}
    
    for inst_dir in instrumentation_dirs:
        # Extract library name from directory name
        lib_name = inst_dir.name.replace('opentelemetry-instrumentation-', '')
        
        # Check for signals in the instrumentation code
        src_dir = inst_dir / 'src' / 'opentelemetry' / 'instrumentation' / lib_name
        if not src_dir.exists():
            src_dir = inst_dir  # Fallback to main directory
            
        signals = check_instrumentation_signals(src_dir)            
        # Get relative path from repo root
        source_path = os.path.relpath(inst_dir, repo_path)
        
        # Create library entry
        supported_libraries[lib_name] =[Instrumentation(lib_name,
                                                 source_path,
                                                 signals,
                                                 f'https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-{lib_name}')]
    
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    sorted_supported_libraries = dict(sorted(supported_libraries.items()))
    
    with open(output_path, 'w') as f:
        yaml.dump(dict(sorted_supported_libraries), f, sort_keys=False)

    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main()
