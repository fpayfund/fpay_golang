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
	"fmt"
	"fpay"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zlog"
)

func main() {
	zlog.SetLevel(zlog.TRACE)
	zlog.SetTagLevel(zlog.SILENCE, "fpay/(*Core)")
	rand.Seed(time.Now().UnixNano())

	settings, err := fpay.ParseCLI()

	if err != nil {
		panic("Commandline params parse failed: " + err.Error())
	}

	if settings.NewAccount {
		fmt.Println(fpay.NewAccount().ToJson())
		return
	} else if settings.AccountPath != "" {
		ac, _ := fpay.LoadAccount(settings.AccountPath)
		fmt.Println(ac.ToJson())
		return
	}

	zlog.Infoln("FPAY is starting up.")
	defer zlog.Infoln("FPAY is shutdown.")

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	defer signal.Stop(osSignal)

	// TODO: 设置settings参数

	fpayService, err := fpay.New(settings)
	if err != nil {
		return
	}

	err = fpayService.Startup()
	if err != nil {
		return
	}

	defer fpayService.Shutdown()

	<-osSignal
	return
}
