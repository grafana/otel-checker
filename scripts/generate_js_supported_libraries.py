#!/usr/bin/env python3

import os
import re
import sys
import yaml
import argparse
from pathlib import Path
from typing import Optional, Dict, Any

def create_result(library_name: str, link: str, version_range: str, dir_name: str) -> Dict[str, Any]:
    """Create a standardized result dictionary."""
    return {
        'name': library_name,
        'link': link,
        'version_range': version_range,
        'src_path': f"plugins/node/{dir_name}"
    }

def get_repo_link(dir_name: str) -> str:
    """Get the repository link for a library."""
    return f"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/{dir_name}"

def extract_supported_versions(readme_path: Path) -> Optional[Dict[str, Any]]:
    """Extract supported versions from a README.md file."""
    try:
        with open(readme_path, 'r') as f:
            content = f.read()

        # Look for the Supported Versions section (case insensitive)
        versions_match = re.search(r'#{2,3}\s+Supported\s+Versions\n\n(.*?)(?:\n\n|$)', content, re.DOTALL | re.IGNORECASE)
        if not versions_match:
            print(f"Warning: No Supported Versions section found in {readme_path}", file=sys.stderr)
            return None

        versions_text = versions_match.group(1)
        dir_name = readme_path.parent.name
        library_name = dir_name.replace('instrumentation-', '')
        repo_link = get_repo_link(dir_name)

        # Define patterns to match different version formats
        patterns = [
            # Pattern 1: [`library`](link) version(s) `>=0.5.5 <1`
            (r'\[`(.*?)`\]\((.*?)\)\s+version[s]?\s+`(.*?)`',
             lambda m: create_result(m.group(1), repo_link, m.group(3), dir_name)),

            # Pattern 2: Node.js `>=14`
            (r'Node\.js\s+`(.*?)`',
             lambda m: create_result(library_name, repo_link, m.group(1), dir_name)),

            # Pattern 3: Library `>=1.0.0`
            (r'`(.*?)`\s+`(.*?)`',
             lambda m: create_result(m.group(1), repo_link, m.group(2), dir_name)),

            # Pattern 4: - Library `>=1.0.0`
            (r'-\s+`(.*?)`\s+`(.*?)`',
             lambda m: create_result(m.group(1), repo_link, m.group(2), dir_name)),

            # Pattern 5: - [library](link) version(s) `>=1.0.0`
            (r'-\s+\[(.*?)\]\((.*?)\)\s+version[s]?\s+`(.*?)`',
             lambda m: create_result(m.group(1), repo_link, m.group(3), dir_name)),

            # Pattern 6: - [library](link) `>=1.0.0` (without "versions" word)
            (r'-\s+\[(.*?)\]\((.*?)\)\s+`(.*?)`',
             lambda m: create_result(m.group(1), repo_link, m.group(3), dir_name)),

            # Pattern 7: "regardless of versions" or similar
            (r'(?:\[`(.*?)`\]\((.*?)\)\s+)?(?:regardless of versions|all versions|any version)',
             lambda m: create_result(
                 m.group(1) if m.group(1) else library_name,
                 repo_link,
                 '>=0.0.0',
                 dir_name
             ))
        ]

        # Try each pattern
        for pattern, handler in patterns:
            match = re.search(pattern, versions_text, re.IGNORECASE)
            if match:
                return handler(match)

        print(f"Warning: Could not parse version information in {readme_path}", file=sys.stderr)
        return None
    except Exception as e:
        print(f"Error processing {readme_path}: {e}", file=sys.stderr)
        return None

def convert_version_range(version_range: str) -> str:
    """Convert version range string to YAML format."""
    parts = version_range.split()
    if len(parts) == 2:
        min_version = parts[0].replace('>=', '')
        max_version = parts[1].replace('<', '')
        return f"[{min_version},{max_version})"
    elif len(parts) == 1:
        if parts[0].startswith('>='):
            return f"[{parts[0].replace('>=', '')},)"
        elif parts[0].startswith('<'):
            return f"[,{parts[0].replace('<', '')})"
    return version_range

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from OpenTelemetry JS Contrib repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry JS Contrib repository')
    parser.add_argument('--output', '-o', default='checks/sdk/js/supported-libraries.yaml',
                      help='Output path for the YAML file (default: checks/sdk/js/supported-libraries.yaml)')
    args = parser.parse_args()

    plugins_dir = Path(args.repo_dir) / "plugins/node"
    if not plugins_dir.exists():
        print(f"Error: {plugins_dir} does not exist", file=sys.stderr)
        sys.exit(1)

    supported_libraries = {}
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
                        'library': [convert_version_range(result['version_range'])]
                    }
                }]
            }

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    with open(output_path, 'w') as f:
        yaml.dump(supported_libraries, f, sort_keys=False)

    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main()
