#!/bin/bash

# Script para recuperar el stack de monitoreo cuando se desconectan las redes
# Uso: ./fix-monitoring-stack.sh
# Autor: Flighthours Team
# Fecha: 2025-12-13

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Reparación del Stack de Monitoreo Flighthours${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Navegar al directorio del proyecto
cd "$(dirname "$0")/../.."

echo -e "${YELLOW}📋 Paso 1: Deteniendo servicios de monitoreo...${NC}"
docker-compose -f docker-compose.grafana.yml down --remove-orphans 2>/dev/null || true
echo -e "${GREEN}✅ Servicios detenidos${NC}"
echo ""

echo -e "${YELLOW}📋 Paso 2: Levantando servicios en orden correcto...${NC}"
echo -e "${BLUE}   → Iniciando Loki...${NC}"
docker-compose -f docker-compose.grafana.yml up -d loki

echo -e "${BLUE}   → Iniciando Prometheus...${NC}"
docker-compose -f docker-compose.grafana.yml up -d prometheus

echo -e "${BLUE}   → Iniciando Promtail...${NC}"
docker-compose -f docker-compose.grafana.yml up -d promtail

echo -e "${BLUE}   → Iniciando Grafana...${NC}"
docker-compose -f docker-compose.grafana.yml up -d grafana

echo -e "${GREEN}✅ Todos los servicios iniciados${NC}"
echo ""

# Esperar a que los servicios estén listos
echo -e "${YELLOW}⏳ Esperando a que los servicios estén listos...${NC}"
sleep 3

echo -e "${YELLOW}📋 Paso 3: Verificando estado de la red...${NC}"
echo ""

# Función para verificar red de un contenedor
check_network() {
    local container=$1
    local network_info=$(docker inspect $container --format='{{range $net, $conf := .NetworkSettings.Networks}}{{$net}}{{end}}' 2>/dev/null || echo "NONE")

    if [[ "$network_info" == *"flighthours-network"* ]]; then
        echo -e "  ${GREEN}✅ $container → conectado a flighthours-network${NC}"
        return 0
    else
        echo -e "  ${RED}❌ $container → NO conectado (red: $network_info)${NC}"
        return 1
    fi
}

all_ok=true

# Verificar cada contenedor
check_network "flighthours-prometheus" || all_ok=false
check_network "flighthours-grafana" || all_ok=false
check_network "flighthours-loki" || all_ok=false
check_network "flighthours-promtail" || all_ok=false

echo ""

echo -e "${YELLOW}📋 Paso 4: Verificando conectividad...${NC}"
echo ""

# Verificar que Promtail puede resolver 'loki'
echo -e "${BLUE}   → Verificando Promtail → Loki...${NC}"
if docker exec motogo-promtail cat /etc/promtail/config.yml | grep -q "http://loki:3100"; then
    echo -e "     ${GREEN}✅ Configuración correcta${NC}"
else
    echo -e "     ${RED}❌ Configuración incorrecta${NC}"
    all_ok=false
fi

# Verificar que Grafana puede resolver 'prometheus'
echo -e "${BLUE}   → Verificando Grafana → Prometheus...${NC}"
if docker exec motogo-grafana cat /etc/grafana/provisioning/datasources/prometheus.yml | grep -q "http://prometheus:9090"; then
    echo -e "     ${GREEN}✅ Configuración correcta${NC}"
else
    echo -e "     ${RED}❌ Configuración incorrecta${NC}"
    all_ok=false
fi

# Verificar que Grafana puede resolver 'loki'
echo -e "${BLUE}   → Verificando Grafana → Loki...${NC}"
loki_check=$(docker exec motogo-grafana cat /etc/grafana/provisioning/datasources/loki.yml 2>/dev/null | grep -q "http://loki:3100" && echo "ok" || echo "fail")
if [ "$loki_check" == "ok" ]; then
    echo -e "     ${GREEN}✅ Configuración correcta${NC}"
else
    echo -e "     ${YELLOW}⚠️  Datasource no encontrado (puede ser normal si no está configurado)${NC}"
fi

echo ""

echo -e "${YELLOW}📋 Paso 5: Estado final de servicios...${NC}"
echo ""
docker ps --filter "name=flighthours-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ "$all_ok" = true ]; then
    echo -e "${GREEN}✅ Stack de monitoreo reparado exitosamente${NC}"
    echo ""
    echo -e "Accesos:"
    echo -e "  • Grafana:    ${BLUE}http://localhost:3000${NC} (admin/admin)"
    echo -e "  • Prometheus: ${BLUE}http://localhost:9091${NC}"
    echo -e "  • Loki:       ${BLUE}http://localhost:3100${NC}"
else
    echo -e "${RED}⚠️  Algunos servicios tienen problemas de conectividad${NC}"
    echo -e "${YELLOW}Ejecuta este script nuevamente o verifica los logs:${NC}"
    echo -e "  docker logs flighthours-prometheus"
    echo -e "  docker logs flighthours-loki"
    echo -e "  docker logs flighthours-promtail"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
