#!/usr/bin/env python3

import os
import re
import subprocess
import sys
from pathlib import Path

def clone_repo():
    """Clone the OpenTelemetry JS Contrib repository if it doesn't exist."""
    if not os.path.exists("opentelemetry-js-contrib"):
        print("Cloning OpenTelemetry JS Contrib repository...")
        subprocess.run(["git", "clone", "https://github.com/open-telemetry/opentelemetry-js-contrib.git"], check=True)

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
        
        return {
            'name': library_name,
            'link': link,
            'version_range': version_range
        }
    except Exception as e:
        print(f"Error processing {readme_path}: {e}", file=sys.stderr)
        return None

def main():
    # Clone the repository
    #clone_repo()
    
    # Path to the plugins directory
    plugins_dir = Path("opentelemetry-js-contrib/plugins/node")
    if not plugins_dir.exists():
        print(f"Error: {plugins_dir} does not exist", file=sys.stderr)
        sys.exit(1)
    
    # Collect all supported libraries
    supported_libraries = []
    
    # Process each instrumentation directory
    for item in plugins_dir.iterdir():
        if not item.is_dir():
            continue
            
        readme_path = item / "README.md"
        if not readme_path.exists():
            continue
            
        result = extract_supported_versions(readme_path)
        if result:
            supported_libraries.append(result)
    
    # Generate the supported-libraries.md file
    output_path = Path("checks/sdk/js/supported-libraries.md")
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w') as f:
        for lib in supported_libraries:
            f.write(f"[{lib['name']}]({lib['link']}) versions `{lib['version_range']}`\n")
    
    print(f"Generated {output_path} with {len(supported_libraries)} supported libraries")

if __name__ == "__main__":
    main() 