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

type Server struct {
	laddr *net.TCPAddr
	acpt  *Acceptor
	fd    *Finder
}

func (this *Server) loadHandlers() {}

func NewServer(saddr string, officers []string) (s *Server) {
	var err error
	s = new(Server)
	s.laddr, err = net.ResolveTCPAddr("tcp", saddr)
	if err != nil {
		zlog.Errorln("Address %s resolved failed: " + err.Error())
		return nil
	}

	s.acpt = NewAcceptor(s.laddr)
	s.fd = NewFinder(s.laddr, officers)
	return
}

func (this *Server) Startup() (err error) {
	err = this.acpt.Startup()

	if err == nil {
		err = this.fd.Startup()
	}
	return
}

func (this *Server) Shutdown() {
	this.acpt.Shutdown()
	this.fd.Shutdown()
}
