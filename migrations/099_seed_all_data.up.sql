-- ============================================
-- 099. MEGA SEED TEST DATA (MULTI-TENANCY)
-- ============================================

-- 0. Make user_id nullable in parents table and add missing columns
ALTER TABLE parents ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE parents ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE parents ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE parents ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE parents ADD COLUMN IF NOT EXISTS occupation VARCHAR(100);

-- Add is_primary column to student_parents if not exists
ALTER TABLE student_parents ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT false;

-- 1. SCHOOLS
INSERT INTO schools (id, code, name, npsn, school_level, email, status) VALUES
(1, 'SCH001', 'SMA Negeri 1 Jakarta', '12345678', 'SMA', 'info@sman1jkt.sch.id', 'active'),
(2, 'SCH002', 'SMA Negeri 2 Jakarta', '87654321', 'SMA', 'info@sman2jkt.sch.id', 'active')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, npsn = EXCLUDED.npsn;

-- 2. APP PACKAGES
INSERT INTO app_packages (code, name, price_monthly, price_yearly, max_students, max_teachers) VALUES
('PKG-BASIC', 'Basic Package', 500000, 5000000, 500, 50),
('PKG-PRO', 'Pro Package', 1000000, 10000000, 2000, 200)
ON CONFLICT (code) DO NOTHING;

