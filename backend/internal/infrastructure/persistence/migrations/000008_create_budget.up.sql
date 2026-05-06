-- 000008_create_budget.up.sql

CREATE TYPE budget_status AS ENUM ('draft', 'active', 'closed');
CREATE TYPE allocation_status AS ENUM ('pending', 'approved', 'rejected');

-- Kategori anggaran
CREATE TABLE budget_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

-- Master anggaran per periode
CREATE TABLE budgets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    budget_category_id UUID NOT NULL REFERENCES budget_categories(id),
    department_id   UUID REFERENCES departments(id),
    period_month    SMALLINT,          -- NULL = anggaran tahunan
    period_year     SMALLINT NOT NULL,
    amount          NUMERIC(15,2) NOT NULL DEFAULT 0,
    used_amount     NUMERIC(15,2) NOT NULL DEFAULT 0,   -- terpakai
    status          budget_status NOT NULL DEFAULT 'draft',
    notes           TEXT,
    created_by      UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Pengajuan penggunaan anggaran
CREATE TABLE budget_allocations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    budget_id   UUID NOT NULL REFERENCES budgets(id),
    amount      NUMERIC(15,2) NOT NULL,
    description TEXT NOT NULL,
    status      allocation_status NOT NULL DEFAULT 'pending',
    requested_by UUID NOT NULL REFERENCES users(id),
    approved_by  UUID REFERENCES users(id),
    approved_at  TIMESTAMPTZ,
    rejection_note TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_budgets_tenant ON budgets(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budgets_category ON budgets(budget_category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_allocations_budget ON budget_allocations(budget_id) WHERE deleted_at IS NULL;
