Entity Relationship Diagram (ERD)
Web-Based Academic Information System with SaaS Capability

1. APPLICATION & LICENSE MANAGEMENT
schools
id (PK)
code (UQ)
name
email
phone
address
city
province
country
school_level (ENUM: SD, SMP, SMA)
logo
status (ENUM: active, inactive, suspended)
subscription_start_date
subscription_end_date
created_at
updated_at
app_packages
id (PK)
code (UQ)
name
description
price_monthly
price_yearly
max_students
max_teachers
features (JSON)
is_active
created_at
updated_at
school_licenses
id (PK)
school_id (FK → schools)
app_package_id (FK → app_packages)
license_key (UQ)
start_date
end_date
status (ENUM: active, expired, suspended)
created_at
updated_at
payments
id (PK)
school_id (FK → schools)
school_license_id (FK → school_licenses)
amount
payment_method (ENUM: bank_transfer, credit_card, etc)
payment_date
reference_number
proof_file_url
status (ENUM: pending, verified, rejected)
notes
created_at
updated_at
payment_verifications
id (PK)
payment_id (FK → payments)
verified_by (FK → users)
verification_date
status (ENUM: approved, rejected)
notes
created_at
updated_at
school_modules
id (PK)
school_id (FK → schools)
module_name (ENUM: academic, attendance, finance, hr, library, transport, dormitory, cms)
is_enabled
enabled_at
created_at
updated_at

2. USER & ACCESS CONTROL
users
id (PK)
school_id (FK → schools, NULL for Super Admin)
username (UQ)
email (UQ)
password_hash
full_name
phone
avatar
user_type (ENUM: super_admin, school_admin, teacher, student, parent, staff)
is_active
last_login_at
created_at
updated_at
roles
id (PK)
school_id (FK → schools, NULL for system roles)
code (UQ)
name
description
is_system_role
created_at
updated_at
permissions
id (PK)
permission_code (UQ)
module_name
permission_name
description
created_at
updated_at
role_permissions (PIVOT)
id (PK)
role_id (FK → roles)
permission_id (FK → permissions)
created_at
user_roles (PIVOT)
id (PK)
user_id (FK → users)
role_id (FK → roles)
assigned_at
assigned_by (FK → users)
created_at
activity_logs
id (PK)
user_id (FK → users)
school_id (FK → schools)
action
module
description
ip_address
user_agent
created_at
deleted_records
id (PK)
user_id (FK → users)
school_id (FK → schools)
table_name
record_id
record_data (JSON)
deleted_reason
deleted_at
created_at

3. ACADEMIC MANAGEMENT
academic_sessions
id (PK)
school_id (FK → schools)
name
code
start_date
end_date
is_active
created_at
updated_at
classes
id (PK)
school_id (FK → schools)
name
code
level (ENUM: SD, SMP, SMA)
grade_number (1-12)
description
created_at
updated_at
sections
id (PK)
class_id (FK → classes)
academic_session_id (FK → academic_sessions)
name
code
room_number
capacity
teacher_id (FK → teachers, homeroom teacher)
created_at
updated_at
subjects
id (PK)
school_id (FK → schools)
name
code
description
credit_hours
created_at
updated_at
class_subjects (PIVOT)
id (PK)
class_id (FK → classes)
subject_id (FK → subjects)
academic_session_id (FK → academic_sessions)
created_at
teachers
id (PK)
user_id (FK → users)
school_id (FK → schools)
employee_number (UQ)
full_name
date_of_birth
gender
phone
email
address
qualification
specialization
join_date
status (ENUM: active, inactive, resigned)
created_at
updated_at
teacher_subjects (PIVOT)
id (PK)
teacher_id (FK → teachers)
section_id (FK → sections)
subject_id (FK → subjects)
academic_session_id (FK → academic_sessions)
assigned_at
created_at
students
id (PK)
user_id (FK → users, NULL for SD & SMP)
school_id (FK → schools)
student_number (UQ)
full_name
date_of_birth
gender
blood_type
phone
email
address
admission_date
status (ENUM: active, inactive, graduated, transferred)
created_at
updated_at
student_sections (PIVOT)
id (PK)
student_id (FK → students)
section_id (FK → sections)
academic_session_id (FK → academic_sessions)
roll_number
enrollment_date
status (ENUM: active, promoted, transferred)
created_at
parents
id (PK)
user_id (FK → users)
school_id (FK → schools)
full_name
phone
email
address
occupation
created_at
updated_at
student_parents (PIVOT)
id (PK)
student_id (FK → students)
parent_id (FK → parents)
relationship (ENUM: father, mother, guardian)
is_primary
created_at
schedules
id (PK)
section_id (FK → sections)
subject_id (FK → subjects)
teacher_id (FK → teachers)
academic_session_id (FK → academic_sessions)
day_of_week (ENUM: monday-sunday)
start_time
end_time
room_number
created_at
updated_at
assignments
id (PK)
section_id (FK → sections)
subject_id (FK → subjects)
teacher_id (FK → teachers)
academic_session_id (FK → academic_sessions)
title
description
assignment_date
due_date
max_score
attachment_url
created_at
updated_at
assignment_submissions
id (PK)
assignment_id (FK → assignments)
student_id (FK → students)
submission_date
submission_file_url
notes
score
feedback
graded_by (FK → teachers)
graded_at
created_at
updated_at
exams
id (PK)
section_id (FK → sections)
subject_id (FK → subjects)
academic_session_id (FK → academic_sessions)
exam_type (ENUM: daily, midterm, final, practice)
title
description
exam_date
duration_minutes
max_score
created_at
updated_at
exam_marks
id (PK)
exam_id (FK → exams)
student_id (FK → students)
score
notes
entered_by (FK → teachers)
entered_at
created_at
updated_at
grades
id (PK)
school_id (FK → schools)
academic_session_id (FK → academic_sessions)
grade_name (A, B, C, D, E)
min_score
max_score
grade_point
description
created_at
updated_at
student_promotions
id (PK)
student_id (FK → students)
from_section_id (FK → sections)
to_section_id (FK → sections)
from_session_id (FK → academic_sessions)
to_session_id (FK → academic_sessions)
promotion_status (ENUM: promoted, retained, graduated)
promotion_date
notes
created_by (FK → users)
created_at
updated_at
syllabi
id (PK)
subject_id (FK → subjects)
class_id (FK → classes)
academic_session_id (FK → academic_sessions)
title
description
content
file_url
created_at
updated_at

