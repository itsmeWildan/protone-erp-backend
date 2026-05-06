-- 000006_create_attendance.up.sql

CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'late', 'early_leave', 'on_leave', 'holiday');
CREATE TYPE overtime_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE attendances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id),
    date            DATE NOT NULL,
    check_in        TIMESTAMPTZ,
    check_out       TIMESTAMPTZ,
    status          attendance_status NOT NULL DEFAULT 'absent',
    location_in     TEXT,                    -- koordinat GPS check-in
    location_out    TEXT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(employee_id, date)
);

CREATE TABLE overtimes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL REFERENCES employees(id),
    date            DATE NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL,           -- durasi dalam menit
    reason          TEXT NOT NULL,
    status          overtime_status NOT NULL DEFAULT 'pending',
    approved_by     UUID REFERENCES employees(id),
    approved_at     TIMESTAMPTZ,
    rejection_note  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_attendances_employee_date ON attendances(employee_id, date);
CREATE INDEX idx_attendances_tenant_date ON attendances(tenant_id, date);
CREATE INDEX idx_overtimes_employee ON overtimes(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_overtimes_status ON overtimes(tenant_id, status) WHERE deleted_at IS NULL;
