-- Migración manual: agrega oaci_code a `airport`, para buscar por código OACI
-- igual que ya funciona por iata_code.
-- Corre esto en DBeaver contra tu DB local.

ALTER TABLE airport ADD COLUMN oaci_code VARCHAR(4) NULL AFTER iata_code;

-- Opcional: si quieres rellenar los OACI de los aeropuertos que ya insertaste,
-- hazlo con UPDATE por iata_code, por ejemplo:
-- UPDATE airport SET oaci_code = 'SKBO' WHERE iata_code = 'BOG';
-- UPDATE airport SET oaci_code = 'SKBQ' WHERE iata_code = 'BAQ';

-- Verificación
SELECT iata_code, oaci_code, name FROM airport ORDER BY name;
