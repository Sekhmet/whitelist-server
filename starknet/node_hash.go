package starknet

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
	pedersenhash "github.com/consensys/gnark-crypto/ecc/stark-curve/pedersen-hash"
)

func NodeHash(a, b *big.Int) *big.Int {
	var left, right *big.Int

	if (a.Cmp(b)) > 0 {
		left = a
		right = b
	} else {
		left = b
		right = a
	}

	leftFp := new(fp.Element).SetBigInt(left)
	rightFp := new(fp.Element).SetBigInt(right)

	fp := pedersenhash.Pedersen(leftFp, rightFp)
	return fp.BigInt(new(big.Int))
}
