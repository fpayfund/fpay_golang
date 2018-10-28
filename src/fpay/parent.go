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
	"time"
	"zlog"
)

type Parent struct {
	Core
	Ctx     *FPAY
	Saddr   string
	Raddr   *net.TCPAddr
	Conn    *net.TCPConn
	Account *Account
}

func ParentNew(ctx *FPAY, raddr *net.TCPAddr) (prt *Parent) {
	prt = new(Parent)
	prt.Init(prt)
	prt.Ctx = ctx
	prt.Raddr = raddr
	prt.Saddr = raddr.String()
	return
}

func (this *Parent) PreLoop() (err error) {
	zlog.Infof("Try to connect to %s.\n", this.Saddr)

	this.Conn, err = net.DialTCP("tcp", nil, this.Raddr)
	if err != nil {
		zlog.Warningf("%s connect failed.\n", this.Saddr)
	} else {
		zlog.Infof("%s connection established.")
	}
	return
}

func (this *Parent) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		zlog.Tracef("%s received.\n", CMDS[cmd])
		if cmd == CMD_SHUT {
			return false
		}
	default:
		<-time.After(500 * time.Millisecond)
	}
	return true
}

func (this *Parent) AftLoop() {
	zlog.Infof("Closing connection with %s.\n", this.Saddr)

	this.Conn.Close()

	zlog.Infof("connection with %s closed.\n", this.Saddr)
}
