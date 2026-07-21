-- Corrige el status de la airline_route BOG->BAQ que insertamos con 'active'.
-- En este proyecto airline_route.status usa '1' (activo) / '0' (inactivo),
-- no las palabras 'active'/'inactive' (esas se usan en otras tablas, pero no en esta).
-- Confirmado contra la data semilla real del repo (docs/FlightHours/.../data-flighthours.sql).

UPDATE airline_route ar
JOIN route r ON r.id = ar.route_id
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
SET ar.status = '1'
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';

-- Verificación
SELECT ar.id, a_o.iata_code AS origin, a_d.iata_code AS destination, al.airline_name, ar.status
FROM airline_route ar
JOIN route r ON r.id = ar.route_id
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
JOIN airline al ON al.id = ar.airline_id
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';
