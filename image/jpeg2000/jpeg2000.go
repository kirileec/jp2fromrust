package jpeg2000

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"strconv"
)

// A FormatError reports that the input is not a valid JPEG2000.
type FormatError string

func (e FormatError) Error() string { return "invalid JPEG2000 format: " + string(e) }

// An UnsupportedError reports that the input uses a valid but unimplemented JPEG feature.
type UnsupportedError string

func (e UnsupportedError) Error() string { return "unsupported JPEG2000 feature: " + string(e) }

type bits struct {
	a uint32
	m uint32
	n int32
}

type siz struct {
	profile        uint16
	width          uint32
	height         uint32
	offsetX        uint32
	offsetY        uint32
	tile_width     uint32
	tile_height    uint32
	tile_offsetX   uint32
	tile_offsetY   uint32
	num_components uint16
	component_size int
}

type cod struct {
	scod             byte
	prog_order       byte
	quality_layers   uint16
	mct              uint8 // 0 or 1
	decomp_levels    uint8
	codeblock_width  int
	codeblock_height int
	codeblock_style  int
	wavelet_filter   int // 0 or 1
	precints         int // optional
}

type ihdr struct {
	vers   uint16 // major / minor version
	nc     uint16 // number of components; non-0
	height uint32 // image  height; non-0
	width  uint32 // image width; non-0
	bpc    uint8  // bits per component; signed
	c      uint8  // compression type; must be 7
	unkc   uint8  // 0 known; 1 unknown
	ipr    uint8  // 0 or 1
}

type colr struct {
	meth    uint   // 1 - Enumerated; 2 - Restricted ICC profile
	prec    uint   // 0; precedence
	approx  uint   // 0; color space approximation
	enumCs  uint32 // Enumereated color space meth 1 only; 16 = sRGB; 17 = greyscale
	profile []byte // ICC profile
}

type ftyp struct {
	br   string
	minV uint32
	cl   []string
}

type decoder struct {
	r             io.Reader
	bits          bits
	width, height int
	tmp           [4096]byte
	ihdr          ihdr
	colr          colr
	ftyp          ftyp
}

type j2k_box struct {
	r        io.Reader
	box_type string
	length   int64
	tmp      [4096]byte
}

type j2k_segment struct {
	r            io.Reader
	segment_type string
	length       int64
	siz          siz
	cod          cod
}

func (b *j2k_box) reader(r io.Reader) io.Reader {
	if b.length == 0 {
		return r
	}

	return io.LimitReader(r, b.length)
}

func (b *j2k_segment) reader(r io.Reader) io.Reader {
	return io.LimitReader(r, b.length)
}

func debugOffset(r io.Reader) {
	switch r := r.(type) {
	case io.Seeker:
		pos, _ := r.Seek(0, io.SeekCurrent)
		fmt.Println("@ Offset " + strconv.FormatInt(pos, 10))
	default:

	}
}

