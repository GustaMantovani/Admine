"""Shared configurations for Admine Bot project tests."""

import sys
from pathlib import Path

# Add src directory to path to allow module imports
src_path = Path(__file__).parent.parent / "src"
sys.path.insert(0, str(src_path))
