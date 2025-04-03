package main

import (
	"fmt"

	"math/big"
	"testing"

	"github.com/Sekhmet/whitelist-server/evm"
	"github.com/Sekhmet/whitelist-server/starknet"
)

func getStarnetTree() []*big.Int {
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

	return GenerateMerkleTree(leaves, starknet.NodeHash, false)
}

func getEvmTree() []*big.Int {
	var leaves []Leaf

	for i := range 20 {
		leaf := &evm.Leaf{
			Address:     fmt.Sprintf("0x%x", i),
			VotingPower: *big.NewInt(int64(i)),
		}
		leaves = append(leaves, leaf)
	}

	return GenerateMerkleTree(leaves, evm.NodeHash, true)
}

func TestGenerateMerkleTreeStarknet(t *testing.T) {
	tree := getStarnetTree()

	got := tree[0]
	want, _ := new(big.Int).SetString("0xbfddde52fc7d24a63693fb4dfa257571238e2d654aecbe6bc26f067e770bc5", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("Merkletree root = 0x%x, want 0x%x", got, want)
	}
}

func TestGetMerkleProofStarknet(t *testing.T) {
	tree := getStarnetTree()

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

func TestGenerateMerkleTreeEvm(t *testing.T) {
	tree := getEvmTree()

	got := tree[0]
	want, _ := new(big.Int).SetString("0xacc4e698f6fd65fbe01a589f241bbdae1739e556bf2e0510b844acbdecda2fdf", 0)

	if got.Cmp(want) != 0 {
		t.Errorf("Merkletree root = 0x%x, want 0x%x", got, want)
	}
}

func TestGetMerkleProofEvm(t *testing.T) {
	tree := getEvmTree()

	output := []string{
		"0xae1f6f36060f166f063fb01d63adab80297f56b5a444cab19384c535141dbd8b",
		"0xd152f7eda6853d2bfe5583d51fc737d9e0c2689b23b67dc2d8711acda31bc6f6",
		"0x390dc37559610bf6e51c8c3672ad37d054380df794014e600c1a8684836d5027",
		"0xaf1879534a42e21defe29b1acf1c2b3f5ab1cdffba47d49d86f4330902cd7e8c",
	}

	want := make([]*big.Int, len(output))
	for i, o := range output {
		want[i], _ = new(big.Int).SetString(o, 0)
	}

	// This is updated index after sorting.
	// The original index was 5. Needs to be updated after sorting is fully implemented.
	leafIndex := 10
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
