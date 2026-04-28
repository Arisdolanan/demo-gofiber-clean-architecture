-- Down Migration: Remove subject_id from student_attendance
ALTER TABLE student_attendance DROP CONSTRAINT IF EXISTS student_attendance_unique_per_subject;
ALTER TABLE student_attendance ADD CONSTRAINT student_attendance_student_id_section_id_attendance_date_key 
    UNIQUE(student_id, section_id, attendance_date);
ALTER TABLE student_attendance DROP COLUMN IF EXISTS subject_id;
