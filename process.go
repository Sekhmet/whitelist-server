package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Sekhmet/whitelist-server/starknet"
)

type Request struct {
	id      string
	network string
	entries []string
}

func ProcessRequest(r *Request, db *sql.DB) error {
	log.Printf("Processing request: %v", r.id)

	var leaves []Leaf
	for _, entry := range r.entries {
		exploded := strings.Split(entry, ":")
		if len(exploded) != 2 {
			return errors.New("invalid payload format")
		}

		address := exploded[0]
		addressType := starknet.AddressTypeStarknet
		if len(address) == 42 {
			addressType = starknet.AddressTypeEthereum
		}

		votingPower, success := new(big.Int).SetString(exploded[1], 0)
		if !success {
			return errors.New("invalid voting power")
		}

		leaf := &starknet.Leaf{
			AddressType: addressType,
			Address:     address,
			VotingPower: *votingPower,
		}
		leaves = append(leaves, leaf)
	}

	tree := NewMerkleTree(leaves, starknet.NodeHash, false)
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
