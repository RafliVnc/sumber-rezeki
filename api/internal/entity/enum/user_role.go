package enum

type UserRole string

const (
	SuperAdmin    UserRole = "SUPER_ADMIN"
	Owner         UserRole = "OWNER"
	WarehouseHead UserRole = "WAREHOUSE_HEAD"
	Treasurer     UserRole = "TREASURER"
)
