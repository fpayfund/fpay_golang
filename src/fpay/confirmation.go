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

// 1.1.2. 确认请求
type Confirmation struct {
	version                  []byte /* 4位协议版本，确定确认的处理方式 */
	encryption               uint16 /* 加密类型，确定PublicKey和Signature处理方式 */
	address                  []byte /* 20位工作节点地址 */
	publicKey                []byte /* 64位工作节点公钥 */
	blockNounce              uint64 /* 工作节点区块高度 */
	balanceSnapshot          uint64 /* 余额快照，工作节点的查询结果 */
	balanceBlockNounce       uint64 /* 余额快照对应的区块高度 */
	toBalanceSnapshot        uint64 /* 接收节点余额快照，工作节点的查询结果 */
	toBalanceBlockNounce     uint64 /* 接收节点余额快照对应的区块高度 */
	toFPAYBalanceSnapshot    uint64 /* 接收节点FPAY余额快照，工作节点的查询结果 */
	toFPAYBalanceBlockNounce uint64 /* 接收节点FPAY余额快照对应的区块高度 */
	workerBalanceSnapshot    uint64 /* 工作节点FPAY余额快照 */
	nextWorker               []byte /* 20位下一个工作节点的地址。用于避免女巫攻击 */
	status                   uint8  /* 检查结果，0为通过，1为放弃，其余为失败 */
	signature                []byte /* 64位签名，包含前面所有数据 */
	nextWorkerIP             []byte /* 4位下一个工作节点IP。只广播不存储 */
	nextWorkerPort           uint32 /* 下一个工作节点端口。只广播不存储 */
}

// 创建一个Confirmation
func ConfirmationNew(last *Confirmation) (c *Confirmation) {
	c = new(Confirmation)
	return
}

// 反序列化，给子类使用
func ConfirmationUnmarshal(reader io.Reader) (c *Confirmation, err error) {
	buf := make([]byte, 247)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		zlog.Warningln("io.ReadFull failed: " + err.Error())
		return
	}

	c = new(Confirmation)
	c.version = buf[:4]
	encryption, _ := binary.Uvarint(buf[4:6])
	c.encryption = uint16(encryption)
	c.address = buf[6:26]
	c.publicKey = buf[26:90]
	c.blockNounce, _ = binary.Uvarint(buf[90:98])
	c.balanceSnapshot, _ = binary.Uvarint(buf[98:106])
	c.balanceBlockNounce, _ = binary.Uvarint(buf[106:114])
	c.toBalanceSnapshot, _ = binary.Uvarint(buf[114:122])
	c.toBalanceBlockNounce, _ = binary.Uvarint(buf[122:130])
	c.toFPAYBalanceSnapshot, _ = binary.Uvarint(buf[130:138])
	c.toFPAYBalanceBlockNounce, _ = binary.Uvarint(buf[138:146])
	c.workerBalanceSnapshot, _ = binary.Uvarint(buf[146:154])
	c.nextWorker = buf[154:174]
	status, _ := binary.Uvarint(buf[174:175])
	c.status = uint8(status)
	c.signature = buf[175:239]
	c.nextWorkerIP = buf[239:243]
	nextWorkerPort, _ := binary.Uvarint(buf[243:247])
	c.nextWorkerPort = uint32(nextWorkerPort)

	return
}

func (this *Confirmation) Marshal(writer io.Writer) (err error) {
	buf := make([]byte, 247)
	binary.PutUvarint(buf, uint64(this.encryption))
	copy(buf[2:], this.version)
	copy(buf[6:], this.address)
	copy(buf[26:], this.publicKey)
	binary.PutUvarint(buf[90:], this.blockNounce)
	binary.PutUvarint(buf[98:], this.balanceSnapshot)
	binary.PutUvarint(buf[106:], this.balanceBlockNounce)
	binary.PutUvarint(buf[114:], this.toBalanceSnapshot)
	binary.PutUvarint(buf[122:], this.toBalanceBlockNounce)
	binary.PutUvarint(buf[130:], this.toFPAYBalanceSnapshot)
	binary.PutUvarint(buf[138:], this.toFPAYBalanceBlockNounce)
	binary.PutUvarint(buf[146:], this.workerBalanceSnapshot)
	copy(buf[154:], this.nextWorker)
	binary.PutUvarint(buf[174:], uint64(this.status))
	copy(buf[175:], this.signature)
	copy(buf[239:], this.nextWorkerIP)
	binary.PutUvarint(buf[243:], uint64(this.nextWorkerPort))

	_, err = writer.Write(buf)
	return
}
