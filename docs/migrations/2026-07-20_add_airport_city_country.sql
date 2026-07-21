-- Migración manual: agrega city/country a `airport`.
-- No hay herramienta de migraciones en este repo (solo dumps estáticos en documents/sql/),
-- así que este script se corre a mano contra la DB de producción.
-- Después de correrlo, documents/sql/current_schema*.sql y tables.sql ya reflejan el nuevo esquema.

ALTER TABLE airport ADD COLUMN city VARCHAR(50) NULL AFTER name;
ALTER TABLE airport ADD COLUMN country VARCHAR(50) NULL AFTER city;

UPDATE airport SET city = 'Cali', country = 'Colombia' WHERE iata_code = 'CLO';
UPDATE airport SET city = 'Miami', country = 'Estados Unidos' WHERE iata_code = 'MIA';
UPDATE airport SET city = 'Pereira', country = 'Colombia' WHERE iata_code = 'PEI';
UPDATE airport SET city = 'Bogotá', country = 'Colombia' WHERE iata_code = 'BOG';
UPDATE airport SET city = 'Santa Marta', country = 'Colombia' WHERE iata_code = 'SMR';
UPDATE airport SET city = 'Pasto', country = 'Colombia' WHERE iata_code = 'PSO';
UPDATE airport SET city = 'Montería', country = 'Colombia' WHERE iata_code = 'MTR';
UPDATE airport SET city = 'Bucaramanga', country = 'Colombia' WHERE iata_code = 'BGA';
UPDATE airport SET city = 'Cúcuta', country = 'Colombia' WHERE iata_code = 'CUC';
UPDATE airport SET city = 'Yopal', country = 'Colombia' WHERE iata_code = 'EYP';
UPDATE airport SET city = 'Barranquilla', country = 'Colombia' WHERE iata_code = 'BAQ';
UPDATE airport SET city = 'Armenia', country = 'Colombia' WHERE iata_code = 'AXM';
UPDATE airport SET city = 'San Andrés', country = 'Colombia' WHERE iata_code = 'ADZ';
UPDATE airport SET city = 'Rionegro', country = 'Colombia' WHERE iata_code = 'MDE';
UPDATE airport SET city = 'Cartagena', country = 'Colombia' WHERE iata_code = 'CTG';

-- Verificación
SELECT iata_code, name, city, country, airport_type FROM airport ORDER BY name;