-- 3. SCHOOL LICENSES
INSERT INTO school_licenses (school_id, app_package_id, license_key, start_date, end_date, status) VALUES
(1, 1, 'LIC-SCH001-BASIC', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 year', 'active'),
(2, 2, 'LIC-SCH002-PRO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 year', 'active')
ON CONFLICT (license_key) DO NOTHING;

-- 4. PERMISSIONS (System Global & Granular CRUD)
TRUNCATE TABLE permissions RESTART IDENTITY CASCADE;

INSERT INTO permissions (permission_code, module_name, permission_name, description) VALUES
-- AKADEMIK
('akademik.kelas.read', 'Akademik', 'Lihat Kelas', 'Melihat daftar kelas'),
('akademik.kelas.create', 'Akademik', 'Tambah Kelas', 'Membuat kelas baru'),
('akademik.kelas.update', 'Akademik', 'Edit Kelas', 'Mengubah data kelas'),
('akademik.kelas.delete', 'Akademik', 'Hapus Kelas', 'Menghapus data kelas'),

('akademik.mapel.read', 'Akademik', 'Lihat Mata Pelajaran', 'Melihat daftar mata pelajaran'),
('akademik.mapel.create', 'Akademik', 'Tambah Mata Pelajaran', 'Menambah mata pelajaran baru'),
('akademik.mapel.update', 'Akademik', 'Edit Mata Pelajaran', 'Mengubah mata pelajaran'),

('akademik.jadwal.read', 'Akademik', 'Lihat Jadwal', 'Melihat jadwal pelajaran'),
('akademik.silabus.read', 'Akademik', 'Lihat Silabus', 'Melihat kurikulum dan silabus'),
('akademik.ujian.read', 'Akademik', 'Lihat Ujian', 'Melihat daftar ujian'),
('akademik.penilaian.read', 'Akademik', 'Lihat Penilaian', 'Melihat data tugas dan penilaian'),
('akademik.nilai_rapor.read', 'Akademik', 'Lihat Nilai/Rapor', 'Melihat nilai akhir dan rapor'),
('akademik.kenaikan_kelas.read', 'Akademik', 'Lihat Kenaikan', 'Manajemen kenaikan kelas'),
('akademik.aktivitas.read', 'Akademik', 'Akses Aktivitas', 'Membuka menu grup aktivitas akademik'),
('akademik.grading_master.read', 'Akademik', 'Akses Menu Penilaian', 'Membuka menu grup penilaian'),
('akademik.config.read', 'Akademik', 'Konfigurasi Akademik', 'Mengatur jenis penilaian dan konfigurasi rapor'),
('akademik.jadwal.create', 'Akademik', 'Tambah Jadwal', 'Mengatur jadwal pelajaran'),

('akademik.nilai.read', 'Akademik', 'Lihat Nilai', 'Melihat nilai siswa'),
('akademik.nilai.update', 'Akademik', 'Input Nilai', 'Menginput/mengubah nilai siswa'),

-- KEUANGAN
('keuangan.spp.read', 'Keuangan', 'Lihat Pembayaran SPP', 'Melihat riwayat pembayaran'),
('keuangan.spp.create', 'Keuangan', 'Tambah Pembayaran', 'Mencatat pembayaran baru'),
('keuangan.laporan.read', 'Keuangan', 'Lihat Laporan Keuangan', 'Melihat ringkasan keuangan'),

-- MANAJEMEN USER & AKSES
('manajemen_user.read', 'Manajemen User', 'Lihat Daftar User', 'Melihat semua user sistem'),
('manajemen_user.create', 'Manajemen User', 'Tambah User', 'Mendaftarkan user baru'),
('manajemen_user.delete', 'Manajemen User', 'Hapus User', 'Menghapus user dari sistem'),

('role_permission.read', 'Role & Permission', 'Lihat Role', 'Melihat daftar jabatan/role'),
('role_permission.update', 'Role & Permission', 'Kelola Izin Role', 'Mengatur matrix permission per role'),

-- DATA MASTER
('data_guru.read', 'Data Master', 'Lihat Data Guru', 'Melihat profil guru'),
('data_guru.create', 'Data Master', 'Tambah Guru', 'Menambah data guru baru'),
('data_siswa.read', 'Data Master', 'Lihat Data Siswa', 'Melihat profil siswa'),
('data_siswa.create', 'Data Master', 'Tambah Siswa', 'Menambah data siswa baru'),

-- OPERASIONAL & ABSENSI
('operasional.read', 'Operasional', 'Akses Menu Operasional', 'Melihat menu operasional'),
('absensi_guru.read', 'Operasional', 'Lihat Absensi Guru', 'Melihat riwayat absensi guru'),
('absensi_guru.create', 'Operasional', 'Input Absensi Guru', 'Mencatat absensi guru'),
('absensi_siswa.read', 'Operasional', 'Lihat Absensi Siswa', 'Melihat riwayat absensi siswa'),
('absensi_siswa.create', 'Operasional', 'Input Absensi Siswa', 'Mencatat absensi siswa'),
('absensi_staf.read', 'Operasional', 'Lihat Absensi Staf', 'Melihat riwayat absensi staf'),
('absensi_staf.create', 'Operasional', 'Input Absensi Staf', 'Mencatat absensi staf')
ON CONFLICT (permission_code) DO NOTHING;

-- 5. ROLES
TRUNCATE TABLE roles RESTART IDENTITY CASCADE;

INSERT INTO roles (school_id, code, name, description, is_system_role) VALUES
(NULL, 'super_admin', 'Super Administrator', 'Pemilik sistem dengan akses tak terbatas ke seluruh fitur dan pengaturan multi-sekolah', true),
(1, 'school_admin', 'Administrator Sekolah', 'Manajemen operasional sekolah, pengaturan user, dan konfigurasi kurikulum', false),
(1, 'teacher', 'Guru / Tenaga Pengajar', 'Manajemen kegiatan belajar mengajar, input nilai, absensi siswa, dan materi pelajaran', false),
(1, 'staff', 'Staf Administrasi', 'Manajemen data operasional, absensi pegawai, dan administrasi perkantoran', false),
(1, 'student', 'Siswa', 'Akses ke jadwal pelajaran, materi, tugas, dan melihat hasil studi (rapor)', false),
(2, 'school_admin', 'Administrator Sekolah', 'Manajemen operasional sekolah (Sekolah 2)', false),
(2, 'teacher', 'Guru / Tenaga Pengajar', 'Manajemen kegiatan belajar mengajar (Sekolah 2)', false)
ON CONFLICT (school_id, code) DO NOTHING;

-- 6. ROLE PERMISSIONS
-- Assign all permissions to Super Admin and School Admin
INSERT INTO role_permissions (school_id, role_id, permission_id)
SELECT r.school_id, r.id, p.id FROM roles r, permissions p
WHERE r.code IN ('super_admin', 'school_admin')
ON CONFLICT DO NOTHING;

-- Teacher: Akses Akademik + Absensi Guru & Absensi Siswa
INSERT INTO role_permissions (school_id, role_id, permission_id)
SELECT r.school_id, r.id, p.id FROM roles r, permissions p
WHERE r.code = 'teacher' AND (
    (p.permission_code LIKE 'akademik.%' AND p.permission_code NOT LIKE 'akademik.config%') OR 
    p.permission_code LIKE 'absensi_guru.%' OR 
    p.permission_code LIKE 'absensi_siswa.%' OR
    p.permission_code = 'operasional.read'
)
ON CONFLICT DO NOTHING;

-- Staff: Akses Absensi Staf + Operasional
INSERT INTO role_permissions (school_id, role_id, permission_id)
SELECT r.school_id, r.id, p.id FROM roles r, permissions p
WHERE r.code = 'staff' AND (
    p.permission_code LIKE 'absensi_staf.%' OR 
    p.permission_code = 'operasional.read'
)
ON CONFLICT DO NOTHING;

-- 7. USERS
-- Note: Password for all users is "@SuperIndo1" 
-- Using bcrypt hash: $2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC
INSERT INTO users (email, password, username, user_type, school_id, email_verified_at) VALUES
('admin@system.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'superadmin', 'super_admin', ARRAY[1, 2], CURRENT_TIMESTAMP),
('admin1@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'admin1', 'admin', ARRAY[1], CURRENT_TIMESTAMP),
('admin2@school2.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'admin2', 'admin', ARRAY[2], CURRENT_TIMESTAMP),
('teacher1@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'teacher1', 'teacher', ARRAY[1], CURRENT_TIMESTAMP),
('teacher2@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'teacher2', 'teacher', ARRAY[1], CURRENT_TIMESTAMP),
('teacher3@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'teacher3', 'teacher', ARRAY[1], CURRENT_TIMESTAMP),
('teacher4@school2.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'teacher4', 'teacher', ARRAY[2], CURRENT_TIMESTAMP),
('staff1@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'staff1', 'staff', ARRAY[1], CURRENT_TIMESTAMP),
('staff2@school2.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'staff2', 'staff', ARRAY[2], CURRENT_TIMESTAMP),
('parent1@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'parent1', 'parent', ARRAY[1], CURRENT_TIMESTAMP),
('parent2@school1.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'parent2', 'parent', ARRAY[1], CURRENT_TIMESTAMP),
('parent3@school2.com', '$2a$10$lY6rsGhNtb4DkVf1ihy4hu6nowyLJu4YWWkWKNSMsqrsrSW31HepC', 'parent3', 'parent', ARRAY[2], CURRENT_TIMESTAMP)
ON CONFLICT (username) DO UPDATE SET school_id = EXCLUDED.school_id;

-- 7.1 USER ROLES (RBAC Mapping)
INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'superadmin' AND r.code = 'super_admin' AND r.school_id IS NULL
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'admin1' AND r.code = 'school_admin' AND r.school_id = 1
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'admin2' AND r.code = 'school_admin' AND r.school_id = 2
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'teacher1' AND r.code = 'teacher' AND r.school_id = 1
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'staff1' AND r.code = 'staff' AND r.school_id = 1
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (school_id, user_id, role_id)
SELECT r.school_id, u.id, r.id FROM users u, roles r 
WHERE u.username = 'parent1' AND r.code = 'parent' AND r.school_id = 1
ON CONFLICT DO NOTHING;

-- 8. ACADEMIC SESSIONS
INSERT INTO academic_sessions (id, school_id, code, name, start_date, end_date, is_active) VALUES
(1, 1, '2023-ODD-S1', '2023/2024 Ganjil', '2023-07-01', '2023-12-31', true),
(2, 2, '2023-ODD-S2', '2023/2024 Ganjil', '2023-07-01', '2023-12-31', true)
ON CONFLICT (id) DO NOTHING;

-- 9. CLASSES
INSERT INTO classes (id, school_id, code, name, level, grade_number) VALUES
(1, 1, '10-S1', 'Kelas 10', 'SMA', 10),
(2, 1, '11-S1', 'Kelas 11', 'SMA', 11),
(3, 1, '12-S1', 'Kelas 12', 'SMA', 12),
(4, 2, '10-S2', 'Kelas 10', 'SMA', 10)
ON CONFLICT (id) DO NOTHING;

-- 10. SUBJECTS
INSERT INTO subjects (id, school_id, code, name, credit_hours) VALUES
(1, 1, 'MATH-10', 'Matematika', 4),
(2, 1, 'PHYS-10', 'Fisika', 3),
(3, 1, 'BIOL-10', 'Biologi', 3),
(4, 1, 'CHEM-10', 'Kimia', 3),
(5, 1, 'IND-10', 'Bahasa Indonesia', 2),
(6, 2, 'MATH-10', 'Matematika', 4)
ON CONFLICT (id) DO NOTHING;

-- 11. SECTIONS
INSERT INTO sections (id, school_id, class_id, academic_session_id, code, name, capacity) VALUES
(1, 1, 1, 1, '10-A', '10-A', 36),
(2, 1, 2, 1, '11-A', '11-A', 36),
(3, 1, 3, 1, '12-A', '12-A', 36),
(4, 2, 4, 2, '10-A', '10-A', 36)
ON CONFLICT (id) DO NOTHING;

-- 12. TEACHERS
INSERT INTO teachers (id, user_id, school_id, employee_number, full_name, status)
SELECT 1, id, 1, 'EMP001', 'Dr. Budi Utomo', 'active' FROM users WHERE username = 'teacher1'
ON CONFLICT (id) DO NOTHING;

INSERT INTO teachers (id, user_id, school_id, employee_number, full_name, status)
SELECT 2, id, 1, 'EMP002', 'Dra. Siti Aminah', 'active' FROM users WHERE username = 'teacher2'
ON CONFLICT (id) DO NOTHING;

INSERT INTO teachers (id, user_id, school_id, employee_number, full_name, status)
SELECT 3, id, 1, 'EMP003', 'Ahmad Dani, M.Pd', 'active' FROM users WHERE username = 'teacher3'
ON CONFLICT (id) DO NOTHING;

INSERT INTO teachers (id, user_id, school_id, employee_number, full_name, status)
SELECT 4, id, 2, 'EMP004', 'Asep Saepudin', 'active' FROM users WHERE username = 'teacher4'
ON CONFLICT (id) DO NOTHING;

-- 13. STAFF (Additional Staff for School 1)
INSERT INTO staff (id, user_id, school_id, employee_number, full_name, status, position, department) VALUES
(3, NULL, 1, 'STF003', 'Hendra Wijaya', 'active', 'Finance Staff', 'Finance'),
(4, NULL, 1, 'STF004', 'Rina Susilowati', 'active', 'HR Staff', 'Human Resources'),
(5, NULL, 1, 'STF005', 'Dedi Kurniawan', 'active', 'Library Staff', 'Library'),
(6, NULL, 1, 'STF006', 'Siti Rahayu', 'active', 'Counselor', 'Student Affairs'),
(7, NULL, 1, 'STF007', 'Muhamad Faisal', 'active', 'Security', 'Security'),
(8, NULL, 1, 'STF008', 'Nurul Hidayah', 'active', 'Lab Assistant', 'Laboratory'),
(9, NULL, 1, 'STF009', 'Ahmad Rizal', 'active', 'Administration', 'Administration'),
(10, NULL, 1, 'STF010', 'Dewi Kartika', 'active', 'Registrar', 'Administration')
ON CONFLICT (id) DO NOTHING;

-- 13. STAFF (Additional Staff for School 2)
INSERT INTO staff (id, user_id, school_id, employee_number, full_name, status, position, department) VALUES
(11, NULL, 2, 'STF011', 'Tono Saputra', 'active', 'Finance Staff', 'Finance'),
(12, NULL, 2, 'STF012', 'Yanti Marlina', 'active', 'HR Staff', 'Human Resources'),
(13, NULL, 2, 'STF013', 'Baba Abdul', 'active', 'Security', 'Security'),
(14, NULL, 2, 'STF014', 'Ika Puspita', 'active', 'Library Staff', 'Library')
ON CONFLICT (id) DO NOTHING;

-- 13.1 PARENTS - Menggunakan ID 1-15 (replace ID asli yang tidak ada)
INSERT INTO parents (id, user_id, school_id, full_name, phone, email, occupation, address) VALUES
(1, NULL, 1, 'Bapak Suherman', '08123456789', 'suherman@email.com', 'Wiraswasta', 'Jl. Sudirman No. 5, Jakarta'),
(2, NULL, 1, 'Ibu Sumiati', '08129876543', 'sumiati@email.com', 'Ibu Rumah Tangga', 'Jl. Thamrin No. 10, Jakarta'),
(3, NULL, 2, 'Bapak Jajang', '08131122334', 'jajang@email.com', 'Petani', 'Jl. Braga No. 3, Bandung'),
(4, NULL, 1, 'Bapak H. Ahmad Yani', '081234567810', 'ahmadyani@email.com', 'Wiraswasta', 'Jl. Merdeka No. 10, Jakarta'),
(5, NULL, 1, 'Ibu Dra. Hj. Aminah', '081234567811', 'aminah@email.com', 'Guru', 'Jl. Pendidikan No. 5, Jakarta'),
(6, NULL, 1, 'Bapak Ir. Budi Santoso', '081234567812', 'budisantoso@email.com', 'Wiraswasta', 'Jl. Sudirman No. 15, Jakarta'),
(7, NULL, 1, 'Ibu Hj. Siti Rohanah', '081234567813', 'sitirohanah@email.com', 'Ibu Rumah Tangga', 'Jl. Dago No. 20, Bandung'),
(8, NULL, 1, 'Bapak Dr. Anwar', '081234567814', 'anwar@email.com', 'Dokter', 'Jl. Asia Afrika No. 8, Jakarta'),
(9, NULL, 1, 'Ibu Dr. Rina', '081234567815', 'rina@email.com', 'Dokter', 'Jl. Thamrin No. 12, Jakarta'),
(10, NULL, 1, 'Bapak M. Yusuf', '081234567816', 'yusuf@email.com', 'Pegawai Negeri', 'Jl. Gatot Subroto No. 25, Jakarta'),
(11, NULL, 1, 'Ibu Sri Wahyuni', '081234567817', 'sriwahyuni@email.com', 'Guru', 'Jl. Diponegoro No. 30, Jakarta'),
(12, NULL, 1, 'Bapak Fahmi Idris', '081234567818', 'fahmi@email.com', 'Wiraswasta', 'Jl. A.Yani No. 18, Jakarta'),
(13, NULL, 1, 'Ibu Lilis Sumarni', '081234567819', 'lilis@email.com', 'Ibu Rumah Tangga', 'Jl. Ciledug No. 22, Jakarta'),
(14, NULL, 2, 'Bapak Ujang Sulaeman', '081234567820', 'ujang@email.com', 'Wiraswasta', 'Jl. Braga No. 5, Bandung'),
(15, NULL, 2, 'Ibu Cucu Cahyati', '081234567821', 'cucu@email.com', 'Ibu Rumah Tangga', 'Jl. Asia Afrika No. 10, Bandung')
ON CONFLICT (id) DO NOTHING;

-- 14. STUDENTS
INSERT INTO students (id, user_id, school_id, student_number, full_name, status, nis, nisn, date_of_birth, gender) VALUES
(1, NULL, 1, 'STU001', 'Siti Aminah', 'active', '1001', '0012345671', '2008-05-15', 'female'),
(2, NULL, 1, 'STU002', 'Budi Santoso', 'active', '1002', '0012345672', '2008-08-20', 'male'),
(3, NULL, 1, 'STU003', 'Ani Wijaya', 'active', '1101', '0012345673', '2007-03-10', 'female'),
(4, NULL, 1, 'STU004', 'Joko Susilo', 'active', '1102', '0012345674', '2007-11-25', 'male'),
(5, NULL, 1, 'STU005', 'Rina Kartika', 'active', '1201', '0012345675', '2006-01-05', 'female'),
(6, NULL, 2, 'STU006', 'Ujang Memet', 'active', '2001', '0012345676', '2008-02-14', 'male'),
(7, NULL, 2, 'STU007', 'Neneng Geulis', 'active', '2002', '0012345677', '2008-09-09', 'female'),
(8, NULL, 1, 'STU008', 'Deden Hamzah', 'active', '1003', '0012345678', '2008-06-20', 'male'),
(9, NULL, 1, 'STU009', 'Euis Marlina', 'active', '1004', '0012345679', '2008-07-21', 'female'),
(10, NULL, 1, 'STU010', 'Fajar Sidik', 'active', '1103', '0012345680', '2007-12-12', 'male'),
(11, NULL, 1, 'STU011', 'Gita Gutawa', 'active', '1202', '0012345681', '2006-02-14', 'female'),
(12, NULL, 2, 'STU012', 'Cecep Gorbacep', 'active', '2003', '0012345682', '2008-04-01', 'male')
ON CONFLICT (id) DO NOTHING;

-- 14.1 STUDENT PARENT LINK
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 1, 1, 'father', true),
(1, 2, 2, 'mother', false),
(2, 6, 3, 'father', true)
ON CONFLICT DO NOTHING;

-- 14.2 ADDITIONAL STUDENT PARENT LINKS (must come AFTER students)
-- School 1 Students: 1(Siti Aminah), 2(Budi Santoso), 3(Ani Wijaya), 4(Joko Susilo), 5(Rina Kartika), 8(Deden Hamzah), 9(Euis Marlina), 10(Fajar Sidik), 11(Gita Gutawa)
-- School 2 Students: 6(Ujang Memet), 7(Neneng Geulis), 12(Cecep Gorbacep)

-- Student 3 (Ani Wijaya) -> Parent 4 (father) & Parent 5 (mother)
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 3, 4, 'father', true),
(1, 3, 5, 'mother', false)
ON CONFLICT DO NOTHING;

-- Student 4 (Joko Susilo) -> Parent 6 (father) & Parent 7 (mother)
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 4, 6, 'father', true),
(1, 4, 7, 'mother', false)
ON CONFLICT DO NOTHING;

-- Student 5 (Rina Kartika) -> Parent 8 (father) & Parent 9 (mother)
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 5, 8, 'father', true),
(1, 5, 9, 'mother', false)
ON CONFLICT DO NOTHING;

-- Student 8 (Deden Hamzah) -> Parent 10 (father) & Parent 11 (mother)
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 8, 10, 'father', true),
(1, 8, 11, 'mother', false)
ON CONFLICT DO NOTHING;

-- Student 9 (Euis Marlina) -> Parent 12 (father) & Parent 13 (mother)
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 9, 12, 'father', true),
(1, 9, 13, 'mother', false)
ON CONFLICT DO NOTHING;

-- Student 10 (Fajar Sidik) -> linked to existing parent 1
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 10, 1, 'father', true)
ON CONFLICT DO NOTHING;

-- Student 11 (Gita Gutawa) -> linked to existing parent 2
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(1, 11, 2, 'father', true)
ON CONFLICT DO NOTHING;

-- School 2 - Student 7 (Neneng Geulis) -> Parent 14 & 15
INSERT INTO student_parents (school_id, student_id, parent_id, relationship, is_primary) VALUES
(2, 7, 14, 'father', true),
(2, 7, 15, 'mother', false)
ON CONFLICT DO NOTHING;

-- 15. STUDENT SECTIONS (Enrollment)
INSERT INTO student_sections (school_id, student_id, section_id, academic_session_id, roll_number, enrollment_date) VALUES
(1, 1, 1, 1, '01', '2023-07-10'),
(1, 2, 1, 1, '02', '2023-07-10'),
(1, 8, 1, 1, '03', '2023-07-10'),
(1, 9, 1, 1, '04', '2023-07-10'),
(1, 3, 2, 1, '01', '2023-07-10'),
(1, 4, 2, 1, '02', '2023-07-10'),
(1, 10, 2, 1, '03', '2023-07-10'),
(1, 5, 3, 1, '01', '2023-07-10'),
(1, 11, 3, 1, '02', '2023-07-10'),
(2, 6, 4, 2, '01', '2023-07-10'),
(2, 7, 4, 2, '02', '2023-07-10'),
(2, 12, 4, 2, '03', '2023-07-10')
ON CONFLICT DO NOTHING;

-- 16. SCHEDULES
INSERT INTO schedules (school_id, section_id, subject_id, teacher_id, academic_session_id, day_of_week, start_time, end_time, room_number) VALUES
(1, 1, 1, 1, 1, 1, '07:30:00', '09:00:00', 'R101'),
(1, 1, 2, 2, 1, 1, '09:30:00', '11:00:00', 'R101'),
(1, 2, 1, 1, 1, 2, '07:30:00', '09:00:00', 'R201'),
(1, 3, 1, 3, 1, 3, '07:30:00', '09:00:00', 'R301'),
(2, 4, 6, 4, 2, 1, '07:30:00', '09:00:00', 'R101')
ON CONFLICT DO NOTHING;

-- 17. EXAMS
INSERT INTO exams (id, school_id, section_id, subject_id, academic_session_id, exam_type, exam_date, max_score, title) VALUES
(1, 1, 1, 1, 1, 'Midterm', '2023-10-15', 100, 'UTS Matematika Ganjil'),
(2, 1, 1, 2, 1, 'Midterm', '2023-10-16', 100, 'UTS Fisika Ganjil'),
(3, 1, 2, 1, 1, 'Midterm', '2023-10-15', 100, 'UTS Matematika Ganjil'),
(4, 1, 3, 1, 1, 'Midterm', '2023-10-15', 100, 'UTS Matematika Ganjil'),
(5, 2, 4, 6, 2, 'Midterm', '2023-10-15', 100, 'UTS Matematika Ganjil')
ON CONFLICT (id) DO NOTHING;

-- 18. EXAM MARKS
INSERT INTO exam_marks (school_id, exam_id, student_id, score, notes) VALUES
(1, 1, 1, 85.5, 'Bagus'),
(1, 1, 2, 70.0, 'Cukup'),
(1, 1, 8, 92.0, 'Sangat Memuaskan'),
(1, 1, 9, 78.5, 'Baik'),
(1, 2, 1, 90.0, 'Sangat Bagus'),
(1, 3, 3, 88.0, 'Mantap'),
(1, 3, 4, 65.0, 'Perlu Belajar'),
(1, 3, 10, 82.0, 'Bagus'),
(1, 4, 5, 95.0, 'Istimewa'),
(1, 4, 11, 89.0, 'Bagus Sekali'),
(2, 5, 6, 80.0, 'Baik'),
(2, 5, 7, 75.0, 'Cukup'),
(2, 5, 12, 85.0, 'Bagus')
ON CONFLICT (exam_id, student_id) DO NOTHING;

-- 19. ATTENDANCE
-- Student Attendance
INSERT INTO student_attendance (school_id, student_id, section_id, academic_session_id, attendance_date, status) VALUES
(1, 1, 1, 1, CURRENT_DATE, 'Present'),
(1, 2, 1, 1, CURRENT_DATE, 'Present'),
(1, 8, 1, 1, CURRENT_DATE, 'Present'),
(1, 9, 1, 1, CURRENT_DATE, 'Present'),
(1, 3, 2, 1, CURRENT_DATE, 'Late'),
(1, 4, 2, 1, CURRENT_DATE, 'Present'),
(1, 10, 2, 1, CURRENT_DATE, 'Sick'),
(1, 5, 3, 1, CURRENT_DATE, 'Sick'),
(1, 11, 3, 1, CURRENT_DATE, 'Present'),
(2, 6, 4, 2, CURRENT_DATE, 'Present'),
(2, 7, 4, 2, CURRENT_DATE, 'Present'),
(2, 12, 4, 2, CURRENT_DATE, 'Absent')
ON CONFLICT DO NOTHING;

-- Teacher Attendance
INSERT INTO teacher_attendance (school_id, teacher_id, attendance_date, status) VALUES
(1, 1, CURRENT_DATE, 'Present'),
(1, 2, CURRENT_DATE, 'Present'),
(1, 3, CURRENT_DATE, 'Present'),
(2, 4, CURRENT_DATE, 'Present')
ON CONFLICT DO NOTHING;

-- Staff Attendance
INSERT INTO staff_attendance (school_id, employee_id, attendance_date, status, check_in_time, check_out_time) VALUES
(1, 3, CURRENT_DATE, 'present', '07:30:00', '16:00:00'),
(2, 11, CURRENT_DATE, 'present', '07:35:00', '16:05:00')
ON CONFLICT DO NOTHING;

-- 20. CLASS SUBJECTS
INSERT INTO class_subjects (school_id, class_id, subject_id, academic_session_id) VALUES
(1, 1, 1, 1), (1, 1, 2, 1), (1, 1, 3, 1), (1, 1, 4, 1), (1, 1, 5, 1),
(1, 2, 1, 1), (1, 2, 2, 1),
(1, 3, 1, 1),
(2, 4, 6, 2)
ON CONFLICT DO NOTHING;

-- 21. TEACHER SUBJECTS
INSERT INTO teacher_subjects (school_id, teacher_id, section_id, subject_id, academic_session_id) VALUES
(1, 1, 1, 1, 1),
(1, 2, 1, 2, 1),
(1, 1, 2, 1, 1),
(1, 3, 3, 1, 1),
(2, 4, 4, 6, 2)
ON CONFLICT DO NOTHING;

-- 22. SETTINGS
INSERT INTO settings (school_id, setting_key, setting_value, group_name, description) VALUES
(1, 'school_name', 'SMA Negeri 1 Jakarta', 'umum', 'Nama sekolah'),
(1, 'school_address', 'Jl. Budi Utomo No.7, Jakarta', 'umum', 'Alamat sekolah'),
(1, 'school_email', 'info@sman1jkt.sch.id', 'umum', 'Email resmi sekolah'),
(1, 'school_phone', '021-1234567', 'umum', 'Nomor telepon sekolah'),
(1, 'academic_year', '2023/2024', 'akademik', 'Tahun ajaran aktif'),
(1, 'principal_name', 'H. Akhmad Fauzi, M.Pd', 'akademik', 'Nama Kepala Sekolah'),
(1, 'principal_nip', '197501012000011001', 'akademik', 'NIP Kepala Sekolah'),
(1, 'theme_color', '#1a73e8', 'tampilan', 'Warna utama aplikasi'),
(1, 'sidebar_style', 'expanded', 'tampilan', 'Gaya sidebar aplikasi'),
(1, 'backup_schedule', 'daily', 'backup', 'Jadwal backup otomatis'),
(1, 'backup_retention_days', '30', 'backup', 'Lama penyimpanan backup'),
(1, 'integrasi_whatsapp', '{"endpoint": "https://wa.sman1jkt.sch.id/api", "token": "WA-KEY-12345", "isActive": true}', 'integrasi', 'WhatsApp Gateway'),
(1, 'integrasi_sms', '{"endpoint": "https://api.zenziva.net/v1", "token": "SMS-KEY-67890", "isActive": false}', 'integrasi', 'SMS Gateway'),
(1, 'integrasi_smtp', '{"endpoint": "smtp.gmail.com", "token": "smtp_pass123", "isActive": true}', 'integrasi', 'Email Server SMTP'),
(1, 'integrasi_midtrans', '{"endpoint": "https://api.midtrans.com", "token": "MT-KEY-abc123", "isActive": true}', 'integrasi', 'Midtrans Payment'),
(1, 'integrasi_xendit', '{"endpoint": "https://api.xendit.co/v2", "token": "XN-KEY-xyz789", "isActive": false}', 'integrasi', 'Xendit Payment'),
(1, 'integrasi_telegram', '{"endpoint": "https://api.telegram.org", "token": "TG-BOT-TOKEN", "isActive": false}', 'integrasi', 'Telegram Bot'),
(1, 'integrasi_zoom', '{"endpoint": "https://api.zoom.us/v2", "token": "ZOOM-JWT-TOKEN", "isActive": false}', 'integrasi', 'Zoom Meeting'),
(1, 'integrasi_google_calendar', '{"endpoint": "https://www.googleapis.com/calendar/v3", "token": "GOOGLE-API-KEY", "isActive": false}', 'integrasi', 'Google Calendar'),
(1, 'system_api_key', 'ak_live_default_key', 'integrasi', 'System API Key'),
(1, 'subscription_token', 'sub_active_12345', 'integrasi', 'Subscription Token'),
(1, 'dev_tools_url', '/swagger', 'integrasi', 'Dev Tools URL'),
(2, 'school_name', 'SMA Negeri 2 Jakarta', 'umum', 'Nama sekolah'),
(2, 'school_address', 'Jl. Gajah Mada No.1, Jakarta', 'umum', 'Alamat sekolah'),
(2, 'academic_year', '2023/2024', 'akademik', 'Tahun ajaran aktif')
ON CONFLICT (school_id, setting_key) DO NOTHING;

-- 23. FILES
INSERT INTO files (id, school_id, user_id, filename, original_name, path, mime_type, size, is_public) VALUES
(1, 1, 1, 'logo_sman1.png', 'logo.png', 'storage/school_1/logo.png', 'image/png', 51200, true),
(2, 1, 1, 'student_import_template.xlsx', 'template.xlsx', 'storage/school_1/templates/student.xlsx', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 10240, false),
(3, 2, 3, 'logo_sman2.png', 'logo.png', 'storage/school_2/logo.png', 'image/png', 48000, true)
ON CONFLICT (id) DO NOTHING;

-- 24. BACKUPS
INSERT INTO backups (id, school_id, filename, storage_path, size_bytes, status) VALUES
(1, 1, 'backup_20231001.sql', 'storage/school_1/backups/backup_20231001.sql', 1048576, 'success'),
(2, 1, 'backup_20231015.sql', 'storage/school_1/backups/backup_20231015.sql', 1052000, 'success'),
(3, 2, 'backup_20231001.sql', 'storage/school_2/backups/backup_20231001.sql', 1048576, 'success')
ON CONFLICT (id) DO NOTHING;

-- 25. ACTIVITY LOGS
INSERT INTO activity_logs (school_id, user_id, action, module, description) VALUES
(1, 1, 'LOGIN', 'Auth', 'Super Admin logged in from web'),
(1, 2, 'CREATE_STUDENT', 'Student', 'Admin created new student record: Siti Aminah'),
(1, 2, 'CREATE_STUDENT', 'Student', 'Admin created new student record: Budi Santoso'),
(1, 4, 'MARK_ATTENDANCE', 'Attendance', 'Teacher marked attendance for 10-A'),
(2, 3, 'UPDATE_SETTING', 'Setting', 'School Admin updated school address')
ON CONFLICT DO NOTHING;

-- 26. NOTIFICATIONS
-- First ensure deleted_at column exists
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- Insert notification seed data for School 1 (user_id 2 = admin1)
INSERT INTO notifications (user_id, school_id, title, message, reference_type, reference_id, is_read, created_at) VALUES
(2, 1, 'Selamat Datang', 'Selamat datang di Sistem Akademik SMA Negeri 1 Jakarta', 'system', NULL, false, CURRENT_TIMESTAMP - INTERVAL '1 day'),
(2, 1, 'Update Sistem', 'Sistem akan undergo pemeliharaan pada malam hari', 'system', NULL, false, CURRENT_TIMESTAMP - INTERVAL '12 hours'),
(2, 1, 'Absensi Hari Ini', 'Absensi hari ini telah dicatat untuk 45 siswa', 'attendance', NULL, true, CURRENT_TIMESTAMP - INTERVAL '5 hours'),
(2, 1, 'Raport Siap', 'Raport semester ganjil telah siap untuk didistribusikan', 'academic', NULL, false, CURRENT_TIMESTAMP - INTERVAL '2 hours'),
(2, 1, 'Pembayaran SPP', 'Pembayaran SPP bulan November telah lunas untuk 80% siswa', 'payment', NULL, true, CURRENT_TIMESTAMP - INTERVAL '1 day')
ON CONFLICT DO NOTHING;

-- Insert notification seed data for School 2 (user_id 3 = admin2)
INSERT INTO notifications (user_id, school_id, title, message, reference_type, reference_id, is_read, created_at) VALUES
(3, 2, 'Selamat Datang', 'Selamat datang di Sistem Akademik SMA Negeri 2 Jakarta', 'system', NULL, false, CURRENT_TIMESTAMP - INTERVAL '1 day'),
(3, 2, 'Jadwal Ujian', 'Jadwal UTS semester ganjil telah dirilis', 'academic', NULL, false, CURRENT_TIMESTAMP - INTERVAL '6 hours'),
(3, 2, 'Backup Selesai', 'Backup database berhasil dilakukan', 'backup', NULL, true, CURRENT_TIMESTAMP - INTERVAL '1 day')
ON CONFLICT DO NOTHING;

-- Insert notification seed data for Super Admin (user_id 1)
INSERT INTO notifications (user_id, school_id, title, message, reference_type, reference_id, is_read, created_at) VALUES
(1, NULL, 'Notifikasi Multi-Sekolah', 'Ada 2 sekolah yang memerlukan konfirmasi', 'system', NULL, false, CURRENT_TIMESTAMP - INTERVAL '3 hours'),
(1, NULL, 'Laporan Bulanan', 'Laporan bulanan Oktober telah tersedia', 'report', NULL, false, CURRENT_TIMESTAMP - INTERVAL '1 day')
ON CONFLICT DO NOTHING;


-- SEED DATA FROM OTHER MIGRATIONS

INSERT INTO backups (school_id, filename, storage_path, size_bytes, status, created_at) VALUES (1, 'backup_system_initial.sql', 'storage/backup/backup_system_initial.sql', 1048576, 'success', NOW() - INTERVAL '2 days'), (1, 'backup_system_weekly.sql', 'storage/backup/backup_system_weekly.sql', 5242880, 'success', NOW() - INTERVAL '7 days') ON CONFLICT DO NOTHING;

INSERT INTO integration_definitions (school_id, code, name, provider, category, description, is_system) VALUES (1, 'whatsapp', 'WhatsApp Gateway', 'Wablas', 'messaging', 'WhatsApp messaging service for notifications', true), (1, 'sms', 'SMS Gateway', 'Zenziva', 'messaging', 'SMS notification service', true), (1, 'smtp', 'Email Server (SMTP)', 'Gmail', 'messaging', 'Email delivery via SMTP', true), (1, 'midtrans', 'Midtrans', 'Midtrans', 'payment', 'Payment gateway integration', true), (1, 'xendit', 'Xendit', 'Xendit', 'payment', 'Alternative payment gateway', true), (1, 'telegram', 'Telegram Bot', 'Telegram', 'messaging', 'Telegram bot for notifications', true), (1, 'zoom', 'Zoom Meeting', 'Zoom', 'meeting', 'Video conferencing integration', true), (1, 'google_calendar', 'Google Calendar', 'Google', 'calendar', 'Calendar sync with Google Calendar', true) ON CONFLICT (school_id, code) DO NOTHING;

INSERT INTO integration_definitions (school_id, code, name, provider, category, description, is_system) VALUES (2, 'whatsapp', 'WhatsApp Gateway', 'Wablas', 'messaging', 'WhatsApp messaging service for notifications', true), (2, 'sms', 'SMS Gateway', 'Zenziva', 'messaging', 'SMS notification service', true), (2, 'smtp', 'Email Server (SMTP)', 'Gmail', 'messaging', 'Email delivery via SMTP', true), (2, 'midtrans', 'Midtrans', 'Midtrans', 'payment', 'Payment gateway integration', true), (2, 'xendit', 'Xendit', 'Xendit', 'payment', 'Alternative payment gateway', true), (2, 'telegram', 'Telegram Bot', 'Telegram', 'messaging', 'Telegram bot for notifications', true), (2, 'zoom', 'Zoom Meeting', 'Zoom', 'meeting', 'Video conferencing integration', true), (2, 'google_calendar', 'Google Calendar', 'Google', 'calendar', 'Calendar sync with Google Calendar', true) ON CONFLICT (school_id, code) DO NOTHING;
