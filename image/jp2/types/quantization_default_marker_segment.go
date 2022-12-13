package types

import "gitee.com/kayi-cloud/gdimse/example/image/jp2/types/quantization_style"

type QuantizationDefaultMarkerSegment struct {
	length         []uint8
	style          []uint8
	stepSizeValues []uint8
}

func (q *QuantizationDefaultMarkerSegment) Length() uint16 {
	panic("")
}

func (q *QuantizationDefaultMarkerSegment) Style() quantization_style.QuantizationStyle {
	return quantization_style.NewQuantizationStyle(q.style[0])
}

func (q *QuantizationDefaultMarkerSegment) StepSizeValues() {
	panic("")
}

func (q *QuantizationDefaultMarkerSegment) NoGuardbits() uint8 {
	return q.style[0] >> 5
}
