package evm

import (
	"math/big"
	"testing"
)

func TestNodeHash(t *testing.T) {
	a, _ := new(big.Int).SetString("0x2df2ce6efd5635498c6fbc578f885d6ace29517e6f11e507d1c17dcb86d9ddd4", 0)
	b, _ := new(big.Int).SetString("0x042f4dc68248f096de5b373868763012bca8ff8c67af7fdbf501da62f38d02cc", 0)

	got := NodeHash(a, b)

	want, _ := new(big.Int).SetString("0xa3efc99720053662d6a13e9ee82b3bd977d544474ff4988e228b511ed7876791", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("NodeHash(0x%x, 0x%x) = 0x%x, want 0x%x", a, b, got, want)
	}
}
