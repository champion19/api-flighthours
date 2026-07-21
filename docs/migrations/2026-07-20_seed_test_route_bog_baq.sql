-- Script de prueba manual (no es parte de una migración real): crea una
-- ruta BOG -> BAQ y la vincula a una aerolínea, para poder probar el nuevo
-- selector de origen/destino sin tener que montar el admin de rutas.
-- Ajusta el WHERE del @airline_id si tu aerolínea de prueba es otra.

SET @origin_id = (SELECT id FROM airport WHERE iata_code = 'BOG' LIMIT 1);
SET @dest_id   = (SELECT id FROM airport WHERE iata_code = 'BAQ' LIMIT 1);
SET @route_id  = UUID();
SET @airline_id = (SELECT id FROM airline LIMIT 1); -- <-- ajusta si hace falta

-- route.airport_type describe la ruta en sí (Nacional porque BOG y BAQ son
-- ambos colombianos), no cada aeropuerto por separado.
INSERT INTO route (id, origin_airport_id, destination_airport_id, airport_type, estimated_flight_time)
VALUES (@route_id, @origin_id, @dest_id, 'Nacional', '01:10:00');

INSERT INTO airline_route (id, route_id, airline_id, status)
VALUES (UUID(), @route_id, @airline_id, '1'); -- status usa '1'/'0', no 'active'/'inactive'

-- Verificación
SELECT ar.id, a_o.iata_code AS origin, a_d.iata_code AS destination, al.airline_name, ar.status
FROM airline_route ar
JOIN route r ON r.id = ar.route_id
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
JOIN airline al ON al.id = ar.airline_id
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';
