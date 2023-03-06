package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

func TestConnect(t *testing.T) {
	client := NewClient(898286)
	defer client.Close()
	try.To(client.Connect())
	for {
		_, b := try.To2(client.Conn.Read(context.Background()))
		go func(b []byte) {
			defer err2.Catch(func(err error) {
				log.Println("err", err)
			})
			pkt := try.To1(Decode(b))
			if pkt.Operation != OpreationMessage {
				return
			}
			var msg Message
			try.To(json.Unmarshal(pkt.Body, &msg))
			if msg.CMD != CMD_DANMU_MSG {
				return
			}
			fmt.Println("danmu", string(pkt.Body))
		}(b)
	}
}
