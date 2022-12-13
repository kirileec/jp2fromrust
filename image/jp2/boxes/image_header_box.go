package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type ImageHeaderBox struct {
	*SignatureBox
	Height               box_type.BoxType
	Width                box_type.BoxType
	ComponentsNum        []byte
	ComponentsBits       []byte
	CompressionType      []byte
	ColorsPaceUnknown    []byte
	IntellectualProperty []byte
}

func (i *ImageHeaderBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_IMAGE_HEADER
}

func (i *ImageHeaderBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(i.Height)
	reader.Read(i.Width)
	reader.Read(i.ComponentsNum)
	reader.Read(i.ComponentsBits)
	reader.Read(i.CompressionType)
	reader.Read(i.ColorsPaceUnknown)
	reader.Read(i.IntellectualProperty)

	return nil
}

func (i *ImageHeaderBox) GetHeight() uint32 {
	return binary.BigEndian.Uint32(i.Height)
}

func (i *ImageHeaderBox) GetWidth() uint32 {
	return binary.BigEndian.Uint32(i.Width)
}
func (i *ImageHeaderBox) GetComponentsNum() uint16 {
	return binary.BigEndian.Uint16(i.ComponentsNum)
}
func (i *ImageHeaderBox) GetComponentsBits() uint8 {
	if i.ComponentsBits[0] == 255 {
		return i.ComponentsBits[0]
	} else if i.ComponentsBits[0] <= 37 {
		return i.ComponentsBits[0] + 1
	} else {

	}

	return 0

}
func (i *ImageHeaderBox) GetCompressionType() uint8 {
	return i.CompressionType[0]
}

func (i *ImageHeaderBox) GetColorSpaceUnknown() uint8 {
	return i.ColorsPaceUnknown[0]
}

func (i *ImageHeaderBox) GetIntellectualProperty() uint8 {
	return i.IntellectualProperty[0]
}
