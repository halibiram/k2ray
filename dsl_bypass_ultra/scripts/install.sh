#!/bin/bash

# ==============================================================================
# DSL Bypass Ultra System - One-Click Installer (GÖREV 5 - Kararlı Sürüm)
# ==============================================================================
# Bu betik, sanal ortam oluşturmadan, doğrudan sistem Python'u üzerine
# kurulum yaparak ortam sorunlarını çözer.
# ==============================================================================

# Betiğin kendi konumuna göre mutlak yolları belirle.
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROJECT_ROOT_DIR=$(cd -- "$SCRIPT_DIR/.." &> /dev/null && pwd)
CONFIG_FILE="$PROJECT_ROOT_DIR/config.yaml"

echo "================================================="
echo "🚀 Starting DSL Bypass Ultra System Installation (No-venv mode)"
echo "================================================="

# --- Adım 1: Python 3'ü kontrol et ---
echo -e "\n[STEP 1/4] Checking for Python 3..."
if ! command -v python3 &> /dev/null; then
    echo "❌ CRITICAL ERROR: Python 3 is not available. Please install it."
    exit 1
fi
echo "✅ Python 3 found."

# --- Adım 2: Bağımlılıkları kur ---
echo -e "\n[STEP 2/4] Installing dependencies globally..."
pip install -r "$PROJECT_ROOT_DIR/requirements.txt"
if [ $? -ne 0 ]; then
    echo "❌ ERROR: Failed to install Python dependencies."
    exit 1
fi
echo "✅ Dependencies installed successfully."

# --- Adım 3: Modemi bul ve yapılandırma dosyasını oluştur ---
echo -e "\n[STEP 3/4] Discovering modem and generating config file..."
MODEM_IP=$(python3 "$SCRIPT_DIR/run_scanner.py" | tail -n 1)
if [[ -z "$MODEM_IP" ]]; then
    echo "⚠️ WARNING: Could not automatically discover modem IP. Using default."
    MODEM_IP="192.168.1.1"
fi
echo "-> Modem IP found/set to: $MODEM_IP"

python3 "$SCRIPT_DIR/run_config_generator.py" "$MODEM_IP"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ ERROR: Failed to create 'config.yaml'."
    exit 1
fi
echo "✅ Configuration file 'config.yaml' created successfully."

# --- Adım 4: Başlatma betiğini oluştur ---
echo -e "\n[STEP 4/4] Creating start script..."
cat > "$PROJECT_ROOT_DIR/start.sh" << EOL
#!/bin/bash
# Bu betik, sistemi sistemin Python'u ile başlatır.
BASH_DIR=\$(cd -- "\$(dirname -- "\${BASH_SOURCE[0]}")" &> /dev/null && pwd)
MAIN_PY="\$BASH_DIR/main.py"

echo "Starting DSL Bypass Ultra System..."
python3 "\$MAIN_PY" \$@
EOL

chmod +x "$PROJECT_ROOT_DIR/start.sh"
echo "✅ 'start.sh' script created successfully."

echo -e "\n================================================="
echo "✅🎉 INSTALLATION COMPLETE! 🎉✅"
echo "================================================="
echo "NEXT STEPS:"
echo "1. IMPORTANT: Edit 'config.yaml' in the project root to set your modem's password."
echo "   -> vim $CONFIG_FILE"
echo ""
echo "2. To run the system, execute from the project root:"
echo "   -> ./start.sh"
echo "================================================="