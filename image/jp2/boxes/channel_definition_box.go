package boxes

import (
	"encoding/binary"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/box_type"
	"gitee.com/kayi-cloud/gdimse/example/image/jp2/types"
	"gitee.com/kayi-cloud/gdimse/lib/buffer"
)

type ChannelDefinitionBox struct {
	*SignatureBox
	Channels []types.Channel
}

func (c *ChannelDefinitionBox) Identifier() box_type.BoxType {
	return box_type.BOX_TYPE_CHANNEL_DEFINITION
}

func (c *ChannelDefinitionBox) GetLength() uint64 {
	return c.Length
}

func (c *ChannelDefinitionBox) GetOffset() uint64 {
	return c.Offset
}

func (c *ChannelDefinitionBox) Decode(reader *buffer.ByteBuffer) error {
	noChannelDescriptions := make([]byte, 2)
	reader.Read(noChannelDescriptions)

	size := binary.BigEndian.Uint16(noChannelDescriptions)
	channels := make([]types.Channel, size)
	for size > 0 {
		channel := types.Channel{
			ChannelIndex:       make([]byte, 2),
			ChannelType:        make([]byte, 2),
			ChannelAssociation: make([]byte, 2),
		}
		reader.Read(channel.ChannelIndex)
		reader.Read(channel.ChannelType)
		reader.Read(channel.ChannelAssociation)

		channels = append(channels, channel)

		size -= 1
	}
	c.Channels = channels
	return nil

}
