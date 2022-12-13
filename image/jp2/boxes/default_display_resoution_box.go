package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/conv"

	"math"
)

type DefaultDisplayResolutionBox struct {
	*SignatureBox
	VerticalDisplayGridResolutionNumerator     []byte
	VerticalDisplayGridResolutionDenominator   []byte
	HorizontalDisplayGridResolutionNumerator   []byte
	HorizontalDisplayGridResolutionDenominator []byte
	VerticalDisplayGridResolutionExponent      []byte
	HorizontalDisplayGridResolutionExponent    []byte
}

func (d *DefaultDisplayResolutionBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_DEFAULT_DISPLAY_RESOLUTION
}

func (d *DefaultDisplayResolutionBox) GetLength() uint64 {
	return d.Length
}

func (d *DefaultDisplayResolutionBox) GetOffset() uint64 {
	return d.Offset
}

func (d *DefaultDisplayResolutionBox) Decode(reader *buffer.ByteBuffer) error {
	reader.Read(d.VerticalDisplayGridResolutionNumerator)
	reader.Read(d.VerticalDisplayGridResolutionDenominator)
	reader.Read(d.HorizontalDisplayGridResolutionNumerator)
	reader.Read(d.HorizontalDisplayGridResolutionDenominator)
	reader.Read(d.VerticalDisplayGridResolutionExponent)
	reader.Read(d.HorizontalDisplayGridResolutionExponent)

	return nil
}

func (d *DefaultDisplayResolutionBox) GetVerticalDisplayGridResolutionNumerator() uint16 {
	return binary.BigEndian.Uint16(d.VerticalDisplayGridResolutionNumerator)
}

func (d *DefaultDisplayResolutionBox) GetVerticalDisplayGridResolutionDenominator() uint16 {
	return binary.BigEndian.Uint16(d.VerticalDisplayGridResolutionDenominator)
}

func (d *DefaultDisplayResolutionBox) GetHorizontalDisplayGridResolutionNumerator() uint16 {
	return binary.BigEndian.Uint16(d.HorizontalDisplayGridResolutionNumerator)
}

func (d *DefaultDisplayResolutionBox) GetHorizontalDisplayGridResolutionDenominator() uint16 {
	return binary.BigEndian.Uint16(d.HorizontalDisplayGridResolutionDenominator)
}
func (d *DefaultDisplayResolutionBox) GetVerticalDisplayGridResolutionExponent() int8 {
	return conv.Int8(binary.BigEndian.Uint16(d.VerticalDisplayGridResolutionExponent))
}
func (d *DefaultDisplayResolutionBox) GetHorizontalDisplayGridResolutionExponent() int8 {
	return conv.Int8(binary.BigEndian.Uint16(d.HorizontalDisplayGridResolutionExponent))
}

// GetVerticalDisplayGridResolution
// VRd = VRdN/VRdD * 10^VRdE
func (d *DefaultDisplayResolutionBox) GetVerticalDisplayGridResolution() uint64 {
	return conv.Uint64(conv.Float64(d.GetVerticalDisplayGridResolutionNumerator()) / conv.Float64(d.GetVerticalDisplayGridResolutionDenominator()) * math.Pow(conv.Float64(10), conv.Float64(d.GetVerticalDisplayGridResolutionExponent())))
}

// HorizontalDisplayGridResolution
// HRd = HRdN/HRdD * 10^HRdE
func (d *DefaultDisplayResolutionBox) HorizontalDisplayGridResolution() uint64 {
	return conv.Uint64(conv.Float64(d.GetHorizontalDisplayGridResolutionNumerator()) / conv.Float64(d.GetHorizontalDisplayGridResolutionDenominator()) * math.Pow(conv.Float64(10), conv.Float64(d.GetHorizontalDisplayGridResolutionExponent())))
}
