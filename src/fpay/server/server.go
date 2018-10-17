/* The MIT License (MIT)
Copyright © 2018 by Atlas Lee(atlas@fpay.io)

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the “Software”),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"fpay/base"
	"io"
	"math/rand"
	"net"
	"time"
	"zlog"
)

type Handler interface {
	Handle(rw io.ReadWriter) (err error)
	Close()
}

type Connection struct {
	state chan uint8
	conn  *net.Conn
	buf   []byte
}

type Server struct {
	state       chan uint8
	conns       map[*net.TCPConn]chan uint8
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
	handlers    map[string]Handler
}

const (
	CMD_SHUT = iota
	STATE_READY
	STATE_DONE
	STATE_CLOSED
)

var STATE_NAMES []string = []string{"CMD_SHUT", "STATE_READY", "STATE_DONE", "STATE_CLOSED"}

func New(addr string) (s *Server, err error) {
	s = new(Server)
	s.state = make(chan uint8, 1)
	s.conns = make(map[*net.TCPConn]chan uint8, 100)
	s.handlers = make(map[string]Handler, 100)

	s.tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		zlog.Fatalf("TCP address %s resolution failed.\n", addr)
	}
	return
}

func (this *Server) ReaderLoop(conn *net.TCPConn, state chan uint8) {
	saddr := conn.RemoteAddr().String()
	zlog.Debugf("ReaderLoop for %s is starting.\n", saddr)
	defer zlog.Debugf("Connection for %s closed.\n", saddr)
	defer conn.Close()

	var s uint8
	for {
		select {
		case s = <-state:
			zlog.Tracef("%s received.\n", STATE_NAMES[s])
			break
		default:
			c, err := base.Unmarshal(conn)
			if err != nil {
				if err != io.EOF {
					zlog.Warningln("Failed to unmarshal: " + err.Error())
				}
				break
			}
			zlog.Traceln(c)
			/*
				protocol := string(c.Protocol)
				handler, ok := this.handlers[protocol]
				if !ok {
					zlog.Warningln("Unspport protocol: " + protocol)
				}

				err = handler.Handle(conn)
				if err != nil {
					break
				}
			*/
			// 暂时没可读数据，延长检查时间
			// 平均会造成2.5毫秒左右的处理延时
			// TODO: 该参数应该可以根据不同角色实现动态调整
			<-time.After(time.Duration(rand.Intn(10*1000*1000)) * time.Nanosecond)
		}
	}
	state <- STATE_CLOSED
}

func (this *Server) AcceptorLoop() {
	zlog.Debugln("AcceptorLoop is starting.\n")
	defer zlog.Debugln("AcceptorLoop closed.\n")
	var saddr string
	for {
		select {
		case state := <-this.state:
			zlog.Tracef("%s received.\n", STATE_NAMES[state])
			switch state {
			case CMD_SHUT:
				break
			default:
				continue
			}
		default:
			conn, err := this.tcpListener.AcceptTCP()
			if conn == nil {
				break
			}

			saddr = conn.RemoteAddr().String()

			if err != nil {
				zlog.Warningf("%s connect failed: %s.\n", saddr, err.Error())
				conn.Close()
				zlog.Debugf("%s connection closed.\n", saddr)
				continue
			}

			s := make(chan uint8, 1)
			this.conns[conn] = s

			go this.ReaderLoop(conn, s)
		}
		<-time.After(10 * time.Millisecond)
	}
}

func (this *Server) FinderLoop() {

}

func (this *Server) TransferLoop() {

}

func (this *Server) ReceiverLoop() {

}

func (this *Server) BroadcasterLoop() {

}

func (this *Server) Startup() (err error) {
	zlog.Infoln("Server is starting up.")

	saddr := this.tcpAddr.String()

	this.tcpListener, err = net.ListenTCP("tcp", this.tcpAddr)
	if err != nil {
		zlog.Fatalf("Address %s binding failed.\n", saddr)
		this.tcpListener.Close()
		return
	}
	zlog.Debugf("Address %s is listening.\n", saddr)
	this.state <- STATE_READY
	zlog.Traceln("STATE_READY send to main.")

	go this.AcceptorLoop()

	return
}

func (this *Server) Shutdown() (err error) {
	zlog.Infoln("Server is Shutting down.")

	saddr := this.tcpAddr.String()

	err = this.tcpListener.Close()
	zlog.Debugf("Address %s already closed.\n", saddr)

	this.state <- CMD_SHUT
	zlog.Debugf("%s is send to server.\n", saddr)

	<-this.state
	zlog.Infoln("Server already closed.")
	return
}
