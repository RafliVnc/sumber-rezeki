enum EmployeeRole {
  EMPLOYEE_WAREHOUSE_HEAD = "WAREHOUSE_HEAD",
  EMPLOYEE_TREASURER = "TREASURER",
  SALES = "SALES",
  DRIVER = "DRIVER",
  HELPER = "HELPER",
  STAFF = "STAFF",
}

type Employee = {
  id: number;
  name: string;
  salary: float;
  supervisorId: number;
  role: EmployeeRole;
  Sales?: Sales;
};
