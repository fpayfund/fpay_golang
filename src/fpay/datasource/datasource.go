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

package datasource

import (
	"fpay/cli"
	"time"
	"zlog"
)

type DataSource struct {
	in, out chan uint8
}

func New(settings *cli.Settings) (ds *DataSource) {
	ds = new(DataSource)
	ds.in = make(chan uint8, 1)
	ds.out = make(chan uint8, 1)
	return
}

func (this *DataSource) loop() {
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

func (this *DataSource) Startup() {
	zlog.Infoln("DataSource service is starting up.")

	go this.loop()
}

func (this *DataSource) Shutdown() {
	zlog.Infoln("DataSource service is Shutting down.")
	this.in <- 0
	<-this.out
	zlog.Infoln("DataSource service already closed.")
}
