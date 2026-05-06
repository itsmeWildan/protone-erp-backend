-- 000005_create_leave_management.up.sql

CREATE TYPE leave_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled');

CREATE TABLE leave_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,           -- Cuti Tahunan, Cuti Sakit, dll
    max_days        SMALLINT NOT NULL DEFAULT 0,     -- 0 = unlimited
    is_paid         BOOLEAN NOT NULL DEFAULT TRUE,
    requires_doc    BOOLEAN NOT NULL DEFAULT FALSE,  -- perlu dokumen pendukung
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

-- Saldo cuti per karyawan per tahun
CREATE TABLE leave_balances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id   UUID NOT NULL REFERENCES leave_types(id),
    year            SMALLINT NOT NULL,
    total_days      SMALLINT NOT NULL DEFAULT 0,    -- jatah
    used_days       SMALLINT NOT NULL DEFAULT 0,    -- sudah terpakai
    pending_days    SMALLINT NOT NULL DEFAULT 0,    -- sedang diajukan
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(employee_id, leave_type_id, year)
);

-- Pengajuan cuti
CREATE TABLE leaves (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id),
    leave_type_id   UUID NOT NULL REFERENCES leave_types(id),
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    total_days      SMALLINT NOT NULL,
    reason          TEXT NOT NULL,
    status          leave_status NOT NULL DEFAULT 'pending',
    approved_by     UUID REFERENCES employees(id),
    approved_at     TIMESTAMPTZ,
    rejection_note  TEXT,
    attachment_url  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_leave_types_tenant ON leave_types(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_leave_balances_employee ON leave_balances(employee_id);
CREATE INDEX idx_leaves_employee ON leaves(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_leaves_status ON leaves(tenant_id, status) WHERE deleted_at IS NULL;
