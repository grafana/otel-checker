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
import xml.etree.ElementTree as ET

TRACE_PATTERNS = [r"public static\s+(?:async\s+)?(?:Task<)?TracerProviderBuilder(?:>)?\s+\w+\s*\(this\s+TracerProviderBuilder"]
METRIC_PATTERNS = [r"public static\s+(?:async\s+)?(?:Task<)?MeterProviderBuilder(?:>)?\s+\w+\s*\(this\s+MeterProviderBuilder"]

def parse_csproj(csproj_path: Path) -> Optional[Dict[str, str]]:
    try:
        tree = ET.parse(str(csproj_path))
        root = tree.getroot()
        
        name = None
        version = None

        # Common tags for package name
        name_tags = ['AssemblyName', 'PackageId']
        for tag in name_tags:
            element = root.find(f".//{tag}")
            if element is not None and element.text:
                name = element.text.strip()
                break
        # Fallback to file name if no tag found
        name = name or csproj_path.stem

        # Common tags for package version
        version_tags = ['Version', 'PackageVersion']
        for tag in version_tags:
            element = root.find(f".//{tag}")
            if element is not None and element.text:
                version = element.text.strip()
                break
        
        if name:
             return {'name': name, 'version': version or 'unknown'}
        return None

    except FileNotFoundError:
        print(f"Error: .csproj file not found at {csproj_path}", file=sys.stderr)
        return None
    except ET.ParseError:
        print(f"Error: Could not parse XML in .csproj file at {csproj_path}", file=sys.stderr)
        return None
    return None # Default return if critical info missing or error

def main():
    parser = argparse.ArgumentParser(description='Generate supported libraries YAML file from OpenTelemetry Go Contrib repository')
    parser.add_argument('repo_dir', help='Path to the OpenTelemetry Go Contrib repository')
    parser.add_argument('--output', '-o', default='checks/sdk/go/supported-libraries.yaml',
                        help='Output path for the YAML file (default: checks/sdk/go/supported-libraries.yaml)')
    args = parser.parse_args()

    repo_path = Path(args.repo_dir)
    supported_libraries = {}
    src_root_dir = repo_path / "src"

    if not src_root_dir.is_dir():
        print(f"Error: Source directory {src_root_dir} not found. Please check the repository path.", file=sys.stderr)
        sys.exit(1)

    csproj_files = list(src_root_dir.glob("**/*.csproj"))

    for csproj_file in csproj_files:
        print(f"Processing {csproj_file}")
        try:
            metadata = parse_csproj(csproj_file)
            if not metadata:
                print(f"Skipping {csproj_file}: No valid metadata found", file=sys.stderr)
                continue
            # Fall back to csproj file name if metadata is missing
            library_name = metadata.get('name', csproj_file.stem)
            library_version = metadata.get('version')
            project_dir = csproj_file.parent
            cs_files = list(project_dir.glob("**/*.cs"))

            for cs_file in cs_files:
                signals = signals_match_file(cs_file,
                                   metric_patterns=METRIC_PATTERNS,
                                   trace_patterns=TRACE_PATTERNS)
                if library_name and (signals.traces or signals.metrics):
                    source_path = project_dir.relative_to(repo_path).as_posix()
                    github_link = f"https://github.com/open-telemetry/opentelemetry-dotnet-contrib/tree/main/{relative_source_path}"

                    entry = Instrumentation(
                        name=library_name,
                        source_path=source_path,
                        signals=signals.to_dict(),
                        link=github_link,
                        target_versions_library=[library_version]
                    )

                    if library_name not in supported_libraries:
                        supported_libraries[library_name] = [entry]
                    else:
                        supported_libraries[library_name].append(entry)
        except Exception as e:
            print(f"Error processing {csproj_file}: {e}", file=sys.stderr)

    print(f"Supported libraries will be saved to {args.output}")
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


