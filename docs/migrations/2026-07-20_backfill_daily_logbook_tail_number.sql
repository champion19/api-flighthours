-- Backfill: rellena daily_logbook.tail_number_id para bitácoras creadas antes
-- de que el campo existiera, usando la matrícula de sus vuelos ya guardados
-- en daily_logbook_detail.actual_tail_number_id.
--
-- Solo rellena bitácoras donde TODOS sus vuelos comparten la misma matrícula
-- (caso inequívoco). Si una bitácora tiene vuelos con matrículas distintas,
-- se deja en NULL a propósito — no hay una única respuesta correcta ahí.

UPDATE daily_logbook dl
JOIN (
    SELECT daily_logbook_id, MIN(actual_tail_number_id) AS tail_number_id
    FROM daily_logbook_detail
    GROUP BY daily_logbook_id
    HAVING COUNT(DISTINCT actual_tail_number_id) = 1
) single_tail ON single_tail.daily_logbook_id = dl.id
SET dl.tail_number_id = single_tail.tail_number_id
WHERE dl.tail_number_id IS NULL;

-- Verificación: cuántas quedaron rellenas
SELECT COUNT(*) AS backfilled FROM daily_logbook WHERE tail_number_id IS NOT NULL;

-- Bitácoras que NO se pudieron rellenar automáticamente (0 vuelos, o vuelos
-- con matrículas distintas entre sí) — revísalas a mano si hace falta.
SELECT
    dl.id,
    dl.log_date,
    dl.book_page,
    COUNT(dld.id) AS flight_count,
    COUNT(DISTINCT dld.actual_tail_number_id) AS distinct_tail_numbers
FROM daily_logbook dl
LEFT JOIN daily_logbook_detail dld ON dld.daily_logbook_id = dl.id
WHERE dl.tail_number_id IS NULL
GROUP BY dl.id, dl.log_date, dl.book_page
ORDER BY dl.log_date DESC;
