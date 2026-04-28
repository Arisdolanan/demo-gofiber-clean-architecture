-- Up Migration: Staff Table
-- This migration creates the staff table before attendance tracking

-- Create staff_status enum if not exists
DO $$ BEGIN
    CREATE TYPE staff_status AS ENUM ('active', 'inactive');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create staff table
CREATE TABLE IF NOT EXISTS staff (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    employee_number VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    date_of_birth DATE,
    gender VARCHAR(10) DEFAULT 'male',
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    position VARCHAR(100),
    department VARCHAR(100),
    join_date DATE,
    status staff_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id),
    UNIQUE(school_id, employee_number)
);

-- Index for faster queries
CREATE INDEX IF NOT EXISTS idx_staff_school_id ON staff(school_id);
CREATE INDEX IF NOT EXISTS idx_staff_status ON staff(status);
CREATE INDEX IF NOT EXISTS idx_staff_deleted_at ON staff(deleted_at);

-- Trigger for updated_at
DO $$ BEGIN
    CREATE TRIGGER update_staff_updated_at 
        BEFORE UPDATE ON staff 
        FOR EACH ROW 
        EXECUTE PROCEDURE update_updated_at_column();
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add comments
COMMENT ON TABLE staff IS 'School staff members (non-teaching)';
COMMENT ON COLUMN staff.employee_number IS 'Unique employee identification number within a school';
COMMENT ON COLUMN staff.position IS 'Job title or position';
COMMENT ON COLUMN staff.department IS 'Department or unit';
