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

package server

import (
	"container/list"
	"fpay/base"
	"github.com/go-redis/redis"
	"io"
	"math/rand"
	"net"
	"time"
	"zlog"
)

type Handler interface {
	Handle(rw io.ReadWriter) (err error)
	Close()
}

type Worker struct {
	state chan uint8
	conn  *net.TCPConn
}

func newWorker(conn *net.TCPConn) (w *Worker) {
	w = new(Worker)
	w.state = make(chan uint8, 1)
	w.conn = conn
	return
}

type Parent struct {
	Worker
	addr       *net.TCPAddr
	lastState  uint8 //
	lastAlive  int64
	lastHeight uint64
}

const (
	PARENT_STATE_UNAVAILABLE = 1
	PARENT_STATE_REJECTED
	PARENT_STATE_CLOSED
	PARENT_STATE_RESERVED
)

var PARENT_STATES []string = []string{"CMD_SHUT", "PARENT_STATE_UNAVAILABLE", "PARENT_STATE_REJECTED", "PARENT_STATE_CLOSED", "PARENT_STATE_RESERVED", "PARENT_STATE_MAJOR"}

func newParent(addr *net.TCPAddr) (p *Parent) {
	p = new(Parent)
	p.state = make(chan uint8, 1)
	p.lastState = PARENT_STATE_UNAVAILABLE
	p.addr = addr
	return
}

type Server struct {
	acceptorState, finderState chan uint8
	workers                    map[*net.TCPConn]*Worker
	children                   map[*net.TCPConn]*Worker
	parent                     *Parent
	reservedParents            []*Parent
	unusedParents              *list.List
	prepParents                *list.List
	tcpAddr                    *net.TCPAddr
	tcpListener                *net.TCPListener
	handlers                   map[string]Handler
}

const (
	CMD_SHUT = iota
	STATE_READY
	STATE_DONE
	STATE_CLOSED
)

var SERVER_STATES []string = []string{"CMD_SHUT", "STATE_READY", "STATE_DONE", "STATE_CLOSED"}

func New(saddr string, officers []string) (s *Server, err error) {
	s = new(Server)
	s.acceptorState = make(chan uint8, 1)
	s.finderState = make(chan uint8, 1)
	s.workers = make(map[*net.TCPConn]*Worker, 100)
	s.children = make(map[*net.TCPConn]*Worker, 100)
	s.reservedParents = make([]*Parent, 0, 20)
	s.unusedParents = list.New()
	s.prepParents = list.New()
	s.handlers = make(map[string]Handler, 10)

	s.tcpAddr, err = net.ResolveTCPAddr("tcp", saddr)
	if err != nil {
		zlog.Fatalf("TCP address %s resolution failed.\n", saddr)
	}

	for _, saddr = range officers {
		addr, err := net.ResolveTCPAddr("tcp", saddr)
		if err != nil {
			continue
		}

		p := newParent(addr)

		s.prepParents.PushBack(p)
	}

	return
}

func (this *Server) checkAcceptable(addr *net.TCPAddr) (isAcceptable bool) {
	// TODO: 检查该客户是否满足建立连接的条件
	return true
}

func (this *Server) sendRecommendationList(conn *net.TCPConn) {
	// TODO: 发送推荐列表
}

func (this *Server) createChild(conn *net.TCPConn) {
	w := newWorker(conn)
	this.children[conn] = w
	go this.requestLoop(conn, w)
}

func (this *Server) releaseAll() {
	var saddr string

	zlog.Debugln("Children is going to be released.")
	for _, w := range this.children {
		w.state <- CMD_SHUT
		zlog.Traceln("CMD_SHUT is sended to %s" + w.conn.RemoteAddr().String())
	}

	zlog.Debugln("Parents is going to be released.")
	for _, w := range this.reservedParents {
		w.state <- CMD_SHUT
		zlog.Traceln("CMD_SHUT is sended to %s" + w.conn.RemoteAddr().String())
	}

	for _, w := range this.children {
		select {
		case s := <-w.state:
			if s == STATE_CLOSED {
				saddr = w.conn.RemoteAddr().String()

				zlog.Traceln("STATE_CLOSED is received from %s" + saddr)
				w.conn.Close()
				zlog.Tracef("Connection %s closed.\n", saddr)
			}
		default:
			<-time.After(10 * time.Millisecond)
		}
	}
	zlog.Debugln("All children released.")

	for _, w := range this.reservedParents {
		select {
		case s := <-w.state:
			if s == STATE_CLOSED {
				saddr = w.conn.RemoteAddr().String()

				zlog.Traceln("STATE_CLOSED is received from %s" + saddr)
				w.conn.Close()
				zlog.Tracef("Connection %s closed.\n", saddr)
			}
		default:
			<-time.After(10 * time.Millisecond)
		}
	}
	zlog.Debugln("All parents released.")
}

