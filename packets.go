package ws

import (
	"encoding/json"
	"fmt"

	"github.com/lainio/err2/try"
)

func NewJSONPacket(body []byte) (pkt Packet) {
	pkt = Packet{
		PacketLength:    16 + uint32(len(body)),
		HeaderLength:    16,
		ProtocolVersion: ProtocolJSON,
		Operation:       OpreationUserAuthentication,
		SequenceID:      1,
		Body:            body,
	}
	return
}

func NewPingPacket() Packet {
	return Packet{
		PacketLength:    16,
		HeaderLength:    16,
		ProtocolVersion: ProtocolJSON,
		Operation:       OpreationHeartbeat,
		SequenceID:      1,
	}
}

func NewConnectPacket(roomid int) Packet {
	er := EnterRoom{
		UID:       ref(0),
		RoomID:    roomid,
		ProtoVer:  ref(2),
		Platform:  ref("web"),
		ClientVer: ref("2.0.11"),
		Type:      ref(2),
	}
	body := try.To1(json.Marshal(er))
	return NewJSONPacket(body)
}

func DecodeBody(body []byte) []Packet {
	arr := make([]Packet, 0)
	for {
		h := body[:16]
		if len(h) < 16 {
			break
		}
		pkt := try.To1(Decode(h))
		fmt.Printf("pkts %+v\n", pkt)
	}
	return arr
}

func ref[T any](t T) *T { return &t }
