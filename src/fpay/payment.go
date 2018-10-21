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

// 1.1.4. 支付请求
type Payment struct {
	Core
	version            [4]byte  /* 协议版本，结合协议类型确定协议处理方式 */
	encryption         uint16   /* 加密类型，确定PublicKey和Signature处理方式 */
	from               [20]byte /* 支付发起者地址 */
	publicKey          [64]byte /* 支付者公钥 */
	timestamp          uint32   /* 支付时间戳。仅供参考 */
	amount             uint64   /* 支付金额 */
	currency           uint32   /* 货币类型。0为FPAY，其余为各种FxToken */
	currencyVersion    [4]byte  /* 余额版本。在等比调节后结算跨版本以后的真实余额 */
	balanceSnapshot    uint64   /* 余额快照。对应钱包地址和支付类型 */
	balanceBlockNounce uint64   /* 余额区块高度 */
	balanceVersion     [4]byte  /* 余额版本 */
	to                 [20]byte /* 支付接收者地址 */
	nextWorker         [20]byte /* 下一个工作节点的地址。用于避免女巫攻击 */
	signature          [64]byte /* 签名，包含前面所有数据 */
	nextWorkerIP       [4]byte  /* 下一个工作节点IP。只广播不存储 */
	nextWorkerPort     uint32   /* 下一个工作节点端口。只广播不存储 */
	confirmationAmount uint8    /* 确认数量 */
}
