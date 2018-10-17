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

package base

import (
	"crypto/rand"
	"errors"
	"io"
	"zlog"
)

// 基类
type Base struct {
	Name     []byte /* "FPAY"字符 */
	Version  []byte /* 版本号 */
	Id       []byte /* id，随机数 */
	Protocol []byte /* 协议类型，用于扩展。不同的协议用不同的插件处理，长度为15 */
}

var PROTOCOL_NAME string = "FPAY"
var PROTOCOL_VERSION []byte = []byte{0, 0, 1, 0}

// 初始化Core，给子类使用
func (this *Base) Init(s string) {
	this.Protocol = []byte(s)
	this.Name = []byte(PROTOCOL_NAME)
	this.Version = PROTOCOL_VERSION

	this.Id = make([]byte, 32)
	_, err := rand.Read(this.Id)
	if err != nil {
		panic("crypto/rand.Read failure: " + err.Error())
	}
}

// 创建一个Core
func New(s string) (c *Base) {
	c = new(Base)
	c.Init(s)
	return
}

// 反序列化，给子类使用
func Unmarshal(reader io.Reader) (c *Base, err error) {
	buf := make([]byte, 56)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		zlog.Warningln("io.ReadFull failed: " + err.Error())
		return
	}

	c = new(Base)
	c.Name = buf[:4]
	name := string(c.Name)
	if name != PROTOCOL_NAME {
		errors.New("Unsupport protocol: " + name)
		zlog.Warningf("Wrong c.Name: [%v %v %v %v].\n", c.Name[0], c.Name[1], c.Name[2], c.Name[3])
		return
	}

	c.Version = buf[4:8]
	c.Id = buf[8:40]
	c.Protocol = buf[40:56]

	zlog.Tracef("Base:{Name:%s, Version:%v, Id:%v, Protocol:%v} unmarshaled.\n", string(c.Name), c.Version, c.Id, string(c.Protocol))

	return
}

func (this *Base) Marshal(writer io.Writer) (err error) {
	buf := make([]byte, 56)
	copy(buf, this.Name)
	copy(buf[4:8], this.Version)
	copy(buf[8:40], this.Id)
	copy(buf[40:56], this.Protocol)

	_, err = writer.Write(buf)
	return
}
