package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	Address, _ = abi.NewType("address", "", nil)
	Uint96, _  = abi.NewType("uint96", "", nil)

	args = abi.Arguments{
		{
			Name: "address",
			Type: Address,
		},
		{
			Name: "votingPower",
			Type: Uint96,
		},
	}
)

type Leaf struct {
	Address     string
	VotingPower big.Int
}

func (l *Leaf) Hash() *big.Int {
	data, _ := args.Pack(common.HexToAddress(l.Address), &l.VotingPower)

	return new(big.Int).SetBytes(crypto.Keccak256(crypto.Keccak256(data)))
}
