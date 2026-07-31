-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: el cliente pidió 3 cargos más para crew_role (además de
-- captain/first officer): Instructor, Line Check Captain y Safety Pilot.
-- Son roles que el piloto vuela EN EL VUELO — no tienen relación con
-- employee.role (admin/pilot, ligado a Keycloak, no se toca).
--
-- Además, se agrega crew_role a `daily_logbook` (la bitácora/página, padre)
-- para que se pregunte UNA VEZ al crear la bitácora ("New Logbook Entry") y
-- sirva de default para todos sus vuelos — mismo patrón ya usado con
-- tail_number_id (ver 2026-07-20_add_daily_logbook_tail_number.sql). Cada
-- vuelo sigue pudiendo cambiar su propio crew_role sin afectar el default
-- de la bitácora.

ALTER TABLE daily_logbook_detail
  MODIFY COLUMN crew_role ENUM('captain', 'first officer', 'instructor', 'line check captain', 'safety pilot') DEFAULT NULL;

ALTER TABLE daily_logbook
  ADD COLUMN crew_role ENUM('captain', 'first officer', 'instructor', 'line check captain', 'safety pilot') DEFAULT NULL AFTER tail_number_id;

-- Verificación: la columna nueva debe existir en daily_logbook (NULL en bitácoras existentes, no se hace backfill).
SELECT id, tail_number_id, crew_role FROM daily_logbook LIMIT 5;

-- Verificación: el ENUM de daily_logbook_detail debe aceptar los 5 valores.
SHOW COLUMNS FROM daily_logbook_detail LIKE 'crew_role';
