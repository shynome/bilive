package bilive

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/google/brotli/go/cbrotli"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

// 参考自 https://github.com/lovelyyoshino/Bilibili-Live-API/blob/master/API.WebSocket.md
type Packet struct {
	PacketLength    uint32          // 数据包长度
	HeaderLength    uint16          // 数据包头部长度（固定为 16）
	ProtocolVersion ProtocolVersion //
	Operation       OpreationType   //
	SequenceID      uint32          // 数据包头部长度（固定为 1）
	Body            []byte          //
}

func (pkt *Packet) Header() (h []byte) {
	h = make([]byte, 16)
	var au32 = binary.BigEndian.PutUint32
	var au16 = binary.BigEndian.PutUint16
	au32(h[0:4], pkt.PacketLength)
	au16(h[4:6], pkt.HeaderLength)
	au16(h[6:8], uint16(pkt.ProtocolVersion))
	au32(h[8:12], uint32(pkt.Operation))
	au32(h[12:16], pkt.SequenceID)
	return
}

func (pkt *Packet) Bytes() (h []byte) {
	h = pkt.Header()
	h = append(h, pkt.Body...)
	return h
}

func DecodeHeader(h []byte) (pkt *Packet) {
	var u32 = binary.BigEndian.Uint32
	var u16 = binary.BigEndian.Uint16
	pkt = &Packet{
		PacketLength:    u32(h[0:4]),
		HeaderLength:    u16(h[4:6]),
		ProtocolVersion: ProtocolVersion(u16(h[6:8])),
		Operation:       OpreationType(u32(h[8:12])),
		SequenceID:      u32(h[12:16]),
		Body:            h,
	}
	return pkt
}

func (pkt *Packet) Decode() (npkt *Packet, err error) {
	switch pkt.ProtocolVersion {
	case ProtocolNormalBuffer:
		pkt.Body = pkt.Body[pkt.HeaderLength:pkt.PacketLength]
		npkt = pkt
	case ProtocolInflateBuffer:
		r := bytes.NewBuffer(pkt.Body[pkt.HeaderLength:])
		zr := try.To1(zlib.NewReader(r))
		defer zr.Close()
		body := try.To1(io.ReadAll(zr))
		npkt = DecodeHeader(body)
		return
	case ProtocolBrotliBuffer:
		body := try.To1(cbrotli.Decode(pkt.Body[pkt.HeaderLength:]))
		npkt = DecodeHeader(body)
		return
	default:
		err = fmt.Errorf("unkonw protover %d", pkt.ProtocolVersion)
	}
	return
}

func DecodePackets(body []byte) (pkts []*Packet, err error) {
	defer err2.Handle(&err)
	var wg sync.WaitGroup
	for {
		next := func() (next bool) {
			pkt := DecodeHeader(body[:16])
			pkt.Body = body[:pkt.PacketLength]
			defer func() {
				if len(body) > int(pkt.PacketLength) {
					body = body[pkt.PacketLength:]
					next = true
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				switch pkt.ProtocolVersion {
				default:
					pkt = try.To1(pkt.Decode())
					pkts = append(pkts, pkt)
				case ProtocolBrotliBuffer:
					fallthrough
				case ProtocolInflateBuffer:
					pkt = try.To1(pkt.Decode())
					npkts := try.To1(DecodePackets(pkt.Body))
					pkts = append(pkts, npkts...)
				}
			}()
			return
		}()
		if !next {
			break
		}
	}
	wg.Wait()
	return
}

type ProtocolVersion uint16

const (
	ProtocolNormalBuffer   ProtocolVersion = 0 // 未压缩的buffer
	ProtocolInt32BigEndian ProtocolVersion = 1 // Body 内容为房间人气值
	ProtocolInflateBuffer  ProtocolVersion = 2 // 压缩过的 Buffer，Body 内容需要用zlib.inflate解压出一个新的数据包，然后从数据包格式那一步重新操作一遍
	ProtocolBrotliBuffer   ProtocolVersion = 3 // 压缩信息,需要brotli解压,然后从数据包格式 那一步重新操作一遍
)

type OpreationType uint32

const (
	OpreationHeartbeat          OpreationType = 2 // 客户端 - 心跳 - (空) - 不发送心跳包，70 秒之后会断开连接，通常每 30 秒发送 1 次
	OpreationHeartbeatReply     OpreationType = 3 // 服务器 - 心跳回应 - Int 32 Big Endian - Body 内容为房间人气值
	OpreationMessage            OpreationType = 5 // 服务器 - 通知 - JSON - 弹幕、广播等全部信息
	OpreationUserAuthentication OpreationType = 7 // 客户端	- 进房 - JSON - WebSocket 连接成功后的发送的第一个数据包，发送要进入房间 ID
	OpreationConnectSuccess     OpreationType = 8 // 服务器 - 进房回应 - (空)
)

// 进房-json-内容 https://github.com/lovelyyoshino/Bilibili-Live-API/blob/master/API.WebSocket.md#进房-json-内容
type EnterRoom struct {
	ClientVer *string `json:"clientver,omitempty"` // 例如 "1.5.10.1"
	Platform  *string `json:"platform,omitempty"`  // 例如 "web"
	ProtoVer  *int    `json:"protover,omitempty"`  // 1 或者 2. protover 为 1 时不会使用zlib压缩，为 2 时会发送带有zlib压缩的包，也就是数据包协议为 2
	RoomID    int     `json:"roomid"`              // 房间长 ID，可以通过 room_init API 获取
	UID       *int    `json:"uid,omitempty"`       // uin，可以通过 getUserInfo API 获取
	Type      *int    `json:"type,omitempty"`      // 不知道啥，总之写 2
}
