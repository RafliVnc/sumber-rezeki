-- CreateEnum
CREATE TYPE "UserRole" AS ENUM ('SUPER_ADMIN', 'OWNER', 'WAREHOUSE_HEAD', 'TREASURER');
CREATE TYPE "PeriodType" AS ENUM ('WEEKLY', 'MONTHLY');
CREATE TYPE "EmployeeRole" AS ENUM ('WAREHOUSE_HEAD', 'SALES', 'DRIVER', 'HELPER', 'TREASURER', 'STAFF');
CREATE TYPE "AttendanceStatus" AS ENUM ('PRESENT', 'ABSENT');
CREATE TYPE "PayrollModule" AS ENUM ('SALES_INCENTIVE', 'OPERATIONAL');

-- CreateTable: users
CREATE TABLE "users" (
    "id" VARCHAR(200) PRIMARY KEY,
    "name" VARCHAR(200) NOT NULL,
    "username" VARCHAR(200) NOT NULL UNIQUE,
    "phone" VARCHAR(20) NOT NULL UNIQUE,
    "role" "UserRole" NOT NULL,
    "password" VARCHAR(200) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

-- CreateTable: periods
CREATE TABLE "periods" (
    "id" SERIAL PRIMARY KEY,
    "type" "PeriodType" NOT NULL,
    "start_date" TIMESTAMP(3) NOT NULL,
    "end_date" TIMESTAMP(3) NOT NULL,
    "week_number" INTEGER NOT NULL,
    "month" INTEGER NOT NULL,
    "year" INTEGER NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "is_closed" BOOLEAN NOT NULL DEFAULT false,
    "closed_at" TIMESTAMP(3),
    "closed_by" VARCHAR REFERENCES "users"("id") ON DELETE SET NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

-- CreateTable: period_closures
CREATE TABLE "period_closures" (
    "id" SERIAL PRIMARY KEY,
    "module_name" VARCHAR(100) NOT NULL,
    "notes" TEXT NOT NULL,
    "is_closed" BOOLEAN NOT NULL DEFAULT false,
    "closed_at" TIMESTAMP(3),
    "closed_by" VARCHAR REFERENCES "users"("id") ON DELETE SET NULL,
    "print_count" INTEGER NOT NULL DEFAULT 0,
    "period_id" INTEGER NOT NULL REFERENCES "periods"("id") ON DELETE SET NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3),
    UNIQUE ("period_id", "module_name")
);

-- CreateTable: employees
CREATE TABLE "employees" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "salary" DECIMAL(12,2) NOT NULL,
    "role" "EmployeeRole" NOT NULL,
    "join_date" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "supervisor_id" INTEGER REFERENCES "employees"("id") ON DELETE SET NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

-- CreateTable: employee_attendances
CREATE TABLE "employee_attendances" (
    "id" SERIAL PRIMARY KEY,
    "date" DATE NOT NULL,
    "status" "AttendanceStatus" NOT NULL,
    "employee_id" INTEGER NOT NULL REFERENCES "employees"("id") ON DELETE RESTRICT,
    "period_id" INTEGER NOT NULL REFERENCES "periods"("id") ON DELETE RESTRICT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3),
    UNIQUE ("date", "employee_id")
);

-- CreateTable: payrolls
CREATE TABLE "payrolls" (
    "id" SERIAL PRIMARY KEY,
    "base_salary" DECIMAL(12,2) NOT NULL,
    "attendance_days" INTEGER NOT NULL,
    "deductions" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "bonuses" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "module_type" "PayrollModule" NOT NULL,
    "notes" TEXT,
    "is_paid" BOOLEAN NOT NULL DEFAULT false,
    "paid_at" TIMESTAMP(3),
    "employee_id" INTEGER NOT NULL REFERENCES "employees"("id") ON DELETE RESTRICT,
    "period_id" INTEGER NOT NULL REFERENCES "periods"("id") ON DELETE RESTRICT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3),
    UNIQUE ("employee_id", "period_id")
);

-- CreateTable: sales
CREATE TABLE "sales" (
    "id" SERIAL PRIMARY KEY,
    "phone" VARCHAR(20) NOT NULL,
    "employee_id" INTEGER NOT NULL UNIQUE REFERENCES "employees"("id") ON DELETE RESTRICT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

-- CreateTable: routes
CREATE TABLE "routes" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

-- CreateTable: sales_routes
CREATE TABLE "sales_routes" (
    "id" SERIAL PRIMARY KEY,
    "route_id" INTEGER NOT NULL REFERENCES "routes"("id") ON DELETE RESTRICT,
    "sales_id" INTEGER NOT NULL REFERENCES "sales"("id") ON DELETE RESTRICT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3),
    UNIQUE ("route_id", "sales_id")
);

-- CreateIndex
CREATE INDEX "employees_supervisor_id_idx" ON "employees"("supervisor_id");
CREATE INDEX "employees_role_idx" ON "employees"("role");
CREATE INDEX "employee_attendances_employee_id_idx" ON "employee_attendances"("employee_id");
CREATE INDEX "employee_attendances_date_idx" ON "employee_attendances"("date");
CREATE INDEX "payrolls_employee_id_idx" ON "payrolls"("employee_id");
CREATE INDEX "payrolls_period_id_idx" ON "payrolls"("period_id");
CREATE INDEX "sales_routes_sales_id_idx" ON "sales_routes"("sales_id");
CREATE INDEX "sales_routes_route_id_idx" ON "sales_routes"("route_id");

-- Seed Admin
INSERT INTO users (id, name, username, phone, role, password) 
VALUES ('7f5dc73e-e097-4e8e-ba7c-5ed828fabc74', 'Super Admin', 'superadmin', '0888888888', 'SUPER_ADMIN', '$2a$10$ibmW.pidc9yRifFckQJFZ.1Hs4BEkf8.B.b5.xSJio3nsR.3y7rQO');