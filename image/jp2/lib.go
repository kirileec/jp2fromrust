package jp2

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/boxes"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
	"github.com/linxlib/logs"
	"io"
)

const (
	COMPRESSION_TYPE_WAVELET           = byte(7)
	CHANNEL_TYPE_COLOR_IMAGE_DATA      = uint16(0)
	CHANNEL_TYPE_OPACITY_DATA          = uint16(1)
	CHANNEL_TYPE_PREMULTIPLIED_OPACITY = uint16(3)
)

type IJBox interface {
	Identifier() box_type.BoxType
	GetLength() uint64
	GetOffset() uint64
	Decode(reader *buffer.ByteBuffer) error
}

type JP2File struct {
	length                uint64
	signature             *boxes.SignatureBox
	fileType              *boxes.FileTypeBox
	header                *boxes.HeaderSuperBox
	contiguousCodeStreams []*boxes.ContiguousCodeStreamBox
	xml                   []*boxes.XMLBox
	uuid                  []*boxes.UUIDBox
}

func (J *JP2File) Length() uint64 {
	return J.length
}

func (J *JP2File) Signature() *boxes.SignatureBox {
	return J.signature
}

func (J *JP2File) GetFileType() *boxes.FileTypeBox {
	return J.fileType
}

func (J *JP2File) GetHeader() *boxes.HeaderSuperBox {
	return J.header
}

func (J *JP2File) ContiguousCodeStreams() []*boxes.ContiguousCodeStreamBox {
	return J.contiguousCodeStreams
}

func (J *JP2File) Xml() []*boxes.XMLBox {
	return J.xml
}

func (J *JP2File) GetUUID() []*boxes.UUIDBox {
	return J.uuid
}

