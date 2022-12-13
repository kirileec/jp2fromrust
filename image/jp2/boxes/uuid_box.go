package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type UUIDBox struct {
	*SignatureBox
	UUID []byte
	Data []uint8
}

func (i *UUIDBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_UUID
}

func (i *UUIDBox) GetLength() uint64 {
	return i.Length
}

func (i *UUIDBox) GetOffset() uint64 {
	return i.Offset
}

func (i *UUIDBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(i.UUID)

	i.Data = make([]byte, int(i.Length)-len(i.UUID))
	reader.Read(i.Data)
	return nil
}
