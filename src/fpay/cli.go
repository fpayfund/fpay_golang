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
	"os"
	"strings"
	"zlog"
)

type Settings struct {
	Args    []string
	TCPAddr string
}

// 命令行参数解释
// -T:  tcp协议监听地址
//      用法: -T0.0.0.0:8080 或者 -T 0.0.0.0:8080
func Parse() (settings *Settings, err error) {
	settings = new(Settings)
	settings.Args = os.Args

	// TODO: 正式解释参数
	for i := 0; i < len(settings.Args); i++ {
		arg := settings.Args[i]
		l := len(arg)

		if i == 0 || l <= 1 || (!strings.HasPrefix(arg, "-")) {
			continue
		}

		switch []byte(arg)[1] {
		case []byte("T")[0]:
			if l == 2 {
				i++
				settings.TCPAddr = settings.Args[i]
			} else {
				settings.TCPAddr = string([]byte(arg)[2:l])
			}
			zlog.Tracef("settings.TCPAddr=%s\n", settings.TCPAddr)
		}
	}
	return
}
