package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"
	"math"
)

type CaptureResolutionBox struct {
	*SignatureBox
	VerticalCaptureGridResolutionNumerator     []byte
	VerticalCaptureGridResolutionDenominator   []byte
	HorizontalCaptureGridResolutionNumerator   []byte
	HorizontalCaptureGridResolutionDenominator []byte
	VerticalCaptureGridResolutionExponent      []byte
	HorizontalCaptureGridResolutionExponent    []byte
}

func (c *CaptureResolutionBox) GetVerticalCaptureGridResolutionNumerator() uint16 {
	return binary.BigEndian.Uint16(c.VerticalCaptureGridResolutionNumerator)
}
func (c *CaptureResolutionBox) GetVerticalCaptureGridResolutionDenominator() uint16 {
	return binary.BigEndian.Uint16(c.VerticalCaptureGridResolutionDenominator)
}
func (c *CaptureResolutionBox) GetHorizontalCaptureGridResolutionNumerator() uint16 {
	return binary.BigEndian.Uint16(c.HorizontalCaptureGridResolutionNumerator)
}
func (c *CaptureResolutionBox) GetHorizontalCaptureGridResolutionDenominator() uint16 {
	return binary.BigEndian.Uint16(c.HorizontalCaptureGridResolutionDenominator)
}
func (c *CaptureResolutionBox) GetVerticalCaptureGridResolutionExponent() int8 {
	return conv.Int8(binary.BigEndian.Uint16(c.VerticalCaptureGridResolutionExponent))
}
func (c *CaptureResolutionBox) GetHorizontalCaptureGridResolutionExponent() int8 {
	return conv.Int8(binary.BigEndian.Uint16(c.HorizontalCaptureGridResolutionExponent))
}

func (c *CaptureResolutionBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_CAPTURE_RESOLUTION
}

func (c *CaptureResolutionBox) GetLength() uint64 {
	return c.Length
}

func (c *CaptureResolutionBox) GetOffset() uint64 {
	return c.Offset
}

func (c *CaptureResolutionBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(c.VerticalCaptureGridResolutionNumerator)
	reader.Read(c.VerticalCaptureGridResolutionDenominator)
	reader.Read(c.HorizontalCaptureGridResolutionNumerator)
	reader.Read(c.HorizontalCaptureGridResolutionDenominator)
	reader.Read(c.VerticalCaptureGridResolutionExponent)
	reader.Read(c.HorizontalCaptureGridResolutionExponent)
	return nil
}

func (c *CaptureResolutionBox) VerticalResolutionCapture() float64 {
	verticalResolutionCapture := conv.Float64(c.GetVerticalCaptureGridResolutionNumerator()) / conv.Float64(c.GetVerticalCaptureGridResolutionDenominator())
	verticalResolutionCapture *= math.Pow(conv.Float64(10), conv.Float64(c.GetVerticalCaptureGridResolutionExponent()))
	return verticalResolutionCapture
}

func (c *CaptureResolutionBox) HorizontalResolutionCapture() float64 {
	horizontalResolutionCapture := conv.Float64(c.GetHorizontalCaptureGridResolutionNumerator()) / conv.Float64(c.GetHorizontalCaptureGridResolutionDenominator())
	horizontalResolutionCapture *= math.Pow(conv.Float64(10), conv.Float64(c.GetHorizontalCaptureGridResolutionExponent()))
	return horizontalResolutionCapture
}
