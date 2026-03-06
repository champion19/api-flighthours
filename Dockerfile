# ============================================
# FlightHours Backend API — Multi-Stage Dockerfile
# ============================================
# Stage 1: Build (golang:1.25-alpine ~300MB)
# Stage 2: Runtime (distroless ~5MB)
# Final image: ~10-15MB
# ============================================

# ── Stage 1: Build ───────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache de dependencias (solo se re-descarga si go.mod/go.sum cambian)
COPY go.mod go.sum ./
RUN go mod download

# Compilar binario estático
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -o /flighthours-api ./cmd

# ── Stage 2: Runtime (imagen mínima) ─────────
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copiar binario compilado
COPY --from=builder /flighthours-api /app/flighthours-api

# Copiar go.mod (requerido por FindModuleRoot() en runtime)
COPY --from=builder /app/go.mod /app/go.mod

# Copiar archivos necesarios en runtime
COPY --from=builder /app/config /app/config
COPY --from=builder /app/platform/schema/json_schema /app/platform/schema/json_schema
COPY --from=builder /app/platform/swaggo /app/platform/swaggo

EXPOSE 8081

ENTRYPOINT ["/app/flighthours-api"]
