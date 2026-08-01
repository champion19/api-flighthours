-- Migración manual: agrega created_at (timestamp real de creación) a `daily_logbook`.
--
-- Motivo: el frontend necesita saber cuál es la bitácora creada MÁS RECIENTEMENTE
-- para defaultear el campo Crew Role en "New Logbook Entry". log_date y book_page
-- no sirven para esto — ambos son datos que el piloto edita libremente (log_date
-- puede ser una fecha pasada al cargar un vuelo atrasado; book_page es un código de
-- página física del logbook, no un contador secuencial de la app), así que ninguno
-- garantiza reflejar el orden real de creación. created_at lo pone la base de datos
-- sola al insertar, sin que el piloto lo toque, así que es el único dato confiable.
--
-- No hay herramienta de migraciones en este repo (solo dumps estáticos en documents/sql/),
-- así que este script se corre a mano contra la DB de producción.
-- Después de correrlo, documents/sql/current_schema*.sql y tables.sql deberían reflejar
-- el nuevo esquema (ya estaban desactualizados respecto a crew_role antes de este cambio).

ALTER TABLE daily_logbook
  ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER crew_role;

-- Verificación
SELECT id, log_date, book_page, crew_role, created_at FROM daily_logbook ORDER BY created_at DESC LIMIT 20;
