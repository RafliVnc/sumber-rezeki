enum UserRole {
  SUPER_ADMIN = "SUPER_ADMIN",
  OWNER = "OWNER",
  WAREHOUSE_HEAD = "WAREHOUSE_HEAD",
  TREASURER = "TREASURER",
}

type User = {
  id: string;
  name: string;
  username: string;
  password?: string;
  role: UserRole;
};
