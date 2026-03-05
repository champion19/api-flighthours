#!/bin/bash

# ============================================================
# EXPORTACIÓN DE ESTRUCTURAS DE BASES DE DATOS - FLIGHTHOURS
# ============================================================
#
# Arquitectura:
# - Servidor MySQL: Docker container -> mysql-flighthours
# - Base de datos de negocio: flightDb
# - Base de datos de autenticación (Keycloak): keycloakDb
#
# Este script exporta SOLO LA ESTRUCTURA (sin datos)
# Los archivos se guardan en documents/migracion/
# Ideal para versionar en Git.
#
# ============================================================

set -e

MYSQL_CONTAINER="mysql-flighthours"
MYSQL_USER="root"
MYSQL_PASSWORD="1997"

# Directorio raíz del proyecto (relativo al script)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

OUTPUT_NEGOCIO="$PROJECT_ROOT/documents/migracion/negocio-dump"
OUTPUT_KEYCLOAK="$PROJECT_ROOT/documents/migracion/keycloak-dump"

# Crear directorios si no existen
mkdir -p "$OUTPUT_NEGOCIO"
mkdir -p "$OUTPUT_KEYCLOAK"

# Verificar que el contenedor MySQL está corriendo
if ! docker ps --format '{{.Names}}' | grep -q "^${MYSQL_CONTAINER}$"; then
  echo "❌ Error: El contenedor '$MYSQL_CONTAINER' no está corriendo." >&2
  echo "   Inicia el contenedor antes de ejecutar este script." >&2
  exit 1
fi

# ============================================================
# 1️⃣ Exportar estructura base de negocio (flightDb)
# ============================================================

echo "📦 Exportando estructura de negocio (flightDb)..."

docker exec "$MYSQL_CONTAINER" \
  mysqldump -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" \
  --no-data \
  flightDb 2>/dev/null > "$OUTPUT_NEGOCIO/schema-flighthours.sql"

echo "✔ Archivo generado: documents/migracion/negocio-dump/schema-flighthours.sql"

# ============================================================
# 2️⃣ Exportar estructura base de Keycloak (keycloakDb)
# ============================================================

echo "🔐 Exportando estructura de autenticación (keycloakDb)..."

docker exec "$MYSQL_CONTAINER" \
  mysqldump -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" \
  --no-data \
  keycloakDb 2>/dev/null > "$OUTPUT_KEYCLOAK/schema-keycloak.sql"

echo "✔ Archivo generado: documents/migracion/keycloak-dump/schema-keycloak.sql"

# ============================================================
# ✅ Resumen
# ============================================================

echo ""
echo "============================================================"
echo "  ✅ Proceso finalizado correctamente"
echo "============================================================"
echo "  📁 Negocio:   $OUTPUT_NEGOCIO/schema-flighthours.sql"
echo "  📁 Keycloak:  $OUTPUT_KEYCLOAK/schema-keycloak.sql"
echo "============================================================"
