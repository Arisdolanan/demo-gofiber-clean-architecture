-- Down Migration: Remove Comprehensive Attendance Tracking

-- 1. Remove enhancements from student_attendance
ALTER TABLE student_attendance 
    DROP COLUMN IF EXISTS check_in_time,
    DROP COLUMN IF EXISTS check_out_time,
    DROP COLUMN IF EXISTS check_in_location,
    DROP COLUMN IF EXISTS check_out_location,
    DROP COLUMN IF EXISTS check_in_ip_address,
    DROP COLUMN IF EXISTS check_out_ip_address,
    DROP COLUMN IF EXISTS marked_by,
    DROP COLUMN IF EXISTS notes;

-- 2. Remove enhancements from teacher_attendance
ALTER TABLE teacher_attendance 
    DROP COLUMN IF EXISTS check_in_time,
    DROP COLUMN IF EXISTS check_out_time,
    DROP COLUMN IF EXISTS check_in_location,
    DROP COLUMN IF EXISTS check_out_location,
    DROP COLUMN IF EXISTS check_in_ip_address,
    DROP COLUMN IF EXISTS check_out_ip_address,
    DROP COLUMN IF EXISTS notes;

-- 3. Remove indexes
DROP INDEX IF EXISTS idx_student_attendance_date;
DROP INDEX IF EXISTS idx_teacher_attendance_date;

-- 4. Drop staff_attendance table and related objects
DROP TRIGGER IF EXISTS update_staff_attendance_updated_at ON staff_attendance;
DROP INDEX IF EXISTS idx_staff_attendance_employee_date;
DROP INDEX IF EXISTS idx_staff_attendance_status;
DROP INDEX IF EXISTS idx_staff_attendance_deleted_at;
DROP TABLE IF EXISTS staff_attendance;