func (d *decoder) decode(r io.Reader, configOnly bool) (image.Image, error) {
	d.r = r
	debugOffset(d.r)
	box, err := d.nextBox()

	if err != nil {
		return nil, err
	}

	if box.length != 4 && box.box_type != string(rfc3745Magic[4:8]) {
		return nil, FormatError("expected a signature box; found " + box.box_type)
	}

	// check for the JP2 header
	if _, err := io.ReadFull(d.r, d.tmp[:box.length]); err != nil {
		return nil, err
	}

	if string(d.tmp[:box.length]) != string(rfc3745Magic[8:12]) {
		return nil, FormatError("couldn't find magic number in signature")
	}

	// format:
	//   4 bytes (length)
	//   ^length bytes (first 4 are the box type, remaining are the goods)
	// repeat forever
	for {
		box, err := d.nextBox()

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		fmt.Println("JP2 box: type=\"" + box.box_type + "\"; length=" + strconv.FormatInt(box.length, 10))
		switch box.box_type {
		case "ftyp":
			if _, err := io.ReadFull(box.r, d.tmp[:box.length]); err != nil {
				return nil, err
			}

			d.ftyp.br = string(d.tmp[:4])
			if d.ftyp.br != "jp2\x20" {
				return nil, FormatError("couldn't find magic string in profile box BR")
			}

			d.ftyp.minV = binary.BigEndian.Uint32(d.tmp[4:8])
			// expecting 0

			found := false

			for i := int64(2); i < (box.length / 4); i++ {
				str := string(d.tmp[4*i : 4*(i+1)])
				d.ftyp.cl = append(d.ftyp.cl, str)
				found = found || str == "jp2\x20"
			}

			if found == false {
				return nil, FormatError("couldn't find magic string in profile box CL")
			}
			fmt.Printf("FTYP:  %+v\n", d.ftyp)

		case string(RESCBox):
		case string(RESDBox):
		case string(JP2HBox):
			// super box!
			for {
				innerBox, err := box.nextBox()

				if err != nil {
					if err == io.EOF {
						break
					}
				}
				fmt.Println("JP2H box: type=\"" + innerBox.box_type + "\"; length=" + strconv.FormatInt(innerBox.length, 10))
				switch innerBox.box_type {
				case string(IHDRBox):
					if _, err := io.ReadFull(innerBox.r, d.tmp[:innerBox.length]); err != nil {
						return nil, err
					}

					d.ihdr.vers = binary.BigEndian.Uint16(d.tmp[0:2])
					d.ihdr.nc = binary.BigEndian.Uint16(d.tmp[2:4])
					d.ihdr.height = binary.BigEndian.Uint32(d.tmp[4:8])
					d.ihdr.width = binary.BigEndian.Uint32(d.tmp[8:12])
					d.ihdr.bpc = uint8(d.tmp[12])
					d.ihdr.c = uint8(d.tmp[13])
					d.ihdr.unkc = uint8(d.tmp[14])
					d.ihdr.ipr = uint8(d.tmp[15])
					fmt.Printf("IHDR:  %+v\n", d.ihdr)
				case string(BPCCBox):
					// duplicated by IHDR?
					// 1 byte for each component
				case string(COLRBox):
					// ignore duplicates

					if _, err := io.ReadFull(innerBox.r, d.tmp[:3]); err != nil {
						return nil, err
					}

					d.colr.meth = uint(d.tmp[0])
					d.colr.prec = uint(d.tmp[1])
					d.colr.approx = uint(d.tmp[2])

					if d.colr.meth == 1 {

					} else if d.colr.meth == 2 {

					}
					// meth = 1 (CIELab)
					//   4 bytes each of rl, ol, ra, oa, rb, ob, il
					// meth = 2 (ICC)
					// read the rest; 1 byte for every icc value
					// meth (ignore, read the rest)
					fmt.Printf("COLR:  %+v\n", d.colr)
				case string(PCLRBox):
					// 2 bytes for num entries
					// 1 byte NPC (channel)
					// for each channel, 1 byte (Bi) for size
					// for each entry, for each channel, bytes given by above
				case string(CMAPBox):
					// as most 1.
					// need PCLR first
					// for each channel:
					//    - 2 bytes of CMP^i
					//    - 1 byte of MTYPE^i
					//    - 1 byte of PCOL^i
				case string(CDEFBox):
					// at  most 1.
					// 2 bytes of N
					// for each:
					//   - 2 bytes of Cn^i
					//   - 2 bytes of  Typ^i
					//   - 2 bytes of Asoc^i
				case string(RESBox):
				}
				innerBox.consume()
			}
		case string(JP2IBox):
		case string(XMLBox):
		case string(UUIDBox):
			if _, err := io.ReadFull(box.r, d.tmp[:16]); err != nil {
				return nil, err
			}
		case string(UINFBox):
		case string(ULSTBox):
		case string(URLBox):
		case string(JP2CBox):
			// soc
			if _, err := io.ReadFull(box.r, d.tmp[:2]); err != nil {
				return nil, err
			}

			if string(d.tmp[:2]) != "\xFF\x4F" {
				return nil, FormatError("No SOC")
			}

			for {
				segment, err := box.nextSegment()

				if err != nil {
					if err == io.EOF {
						break
					}
				}
				fmt.Printf("Segment type=0x%x; length=%d\n", segment.segment_type, segment.length)
				psot := int64(0)
				switch segment.segment_type {
				case "\xFF\x51": // SIZ
					if _, err := io.ReadFull(segment.r, d.tmp[:segment.length]); err != nil {
						return nil, err
					}

					segment.siz.profile = binary.BigEndian.Uint16(d.tmp[0:2])
					// width (4 bytes)
					segment.siz.width = binary.BigEndian.Uint32(d.tmp[2:6])
					// height (4 bytes)
					segment.siz.height = binary.BigEndian.Uint32(d.tmp[6:10])
					segment.siz.offsetX = binary.BigEndian.Uint32(d.tmp[10:14])
					segment.siz.offsetY = binary.BigEndian.Uint32(d.tmp[14:18])
					// tile width (4 bytes)
					segment.siz.tile_width = binary.BigEndian.Uint32(d.tmp[18:22])
					// tile height (4 bytes)
					segment.siz.tile_height = binary.BigEndian.Uint32(d.tmp[22:26])
					segment.siz.tile_offsetX = binary.BigEndian.Uint32(d.tmp[26:30])
					segment.siz.tile_offsetY = binary.BigEndian.Uint32(d.tmp[30:34])
					// num components (2 bytes)
					segment.siz.num_components = binary.BigEndian.Uint16(d.tmp[34:36])
					// for each component: precision (1 byte), subsampling x (1 byte), subsampling y (1 byte)
					// component size (1 byte)
					segment.siz.component_size = int(d.tmp[36]) + 1
					fmt.Printf("SIZ:  %+v\n", segment.siz)
				case "\xFF\x52": // COD
					if _, err := io.ReadFull(segment.r, d.tmp[:segment.length]); err != nil {
						return nil, err
					}
					segment.cod.scod = d.tmp[0]
					segment.cod.prog_order = d.tmp[1]
					segment.cod.quality_layers = binary.BigEndian.Uint16(d.tmp[2:4])
					segment.cod.mct = uint8(d.tmp[4])
					segment.cod.decomp_levels = uint8(d.tmp[5])
					segment.cod.codeblock_width = int(d.tmp[6])
					segment.cod.codeblock_height = int(d.tmp[7])
					segment.cod.codeblock_style = int(d.tmp[8])
					segment.cod.codeblock_style = int(d.tmp[8])
					segment.cod.wavelet_filter = int(d.tmp[9])
					fmt.Printf("COD:  %+v\n", segment.cod)
					// maybe: segment.cod.precints = int(d.tmp[10])
				case "\xFF\x53": // COC
				case "\xFF\x5C": // QCD
				case "\xFF\x5D": // QCC
				case "\xFF\x90": // SOT
					fmt.Println("[SOT]")
					if _, err := io.ReadFull(segment.r, d.tmp[:segment.length]); err != nil {
						return nil, err
					}
					fmt.Println(d.tmp[:segment.length])
					psot = int64(binary.BigEndian.Uint32(d.tmp[2:6]))
					fmt.Println("PSot:" + strconv.FormatInt(psot, 10))
				case "\xFF\x93": // SOT
					fmt.Println("[SOD]")
					segment.consume()
					box.consume()
					d.seek(psot - 2)
					break
					// if _, err := io.ReadFull(box.r, d.tmp[:8]); err != nil {
					//   return nil, err
					// }
					// fmt.Println(d.tmp[:8])

					// if _, err := io.ReadFull(box.r, d.tmp[:2]); err != nil {
					//   return nil, err
					// }
					//
					// if (string(d.tmp[:2]) == "\xFF\x90") {
					//   fmt.Println("[SOT]")
					// }

					// Lsot 2 bytes (already accounted for by segment length)
					// Isot 2 bytes
					// Psot 32bits
					// TPsot 8bits
					// TNsot 8bits
					// possibly some other headers
					// SOD
					//
					// READ TO SOD
					// look for optional markers
					// SOD means the rest of the segment is data
				default:
				}

				segment.consume()
			}

			// components
			// 2 bytes of name
			// 2 bytes of length
			// stuff

			// must be 38+ bytes

			// cod and/or coc
			// some other segment type

			switch r := d.r.(type) {
			case io.Seeker:
				readBytes, _ := r.Seek(0, io.SeekEnd)
				fmt.Println(readBytes)
			default:
				readBytes, _ := io.Copy(ioutil.Discard, d.r)
				fmt.Println(readBytes)
			}
		default:
		}

		if _, err := box.consume(); err != nil {
			return nil, nil
		}
	}

	return nil, nil
}

