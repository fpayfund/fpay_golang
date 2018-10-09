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

package core

import (
	"crypto/rand"
)

// 基类
type Core struct {
	id       [32]byte /* id，随机数 */
	protocol uint16   /* 协议类型，用于扩展，不同的协议用不同的插件处理 */
}

// 初始化Core，给子类使用
func (this *Core) Init(protocol uint16) {
	this.protocol = protocol

	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		panic("crypto/rand.Read failure: " + err.Error())
	}

	for i := 0; i < 32; i++ {
		this.id[i] = id[i]
	}
}

// 序列化，给子类使用
func (this *Core) ToByte() (data []byte) {
	// TODO
	panic("Unsupport function")
}

// 创建一个Core
func Create(protocol uint16) (core *Core) {
	core = new(Core)
	core.Init(protocol)
	return
}

// 反序列化，给子类使用
func Parse(data []byte) (core *Core) {
	// TODO
	panic("Unsupport function")
}
