-- QRAP Schema Rollback

DROP TRIGGER IF EXISTS trg_assessments_updated_at ON assessments;
DROP TRIGGER IF EXISTS trg_organizations_updated_at ON organizations;
DROP FUNCTION IF EXISTS qrap_update_updated_at();

DROP TABLE IF EXISTS qrap_audit_log;
DROP TABLE IF EXISTS findings;
DROP TABLE IF EXISTS assessments;
DROP TABLE IF EXISTS organizations;

DROP TYPE IF EXISTS finding_category;
DROP TYPE IF EXISTS assessment_status;
DROP TYPE IF EXISTS risk_level;
