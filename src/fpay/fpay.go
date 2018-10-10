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
	"fpay/cache"
	"fpay/node"
	"zlog"
)

func Run(settings *Settings) (err error) {
	zlog.Infoln("FPAY Service is starting up.")

	cacheService, err := cache.New()
	if err != nil || cacheService == nil {
		zlog.Fatalln("Cache service Create failed: " + err.Error())
		return
	}

	// TODO: 设定cache 参数

	cacheService.Startup()
	message, ok := <-cacheService.MessageOut

	if ok && message == cache.READY {
		err = node.Run(settings, cacheService)

	} else {
		if !ok {
			zlog.Errorln("The MessageOut chan closed unexpectedly.")
		} else {
			zlog.Errorln("The cache state is incorrect.")
		}

		zlog.Fatalln("Cache service startup failed.")

	}

	if err != nil {
		zlog.Fatalln("FPAY node startup failed.")
	}

	cacheService.Shutdown()

	zlog.Infoln("FPAY Service is shutdown.")
	return
}
