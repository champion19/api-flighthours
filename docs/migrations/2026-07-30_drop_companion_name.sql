-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: `companion_name` (texto libre) fue reemplazado por la
-- tripulación estructurada (crew_member + daily_logbook_detail_crew,
-- ver 2026-07-30_add_crew_member_and_assignments.sql). El backend dejó de
-- escribir este campo desde entonces y quedó siempre en NULL en todos los
-- vuelos (confirmado: 0 de 5 vuelos locales tenían valor). El cliente
-- confirmó que la columna ya no se necesita.

ALTER TABLE daily_logbook_detail
  DROP COLUMN companion_name;

-- Verificación: la columna ya no debe existir.
SHOW COLUMNS FROM daily_logbook_detail LIKE 'companion_name';
