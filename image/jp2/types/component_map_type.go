package types

var (
	COMPONENT_MAP_TYPE_DIRECT  = byte(1)
	COMPONENT_MAP_TYPE_PALETTE = byte(2)
)

type ComponentMapType byte

const (
	CMTDirect ComponentMapType = iota
	CMTPalette
	CMTReserved
)

func NewComponentMapType(value byte) ComponentMapType {
	switch value {
	case COMPONENT_MAP_TYPE_DIRECT:
		return CMTDirect
	case COMPONENT_MAP_TYPE_PALETTE:
		return CMTPalette
	default:
		return CMTReserved
	}
}