4. ATTENDANCE
student_attendance
id (PK)
student_id (FK → students)
section_id (FK → sections)
academic_session_id (FK → academic_sessions)
attendance_date
status (ENUM: present, absent, late, sick, permission)
notes
marked_by (FK → teachers)
created_at
updated_at
teacher_attendance
id (PK)
teacher_id (FK → teachers)
attendance_date
check_in_time
check_out_time
status (ENUM: present, absent, late, sick, permission)
notes
created_at
updated_at
staff_attendance
id (PK)
employee_id (FK → employees)
attendance_date
check_in_time
check_out_time
status (ENUM: present, absent, late, sick, permission)
notes
created_at
updated_at

5. HR & STAFF
employees
id (PK)
school_id (FK → schools)
user_id (FK → users)
employee_number (UQ)
full_name
date_of_birth
gender
phone
email
address
position_id (FK → positions)
join_date
employment_type (ENUM: full_time, part_time, contract)
status (ENUM: active, inactive, resigned)
created_at
updated_at
positions
id (PK)
school_id (FK → schools)
title
code (UQ)
description
department
created_at
updated_at
payrolls
id (PK)
employee_id (FK → employees)
school_id (FK → schools)
period_month
period_year
basic_salary
allowances (JSON)
deductions (JSON)
net_salary
payment_date
payment_status (ENUM: pending, paid, cancelled)
created_at
updated_at
salary_slips
id (PK)
payroll_id (FK → payrolls)
employee_id (FK → employees)
slip_number (UQ)
file_url
generated_at
created_at
updated_at

6. FINANCE
accounts
id (PK)
school_id (FK → schools)
account_name
account_code (UQ)
account_type (ENUM: asset, liability, equity, revenue, expense)
balance
description
created_at
updated_at
income
id (PK)
school_id (FK → schools)
account_id (FK → accounts)
income_category (ENUM: tuition, registration, donation, other)
amount
income_date
description
reference_number
created_by (FK → users)
created_at
updated_at
expenses
id (PK)
school_id (FK → schools)
account_id (FK → accounts)
expense_category (ENUM: salary, utilities, maintenance, supplies, other)
amount
expense_date
description
receipt_url
created_by (FK → users)
created_at
updated_at
invoices
id (PK)
school_id (FK → schools)
student_id (FK → students)
academic_session_id (FK → academic_sessions)
invoice_number (UQ)
invoice_date
due_date
total_amount
paid_amount
status (ENUM: unpaid, partial, paid, cancelled)
notes
created_at
updated_at
invoice_items
id (PK)
invoice_id (FK → invoices)
description
fee_type (ENUM: tuition, book, uniform, transport, dormitory, other)
quantity
unit_price
total_price
created_at
updated_at
student_payments
id (PK)
invoice_id (FK → invoices)
student_id (FK → students)
payment_date
amount
payment_method (ENUM: cash, bank_transfer, credit_card)
reference_number
notes
received_by (FK → users)
created_at
updated_at

