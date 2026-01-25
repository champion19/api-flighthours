FROM alpine:3.19

LABEL maintainer="Flighthours Team"
LABEL description="Log rotation and volume cleanup service for Flighthours Backend"
LABEL version="1.0"

RUN apk add --no-cache \
    logrotate \
    bash \
    tzdata \
    coreutils \
    findutils \
    && rm -rf /var/cache/apk/*

ENV TZ=America/Bogota
ENV RETENTION_DAYS=7

RUN mkdir -p /var/lib/logrotate \
    && mkdir -p /var/log/flighthours-backend \
    && mkdir -p /etc/logrotate.d

COPY logrotate.conf /etc/logrotate.conf
COPY logrotate.d/flighthours-backend /etc/logrotate.d/flighthours-backend
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY cleanup-volumes.sh /usr/local/bin/cleanup-volumes.sh

RUN chmod +x /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/cleanup-volumes.sh

RUN touch /var/log/logrotate.log \
    && touch /var/log/volume-cleanup.log

VOLUME ["/var/log/flighthours-backend"]

HEALTHCHECK --interval=1h --timeout=10s --start-period=30s --retries=3 \
    CMD [ -f /var/log/logrotate.log ] && [ -f /var/log/volume-cleanup.log ] || exit 1

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
