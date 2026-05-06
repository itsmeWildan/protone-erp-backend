-- 000011_seed_leave_data.up.sql

-- 1. Buat Jenis Cuti
INSERT INTO leave_types (id, tenant_id, code, name, max_days, is_paid)
VALUES ('550e8400-e29b-41d4-a716-446655440003', (SELECT id FROM tenants LIMIT 1), 'AL', 'Annual Leave', 12, TRUE);

-- 2. Buat Saldo Cuti untuk Karyawan (Budi Santoso)
INSERT INTO leave_balances (id, tenant_id, employee_id, leave_type_id, year, total_days, used_days, pending_days)
VALUES (
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1),
    (SELECT id FROM employees LIMIT 1),
    '550e8400-e29b-41d4-a716-446655440003',
    2024,
    12,
    0,
    0
);
