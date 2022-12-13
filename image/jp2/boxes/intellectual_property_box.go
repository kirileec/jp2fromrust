package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type IntellectualPropertyBox struct {
	*SignatureBox
	Data []byte
}

func (i *IntellectualPropertyBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_INTELLECTUAL_PROPERTY
}
func (i *IntellectualPropertyBox) Decode(reader *buffer.ByteBuffer) error {
	i.Data = make([]byte, i.Length)
	reader.Read(i.Data)
	return nil
}
