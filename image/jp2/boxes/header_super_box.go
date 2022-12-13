package boxes

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types/enumerated_color_space"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/logs"
	"io"
)

type HeaderSuperBox struct {
	*SignatureBox
	ImageHeaderBox          *ImageHeaderBox
	BitsPerComponentBox     *BitsPerComponentBox
	ColorSpecificationBoxes []*ColorSpecificationBox
	PaletteBox              *PaletteBox
	ComponentMappingBox     *ComponentMappingBox
	ChannelDefinitionBox    *ChannelDefinitionBox
	ResolutionBox           *ResolutionSuperBox
}

func (h *HeaderSuperBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_HEADER
}

func (h *HeaderSuperBox) Decode(reader *buffer.ByteBuffer) error {
	boxHeader, _ := box_type.DecodeBoxHeader(reader)
	if !h.ImageHeaderBox.Identifier().Equal(boxHeader.BoxType) {
		return errors.New("box unexpected")
	}
	h.ImageHeaderBox.Length = boxHeader.BoxLength
	h.ImageHeaderBox.Offset = uint64(reader.ReadOffset())
	err := h.ImageHeaderBox.Decode(reader)
	if err != nil {
		return err
	}

	for {
		logs.Infof("currnet offset: %d", reader.ReadOffset())
		boxHeader1, _ := box_type.DecodeBoxHeader(reader)
		logs.Info(boxHeader1)
		bt := box_type.NewBoxTypes(boxHeader1.BoxType)
		switch bt {
		case box_type.ImageHeader:
			logs.Warn("ImageHeaderBox found in other place, ignoring")
		case box_type.ColorSpecification:
			colorSpecificationBox := &ColorSpecificationBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				Method:                  make([]byte, 1),
				Precedence:              make([]byte, 1),
				ColorSpaceApproximation: make([]byte, 1),
				EnumeratedColorSpace:    enumerated_color_space.ENUMERATED_COLOR_SPACE_UNKNOWN,
				RestrictedIccProfile:    make([]byte, 0),
			}
			logs.Infof("decode color specification box start at %d", colorSpecificationBox.Offset)
			err := colorSpecificationBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode color specification box finish at %d", reader.ReadOffset())
			h.ColorSpecificationBoxes = append(h.ColorSpecificationBoxes, colorSpecificationBox)
		case box_type.BitsPerComponent:
			if h.BitsPerComponentBox != nil {
				return errors.New("box duplicate")
			}
			componentsNum := h.ImageHeaderBox.GetComponentsNum()
			bitsPerComponentBox := &BitsPerComponentBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				ComponentsNum:    componentsNum,
				BitsPerComponent: make([]byte, componentsNum),
			}
			logs.Infof("decode BitsPerComponentBox start at %d", bitsPerComponentBox.Offset)
			err := bitsPerComponentBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode BitsPerComponentBox finish at %d", reader.ReadOffset())
			h.BitsPerComponentBox = bitsPerComponentBox

		case box_type.Palette:
			if h.PaletteBox != nil {
				return errors.New("box duplicate")
			}
			paletteBox := &PaletteBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				NumEntries:          make([]byte, 2),
				NumComponents:       make([]byte, 1),
				GeneratedComponents: make([]types.GeneratedComponent, 0),
			}
			logs.Infof("decode PaletteBox start at %d", paletteBox.Offset)
			err := paletteBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode PaletteBox finish at %d", reader.ReadOffset())
		case box_type.ComponentMapping:
			if h.ComponentMappingBox != nil {
				return errors.New("box duplicate")
			}
			componentMappingBox := &ComponentMappingBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				Mapping: make([]types.ComponentMap, 0),
			}
			logs.Infof("decode ComponentMappingBox start at %d", componentMappingBox.Offset)
			err := componentMappingBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode ComponentMappingBox finish at %d", reader.ReadOffset())
			h.ComponentMappingBox = componentMappingBox
		case box_type.ChannelDefinition:
			if h.ChannelDefinitionBox != nil {
				return errors.New("box duplicate")
			}
			channelDefinitionBox := &ChannelDefinitionBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				Channels: make([]types.Channel, 0),
			}
			logs.Infof("decode ChannelDefinitionBox start at %d", channelDefinitionBox.Offset)
			err := channelDefinitionBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode ChannelDefinitionBox finish at %d", reader.ReadOffset())
			h.ChannelDefinitionBox = channelDefinitionBox
		case box_type.Resolution:
			if h.ResolutionBox != nil {
				return errors.New("box duplicate")
			}
			resolutionBox := &ResolutionSuperBox{
				SignatureBox: &SignatureBox{
					Length: boxHeader1.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				CaptureResolutionBox:        nil,
				DefaultDisplayResolutionBox: nil,
			}
			logs.Infof("decode ResolutionBox start at %d", resolutionBox.Offset)
			err := resolutionBox.Decode(reader)
			if err != nil {
				return err
			}
			logs.Infof("decode ResolutionBox finish at %d", reader.ReadOffset())
			h.ResolutionBox = resolutionBox
		case box_type.Unknown:
			logs.Warn("unknown box type")
			goto BREAK
		default: // End of header but recognised new box type
			logs.Infof("currnet offset: %d %d", reader.ReadOffset(), bt)
			_, err := reader.Seek(-(int64(boxHeader1.HeaderLength)), io.SeekCurrent)
			if err != nil {
				return err
			}
			logs.Infof("currnet offset: %d", reader.ReadOffset())
			goto BREAK
		}

	}
BREAK:

	if len(h.ColorSpecificationBoxes) == 0 {
		return errors.New("box malformed")
	}
	return nil
}
