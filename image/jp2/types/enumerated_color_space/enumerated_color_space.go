package enumerated_color_space

import "gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"

type EnumeratedColorSpace []byte

var (
	ENUMERATED_COLOR_SPACE_UNKNOWN   EnumeratedColorSpace = []byte{0, 0, 0, 0}
	ENUMERATED_COLOR_SPACE_SRGB      EnumeratedColorSpace = []byte{0, 0, 0, 16}
	ENUMERATED_COLOR_SPACE_GREYSCALE EnumeratedColorSpace = []byte{0, 0, 0, 17}
)

type EnumeratedColorSpaces byte

const (
	SRGB EnumeratedColorSpaces = iota
	Greyscale
	Reserved
)

func NewEnumeratedColorSpaces(value []byte) EnumeratedColorSpaces {
	b := box_type.BoxType{}
	b.From(value)
	if b.Equal(ENUMERATED_COLOR_SPACE_SRGB) {
		return SRGB
	} else if b.Equal(ENUMERATED_COLOR_SPACE_GREYSCALE) {
		return Greyscale
	} else {
		return Reserved
	}

}
