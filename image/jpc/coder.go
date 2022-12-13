package jpc

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/boxes"
	"github.com/linxlib/conv"
	"github.com/linxlib/logs"
)

type Register uint32
type Interval uint16
type Index uint64

var (
	QE = []uint16{
		0x5601, 0x3401, 0x1801, 0x0ac1, 0x0521, 0x0221, 0x5601, 0x5401, 0x4801, 0x3801, 0x3001, 0x2401,
		0x1c01, 0x1601, 0x5601, 0x5401, 0x5101, 0x4801, 0x3801, 0x3401, 0x3001, 0x2801, 0x2401, 0x2201,
		0x1c01, 0x1801, 0x1601, 0x1401, 0x1201, 0x1101, 0x0ac1, 0x09c1, 0x08a1, 0x0521, 0x0441, 0x02a1,
		0x0221, 0x0141, 0x0111, 0x0085, 0x0049, 0x0025, 0x0015, 0x0009, 0x0005, 0x0001, 0x5601,
	}
	NEXT_MPS = []byte{
		1, 2, 3, 4, 5, 38, 7, 8, 9, 10, 11, 12, 13, 29, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
		27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 45, 46,
	}
	NEXT_LPS = []byte{
		1, 6, 9, 12, 29, 33, 6, 14, 14, 14, 17, 18, 20, 21, 14, 14, 15, 16, 17, 18, 19, 19, 20, 21, 22,
		23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 46,
	}
	SWITCH_LM = []byte{
		1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	CONTEXT_INITIAL = []byte{
		CONTEXT_UNIFORM,
		CONTEXT_RUN_LENGTH,
		CONTEXT_ALL_ZERO_NEIGHBORS,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	}
)

// Table D-7 - Initial states for all contexts
const CONTEXT_UNIFORM uint8 = 46
const CONTEXT_RUN_LENGTH uint8 = 3
const CONTEXT_ALL_ZERO_NEIGHBORS uint8 = 4

func print_c(n uint32) {
	logs.Infof("\"high: 0x%X , low: 0x%X\"", c_high(n), c_low(n))
}

func c_high(n uint32) uint16 {
	return conv.Uint16(n >> 16)
}
func c_low(n uint32) uint16 {
	return conv.Uint16(n << 16 >> 16)
}

func bytein() uint8 {
	return 0
}

type ImageAndTileSizeMarkerSegment struct {
	*boxes.SignatureBox
	DecoderCapabilities []byte
	ReferenceGridWidth  []byte
	ReferenceGridHeight []byte

	ImageHorizontalOffset []byte
	ImageVerticalOffset   []byte
	ReferenceTileWidth    []byte
	ReferenceTileHeight   []byte
	TileHorizontalOffset  []byte
	TileVerticalOffset    []byte
	NoComponents          []byte
	Precision             [][]byte

	HorizontalSeparation [][]byte
	VerticalSeparation   [][]byte
}

func (i *ImageAndTileSizeMarkerSegment) GetDecoderCapabilities() uint16 {
	return binary.BigEndian.Uint16(i.DecoderCapabilities)
}

func (i *ImageAndTileSizeMarkerSegment) GetReferenceGridWidth() uint32 {
	return binary.BigEndian.Uint32(i.ReferenceGridWidth)
}
func (i *ImageAndTileSizeMarkerSegment) GetReferenceGridHeight() uint32 {
	return binary.BigEndian.Uint32(i.ReferenceGridHeight)
}
func (i *ImageAndTileSizeMarkerSegment) GetImageHorizontalOffset() uint32 {
	return binary.BigEndian.Uint32(i.ImageHorizontalOffset)
}
func (i *ImageAndTileSizeMarkerSegment) GetImageVerticalOffset() uint32 {
	return binary.BigEndian.Uint32(i.ImageVerticalOffset)
}
func (i *ImageAndTileSizeMarkerSegment) GetReferenceTileWidth() uint32 {
	return binary.BigEndian.Uint32(i.ReferenceTileWidth)
}
func (i *ImageAndTileSizeMarkerSegment) GetReferenceTileHeight() uint32 {
	return binary.BigEndian.Uint32(i.ReferenceTileHeight)
}
func (i *ImageAndTileSizeMarkerSegment) GetTileHorizontalOffset() uint32 {
	return binary.BigEndian.Uint32(i.TileHorizontalOffset)
}
func (i *ImageAndTileSizeMarkerSegment) GetTileVerticalOffset() uint32 {
	return binary.BigEndian.Uint32(i.TileVerticalOffset)
}

func (i *ImageAndTileSizeMarkerSegment) GetNoComponents() uint16 {
	return binary.BigEndian.Uint16(i.NoComponents)
}
func (i *ImageAndTileSizeMarkerSegment) GetPrecision(i1 int64) int16 {
	precision := i.HorizontalSeparation[i1]
	signedness := precision[0] >> 7
	switch signedness {
	case 0:
		return conv.Int16(uint8(binary.BigEndian.Uint16(precision)))
	case 1:
		return conv.Int16(int8(binary.BigEndian.Uint16(precision)))
	default:
		return conv.Int16(precision[0])
	}
}
