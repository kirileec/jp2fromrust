package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type XMLBox struct {
	*SignatureBox
	XML []uint8
}

func (x *XMLBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_XML
}

func (x *XMLBox) GetLength() uint64 {
	return x.Length
}

func (x *XMLBox) GetOffset() uint64 {
	return x.Offset
}

func (x *XMLBox) Decode(reader *buffer.ByteBuffer) error {
	x.XML = make([]byte, x.Length)
	reader.Read(x.XML)
	return nil
}

func (x *XMLBox) String() string {
	return ""
}
