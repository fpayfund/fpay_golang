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

type Child struct {
	Core
	Ctx   *FPAY
	Addr  *net.TCPAddr
	Conn  *net.TCPConn
	Brcst *BroadCaster
}

func ChildNew(ctx *FPAY, conn *net.TCPConn) (chd *Child) {
	chd = new(Child)
	chd.Init(chd)
	chd.Ctx = ctx
	chd.Addr, _ = net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	chd.Conn = conn
	chd.Brcst = BroadCasterNew(ctx, conn)
	return
}

func (this *Child) PreLoop() (err error) {
	return
}

func (this *Child) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		zlog.Tracef("%s received.\n", CMDS[cmd])
		if cmd == CMD_SHUT {
			return false
		}
	default:
		<-time.After(500 * time.Millisecond)
	}
	return
}

func (this *Child) AftLoop() {
}
