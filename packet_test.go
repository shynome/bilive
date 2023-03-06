package ws

import (
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

func TestJSONPacket(t *testing.T) {
	pkt := NewJSONPacket([]byte("hello world"))
	pkt2 := try.To1(Decode(pkt.Bytes()))
	assert.DeepEqual(pkt2, pkt)
}
