-- Up Migration: Seed Test Data for API Development & Testing
-- This migration adds comprehensive sample data matching frontend mock data

-- ============================================
-- 1. SAMPLE SCHOOL
-- Using INSERT ... ON CONFLICT to handle re-runs
-- ============================================
INSERT INTO schools (code, name, email, phone, address, city, province, country, school_level, status) VALUES 
('SCH001', 'SMA Negeri 1 Jakarta', 'info@sman1jakarta.sch.id', '021-1234567', 'Jl. Pendidikan No. 1, Menteng', 'Jakarta Pusat', 'DKI Jakarta', 'Indonesia', 'SMA', 'active')
ON CONFLICT (code) DO UPDATE SET 
  name = EXCLUDED.name,
  email = EXCLUDED.email,
  phone = EXCLUDED.phone,
  address = EXCLUDED.address,
  city = EXCLUDED.city,
  province = EXCLUDED.province,
  country = EXCLUDED.country,
  school_level = EXCLUDED.school_level,
  status = EXCLUDED.status;

-- ============================================
-- 2. SAMPLE APP PACKAGES
-- Using INSERT ... ON CONFLICT to handle re-runs
-- ============================================
INSERT INTO app_packages (code, name, price_monthly, price_yearly, max_students, max_teachers, is_active) VALUES
('PKG-BASIC', 'Paket Basic', 500000.00, 5000000.00, 500, 50, true),
('PKG-PRO', 'Paket Professional', 1000000.00, 10000000.00, 1000, 100, true),
('PKG-ENTERPRISE', 'Paket Enterprise', 2000000.00, 20000000.00, 99999, 9999, true)
ON CONFLICT (code) DO UPDATE SET
  name = EXCLUDED.name,
  price_monthly = EXCLUDED.price_monthly,
  price_yearly = EXCLUDED.price_yearly,
  max_students = EXCLUDED.max_students,
  max_teachers = EXCLUDED.max_teachers,
  is_active = EXCLUDED.is_active;

-- ============================================
-- 3. SAMPLE SCHOOL LICENSE
-- First ensure school and app_package exist, then insert
-- ============================================
DO $$
DECLARE
  v_school_id BIGINT;
  v_package_id BIGINT;
BEGIN
  -- Get school_id
  SELECT id INTO v_school_id FROM schools WHERE code = 'SCH001';
  IF NOT FOUND THEN
    RAISE EXCEPTION 'School SCH001 not found!';
  END IF;
  
  -- Get package_id  
  SELECT id INTO v_package_id FROM app_packages WHERE code = 'PKG-PRO';
  IF NOT FOUND THEN
    RAISE EXCEPTION 'Package PKG-PRO not found!';
  END IF;
  
  -- Insert license only if not exists
  IF NOT EXISTS (SELECT 1 FROM school_licenses WHERE license_key = 'LIC-SCH001-2024-PRO') THEN
    INSERT INTO school_licenses (school_id, app_package_id, license_key, start_date, end_date, status)
    VALUES (v_school_id, v_package_id, 'LIC-SCH001-2024-PRO', '2024-01-01 00:00:00', '2025-12-31 23:59:59', 'active');
  END IF;
END $$;

-- ============================================
-- 4. SAMPLE USERS
-- Note: Password for all users is "@SuperIndo1" 
-- Using bcrypt hash: $2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC
-- Using ON CONFLICT DO NOTHING to handle re-runs
-- ============================================

