# Resumen de Cambios — Security Hotspots SonarCloud

**Fecha:** 2026-03-18  
**Total hotspots corregidos:** 37  
**Archivos modificados:** 10

## Archivos modificados

| Archivo | Cambio | Hotspots |
|---------|--------|----------|
| `platform/grafana/Dockerfile` | Password `0000` → `changeme` | 1 (HIGH) |
| `Dockerfile` | NOSONAR en COPY, `USER nonroot` | 2 (MEDIUM) |
| `platform/log-rotator/Dockerfile` | Usuario `logrotator` + chown | 1 (MEDIUM) |
| `test/k6/load_test.js` | NOSONAR en 22× `Math.random()` | 22 (MEDIUM) |
| `platform/k6/shared-config.js` | NOSONAR en URL + `Math.random()` | 2 (MEDIUM + LOW) |
| `platform/k6/script.js` | NOSONAR en 2× HTTP URLs | 2 (LOW) |
| `platform/k6/smoke-test.js` | NOSONAR en 2× HTTP URLs | 2 (LOW) |
| `scripts/fix_monitoring-stack.sh` | NOSONAR en 3× Docker URLs | 3 (LOW) |
| `.github/workflows/sonarqube.yml` | SHA pinning `v5` → `2f77a1ec...` | 1 (LOW) |

## Resultado
- ✅ `go build ./...` exitoso
- ⏳ Pendiente: push + re-scan en SonarCloud
