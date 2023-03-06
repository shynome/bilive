package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"nhooyr.io/websocket"
)

type Client struct {
	RoomID int
	Conn   *websocket.Conn

	closeKeepAlive context.CancelFunc
}

func NewClient(roomid int) *Client {
	return &Client{
		RoomID: roomid,
	}
}

func (c *Client) Connect() (err error) {
	defer err2.Handle(&err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, _ := try.To2(websocket.Dial(ctx, "wss://broadcastlv.chat.bilibili.com/sub", nil))
	c.Conn = conn
	go func() {
		pkt := NewConnectPacket(c.RoomID)
		c.WritePacket(pkt)
	}()
	pkt := try.To1(c.ReadPacket())
	if pkt.Operation != OpreationConnectSuccess {
		return fmt.Errorf("first packet op must be %d, but got %d", OpreationConnectSuccess, pkt.Operation)
	}
	go c.keepAlive()
	return
}

// Decode 是解码操作会耗时, 追求性能的话移动到协程里
func (c *Client) ReadPacket() (pkt Packet, err error) {
	defer err2.Handle(&err)
	ctx := context.Background()
	_, e := try.To2(c.Conn.Read(ctx))
	pkt = try.To1(Decode(e))
	return
}

func (c *Client) WritePacket(pkt Packet) (err error) {
	ctx := context.Background()
	err = c.Conn.Write(ctx, websocket.MessageBinary, pkt.Bytes())
	return
}

func (c *Client) keepAlive() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.closeKeepAlive = cancel
	t := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ping := NewPingPacket()
			c.WritePacket(ping)
		}
	}
}

func (c *Client) Close() (err error) {
	if c.closeKeepAlive != nil {
		c.closeKeepAlive()
	}
	return
}
