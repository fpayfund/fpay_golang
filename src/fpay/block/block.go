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

package block

import (
	"fpay/base"
)

type Block struct {
	blockId   [32]byte /* Id，由上一区块的Id和本区块的签名组合而成 */
	sequence  uint64   /* 区块顺序 */
	worker    [20]byte /* 出块者地址 */
	publicKey [64]byte /* 出块者公钥 */
	timestamp uint32   /* 出块时间戳。仅供参考 */
	signature [64]byte /* 签名，包含前面所有数据 */
}

// 创建一个Block
func New(last *Block) (b *Block) {
	b = new(Block)
}

// 反序列化，给子类使用
func Unmarshal(reader io.Reader) (c *Block, err error) {
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
