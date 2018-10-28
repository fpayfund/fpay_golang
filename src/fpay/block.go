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
	"encoding/binary"
	"io"
	"zlog"
)

type Block struct {
	blockId   []byte /* 32位Id，由上一区块的Id和本区块的签名组合而成 */
	sequence  uint64 /* 区块顺序 */
	worker    []byte /* 20位出块者地址 */
	publicKey []byte /* 64位出块者公钥 */
	timestamp uint32 /* 出块时间戳。仅供参考 */
	signature []byte /* 64位签名，包含前面所有数据 */
}

// 创建一个Block
func BlockNew(a *Account, lb *Block, pms []*Payment) (b *Block) {
	b = new(Block)
	return
}

// 反序列化，给子类使用
func BlockUnmarshal(reader io.Reader) (b *Block, err error) {
	buf := make([]byte, 192)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		zlog.Warningln("io.ReadFull failed: " + err.Error())
		return
	}

	b = new(Block)
	b.blockId = buf[:32]
	b.sequence, _ = binary.Uvarint(buf[32:40])
	b.worker = buf[40:60]
	b.publicKey = buf[60:124]
	timestamp, _ := binary.Uvarint(buf[124:128])
	b.timestamp = uint32(timestamp)
	b.signature = buf[128:192]

	return
}

func (this *Block) Marshal(writer io.Writer) (err error) {
	buf := make([]byte, 192)
	copy(buf, this.blockId)
	binary.PutUvarint(buf[32:], this.sequence)
	copy(buf[40:], this.worker)
	copy(buf[60:], this.publicKey)
	binary.PutUvarint(buf[124:], uint64(this.timestamp))
	copy(buf[128:], this.signature)

	_, err = writer.Write(buf)
	return
}
