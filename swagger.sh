#!/bin/bash

# Script auxiliar para trabajar con Swagger en Flighthours API
# Uso: ./swagger.sh [comando]

set -e

SWAGGER_DIR="platform/swaggo"
DOCKER_COMPOSE_FILE="docker-compose.swagger.yml"
MAIN_FILE="cmd/main.go"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar ayuda
show_help() {
    echo -e "${BLUE}=== Flighthours API - Swagger Helper ===${NC}"
    echo ""
    echo "Uso: ./swagger.sh [comando]"
    echo ""
    echo "Comandos disponibles:"
    echo ""
    echo -e "  ${GREEN}generate${NC}        Generar documentación Swagger desde anotaciones"
    echo -e "  ${GREEN}ui${NC}              Iniciar Swagger UI en Docker (puerto 8080)"
    echo -e "  ${GREEN}ui-stop${NC}         Detener Swagger UI"
    echo -e "  ${GREEN}serve${NC}           Iniciar API con Swagger UI integrado (puerto 8081)"
    echo -e "  ${GREEN}validate${NC}        Validar la documentación generada"
    echo -e "  ${GREEN}format${NC}          Formatear anotaciones en el código"
    echo -e "  ${GREEN}clean${NC}           Limpiar archivos generados"
    echo -e "  ${GREEN}install${NC}         Instalar swag CLI localmente"
    echo -e "  ${GREEN}docker-build${NC}    Construir imagen Docker del generador"
    echo -e "  ${GREEN}open${NC}            Abrir Swagger UI en el navegador"
    echo -e "  ${GREEN}help${NC}            Mostrar esta ayuda"
    echo ""
    return 0
}

# Función para generar documentación
generate_docs() {
    echo -e "${BLUE}📝 Generando documentación Swagger...${NC}"

    if command -v swag &> /dev/null; then
        echo -e "${GREEN}Usando swag CLI local${NC}"
        swag init -g "$MAIN_FILE" -o "$SWAGGER_DIR" --parseDependency --parseInternal
    else
        echo -e "${YELLOW}swag CLI no encontrado. Usando Docker...${NC}"
        if [[ ! "$(docker images -q flighthours-swag 2> /dev/null)" ]]; then
            echo -e "${YELLOW}Construyendo imagen Docker...${NC}"
            docker build -f "$SWAGGER_DIR/Dockerfile.swag" -t flighthours-swag .
        fi
        docker run --rm -v "$(pwd):/app" flighthours-swag \
            init -g "$MAIN_FILE" -o "$SWAGGER_DIR" --parseDependency --parseInternal
    fi

    echo -e "${GREEN}✓ Documentación generada en $SWAGGER_DIR/${NC}"
    echo -e "${BLUE}Archivos generados:${NC}"
    ls -lh "$SWAGGER_DIR"/{docs.go,swagger.json,swagger.yaml} 2>/dev/null || true
    return 0
}

# Función para iniciar Swagger UI
start_ui() {
    echo -e "${BLUE}🚀 Iniciando Swagger UI en Docker...${NC}"
    docker compose -f "$DOCKER_COMPOSE_FILE" up -d swagger-ui
    echo -e "${GREEN}✓ Swagger UI disponible en: http://localhost:8082${NC}"
    echo -e "${YELLOW}💡 Tip: Usa './swagger.sh open' para abrir en el navegador${NC}"
    return 0
}

# Función para detener Swagger UI
stop_ui() {
    echo -e "${BLUE}🛑 Deteniendo Swagger UI...${NC}"
    docker compose -f "$DOCKER_COMPOSE_FILE" down
    echo -e "${GREEN}✓ Swagger UI detenido${NC}"
    return 0
}

# Función para iniciar la API
serve_api() {
    echo -e "${BLUE}🚀 Iniciando Flighthours API con Swagger integrado...${NC}"
    echo -e "${YELLOW}Asegúrate de haber generado la documentación primero${NC}"
    echo ""
    go run "$MAIN_FILE"
    return 0
}