func (b *j2k_box) consume() (int64, error) {
	switch r := b.r.(type) {
	case io.Seeker:
		return r.Seek(0, io.SeekEnd)
	default:
		return io.Copy(ioutil.Discard, b.r)
	}
}

func (b *j2k_segment) consume() (int64, error) {
	switch r := b.r.(type) {
	case io.Seeker:
		return r.Seek(0, io.SeekEnd)
	default:
		return io.Copy(ioutil.Discard, b.r)
	}
}

func (d *decoder) seek(length int64) (int64, error) {
	switch r := d.r.(type) {
	case io.Seeker:
		return r.Seek(length, io.SeekCurrent)
	default:
		return io.CopyN(ioutil.Discard, d.r, length)
	}
}

func (d *j2k_box) seek(length int64) (int64, error) {
	switch r := d.r.(type) {
	case io.Seeker:
		return r.Seek(length, io.SeekCurrent)
	default:
		return io.CopyN(ioutil.Discard, d.r, length)
	}
}

func (d *decoder) nextBox() (j2k_box, error) {
	var box j2k_box

	debugOffset(d.r)

	n, err := io.ReadFull(d.r, d.tmp[:8])
	if err != nil {
		return box, err
	}

	if n == 0 {
		return box, nil
	}

	declaredLength := int32(binary.BigEndian.Uint32(d.tmp[0:4]))
	box.box_type = string(d.tmp[4:8])

	if declaredLength == 0 {
		box.length = 0
	} else if declaredLength == 1 {
		io.ReadFull(d.r, d.tmp[8:16])
		res := binary.BigEndian.Uint64(d.tmp[8:16])
		box.length = int64(res) - 8
	} else {
		box.length = int64(declaredLength) - 8
	}

	box.r = box.reader(d.r)

	return box, nil
}

