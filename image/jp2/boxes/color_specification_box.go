package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types/enumerated_color_space"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types/method"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"
	"github.com/linxlib/logs"
)

type ColorSpecificationBox struct {
	*SignatureBox
	Method                  []byte
	Precedence              []byte
	ColorSpaceApproximation []byte
	EnumeratedColorSpace    enumerated_color_space.EnumeratedColorSpace
	RestrictedIccProfile    []uint8
}

func (c *ColorSpecificationBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_COLOR_SPECIFICATION
}

func (c *ColorSpecificationBox) GetLength() uint64 {
	return c.Length
}

func (c *ColorSpecificationBox) GetOffset() uint64 {
	return c.Offset
}

func (c *ColorSpecificationBox) Decode(reader *buffer.ByteBuffer) error {
	_, err := reader.Read(c.Method)
	if err != nil {
		return err
	}
	_, err = reader.Read(c.Precedence)
	if err != nil {
		return err
	}
	_, err = reader.Read(c.ColorSpaceApproximation)
	if err != nil {
		return err
	}
	if c.GetPrecedence() != 0 {
		logs.Warnf("Precedence %d Unexpected", c.GetPrecedence())
	}
	if c.GetColorSpaceApproximation() != 0 {
		logs.Warnf("Color-space Approximation %d Unexpected", c.GetColorSpaceApproximation())
	}

	logs.Debugf("Method: %s Precedence: %d ColorSpaceApproximation: %d", c.GetMethod(), c.GetPrecedence(), c.GetColorSpaceApproximation())

	switch c.GetMethod() {
	case method.EnumeratedColorSpace:
		_, err := reader.Read(c.EnumeratedColorSpace)
		if err != nil {
			return err
		}
	case method.RestrictedICCProfile:
		_, err := reader.Read(c.RestrictedIccProfile)
		if err != nil {
			return err
		}
	case method.Reserved:

	}

	return nil

}

func (c *ColorSpecificationBox) GetMethod() method.Methods {
	return method.NewMethods(c.Method)
}
func (c *ColorSpecificationBox) GetPrecedence() int8 {
	return conv.Int8(c.Precedence[0])
}
func (c *ColorSpecificationBox) GetColorSpaceApproximation() uint8 {
	return c.ColorSpaceApproximation[0]
}

func (c *ColorSpecificationBox) GetEnumeratedColorSpace() uint32 {
	return binary.BigEndian.Uint32(c.EnumeratedColorSpace)
}
