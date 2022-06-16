package mvconst

const (
	PAY_TYPE_CPI = iota + 1
	PAY_TYPE_CPC
	PAY_TYPE_CPM
	PAY_TYPE_CPA
	PAY_TYPE_CPE
)

var InstallPayTypes []int = []int{
	PAY_TYPE_CPI,
	PAY_TYPE_CPA,
	PAY_TYPE_CPE,
}