func DecodeJp2(reader *buffer.ByteBuffer) *JP2File {
	boxHeader, _ := box_type.DecodeBoxHeader(reader)

	signatureBox := &boxes.SignatureBox{}
	if !signatureBox.Identifier().Equal(boxHeader.BoxType) {
		logs.Error("box unexpected")
		return nil
	}
	signatureBox.SetLength(boxHeader.BoxLength)
	signatureBox.SetOffset(uint64(reader.ReadOffset()))
	logs.Infof("SignatureBox start at %d", signatureBox.GetLength())
	err := signatureBox.Decode(reader)
	if err != nil {
		logs.Println("解码签名Box失败")
		return nil
	}
	logs.Infof("SignatureBox finish at %d", reader.ReadOffset())

	boxHeader1, _ := box_type.DecodeBoxHeader(reader)
	fileTypeBox := &boxes.FileTypeBox{
		SignatureBox: &boxes.SignatureBox{
			Length: boxHeader1.BoxLength,
			Offset: uint64(reader.ReadOffset()),
		},
		Brand:             make(box_type.BoxType, 4),
		MinVersion:        make(box_type.BoxType, 4),
		CompatibilityList: make([]box_type.BoxType, 0),
	}
	if !fileTypeBox.Identifier().Equal(boxHeader1.BoxType) {
		logs.Error("FileTypeBox box unexpected")
		return nil
	}
	logs.Infof("FileTypeBox start at %d", fileTypeBox.Offset)
	err = fileTypeBox.Decode(reader)
	if err != nil {
		logs.Errorln("解码文件类型Box失败", err)
		return nil
	}
	logs.Infof("FileTypeBox finish at %d", reader.ReadOffset())

	var headerBoxOption *boxes.HeaderSuperBox
	contiguousCodeStreamBoxes := make([]*boxes.ContiguousCodeStreamBox, 0)
	xmlBoxes := make([]*boxes.XMLBox, 0)
	uuidBoxes := make([]*boxes.UUIDBox, 0)
	uuidInfoBoxes := make([]*boxes.UUIDInfoSuperBox, 0)
	var currentUUidInfoBox *boxes.UUIDInfoSuperBox
	logs.Infof("loop to decode boxes")
	for {
		boxHeader2, err := box_type.DecodeBoxHeader(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logs.Errorf("boxHeader2 failed")
			break
		}
		bt := box_type.NewBoxTypes(boxHeader2.BoxType)
		switch bt {
		case box_type.Header:
			logs.Infof("HeaderSuperBox start at %d", reader.ReadOffset())
			headerBox := &boxes.HeaderSuperBox{
				SignatureBox: &boxes.SignatureBox{
					Length: 0,
					Offset: 0,
				},
				ImageHeaderBox: &boxes.ImageHeaderBox{
					SignatureBox:         &boxes.SignatureBox{},
					Height:               make([]byte, 4),
					Width:                make([]byte, 4),
					ComponentsNum:        make([]byte, 2),
					ComponentsBits:       make([]byte, 1),
					CompressionType:      make([]byte, 1),
					ColorsPaceUnknown:    make([]byte, 1),
					IntellectualProperty: make([]byte, 1),
				},
				BitsPerComponentBox: &boxes.BitsPerComponentBox{
					SignatureBox:     &boxes.SignatureBox{},
					ComponentsNum:    0,
					BitsPerComponent: make([]byte, 0),
				},
				ColorSpecificationBoxes: make([]*boxes.ColorSpecificationBox, 0),
				PaletteBox: &boxes.PaletteBox{
					SignatureBox:        &boxes.SignatureBox{},
					NumEntries:          make([]byte, 2),
					NumComponents:       make([]byte, 1),
					GeneratedComponents: make([]types.GeneratedComponent, 0),
				},
				ComponentMappingBox: &boxes.ComponentMappingBox{
					SignatureBox: &boxes.SignatureBox{},
					Mapping:      make([]types.ComponentMap, 0),
				},
				ChannelDefinitionBox: &boxes.ChannelDefinitionBox{
					SignatureBox: &boxes.SignatureBox{},
					Channels:     make([]types.Channel, 0),
				},
				ResolutionBox: &boxes.ResolutionSuperBox{
					SignatureBox:                &boxes.SignatureBox{},
					CaptureResolutionBox:        nil,
					DefaultDisplayResolutionBox: nil,
				},
			}
			headerBox.Length = boxHeader2.BoxLength
			headerBox.Offset = uint64(reader.ReadOffset())
			err := headerBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			headerBoxOption = headerBox
			logs.Infof("HeaderSuperBox finish at %d", reader.ReadOffset())
		case box_type.IntellectualProperty:
			intellectualPropertyBox := &boxes.IntellectualPropertyBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				Data: make([]byte, boxHeader2.BoxLength),
			}
			logs.Infof("IntellectualPropertyBox start at %d", intellectualPropertyBox.Offset)
			err := intellectualPropertyBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			logs.Infof("IntellectualPropertyBox finish at %d", reader.ReadOffset())
		case box_type.XML:
			xmlBox := &boxes.XMLBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				XML: make([]byte, boxHeader2.BoxLength),
			}
			err := xmlBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			xmlBoxes = append(xmlBoxes, xmlBox)

		case box_type.UUID:
			uuidBox := &boxes.UUIDBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
			}
			err := uuidBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			uuidBoxes = append(uuidBoxes, uuidBox)

		case box_type.UUIDInfo:
			uuidInfoBox := &boxes.UUIDInfoSuperBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
			}
			err := uuidInfoBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			if currentUUidInfoBox != nil {
				uuidInfoBoxes = append(uuidInfoBoxes, uuidInfoBox)
			}
			//TODO: clone
			currentUUidInfoBox = uuidInfoBox

		case box_type.UUIDList:

			uuidListBox := &boxes.UUIDListBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				NumberOfUuids: make([]byte, 2),
				IDs:           make([]byte, 16),
			}
			err := uuidListBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			if currentUUidInfoBox != nil {
				currentUUidInfoBox.UUIDList = append(currentUUidInfoBox.UUIDList, uuidListBox)
			} else {
				logs.Error("box missing")
				return nil
			}

		case box_type.DataEntryURL:
			dataEntryUrlBox := &boxes.DataEntryURLBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
				Version:  make([]byte, 1),
				Flags:    make([]byte, 3),
				Location: make([]byte, boxHeader2.BoxLength-4),
			}
			err := dataEntryUrlBox.Decode(reader)
			if err != nil {
				return nil
			}
			if currentUUidInfoBox != nil {
				currentUUidInfoBox.DataEntryUrlBox = append(currentUUidInfoBox.DataEntryUrlBox, dataEntryUrlBox)
			} else {
				logs.Error("box missing")
				return nil
			}

		case box_type.ContiguousCodeStream:
			if headerBoxOption == nil {
				logs.Error("box unexpected")
				return nil
			}

			cBox := &boxes.ContiguousCodeStreamBox{
				SignatureBox: &boxes.SignatureBox{
					Length: boxHeader2.BoxLength,
					Offset: uint64(reader.ReadOffset()),
				},
			}
			err := cBox.Decode(reader)
			if err != nil {
				logs.Error(err)
				return nil
			}
			contiguousCodeStreamBoxes = append(contiguousCodeStreamBoxes, cBox)

		default:
			logs.Errorf("unexpected box type:%d", bt)
			return nil
		}

	}
	if currentUUidInfoBox != nil {
		uuidInfoBoxes = append(uuidInfoBoxes, currentUUidInfoBox)
	}

	jp2 := &JP2File{
		length:                uint64(reader.ReadOffset()),
		signature:             signatureBox,
		fileType:              fileTypeBox,
		header:                headerBoxOption,
		contiguousCodeStreams: contiguousCodeStreamBoxes,
		xml:                   xmlBoxes,
		uuid:                  uuidBoxes,
	}
	return jp2

}
