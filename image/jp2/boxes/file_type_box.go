package boxes

import (
	"errors"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type FileTypeBox struct {
	*SignatureBox
	Brand             box_type.BoxType
	MinVersion        box_type.BoxType
	CompatibilityList box_type.CompatibilityList
}

func (f *FileTypeBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_FILE_TYPE
}

func (f *FileTypeBox) GetLength() uint64 {
	return f.Length
}

func (f *FileTypeBox) GetOffset() uint64 {
	return f.Offset
}

func (f *FileTypeBox) Decode(reader *buffer.ByteBuffer) error {
	var buf = make([]byte, 4)
	_, err := reader.Read(buf)
	if err != nil {
		return err
	}
	f.Brand.From(buf)
	if box_type.BRAND_JPX.Equal(buf) {
		return errors.New("BRAND_JPX is not supported")
	}
	if !box_type.BRAND_JP2.Equal(buf) {
		return errors.New("invalid Brand")
	}

	_, err = reader.Read(f.MinVersion)
	if err != nil {
		return err
	}
	size := (f.Length - 8) / 4
	for size > 0 {
		_, err := reader.Read(buf)
		if err != nil {
			return err
		}
		f.CompatibilityList = box_type.ExtendFromSlice(f.CompatibilityList, buf)
		size -= 1
	}
	flag := false
	for _, boxType := range f.CompatibilityList {
		if box_type.BRAND_JP2.EqualA(boxType) {
			flag = true
		}
	}
	if !flag {
		return errors.New("not compatible")
	}

	return nil
}
