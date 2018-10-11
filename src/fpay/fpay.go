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
	"fpay/cli"
	"fpay/datasource"
	"fpay/monitor"
	"net"
	"time"
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

type FPAY struct {
	// FPAY服务状态:
	//   STARTING:
	//     启动状态。负责启动系统监控服务、数据服务、监听服务和节点发现服务，并根据节点发现的结果切换到不同的状态

	//   SHUTTING:
	//     结束状态。负责通知并等待各个服务关闭

	//   BOOKKEEPER:
	//     出块者状态。负责汇集请求并出块

	//   REVIEWER:
	//     评审者状态。负责汇集请求，提交给出块者，同时评审出块者的工作成果并广播

	//   TRANSFERER:
	//     传送者状态。负责汇集请求，提交给上级传送者，同时广播

	//   TOP_TRANSFERER:
	//     顶级传送者状态。负责汇集请求，提交给评审者，同时广播。顶级传送者可以审核整个出块者和评审者的工作业绩

	//   RECEIVER:
	//     传送者状态。负责汇集请求，提交给上级传送者，同时广播

	//   PAYER:
	//     支付者状态。负责提交请求。只接收与自己相关的信息
	state                                   uint8
	in, out                                 chan uint8 // 命令输入、输出队列
	nodes, availableNodes, unavailableNodes []string   // 节点，可用节点，不可用节点
	settings                                *cli.Settings
	datasource                              *datasource.DataSource
	monitor                                 *monitor.Monitor
	tcpAddr                                 *net.TCPAddr
	tcpListener                             *net.TCPListener
	tcpConnections                          []*net.TCPConn
}

// FPAY官方启动节点
var baseNodes = []string{
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

func New(settings *cli.Settings) (fs *FPAY) {
	fs = new(FPAY)
	fs.settings = settings
	fs.in = make(chan uint8, 1)
	fs.out = make(chan uint8, 1)
	fs.datasource = datasource.New(settings)
	fs.monitor = monitor.New(settings)
	fs.nodes = make([]string, 100)
	fs.availableNodes = make([]string, 50)
	fs.unavailableNodes = make([]string, 10)
	return
}

func (this *FPAY) starting() {
	zlog.Traceln("STARTING")

}

func (this *FPAY) bookkeeper() {
	zlog.Traceln("BOOKKEEPER")
}

func (this *FPAY) reviewer() {
	zlog.Traceln("REVIEWER")
}

func (this *FPAY) transferer() {
	zlog.Traceln("TRANSFERER")
}

func (this *FPAY) topTransferer() {
	zlog.Traceln("TOP_TRANSFERER")
}

func (this *FPAY) receiver() {
	zlog.Traceln("RECEIVER")
}

func (this *FPAY) payer() {
	zlog.Traceln("PAYER")
}

func (this *FPAY) shutting() {
	zlog.Traceln("SHUTTING")
	this.out <- 0
}

func (this *FPAY) loop() {
	for {
		select {
		case <-this.in:
			this.state = SHUTTING
		case <-time.After(200 * time.Millisecond):
			switch this.state {
			case STARTING:
				this.starting()
			case BOOKKEEPER:
				this.bookkeeper()
			case REVIEWER:
				this.reviewer()
			case TRANSFERER:
				this.transferer()
			case TOP_TRANSFERER:
				this.topTransferer()
			case RECEIVER:
				this.receiver()
			case PAYER:
				this.payer()
			case SHUTTING:
				this.shutting()
			default:
				zlog.Errorln("Unsupport state: %u\n", this.state)
				this.state = SHUTTING
			}
		}
	}
}

func (this *FPAY) Startup() (err error) {
	zlog.Infoln("FPAY service is starting up.")

	this.datasource.Startup()
	defer func() {
		if err != nil {
			this.datasource.Shutdown()
		}
	}()

	this.monitor.Startup()
	defer func() {
		if err != nil {
			this.monitor.Shutdown()
		}
	}()

	this.tcpAddr, err = net.ResolveTCPAddr("tcp", this.settings.TCPAddr)
	if err != nil {
		zlog.Fatalf("TCP address %s resolution failed.\n", this.settings.TCPAddr)
		return
	}

	this.tcpListener, err = net.ListenTCP("tcp", this.tcpAddr)
	if err != nil {
		zlog.Fatalf("TCP address %s listening failed.\n", this.settings.TCPAddr)
		return
	}
	defer func() {
		if err != nil {
			this.tcpListener.Close()
		}
	}()

	zlog.Debugf("TCP listener at %s started.\n", this.settings.TCPAddr)

	go this.loop()
	return
}

func (this *FPAY) Shutdown() {
	zlog.Infoln("FPAY service is Shutting down.")

	if this.tcpListener != nil {
		this.tcpListener.Close()
		zlog.Debugf("TCP listener at %s closed.\n", this.settings.TCPAddr)
	}

	defer zlog.Traceln("FPAY service already closed.")
	defer this.datasource.Shutdown()
	defer this.monitor.Shutdown()

	this.in <- 0
}
