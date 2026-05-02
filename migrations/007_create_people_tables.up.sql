-- Up Migration: People Tables

CREATE TYPE teacher_status AS ENUM ('active', 'inactive');
CREATE TYPE student_status AS ENUM ('active', 'graduated');

CREATE TABLE IF NOT EXISTS teachers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    employee_number VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    status teacher_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, employee_number)
);

CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- NULL for lower grades usually
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_number VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    status student_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, student_number)
);

CREATE TABLE IF NOT EXISTS parents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS student_parents (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    parent_id BIGINT NOT NULL REFERENCES parents(id) ON DELETE CASCADE,
    relationship VARCHAR(50),
    UNIQUE(student_id, parent_id)
);

CREATE TABLE IF NOT EXISTS student_sections (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    UNIQUE(student_id, section_id, academic_session_id)
);

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

-- Add NIS and NISN columns to students table
ALTER TABLE students ADD COLUMN nis VARCHAR(50) DEFAULT '';
ALTER TABLE students ADD COLUMN nisn VARCHAR(50) DEFAULT '';

-- Add foreign key back to sections for homeroom_teacher_id
ALTER TABLE sections 
ADD CONSTRAINT fk_sections_homeroom_teacher 
FOREIGN KEY (homeroom_teacher_id) REFERENCES teachers(id) ON DELETE SET NULL;

-- Triggers for updated_at
CREATE TRIGGER update_teachers_updated_at BEFORE UPDATE ON teachers FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_students_updated_at BEFORE UPDATE ON students FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_parents_updated_at BEFORE UPDATE ON parents FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


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
