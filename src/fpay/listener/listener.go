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

package listener

import (
	"net"
	"zlog"
)

type Listener struct {
	in, out     chan uint8
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
}

func New(addr string) (listener *Listener, err error) {
	listener = new(Listener)
	listener.in = make(chan uint8, 1)
	listener.out = make(chan uint8, 1)

	listener.tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		zlog.Fatalf("TCP address %s resolution failed.\n", addr)
	}
	return
}

func (this *Listener) loop() {
	for {
		select {
		case <-this.in:
			this.out <- 0
			break
		default:
			conn, err := this.tcpListener.AcceptTCP()
			if err != nil {
				break
			}

			b := make([]byte, 2048)
			_, err = conn.Read(b)
		}
	}
}

func (this *Listener) Startup() (err error) {
	zlog.Infoln("Listener service is starting up.")

	this.tcpListener, err = net.ListenTCP("tcp", this.tcpAddr)
	if err == nil {
		zlog.Debugf("Address %s is listening.\n", this.tcpAddr.String())
	} else {
		zlog.Fatalf("Address %s binding failed.\n", this.tcpAddr.String())
		this.tcpListener.Close()
		return
	}

	go this.loop()

	return
}

func (this *Listener) Shutdown() (err error) {
	zlog.Infoln("Listener service is Shutting down.")

	err = this.tcpListener.Close()
	zlog.Debugf("Address %s already closed.\n", this.tcpAddr.String())

	this.in <- 0
	<-this.out
	zlog.Infoln("Listener service already closed.")
	return
}
