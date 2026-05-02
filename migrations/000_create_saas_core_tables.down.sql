-- Down Migration: SaaS Core Tables

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS school_licenses;
DROP TABLE IF EXISTS app_packages;
DROP TABLE IF EXISTS schools;

DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS license_status;
DROP TYPE IF EXISTS school_status;
DROP TYPE IF EXISTS school_level;
