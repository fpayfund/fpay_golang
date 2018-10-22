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
	"container/list"
	"net"
	"time"
	"zlog"
)

type Finder struct {
	Core
	laddr      *net.TCPAddr
	officers   []*net.TCPAddr
	rsvParents map[string]*Parent
	preAddrs   *list.List
	unusAddrs  *list.List
}

func NewFinder(laddr *net.TCPAddr, officers []string) (fd *Finder) {
	fd = new(Finder)
	fd.Init(fd)
	fd.laddr = laddr
	fd.officers = make([]*net.TCPAddr, 0, len(officers))
	fd.rsvParents = make(map[string]*Parent)
	fd.preAddrs = list.New()
	fd.unusAddrs = list.New()

	for _, officer := range officers {
		oaddr, err := net.ResolveTCPAddr("tcp", officer)
		if err == nil {
			fd.officers = append(fd.officers, oaddr)
		}

		if oaddr.String() != laddr.String() {
			fd.preAddrs.PushBack(oaddr)
		}
	}
	return
}

func (this *Finder) PreLoop() (err error) {
	zlog.Infoln("Starting up.")
	return nil
}

func (this *Finder) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		zlog.Tracef("%s received.\n", CMDS[cmd])
		if cmd == CMD_SHUT {
			return false
		}
	default:
		if this.preAddrs.Len() == 0 {
			zlog.Traceln("Looping.")
			<-time.After(500 * time.Millisecond)
			return true
		}

		raddr, ok := this.preAddrs.Remove(this.preAddrs.Front()).(*net.TCPAddr)
		if !ok {
			panic("Impossible.")
		}

		p := NewParent(raddr)
		err := p.Startup()
		if err != nil {
			this.unusAddrs.PushBack(raddr)
		} else {
			this.rsvParents[raddr.String()] = p
		}

	}
	return true
}

func (this *Finder) AftLoop() {
	zlog.Infoln("Shutting down.")

	for _, p := range this.rsvParents {
		p.Shutdown()
	}

	zlog.Infoln("Closed.")
}
