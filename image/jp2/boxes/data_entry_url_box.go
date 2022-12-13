package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type DataEntryURLBox struct {
	*SignatureBox
	Version  []uint8
	Flags    []uint8
	Location []uint8
}

func (d *DataEntryURLBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_DATA_ENTRY_URL
}

func (d *DataEntryURLBox) GetLength() uint64 {
	return d.Length
}

func (d *DataEntryURLBox) GetLocation() string {
	return string(d.Location)
}

func (d *DataEntryURLBox) GetOffset() uint64 {
	return d.Offset
}

func (d *DataEntryURLBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(d.Version)
	reader.Read(d.Flags)
	size := d.Length - 4

	for size > 0 {
		buffer := make([]byte, 1)
		reader.Read(buffer)
		d.Location = buffer
		size -= 1
	}
	return nil

}
