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

package fpay /* 节点 */

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"zlog"
)

type Account struct {
	EncryptType    uint16     /* 加密类型 */
	PrivateKey     [32]byte   /* 私钥 */
	PublicKey      [64]byte   /* 公钥 */
	Address        [20]byte   /* 地址 */
	MnemonicsWords [32]string /* 助记词 */
}

var MnemonicsWords = []string{
	"a", "an", "AM", "are", "ask", "age", "and", "ago", "arm", "aid", "add",
	"be", "big", "beg", "bar", "bid", "bag", "ben", "bed", "bad", "blue", "born",
	"cup", "cap", "col", "con", "can", "car", "com", "cold", "camp", "calm",
	"do", "die", "due", "dog", "don", "dim", "deg", "did", "def", "deep",
	"eg", "ear", "eph", "end", "eye", "ease", "earn", "echo", "earl", "emma",
	"fun", "foe", "far", "fed", "fro", "fee", "fire", "free", "feel", "fell",
	"go", "gen", "gov", "god", "gun", "gold", "glad", "game", "gate", "gain",
	"he", "hid", "heb", "had", "him", "high", "hear", "help", "held", "hold",
	"I", "if", "in", "ie", "inn", "imp", "ill", "ind", "inc", "ice", "isa",
	"jos", "job", "jim", "jer", "juan", "joke", "jump", "jove", "jail", "jeff",
	"kick", "khan", "karl", "kirk", "keel", "knob", "kine", "kite", "kong", "kith",
	"leg", "lee", "ltd", "led", "lev", "lie", "lad", "line", "late", "laid",
	"me", "mud", "man", "men", "mad", "mere", "move", "main", "moon", "mile",
	"no", "non", "ned", "nat", "now", "nine", "nose", "noon", "nigh", "nile",
	"on", "oil", "odd", "owe", "oak", "oid", "ohg", "one", "own", "old", "off",
	"pm", "pen", "pub", "phd", "pre", "pro", "paul", "pope", "poem", "pipe",
	"quid", "quod", "quito", "qualm", "quaff", "quire", "quake", "qualf", "quash", "quaxk",
	"rid", "run", "red", "rom", "ran", "rise", "rule", "rock", "rate", "rain",
	"so", "sad", "sin", "sum", "sam", "ste", "she", "see", "son", "sat", "sea", "sun",
	"to", "top", "tea", "the", "two", "too", "ten", "tom", "thee", "town",
	"up", "us", "use", "unco", "undo", "unto", "upon", "urge", "used", "utah",
	"vol", "viz", "veil", "vice", "void", "vale", "vase", "vera", "veal", "veto",
	"we", "wit", "won", "who", "war", "win", "web", "wish", "wind", "wild",
	"yule", "yolk", "yuan", "yelp", "yarn", "yawl", "yank", "yowl", "year", "york",
	"zeal", "zone", "zion", "zinc", "zero", "zola", "zeke", "zing", "zoom", "zebra"}

func NewAccount() (account *Account) {
	random := make([]byte, 40)
	_, err := rand.Read(random)
	if err != nil {
		panic("crypto/rand.Read failure: " + err.Error())
	}

	account = new(Account)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(random))
	if err != nil {
		panic("crypto/ecdsa.GenerateKey failure: " + err.Error())
	}

	for n, v := range privateKey.D.Bytes() {
		account.PrivateKey[n] = v
	}

	for n, v := range privateKey.PublicKey.X.Bytes() {
		account.PublicKey[n] = v
	}

	for n, v := range privateKey.PublicKey.Y.Bytes() {
		account.PublicKey[n+32] = v
	}

	account.Address = sha1.Sum(account.PublicKey[:]) // TODO: 还要加一层ripemd160

	for n, v := range account.PrivateKey {
		account.MnemonicsWords[n] = MnemonicsWords[v]
	}
	zlog.Traceln(account)
	return
}
