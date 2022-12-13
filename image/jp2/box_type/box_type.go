package box_type

type BoxType []byte

func (t BoxType) Equal(other []byte) bool {
	flag := 0
	for i := 0; i < 4; i++ {
		tmp := other[i] == t[i]
		if tmp {
			flag++
		}
	}
	return flag == 4
}
func (t BoxType) EqualA(other []byte) bool {
	flag := false
	for i := 0; i < 4; i++ {
		flag = other[i] == t[i]
	}
	return flag
}

func (t BoxType) From(other []byte) {
	for i := 0; i < 4; i++ {
		t[i] = other[i]
	}
}

// jP\040\040 (0x6A50 2020)
var (
	BOX_TYPE_SIGNATURE                  = BoxType{106, 80, 32, 32}
	BOX_TYPE_FILE_TYPE                  = BoxType{102, 116, 121, 112}
	BOX_TYPE_HEADER                     = BoxType{106, 112, 50, 104}
	BOX_TYPE_IMAGE_HEADER               = BoxType{105, 104, 100, 114}
	BOX_TYPE_BITS_PER_COMPONENT         = BoxType{98, 112, 99, 99}
	BOX_TYPE_COLOR_SPECIFICATION        = BoxType{99, 111, 108, 114}
	BOX_TYPE_PALETTE                    = BoxType{112, 99, 108, 114}
	BOX_TYPE_COMPONENT_MAPPING          = BoxType{99, 109, 97, 112}
	BOX_TYPE_CHANNEL_DEFINITION         = BoxType{99, 100, 101, 102}
	BOX_TYPE_RESOLUTION                 = BoxType{114, 101, 115, 32}
	BOX_TYPE_CAPTURE_RESOLUTION         = BoxType{114, 101, 115, 99}
	BOX_TYPE_DEFAULT_DISPLAY_RESOLUTION = BoxType{114, 101, 115, 100}
	BOX_TYPE_CONTIGUOUS_CODESTREAM      = BoxType{106, 112, 50, 99}
	BOX_TYPE_INTELLECTUAL_PROPERTY      = BoxType{106, 112, 50, 105}
	BOX_TYPE_XML                        = BoxType{120, 109, 108, 32}
	BOX_TYPE_UUID                       = BoxType{117, 117, 105, 100}
	BOX_TYPE_UUID_INFO                  = BoxType{117, 105, 110, 102}
	BOX_TYPE_UUID_LIST                  = BoxType{117, 108, 115, 116}
	BOX_TYPE_DATA_ENTRY_URL             = BoxType{117, 114, 108, 32}
)

type BoxTypes int

func (b BoxTypes) String() string {
	switch b {
	case Signature:
		return "Signature"
	case FileType:
		return "FileType"
	case Header:
		return "Header"
	case ImageHeader:
		return "ImageHeader"
	case BitsPerComponent:
		return "BitsPerComponent"
	case ColorSpecification:
		return "ColorSpecification"
	case Palette:
		return "Palette"
	case ComponentMapping:
		return "ComponentMapping"
	case ChannelDefinition:
		return "ChannelDefinition"
	case Resolution:
		return "Resolution"
	case CaptureResolution:
		return "CaptureResolution"
	case DefaultDisplayResolution:
		return "DefaultDisplayResolution"
	case ContiguousCodeStream:
		return "ContiguousCodeStream"
	case IntellectualProperty:
		return "GetIntellectualProperty"
	case XML:
		return "XML"
	case UUID:
		return "UUID"
	case UUIDInfo:
		return "UUIDInfo"
	case UUIDList:
		return "UUIDList"
	case DataEntryURL:
		return "DataEntryURL"
	case Unknown:
		return "Unknown"
	default:
		return "xxxxx"
	}
}

const (
	Signature BoxTypes = iota
	FileType
	Header
	ImageHeader
	BitsPerComponent
	ColorSpecification
	Palette
	ComponentMapping
	ChannelDefinition
	Resolution
	CaptureResolution
	DefaultDisplayResolution
	ContiguousCodeStream
	IntellectualProperty
	XML
	UUID
	UUIDInfo
	UUIDList
	DataEntryURL
	Unknown
)

func NewBoxTypes(value BoxType) BoxTypes {
	if value.Equal(BOX_TYPE_SIGNATURE) {
		return Signature
	} else if value.Equal(BOX_TYPE_FILE_TYPE) {
		return FileType
	} else if value.Equal(BOX_TYPE_HEADER) {
		return Header
	} else if value.Equal(BOX_TYPE_IMAGE_HEADER) {
		return ImageHeader
	} else if value.Equal(BOX_TYPE_BITS_PER_COMPONENT) {
		return BitsPerComponent
	} else if value.Equal(BOX_TYPE_COLOR_SPECIFICATION) {
		return ColorSpecification
	} else if value.Equal(BOX_TYPE_PALETTE) {
		return Palette
	} else if value.Equal(BOX_TYPE_COMPONENT_MAPPING) {
		return ComponentMapping
	} else if value.Equal(BOX_TYPE_CHANNEL_DEFINITION) {
		return ChannelDefinition
	} else if value.Equal(BOX_TYPE_RESOLUTION) {
		return Resolution
	} else if value.Equal(BOX_TYPE_CAPTURE_RESOLUTION) {
		return CaptureResolution
	} else if value.Equal(BOX_TYPE_DEFAULT_DISPLAY_RESOLUTION) {
		return DefaultDisplayResolution
	} else if value.Equal(BOX_TYPE_CONTIGUOUS_CODESTREAM) {
		return ContiguousCodeStream
	} else if value.Equal(BOX_TYPE_INTELLECTUAL_PROPERTY) {
		return IntellectualProperty
	} else if value.Equal(BOX_TYPE_XML) {
		return XML
	} else if value.Equal(BOX_TYPE_UUID) {
		return UUID
	} else if value.Equal(BOX_TYPE_UUID_INFO) {
		return UUIDInfo
	} else if value.Equal(BOX_TYPE_UUID_LIST) {
		return UUIDList
	} else if value.Equal(BOX_TYPE_DATA_ENTRY_URL) {
		return DataEntryURL
	} else {
		return Unknown
	}
}
