package main

import (
	"errors"
	"math"
	"math/big"
	"slices"
)

type NodeHash func(left, right *big.Int) *big.Int

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

func GenerateMerkleTree(leaves []Leaf, nodeHash NodeHash, sortLeaves bool) []*big.Int {
	if len(leaves) == 0 {
		return nil
	}

	sortedHashes := make([]*big.Int, len(leaves))
	for i, leaf := range leaves {
		sortedHashes[i] = leaf.Hash()
	}

	if sortLeaves {
		slices.SortFunc(sortedHashes, func(a, b *big.Int) int {
			return a.Cmp(b)
		})
	}

	tree := make([]*big.Int, 2*len(leaves)-1)
	for i, hash := range sortedHashes {
		tree[len(tree)-1-i] = hash
	}

	for i := len(tree) - len(leaves) - 1; i >= 0; i-- {
		leftChildIndex := getLeftChildIndex(i)
		rightChildIndex := getRightChildIndex(i)

		leftChild := tree[leftChildIndex]
		rightChild := tree[rightChildIndex]

		tree[i] = nodeHash(leftChild, rightChild)
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