# Función para validar documentación
validate_docs() {
    echo -e "${BLUE}🔍 Validando documentación Swagger...${NC}"

    if [[ ! -f "$SWAGGER_DIR/swagger.json" ]]; then
        echo -e "${RED}✗ Error: swagger.json no existe. Ejecuta './swagger.sh generate' primero${NC}" >&2
        exit 1
    fi

    # Validar JSON
    if command -v jq &> /dev/null; then
        if jq empty "$SWAGGER_DIR/swagger.json" 2>/dev/null; then
            echo -e "${GREEN}✓ swagger.json es válido${NC}"
        else
            echo -e "${RED}✗ swagger.json tiene errores de sintaxis${NC}"
            exit 1
        fi
    else
        echo -e "${YELLOW}⚠ jq no está instalado. Saltando validación JSON${NC}"
    fi

    # Verificar que docs.go existe
    if [[ -f "$SWAGGER_DIR/docs.go" ]]; then
        echo -e "${GREEN}✓ docs.go existe${NC}"
    else
        echo -e "${RED}✗ docs.go no existe${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Documentación válida${NC}"
    return 0
}

# Función para formatear anotaciones
format_annotations() {
    echo -e "${BLUE}✨ Formateando anotaciones Swagger...${NC}"

    if command -v swag &> /dev/null; then
        swag fmt
        echo -e "${GREEN}✓ Anotaciones formateadas${NC}"
    else
        echo -e "${YELLOW}⚠ swag CLI no encontrado. Instálalo con './swagger.sh install'${NC}"
        exit 1
    fi
    return 0
}

# Función para limpiar archivos generados
clean_docs() {
    echo -e "${BLUE}🧹 Limpiando archivos generados...${NC}"

    read -p "¿Estás seguro de eliminar docs.go, swagger.json y swagger.yaml? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -f "$SWAGGER_DIR/docs.go"
        rm -f "$SWAGGER_DIR/swagger.json"
        # No eliminamos swagger.yaml porque contiene la especificación manual
        echo -e "${GREEN}✓ Archivos generados eliminados${NC}"
    else
        echo -e "${YELLOW}Operación cancelada${NC}"
    fi
    return 0
}

# Función para instalar swag CLI
install_swag() {
    echo -e "${BLUE}📦 Instalando swag CLI...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✓ swag CLI instalado${NC}"
    echo -e "${BLUE}Versión:${NC}"
    swag --version
    return 0
}

# Función para construir imagen Docker
build_docker() {
    echo -e "${BLUE}🐳 Construyendo imagen Docker del generador...${NC}"
    docker build -f "$SWAGGER_DIR/Dockerfile.swag" -t flighthours-swag .
    echo -e "${GREEN}✓ Imagen Docker construida: flighthours-swag${NC}"
    return 0
}

# Función para abrir en navegador
open_browser() {
    echo -e "${BLUE}🌐 Abriendo Swagger UI en el navegador...${NC}"

    # Detectar sistema operativo
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open http://localhost:8082
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open http://localhost:8082
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        start http://localhost:8082
    else
        echo -e "${YELLOW}⚠ No se pudo detectar el sistema operativo${NC}"
        echo -e "${BLUE}Abre manualmente: http://localhost:8082${NC}"
    fi
    return 0
}

# Procesar comandos
case "${1:-help}" in
    generate|gen|g)
        generate_docs
        ;;
    ui|start)
        start_ui
        ;;
    ui-stop|stop)
        stop_ui
        ;;
    serve|run)
        serve_api
        ;;
    validate|check)
        validate_docs
        ;;
    format|fmt)
        format_annotations
        ;;
    clean)
        clean_docs
        ;;
    install)
        install_swag
        ;;
    docker-build|build)
        build_docker
        ;;
    open|browse)
        open_browser
        ;;
    help|--help|-h|*)
        show_help
        ;;
esac
