package box_type

import (
	"encoding/binary"
	"errors"
	"fmt"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"
	"github.com/linxlib/logs"
	"io"
)

type BoxHeader struct {
	BoxLength    uint64
	BoxType      []byte
	HeaderLength uint8
}

func (b BoxHeader) String() string {
	return fmt.Sprintf("BoxType: %s BoxLength: %d HeaderLength: %d", NewBoxTypes(b.BoxType), b.BoxLength, b.HeaderLength)
}

func DecodeBoxHeader(reader *buffer.ByteBuffer) (*BoxHeader, error) {

	headerLength := byte(8)
	boxLength := make(BoxType, 4)
	boxType := make(BoxType, 4)
	_, err := reader.Read(boxLength)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, err
		}
		logs.Error(err)
		return nil, err
	}
	boxLengthValue := conv.Uint64(binary.BigEndian.Uint32(boxLength))
	if boxLengthValue == 0 {
		_, err := reader.Read(boxType)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
	} else if boxLengthValue == 1 {
		_, err := reader.Read(boxType)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		xlLength := make([]byte, 8)
		_, err = reader.Read(xlLength)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		boxLengthValue = binary.BigEndian.Uint64(xlLength) - 16
		headerLength = 16
	} else if boxLengthValue <= 7 {
		panic(fmt.Sprintf("unsupported reserved box length %d", boxLengthValue))
	} else {
		_, err := reader.Read(boxType)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		boxLengthValue = boxLengthValue - 8
	}

	return &BoxHeader{
		BoxLength:    boxLengthValue,
		BoxType:      boxType,
		HeaderLength: headerLength,
	}, nil

}
