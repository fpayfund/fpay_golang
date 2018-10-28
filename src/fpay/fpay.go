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
	"net"
	"sync"
	"zlog"
)

const (
	STARTING uint8 = iota
	BOOKKEEPER
	REVIEWER
	TRANSFERER
	TOP_TRANSFERER
	RECEIVER
	PAYER
	SHUTTING
)

var VERSION []byte = []byte{0, 0, 1, 0}

type FPAY struct {
	Core
	Version  []byte            // 4位版本号
	Settings *Settings         // 初始设置，启动参数
	Officers []string          // 官方节点
	Laddr    *net.TCPAddr      // 监听端口
	Lsn      *net.TCPListener  // 监听队列
	DB       *Cache            // 缓存
	Fd       *Finder           // 节点发现go程
	Parent   *Parent           // 父节点
	Children map[string]*Child // 子节点go程
	State    uint8             // 当前状态
	Locker   *sync.Mutex       // 状态锁
}

// FPAY官方启动节点
var Officers = []string{
	"127.0.0.1:8080",
	"127.0.0.1:8081",
	"127.0.0.1:8082",
	"127.0.0.1:8083",
	"127.0.0.1:8084",
	"127.0.0.1:8085",
	"127.0.0.1:8086",
	"127.0.0.1:8087",
	"127.0.0.1:8088",
	"127.0.0.1:8089"}

func FPAYNew(settings *Settings) (fs *FPAY, err error) {
	fs = new(FPAY)
	fs.Init(fs)
	fs.Version = VERSION
	fs.Settings = settings
	fs.Laddr, err = net.ResolveTCPAddr("tcp", settings.Laddr)
	if err != nil {
		zlog.Errorln("Address %s resolved failed: " + err.Error())
		return nil, err
	}

	fs.Officers = Officers
	fs.Children = make(map[string]*Child, 100)
	fs.Locker = new(sync.Mutex)
	fs.Fd = FinderNew(fs)
	return
}

func (this *FPAY) checkAvailable(addr *net.TCPAddr) (isAvailable bool) {
	return true
}

func (this *FPAY) sendContext(conn *net.TCPConn) {}

func (this *FPAY) sendRecommandList(conn *net.TCPConn) {}

func (this *FPAY) PreLoop() (err error) {
	zlog.Infoln("Connecting Database.")

	err = this.DB.Startup()

	zlog.Infoln("Database connected.")

	zlog.Infoln("Binding address: " + this.Settings.Laddr)

	this.Locker.Lock()
	this.Lsn, err = net.ListenTCP("tcp", this.Laddr)
	if err != nil {
		zlog.Fatalf("Address %s binding failed.\n", this.Settings.Laddr)
		if this.Lsn != nil {
			this.Lsn.Close()
			this.Locker.Unlock()
		}
		return
	}
	this.Locker.Unlock()
	zlog.Debugf("Address %s is listening.\n", this.Settings.Laddr)
	zlog.Infoln("Address binded.")

	err = this.Fd.Startup()

	zlog.Infoln("Started.")
	return
}

func (this *FPAY) Loop() (isContinue bool) {
	select {
	case cmd := <-this.Command:
		switch cmd {
		case CMD_SHUT:
			zlog.Traceln("CMD_SHUT received.")
			return false
		default:
			zlog.Warningf("Unsupport %s received.\n", CMDS[cmd])
			return true
		}
	default:
		zlog.Traceln("Waiting for next connection.")
		conn, err := this.Lsn.AcceptTCP()
		if err != nil {
			if err == io.EOF {
				zlog.Warningf("Connection %s closed.\n", this.Settings.Laddr)
			}
			return false
		}

		raddr, _ := net.ResolveTCPAddr("tcp", this.Settings.Laddr)

		// TODO:
		// 检查请求者是否具备连接资格
		ok := this.checkAvailable(raddr)
		if !ok {
			// 如果不具备，则返回推荐列表，并关闭连接
			this.sendRecommandList(conn)
			conn.Close()
		} else {
			// 如果具备，创建API线程接收请求以及创建广播线程推送数据
			this.sendContext(conn)
			chd := ChildNew(this, conn)
			chd.Startup()
			this.Children[raddr.String()] = chd
		}
		return true
	}
}

func (this *FPAY) AftLoop() {
	this.Fd.Shutdown()

	zlog.Infoln("Closing child connections.")

	for _, child := range this.Children {
		child.Shutdown()
	}

	zlog.Infoln("All child connections closed.")
}

func (this *FPAY) Shutdown() {
	zlog.Infoln("Shutting down.")

	this.Locker.Lock()
	if this.Lsn != nil {
		this.Lsn.Close()
	}
	this.Locker.Unlock()

	this.Core.Shutdown()

	zlog.Infoln("Closed.")
}
