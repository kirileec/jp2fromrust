package types

type ComponentMap struct {
	Component   []byte
	MappingType ComponentMapType
	Palette     []byte
}

func (c *ComponentMap) GetComponent() uint16 {
	panic("")
}
func (c *ComponentMap) GetMappingType() uint8 {
	panic("")
}

func (c *ComponentMap) GetPalette() uint8 {
	return c.Palette[0]
}
