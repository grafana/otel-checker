#!/usr/bin/env python3

import os
import re
import sys
import yaml
import argparse
from pathlib import Path

def extract_supported_versions(readme_path):
    """Extract supported versions from a README.md file."""
    try:
        with open(readme_path, 'r') as f:
            content = f.read()
            
        # Look for the Supported Versions section
        versions_match = re.search(r'## Supported Versions\n\n(.*?)(?:\n\n|$)', content, re.DOTALL)
        if not versions_match:
            return None
            
        versions_text = versions_match.group(1)
        
        # Extract library name and version range
        # Format: [`library`](link) versions `>=0.5.5 <1`
        version_match = re.search(r'\[`(.*?)`\]\((.*?)\)\s+versions\s+`(.*?)`', versions_text)
        if not version_match:
            return None
            
        library_name = version_match.group(1)
        link = version_match.group(2)
        version_range = version_match.group(3)
        
        # Get the directory name for srcPath
        dir_name = readme_path.parent.name
        
        return {
            'name': library_name,
            'link': link,
            'version_range': version_range,
            'src_path': f"plugins/node/{dir_name}"
        }
    except Exception as e:
        print(f"Error processing {readme_path}: {e}", file=sys.stderr)
        return None

def convert_version_range(version_range):
    """Convert version range string to YAML format."""
    # Example: ">=0.5.5 <1" -> [0.5.5,1)
    parts = version_range.split()
    if len(parts) == 2:
        min_version = parts[0].replace('>=', '')
        max_version = parts[1].replace('<', '')
        return f"[{min_version},{max_version})"
    return version_range

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from OpenTelemetry JS Contrib repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry JS Contrib repository')
    parser.add_argument('--output', '-o', default='checks/sdk/js/supported-libraries.yaml',
                      help='Output path for the YAML file (default: checks/sdk/js/supported-libraries.yaml)')
    args = parser.parse_args()
    
    # Path to the plugins directory
    plugins_dir = Path(args.repo_dir) / "plugins/node"
    if not plugins_dir.exists():
        print(f"Error: {plugins_dir} does not exist", file=sys.stderr)
        sys.exit(1)
    
    # Collect all supported libraries
    supported_libraries = {}
    
    # Process each instrumentation directory
    for item in plugins_dir.iterdir():
        if not item.is_dir():
            continue
            
        readme_path = item / "README.md"
        if not readme_path.exists():
            continue
            
        result = extract_supported_versions(readme_path)
        if result:
            library_name = result['name']
            supported_libraries[library_name] = {
                'instrumentations': [{
                    'name': library_name,
                    'srcPath': result['src_path'],
                    'target_versions': {
                        'LIBRARY': [convert_version_range(result['version_range'])]
                    }
                }]
            }
    
    # Generate the supported-libraries.yaml file
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w') as f:
        yaml.dump(supported_libraries, f, sort_keys=False)
    
    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main() 