package bilive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type MsgDanmu struct {
	Info []any `json:"info"`
}

func TestConnect(t *testing.T) {
	err2.SetErrorTracer(os.Stderr)
	client := NewClient(898286)
	defer client.Close()
	try.To(client.Connect())
	for {
		_, b := try.To2(client.Conn.Read(context.Background()))
		go func(b []byte) {
			defer err2.CatchTrace(func(err error) {
				log.Println("err", err)
			})
			pkts := try.To1(DecodePackets(b))
			for _, pkt := range pkts {
				go func(pkt *Packet) {
					defer err2.Catch(func(err error) {
						log.Printf("err %+v\n", pkt)
					})
					if pkt.Operation != OpreationMessage {
						return
					}
					var msg Message
					try.To(json.Unmarshal(pkt.Body, &msg))
					if msg.CMD != CMD_DANMU_MSG {
						return
					}
					var danmu MsgDanmu
					try.To(json.Unmarshal(pkt.Body, &danmu))
					// fmt.Println(string(pkt.Body))
					fmt.Fprintln(os.Stderr, "danmu", danmu.Info[2].([]any)[1], danmu.Info[1])
				}(pkt)
			}
		}(b)
	}
}
