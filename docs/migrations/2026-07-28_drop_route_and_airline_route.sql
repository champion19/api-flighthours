-- Migración manual: corre esto en DBeaver / mysql client.
--
-- Contexto: paso final de la eliminación de "route"/"airline_route".
-- SOLO correr después de:
--   1) Haber aplicado 2026-07-28_add_direct_airports_to_daily_logbook_detail.sql
--      y confirmado que vuelos_sin_origen_destino = 0.
--   2) Haber desplegado el código nuevo del backend (que ya no usa
--      airline_route_id) y del frontend (que ya no lo envía).
--   3) Haber verificado en producción que crear/editar/listar vuelos
--      funciona correctamente con origin_airport_id/destination_airport_id.
--
-- Es IRREVERSIBLE (DROP TABLE). Ya existe backup de la DB en Dropbox/nube,
-- pero de todas formas confirmar antes de ejecutar.

-- 0) Verificación previa (debe dar 0 antes de continuar).
SELECT COUNT(*) AS vuelos_sin_origen_destino
FROM daily_logbook_detail
WHERE origin_airport_id IS NULL OR destination_airport_id IS NULL;

-- 1) Vuelve NOT NULL las columnas nuevas y quita el vínculo viejo.
ALTER TABLE daily_logbook_detail
  MODIFY origin_airport_id varchar(36) NOT NULL,
  MODIFY destination_airport_id varchar(36) NOT NULL,
  DROP FOREIGN KEY fk_detail_airline_route,
  DROP COLUMN airline_route_id;

-- 2) Elimina las tablas de route/airline_route.
DROP TABLE airline_route;
DROP TABLE route;
