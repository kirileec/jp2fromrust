package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"
	"io"
)

type ContiguousCodeStreamBox struct {
	*SignatureBox
}

func (c *ContiguousCodeStreamBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_CONTIGUOUS_CODESTREAM
}

func (c *ContiguousCodeStreamBox) GetLength() uint64 {
	return c.Length
}

func (c *ContiguousCodeStreamBox) GetOffset() uint64 {
	return c.Offset
}

func (c *ContiguousCodeStreamBox) Decode(reader *buffer.ByteBuffer) error {
	if c.Length == 0 {
		reader.Seek(0, io.SeekEnd)
		c.Length = conv.Uint64(uint64(reader.ReadOffset()) - c.Offset)
	} else {
		reader.Seek(int64(c.Length), io.SeekCurrent)
	}
	return nil
}
