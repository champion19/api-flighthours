-- Migración manual: corre esto en DBeaver contra tu DB local.
--
-- 1) Normaliza `airline_route.status` (hoy guardado como el booleano de Go,
--    ej. '1'/'0' o 'True'/'False') hacia el texto real que ahora usa el
--    backend: 'active' / 'inactive' / 'pending'.
--    'pending' es nuevo — nunca existía antes de esta migración, así que no
--    hay filas que migrar hacia ese valor, solo aparecerá a futuro cuando
--    el sistema auto-cree un airline_route al resolver una ruta.
--
-- 2) Agrega un UNIQUE(route_id, airline_id) a `airline_route` para que no
--    puedan existir dos vínculos duplicados entre la misma ruta y la misma
--    aerolínea — necesario para que el auto-create sea seguro ante
--    solicitudes concurrentes.

-- Verificación previa: revisa qué valores tienes hoy antes de normalizar.
SELECT status, COUNT(*) FROM airline_route GROUP BY status;

-- Normalización (cubre '1', 'true', 'True', 'TRUE', 't', 'T' -> 'active';
-- '0', 'false', 'False', 'FALSE', 'f', 'F' -> 'inactive').
UPDATE airline_route
SET status = 'active'
WHERE LOWER(status) IN ('1', 'true', 't');

UPDATE airline_route
SET status = 'inactive'
WHERE LOWER(status) IN ('0', 'false', 'f');

-- Verificación: no debería quedar ninguna fila fuera de los 3 valores válidos.
SELECT status, COUNT(*) FROM airline_route
WHERE status NOT IN ('active', 'pending', 'inactive')
GROUP BY status;

-- Antes de agregar el UNIQUE, revisa si ya existen duplicados
-- (route_id, airline_id) — si esta consulta devuelve filas, bórralas a mano
-- primero (deja la que prefieras conservar) o el ALTER de abajo fallará.
SELECT route_id, airline_id, COUNT(*) AS total
FROM airline_route
GROUP BY route_id, airline_id
HAVING COUNT(*) > 1;

ALTER TABLE airline_route
  ADD CONSTRAINT uq_airline_route_route_airline UNIQUE (route_id, airline_id);

-- Mensajes del sistema para el nuevo flujo "resolver/crear ruta pendiente"
-- (RUT_AIR_RES_*), mismo patrón que el resto de mensajes de airline_route.
INSERT INTO system_messages (id, message_code, type, module, category, message_title, message_content, is_active, created_at, updated_at)
VALUES (
    UUID(),
    'RUT_AIR_RES_EXI_04401',
    'EXITO',
    'airline_route',
    'system',
    'Airline Route Link Found',
    'The airline already has a link for this route.',
    TRUE,
    now(),
    now()
);

INSERT INTO system_messages (id, message_code, type, module, category, message_title, message_content, is_active, created_at, updated_at)
VALUES (
    UUID(),
    'RUT_AIR_RES_EXI_04402',
    'EXITO',
    'airline_route',
    'system',
    'Airline Route Link Requested',
    'A new link between this airline and route was created and is pending admin approval.',
    TRUE,
    now(),
    now()
);

INSERT INTO system_messages (id, message_code, type, module, category, message_title, message_content, is_active, created_at, updated_at)
VALUES (
    UUID(),
    'RUT_AIR_RES_ERR_04403',
    'ERROR',
    'airline_route',
    'system',
    'Airline Route Resolve Error',
    'A technical error occurred while resolving or creating the airline route link.',
    TRUE,
    now(),
    now()
);

-- Verificación final
SELECT status, COUNT(*) FROM airline_route GROUP BY status;
SELECT message_code, message_title FROM system_messages WHERE message_code LIKE 'RUT_AIR_RES_%';
