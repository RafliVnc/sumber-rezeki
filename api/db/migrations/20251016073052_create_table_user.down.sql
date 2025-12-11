-- DropIndex
DROP INDEX IF EXISTS "sales_routes_route_id_idx";
DROP INDEX IF EXISTS "sales_routes_sales_id_idx";
DROP INDEX IF EXISTS "payrolls_period_id_idx";
DROP INDEX IF EXISTS "payrolls_employee_id_idx";
DROP INDEX IF EXISTS "employee_attendances_date_idx";
DROP INDEX IF EXISTS "employee_attendances_employee_id_idx";
DROP INDEX IF EXISTS "employees_role_idx";
DROP INDEX IF EXISTS "employees_supervisor_id_idx";

-- DropTable (urutan terbalik karena foreign key dependencies)
DROP TABLE IF EXISTS "sales_routes";
DROP TABLE IF EXISTS "routes";
DROP TABLE IF EXISTS "sales";
DROP TABLE IF EXISTS "payrolls";
DROP TABLE IF EXISTS "employee_attendances";
DROP TABLE IF EXISTS "employees";
DROP TABLE IF EXISTS "period_closures";
DROP TABLE IF EXISTS "periods";
DROP TABLE IF EXISTS "users";

-- DropEnum
DROP TYPE IF EXISTS "PayrollModule";
DROP TYPE IF EXISTS "AttendanceStatus";
DROP TYPE IF EXISTS "EmployeeRole";
DROP TYPE IF EXISTS "PeriodType";
DROP TYPE IF EXISTS "UserRole";