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
	"errors"
	"fmt"
	"zlog"
)

type CoreI interface {
	PreLoop() error
	Loop() bool
	AftLoop()
	Run()
	Startup() error
	Shutdown()
}

type Core struct {
	CoreI
	Command, State chan uint8
}

const (
	CMD_SHUT = iota
)

const (
	STATE_READY = iota
	STATE_FAILED
	STATE_CLOSED
)

var CMDS []string = []string{"CMD_SHUT"}
var STATES []string = []string{"STATE_READY", "STATE_FAILED", "STATE_CLOSED"}

func (this *Core) Init(c CoreI) {
	this.CoreI = c
	this.Command = make(chan uint8, 1)
	this.State = make(chan uint8, 1)
}

func (this *Core) Run() {
	err := this.PreLoop()
	if err != nil {
		this.State <- STATE_FAILED
		zlog.Errorln("PreLoop failed: " + err.Error())
		return
	}

	this.State <- STATE_READY

	for {
		ok := this.Loop()
		if !ok {
			break
		}
	}

	this.AftLoop()

	this.State <- STATE_CLOSED
}

func (this *Core) Startup() (err error) {
	zlog.Infoln("Starting up.")

	go this.Run()

	s := <-this.State
	switch s {
	case STATE_READY:
		zlog.Infoln("Started.")
		return
	default:
		zlog.Debugf("%s received.\n", STATES[s])
		zlog.Errorln("Failed to start.")

		return errors.New(fmt.Sprintf("Unexpected state %s received.", STATES[s]))
	}
}

func (this *Core) Shutdown() {
	zlog.Infoln("Shutting down.")
	this.Command <- CMD_SHUT
	zlog.Traceln("CMD_SHUT sent.")

	s := <-this.State
	switch s {
	case STATE_CLOSED:
		zlog.Infoln("Closed.")
	default:
		zlog.Debugf("%s received.\n", STATES[s])
		zlog.Warningln("Closed abnormally.")
	}
}
