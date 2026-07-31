-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: se agrega tripulación estructurada (Primer Oficial + tripulación
-- de cabina) a cada vuelo de la bitácora. Hasta ahora solo existía
-- `daily_logbook_detail.companion_name` (texto libre, sin catálogo ni
-- reutilización). Estas dos tablas nuevas reemplazan ese flujo hacia
-- adelante; `companion_name` NO se toca ni se borra (queda para lectura
-- histórica de registros viejos, el backend deja de escribirlo).
--
-- `crew_member`: roster de personas por piloto (cada piloto solo ve/busca
-- su propio roster). El UNIQUE (employee_id, name) es lo que permite
-- "agregar si no existe" de forma atómica (INSERT ... ON DUPLICATE KEY).
--
-- `daily_logbook_detail_crew`: asignación de tripulación por vuelo (varias
-- filas por vuelo, sin tope — cubre Primer Oficial + N tripulantes de
-- cabina). El capitán NO se registra aquí: es el propio piloto autenticado,
-- ya cubierto por daily_logbook_detail.crew_role.

CREATE TABLE crew_member (
  id varchar(36) NOT NULL,
  employee_id varchar(36) NOT NULL,
  name varchar(150) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_crew_member_employee_name (employee_id, name),
  CONSTRAINT fk_crew_member_employee FOREIGN KEY (employee_id) REFERENCES employee (id)
);

CREATE TABLE daily_logbook_detail_crew (
  id varchar(36) NOT NULL,
  daily_logbook_detail_id varchar(36) NOT NULL,
  crew_member_id varchar(36) NOT NULL,
  role enum('first_officer', 'purser', 'flight_attendant') NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_dldc_detail (daily_logbook_detail_id),
  CONSTRAINT fk_dldc_detail FOREIGN KEY (daily_logbook_detail_id) REFERENCES daily_logbook_detail (id) ON DELETE CASCADE,
  CONSTRAINT fk_dldc_crew_member FOREIGN KEY (crew_member_id) REFERENCES crew_member (id)
);

-- Verificación: ambas tablas deben existir y estar vacías recién creadas.
SELECT COUNT(*) AS total_crew_member FROM crew_member;
SELECT COUNT(*) AS total_daily_logbook_detail_crew FROM daily_logbook_detail_crew;

-- ============================================================
-- Crew Member Module (TRI_*) - System Messages
-- Mismo patrón que documents/sql/add_crew_member_type_messages.sql
-- y add_tail_number_messages.sql (type EXITO/ERROR, no http_status
-- column — el status HTTP vive en platform/cache/messaging/cache.go).
-- ============================================================
INSERT INTO system_messages (id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at)
VALUES
(
    UUID(),
    'TRI_AGR_EXI_05101',
    'EXITO',
    'final_user',
    'crew_member',
    'Crew Member Registered',
    'The crew member has been added to your roster (or already existed).',
    true,
    NOW(),
    NOW()
),
(
    UUID(),
    'TRI_CON_EXI_05201',
    'EXITO',
    'final_user',
    'crew_member',
    'Crew Member Search Successful',
    'The crew member search was completed successfully.',
    true,
    NOW(),
    NOW()
);

-- Verificación: deben existir exactamente estos 2 códigos nuevos.
SELECT message_code, type, module FROM system_messages WHERE message_code IN ('TRI_AGR_EXI_05101', 'TRI_CON_EXI_05201');
