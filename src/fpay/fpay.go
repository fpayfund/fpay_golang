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

type FPAY struct {
	in, out    chan uint8
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

func (this *FPAY) loop() {
	for {
		select {
		case <-this.in:
			this.out <- 0
			break
		case <-time.After(10 * time.Second):
			zlog.Infoln("10s runout")
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
	<-this.out
	zlog.Infoln("FPAY service already closed.")
}
