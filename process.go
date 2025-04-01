package main

import (
	"errors"
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

func ProcessRequest(r *Request) error {
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

	tree := GenerateMerkleTree(leaves)
	root := tree[0]

	log.Printf("[%s] Merkle root 0x%x", r.id, root)

	return nil
}
