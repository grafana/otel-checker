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

        # Look for the Supported Versions section (case insensitive)
        # Handle both ## and ### headers, and allow for different spacing
        versions_match = re.search(r'#{2,3}\s+Supported\s+Versions\n\n(.*?)(?:\n\n|$)', content, re.DOTALL | re.IGNORECASE)
        if not versions_match:
            print(f"Warning: No Supported Versions section found in {readme_path}", file=sys.stderr)
            return None

        versions_text = versions_match.group(1)

        # Get the directory name for srcPath
        dir_name = readme_path.parent.name

        # Try to extract library name from directory name
        # e.g., instrumentation-fs -> fs
        library_name = dir_name.replace('instrumentation-', '')

        # Pattern 1: [`library`](link) version(s) `>=0.5.5 <1`
        version_match = re.search(r'\[`(.*?)`\]\((.*?)\)\s+version[s]?\s+`(.*?)`', versions_text)
        if version_match:
            library_name = version_match.group(1)
            link = version_match.group(2)
            version_range = version_match.group(3)
            return {
                'name': library_name,
                'link': link,
                'version_range': version_range,
                'src_path': f"plugins/node/{dir_name}"
            }

        # Pattern 2: Node.js `>=14`
        node_match = re.search(r'Node\.js\s+`(.*?)`', versions_text)
        if node_match:
            version_range = node_match.group(1)
            return {
                'name': library_name,
                'link': f"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/{dir_name}",
                'version_range': version_range,
                'src_path': f"plugins/node/{dir_name}"
            }

        # Pattern 3: Library `>=1.0.0`
        lib_match = re.search(r'`(.*?)`\s+`(.*?)`', versions_text)
        if lib_match:
            library_name = lib_match.group(1)
            version_range = lib_match.group(2)
            return {
                'name': library_name,
                'link': f"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/{dir_name}",
                'version_range': version_range,
                'src_path': f"plugins/node/{dir_name}"
            }

        # Pattern 4: - Library `>=1.0.0`
        list_match = re.search(r'-\s+`(.*?)`\s+`(.*?)`', versions_text)
        if list_match:
            library_name = list_match.group(1)
            version_range = list_match.group(2)
            return {
                'name': library_name,
                'link': f"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/{dir_name}",
                'version_range': version_range,
                'src_path': f"plugins/node/{dir_name}"
            }

        # Pattern 5: - [library](link) `>=1.0.0`
        link_list_match = re.search(r'-\s+\[(.*?)\]\((.*?)\)\s+`(.*?)`', versions_text)
        if link_list_match:
            library_name = link_list_match.group(1)
            link = link_list_match.group(2)
            version_range = link_list_match.group(3)
            return {
                'name': library_name,
                'link': link,
                'version_range': version_range,
                'src_path': f"plugins/node/{dir_name}"
            }

        # Pattern 6: "regardless of versions" or similar
        all_versions_match = re.search(r'(?:\[`(.*?)`\]\((.*?)\)\s+)?(?:regardless of versions|all versions|any version)', versions_text, re.IGNORECASE)
        if all_versions_match:
            # If we have a library name and link, use them, otherwise use the directory name
            if all_versions_match.group(1):
                library_name = all_versions_match.group(1)
                link = all_versions_match.group(2)
            else:
                link = f"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/{dir_name}"
            
            return {
                'name': library_name,
                'link': link,
                'version_range': '>=0.0.0',  # This will be converted to [0.0.0,) in convert_version_range
                'src_path': f"plugins/node/{dir_name}"
            }

        print(f"Warning: Could not parse version information in {readme_path}", file=sys.stderr)
        return None
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
    elif len(parts) == 1:
        # Handle single version constraints
        if parts[0].startswith('>='):
            return f"[{parts[0].replace('>=', '')},)"
        elif parts[0].startswith('<'):
            return f"[,{parts[0].replace('<', '')})"
        elif parts[0].startswith('~'):
            # For tilde ranges, we'll use the same version for both bounds
            version = parts[0].replace('~', '')
            return f"[{version},{version})"
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
            print(f"Warning: No README.md found in {item}", file=sys.stderr)
            continue

        result = extract_supported_versions(readme_path)
        if result:
            library_name = result['name']
            supported_libraries[library_name] = {
                'instrumentations': [{
                    'name': library_name,
                    'srcPath': result['src_path'],
                    'link': result['link'],
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