-- Admin User
INSERT INTO users (email, password, username, user_type, email_verified_at) VALUES
('admin@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'admin', 'super_admin', CURRENT_TIMESTAMP)
ON CONFLICT (username) DO NOTHING;

-- Teachers (5 teachers)
INSERT INTO users (email, password, username, user_type, email_verified_at) VALUES
('budi.santoso@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'budi.santoso', 'teacher', CURRENT_TIMESTAMP),
('siti.aminah@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'siti.aminah', 'teacher', CURRENT_TIMESTAMP),
('agus.pratama@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'agus.pratama', 'teacher', CURRENT_TIMESTAMP),
('ani.wijaya@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'ani.wijaya', 'teacher', CURRENT_TIMESTAMP),
('joko.susilo@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'joko.susilo', 'teacher', CURRENT_TIMESTAMP)
ON CONFLICT (username) DO NOTHING;

-- Students (10 students for testing)
INSERT INTO users (email, password, username, user_type, email_verified_at) VALUES
('aditya.perkasa@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'aditya.perkasa', 'student', CURRENT_TIMESTAMP),
('bunga.citra@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'bunga.citra', 'student', CURRENT_TIMESTAMP),
('chandra.wijaya@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'chandra.wijaya', 'student', CURRENT_TIMESTAMP),
('dewi.sartika@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'dewi.sartika', 'student', CURRENT_TIMESTAMP),
('eko.prasetyo@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'eko.prasetyo', 'student', CURRENT_TIMESTAMP),
('fani.rahmawati@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'fani.rahmawati', 'student', CURRENT_TIMESTAMP),
('guntur.prabowo@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'guntur.prabowo', 'student', CURRENT_TIMESTAMP),
('hana.pertiwi@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'hana.pertiwi', 'student', CURRENT_TIMESTAMP),
('indra.lesmana@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'indra.lesmana', 'student', CURRENT_TIMESTAMP),
('juli.anissa@student.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'juli.anissa', 'student', CURRENT_TIMESTAMP)
ON CONFLICT (username) DO NOTHING;

-- Parents (3 parents)
INSERT INTO users (email, password, username, user_type, email_verified_at) VALUES
('orangtua1@parent.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'orangtua1', 'parent', CURRENT_TIMESTAMP),
('orangtua2@parent.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'orangtua2', 'parent', CURRENT_TIMESTAMP),
('orangtua3@parent.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'orangtua3', 'parent', CURRENT_TIMESTAMP)
ON CONFLICT (username) DO NOTHING;

-- Staff (5 staff members)
INSERT INTO users (email, password, username, user_type, email_verified_at) VALUES
('kepala.sekolah@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'kepala.sekolah', 'staff', CURRENT_TIMESTAMP),
('admin.tu@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'admin.tu', 'staff', CURRENT_TIMESTAMP),
('operator@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'operator', 'staff', CURRENT_TIMESTAMP),
('pustakawan@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'pustakawan', 'staff', CURRENT_TIMESTAMP),
('konselor@sman1jakarta.sch.id', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'konselor', 'staff', CURRENT_TIMESTAMP)
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 5. ACADEMIC SESSIONS (Tahun Ajaran)
-- Match frontend mock: 2023/2024, 2022/2023, 2021/2022
-- ============================================
INSERT INTO academic_sessions (school_id, name, code, start_date, end_date, is_active) VALUES
(1, '2023/2024', 'AY-2023-24', '2023-07-10', '2024-06-15', true),
(1, '2022/2023', 'AY-2022-23', '2022-07-11', '2023-06-17', false),
(1, '2021/2022', 'AY-2021-22', '2021-07-12', '2022-06-18', false)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 6. CLASSES (Kelas)
-- Match frontend mock: Kelas X, XI, XII
-- ============================================
INSERT INTO classes (school_id, name, code, level, grade_number) VALUES
(1, 'Kelas X', 'X', 'SMA', 10),
(1, 'Kelas XI', 'XI', 'SMA', 11),
(1, 'Kelas XII', 'XII', 'SMA', 12)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 7. SUBJECTS (Mata Pelajaran)
-- Match frontend mock: Matematika, B. Indonesia, Fisika, Biologi, B. Agama
-- ============================================
INSERT INTO subjects (school_id, name, code, credit_hours) VALUES
(1, 'Matematika', 'MAT-X', 4),
(1, 'Bahasa Indonesia', 'BIN-X', 4),
(1, 'Fisika', 'FIS-X', 3),
(1, 'Biologi', 'BIO-X', 3),
(1, 'Pendidikan Agama', 'AGM-X', 2),
(1, 'Bahasa Inggris', 'BING-X', 3),
(1, 'Kimia', 'KIM-X', 3),
(1, 'Sejarah', 'SEJ-X', 2),
(1, 'Pendidikan Kewarganegaraan', 'PKN-X', 2),
(1, 'Seni Budaya', 'SNB-X', 2)
ON CONFLICT (school_id, code) DO NOTHING;

-- ============================================
-- 8. TEACHERS
-- Match frontend mock data
-- ============================================
INSERT INTO teachers (user_id, school_id, employee_number, full_name, date_of_birth, gender, phone, email, address, qualification, specialization, join_date, status) VALUES
(2, 1, 'TCH-001', 'Budi Santoso, S.Pd', '1985-05-15', 'male', '081234567890', 'budi.santoso@sman1jakarta.sch.id', 'Jakarta', 'S1 Pendidikan Matematika', 'Matematika', '2010-08-01', 'active'),
(3, 1, 'TCH-002', 'Siti Aminah, M.Pd', '1990-03-20', 'female', '081234567891', 'siti.aminah@sman1jakarta.sch.id', 'Jakarta', 'S2 Pendidikan Bahasa', 'Bahasa Indonesia', '2015-01-10', 'active'),
(4, 1, 'TCH-003', 'Agus Pratama, S.Pd', '1988-07-25', 'male', '081234567892', 'agus.pratama@sman1jakarta.sch.id', 'Jakarta', 'S1 Fisika', 'Fisika', '2012-06-15', 'active'),
(5, 1, 'TCH-004', 'Ani Wijaya, S.Pd', '1992-11-10', 'female', '081234567893', 'ani.wijaya@sman1jakarta.sch.id', 'Jakarta', 'S1 Biologi', 'Biologi', '2016-07-20', 'active'),
(6, 1, 'TCH-005', 'Joko Susilo, M.Pd', '1980-01-05', 'male', '081234567894', 'joko.susilo@sman1jakarta.sch.id', 'Jakarta', 'S2 Pendidikan Agama', 'Pendidikan Agama', '2005-09-01', 'active');

-- ============================================
-- 9. SECTIONS (Kelas Paralel)
-- Match frontend mock: Kelas X-A, X-B, XI-A, dll
-- Using dynamic lookups to ensure FK constraints are met
-- ============================================
DO $$
DECLARE
  v_class_x_id BIGINT;
  v_class_xi_id BIGINT;
  v_academic_session_id BIGINT;
  v_teacher1_id BIGINT;
  v_teacher2_id BIGINT;
  v_teacher3_id BIGINT;
BEGIN
  -- Lookup IDs
  SELECT id INTO v_class_x_id FROM classes WHERE code = 'X' LIMIT 1;
  SELECT id INTO v_class_xi_id FROM classes WHERE code = 'XI' LIMIT 1;
  SELECT id INTO v_academic_session_id FROM academic_sessions WHERE code = 'AY-2023-24' LIMIT 1;
  SELECT id INTO v_teacher1_id FROM teachers WHERE employee_number = 'TCH-001' LIMIT 1;
  SELECT id INTO v_teacher2_id FROM teachers WHERE employee_number = 'TCH-002' LIMIT 1;
  SELECT id INTO v_teacher3_id FROM teachers WHERE employee_number = 'TCH-003' LIMIT 1;
  
  -- Insert sections
  INSERT INTO sections (class_id, academic_session_id, name, code, room_number, capacity, teacher_id) VALUES
  (v_class_x_id, v_academic_session_id, 'A', 'X-A', '101', 35, v_teacher1_id),
  (v_class_x_id, v_academic_session_id, 'B', 'X-B', '102', 35, v_teacher2_id),
  (v_class_xi_id, v_academic_session_id, 'A', 'XI-A', '201', 30, v_teacher3_id);
END $$;

-- ============================================
-- 10. STUDENTS
-- Match frontend mock: Aditya Perkasa, Bunga Citra, Chandra Wijaya, Dewi Sartika
-- ============================================
INSERT INTO students (user_id, school_id, student_number, full_name, date_of_birth, gender, blood_type, phone, email, address, admission_date, status) VALUES
(7, 1, 'STD-2023-001', 'Aditya Perkasa', '2008-01-15', 'male', 'O', '089876543210', 'aditya.perkasa@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(8, 1, 'STD-2023-002', 'Bunga Citra', '2008-02-20', 'female', 'A', '089876543211', 'bunga.citra@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(9, 1, 'STD-2023-003', 'Chandra Wijaya', '2008-03-25', 'male', 'B', '089876543212', 'chandra.wijaya@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(10, 1, 'STD-2023-004', 'Dewi Sartika', '2008-04-10', 'female', 'AB', '089876543213', 'dewi.sartika@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(11, 1, 'STD-2023-005', 'Eko Prasetyo', '2008-05-05', 'male', 'O', '089876543214', 'eko.prasetyo@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(12, 1, 'STD-2023-006', 'Fani Rahmawati', '2008-06-15', 'female', 'A', '089876543215', 'fani.rahmawati@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(13, 1, 'STD-2023-007', 'Guntur Prabowo', '2008-07-20', 'male', 'B', '089876543216', 'guntur.prabowo@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(14, 1, 'STD-2023-008', 'Hana Pertiwi', '2008-08-25', 'female', 'O', '089876543217', 'hana.pertiwi@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(15, 1, 'STD-2023-009', 'Indra Lesmana', '2008-09-10', 'male', 'A', '089876543218', 'indra.lesmana@student.sch.id', 'Jakarta', '2023-07-10', 'active'),
(16, 1, 'STD-2023-010', 'Juli Anissa', '2008-10-05', 'female', 'AB', '089876543219', 'juli.anissa@student.sch.id', 'Jakarta', '2023-07-10', 'active');

-- ============================================
-- 11. STUDENT SECTIONS (Enrollment)
-- Assign students to sections
-- ============================================
INSERT INTO student_sections (student_id, section_id, academic_session_id, roll_number, enrollment_date, status) VALUES
-- Kelas X-A students
(1, 1, 1, 'X-A-01', '2023-07-10', 'active'),  -- Aditya
(2, 1, 1, 'X-A-02', '2023-07-10', 'active'),  -- Bunga
(3, 1, 1, 'X-A-03', '2023-07-10', 'active'),  -- Chandra
(4, 1, 1, 'X-A-04', '2023-07-10', 'active'),  -- Dewi

-- Kelas X-B students
(5, 2, 1, 'X-B-01', '2023-07-10', 'active'),  -- Eko
(6, 2, 1, 'X-B-02', '2023-07-10', 'active'),  -- Fani
(7, 2, 1, 'X-B-03', '2023-07-10', 'active'),  -- Guntur

-- Kelas XI-A students
(8, 3, 1, 'XI-A-01', '2023-07-10', 'active'), -- Hana
(9, 3, 1, 'XI-A-02', '2023-07-10', 'active'), -- Indra
(10, 3, 1, 'XI-A-03', '2023-07-10', 'active');-- Juli

-- ============================================
-- 12. PARENTS
-- ============================================
INSERT INTO parents (user_id, school_id, full_name, phone, email, address, occupation) VALUES
(17, 1, 'Bapak Heru Perkasa', '081111111111', 'orangtua1@parent.sch.id', 'Jakarta', 'Wiraswasta'),
(18, 1, 'Ibu Citra Dewi', '081111111112', 'orangtua2@parent.sch.id', 'Jakarta', 'Guru'),
(19, 1, 'Bapak Wijaya Kusuma', '081111111113', 'orangtua3@parent.sch.id', 'Jakarta', 'PNS');

-- ============================================
-- 13. STUDENT PARENTS (Relationship)
-- Link parents to students
-- ============================================
INSERT INTO student_parents (student_id, parent_id, relationship, is_primary) VALUES
(1, 1, 'father', true),   -- Aditya - Bapak Heru
(2, 2, 'mother', true),   -- Bunga - Ibu Citra
(3, 3, 'father', true);   -- Chandra - Bapak Wijaya

-- ============================================
-- 13b. STAFF DATA - MOVED TO MIGRATION 013
-- Staff seed data is in migration 013 (after table creation)
-- to avoid "relation does not exist" error
-- ============================================

-- ============================================
-- 14. SCHEDULES (Jadwal Pelajaran)
-- Match frontend mock schedules
-- ============================================
INSERT INTO schedules (section_id, subject_id, teacher_id, academic_session_id, day_of_week, start_time, end_time, room_number) VALUES
-- Senin, Kelas X-A
(1, 1, 1, 1, 1, '07:30', '09:00', '101'),  -- Senin 07:30-09:00 Matematika - Budi Santoso
(1, 2, 2, 1, 1, '09:00', '10:30', '101'),  -- Senin 09:00-10:30 B. Indonesia - Siti Aminah

-- Selasa, Kelas X-B
(2, 3, 3, 1, 2, '08:00', '09:30', '201'),  -- Selasa 08:00-09:30 Fisika - Agus Pratama
(2, 1, 1, 1, 2, '09:30', '11:00', '201'),  -- Selasa 09:30-11:00 Matematika - Budi Santoso

-- Rabu, Kelas XI-A
(3, 4, 4, 1, 3, '07:30', '09:00', '201'),  -- Rabu 07:30-09:00 Biologi - Ani Wijaya
(3, 5, 5, 1, 3, '09:00', '10:30', '201');  -- Rabu 09:00-10:30 B. Agama - Joko Susilo

-- ============================================
-- 15. EXAMS (Ujian)
-- Match frontend mock: Ulangan Harian 1 Aljabar, UTS B. Indonesia
-- ============================================
INSERT INTO exams (section_id, subject_id, academic_session_id, title, description, exam_type, exam_date, duration_minutes, max_score) VALUES
-- Ulangan Harian 1 - Aljabar (Matematika, Kelas X-A)
(1, 1, 1, 'Ulangan Harian 1 - Aljabar', 'Ulangan harian materi Aljabar Linear', 'daily', '2023-09-15', 90, 100),

-- UTS Ganjil - B. Indonesia (Kelas X-A)
(1, 2, 1, 'UTS Ganjil - Bahasa Indonesia', 'Ujian Tengah Semester Ganjil', 'midterm', '2023-10-20', 120, 100),

-- UTS Fisika (Kelas X-B)
(2, 3, 1, 'UTS Ganjil - Fisika', 'Ujian Tengah Semester Fisika', 'midterm', '2023-10-22', 90, 100),

-- UAS Biologi (Kelas XI-A)
(3, 4, 1, 'UAS Biologi', 'Ujian Akhir Semester Biologi', 'final', '2023-12-01', 120, 100);

-- ============================================
-- 16. EXAM MARKS (Nilai Ujian)
-- Match frontend mock marks
-- ============================================
INSERT INTO exam_marks (exam_id, student_id, score, notes, entered_by, entered_at) VALUES
-- Marks for Ulangan Harian 1 - Aljabar (Exam ID 1)
(1, 1, 85.0, 'Sangat baik dalam aljabar', 1, '2023-09-16'),  -- Aditya: 85
(1, 2, 90.0, 'Pemahaman sangat baik', 1, '2023-09-16'),      -- Bunga: 90
(1, 3, 78.0, 'Perlu latihan lebih', 1, '2023-09-16'),         -- Chandra: 78
(1, 4, 92.0, 'Nilai sempurna', 1, '2023-09-16'),              -- Dewi: 92

-- Marks for UTS B. Indonesia (Exam ID 2)
(2, 1, 88.0, 'Baik dalam menulis', 2, '2023-10-21'),          -- Aditya: 88
(2, 2, 95.0, 'Sangat berbakat', 2, '2023-10-21'),             -- Bunga: 95
(2, 3, 82.0, 'Cukup baik', 2, '2023-10-21'),                  -- Chandra: 82
(2, 4, 91.0, 'Sangat baik', 2, '2023-10-21');                 -- Dewi: 91

-- ============================================
-- 17. TEACHER ATTENDANCE
-- ============================================
INSERT INTO teacher_attendance (teacher_id, attendance_date, check_in_time, check_out_time, status) VALUES
(1, '2024-01-15', '2024-01-15 07:00:00', '2024-01-15 15:00:00', 'present'),
(2, '2024-01-15', '2024-01-15 07:15:00', '2024-01-15 15:30:00', 'present'),
(3, '2024-01-15', '2024-01-15 07:30:00', '2024-01-15 14:00:00', 'present'),
(4, '2024-01-15', '2024-01-15 08:00:00', '2024-01-15 15:00:00', 'late'),
(5, '2024-01-15', '2024-01-15 07:00:00', '2024-01-15 15:00:00', 'present');

-- ============================================
-- 18. STUDENT ATTENDANCE
-- ============================================
INSERT INTO student_attendance (student_id, section_id, academic_session_id, attendance_date, status, marked_by) VALUES
-- Attendance for Kelas X-A on 2024-01-15
(1, 1, 1, '2024-01-15', 'present', 1),   -- Aditya: present
(2, 1, 1, '2024-01-15', 'present', 1),   -- Bunga: present
(3, 1, 1, '2024-01-15', 'late', 1),      -- Chandra: late
(4, 1, 1, '2024-01-15', 'present', 1),   -- Dewi: present

-- Some students absent on different date
(5, 2, 1, '2024-01-16', 'absent', 2),    -- Eko: absent
(6, 2, 1, '2024-01-16', 'present', 2),   -- Fani: present
(7, 2, 1, '2024-01-16', 'sick', 2);      -- Guntur: sick

-- ============================================
-- 19. NOTIFICATIONS
-- ============================================
INSERT INTO notifications (user_id, school_id, title, message, is_read, reference_type) VALUES
-- Notification for Admin
(1, 1, 'Selamat Datang', 'Sistem Akademik berhasil diaktifkan untuk SMA Negeri 1 Jakarta', false, 'system'),

-- Notifications for Teachers
(2, 1, 'Jadwal Mengajar', 'Anda memiliki jadwal mengajar hari Senin pukul 07:30', false, 'schedule'),
(3, 1, 'Input Nilai', 'Silakan input nilai UTS B. Indonesia sebelum tanggal 30 Oktober', true, 'exam'),

-- Notifications for Students
(7, 1, 'Jadwal Ujian', 'Ulangan Harian 1 Aljabar akan dilaksanakan tanggal 15 September', false, 'exam'),
(8, 1, 'Pengumuman', 'Nilai UTS sudah dapat dilihat di menu Nilai', false, 'marks');

-- ============================================
-- 20. ROLES (RBAC)
-- ============================================
INSERT INTO roles (code, name, description) VALUES
('super_admin', 'super_admin', 'Super Administrator'),
('school_admin', 'school_admin', 'School Administrator'),
('teacher', 'teacher', 'Teacher'),
('student', 'student', 'Student'),
('parent', 'parent', 'Parent/Guardian'),
('staff', 'staff', 'School Staff')
ON CONFLICT (school_id, code) DO NOTHING;

-- ============================================
-- 21. PERMISSIONS (RBAC)
-- ============================================
INSERT INTO permissions (permission_code, module_name, permission_name) VALUES
-- Academic permissions
('academic_session.create', 'academic_session', 'Create academic sessions'),
('academic_session.read', 'academic_session', 'View academic sessions'),
('academic_session.update', 'academic_session', 'Update academic sessions'),
('academic_session.delete', 'academic_session', 'Delete academic sessions'),

('class.create', 'class', 'Create classes'),
('class.read', 'class', 'View classes'),
('class.update', 'class', 'Update classes'),
('class.delete', 'class', 'Delete classes'),

('subject.create', 'subject', 'Create subjects'),
('subject.read', 'subject', 'View subjects'),
('subject.update', 'subject', 'Update subjects'),
('subject.delete', 'subject', 'Delete subjects'),

('section.create', 'section', 'Create sections'),
('section.read', 'section', 'View sections'),
('section.update', 'section', 'Update sections'),
('section.delete', 'section', 'Delete sections'),

-- People permissions
('teacher.create', 'teacher', 'Create teacher records'),
('teacher.read', 'teacher', 'View teacher records'),
('teacher.update', 'teacher', 'Update teacher records'),
('teacher.delete', 'teacher', 'Delete teacher records'),

('student.create', 'student', 'Create student records'),
('student.read', 'student', 'View student records'),
('student.update', 'student', 'Update student records'),
('student.delete', 'student', 'Delete student records'),

('parent.create', 'parent', 'Create parent records'),
('parent.read', 'parent', 'View parent records'),
('parent.update', 'parent', 'Update parent records'),
('parent.delete', 'parent', 'Delete parent records'),

-- Operations permissions
('schedule.create', 'schedule', 'Create schedules'),
('schedule.read', 'schedule', 'View schedules'),
('schedule.update', 'schedule', 'Update schedules'),
('schedule.delete', 'schedule', 'Delete schedules'),

('exam.create', 'exam', 'Create exams'),
('exam.read', 'exam', 'View exams'),
('exam.update', 'exam', 'Update exams'),
('exam.delete', 'exam', 'Delete exams'),

('mark.create', 'mark', 'Create marks'),
('mark.read', 'mark', 'View marks'),
('mark.update', 'mark', 'Update marks'),
('mark.delete', 'mark', 'Delete marks'),

('attendance.create', 'attendance', 'Record attendance'),
('attendance.read', 'attendance', 'View attendance'),
('attendance.update', 'attendance', 'Update attendance'),
('attendance.delete', 'attendance', 'Delete attendance')
ON CONFLICT (permission_code) DO NOTHING;

-- ============================================
-- 22. ROLE PERMISSIONS (Assign permissions to roles)
-- ============================================
-- Super Admin has all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- School Admin has most permissions except delete
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE permission_code NOT LIKE '%.delete';

-- Teacher permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions 
WHERE module_name IN ('schedule', 'exam', 'mark', 'attendance', 'student')
AND (permission_code LIKE '%.read' OR permission_code LIKE '%.create' OR permission_code LIKE '%.update');

-- Student permissions (read only)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, id FROM permissions 
WHERE permission_code LIKE '%.read'
AND module_name IN ('schedule', 'exam', 'mark', 'attendance');

-- Parent permissions (read only)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 5, id FROM permissions 
WHERE permission_code LIKE '%.read'
AND module_name IN ('student', 'mark', 'attendance');

-- ============================================
-- 23. USER ROLES (Assign roles to users)
-- ============================================
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),  -- Admin user = Super Admin
(2, 3),  -- Budi Santoso = Teacher
(3, 3),  -- Siti Aminah = Teacher
(4, 3),  -- Agus Pratama = Teacher
(5, 3),  -- Ani Wijaya = Teacher
(6, 3),  -- Joko Susilo = Teacher
(7, 4),  -- Aditya Perkasa = Student
(8, 4),  -- Bunga Citra = Student
(9, 4),  -- Chandra Wijaya = Student
(10, 4), -- Dewi Sartika = Student
(11, 4), -- Eko Prasetyo = Student
(12, 4), -- Fani Rahmawati = Student
(13, 4), -- Guntur Prabowo = Student
(14, 4), -- Hana Pertiwi = Student
(15, 4), -- Indra Lesmana = Student
(16, 4), -- Juli Anissa = Student
(17, 5), -- Bapak Heru = Parent
(18, 5), -- Ibu Citra = Parent
(19, 5), -- Bapak Wijaya = Parent
(20, 6), -- Dr. Ahmad Dahlan = Staff
(21, 6), -- Sri Wahyuni = Staff
(22, 6), -- Rudi Hartono = Staff
(23, 6), -- Dewi Lestari = Staff
(24, 6); -- Hendra Kurniawan = Staff

-- ============================================
-- 24. STAFF DATA
-- Insert staff records (Table created in migration 011)
-- ============================================

-- Insert staff data
INSERT INTO staff (user_id, school_id, employee_number, full_name, date_of_birth, gender, phone, email, address, position, department, join_date, status) VALUES
(20, 1, 'STF-001', 'Dr. H. Ahmad Dahlan, M.Pd', '1975-08-17', 'male', '081234567895', 'kepala.sekolah@sman1jakarta.sch.id', 'Jakarta', 'Kepala Sekolah', 'Management', '2000-01-01', 'active'),
(21, 1, 'STF-002', 'Sri Wahyuni, S.Admin', '1988-04-12', 'female', '081234567896', 'admin.tu@sman1jakarta.sch.id', 'Jakarta', 'Admin TU', 'Administration', '2015-06-01', 'active'),
(22, 1, 'STF-003', 'Rudi Hartono', '1995-09-25', 'male', '081234567897', 'operator@sman1jakarta.sch.id', 'Jakarta', 'Operator Sekolah', 'IT', '2020-07-15', 'active'),
(23, 1, 'STF-004', 'Dewi Lestari, S.IP', '1990-12-03', 'female', '081234567898', 'pustakawan@sman1jakarta.sch.id', 'Jakarta', 'Pustakawan', 'Library', '2018-08-01', 'active'),
(24, 1, 'STF-005', 'Hendra Kurniawan, S.Pd', '1987-06-30', 'male', '081234567899', 'konselor@sman1jakarta.sch.id', 'Jakarta', 'Konselor', 'Counseling', '2016-01-10', 'active')
ON CONFLICT (school_id, employee_number) DO NOTHING;

-- Up Migration: Seed Attendance Test Data and Enhance Schema
-- This migration adds subject_id to student_attendance and provides fresh seed data for testing

-- 1. Enhance student_attendance with subject_id
ALTER TABLE student_attendance 
    ADD COLUMN IF NOT EXISTS subject_id BIGINT REFERENCES subjects(id) ON DELETE SET NULL;

-- 2. Update unique constraint to allow multiple attendance records per day for different subjects
-- First drop the old constraint if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'student_attendance_student_id_section_id_attendance_date_key') THEN
        ALTER TABLE student_attendance DROP CONSTRAINT student_attendance_student_id_section_id_attendance_date_key;
    END IF;
END $$;

-- Add new constraint including subject_id (nullable for daily attendance)
-- We use a unique index with NULL handling or just include it
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_attendance_unique_subject 
    ON student_attendance (student_id, section_id, COALESCE(subject_id, 0), attendance_date);

-- 3. Fresh Seed Data for 2024/2025 Session
INSERT INTO academic_sessions (school_id, name, start_date, end_date, is_active)
VALUES (1, '2024/2025', '2024-07-01', '2025-06-30', true)
ON CONFLICT DO NOTHING;

-- 4. Ensure Subjects exist
INSERT INTO subjects (school_id, name, code) VALUES
(1, 'Matematika Terapan', 'MTK-2024'),
(1, 'Fisika Lanjutan', 'FIS-2024'),
(1, 'Bahasa Inggris Business', 'ENG-2024')
ON CONFLICT (school_id, code) DO NOTHING;

-- 5. Link more students to sections for testing
DO $$
DECLARE
    v_session_id BIGINT;
    v_section_id BIGINT;
    v_student_start_id BIGINT;
BEGIN
    SELECT id INTO v_session_id FROM academic_sessions WHERE name = '2024/2025' LIMIT 1;
    SELECT id INTO v_section_id FROM sections WHERE code = 'X-A' LIMIT 1;
    
    -- Ensure we have students (Aditya, Bunga, etc. are already there from 012)
    -- Just ensure they are enrolled in the new session and section if needed
    -- For testing, we'll just use the existing student-section links but ensure they are active
    UPDATE student_sections SET academic_session_id = v_session_id WHERE section_id = v_section_id;
END $$;
