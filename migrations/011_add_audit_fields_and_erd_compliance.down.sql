-- Down Migration: Revert audit fields and ERD compliance changes

-- IMPORTANT: This migration only reverts structural changes (columns).
-- It DOES NOT delete data created by subsequent migrations.
-- Data will remain intact to allow re-migration without issues.

-- =====================================================
-- Remove Pivot Tables
-- =====================================================
DROP TABLE IF EXISTS teacher_subjects CASCADE;
DROP TABLE IF EXISTS class_subjects CASCADE;

-- =====================================================
-- Remove indexes
-- =====================================================
DROP INDEX IF EXISTS idx_attendance_student_date;
DROP INDEX IF EXISTS idx_exams_section_date;
DROP INDEX IF EXISTS idx_schedules_section_day;
DROP INDEX IF EXISTS idx_student_sections_student_session;
DROP INDEX IF EXISTS idx_schedules_deleted_at;
DROP INDEX IF EXISTS idx_exams_deleted_at;
DROP INDEX IF EXISTS idx_parents_deleted_at;
DROP INDEX IF EXISTS idx_students_deleted_at;
DROP INDEX IF EXISTS idx_teachers_deleted_at;
DROP INDEX IF EXISTS idx_schools_deleted_at;

-- Note: Removing columns is commented out to preserve data
-- If you want to fully revert, uncomment the ALTER TABLE DROP COLUMN commands

/*
-- Revert Activity Logs
ALTER TABLE activity_logs 
    DROP COLUMN IF EXISTS user_agent,
    DROP COLUMN IF EXISTS ip_address;

-- Revert User Roles
ALTER TABLE user_roles 
    DROP COLUMN IF EXISTS assigned_by,
    DROP COLUMN IF EXISTS assigned_at;

-- Revert Permissions
ALTER TABLE permissions 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS description;

-- Revert Roles
ALTER TABLE roles 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS is_system_role,
    DROP COLUMN IF EXISTS description;

-- Revert Teacher Attendance
ALTER TABLE teacher_attendance 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS check_out_time,
    DROP COLUMN IF EXISTS check_in_time;

-- Revert Student Attendance
ALTER TABLE student_attendance 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS marked_by,
    DROP COLUMN IF EXISTS notes;

-- Revert Exam Marks
ALTER TABLE exam_marks 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS entered_at,
    DROP COLUMN IF EXISTS entered_by,
    DROP COLUMN IF EXISTS notes;

-- Revert Exams
ALTER TABLE exams 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS duration_minutes,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS title;

-- Revert Schedules
ALTER TABLE schedules 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS room_number;

-- Revert Student Parents
ALTER TABLE student_parents 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS is_primary;

-- Revert Student Sections
ALTER TABLE student_sections 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS enrollment_date,
    DROP COLUMN IF EXISTS roll_number;

-- Revert Parents
ALTER TABLE parents 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS occupation,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS phone;

-- Revert Students
ALTER TABLE students 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS admission_date,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS blood_type,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS date_of_birth;

-- Revert Teachers
ALTER TABLE teachers 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS join_date,
    DROP COLUMN IF EXISTS specialization,
    DROP COLUMN IF EXISTS qualification,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS date_of_birth;

-- Revert Subjects
ALTER TABLE subjects 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS credit_hours,
    DROP COLUMN IF EXISTS description;

-- Revert Sections (including rename)
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'sections' AND column_name = 'teacher_id'
    ) THEN
        ALTER TABLE sections RENAME COLUMN teacher_id TO homeroom_teacher_id;
    END IF;
END $$;

ALTER TABLE sections 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS capacity,
    DROP COLUMN IF EXISTS room_number,
    DROP COLUMN IF EXISTS code;

-- Revert Classes
ALTER TABLE classes 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS code;

-- Revert Academic Sessions
ALTER TABLE academic_sessions 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS code;

-- Revert Payments
ALTER TABLE payments 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS notes;

-- Revert School Licenses
ALTER TABLE school_licenses 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at;

-- Revert App Packages
ALTER TABLE app_packages 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS description;

-- Revert Schools
ALTER TABLE schools 
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at;
*/
