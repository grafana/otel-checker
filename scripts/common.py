from pathlib import Path
import re
import sys
from typing import Dict, Any, List, Optional

class Signals:
    """Class to hold signals for metrics and traces."""

    def __init__(self, metrics: bool = False, traces: bool = False):
        self.metrics = metrics
        self.traces = traces

    def __repr__(self) -> str:
        return f"Signals(metrics={self.metrics}, traces={self.traces})"

    def to_dict(self) -> Dict[str, bool]:
        return {'metrics': self.metrics, 'traces': self.traces}
    
    def update(self, other: 'Signals') -> None:
        """Update current signals with another Signals instance."""
        self.metrics = self.metrics or other.metrics
        self.traces = self.traces or other.traces

class Instrumentation:
    def __init__(self, name: str, source_path: str, signals: Signals, link: str, target_versions_library: List[str] = []):
        self.name = name
        self.source_path = source_path
        self.signals = signals
        self.link = link
        self.target_versions_library = target_versions_library

    def to_dict(self) -> Dict[str, Any]:
        return {
            'name': self.name,
            'source_path': self.source_path,
            'signals': self.signals.to_dict(),
            'link': self.link,
            'target_versions': {
                'library': self.target_versions_library
            }
        }

def signals_match_file(file: Path, metric_patterns: List[str] = [], trace_patterns: List[str] = []) -> Signals:
    """Check if the file contains any metric or trace patterns."""
    if not file.exists():
        print(f"Warning: {file} does not exist", file=sys.stderr)
        return Signals()

    signals = Signals()
    
    with open(file, 'r', encoding='utf-8') as f:
        content = f.read()
        
    if any(re.search(pattern, content) for pattern in metric_patterns):
        signals.metrics = True
        
    if any(re.search(pattern, content) for pattern in trace_patterns):
        signals.traces = True

    return signals
