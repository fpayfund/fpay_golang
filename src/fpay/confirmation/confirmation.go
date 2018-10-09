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

package confirmation

// 1.1.2. 确认请求
type Confirmation struct {
	encryptType              uint16   /* 加密类型，确定PublicKey和Signature处理方式 */
	version                  [4]byte  /* 协议版本，确定确认的处理方式 */
	address                  [20]byte /* 工作节点地址 */
	publicKey                [64]byte /* 工作节点公钥 */
	blockNounce              uint64   /* 工作节点区块高度 */
	balanceSnapshot          uint64   /* 余额快照，工作节点的查询结果 */
	balanceBlockNounce       uint64   /* 余额快照对应的区块高度 */
	toBalanceSnapshot        uint64   /* 接收节点余额快照，工作节点的查询结果 */
	toBalanceBlockNounce     uint64   /* 接收节点余额快照对应的区块高度 */
	toFPAYBalanceSnapshot    uint64   /* 接收节点FPAY余额快照，工作节点的查询结果 */
	toFPAYBalanceBlockNounce uint64   /* 接收节点FPAY余额快照对应的区块高度 */
	workerBalanceSnapshot    uint64   /* 工作节点FPAY余额快照 */
	nextWorker               [20]byte /* 下一个工作节点的地址。用于避免女巫攻击 */
	status                   uint8    /* 检查结果，0为通过，1为放弃，其余为失败 */
	signature                [64]byte /* 签名，包含前面所有数据 */
	nextWorkerIP             [4]byte  /* 下一个工作节点IP。只广播不存储 */
	nextWorkerPort           uint32   /* 下一个工作节点端口。只广播不存储 */
}
