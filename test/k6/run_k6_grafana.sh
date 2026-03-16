#!/bin/bash
# =====================================================
# FlightHours — k6 Load Test Runner con Grafana
# =====================================================
# Ejecuta tests de carga k6 con output a Prometheus,
# para visualización en tiempo real en Grafana.
#
# Uso:
#   ./test/k6/run_k6_grafana.sh              → escenario 'thousand' (default)
#   ./test/k6/run_k6_grafana.sh smoke         → smoke test
#   ./test/k6/run_k6_grafana.sh load          → carga normal
#   ./test/k6/run_k6_grafana.sh stress        → estrés
#   ./test/k6/run_k6_grafana.sh thousand      → 1000 VUs
#
# Requisitos:
#   - k6 instalado: brew install k6
#   - Docker compose de Grafana corriendo:
#       docker compose -f docker-compose.grafana.yml up -d
#   - Backend corriendo en localhost:8081
# =====================================================

set -euo pipefail

SCENARIO="${1:-thousand}"
PROMETHEUS_URL="${K6_PROMETHEUS_RW_SERVER_URL:-http://localhost:9090/api/v1/write}"
BASE_URL="${BASE_URL:-http://localhost:8081}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}============================================${NC}"
echo -e "${CYAN} FlightHours — k6 Load Testing + Grafana${NC}"
echo -e "${CYAN}============================================${NC}"
echo ""
echo -e "${GREEN}📊 Escenario:${NC}   $SCENARIO"
echo -e "${GREEN}🎯 Backend:${NC}     $BASE_URL"
echo -e "${GREEN}📡 Prometheus:${NC}  $PROMETHEUS_URL"
echo -e "${GREEN}📈 Grafana:${NC}     http://localhost:3000 (admin/0000)"
echo ""

# Verificar k6
if ! command -v k6 &> /dev/null; then
    echo -e "${RED}❌ k6 no está instalado. Instálalo con:${NC}"
    echo "   brew install k6"
    exit 1
fi

# Verificar que Prometheus esté corriendo
echo -e "${YELLOW}⏳ Verificando Prometheus...${NC}"
if curl -s --max-time 3 "http://localhost:9090/-/ready" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Prometheus OK${NC}"
else
    echo -e "${RED}⚠️  Prometheus no responde en localhost:9090${NC}"
    echo -e "${YELLOW}    Levántalo con: docker compose -f docker-compose.grafana.yml up -d${NC}"
    echo ""
    read -p "¿Continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Verificar que el backend esté corriendo
echo -e "${YELLOW}⏳ Verificando backend...${NC}"
if curl -s --max-time 3 "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend OK${NC}"
else
    echo -e "${RED}⚠️  Backend no responde en $BASE_URL${NC}"
    echo ""
    read -p "¿Continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo -e "${CYAN}🚀 Lanzando k6 — escenario: ${SCENARIO}${NC}"
echo -e "${YELLOW}   Abre Grafana en http://localhost:3000 para ver métricas en tiempo real${NC}"
echo -e "${YELLOW}   Dashboard: 'k6 Load Testing — FlightHours'${NC}"
echo ""

# Ejecutar k6 con output a Prometheus
k6 run \
    -o experimental-prometheus-rw \
    --env K6_PROMETHEUS_RW_SERVER_URL="$PROMETHEUS_URL" \
    --env SCENARIO="$SCENARIO" \
    --env BASE_URL="$BASE_URL" \
    "$PROJECT_ROOT/test/k6/load_test.js"

echo ""
echo -e "${GREEN}✅ Test completado — revisa los resultados en Grafana${NC}"
echo -e "${GREEN}   http://localhost:3000/d/k6-flighthours-load-testing${NC}"
