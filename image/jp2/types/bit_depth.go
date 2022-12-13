package types

type BitDepth struct {
	a     bitDepth
	value uint8
}

type bitDepth uint8

const (
	BDSigned bitDepth = iota
	BDUnsigned
	BDReserved
)

func NewBitDepth(b uint8) BitDepth {
	value := uint8(b<<1>>1) + 1
	signedness := b >> 7
	switch signedness {
	case 1:
		return BitDepth{a: BDSigned, value: value}
	case 0:
		return BitDepth{a: BDUnsigned, value: value}
	default:
		return BitDepth{a: BDReserved, value: value}

	}
}

func (b BitDepth) Value() uint8 {
	switch b.a {
	case BDSigned:
		return uint8(b.value)
	case BDUnsigned:
		return uint8(b.value)
	default:
		return uint8(b.value)
	}
}
