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
	"github.com/go-redis/redis"
)

type AccountCachedKey struct {
	address  []byte
	currency uint32
}

type AccountCachedValue struct {
	balance            uint64
	balanceBlockNounce uint64
	balanceVersion     []byte
}

type Cache struct {
	addr   string
	client redis.Client
}

func CacheGet() (c *Cache) {
	return
}

// 通过地址和余额类型查询地址余额————校验支付
func (this *Cache) GetAccountCurrencyBalance(from *Account, currency uint32) (uint64, uint64, []byte, bool) {
	return 0, 0, nil, false
}

// 通过顺序查询区块内容————校验余额
func (this *Cache) GetBlockBySequence(sequence uint64) (*Block, bool) {
	return nil, false
}

// 获取最近一段数量的区块————计算节点权重
func (this *Cache) GetLastBlocks(currency uint32, sequence uint64, amount uint32) ([]*Block, bool) {
	return nil, false
}

// 获取未处理支付信息————检查漏发攻击
func (this *Cache) GetUnprocessedPayments() []*Payment {
	return nil
}

func (this *Cache) Startup() (err error) {
	return nil
}

func (this *Cache) Shutdown() {

}
