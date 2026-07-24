-- Migración manual: corre esto en DBeaver / mysql client.
--
-- Contexto: tras editar manualmente `route` y `airline_route`, quedaron
-- 17 filas en `airline_route` apuntando a un `route_id` que ya no existe
-- (huérfanas, verificado que 0 vuelos en `daily_logbook_detail` las usan),
-- y 16 rutas en `route` sin ningún `airline_route` vinculado.
--
-- 1) Verificación previa: huérfanas en airline_route (route_id inexistente).
SELECT ar.* FROM airline_route ar
LEFT JOIN route r ON r.id = ar.route_id
WHERE r.id IS NULL;

-- 2) Verificación previa: rutas sin ningún airline_route.
SELECT r.id AS route_id, r.origin_airport_id, r.destination_airport_id
FROM route r
LEFT JOIN airline_route ar ON ar.route_id = r.id
WHERE ar.id IS NULL;

-- 3) Borra las huérfanas (confirmado: no las usa ningún vuelo real).
DELETE ar FROM airline_route ar
LEFT JOIN route r ON r.id = ar.route_id
WHERE r.id IS NULL;

-- 4) Crea el airline_route faltante para cada ruta sin vínculo,
--    todas asignadas a Latam (airline_id = 3f9d5e8c-7f0e-4c3a-aabc-91a1f5e986d6).
INSERT INTO airline_route (id, route_id, airline_id, status)
SELECT UUID(), r.id, '3f9d5e8c-7f0e-4c3a-aabc-91a1f5e986d6', 'active'
FROM route r
LEFT JOIN airline_route ar ON ar.route_id = r.id
WHERE ar.id IS NULL;

-- 5) Verificación final: no debe quedar ninguna ruta sin vínculo
--    ni ningún airline_route huérfano.
SELECT COUNT(*) AS rutas_sin_vinculo
FROM route r
LEFT JOIN airline_route ar ON ar.route_id = r.id
WHERE ar.id IS NULL;

SELECT COUNT(*) AS airline_route_huerfanos
FROM airline_route ar
LEFT JOIN route r ON r.id = ar.route_id
WHERE r.id IS NULL;

SELECT COUNT(*) AS total_route, (SELECT COUNT(*) FROM airline_route) AS total_airline_route
FROM route;