7. CMS & COMMUNICATION
pages
id (PK)
school_id (FK → schools)
title
slug (UQ)
content
meta_description
is_published
published_at
created_by (FK → users)
created_at
updated_at
announcements
id (PK)
school_id (FK → schools)
title
content
target_audience (ENUM: all, teachers, students, parents, staff)
priority (ENUM: low, normal, high, urgent)
is_published
published_at
expires_at
created_by (FK → users)
created_at
updated_at
news
id (PK)
school_id (FK → schools)
title
slug
content
featured_image_url
is_published
published_at
created_by (FK → users)
created_at
updated_at
events
id (PK)
school_id (FK → schools)
title
description
event_date
start_time
end_time
location
organizer
is_published
created_by (FK → users)
created_at
updated_at
messages
id (PK)
school_id (FK → schools)
sender_id (FK → users)
recipient_id (FK → users)
subject
body
is_read
read_at
parent_message_id (FK → messages, for threading)
created_at
updated_at
notifications
id (PK)
user_id (FK → users)
school_id (FK → schools)
title
message
notification_type (ENUM: info, warning, success, error)
reference_type (e.g., assignment, exam, payment)
reference_id
is_read
read_at
created_at
updated_at
email_logs
id (PK)
school_id (FK → schools)
recipient_email
subject
body
status (ENUM: sent, failed, pending)
sent_at
error_message
created_at
updated_at
sms_logs
id (PK)
school_id (FK → schools)
recipient_phone
message
provider (ENUM: twilio, nexmo, whatsapp)
status (ENUM: sent, failed, pending)
sent_at
error_message
created_at
updated_at

8. FACILITIES
library_books
id (PK)
school_id (FK → schools)
isbn
title
author
publisher
publication_year
category
quantity
available_quantity
location
created_at
updated_at
library_transactions
id (PK)
book_id (FK → library_books)
borrower_type (ENUM: student, teacher, staff)
borrower_id (student_id or teacher_id or employee_id)
borrow_date
due_date
return_date
status (ENUM: borrowed, returned, overdue)
fine_amount
notes
created_at
updated_at
transport_routes
id (PK)
school_id (FK → schools)
route_name
route_code
description
fare
created_at
updated_at
transport_vehicles
id (PK)
school_id (FK → schools)
route_id (FK → transport_routes)
vehicle_number
vehicle_type (ENUM: bus, van, car)
capacity
driver_name
driver_phone
status (ENUM: active, maintenance, inactive)
created_at
updated_at
student_transport
id (PK)
student_id (FK → students)
vehicle_id (FK → transport_vehicles)
pickup_point
drop_point
start_date
end_date
status (ENUM: active, inactive)
created_at
updated_at
dormitory_rooms
id (PK)
school_id (FK → schools)
room_number
room_type (ENUM: single, double, quad, shared)
capacity
floor
building
status (ENUM: available, occupied, maintenance)
created_at
updated_at
dormitory_assignments
id (PK)
room_id (FK → dormitory_rooms)
student_id (FK → students)
academic_session_id (FK → academic_sessions)
start_date
end_date
status (ENUM: active, inactive)
created_at
updated_at

9. REPORTS (View/Aggregate Tables - Optional)
Reports are typically generated dynamically from existing tables. However, if you need materialized reports for performance:
academic_report_cards (Optional - Generated)
id (PK)
student_id (FK → students)
academic_session_id (FK → academic_sessions)
report_data (JSON)
generated_at
created_at

KEY RELATIONSHIPS SUMMARY
Multi-tenancy
All school-specific tables have school_id (FK → schools)
Super Admin users have school_id = NULL
User Hierarchy
users → user_roles (PIVOT) → roles → role_permissions (PIVOT) → permissions
users.user_type defines base type
students.user_id is NULL for SD & SMP levels
Academic Flow
schools → academic_sessions → sections → students via student_sections
sections → schedules, assignments, exams
teachers → teacher_subjects → teaches specific subjects in sections
Finance Flow
students → invoices → invoice_items
student_payments
schools → accounts → income / expenses
Facilities
library_books → library_transactions (polymorphic borrower)
transport_routes → transport_vehicles → student_transport
dormitory_rooms → dormitory_assignments