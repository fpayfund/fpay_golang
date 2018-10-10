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

package main

import (
	"fpay/cli"
	"fpay/monitor"
	"os"
	"os/signal"
	"syscall"
	"zlog"
)

func main() {
	_, err := cli.Parse()

	if err != nil {
		panic("Commandline params parse failed: " + err.Error())
	}

	zlog.Infoln("FPAY is starting up.")
	defer zlog.Infoln("FPAY is shutdown.")

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	defer signal.Stop(osSignal)

	// TODO: 设置monitor参数
	//
	ms := monitor.New()

	ms.Startup()
	defer ms.Shutdown()

	<-osSignal
	return
}
