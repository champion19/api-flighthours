#!/bin/bash

# Script para recuperar Swagger UI cuando se desconecta de la red
# Uso: ./fix-swagger-ui.sh
# Autor: Flighthours Team
# Fecha: 2025-12-13

set -e

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Reparación de Swagger UI${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Navegar al directorio del proyecto
cd "$(dirname "$0")/../.."

echo -e "${YELLOW}📋 Deteniendo Swagger UI...${NC}"
docker-compose -f docker-compose.swagger-ui.yml down --remove-orphans 2>/dev/null || true
echo -e "${GREEN}✅ Detenido${NC}"
echo ""

echo -e "${YELLOW}📋 Iniciando Swagger UI...${NC}"
docker-compose -f docker-compose.swagger-ui.yml up -d
echo -e "${GREEN}✅ Iniciado${NC}"
echo ""

# Esperar a que el servicio esté listo
echo -e "${YELLOW}⏳ Esperando a que el servicio esté listo...${NC}"
sleep 2

echo -e "${YELLOW}📋 Verificando estado...${NC}"
network_info=$(docker inspect flighthours-swagger-ui --format='{{range $net, $conf := .NetworkSettings.Networks}}{{$net}}{{end}}' 2>/dev/null || echo "NONE")

if [[ "$network_info" == *"flighthours-network"* ]]; then
    echo -e "${GREEN}✅ Swagger UI conectado a flighthours-network${NC}"
    echo ""
    echo -e "${GREEN}✅ Swagger UI reparado exitosamente${NC}"
    echo ""
    echo -e "Acceso: ${BLUE}http://localhost:3001${NC}"
else
    echo -e "${RED}❌ Swagger UI NO está conectado a la red correcta${NC}"
    echo -e "Red actual: $network_info"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
