#!/bin/bash

# Entrypoint para el contenedor log-rotator
# Ejecuta logrotate cada hora y limpieza de volúmenes diariamente

set -e

echo "=== Log Rotator Container Started ==="
echo "Timezone: $(date)"
echo "Retention: 3 días"
echo ""

# Crear archivo de log si no existe
touch /var/log/logrotate.log
touch /var/log/volume-cleanup.log

# Verificar que los directorios estén montados
echo "Verificando volúmenes montados:"
[ -d "/var/log/flighthours-backend" ] && echo "  ✓ Logs del backend: /var/log/flighthours-backend" || echo "  ✗ Logs no montados"
[ -d "/prometheus" ] && echo "  ✓ Prometheus data: /prometheus" || echo "  ✗ Prometheus no montado"
[ -d "/loki" ] && echo "  ✓ Loki data: /loki" || echo "  ✗ Loki no montado"
[ -d "/grafana" ] && echo "  ✓ Grafana data: /grafana" || echo "  ✗ Grafana no montado"
echo ""

# Función para ejecutar logrotate
run_logrotate() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Ejecutando logrotate..."
    /usr/sbin/logrotate -v /etc/logrotate.conf 2>&1 | tee -a /var/log/logrotate.log
}

# Función para ejecutar limpieza de volúmenes
run_volume_cleanup() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Ejecutando limpieza de volúmenes..."
    /usr/local/bin/cleanup-volumes.sh
}

# Ejecutar logrotate inmediatamente al inicio (solo para verificar configuración)
echo "Ejecutando logrotate inicial (verificación)..."
run_logrotate

# Ejecutar limpieza de volúmenes al inicio
echo "Ejecutando limpieza de volúmenes inicial..."
run_volume_cleanup

echo ""
echo "=== Configurando tareas programadas ==="
echo "  • Logrotate: cada hora"
echo "  • Limpieza de volúmenes: diariamente a las 02:00 AM"
echo ""

# Usar un loop infinito en lugar de cron para mejor compatibilidad con Docker
echo "Iniciando loops de monitoreo..."
echo ""

# Loop en background para logrotate (cada hora)
(
  while true; do
    sleep 3600  # 1 hora
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Ejecutando logrotate programado..." >> /var/log/logrotate.log
    /usr/sbin/logrotate /etc/logrotate.conf >> /var/log/logrotate.log 2>&1
  done
) &

# Loop en background para limpieza de volúmenes (cada 24 horas a las 2 AM)
(
  while true; do
    # Calcular segundos hasta las 2 AM del siguiente día
    current_hour=$(date +%H)
    current_min=$(date +%M)

    if [ "$current_hour" -eq 2 ] && [ "$current_min" -lt 10 ]; then
      # Estamos entre 2:00 y 2:10 AM, ejecutar limpieza
      /usr/local/bin/cleanup-volumes.sh
      sleep 600  # Esperar 10 minutos para evitar ejecuciones múltiples
    fi

    sleep 600  # Verificar cada 10 minutos
  done
) &

# Loop en background para heartbeat (cada 6 horas)
(
  while true; do
    sleep 21600  # 6 horas
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Log rotator is alive" >> /var/log/logrotate.log
  done
) &

echo "✅ Todas las tareas programadas iniciadas"
echo "   → Logrotate: cada hora"
echo "   → Limpieza de volúmenes: diariamente a las 02:00 AM"
echo "   → Heartbeat: cada 6 horas"
echo ""
echo "Contenedor en ejecución. Presiona Ctrl+C para detener."
echo ""

# Mantener el contenedor vivo
tail -f /var/log/logrotate.log /var/log/volume-cleanup.log
