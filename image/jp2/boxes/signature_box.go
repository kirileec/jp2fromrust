package boxes

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type SignatureBox struct {
	Length uint64
	Offset uint64
}

func NewSignatureBox(length uint64, offset uint64) *SignatureBox {
	return &SignatureBox{Length: length, Offset: offset}
}

func (s *SignatureBox) SetOffset(offset uint64) {
	s.Offset = offset
}

func (s *SignatureBox) SetLength(length uint64) {
	s.Length = length
}

func (s *SignatureBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_SIGNATURE
}

func (s *SignatureBox) GetLength() uint64 {
	return s.Length
}

func (s *SignatureBox) GetOffset() uint64 {
	return s.Offset
}

func (s *SignatureBox) Decode(reader *buffer.ByteBuffer) error {
	s.Length = 12
	buffer := make([]byte, 4)
	reader.Read(buffer)
	if !SIGNATURE_MAGIC.EqualA(buffer) {
		return errors.New("Invalid signature")
	}
	return nil
}

// <CR><LF><0x87><LF> (0x0D0A 870A).
var SIGNATURE_MAGIC = box_type.BoxType{13, 10, 135, 10}

func (s *SignatureBox) Signature() box_type.BoxType {
	return SIGNATURE_MAGIC
}
