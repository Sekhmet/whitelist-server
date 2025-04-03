package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Sekhmet/whitelist-server/evm"
	"github.com/Sekhmet/whitelist-server/starknet"
)

type Request struct {
	id      string
	network string
	entries []string
}

func GetStarknetTree(r *Request) (*MerkleTree, error) {
	var leaves []Leaf
	for _, entry := range r.entries {
		exploded := strings.Split(entry, ":")
		if len(exploded) != 2 {
			return nil, errors.New("invalid payload format")
		}

		address := exploded[0]
		addressType := starknet.AddressTypeStarknet
		if len(address) == 42 {
			addressType = starknet.AddressTypeEthereum
		}

		votingPower, success := new(big.Int).SetString(exploded[1], 0)
		if !success {
			return nil, errors.New("invalid voting power")
		}

		leaf := &starknet.Leaf{
			AddressType: addressType,
			Address:     address,
			VotingPower: *votingPower,
		}
		leaves = append(leaves, leaf)
	}

	return NewMerkleTree(leaves, starknet.NodeHash, false), nil
}

func GetEvmTree(r *Request) (*MerkleTree, error) {
	var leaves []Leaf
	for _, entry := range r.entries {
		exploded := strings.Split(entry, ":")
		if len(exploded) != 2 {
			return nil, errors.New("invalid payload format")
		}

		address := exploded[0]
		votingPower, success := new(big.Int).SetString(exploded[1], 0)
		if !success {
			return nil, errors.New("invalid voting power")
		}

		leaf := &evm.Leaf{
			Address:     address,
			VotingPower: *votingPower,
		}
		leaves = append(leaves, leaf)
	}

	return NewMerkleTree(leaves, evm.NodeHash, true), nil
}

func ProcessRequest(r *Request, db *sql.DB) error {
	log.Printf("Processing request: %v", r.id)

	var tree *MerkleTree
	var err error

	switch r.network {
	case "starknet":
		tree, err = GetStarknetTree(r)
	case "evm":
		tree, err = GetEvmTree(r)
	default:
		return errors.New("unsupported network")
	}

	if err != nil {
		return err
	}

	root := tree.Root()

	encodedTree, err := json.Marshal(tree)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE merkletree_requests SET processed = true, updated_at = CURRENT_TIMESTAMP, root = $1, tree = $2 WHERE id = $3",
		fmt.Sprintf("0x%x", root), encodedTree, r.id,
	)
	if err != nil {
		return err
	}

	return nil
}