func (this *Server) requestLoop(conn *net.TCPConn, w *Worker) {
	saddr := conn.RemoteAddr().String()
	zlog.Debugf("ReaderLoop for %s is starting.\n", saddr)
	defer zlog.Debugf("Connection for %s closed.\n", saddr)
	defer conn.Close()

	var s uint8
	for {
		select {
		case s = <-w.state:
			zlog.Tracef("%s received.\n", SERVER_STATES[s])
			break
		default:
			c, err := base.Unmarshal(conn)
			if err != nil {
				if err != io.EOF {
					zlog.Warningln("Failed to unmarshal: " + err.Error())
				}
				break
			}
			protocol := string(c.Protocol)
			handler, ok := this.handlers[protocol]
			if !ok {
				zlog.Warningln("Unspport protocol: " + protocol)
			}

			err = handler.Handle(conn)
			if err != nil {
				break
			}
			// 暂时没可读数据，延长检查时间
			// 平均会造成2.5毫秒左右的处理延时
			// TODO: 该参数应该可以根据不同角色实现动态调整
			<-time.After(time.Duration(rand.Intn(10*1000*1000)) * time.Nanosecond)
		}
	}
	w.state <- STATE_CLOSED
}

func (this *Server) acceptorLoop() {
	zlog.Debugln("acceptorLoop is starting.\n")
	defer zlog.Debugln("acceptorLoop closed.\n")
	var saddr string
	for {
		select {
		case state := <-this.acceptorState:
			zlog.Tracef("%s received.\n", SERVER_STATES[state])
			switch state {
			case CMD_SHUT:
				break
			default:
				continue
			}
		default:
			conn, err := this.tcpListener.AcceptTCP()
			saddr = conn.RemoteAddr().String()
			if err != nil {
				zlog.Warningf("Connection from %s failed: %s.\n", saddr, err.Error())
				conn.Close()
				zlog.Debugf("Connection %s closed.\n", saddr)
				continue
			}

			addr, err := net.ResolveTCPAddr("tcp", saddr)
			if err != nil {
				panic("Impossible.")
			}
			ok := this.checkAcceptable(addr)
			if !ok {
				this.sendRecommendationList(conn)
				conn.Close()
				continue
			}

			this.createChild(conn)
		}

		<-time.After(10 * time.Millisecond)
	}

	this.releaseAll()
}

func (this *Server) preparedLoop(p *Parent) {
	// TODO: 同步后备父节点的数据，主要用于更新后备父节点的区块高度
}

func (this *Server) finderLoop() {
	// TODO: 查找后备父节点
	zlog.Infoln("FinderLoop is starting.")
	defer zlog.Infoln("FinderLoop is closed.")

	var err error
	for {
		select {
		case s := <-this.finderState:
			if s == CMD_SHUT {
				break
			}
		default:
			if this.prepParents.Len() == 0 {
				<-time.After(30 * time.Second)
				continue
			}

			p, ok := this.prepParents.Remove(this.prepParents.Front()).(*Parent)
			if !ok {
				panic("Impossible.")
			}

			p.conn, err = net.DialTCP("tcp", nil, p.addr)
			if err != nil {
				p.conn.Close()
				p.lastState = PARENT_STATE_UNAVAILABLE
				this.unusedParents.PushBack(p)
				continue
			}

			go this.receiverLoop(p)
		}

		<-time.After(500 * time.Millisecond)
	}
	zlog.Infoln("FinderLoop start to close.")

	for _, p := range this.reservedParents {
		p.state <- CMD_SHUT
		zlog.Traceln("Send CMD_SHUT to reserved parent: " + p.addr.String())
	}

	for {
		if len(this.reservedParents) == 0 {
			zlog.Debugln("All reserved parents closed.")
			break
		}
		zlog.Tracef("%v reserved parents left.\n", len(this.reservedParents))

		<-time.After(100 * time.Millisecond)
	}

	this.finderState <- STATE_CLOSED
}

func (this *Server) transferLoop() {

}

func (this *Server) receiverLoop() {
	// 接收父节点的数据
}

func (this *Server) broadcasterLoop() {

}

func (this *Server) loadHandlers() {

}

func (this *Server) Startup() (err error) {
	zlog.Infoln("Server is starting up.")

	this.loadHandlers()

	saddr := this.tcpAddr.String()
	this.tcpListener, err = net.ListenTCP("tcp", this.tcpAddr)
	if err != nil {
		zlog.Fatalf("Address %s binding failed.\n", saddr)
		this.tcpListener.Close()
		return
	}
	zlog.Debugf("Address %s is listening.\n", saddr)
	this.acceptorState <- STATE_READY
	zlog.Traceln("STATE_READY send to main.")

	go this.acceptorLoop()
	go this.finderLoop()

	return
}

func (this *Server) Shutdown() (err error) {
	zlog.Infoln("Server is Shutting down.")

	saddr := this.tcpAddr.String()

	err = this.tcpListener.Close()
	zlog.Debugf("Address %s already closed.\n", saddr)

	this.acceptorState <- CMD_SHUT
	zlog.Debugf("State %s is send to AcceptorLoop.\n", saddr)

	this.finderState <- CMD_SHUT
	zlog.Debugf("State %s is send to FinderLoop.\n", saddr)

	for i := 0; i < 2; {
		select {
		case s := <-this.acceptorState:
			if s == STATE_CLOSED {
				i++
				zlog.Debugln("Receive STATE_CLOSED from acceptorLoop.")
				continue
			}
		case s := <-this.finderState:
			if s == STATE_CLOSED {
				i++
				zlog.Debugln("Receive STATE_CLOSED from finderLoop.")
				continue
			}
		default:
			continue
		}
	}
	zlog.Infoln("Server already closed.")
	return
}
