package boxes

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"io"
)

type ResolutionSuperBox struct {
	*SignatureBox

	CaptureResolutionBox        *CaptureResolutionBox
	DefaultDisplayResolutionBox *DefaultDisplayResolutionBox
}

func (s *ResolutionSuperBox) GetCaptureResolutionBox() *CaptureResolutionBox {
	return s.CaptureResolutionBox
}

func (s *ResolutionSuperBox) GetDefaultDisplayResolutionBox() *DefaultDisplayResolutionBox {
	return s.DefaultDisplayResolutionBox
}

func (s *ResolutionSuperBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_RESOLUTION
}

func (s *ResolutionSuperBox) Decode(reader *buffer.ByteBuffer) error {
	boxHeader, _ := box_type.DecodeBoxHeader(reader)
	switch box_type.NewBoxTypes(boxHeader.BoxType) {
	case box_type.CaptureResolution:
		if s.CaptureResolutionBox != nil {
			return errors.New("Box Unexpected")
		}
		captureResolution := &CaptureResolutionBox{
			SignatureBox: &SignatureBox{
				Length: boxHeader.BoxLength,
				Offset: uint64(reader.ReadOffset()),
			},
		}
		captureResolution.Decode(reader)
		s.CaptureResolutionBox = captureResolution

	case box_type.DefaultDisplayResolution:
		if s.DefaultDisplayResolutionBox != nil {
			return errors.New("box unexpected")
		}
		defaultDisplayResolutionBox := &DefaultDisplayResolutionBox{
			SignatureBox: &SignatureBox{
				Length: boxHeader.BoxLength,
				Offset: uint64(reader.ReadOffset()),
			},
		}

		defaultDisplayResolutionBox.Decode(reader)
		s.DefaultDisplayResolutionBox = defaultDisplayResolutionBox

	default:
		reader.Seek(-int64(boxHeader.HeaderLength), io.SeekCurrent)
	}
	if s.CaptureResolutionBox == nil && s.DefaultDisplayResolutionBox == nil {
		return errors.New("box malformed")
	}
	return nil
}
