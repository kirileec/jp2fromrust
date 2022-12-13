package quantization_style

type QuantizationStyle byte

const (
	No QuantizationStyle = iota
	ScalarDerived
	ScalarExpounded
	Reserved
)

func NewQuantizationStyle(b uint8) QuantizationStyle {
	value := b << 3
	value = value >> 3
	switch value {
	case 0x00000000:
		return No
	case 0x00000001:
		return ScalarDerived
	case 0x00000010:
		return ScalarExpounded
	default:
		return Reserved
	}
}
