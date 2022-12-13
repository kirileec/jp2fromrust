package boxes

import (
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type BitsPerComponentBox struct {
	*SignatureBox
	ComponentsNum    uint16
	BitsPerComponent []uint8
}

func NewBitsPerComponentBox(signatureBox *SignatureBox, componentsNum uint16, bitsPerComponent []uint8) *BitsPerComponentBox {
	return &BitsPerComponentBox{SignatureBox: signatureBox, ComponentsNum: componentsNum, BitsPerComponent: bitsPerComponent}
}

func (b *BitsPerComponentBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_BITS_PER_COMPONENT
}

func (b *BitsPerComponentBox) GetLength() uint64 {
	return b.Length
}

func (b *BitsPerComponentBox) GetOffset() uint64 {
	return b.Offset
}

func (b *BitsPerComponentBox) Decode(reader *buffer.ByteBuffer) error {
	_, err := reader.Read(b.BitsPerComponent)
	if err != nil {
		return err
	}
	return nil
}

func (b *BitsPerComponentBox) GetBotsPerComponent() []types.BitDepth {
	ll := make([]types.BitDepth, 0)
	for _, u := range b.BitsPerComponent {
		a := types.NewBitDepth(u)
		ll = append(ll, a)
	}
	return ll
}
