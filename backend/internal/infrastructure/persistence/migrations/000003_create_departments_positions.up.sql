-- 000003_create_departments_positions.up.sql

CREATE TABLE departments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  UUID REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

CREATE TABLE positions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    department_id   UUID NOT NULL REFERENCES departments(id),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    level           SMALLINT NOT NULL DEFAULT 1,  -- 1=Staff, 2=Supervisor, 3=Manager, dst
    description     TEXT,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_departments_tenant ON departments(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_positions_tenant ON positions(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_positions_dept ON positions(department_id) WHERE deleted_at IS NULL;
