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

package fpay

import (
	"net"
)

type Child struct {
	addr  *net.TCPAddr
	conn  *net.TCPConn
	brcst *BroadCaster
	rt    *Router
}

func NewChild(conn *net.TCPConn) (chd *Child) {
	chd = new(Child)
	chd.addr, _ = net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	chd.conn = conn
	chd.brcst = NewBroadCaster(conn)
	chd.rt = NewRouter(conn)
	return
}

func (this *Child) Startup() {
	this.brcst.Startup()
	this.rt.Startup()
}

func (this *Child) Shutdown() {
	this.brcst.Shutdown()
	this.rt.Shutdown()
}

/*
func (this *Server) releaseAll() {
	var saddr string

	zlog.Debugln("Children is going to be released.")
	for _, w := range this.children {
		w.state <- CMD_SHUT
		zlog.Traceln("CMD_SHUT is sended to %s" + w.conn.RemoteAddr().String())
	}

	zlog.Debugln("Parents is going to be released.")
	for _, w := range this.reservedParents {
		w.state <- CMD_SHUT
		zlog.Traceln("CMD_SHUT is sended to %s" + w.conn.RemoteAddr().String())
	}

	for _, w := range this.children {
		select {
		case s := <-w.state:
			if s == STATE_CLOSED {
				saddr = w.conn.RemoteAddr().String()

				zlog.Traceln("STATE_CLOSED is received from %s" + saddr)
				w.conn.Close()
				zlog.Tracef("Connection %s closed.\n", saddr)
			}
		default:
			<-time.After(10 * time.Millisecond)
		}
	}
	zlog.Debugln("All children released.")

	for _, w := range this.reservedParents {
		select {
		case s := <-w.state:
			if s == STATE_CLOSED {
				saddr = w.conn.RemoteAddr().String()

				zlog.Traceln("STATE_CLOSED is received from %s" + saddr)
				w.conn.Close()
				zlog.Tracef("Connection %s closed.\n", saddr)
			}
		default:
			<-time.After(10 * time.Millisecond)
		}
	}
	zlog.Debugln("All parents released.")
}
*/
