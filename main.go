package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/Sekhmet/whitelist-server/starknet"
)

type JsonRpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type GetMerkleRootRequestParams struct {
	Entries []string `json:"entries"`
}

func main() {
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var req JsonRpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case "getMerkleRoot":
			var params GetMerkleRootRequestParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				http.Error(w, "Invalid parameters", http.StatusBadRequest)
				return
			}

			var leaves []Leaf
			for _, entry := range params.Entries {
				exploded := strings.Split(entry, ":")
				if len(exploded) != 2 {
					http.Error(w, "Invalid entry format", http.StatusBadRequest)
					return
				}

				address := exploded[0]
				addressType := starknet.AddressTypeStarknet
				if len(address) == 42 {
					addressType = starknet.AddressTypeEthereum
				}

				votingPower, success := new(big.Int).SetString(exploded[1], 0)
				if !success {
					http.Error(w, "Invalid voting power", http.StatusBadRequest)
					return
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

			fmt.Fprintf(w, "0x%x", root)

		default:
			http.Error(w, "Method not found", http.StatusNotFound)
			return
		}
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
