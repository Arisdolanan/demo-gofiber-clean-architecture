-- Up Migration: Operational Tables

CREATE TABLE IF NOT EXISTS schedules (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    subject_id BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL, -- 0-6 (Sunday-Saturday)
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exams (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    subject_id BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    exam_type VARCHAR(100) NOT NULL,
    exam_date DATE NOT NULL,
    max_score INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS exam_marks (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    exam_id BIGINT NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    score DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(exam_id, student_id)
);

CREATE TABLE IF NOT EXISTS student_attendance (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    academic_session_id BIGINT NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL, -- Present, Absent, Late, Sick
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, section_id, attendance_date)
);

CREATE TABLE IF NOT EXISTS teacher_attendance (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL, -- Present, Absent, Sick, Leave
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(teacher_id, attendance_date)
);

CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id BIGINT REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    reference_type VARCHAR(100), -- Link to specific module (exam, attendance, etc)
    reference_id BIGINT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Triggers for updated_at
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_exams_updated_at BEFORE UPDATE ON exams FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_exam_marks_updated_at BEFORE UPDATE ON exam_marks FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


-- Up Migration: Setup Comprehensive Attendance Tracking

-- 1. Create Staff Attendance Table
CREATE TABLE IF NOT EXISTS staff_attendance (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    check_in_time TIME,
    check_out_time TIME,
    check_in_location VARCHAR(500),
    check_out_location VARCHAR(500),
    check_in_ip_address VARCHAR(45),
    check_out_ip_address VARCHAR(45),
    check_in_device TEXT,
    check_out_device TEXT,
    status VARCHAR(50) NOT NULL, -- present, absent, late, sick, permission
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id),
    UNIQUE(school_id, employee_id, attendance_date)
);

-- 2. Enhance student_attendance (already created in 009)
ALTER TABLE student_attendance 
    ADD COLUMN IF NOT EXISTS subject_id BIGINT REFERENCES subjects(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS check_in_time TIME,
    ADD COLUMN IF NOT EXISTS check_out_time TIME,
    ADD COLUMN IF NOT EXISTS check_in_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_out_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_in_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_out_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_in_device TEXT,
    ADD COLUMN IF NOT EXISTS check_out_device TEXT,
    ADD COLUMN IF NOT EXISTS marked_by BIGINT REFERENCES teachers(id),
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- Update unique constraint for student_attendance
ALTER TABLE student_attendance DROP CONSTRAINT IF EXISTS student_attendance_student_id_section_id_attendance_date_key;
ALTER TABLE student_attendance ADD CONSTRAINT student_attendance_unique_per_subject 
    UNIQUE(student_id, section_id, attendance_date, subject_id);

-- 3. Enhance teacher_attendance
ALTER TABLE teacher_attendance 
    ADD COLUMN IF NOT EXISTS check_in_time TIME,
    ADD COLUMN IF NOT EXISTS check_out_time TIME,
    ADD COLUMN IF NOT EXISTS check_in_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_out_location VARCHAR(500),
    ADD COLUMN IF NOT EXISTS check_in_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_out_ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS check_in_device TEXT,
    ADD COLUMN IF NOT EXISTS check_out_device TEXT,
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- 4. Indexes
CREATE INDEX IF NOT EXISTS idx_staff_attendance_employee_date ON staff_attendance(employee_id, attendance_date);
CREATE INDEX IF NOT EXISTS idx_student_attendance_subject_id ON student_attendance(subject_id);
CREATE INDEX IF NOT EXISTS idx_student_attendance_date ON student_attendance(attendance_date);
CREATE INDEX IF NOT EXISTS idx_teacher_attendance_date ON teacher_attendance(attendance_date);

-- Trigger
CREATE TRIGGER update_staff_attendance_updated_at BEFORE UPDATE ON staff_attendance FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
