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
	"zlog"
)

type Parent struct {
	saddr string
	raddr *net.TCPAddr
	conn  *net.TCPConn
	rcv   *Receiver
	rv    *Reviewer
	trsf  *Transferer
}

func NewParent(raddr *net.TCPAddr) (prt *Parent) {
	prt = new(Parent)
	prt.raddr = raddr
	prt.saddr = raddr.String()
	prt.rcv = NewReceiver(prt.raddr)
	return
}

func (this *Parent) Startup() (err error) {
	zlog.Infof("Connection to %s.\n", this.saddr)

	err = this.rcv.Startup()
	if err != nil {
		zlog.Warningln("Startup failed: " + err.Error())
	}
	return
}

func (this *Parent) Shutdown() {
	this.rcv.Shutdown()
}
