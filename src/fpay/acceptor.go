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
	"io"
	"net"
	"sync"
	"zlog"
)

type Acceptor struct {
	Core
	addr     *net.TCPAddr
	saddr    string
	lsn      *net.TCPListener
	children map[string]*Child
	locker   *sync.Mutex
}

func NewAcceptor(laddr *net.TCPAddr) (act *Acceptor) {
	act = new(Acceptor)
	act.Init(act)
	act.addr = laddr
	act.saddr = laddr.String()
	act.children = make(map[string]*Child, 100)
	act.locker = new(sync.Mutex)
	return
}

func (this *Acceptor) checkAvailable(addr *net.TCPAddr) (isAvailable bool) {
	return true
}

func (this *Acceptor) sendRecommandList(conn *net.TCPConn) {}

func (this *Acceptor) PreLoop() (err error) {
	zlog.Infoln("Binding address: " + this.saddr)

	this.locker.Lock()
	this.lsn, err = net.ListenTCP("tcp", this.addr)
	if err != nil {
		zlog.Fatalf("Address %s binding failed.\n", this.saddr)
		if this.lsn != nil {
			this.lsn.Close()
			this.locker.Unlock()
		}
		return
	}
	this.locker.Unlock()
	zlog.Debugf("Address %s is listening.\n", this.saddr)
	zlog.Infoln("Address binded.")
	return
}

func (this *Acceptor) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		switch cmd {
		case CMD_SHUT:
			zlog.Traceln("CMD_SHUT received.")
			return false
		default:
			zlog.Warningf("Unsupport %s received.\n", CMDS[cmd])
			return true
		}
	default:
		zlog.Traceln("Waiting for next connection.")
		conn, err := this.lsn.AcceptTCP()
		if err != nil {
			if err == io.EOF {
				zlog.Warningf("Connection %s closed.\n", this.saddr)
			}
			return false
		}

		raddr, _ := net.ResolveTCPAddr("tcp", this.saddr)

		// TODO:
		// 检查请求者是否具备连接资格
		ok := this.checkAvailable(raddr)
		if !ok {
			// 如果不具备，则返回推荐列表，并关闭连接
			this.sendRecommandList(conn)
			conn.Close()
		} else {
			// 如果具备，创建API线程接收请求以及创建广播线程推送数据
			chd := NewChild(conn)
			chd.Startup()
			this.children[raddr.String()] = chd
		}
		return true
	}
}

func (this *Acceptor) AftLoop() {
	zlog.Infoln("Closing children.")
	for _, child := range this.children {
		child.Shutdown()
	}
	zlog.Infoln("All children closed.")
}

func (this *Acceptor) Shutdown() {
	zlog.Infoln("Shutting down.")
	this.locker.Lock()
	if this.lsn != nil {
		this.lsn.Close()
	}
	this.locker.Unlock()

	this.Core.Shutdown()
	zlog.Infoln("Closed.")
}
