-- Migration: Add audit fields and ERD-compliant fields to all tables
-- This migration standardizes all tables with created_by, updated_by, deleted_by, deleted_at
-- and adds missing fields according to ERD specification

-- =====================================================
-- 1. SCHOOLS - Add audit fields and deleted_at
-- ===================================================
ALTER TABLE schools 
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 2. APP_PACKAGES - Add description and audit fields
-- =====================================================
ALTER TABLE app_packages 
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 3. SCHOOL_LICENSES - Add audit fields
-- =====================================================
ALTER TABLE school_licenses 
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 4. PAYMENTS - Add notes and audit fields
-- =====================================================
ALTER TABLE payments 
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 5. ACADEMIC_SESSIONS - Add code and audit fields
-- =====================================================
ALTER TABLE academic_sessions 
    ADD COLUMN IF NOT EXISTS code VARCHAR(50) UNIQUE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 6. CLASSES - Add code, description, and audit fields
-- =====================================================
ALTER TABLE classes 
    ADD COLUMN IF NOT EXISTS code VARCHAR(50) UNIQUE,
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
   ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 7. SECTIONS - Add code, room_number, capacity, rename homeroom_teacher_id to teacher_id, add audit fields
-- =====================================================
ALTER TABLE sections 
    ADD COLUMN IF NOT EXISTS code VARCHAR(50),
    ADD COLUMN IF NOT EXISTS room_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS capacity INT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- Rename homeroom_teacher_id to teacher_id for consistency with ERD
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'sections' AND column_name = 'homeroom_teacher_id'
    ) THEN
        ALTER TABLE sections RENAME COLUMN homeroom_teacher_id TO teacher_id;
    END IF;
END $$;

-- =====================================================
-- 8. SUBJECTS - Add description, credit_hours, and audit fields
-- =====================================================
ALTER TABLE subjects 
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS credit_hours INT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 9. TEACHERS - Add personal fields and audit fields
-- =====================================================
ALTER TABLE teachers 
    ADD COLUMN IF NOT EXISTS date_of_birth DATE,
    ADD COLUMN IF NOT EXISTS gender VARCHAR(10),
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS qualification VARCHAR(255),
    ADD COLUMN IF NOT EXISTS specialization VARCHAR(255),
    ADD COLUMN IF NOT EXISTS join_date DATE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 10. STUDENTS - Add personal fields and audit fields
-- =====================================================
ALTER TABLE students 
    ADD COLUMN IF NOT EXISTS date_of_birth DATE,
    ADD COLUMN IF NOT EXISTS gender VARCHAR(10),
    ADD COLUMN IF NOT EXISTS blood_type VARCHAR(5),
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS admission_date DATE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- Update student status enum to include all ERD values
ALTER TYPE student_status ADD VALUE IF NOT EXISTS 'inactive';
ALTER TYPE student_status ADD VALUE IF NOT EXISTS 'transferred';

-- =====================================================
-- 11. PARENTS - Add contact fields and audit fields
-- =====================================================
ALTER TABLE parents 
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS occupation VARCHAR(255),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 12. STUDENT_SECTIONS - Add roll_number, enrollment_date, status, and audit fields
-- =====================================================
ALTER TABLE student_sections 
    ADD COLUMN IF NOT EXISTS roll_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS enrollment_date DATE,
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 13. STUDENT_PARENTS - Add is_primary and audit fields
-- =====================================================
ALTER TABLE student_parents 
    ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 14. SCHEDULES - Add room_number and audit fields
-- =====================================================
ALTER TABLE schedules 
    ADD COLUMN IF NOT EXISTS room_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 15. EXAMS - Add title, description, duration_minutes, and audit fields
-- =====================================================
ALTER TABLE exams 
    ADD COLUMN IF NOT EXISTS title VARCHAR(255),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS duration_minutes INT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 16. EXAM_MARKS - Add notes, entered_by, entered_at, and audit fields
