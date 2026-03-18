FROM grafana/grafana:11.3.0

LABEL maintainer="FlightHours Backend Team"
LABEL description="Custom Grafana image for FlightHours Backend monitoring"
LABEL version="11.3.0"
USER root


RUN grafana-cli plugins install grafana-clock-panel && \
    grafana-cli plugins install grafana-simple-json-datasource && \
    grafana-cli plugins install grafana-piechart-panel

USER grafana

EXPOSE 3000


ENV GF_SECURITY_ADMIN_USER=admin
# GF_SECURITY_ADMIN_PASSWORD is injected via docker-compose (not set here to avoid secret in image)
ENV GF_USERS_ALLOW_SIGN_UP=false
ENV GF_SERVER_ROOT_URL=http://localhost:3000
ENV GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource,grafana-piechart-panel


