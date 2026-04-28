-- Down Migration: Remove Staff Table

DROP TABLE IF EXISTS staff CASCADE;
DROP TYPE IF EXISTS staff_status;
