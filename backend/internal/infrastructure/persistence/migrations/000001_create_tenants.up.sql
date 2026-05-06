-- 000001_create_tenants.up.sql
-- Tenant = perusahaan yang menggunakan sistem ini (multi-tenant)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(50)  NOT NULL UNIQUE,   -- alias unik perusahaan
    email       VARCHAR(255) NOT NULL UNIQUE,   -- email contact perusahaan
    phone       VARCHAR(20),
    address     TEXT,
    logo_url    TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_tenants_slug ON tenants(slug) WHERE deleted_at IS NULL;
