package starknet

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
	pedersenhash "github.com/consensys/gnark-crypto/ecc/stark-curve/pedersen-hash"
)

type AddressType int

var b = new(big.Int)
var UINT_128_MAX = b.Sub(b.Lsh(big.NewInt(1), 128), big.NewInt(1))

const (
	AddressTypeStarknet AddressType = iota
	AddressTypeEthereum
	AddressTypeCustom
)

type Leaf struct {
	AddressType AddressType
	Address     string
	VotingPower big.Int
}

func (l *Leaf) Hash() *big.Int {
	var data []*fp.Element

	item := new(fp.Element).SetInt64(int64(l.AddressType))
	data = append(data, item)

	item, _ = new(fp.Element).SetString(l.Address)
	data = append(data, item)

	low := big.NewInt(0).And(&l.VotingPower, UINT_128_MAX)
	high := big.NewInt(0).Rsh(&l.VotingPower, 128)

	item = new(fp.Element).SetBigInt(low)
	data = append(data, item)

	item = new(fp.Element).SetBigInt(high)
	data = append(data, item)

	hash := pedersenhash.PedersenArray(data...)

	return hash.BigInt(new(big.Int))
}
