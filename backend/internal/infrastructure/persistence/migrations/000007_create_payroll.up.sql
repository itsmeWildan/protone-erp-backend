-- 000007_create_payroll.up.sql

CREATE TYPE payroll_status AS ENUM ('draft', 'processing', 'approved', 'paid', 'cancelled');
CREATE TYPE reimbursement_status AS ENUM ('pending', 'approved', 'rejected', 'paid');

-- Komponen gaji: Tunjangan, Potongan, dll
CREATE TABLE salary_components (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,   -- Tunjangan Transport, BPJS, PPh21, dll
    type            VARCHAR(20) NOT NULL,    -- 'allowance' (tunjangan) | 'deduction' (potongan)
    is_taxable      BOOLEAN NOT NULL DEFAULT TRUE,
    is_fixed        BOOLEAN NOT NULL DEFAULT TRUE,   -- fixed amount atau dihitung
    default_amount  NUMERIC(15,2) NOT NULL DEFAULT 0,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

-- Periode penggajian
CREATE TABLE payroll_periods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    period_month    SMALLINT NOT NULL,    -- 1-12
    period_year     SMALLINT NOT NULL,
    status          payroll_status NOT NULL DEFAULT 'draft',
    total_amount    NUMERIC(15,2) NOT NULL DEFAULT 0,
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    paid_at         TIMESTAMPTZ,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, period_month, period_year)
);

-- Slip gaji per karyawan per periode
CREATE TABLE payroll_slips (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    payroll_period_id   UUID NOT NULL REFERENCES payroll_periods(id),
    employee_id         UUID NOT NULL REFERENCES employees(id),
    basic_salary        NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_allowance     NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_deduction     NUMERIC(15,2) NOT NULL DEFAULT 0,
    net_salary          NUMERIC(15,2) NOT NULL DEFAULT 0,    -- take home pay
    working_days        SMALLINT NOT NULL DEFAULT 0,
    present_days        SMALLINT NOT NULL DEFAULT 0,
    overtime_hours      NUMERIC(5,2) NOT NULL DEFAULT 0,
    overtime_amount     NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(payroll_period_id, employee_id)
);

-- Detail komponen per slip
CREATE TABLE payroll_slip_details (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_slip_id     UUID NOT NULL REFERENCES payroll_slips(id) ON DELETE CASCADE,
    salary_component_id UUID NOT NULL REFERENCES salary_components(id),
    type                VARCHAR(20) NOT NULL,  -- 'allowance' | 'deduction'
    amount              NUMERIC(15,2) NOT NULL DEFAULT 0
);

-- Reimbursement
CREATE TABLE reimbursements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id),
    category        VARCHAR(100) NOT NULL,    -- Transport, Makan, dll
    amount          NUMERIC(15,2) NOT NULL,
    date            DATE NOT NULL,
    description     TEXT NOT NULL,
    receipt_url     TEXT,
    status          reimbursement_status NOT NULL DEFAULT 'pending',
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    rejection_note  TEXT,
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_payroll_periods_tenant ON payroll_periods(tenant_id);
CREATE INDEX idx_payroll_slips_period ON payroll_slips(payroll_period_id);
CREATE INDEX idx_payroll_slips_employee ON payroll_slips(employee_id);
CREATE INDEX idx_reimbursements_employee ON reimbursements(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reimbursements_status ON reimbursements(tenant_id, status) WHERE deleted_at IS NULL;
