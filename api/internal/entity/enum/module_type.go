package enum

type PayrollModule string
type ModuleType string

const (
	SALES_INCENTIVE PayrollModule = "SALES_INCENTIVE"
	OPERATIONAL     PayrollModule = "OPERATIONAL"
)

const (
	EMPLOYEE ModuleType = "EMPLOYEE"
)
