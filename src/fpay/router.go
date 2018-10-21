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
	"math/rand"
	"net"
	"time"
	"zlog"
)

type Router struct {
	Core
	conn  *net.TCPConn
	saddr string
}

func NewRouter(conn *net.TCPConn) (rt *Router) {
	rt = new(Router)
	rt.conn = conn
	return
}

// 需要重写
func (this *Router) PreLoop() (err error) {
	this.saddr = this.conn.RemoteAddr().String()
	return
}

// 需要重写
func (this *Router) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		zlog.Tracef("%s received.\n", CMDS[cmd])
		return false
	default:
		_, err := UnmarshalBase(this.conn)
		if err != nil {
			if err != io.EOF {
				zlog.Warningln("Failed to unmarshal: " + err.Error())
			}
			return false
		}
		/*
			protocol := string(c.Protocol)
			handler, ok := this.handlers[protocol]
			if !ok {
				zlog.Warningln("Unspport protocol: " + protocol)
			}

			err = handler.Handle(conn)
			if err != nil {
				break
			}
		*/
		// 暂时没可读数据，延长检查时间
		// 平均会造成2.5毫秒左右的处理延时
		// TODO: 该参数应该可以根据不同角色实现动态调整
		<-time.After(time.Duration(rand.Intn(10*1000*1000)) * time.Nanosecond)
	}
	return true
}

// 需要重写
func (this *Router) AftLoop() {}
