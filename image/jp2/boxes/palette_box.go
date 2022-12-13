package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type PaletteBox struct {
	*SignatureBox
	NumEntries          []byte
	NumComponents       []byte
	GeneratedComponents []types.GeneratedComponent
}

func (pb *PaletteBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_PALETTE
}

func (pb *PaletteBox) GetLength() uint64 {
	return pb.Length
}

func (pb *PaletteBox) GetOffset() uint64 {
	return pb.Offset
}

func (pb *PaletteBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(pb.NumEntries)
	reader.Read(pb.NumComponents)
	numEntries := int(pb.GetNumEntries())
	pb.GeneratedComponents = make([]types.GeneratedComponent, 0)
	l := int(pb.GetNumComponents())

	for i := 0; i < l; i++ {
		pb.GeneratedComponents = append(pb.GeneratedComponents, types.GeneratedComponent{
			BitDepth: []byte{0},
			Values:   make([]byte, numEntries),
		})
	}

	for _, component := range pb.GeneratedComponents {
		reader.Read(component.BitDepth)
	}
	for _, component := range pb.GeneratedComponents {
		j := 0
		for j < numEntries {
			entry := make([]byte, 1)
			reader.Read(entry)

			component.Values = append(component.Values, entry[0])

			j++
		}
	}
	return nil
}

func (pb *PaletteBox) GetNumEntries() uint16 {
	return binary.BigEndian.Uint16(pb.NumEntries)
}
func (pb *PaletteBox) GetNumComponents() uint8 {
	return pb.NumComponents[0]
}
func (pb *PaletteBox) GetGeneratedComponents() []types.GeneratedComponent {
	return pb.GeneratedComponents
}
