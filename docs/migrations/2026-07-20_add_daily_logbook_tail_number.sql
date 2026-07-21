-- Migración manual: agrega tail_number_id (opcional) a `daily_logbook`.
-- Permite capturar la matrícula una sola vez al crear la bitácora ("New Logbook Entry")
-- y reusarla para todos los vuelos que se agreguen bajo ese book page.
-- No hay herramienta de migraciones en este repo (solo dumps estáticos en documents/sql/),
-- así que este script se corre a mano contra la DB de producción.
-- Después de correrlo, documents/sql/current_schema*.sql y tables.sql ya reflejan el nuevo esquema.

ALTER TABLE daily_logbook ADD COLUMN tail_number_id VARCHAR(36) NULL AFTER status;
ALTER TABLE daily_logbook ADD KEY fk_logbook_tail_number (tail_number_id);
ALTER TABLE daily_logbook ADD CONSTRAINT fk_logbook_tail_number FOREIGN KEY (tail_number_id) REFERENCES tail_number(id);

-- Verificación
SELECT id, log_date, book_page, tail_number_id FROM daily_logbook ORDER BY log_date DESC LIMIT 20;
