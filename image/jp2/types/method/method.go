package method

type Method byte

var (
	METHOD_ENUMERATED_COLOR_SPACE            = byte(1)
	METHOD_ENUMERATED_RESTRICTED_ICC_PROFILE = byte(2)
)

type Methods byte

func (m Methods) String() string {
	switch m {
	case EnumeratedColorSpace:
		return "EnumeratedColorSpace"
	case RestrictedICCProfile:
		return "RestrictedICCProfile"
	default:
		return "Reserved"
	}
}

const (
	EnumeratedColorSpace Methods = iota
	RestrictedICCProfile
	Reserved
)

func NewMethods(value []byte) Methods {
	switch value[0] {
	case METHOD_ENUMERATED_COLOR_SPACE:
		return EnumeratedColorSpace
	case METHOD_ENUMERATED_RESTRICTED_ICC_PROFILE:
		return RestrictedICCProfile
	default:
		return Reserved
	}
}
