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
	"time"
	"zlog"
)

const (
	STARTING uint8 = iota
	BOOKKEEPER
	REVIEWER
	TRANSFERER
	RECEIVER
	PAYER
	SHUTTING
)

type FPAY struct {
	state      uint8      // 状态，BOOKKEEPER, REVIEWER, TRANSFORMER
	in, out    chan uint8 // 命令输入、输出队列
	settings   *cli.Settings
	datasource *datasource.DataSource
	monitor    *monitor.Monitor
}

func New(settings *cli.Settings) (ctx *FPAY) {
	ctx = new(FPAY)
	ctx.in = make(chan uint8, 1)
	ctx.out = make(chan uint8, 1)
	ctx.datasource = datasource.New(settings)
	ctx.monitor = monitor.New(settings)
	return
}

func (this *FPAY) starting() {
	zlog.Infoln("STARTING")
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
			case RECEIVER:
				this.receiver()
			case PAYER:
				this.payer()
			case SHUTTING:
				this.shutting()
			default:
				zlog.Warningf("Wrong state: %u\n", this.state)
			}
		}
	}
}

func (this *FPAY) Startup() {
	zlog.Infoln("FPAY service is starting up.")

	this.datasource.Startup()
	this.monitor.Startup()

	go this.loop()
}

func (this *FPAY) Shutdown() {
	zlog.Infoln("FPAY service is Shutting down.")

	this.in <- 0
	this.monitor.Shutdown()
	this.datasource.Shutdown()
	zlog.Infoln("FPAY service already closed.")
}
