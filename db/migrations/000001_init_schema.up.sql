-- QRAP Initial Schema -- Risk Assessment Platform

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ---- Enum types ----

CREATE TYPE risk_level AS ENUM (
    'CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'INFO'
);

CREATE TYPE assessment_status AS ENUM (
    'DRAFT', 'IN_PROGRESS', 'COMPLETED', 'ARCHIVED'
);

CREATE TYPE finding_category AS ENUM (
    'WEAK_ALGORITHM',
    'SHORT_KEY_LENGTH',
    'DEPRECATED_PROTOCOL',
    'MISSING_PQC',
    'CERTIFICATE_EXPIRY',
    'HARVEST_NOW_DECRYPT_LATER'
);

-- ---- Organizations table ----

CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_by  VARCHAR(255) NOT NULL DEFAULT 'system',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ---- Assessments table ----

CREATE TABLE assessments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    status          assessment_status NOT NULL DEFAULT 'DRAFT',
    overall_risk    risk_level,
    risk_score      DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    target_assets   TEXT[] NOT NULL DEFAULT '{}',
    assets_scanned  INTEGER NOT NULL DEFAULT 0,
    pqc_readiness   DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_by      VARCHAR(255) NOT NULL DEFAULT 'system',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_assessments_org ON assessments (organization_id);
CREATE INDEX idx_assessments_status ON assessments (status);
CREATE INDEX idx_assessments_risk ON assessments (overall_risk);

-- ---- Findings table ----

CREATE TABLE findings (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id           UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    category                finding_category NOT NULL,
    risk_level              risk_level NOT NULL,
    title                   VARCHAR(512) NOT NULL,
    description             TEXT NOT NULL,
    affected_asset          VARCHAR(512) NOT NULL,
    current_algorithm       VARCHAR(100),
    recommended_algorithm   VARCHAR(100),
    remediation             TEXT,
    discovered_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_findings_assessment ON findings (assessment_id);
CREATE INDEX idx_findings_risk_level ON findings (risk_level);
CREATE INDEX idx_findings_category ON findings (category);

-- ---- Audit log ----

CREATE TABLE qrap_audit_log (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID NOT NULL,
    action      VARCHAR(50) NOT NULL,
    actor       VARCHAR(255) NOT NULL,
    details     JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qrap_audit_entity ON qrap_audit_log (entity_type, entity_id);
CREATE INDEX idx_qrap_audit_created ON qrap_audit_log (created_at);

-- ---- Updated-at trigger ----

CREATE OR REPLACE FUNCTION qrap_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION qrap_update_updated_at();

CREATE TRIGGER trg_assessments_updated_at
    BEFORE UPDATE ON assessments
    FOR EACH ROW EXECUTE FUNCTION qrap_update_updated_at();
