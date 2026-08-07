-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: rediseño de "Tripulación" en Daily Logbook Detail. El slot fijo
-- de Primer Oficial (0-1, role='first_officer') se reemplaza por una lista
-- repetible "Tripulación de Mando" (0-N, realmente 2-4) donde cada persona
-- lleva uno de los 5 roles de daily_logbook_detail.crew_role (captain,
-- first officer, instructor, line check captain, safety pilot) — así se
-- registra con qué rol voló cada piloto adicional, no solo que era "el
-- primer oficial". Tripulación de cabina (purser/flight_attendant) no
-- cambia.
--
-- `daily_logbook_detail_crew.role` es un ENUM que hoy solo acepta
-- ('first_officer', 'purser', 'flight_attendant') — hay que ampliarlo para
-- aceptar también los 5 valores de crew_role. Se deja 'first_officer' en el
-- enum (no se borra) porque filas históricas ya lo usan; el backend/frontend
-- simplemente dejan de escribirlo hacia adelante.

ALTER TABLE daily_logbook_detail_crew
  MODIFY COLUMN role enum(
    'first_officer',
    'purser',
    'flight_attendant',
    'captain',
    'first officer',
    'instructor',
    'line check captain',
    'safety pilot'
  ) NOT NULL;

-- Verificación: la definición de la columna debe incluir los 8 valores de arriba.
SHOW COLUMNS FROM daily_logbook_detail_crew LIKE 'role';
