package backbone

import (
	"github.com/anantadwi13/gorong2/internal/legacy/common"
)

type packet struct {
	id         common.Id
	packetType PacketType
	data       []byte
}

func NewPacket(id common.Id, packetType PacketType, data []byte) (Packet, error) {
	if packetType != PacketTypeOpen &&
		packetType != PacketTypeOpenAckSuccess &&
		packetType != PacketTypeOpenAckFailed &&
		packetType != PacketTypeEdgeCreate &&
		packetType != PacketTypeEdgeCreateAckSuccess &&
		packetType != PacketTypeEdgeCreateAckFailed &&
		packetType != PacketTypeEdgeClose &&
		packetType != PacketTypeMessage {
		return nil, ErrPacketInvalid
	}

	if id == nil {
		id = common.GenerateId()
	}

	return &packet{
		id:         id,
		packetType: packetType,
		data:       data,
	}, nil
}

func (p *packet) Id() common.Id {
	return p.id
}

func (p *packet) Type() PacketType {
	return p.packetType
}

func (p *packet) Data() []byte {
	return p.data
}
