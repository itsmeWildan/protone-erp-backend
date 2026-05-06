-- 000014_seed_finance_data.up.sql

-- 1. Tambah Kelompok Akun
INSERT INTO account_groups (id, tenant_id, code, name, type) VALUES 
('550e8400-e29b-41d4-a716-446655440301', (SELECT id FROM tenants LIMIT 1), '10', 'Asset Lancar', 'asset'),
('550e8400-e29b-41d4-a716-446655440302', (SELECT id FROM tenants LIMIT 1), '20', 'Kewajiban Jangka Pendek', 'liability'),
('550e8400-e29b-41d4-a716-446655440303', (SELECT id FROM tenants LIMIT 1), '50', 'Bebas Operasional', 'expense');

-- 2. Tambah Chart of Accounts (COA)
INSERT INTO chart_of_accounts (id, tenant_id, account_group_id, code, name, type, normal_balance) VALUES 
-- Kas & Bank
('550e8400-e29b-41d4-a716-446655440401', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440301', '1-1001', 'Kas & Bank', 'asset', 'debit'),
-- Hutang Gaji
('550e8400-e29b-41d4-a716-446655440402', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440302', '2-1001', 'Hutang Gaji', 'liability', 'credit'),
-- Beban Gaji
('550e8400-e29b-41d4-a716-446655440403', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440303', '5-1001', 'Beban Gaji', 'expense', 'debit'),
-- Beban Reimbursement
('550e8400-e29b-41d4-a716-446655440404', (SELECT id FROM tenants LIMIT 1), '550e8400-e29b-41d4-a716-446655440303', '5-1002', 'Beban Reimbursement', 'expense', 'debit');
