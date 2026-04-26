-- Down Migration: People Tables

ALTER TABLE sections DROP CONSTRAINT IF EXISTS fk_sections_homeroom_teacher;

DROP TABLE IF EXISTS student_sections;
DROP TABLE IF EXISTS student_parents;
DROP TABLE IF EXISTS parents;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS teachers;

DROP TYPE IF EXISTS student_status;
DROP TYPE IF EXISTS teacher_status;
