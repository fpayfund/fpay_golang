package server

import (
	"fpay/base"
	"fpay/server"
	"net"
	"testing"
	"time"
	"zlog"
)

func TestNew(t *testing.T) {

	s, _ := server.New("0.0.0.0:8080")
	s.Startup()

	<-time.After(1 * time.Second)

	raddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	conn, _ := net.DialTCP("tcp", nil, raddr)
	b := base.New("Hello")

	for i := 0; i < 256; i++ {
		b.Marshal(conn)
		b.Version[3] = byte(i)
		zlog.Traceln(b)
		zlog.Debugf("b[%v] marshaled.\n", i)
		<-time.After(10 * time.Millisecond)
	}
}
