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
	Ctx        *FPAY
	Officers   []*net.TCPAddr
	RsvParents map[string]*Parent
	PreAddrs   *list.List
	UnusAddrs  *list.List
}

func FinderNew(ctx *FPAY) (fd *Finder) {
	fd = new(Finder)
	fd.Init(fd)
	fd.Ctx = ctx
	fd.Officers = make([]*net.TCPAddr, 0, len(ctx.Officers))
	fd.RsvParents = make(map[string]*Parent)
	fd.PreAddrs = list.New()
	fd.UnusAddrs = list.New()

	// 如果存在默认父节点，则优先连接
	oaddr, err := net.ResolveTCPAddr("tcp", ctx.Settings.Paddr)
	if err == nil {
		fd.PreAddrs.PushBack(oaddr)
	}

	// 先请求官方节点做服务发现
	for _, officer := range ctx.Officers {
		oaddr, err = net.ResolveTCPAddr("tcp", officer)
		if err == nil {
			fd.Officers = append(fd.Officers, oaddr)
		}

		// 确认本节点是不是官方节点。目前只通过监听的IP地址确认，后期需要找到更准确的办法
		if oaddr.String() != ctx.Settings.Laddr {
			// 如果本节点与官方节点不一致一致，则建立连接
			// 主要目的为了避免自己连接自己
			fd.PreAddrs.PushBack(oaddr)
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
		if this.PreAddrs.Len() == 0 {
			zlog.Traceln("Looping.")
			<-time.After(500 * time.Millisecond)
			return true
		}

		raddr, ok := this.PreAddrs.Remove(this.PreAddrs.Front()).(*net.TCPAddr)
		if !ok {
			panic("Impossible.")
		}

		p := ParentNew(this.Ctx, raddr)
		err := p.Startup()
		if err != nil {
			this.UnusAddrs.PushBack(raddr)
		} else {
			this.RsvParents[raddr.String()] = p
		}

	}
	return true
}

func (this *Finder) AftLoop() {
	zlog.Infoln("Shutting down.")

	for _, p := range this.RsvParents {
		p.Shutdown()
	}

	zlog.Infoln("Closed.")
}
