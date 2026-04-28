-- Up Migration: Setup Comprehensive Attendance Tracking

-- 1. Create Staff Attendance Table
CREATE TABLE IF NOT EXISTS staff_attendance (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    check_in_time TIME,
    check_out_time TIME,
    check_in_location VARCHAR(500),
    check_out_location VARCHAR(500),
    check_in_ip_address VARCHAR(45),
    check_out_ip_address VARCHAR(45),
    status VARCHAR(50) NOT NULL, -- present, absent, late, sick, permission
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id),
    UNIQUE(employee_id, attendance_date)
);

-- Index for faster staff attendance queries
CREATE INDEX IF NOT EXISTS idx_staff_attendance_employee_date ON staff_attendance(employee_id, attendance_date);
CREATE INDEX IF NOT EXISTS idx_staff_attendance_status ON staff_attendance(status);
CREATE INDEX IF NOT EXISTS idx_staff_attendance_deleted_at ON staff_attendance(deleted_at);

-- Trigger for staff_attendance updated_at
CREATE TRIGGER update_staff_attendance_updated_at 
    BEFORE UPDATE ON staff_attendance 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 2. Enhance student_attendance with check-in/out tracking
ALTER TABLE student_attendance 
    ADD COLUMN IF NOT EXISTS check_in_time TIME,
    ADD COLUMN IF NOT EXISTS check_out_time TIME,
    ADD COLUMN IF NOT EXISTS check_in_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_out_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_in_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_out_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS marked_by BIGINT REFERENCES teachers(id),
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- 3. Enhance teacher_attendance with check-in/out tracking
ALTER TABLE teacher_attendance 
    ADD COLUMN IF NOT EXISTS check_in_time TIME,
    ADD COLUMN IF NOT EXISTS check_out_time TIME,
    ADD COLUMN IF NOT EXISTS check_in_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_out_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_in_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_out_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- 4. Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_student_attendance_date ON student_attendance(attendance_date);
CREATE INDEX IF NOT EXISTS idx_teacher_attendance_date ON teacher_attendance(attendance_date);

-- Add comments
COMMENT ON TABLE staff_attendance IS 'Staff attendance tracking records';
COMMENT ON COLUMN staff_attendance.employee_id IS 'Reference to staff table';
COMMENT ON COLUMN staff_attendance.attendance_date IS 'Date of attendance';
COMMENT ON COLUMN staff_attendance.check_in_time IS 'Check-in time';
COMMENT ON COLUMN staff_attendance.check_out_time IS 'Check-out time';
COMMENT ON COLUMN staff_attendance.check_in_location IS 'Check-in location (GPS coordinates or address)';
COMMENT ON COLUMN staff_attendance.check_out_location IS 'Check-out location (GPS coordinates or address)';
COMMENT ON COLUMN staff_attendance.check_in_ip_address IS 'IP address at check-in';
COMMENT ON COLUMN staff_attendance.check_out_ip_address IS 'IP address at check-out';
COMMENT ON COLUMN staff_attendance.status IS 'Attendance status: present, absent, late, sick, permission';
COMMENT ON COLUMN staff_attendance.notes IS 'Additional notes or remarks';

-- Comments for student_attendance
COMMENT ON COLUMN student_attendance.check_in_location IS 'Check-in location (GPS coordinates or address)';
COMMENT ON COLUMN student_attendance.check_out_location IS 'Check-out location (GPS coordinates or address)';
COMMENT ON COLUMN student_attendance.check_in_ip_address IS 'IP address at check-in';
COMMENT ON COLUMN student_attendance.check_out_ip_address IS 'IP address at check-out';

-- Comments for teacher_attendance
COMMENT ON COLUMN teacher_attendance.check_in_location IS 'Check-in location (GPS coordinates or address)';
COMMENT ON COLUMN teacher_attendance.check_out_location IS 'Check-out location (GPS coordinates or address)';
COMMENT ON COLUMN teacher_attendance.check_in_ip_address IS 'IP address at check-in';
COMMENT ON COLUMN teacher_attendance.check_out_ip_address IS 'IP address at check-out';
