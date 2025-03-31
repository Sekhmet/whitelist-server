package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Sekhmet/whitelist-server/starknet"
)

func getTree() []*big.Int {
	var leaves []Leaf

	for i := range 20 {
		value := int64(i + 1)

		leaf := &starknet.Leaf{
			AddressType: starknet.AddressTypeStarknet,
			Address:     fmt.Sprintf("0x%x", value),
			VotingPower: *big.NewInt(value),
		}
		leaves = append(leaves, leaf)
	}

	return GenerateMerkleTree(leaves)
}

func TestGenerateMerkleTree(t *testing.T) {
	tree := getTree()

	got := tree[0]
	want, _ := new(big.Int).SetString("0xbfddde52fc7d24a63693fb4dfa257571238e2d654aecbe6bc26f067e770bc5", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("Merkletree root = 0x%x, want 0x%x", got, want)
	}
}

func TestGetMerkleProof(t *testing.T) {
	tree := getTree()

	output := []string{
		"0x3eca1772359b7a5b248088472ef392716c034e899c510d6e02c0c97704164ab",
		"0x1919a163ca6cb8d28728b24847c269529cf3af4caafb0a2b3e2fd19715f1a8b",
		"0x7a99182fabd949861d469e6e6143187c71ced341c2168de0f62e8e025d76dfc",
		"0x2f0a9bf11b7f4792a3b732f518168d960cfa73d244c020cb88c80293b4fff91",
		"0x2f5a19d2f01021cbd28d1073a68b4123cc56c5811d6da8dea6e0c4c921c0c21",
	}

	want := make([]*big.Int, len(output))
	for i, o := range output {
		want[i], _ = new(big.Int).SetString(o, 0)
	}

	leafIndex := 2
	got, err := GetMerkleProof(tree, leafIndex)
	if err != nil {
		t.Fatalf("GetMerkleProof() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("len(GetMerkleProof()) = %v, want %v", len(got), len(want))
	}

	for i := range got {
		if got[i].Cmp(want[i]) != 0 {
			t.Errorf("GetMerkleProof()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
