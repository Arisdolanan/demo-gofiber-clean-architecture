-- Down Migration: Remove Test Data
-- This migration removes all seed data (in reverse order of dependencies)

-- ============================================
-- 1. USER ROLES
-- ============================================
DELETE FROM user_roles;

-- ============================================
-- 2. ROLE PERMISSIONS
-- ============================================
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE name IN ('super_admin', 'school_admin', 'teacher', 'student', 'parent'));

-- ============================================
-- 3. PERMISSIONS
-- ============================================
DELETE FROM permissions WHERE module_name IN ('academic_session', 'class', 'subject', 'section', 'teacher', 'student', 'parent', 'schedule', 'exam', 'mark', 'attendance');

-- ============================================
-- 4. ROLES
-- ============================================
DELETE FROM roles WHERE name IN ('super_admin', 'school_admin', 'teacher', 'student', 'parent');

-- ============================================
-- 5. NOTIFICATIONS
-- ============================================
DELETE FROM notifications WHERE school_id = 1;

-- ============================================
-- 6. STUDENT ATTENDANCE
-- ============================================
DELETE FROM student_attendance WHERE academic_session_id = 1;

-- ============================================
-- 7. TEACHER ATTENDANCE
-- ============================================
DELETE FROM teacher_attendance WHERE attendance_date >= '2024-01-01';

-- ============================================
-- 8. EXAM MARKS
-- ============================================
DELETE FROM exam_marks WHERE exam_id IN (SELECT id FROM exams WHERE academic_session_id = 1);

-- ============================================
-- 9. EXAMS
-- ============================================
DELETE FROM exams WHERE academic_session_id = 1;

-- ============================================
-- 10. SCHEDULES
-- ============================================
DELETE FROM schedules WHERE academic_session_id = 1;

-- ============================================
-- 11. STUDENT PARENTS
-- ============================================
DELETE FROM student_parents WHERE student_id IN (SELECT id FROM students WHERE school_id = 1);

-- ============================================
-- 12. PARENTS
-- ============================================
DELETE FROM parents WHERE school_id = 1;

-- ============================================
-- 13. STUDENT SECTIONS
-- ============================================
DELETE FROM student_sections WHERE academic_session_id = 1;

-- ============================================
-- 14. STUDENTS
-- ============================================
DELETE FROM students WHERE school_id = 1;

-- ============================================
-- 15. TEACHERS
-- ============================================
DELETE FROM teachers WHERE school_id = 1;

-- ============================================
-- 16. SECTIONS
-- ============================================
DELETE FROM sections WHERE academic_session_id = 1;

-- ============================================
-- 17. SUBJECTS
-- ============================================
DELETE FROM subjects WHERE school_id = 1;

-- ============================================
-- 18. CLASSES
-- ============================================
DELETE FROM classes WHERE school_id = 1;

-- ============================================
-- 19. ACADEMIC SESSIONS
-- ============================================
DELETE FROM academic_sessions WHERE school_id = 1;

-- ============================================
-- 20. USERS (Admin, Teachers, Students, Parents)
-- ============================================
DELETE FROM users WHERE user_type IN ('teacher', 'student', 'parent');

-- ============================================
-- 21. SCHOOL LICENSES
-- ============================================
DELETE FROM school_licenses WHERE school_id = 1;

-- ============================================
-- 22. APP PACKAGES
-- ============================================
DELETE FROM app_packages WHERE code IN ('PKG-BASIC', 'PKG-PRO', 'PKG-ENTERPRISE');

-- ============================================
-- 23. SCHOOL - DO NOT DELETE!
-- Keep school data to allow re-migration
-- ============================================
-- DELETE FROM schools WHERE code = 'SCH001';
