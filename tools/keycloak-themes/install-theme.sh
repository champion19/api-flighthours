#!/bin/bash

# Flighthours Keycloak Theme Installer
# Este script copia el theme personalizado al contenedor de Keycloak

set -e

echo "🏍️  Flighthours - Instalador de Themes para Keycloak"
echo "=================================================="
echo ""

# Variables
CONTAINER_NAME="keycloak-prod"
THEME_NAME="flighthours"
LOCAL_THEME_PATH="./tools/keycloak-themes/flighthours"
KEYCLOAK_THEMES_PATH="/opt/keycloak/themes"

# Verificar que el contenedor esté corriendo
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Error: El contenedor $CONTAINER_NAME no está corriendo."
    echo "   Ejecuta: docker-compose -f docker-compose.keycloak.yml up -d"
    exit 1
fi

echo "✓ Contenedor de Keycloak encontrado"
echo ""

# Verificar que exista la carpeta local del theme
if [ ! -d "$LOCAL_THEME_PATH" ]; then
    echo "❌ Error: No se encontró el theme en $LOCAL_THEME_PATH"
    exit 1
fi

echo "✓ Theme local encontrado"
echo ""

# Limpiar theme anterior si existe
echo "🧹 Limpiando theme anterior..."
docker exec --user root "$CONTAINER_NAME" rm -rf "$KEYCLOAK_THEMES_PATH/$THEME_NAME" 2>/dev/null || true
echo "✓ Limpieza completada"
echo ""

# Copiar theme al contenedor
echo "📦 Copiando theme al contenedor..."
# Copiar el contenido de flighthours/ directamente (incluyendo el .)
docker cp "$LOCAL_THEME_PATH/." "$CONTAINER_NAME:$KEYCLOAK_THEMES_PATH/$THEME_NAME/"

if [ $? -eq 0 ]; then
    echo "✓ Theme copiado exitosamente"
    echo ""
else
    echo "❌ Error copiando el theme"
    exit 1
fi


# Verificar la estructura copiada
echo "🔍 Verificando estructura..."
echo "Email theme:"
docker exec "$CONTAINER_NAME" ls -la "$KEYCLOAK_THEMES_PATH/$THEME_NAME/email/" | head -8
echo ""
echo "Login theme:"
docker exec "$CONTAINER_NAME" ls -la "$KEYCLOAK_THEMES_PATH/$THEME_NAME/login/" | head -8

# Reiniciar Keycloak para que cargue el theme
echo ""
echo "🔄 Reiniciando Keycloak..."
docker restart "$CONTAINER_NAME" > /dev/null 2>&1

echo "✓ Keycloak reiniciado"
echo ""

echo "=================================================="
echo "✅ ¡Instalación completada!"
echo ""
echo "📝 Próximos pasos:"
echo "   1. Espera 30 segundos a que Keycloak inicie"
echo "   2. Ir a: http://localhost:8080"
echo "   3. Login con tus credenciales"
echo "   4. Seleccionar realm 'flighthours'"
echo "   5. Realm Settings → Themes"
echo "   6. Email theme → seleccionar 'flighthours'"
echo "   7. Login theme → seleccionar 'flighthours'"
echo "   8. Save"
echo "   9. Probar: Realm Settings → Email → Test connection"
echo ""
echo "🎉 ¡Disfruta tus nuevos themes!"


