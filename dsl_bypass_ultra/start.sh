#!/bin/bash
# Bu betik, sistemi sistemin Python'u ile başlatır.
BASH_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
MAIN_PY="$BASH_DIR/main.py"

echo "Starting DSL Bypass Ultra System..."
python3 "$MAIN_PY" $@
