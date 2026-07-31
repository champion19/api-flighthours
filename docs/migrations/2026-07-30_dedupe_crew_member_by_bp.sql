-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: SearchCrewMembers y FindOrCreateCrewMember solo comparaban por
-- `name` — el bp (badge/carné) se guardaba pero nunca se usaba para
-- identificar a la persona. Si el piloto tipeaba el nombre distinto entre
-- vuelos (typo, nombre corto, etc.) con el mismo bp, el sistema creaba un
-- crew_member duplicado en vez de reconocer a la misma persona. Caso real
-- detectado: bp 4645117 tiene dos registros ("Laura Peña" y "Laura Marcela
-- P"), confirmado por el cliente como la misma persona.
--
-- Esta migración: 1) repointea daily_logbook_detail_crew de los duplicados
-- hacia el registro más antiguo por (employee_id, bp) [keeper], 2) borra los
-- duplicados, 3) agrega UNIQUE (employee_id, bp) para que esto no vuelva a
-- pasar a nivel de datos (el código ya se corrigió para buscar por bp antes
-- que por name).

-- Paso 1: mapa duplicado -> keeper (el creado primero por employee_id+bp)
CREATE TEMPORARY TABLE crew_member_keepers AS
SELECT employee_id, bp, MIN(created_at) AS keeper_created_at
FROM crew_member
WHERE bp IS NOT NULL AND bp <> ''
GROUP BY employee_id, bp
HAVING COUNT(*) > 1;

CREATE TEMPORARY TABLE crew_member_keeper_ids AS
SELECT cm.employee_id, cm.bp, cm.id AS keeper_id
FROM crew_member cm
JOIN crew_member_keepers k
  ON k.employee_id = cm.employee_id
 AND k.bp = cm.bp
 AND k.keeper_created_at = cm.created_at;

-- Verificación previa: registros que se van a fusionar (debe mostrar el caso
-- de Laura Peña / Laura Marcela P, bp 4645117, si corres esto en la BD local
-- con los datos de prueba actuales).
SELECT cm.id, cm.name, cm.bp, cm.created_at, k.keeper_id
FROM crew_member cm
JOIN crew_member_keeper_ids k ON k.employee_id = cm.employee_id AND k.bp = cm.bp;

-- Paso 2: repointear asignaciones de vuelo de los duplicados hacia el keeper
UPDATE daily_logbook_detail_crew dldc
JOIN crew_member cm ON cm.id = dldc.crew_member_id
JOIN crew_member_keeper_ids k ON k.employee_id = cm.employee_id AND k.bp = cm.bp
SET dldc.crew_member_id = k.keeper_id
WHERE cm.id <> k.keeper_id;

-- Paso 3: borrar los crew_member duplicados ya huérfanos
DELETE cm FROM crew_member cm
JOIN crew_member_keeper_ids k ON k.employee_id = cm.employee_id AND k.bp = cm.bp
WHERE cm.id <> k.keeper_id;

DROP TEMPORARY TABLE crew_member_keeper_ids;
DROP TEMPORARY TABLE crew_member_keepers;

-- Paso 4: prevenir que vuelva a pasar a nivel de datos
ALTER TABLE crew_member
  ADD UNIQUE KEY uq_crew_member_employee_bp (employee_id, bp);

-- Verificación final: no debe haber más bp duplicados por piloto.
SELECT employee_id, bp, COUNT(*) AS total
FROM crew_member
WHERE bp IS NOT NULL AND bp <> ''
GROUP BY employee_id, bp
HAVING total > 1;
