-- 000015_create_overtime.up.sql

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'overtime_status') THEN
        CREATE TYPE overtime_status AS ENUM ('pending', 'approved', 'rejected');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS overtime_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id),
    date            DATE NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    duration_hours  NUMERIC(4,2) NOT NULL,
    reason          TEXT NOT NULL,
    status          overtime_status NOT NULL DEFAULT 'pending',
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    rejection_note  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(employee_id, date)
);

CREATE INDEX IF NOT EXISTS idx_overtime_employee ON overtime_requests(employee_id, status);
CREATE INDEX IF NOT EXISTS idx_overtime_tenant ON overtime_requests(tenant_id, date);