func (d *j2k_box) nextBox() (j2k_box, error) {
	var box j2k_box
	debugOffset(d.r)

	n, err := io.ReadFull(d.r, d.tmp[:8])
	if err != nil {
		return box, err
	}

	if n == 0 {
		return box, nil
	}

	declaredLength := int32(binary.BigEndian.Uint32(d.tmp[0:4]))
	box.box_type = string(d.tmp[4:8])

	if declaredLength == 0 {
		box.length = 0
	} else if declaredLength == 1 {
		io.ReadFull(d.r, d.tmp[8:16])
		res := binary.BigEndian.Uint64(d.tmp[8:16])
		box.length = int64(res) - 8
	} else {
		box.length = int64(declaredLength) - 8
	}

	box.r = box.reader(d.r)

	return box, nil
}

func (d *j2k_box) nextSegment() (j2k_segment, error) {
	var segment j2k_segment

	debugOffset(d.r)
	n, err := io.ReadFull(d.r, d.tmp[:2])

	if err != nil {
		return segment, err
	}

	if n == 0 {
		return segment, nil
	}

	segment.segment_type = string(d.tmp[0:2])

	if segment.segment_type == "\xff\x93" {
		segment.length = 0
		segment.r = segment.reader(d.r)
		return segment, nil
	}

	n, err = io.ReadFull(d.r, d.tmp[2:4])

	if err != nil {
		return segment, err
	}

	if n == 0 {
		return segment, nil
	}
	segment.length = int64(binary.BigEndian.Uint16(d.tmp[2:4])) - 2

	segment.r = segment.reader(d.r)

	return segment, nil
}

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	return d.decode(r, false)
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, nil
}

var (
	rfc3745Magic = []byte("\x00\x00\x00\x0c\x6a\x50\x20\x20\x0d\x0a\x87\x0a")
	jp2Magic     = []byte("\x0d\x0a\x87\x0a")
	j2kMagic     = []byte("\xff\x4f\xff\x51")
	IHDRBox      = []byte("\x69\x68\x64\x72")
	BPCCBox      = []byte("\x62\x70\x63\x63") /* Bits Per Component */
	COLRBox      = []byte("\x63\x6f\x6c\x72") /* Color Specification */
	PCLRBox      = []byte("\x70\x63\x6c\x72") /* Palette */
	CMAPBox      = []byte("\x63\x6d\x61\x70") /* Component Mapping */
	CDEFBox      = []byte("\x63\x64\x65\x66") /* Channel Definition */
	RESBox       = []byte("\x72\x65\x73\x20") /* Resolution */
	RESCBox      = []byte("\x72\x65\x73\x63") /* Capture Resolution */
	RESDBox      = []byte("\x72\x65\x73\x64") /* Default Display Resolution */
	JP2CBox      = []byte("\x6a\x70\x32\x63") /* Contiguous Code Stream */
	JP2HBox      = []byte("\x6a\x70\x32\x68") /* JP2 Header */
	JP2IBox      = []byte("\x6a\x70\x32\x69") /* Intellectual Property */
	XMLBox       = []byte("\x78\x6d\x6c\x20") /* XML */
	UUIDBox      = []byte("\x75\x75\x69\x64") /* UUID */
	UINFBox      = []byte("\x75\x69\x6e\x66") /* UUID Info */
	ULSTBox      = []byte("\x75\x63\x73\x74") /* UUID List */
	URLBox       = []byte("\x75\x72\x6c\x20") /* URL */
)

func init() {
	image.RegisterFormat("jpeg2000", string(rfc3745Magic), Decode, DecodeConfig)
	// image.RegisterFormat("jpeg2000", string(jp2Magic), Decode, DecodeConfig)
	// image.RegisterFormat("jpeg2000", string(j2kMagic), Decode, DecodeConfig)
}
