import yaml
import csv
import os
from pathlib import Path

def process_python_js_style(data, language):
    """Process Python/JS style YAML where libraries are top-level keys."""
    results = []
    for lib_name, lib_data in data.items():
        if not isinstance(lib_data, list) or not lib_data:
            continue
        lib = lib_data[0]  # Take first entry
        if isinstance(lib, dict) and 'signals' in lib:
            results.append({
                'language': language,
                'library_name': lib.get('name', lib_name),
                'link': lib.get('link', ''),
                'metrics': lib.get('signals', {}).get('metrics', False),
                'traces': lib.get('signals', {}).get('traces', False)
            })
    return results

def process_java_style(data, language):
    """Process Java style YAML where libraries are under 'libraries' key."""
    results = []
    libraries = data.get('libraries', {})
    for lib_name, lib_data in libraries.items():
        if not isinstance(lib_data, list) or not lib_data:
            continue
        lib = lib_data[0]  # Take first entry
        if isinstance(lib, dict):
            results.append({
                'language': language,
                'library_name': lib.get('name', lib_name),
                'link': lib.get('link', ''),
                'metrics': lib.get('signals', {}).get('metrics', False),
                'traces': lib.get('signals', {}).get('traces', False)
            })
    return results

def process_yaml_file(file_path):
    """Process a YAML file and return list of library data."""
    language = file_path.parent.name  # Get language from parent directory name
    
    # Skip processing if language is 'sdk'
    if language == 'sdk':
        return []
    
    with open(file_path, 'r', encoding='utf-8') as f:
        try:
            data = yaml.safe_load(f)
            if not data:
                return []
                
            # Skip file header comments that might be included in the data
            if isinstance(data, dict):
                data = {k: v for k, v in data.items() if isinstance(v, (dict, list))}
            
            # Determine format and process accordingly
            if 'libraries' in data:
                return process_java_style(data, language)
            else:
                return process_python_js_style(data, language)
                
        except yaml.YAMLError as e:
            print(f"Error processing {file_path}: {e}")
            return []

def main():
    script_dir = Path(__file__).parent
    root_dir = script_dir.parent
    yaml_files = list(root_dir.rglob('supported-libraries.yaml'))

    if not yaml_files:
        print("No supported-libraries.yaml files found.")
        return
    
    all_libraries = []
    for yaml_file in yaml_files:
        libraries = process_yaml_file(yaml_file)
        all_libraries.extend(libraries)

    if not all_libraries:
        print("No libraries found in the YAML files.")
        return
    
    output_file = root_dir / 'supported_libraries.csv'
    fieldnames = ['language', 'library_name', 'link', 'metrics', 'traces']
    
    with open(output_file, 'w', newline='', encoding='utf-8') as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(all_libraries)
    
    print(f"Generated CSV file at: {output_file}")

if __name__ == '__main__':
    main()
