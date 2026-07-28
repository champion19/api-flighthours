-- Migración manual: corre esto en DBeaver / mysql client.
--
-- Contexto: se elimina el concepto de "route"/"airline_route" — un vuelo ya
-- no requiere una ruta preconfigurada ni aprobada por un admin, se crea
-- eligiendo libremente aeropuerto de origen y destino. Este primer paso
-- agrega las columnas directas a daily_logbook_detail y rellena (backfill)
-- los vuelos existentes desde la cadena airline_route -> route. Es seguro
-- de correr en caliente: no rompe el código actual (las columnas quedan
-- nullable y airline_route_id se conserva por ahora).
--
-- La migración que borra airline_route_id, route y airline_route
-- (2026-07-28_drop_route_and_airline_route.sql) se corre en un paso
-- posterior, junto con el despliegue del código nuevo.

-- 1) Agrega las columnas nuevas (nullable por ahora).
ALTER TABLE daily_logbook_detail
  ADD COLUMN origin_airport_id varchar(36) NULL AFTER airline_route_id,
  ADD COLUMN destination_airport_id varchar(36) NULL AFTER origin_airport_id,
  ADD CONSTRAINT fk_detail_origin_airport FOREIGN KEY (origin_airport_id) REFERENCES airport(id),
  ADD CONSTRAINT fk_detail_destination_airport FOREIGN KEY (destination_airport_id) REFERENCES airport(id);

-- 2) Backfill: copia origen/destino desde la ruta física vinculada.
UPDATE daily_logbook_detail dld
JOIN airline_route alr ON dld.airline_route_id = alr.id
JOIN route r ON alr.route_id = r.id
SET dld.origin_airport_id = r.origin_airport_id,
    dld.destination_airport_id = r.destination_airport_id;

-- 3) Verificación: no debe quedar ningún vuelo sin origen/destino.
--    Si esto no da 0, investigar antes de desplegar el código nuevo o
--    correr la migración final de DROP.
SELECT COUNT(*) AS vuelos_sin_origen_destino
FROM daily_logbook_detail
WHERE origin_airport_id IS NULL OR destination_airport_id IS NULL;
