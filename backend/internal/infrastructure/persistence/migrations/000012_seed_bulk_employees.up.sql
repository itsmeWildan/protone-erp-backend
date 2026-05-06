-- 000012_seed_bulk_employees.up.sql

-- 1. Tambah Departemen Baru
INSERT INTO departments (id, tenant_id, name, code) VALUES 
('550e8400-e29b-41d4-a716-446655440010', (SELECT id FROM tenants LIMIT 1), 'Human Resources', 'HR'),
('550e8400-e29b-41d4-a716-446655440020', (SELECT id FROM tenants LIMIT 1), 'Finance & Accounting', 'FIN'),
('550e8400-e29b-41d4-a716-446655440030', (SELECT id FROM tenants LIMIT 1), 'Marketing', 'MKT');

-- 2. Tambah Posisi Baru
INSERT INTO positions (id, tenant_id, department_id, name, code) VALUES 
('550e8400-e29b-41d4-a716-446655440011', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440010', 'HR Specialist', 'HRS'),
('550e8400-e29b-41d4-a716-446655440012', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440010', 'HR Manager', 'HRM'),
('550e8400-e29b-41d4-a716-446655440021', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440020', 'Accountant', 'ACC'),
('550e8400-e29b-41d4-a716-446655440022', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440020', 'Finance Controller', 'FC'),
('550e8400-e29b-41d4-a716-446655440031', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440030', 'Social Media Specialist', 'SMS'),
('550e8400-e29b-41d4-a716-446655440032', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440030', 'Marketing Manager', 'MKM');

-- 3. Tambah Karyawan Baru (Minimal 2 per departemen)
INSERT INTO employees (id, tenant_id, nik, full_name, email, department_id, position_id, employment_type, status, join_date) VALUES 
-- HR
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'HR001', 'Siti Aminah', 'siti@acme.com', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440011', 'permanent', 'active', '2023-01-15'),
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'HR002', 'Andi Wijaya', 'andi@acme.com', '550e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440012', 'permanent', 'active', '2022-05-20'),
-- Finance
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'FIN001', 'Rina Pratama', 'rina@acme.com', '550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440021', 'contract', 'active', '2024-02-01'),
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'FIN002', 'Dewi Lestari', 'dewi@acme.com', '550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440022', 'permanent', 'active', '2021-11-10'),
-- Marketing
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'MKT001', 'Reza Rahardian', 'reza@acme.com', '550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440031', 'contract', 'active', '2024-03-15'),
(gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'MKT002', 'Maya Septha', 'maya@acme.com', '550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440032', 'permanent', 'active', '2022-09-01');
