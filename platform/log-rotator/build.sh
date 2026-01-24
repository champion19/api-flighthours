#!/bin/bash

set -e

echo "=== Flighthours Log Rotator - Build Script ==="
echo ""

PROJECT_NAME="flighthours-log-rotator"
IMAGE_TAG="latest"

cd "$(dirname "$0")"

echo "📦 Construyendo imagen Docker: ${PROJECT_NAME}:${IMAGE_TAG}"
echo ""

docker build \
    --tag ${PROJECT_NAME}:${IMAGE_TAG} \
    --label "com.flighthours.build-date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    --label "com.flighthours.version=1.0" \
    .

echo ""
echo "✅ Imagen construida exitosamente!"
echo ""
echo "Información de la imagen:"
docker images ${PROJECT_NAME}:${IMAGE_TAG}
echo ""
echo "Para ejecutar el contenedor:"
echo "  docker-compose up -d"
echo ""
echo "Para ver logs:"
echo "  docker-compose logs -f log-rotator"
echo ""
