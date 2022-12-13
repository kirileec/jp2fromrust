package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"
)

type UUIDListBox struct {
	*SignatureBox
	NumberOfUuids []uint8
	IDs           []uint8
}

func (u *UUIDListBox) GetNumberOfUUIDs() int16 {
	return conv.Int16(binary.BigEndian.Uint16(u.NumberOfUuids))
}

func (u *UUIDListBox) GetIds() []string {
	t := make([]string, 0)
	for _, d := range u.IDs {
		t = append(t, string(d))
	}
	return t
}

func (u *UUIDListBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_UUID_LIST
}

func (u *UUIDListBox) GetLength() uint64 {
	return u.Length
}

func (u *UUIDListBox) GetOffset() uint64 {
	return u.Offset
}

func (u *UUIDListBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(u.NumberOfUuids)
	size := u.GetNumberOfUUIDs()

	u.IDs = make([]byte, size)
	buffer := make([]byte, 16)
	for size > 0 {
		reader.Read(buffer)

		for i, b := range buffer {
			u.IDs[i] = b
		}

		size -= 1

	}
	return nil
}
