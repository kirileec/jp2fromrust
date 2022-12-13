package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type UUIDInfoSuperBox struct {
	*SignatureBox
	UUIDList        []*UUIDListBox
	DataEntryUrlBox []*DataEntryURLBox
}

func (i *UUIDInfoSuperBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_UUID_INFO
}
func (i *UUIDInfoSuperBox) Decode(reader *buffer.ByteBuffer) error {
	return nil
}
