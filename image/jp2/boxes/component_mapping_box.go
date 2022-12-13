package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type ComponentMappingBox struct {
	*SignatureBox
	Mapping []types.ComponentMap
}

func (c *ComponentMappingBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_COMPONENT_MAPPING
}

func (c *ComponentMappingBox) GetLength() uint64 {
	return c.Length
}

func (c *ComponentMappingBox) GetOffset() uint64 {
	return c.Offset
}

func (c *ComponentMappingBox) Decode(reader *buffer.ByteBuffer) error {
	index := uint64(0)
	for index < c.Length {
		componentMap := types.ComponentMap{
			Component:   make([]byte, 2),
			MappingType: types.NewComponentMapType(255),
			Palette:     make([]byte, 1),
		}
		reader.Read(componentMap.Component)
		mappingType := make([]byte, 1)
		reader.Read(mappingType)
		componentMap.MappingType = types.NewComponentMapType(mappingType[0])
		reader.Read(componentMap.Palette)
		c.Mapping = append(c.Mapping, componentMap)
		index += 4
	}
	return nil
}

func (c *ComponentMappingBox) GetComponentMap() []types.ComponentMap {
	return c.Mapping
}
