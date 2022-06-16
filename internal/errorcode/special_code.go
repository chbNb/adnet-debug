package errorcode

type SpecialCode int

const (
	EXCEPTION_SPECIAL_ZERO         SpecialCode = 0
	EXCEPTION_SPECIAL_EMPTY_STRING SpecialCode = 1
)

func (specialCode SpecialCode) String() string {
	switch specialCode {
	case EXCEPTION_SPECIAL_ZERO:
		return "0"
	case EXCEPTION_SPECIAL_EMPTY_STRING:
		return ""
	default:
		return "unset"
	}
}

func (specialCode SpecialCode) Error() string {
	return specialCode.String()
}
