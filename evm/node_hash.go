package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func NodeHash(a, b *big.Int) *big.Int {
	var left, right *big.Int

	if (a.Cmp(b)) > 0 {
		left = b
		right = a
	} else {
		left = a
		right = b
	}

	bytes := append(left.Bytes(), right.Bytes()...)

	return new(big.Int).SetBytes(crypto.Keccak256(bytes))
}
