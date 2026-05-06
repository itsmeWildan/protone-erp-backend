-- 000010_seed_initial_data.up.sql

-- Buat satu departemen dummy
INSERT INTO departments (id, tenant_id, name, code)
VALUES ('550e8400-e29b-41d4-a716-446655440001', (SELECT id FROM tenants LIMIT 1), 'Information Technology', 'IT');

-- Buat satu posisi dummy
INSERT INTO positions (id, tenant_id, department_id, name, code)
VALUES ('550e8400-e29b-41d4-a716-446655440002', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440001', 'Software Engineer', 'SWE');
