package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

type MerkleTree struct {
	hashes      []*big.Int
	treeIndices map[int]int
}

type EncodedTree struct {
	Hashes      []string    `json:"hashes"`
	TreeIndices map[int]int `json:"tree_indices"`
}

type IndexedLeaf struct {
	index int
	hash  *big.Int
}

func NewMerkleTree(leaves []Leaf, nodeHash NodeHash, sortLeaves bool) *MerkleTree {
	if len(leaves) == 0 {
		return nil
	}

	indexedLeaves := make([]IndexedLeaf, len(leaves))
	for i, leaf := range leaves {
		indexedLeaves[i] = IndexedLeaf{
			index: i,
			hash:  leaf.Hash(),
		}
	}

	if sortLeaves {
		slices.SortFunc(indexedLeaves, func(a, b IndexedLeaf) int {
			return a.hash.Cmp(b.hash)
		})
	}

	hashes := make([]*big.Int, 2*len(leaves)-1)
	treeIndices := make(map[int]int, len(leaves))

	for i, indexedLeaf := range indexedLeaves {
		treeIndex := len(hashes) - 1 - i

		hashes[treeIndex] = indexedLeaf.hash
		treeIndices[indexedLeaf.index] = treeIndex
	}

	for i := len(hashes) - len(leaves) - 1; i >= 0; i-- {
		leftChildIndex := getLeftChildIndex(i)
		rightChildIndex := getRightChildIndex(i)

		leftChild := hashes[leftChildIndex]
		rightChild := hashes[rightChildIndex]

		hashes[i] = nodeHash(leftChild, rightChild)
	}

	return &MerkleTree{
		hashes:      hashes,
		treeIndices: treeIndices,
	}
}

func (m *MerkleTree) Root() *big.Int {
	return m.hashes[0]
}

func (m *MerkleTree) GetMerkleProof(index int) ([]*big.Int, error) {
	treeIndex := m.treeIndices[index]

	var proof []*big.Int
	for treeIndex > 0 {
		siblingIndex, err := getSiblingIndex(treeIndex)
		if err != nil {
			return nil, err
		}
		proof = append(proof, m.hashes[siblingIndex])
		treeIndex, err = getParentIndex(treeIndex)
		if err != nil {
			return nil, err
		}
	}

	return proof, nil
}

func (m *MerkleTree) MarshalJSON() ([]byte, error) {
	var encodedHashes = make([]string, len(m.hashes))
	for i, node := range m.hashes {
		encodedHashes[i] = fmt.Sprintf("0x%x", node)
	}

	return json.Marshal(EncodedTree{
		Hashes:      encodedHashes,
		TreeIndices: m.treeIndices,
	})
}

func (m *MerkleTree) UnmarshalJSON(data []byte) error {
	var encodedTree EncodedTree

	if string(data)[0] == '[' {
		// Handle legacy format that does not include tree indices
		// This is only used for Starknet merkle trees which are not sorted
		// So we can compute treeIndices

		var hashes []string
		if err := json.Unmarshal(data, &hashes); err != nil {
			return err
		}

		numberOfLeaves := (len(hashes) + 1) / 2

		treeIndices := make(map[int]int, numberOfLeaves)
		for i := range numberOfLeaves {
			treeIndices[i] = len(hashes) - 1 - i
		}

		encodedTree.Hashes = hashes
		encodedTree.TreeIndices = treeIndices
	} else {
		if err := json.Unmarshal(data, &encodedTree); err != nil {
			return err
		}
	}

	hashes := make([]*big.Int, len(encodedTree.Hashes))
	for i, node := range encodedTree.Hashes {
		var success bool
		hashes[i], success = new(big.Int).SetString(node, 0)
		if !success {
			return errors.New("invalid tree format")
		}
	}

	m.hashes = hashes
	m.treeIndices = encodedTree.TreeIndices

	return nil
}
