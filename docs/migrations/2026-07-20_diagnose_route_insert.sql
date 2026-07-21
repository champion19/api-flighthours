-- Diagnóstico: por qué el INSERT de la ruta BOG->BAQ no funcionó.
-- Corre cada SELECT por separado y revisa qué sale vacío o inesperado.

-- 1. ¿Existen los aeropuertos con esos códigos IATA exactos?
SELECT id, iata_code, name FROM airport WHERE iata_code IN ('BOG', 'BAQ');

-- 2. ¿Hay alguna aerolínea creada en la DB local?
SELECT id, airline_name, airline_code, status FROM airline;

-- 3. ¿Ya existe una ruta BOG->BAQ (el UNIQUE origin+destination pudo rechazar el insert)?
SELECT r.id, a_o.iata_code AS origin, a_d.iata_code AS destination
FROM route r
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id
WHERE a_o.iata_code = 'BOG' AND a_d.iata_code = 'BAQ';

-- 4. Todas las rutas que existan, para tener contexto
SELECT r.id, a_o.iata_code AS origin, a_d.iata_code AS destination, r.airport_type
FROM route r
JOIN airport a_o ON a_o.id = r.origin_airport_id
JOIN airport a_d ON a_d.id = r.destination_airport_id;
