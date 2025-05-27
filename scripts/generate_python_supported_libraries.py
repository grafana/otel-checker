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

def check_instrumentation_signals(src_dir: Path) -> Dict[str, bool]:
    """Check if instrumentation supports traces and/or metrics."""
    signals = {}
    
    # First check package.py for explicit metrics support flag
    package_py = src_dir / "package.py"
    if package_py.exists():
        with open(package_py, 'r', encoding='utf-8') as f:
            content = f.read()
            if re.search(r'_supports_metrics\s*=\s*True', content):
                signals['metrics'] = True

    # Walk through Python source files
    for root, _, files in os.walk(src_dir):
        for file in files:
            if not file.endswith(('.py', '.pyi')):
                continue
                
            file_path = Path(root) / file
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            # Check for metrics support (if not already found in package.py)
            if 'metrics' not in signals:
                metric_patterns = [
                    r'from opentelemetry.metrics import',
                    r'create_counter',
                    r'create_up_down_counter',
                    r'create_histogram',
                    r'create_observable_gauge',
                    r'Counter\(',
                    r'UpDownCounter\(',
                    r'Histogram\('
                ]
                if any(re.search(p, content) for p in metric_patterns):
                    signals['metrics'] = True
            
            if 'traces' not in signals:
                # Check for tracing support
                trace_patterns = [
                    r'from opentelemetry.trace import',
                    r'SpanKind',
                    r'start_as_current_span',
                    r'start_span',
                    r'set_span_in_context',
                    r'set_attributes',
                    r'add_event'
                ]
                if any(re.search(p, content) for p in trace_patterns):
                    signals['traces'] = True
                
            if 'traces' in signals and 'metrics' in signals:  # Both signals found
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
        supported_libraries[lib_name] = {
            'instrumentations': [{
                'name': lib_name,
                'source_path': source_path,
                'signals': signals,
                'link': f'https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-{lib_name}'
            }]
        }
    
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    sorted_supported_libraries = dict(sorted(supported_libraries.items()))
    
    with open(output_path, 'w') as f:
        yaml.dump(dict(sorted_supported_libraries), f, sort_keys=False)

    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main()
