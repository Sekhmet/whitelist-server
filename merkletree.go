package main

import (
	"errors"
	"math"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
	pedersenhash "github.com/consensys/gnark-crypto/ecc/stark-curve/pedersen-hash"
)

type Leaf interface {
	Hash() *big.Int
}

func getLeftChildIndex(index int) int {
	return 2*index + 1
}

func getRightChildIndex(index int) int {
	return 2*index + 2
}

func getParentIndex(index int) (int, error) {
	if index > 0 {
		return (index - 1) / 2, nil
	}

	return -1, errors.New("root has no parent")
}

func getSiblingIndex(index int) (int, error) {
	if index > 0 {
		return index - int(math.Pow((-1), float64(index%2))), nil
	}

	return -1, errors.New("root has no siblings")
}

func GenerateMerkleTree(leaves []Leaf) []*big.Int {
	if len(leaves) == 0 {
		return nil
	}

	tree := make([]*big.Int, 2*len(leaves)-1)
	for i, leaf := range leaves {
		tree[len(tree)-1-i] = leaf.Hash()
	}

	for i := len(tree) - len(leaves) - 1; i >= 0; i-- {
		leftChildIndex := getLeftChildIndex(i)
		rightChildIndex := getRightChildIndex(i)

		leftChild := new(fp.Element).SetBigInt(tree[leftChildIndex])
		rightChild := new(fp.Element).SetBigInt(tree[rightChildIndex])

		if leftChild.Cmp(rightChild) > 0 {
			fp := pedersenhash.Pedersen(leftChild, rightChild)
			tree[i] = fp.BigInt(new(big.Int))
		} else {
			fp := pedersenhash.Pedersen(rightChild, leftChild)
			tree[i] = fp.BigInt(new(big.Int))
		}

	}

	return tree
}

func GetMerkleProof(tree []*big.Int, index int) ([]*big.Int, error) {
	treeIndex := len(tree) - 1 - index

	var proof []*big.Int

	for treeIndex > 0 {
		siblingIndex, err := getSiblingIndex(treeIndex)
		if err != nil {
			return nil, err
		}
		proof = append(proof, tree[siblingIndex])
		treeIndex, err = getParentIndex(treeIndex)
		if err != nil {
			return nil, err
		}
	}

	return proof, nil
}
