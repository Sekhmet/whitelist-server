package starknet

import (
	"math/big"
	"testing"
)

func TestHash(t *testing.T) {
	leaf := Leaf{
		AddressType: AddressTypeEthereum,
		Address:     "0x556B14CbdA79A36dC33FcD461a04A5BCb5dC2A70",
		VotingPower: *big.NewInt(42),
	}

	got := leaf.Hash()
	want, _ := new(big.Int).SetString("0x196903245bb2dcafaf9acc391de440ce08a8853b7b1dcbfc670171bb255e119", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("Leaf.Hash() = 0x%x, want 0x%x", got, want)
	}
}
