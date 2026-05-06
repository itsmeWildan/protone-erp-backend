-- 000016_create_budgeting.up.sql

CREATE TABLE department_budgets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    department_id   UUID NOT NULL REFERENCES departments(id),
    month           INT NOT NULL,
    year            INT NOT NULL,
    allocated_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    spent_amount     NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(department_id, month, year)
);

CREATE INDEX idx_budget_lookup ON department_budgets(tenant_id, department_id, month, year);
