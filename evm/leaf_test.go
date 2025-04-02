package evm

import (
	"math/big"
	"testing"
)

func TestHash(t *testing.T) {
	leaf := Leaf{
		Address:     "0x556B14CbdA79A36dC33FcD461a04A5BCb5dC2A70",
		VotingPower: *big.NewInt(21),
	}

	got := leaf.Hash()
	want, _ := new(big.Int).SetString("0xd8c29f38c935b4a569d48ffec67aa6247c90b6598fea89d7bd9415ac50ed7acc", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("Leaf.Hash() = 0x%x, want 0x%x", got, want)
	}
}
