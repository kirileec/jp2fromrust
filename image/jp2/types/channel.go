package types

import "encoding/binary"

type Channel struct {
	ChannelIndex       []byte
	ChannelType        []byte
	ChannelAssociation []byte
}

type ChannelTypes uint16

const (
	ColorImageData ChannelTypes = iota
	Opacity
	PremultipliedOpacity
	Reserved
	Unspecified
)

func NewChannelTypes(value []byte) ChannelTypes {
	panic("")
}

func (c *Channel) ChannelTypeU16() uint16 {
	return binary.BigEndian.Uint16(c.ChannelType)
}
func (c *Channel) GetChannelType() ChannelTypes {
	return NewChannelTypes(c.ChannelType)
}
func (c *Channel) GetChannelIndex() uint16 {
	return binary.BigEndian.Uint16(c.ChannelIndex)
}

func (c *Channel) GetChannelAssociation() uint16 {
	return binary.BigEndian.Uint16(c.ChannelAssociation)
}
