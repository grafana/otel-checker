#!/usr/bin/env python3

import os
import re
import sys
import yaml
import argparse
from pathlib import Path
from typing import Dict, Any, List, Optional

def parse_go_mod_file(file_path: Path) -> Dict[str, Any]:
    """Parse a go.mod file and extract the module name and dependencies."""
    dependencies = {}
    module_name = None
    
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Extract the module name
    module_match = re.search(r'^module\s+(.+)$', content, re.MULTILINE)
    if module_match:
        module_name = module_match.group(1).strip()
    
    # Extract dependencies
    require_pattern = re.compile(r'^\t(.+?)\s+(v[0-9]+\.[0-9]+\.[0-9]+(?:-[a-zA-Z0-9.]+)?)$', re.MULTILINE)
    for match in require_pattern.finditer(content):
        dep_name = match.group(1).strip()
        dep_version = match.group(2).strip()
        dependencies[dep_name] = dep_version
    
    return {
        "module_name": module_name,
        "dependencies": dependencies
    }

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
    for i in range(len(rel_path_parts)):
        # Skip the last part if it looks like an instrumentation name (e.g., otelgrpc)
        if i == len(rel_path_parts) - 1 and rel_path_parts[i].startswith("otel"):
            continue
            
        # Try to build a path from components
        potential_path = "/".join(rel_path_parts[:i+1])
        
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
        try:
            go_mod_data = parse_go_mod_file(go_mod_file)
            rel_path = go_mod_file.relative_to(repo_path)
            
            matching_dep = find_matching_dependency(go_mod_file, go_mod_data)
            
            if matching_dep:
                library_name = matching_dep["name"]
                module_name = matching_dep["module"]
                version = matching_dep["version"]
                
                print(f"Found match for {rel_path}:")
                print(f"  Library: {library_name}")
                print(f"  Version: {version}")
                print(f"  Module: {module_name}")
                
                if library_name not in supported_libraries:
                    supported_libraries[library_name] = {
                        "instrumentations": [{
                            "name": library_name,
                            "srcPath": str(rel_path.parent),
                            "link": module_name,
                            "target_versions": {
                                "library": [f"[{version.lstrip('v')},)"]
                            }
                        }]
                    }
            else:
                print(f"No matching dependency found for {rel_path}")
        except Exception as e:
            print(f"Error processing {go_mod_file}: {e}", file=sys.stderr)
    
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w') as f:
        yaml.dump(supported_libraries, f, sort_keys=False)
    
    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main() 