-- =====================================================
ALTER TABLE exam_marks 
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS entered_by BIGINT REFERENCES teachers(id),
    ADD COLUMN IF NOT EXISTS entered_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 17. STUDENT_ATTENDANCE - Add notes, marked_by, and audit fields
-- =====================================================
ALTER TABLE student_attendance 
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS marked_by BIGINT REFERENCES teachers(id),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 18. TEACHER_ATTENDANCE - Add check times, notes, and audit fields
-- =====================================================
ALTER TABLE teacher_attendance 
    ADD COLUMN IF NOT EXISTS check_in_time TIME,
    ADD COLUMN IF NOT EXISTS check_out_time TIME,
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 19. ROLES - Add description, is_system_role, and audit fields
-- =====================================================
ALTER TABLE roles 
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS is_system_role BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- ===================================================
-- 20. PERMISSIONS - Add description and audit fields
-- =====================================================
ALTER TABLE permissions 
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

-- =====================================================
-- 21. USER_ROLES - Add assigned_at, assigned_by
-- =====================================================
ALTER TABLE user_roles 
    ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS assigned_by BIGINT REFERENCES users(id);

-- =====================================================
-- 22. ACTIVITY_LOGS - Add ip_address, user_agent
-- =====================================================
ALTER TABLE activity_logs 
    ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS user_agent TEXT;

-- =====================================================
-- 23. Create Pivot Tables for Many-to-Many Relationships
-- =====================================================

-- Class Subjects pivot table
CREATE TABLE IF NOT EXISTS class_subjects (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id BIGINT NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id),
    UNIQUE(school_id, class_id, subject_id, academic_session_id)
);

CREATE INDEX IF NOT EXISTS idx_class_subjects_class_id ON class_subjects(class_id);
CREATE INDEX IF NOT EXISTS idx_class_subjects_subject_id ON class_subjects(subject_id);
CREATE INDEX IF NOT EXISTS idx_class_subjects_session_id ON class_subjects(academic_session_id);

-- Teacher Subjects pivot table
CREATE TABLE IF NOT EXISTS teacher_subjects (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    subject_id BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id),
    UNIQUE(school_id, teacher_id, section_id, subject_id, academic_session_id)
);

CREATE INDEX IF NOT EXISTS idx_teacher_subjects_teacher_id ON teacher_subjects(teacher_id);
CREATE INDEX IF NOT EXISTS idx_teacher_subjects_section_id ON teacher_subjects(section_id);
CREATE INDEX IF NOT EXISTS idx_teacher_subjects_subject_id ON teacher_subjects(subject_id);
CREATE INDEX IF NOT EXISTS idx_teacher_subjects_session_id ON teacher_subjects(academic_session_id);

-- =====================================================
-- 24. Create Indexes for Performance
-- =====================================================

-- Deleted_at indexes for soft delete queries
CREATE INDEX IF NOT EXISTS idx_schools_deleted_at ON schools(deleted_at);
CREATE INDEX IF NOT EXISTS idx_teachers_deleted_at ON teachers(deleted_at);
CREATE INDEX IF NOT EXISTS idx_students_deleted_at ON students(deleted_at);
CREATE INDEX IF NOT EXISTS idx_parents_deleted_at ON parents(deleted_at);
CREATE INDEX IF NOT EXISTS idx_exams_deleted_at ON exams(deleted_at);
CREATE INDEX IF NOT EXISTS idx_schedules_deleted_at ON schedules(deleted_at);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_student_sections_student_session ON student_sections(student_id, academic_session_id);
CREATE INDEX IF NOT EXISTS idx_schedules_section_day ON schedules(section_id, day_of_week);
CREATE INDEX IF NOT EXISTS idx_exams_section_date ON exams(section_id, exam_date);
CREATE INDEX IF NOT EXISTS idx_attendance_student_date ON student_attendance(student_id, attendance_date);
