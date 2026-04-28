-- Up Migration: Add subject_id to student_attendance and update unique constraint
ALTER TABLE student_attendance 
    ADD COLUMN IF NOT EXISTS subject_id BIGINT REFERENCES subjects(id) ON DELETE SET NULL;

-- Update unique constraint to include subject_id
-- First drop existing constraint
ALTER TABLE student_attendance DROP CONSTRAINT IF EXISTS student_attendance_student_id_section_id_attendance_date_key;

-- Add new constraint including subject_id
-- Note: PostgreSQL unique constraint treats NULL as distinct, so we use COALESCE if we want to allow one null subject record
-- But in our case, subject_id is mandatory now, so we can just include it.
ALTER TABLE student_attendance ADD CONSTRAINT student_attendance_unique_per_subject 
    UNIQUE(student_id, section_id, attendance_date, subject_id);

-- Add index for subject-based reports
CREATE INDEX IF NOT EXISTS idx_student_attendance_subject_id ON student_attendance(subject_id);
