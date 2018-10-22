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

type Receiver struct {
	Core
	raddr *net.TCPAddr
	saddr string
	conn  *net.TCPConn
}

func NewReceiver(raddr *net.TCPAddr) (rcv *Receiver) {
	rcv = new(Receiver)
	rcv.Init(rcv)
	rcv.raddr = raddr
	rcv.saddr = raddr.String()
	return
}

func (this *Receiver) PreLoop() (err error) {
	zlog.Infof("Connecting to %s.\n", this.saddr)

	this.conn, err = net.DialTCP("tcp", nil, this.raddr)
	if err != nil {
		zlog.Warningf("Connect %s failed.", this.saddr)
	}
	zlog.Infoln("Connected.")
	return
}

// 需要重写
func (this *Receiver) Loop() (isContinue bool) {
	<-time.After(1000 * time.Millisecond)
	return true
}

// 需要重写
func (this *Receiver) AftLoop() {
	this.conn.Close()
}
