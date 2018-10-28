package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"math/big"
)

func main() {
	ran := make([]byte, 40, 40)
	rand.Read(ran)
	prvKey, err := ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(ran))
	if err != nil {
		panic("ecdsa.GenerateKey:" + err.Error())
	}
	fmt.Println("ran", ran, len(ran))
	d := prvKey.D.Bytes()
	fmt.Println("d", d, len(d))
	x := prvKey.X.Bytes()
	fmt.Println("x", x, len(x))
	y := prvKey.Y.Bytes()
	fmt.Println("y", y, len(y))

	h := sha1.Sum([]byte("Hello WorldHello World"))
	fmt.Println("h", h, len(h))
	hh := h[:len(h)]
	fmt.Println("hh", hh, len(hh))

	r, s, err := ecdsa.Sign(bytes.NewReader(ran), prvKey, hh)
	if err != nil {
		panic("ecdsa.Sign:" + err.Error())
	}
	rb := r.Bytes()
	fmt.Println("rb", rb, len(rb))
	sb := s.Bytes()
	fmt.Println("sb", sb, len(sb))

	ri := new(big.Int)
	ri.SetBytes(rb)
	si := new(big.Int)
	si.SetBytes(sb)
	v := ecdsa.Verify(&prvKey.PublicKey, hh, ri, si)
	fmt.Println("v", ri, si, v)
}
