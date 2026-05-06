-- 000004_create_employees.up.sql

CREATE TYPE employee_status AS ENUM ('active', 'inactive', 'terminated', 'on_leave');
CREATE TYPE employee_gender AS ENUM ('male', 'female');
CREATE TYPE marital_status AS ENUM ('single', 'married', 'divorced', 'widowed');
CREATE TYPE employment_type AS ENUM ('permanent', 'contract', 'internship', 'freelance');

CREATE TABLE employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id             UUID REFERENCES users(id),         -- link ke akun login (optional)
    nik                 VARCHAR(50) NOT NULL,              -- Nomor Induk Karyawan
    full_name           VARCHAR(255) NOT NULL,
    email               VARCHAR(255),
    phone               VARCHAR(20),
    gender              employee_gender,
    birth_date          DATE,
    birth_place         VARCHAR(100),
    marital_status      marital_status,
    national_id         VARCHAR(20),                       -- KTP
    tax_id              VARCHAR(20),                       -- NPWP
    address             TEXT,
    department_id       UUID NOT NULL REFERENCES departments(id),
    position_id         UUID NOT NULL REFERENCES positions(id),
    manager_id          UUID REFERENCES employees(id),    -- atasan langsung
    employment_type     employment_type NOT NULL DEFAULT 'permanent',
    status              employee_status NOT NULL DEFAULT 'active',
    join_date           DATE NOT NULL,
    end_date            DATE,                              -- untuk kontrak
    basic_salary        NUMERIC(15,2) NOT NULL DEFAULT 0,
    bank_name           VARCHAR(100),
    bank_account_no     VARCHAR(50),
    bank_account_name   VARCHAR(255),
    bpjs_ketenagakerjaan VARCHAR(20),
    bpjs_kesehatan      VARCHAR(20),
    photo_url           TEXT,
    created_by          UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    UNIQUE(tenant_id, nik)
);

CREATE INDEX idx_employees_tenant ON employees(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_dept ON employees(department_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_position ON employees(position_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_manager ON employees(manager_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_status ON employees(tenant_id, status) WHERE deleted_at IS NULL;
