-- Versión sin variables de sesión (@origin_id, etc.), por si DBeaver no las
-- mantenía entre statements. Todo en un solo INSERT con subconsultas.

-- 0. Confirma que route está realmente vacía
SELECT COUNT(*) AS total_routes FROM route;

-- 1. Crea la ruta BOG -> BAQ directamente, sin variables
INSERT INTO route (id, origin_airport_id, destination_airport_id, airport_type, estimated_flight_time)
SELECT
    UUID(),
    (SELECT id FROM airport WHERE iata_code = 'BOG' LIMIT 1),
    (SELECT id FROM airport WHERE iata_code = 'BAQ' LIMIT 1),
    'Nacional',
    '01:10:00';

-- 2. Verifica que quedó
SELECT r.id, a_o.iata_code AS origin, a_d.iata_code AS destination
FROM route r
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';

-- 3. Vincula esa ruta a tu aerolínea (ajusta el WHERE si tienes varias aerolíneas)
INSERT INTO airline_route (id, route_id, airline_id, status)
SELECT
    UUID(),
    (SELECT r.id FROM route r
       JOIN airport a_o ON a_o.id = r.origin_airport_id
       JOIN airport a_d ON a_d.id = r.destination_airport_id
       WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ' LIMIT 1),
    (SELECT id FROM airline LIMIT 1), -- <-- ajusta si hace falta
    '1'; -- airline_route.status usa '1'/'0', no 'active'/'inactive'

-- 4. Verificación final
SELECT ar.id, a_o.iata_code AS origin, a_d.iata_code AS destination, al.airline_name, ar.status
FROM airline_route ar
JOIN route r ON r.id = ar.route_id
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
JOIN airline al ON al.id = ar.airline_id
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';
