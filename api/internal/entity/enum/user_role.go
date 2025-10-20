package enum

type UserRole string

const (
	SUPER_ADMIN    UserRole = "SUPER_ADMIN"
	OWNER          UserRole = "OWNER"
	WAREHOUSE_HEAD UserRole = "WAREHOUSE_HEAD"
	TREASURER      UserRole = "TREASURER"
)
