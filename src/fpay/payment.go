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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/big"
	"time"
	"zlog"
)

// 1.1.4. 支付请求
type Payment struct {
	version            []byte /* 4位协议版本，结合协议类型确定协议处理方式 */
	encryption         uint16 /* 加密类型，确定PublicKey和Signature处理方式 */
	from               []byte /* 20位支付发起者地址 */
	publicKey          []byte /* 64位支付者公钥 */
	timestamp          uint64 /* 支付时间戳。仅供参考 */
	amount             uint64 /* 支付金额 */
	currency           uint32 /* 货币类型。0为FPAY，其余为各种FxToken */
	currencyVersion    []byte /* 4位余额版本。在等比调节后结算跨版本以后的真实余额 */
	balanceSnapshot    uint64 /* 余额快照。对应钱包地址和支付类型 */
	balanceBlockNounce uint64 /* 余额区块高度 */
	balanceVersion     []byte /* 4位余额版本 */
	to                 []byte /* 20位支付接收者地址 */
	nextWorker         []byte /* 20位下个节点地址 */
	signature          []byte /* 64位签名，包含前面所有数据 */
	nextWorkerIP       []byte /* 4位下一个工作节点IP。只广播不存储 */
	nextWorkerPort     uint32 /* 下一个工作节点端口。只广播不存储 */
	confirmationAmount uint8  /* 确认数量，存储但不校验 */
}

// 创建一个Payment
func PaymentNew(ctx *FPAY, from *Account, encryption uint16, to []byte, amount uint64, currency uint32) (p *Payment, isEnough bool) {
	var ok bool
	p = new(Payment)
	p.version = ctx.Version
	p.encryption = encryption
	p.from = from.Address
	p.publicKey = from.PublicKey
	p.timestamp = uint64(time.Now().Unix())
	p.amount = amount
	p.currency = currency
	p.balanceSnapshot, p.balanceBlockNounce, p.balanceVersion, ok = ctx.DB.GetAccountCurrencyBalance(from, currency)
	// 如果账号不存在或者余额不足，则支付失败
	if !ok || p.balanceSnapshot < amount {
		return nil, false
	}
	p.to = to
	p.nextWorker = ctx.Parent.Account.Address[:20]
	p.signature = p.Sign(ctx)
	p.nextWorkerIP = ctx.Parent.Raddr.IP[:4]
	p.nextWorkerPort = uint32(ctx.Parent.Raddr.Port)
	return true
}

func PaymentUnmarshal(reader io.Reader) (p *Payment, err error) {
	buf := make([]byte, 227)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		zlog.Warningln("io.ReadFull failed: " + err.Error())
		return
	}

	p = new(Payment)
	p.version = buf[:4]
	encryption, _ := binary.Uvarint(buf[4:6])
	p.encryption = uint16(encryption)
	p.from = buf[6:26]
	p.publicKey = buf[26:90]
	p.timestamp, _ = binary.Uvarint(buf[90:98])
	p.amount, _ = binary.Uvarint(buf[98:106])
	currency, _ := binary.Uvarint(buf[106:110])
	p.currency = uint32(currency)
	p.currencyVersion = buf[110:114]
	p.balanceSnapshot, _ = binary.Uvarint(buf[114:122])
	p.balanceBlockNounce, _ = binary.Uvarint(buf[122:130])
	p.balanceVersion = buf[130:134]
	p.to = buf[134:154]
	p.signature = buf[154:218]
	p.nextWorkerIP = buf[218:222]
	nextWorkerPort, _ := binary.Uvarint(buf[222:226])
	p.nextWorkerPort = uint32(nextWorkerPort)
	confirmationAmount, _ := binary.Uvarint(buf[226:227])
	p.confirmationAmount = uint8(confirmationAmount)
	return
}

func (this *Payment) Marshal(writer io.Writer) (err error) {
	buf := make([]byte, 227)
	copy(buf, this.version[:4])
	binary.PutUvarint(buf[4:], uint64(this.encryption))
	copy(buf[6:], this.from[:20])
	copy(buf[26:], this.publicKey[:64])
	binary.PutUvarint(buf[90:], this.timestamp)
	binary.PutUvarint(buf[98:], this.amount)
	binary.PutUvarint(buf[106:], uint64(this.currency))
	copy(buf[110:], this.currencyVersion[:4])
	binary.PutUvarint(buf[114:], this.balanceSnapshot)
	binary.PutUvarint(buf[122:], this.balanceBlockNounce)
	copy(buf[130:], this.balanceVersion[:4])
	copy(buf[134:], this.to[:20])
	copy(buf[154:], this.signature[:64])
	copy(buf[218:], this.nextWorkerIP[:4])
	binary.PutUvarint(buf[222:], uint64(this.nextWorkerPort))
	binary.PutUvarint(buf[226:], uint64(this.confirmationAmount))
	_, err = writer.Write(buf)
	return
}

func (this *Payment) ToAbstract() (abstract []byte) {
	buf := make([]byte, 154)
	copy(buf, this.version[:4])
	binary.PutUvarint(buf[4:], uint64(this.encryption))
	copy(buf[6:], this.from[:20])
	copy(buf[26:], this.publicKey[:64])
	binary.PutUvarint(buf[90:], this.timestamp)
	binary.PutUvarint(buf[98:], this.amount)
	binary.PutUvarint(buf[106:], uint64(this.currency))
	copy(buf[110:], this.currencyVersion[:4])
	binary.PutUvarint(buf[114:], this.balanceSnapshot)
	binary.PutUvarint(buf[122:], this.balanceBlockNounce)
	copy(buf[130:], this.balanceVersion[:4])
	copy(buf[134:], this.to[:20])
	a := sha256.Sum256(buf)
	abstract = a[:32]
	return
}

func (this *Payment) Sign(ctx *FPAY) (signature []byte) {
	r, s, err := ecdsa.Sign(bytes.NewReader(ctx.Parent.Account.Random), ctx.Parent.Account.ToPrivateKey(), this.ToAbstract())

	if err != nil {
		return nil
	}

	signature = make([]byte, 0, 64)
	copy(signature, r.Bytes()[:32])
	copy(signature[32:], s.Bytes()[:32])
	return
}

func (this *Payment) Verify(p []byte) bool {
	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P256()

	pub.X = new(big.Int)
	pub.X.SetBytes(this.publicKey[:32])

	pub.Y = new(big.Int)
	pub.Y.SetBytes(this.publicKey[32:64])

	r := new(big.Int)
	r.SetBytes(this.signature[:32])

	s := new(big.Int)
	s.SetBytes(this.signature[32:64])
	return ecdsa.Verify(pub, this.ToAbstract(), r, s)
}
