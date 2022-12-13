package types

type GeneratedComponent struct {
	BitDepth []byte

	Values []uint8
}

func (g *GeneratedComponent) GetBitDepth() BitDepth {
	panic("")
}
func (g *GeneratedComponent) GetValues() []uint8 {
	panic("")
}
