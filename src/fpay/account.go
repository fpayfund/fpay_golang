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
	"encoding/json"
	"errors"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"math/big"
)

type Account struct {
	EncryptType    uint16   /* 加密类型 */
	Random         []byte   /* 40位随机数 */
	PrivateKey     []byte   /* 32位私钥 */
	PublicKey      []byte   /* 64位公钥 */
	Address        []byte   /* 20位地址 */
	MnemonicsWords []string /* 32个助记词 */
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

func NewWrongDataLengthError() (err error) {
	return errors.New("The length of data is wrong.")
}

func AddressGenerate(pubKey []byte) (addr []byte, err error) {
	if len(pubKey) != 64 {
		err = NewWrongDataLengthError()
		return
	}

	addr = ripemd160.New().Sum(sha1.New().Sum(pubKey)[64:])[20:]
	return
}

func AccountGenerate(random []byte) (a *Account) {
	a = new(Account)
	a.Random = random
	a.MnemonicsWords = make([]string, 0, 40)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(a.Random))
	if err != nil {
		panic("crypto/ecdsa.GenerateKey failure: " + err.Error())
	}

	a.PrivateKey = privateKey.D.Bytes()
	a.PublicKey = make([]byte, 64, 64)
	copy(a.PublicKey, privateKey.X.Bytes())
	copy(a.PublicKey[32:], privateKey.Y.Bytes())
	a.Address, _ = AddressGenerate(a.PublicKey)

	for _, v := range a.Random {
		a.MnemonicsWords = append(a.MnemonicsWords, MnemonicsWords[v])
	}
	return
}

func AccountNew() (a *Account) {
	random := make([]byte, 40)
	_, err := rand.Read(random)
	if err != nil {
		panic("crypto/rand.Read failure: " + err.Error())
	}

	return AccountGenerate(random)
}

func AccountLoad(path string) (ac *Account, err error) {
	fs, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	ac = new(Account)
	err = json.Unmarshal(fs, ac)
	if err != nil {
		return
	}

	if ac.Random == nil && ac.MnemonicsWords == nil {
		err = errors.New("No private data or mnemonics words found.")
		return
	}

	if ac.Random == nil {
		if len(ac.MnemonicsWords) != 40 {
			err = NewWrongDataLengthError()
			return
		}

		ac.Random = make([]byte, 0, 40)
		var n int
		var w, gw string
		for _, w = range ac.MnemonicsWords {
			for n, gw = range MnemonicsWords {
				if w == gw {
					break
				}
			}
			ac.Random = append(ac.Random, byte(n))
		}
	}

	if len(ac.Random) < 40 {
		err = errors.New("The length of data is wrong.")
		return
	}

	ac = AccountGenerate(ac.Random)
	return
}

func AccountsNew(n int) (acs []*Account) {
	acs = make([]*Account, 0, n)
	for i := 0; i < n; i++ {
		acs = append(acs, AccountNew())
	}
	return
}

func AccountsToJson(acs []*Account) (js string) {
	if len(acs) == 1 {
		return acs[0].ToJson()
	}
	bytes, _ := json.Marshal(acs)
	js = string(bytes)
	return
}

func (this *Account) ToJson() (j string) {
	a, _ := json.Marshal(this)
	j = string(a)
	return
}

func (this *Account) ToPrivateKey() (priv *ecdsa.PrivateKey) {
	priv = new(ecdsa.PrivateKey)
	priv.D = new(big.Int)
	priv.D = priv.D.SetBytes(this.PrivateKey)
	return
}

func (this *Account) ToPublicKey() (pub *ecdsa.PublicKey) {
	pub = new(ecdsa.PublicKey)
	pub.Curve = elliptic.P256()
	pub.X = new(big.Int)
	pub.X.SetBytes(this.PublicKey[:32])
	pub.Y = new(big.Int)
	pub.Y.SetBytes(this.PublicKey[32:64])
	return
}
