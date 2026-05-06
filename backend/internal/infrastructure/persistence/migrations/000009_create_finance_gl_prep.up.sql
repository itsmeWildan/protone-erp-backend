-- 000009_create_finance_gl_prep.up.sql
-- Persiapan untuk GL (General Ledger), AP (Accounts Payable), AR (Accounts Receivable)
-- Tabel ini dibuat dari awal agar future GL tidak butuh migrasi besar

CREATE TYPE account_type AS ENUM ('asset', 'liability', 'equity', 'revenue', 'expense');
CREATE TYPE account_normal AS ENUM ('debit', 'credit');
CREATE TYPE journal_status AS ENUM ('draft', 'posted', 'reversed');

-- Kelompok akun COA
CREATE TABLE account_groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    type        account_type NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

-- Chart of Accounts (COA)
CREATE TABLE chart_of_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    account_group_id UUID NOT NULL REFERENCES account_groups(id),
    code            VARCHAR(20) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    type            account_type NOT NULL,
    normal_balance  account_normal NOT NULL,
    is_header       BOOLEAN NOT NULL DEFAULT FALSE,   -- akun induk
    parent_id       UUID REFERENCES chart_of_accounts(id),
    level           SMALLINT NOT NULL DEFAULT 1,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    UNIQUE(tenant_id, code)
);

-- Journal entries (header)
CREATE TABLE journal_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    journal_no      VARCHAR(50) NOT NULL,
    date            DATE NOT NULL,
    description     TEXT NOT NULL,
    status          journal_status NOT NULL DEFAULT 'draft',
    source_type     VARCHAR(50),       -- 'payroll', 'reimbursement', 'manual', dsb
    source_id       UUID,              -- ID dari tabel asal (payroll_period_id, dll)
    total_debit     NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_credit    NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_by      UUID REFERENCES users(id),
    posted_by       UUID REFERENCES users(id),
    posted_at       TIMESTAMPTZ,
    reversed_by     UUID REFERENCES journal_entries(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, journal_no)
);

-- Journal lines (detail debit/kredit)
CREATE TABLE journal_lines (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id    UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    coa_id              UUID NOT NULL REFERENCES chart_of_accounts(id),
    description         TEXT,
    debit               NUMERIC(15,2) NOT NULL DEFAULT 0,
    credit              NUMERIC(15,2) NOT NULL DEFAULT 0,
    line_order          SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_coa_tenant ON chart_of_accounts(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_journal_entries_tenant ON journal_entries(tenant_id);
CREATE INDEX idx_journal_entries_date ON journal_entries(tenant_id, date);
CREATE INDEX idx_journal_entries_source ON journal_entries(source_type, source_id);
CREATE INDEX idx_journal_lines_entry ON journal_lines(journal_entry_id);
CREATE INDEX idx_journal_lines_coa ON journal_lines(coa_id);
