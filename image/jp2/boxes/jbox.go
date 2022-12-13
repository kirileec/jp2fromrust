package boxes

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type JBox struct {
	*SignatureBox
}

func (J *JBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_SIGNATURE
}

func (J *JBox) GetLength() uint64 {
	return J.Length
}

func (J *JBox) GetOffset() uint64 {
	return J.Offset
}

func (J *JBox) Decode(reader *buffer.ByteBuffer) error {
	J.Length = 12
	var buffer = make([]byte, 4)
	_, err := reader.Read(buffer)
	if err != nil {
		return err
	}
	if !SIGNATURE_MAGIC.Equal(buffer) {
		return errors.New("invalid signature")
	}
	return nil
}
