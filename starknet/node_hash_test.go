package starknet

import (
	"math/big"
	"testing"
)

func TestNodeHash(t *testing.T) {
	a := big.NewInt(42)
	b := big.NewInt(43)

	got := NodeHash(a, b)

	want, _ := new(big.Int).SetString("0x48d817afbb700a072d6af577e7050cbf4e7b154b8340db50b02b770f7f54bae", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("NodeHash(0x%x, 0x%x) = 0x%x, want 0x%x", a, b, got, want)
	}
}
