-- Migración manual: corre esto en DBeaver / mysql client (BD local primero).
--
-- Contexto: se agrega el número BP (código interno/badge del tripulante,
-- mismo concepto y formato que airline_employee.bp — numérico, opcional)
-- al roster de crew_member, para poder identificar a la persona más allá
-- del nombre.

ALTER TABLE crew_member
  ADD COLUMN bp varchar(16) NULL AFTER name;

-- Verificación: la columna debe existir y estar vacía en los registros actuales.
SELECT id, employee_id, name, bp FROM crew_member;
