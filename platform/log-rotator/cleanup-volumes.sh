#!/bin/bash

# Script para limpiar datos antiguos de volúmenes Docker
# Ejecuta limpieza cada 7 días
# Autor: Motogo Team
# Fecha: 2025-12-13

set -e

LOG_FILE="/var/log/volume-cleanup.log"
RETENTION_DAYS=30

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "=== Iniciando limpieza de volúmenes Docker ==="

# Función para limpiar datos antiguos de Prometheus
cleanup_prometheus() {
    log "Limpiando datos de Prometheus (>$RETENTION_DAYS días)..."

    # Prometheus TSDB - eliminar bloques antiguos
    if [ -d "/prometheus" ]; then
        local deleted=0
        # Buscar directorios de bloques TSDB (formato: 01XXXXXXXXXX)
        find /prometheus -maxdepth 1 -type d -name "01*" -mtime +$RETENTION_DAYS | while read -r block; do
            log "  Eliminando bloque Prometheus: $(basename "$block")"
            rm -rf "$block"
            ((deleted++))
        done

        # Limpiar WAL antiguo
        if [ -d "/prometheus/wal" ]; then
            find /prometheus/wal -type f -mtime +$RETENTION_DAYS -delete
            log "  WAL limpiado"
        fi

        log "Prometheus: $deleted bloques eliminados"
    else
        log "Directorio /prometheus no encontrado (volumen no montado)"
    fi
}

# Función para limpiar datos antiguos de Loki
cleanup_loki() {
    log "Limpiando datos de Loki (>$RETENTION_DAYS días)..."

    if [ -d "/loki" ]; then
        local deleted=0

        # Limpiar chunks antiguos
        if [ -d "/loki/chunks" ]; then
            find /loki/chunks -type f -mtime +$RETENTION_DAYS | while read -r chunk; do
                log "  Eliminando chunk Loki: $(basename "$chunk")"
                rm -f "$chunk"
                ((deleted++))
            done
        fi

        # Limpiar índices antiguos
        if [ -d "/loki/index" ]; then
            find /loki/index -type f -mtime +$RETENTION_DAYS -delete
        fi

        # Limpiar WAL
        if [ -d "/loki/wal" ]; then
            find /loki/wal -type f -mtime +$RETENTION_DAYS -delete
        fi

        log "Loki: $deleted chunks eliminados"
    else
        log "Directorio /loki no encontrado (volumen no montado)"
    fi
}

# Función para limpiar datos antiguos de Grafana (principalmente dashboards temporales)
cleanup_grafana() {
    log "Limpiando datos temporales de Grafana (>$RETENTION_DAYS días)..."

    if [ -d "/grafana" ]; then
        # Limpiar sesiones antiguas
        if [ -d "/grafana/sessions" ]; then
            find /grafana/sessions -type f -mtime +$RETENTION_DAYS -delete
            log "Sesiones antiguas limpiadas"
        fi

        # Limpiar logs de plugins antiguos
        if [ -d "/grafana/plugins" ]; then
            find /grafana/plugins -name "*.log" -type f -mtime +$RETENTION_DAYS -delete
        fi

        log "Grafana: datos temporales limpiados"
    else
        log "Directorio /grafana no encontrado (volumen no montado)"
    fi
}

# Función para mostrar espacio liberado
show_disk_usage() {
    log "=== Uso de disco de volúmenes ==="

    if [ -d "/prometheus" ]; then
        local size=$(du -sh /prometheus 2>/dev/null | cut -f1)
        log "Prometheus: $size"
    fi

    if [ -d "/loki" ]; then
        local size=$(du -sh /loki 2>/dev/null | cut -f1)
        log "Loki: $size"
    fi

    if [ -d "/grafana" ]; then
        local size=$(du -sh /grafana 2>/dev/null | cut -f1)
        log "Grafana: $size"
    fi
}

# Mostrar uso antes de la limpieza
log "Uso de disco ANTES de la limpieza:"
show_disk_usage

# Ejecutar limpiezas
cleanup_prometheus
cleanup_loki
cleanup_grafana

# Mostrar uso después de la limpieza
log "Uso de disco DESPUÉS de la limpieza:"
show_disk_usage

log "=== Limpieza de volúmenes completada ==="
