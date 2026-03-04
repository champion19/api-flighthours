#!/bin/bash
# ============================================================
# Export Antigravity Extensions
# Genera un archivo con todas las extensiones instaladas
# y un script para reinstalarlas en otro PC.
#
# Uso: ./export-extensions.sh
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
EXTENSIONS_FILE="$SCRIPT_DIR/extensions.txt"
INSTALL_SCRIPT="$SCRIPT_DIR/install-extensions.sh"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "🔍 Obteniendo extensiones de Antigravity..."

# Exportar lista de extensiones
antigravity --list-extensions > "$EXTENSIONS_FILE" 2>/dev/null

if [ $? -ne 0 ]; then
    echo "❌ Error: No se pudo ejecutar 'antigravity --list-extensions'"
    echo "   Asegúrate de que Antigravity esté instalado y el comando 'antigravity' esté en el PATH."
    exit 1
fi

TOTAL=$(wc -l < "$EXTENSIONS_FILE" | tr -d ' ')
echo "✅ Se encontraron $TOTAL extensiones"
echo "📄 Lista guardada en: $EXTENSIONS_FILE"

# Generar script de instalación
cat > "$INSTALL_SCRIPT" << 'HEADER'
#!/bin/bash
# ============================================================
# Install Antigravity Extensions
# Reinstala todas las extensiones exportadas.
#
# Uso: ./install-extensions.sh
#
# Opciones:
#   --dry-run        Solo muestra qué se instalaría, sin instalar
#   --skip-existing  Salta extensiones ya instaladas
# ============================================================

DRY_RUN=false
SKIP_EXISTING=false
INSTALLED=()

for arg in "$@"; do
    case $arg in
        --dry-run) DRY_RUN=true ;;
        --skip-existing) SKIP_EXISTING=true ;;
    esac
done

# Detectar el comando del editor
if command -v antigravity &> /dev/null; then
    EDITOR_CMD="antigravity"
elif command -v code &> /dev/null; then
    EDITOR_CMD="code"
else
    echo "❌ Error: No se encontró 'antigravity' ni 'code' en el PATH"
    exit 1
fi

echo "🔧 Usando comando: $EDITOR_CMD"

if $SKIP_EXISTING; then
    echo "📋 Obteniendo extensiones ya instaladas..."
    INSTALLED=($($EDITOR_CMD --list-extensions 2>/dev/null))
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
EXTENSIONS_FILE="$SCRIPT_DIR/extensions.txt"

if [ ! -f "$EXTENSIONS_FILE" ]; then
    echo "❌ Error: No se encontró $EXTENSIONS_FILE"
    exit 1
fi

TOTAL=$(wc -l < "$EXTENSIONS_FILE" | tr -d ' ')
CURRENT=0
SKIPPED=0
FAILED=0
SUCCESS=0

echo ""
echo "🚀 Instalando $TOTAL extensiones..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

while IFS= read -r extension || [ -n "$extension" ]; do
    [ -z "$extension" ] && continue
    CURRENT=$((CURRENT + 1))

    # Verificar si ya está instalada
    if $SKIP_EXISTING; then
        for installed in "${INSTALLED[@]}"; do
            if [ "$installed" = "$extension" ]; then
                echo "  ⏭️  [$CURRENT/$TOTAL] $extension (ya instalada)"
                SKIPPED=$((SKIPPED + 1))
                continue 2
            fi
        done
    fi

    if $DRY_RUN; then
        echo "  📦 [$CURRENT/$TOTAL] $extension (dry-run)"
        SUCCESS=$((SUCCESS + 1))
    else
        echo -n "  📦 [$CURRENT/$TOTAL] $extension... "
        if $EDITOR_CMD --install-extension "$extension" --force > /dev/null 2>&1; then
            echo "✅"
            SUCCESS=$((SUCCESS + 1))
        else
            echo "❌"
            FAILED=$((FAILED + 1))
        fi
    fi
done < "$EXTENSIONS_FILE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Resumen:"
echo "   ✅ Instaladas: $SUCCESS"
[ $SKIPPED -gt 0 ] && echo "   ⏭️  Omitidas:   $SKIPPED"
[ $FAILED -gt 0 ]  && echo "   ❌ Fallidas:   $FAILED"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
HEADER

chmod +x "$INSTALL_SCRIPT" "$0" 2>/dev/null

echo "📦 Script de instalación generado: $INSTALL_SCRIPT"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Extensiones exportadas ($TOTAL) — $TIMESTAMP"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat "$EXTENSIONS_FILE"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Para reinstalar en otro PC:"
echo "   1. Copia la carpeta '$SCRIPT_DIR' al nuevo PC"
echo "   2. Ejecuta: ./install-extensions.sh"
echo "   3. O para ver qué se instalaría: ./install-extensions.sh --dry-run"
echo "   4. Para saltar las que ya tengas: ./install-extensions.sh --skip-existing"
