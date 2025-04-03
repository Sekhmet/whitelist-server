package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

type JsonRpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Result any `json:"result"`
}

type GenerateMerkleTreeRequestParams struct {
	Network string   `json:"network"`
	Entries []string `json:"entries"`
}

type GetMerkleRootRequestParams struct {
	RequestId string `json:"requestId"`
}

type GetMerkleProofRequestParams struct {
	Root  string `json:"root"`
	Index int    `json:"index"`
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(response)
}

func writeResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	response := SuccessResponse{Result: result}
	json.NewEncoder(w).Encode(response)
}

func NewRpcHandler(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var req JsonRpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case "generateMerkleTree":
			requestId, err := uuid.NewRandom()
			if err != nil {
				writeError(w, err)
				return
			}

			var params GenerateMerkleTreeRequestParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				writeError(w, err)
				return
			}

			if len(params.Entries) == 0 {
				writeError(w, errors.New("entries cannot be empty"))
				return
			}

			if _, err = db.Exec("INSERT INTO merkletree_requests (id, network) VALUES ($1, $2)", requestId.String(), params.Network); err != nil {
				writeError(w, err)
				return
			}

			request := &Request{
				id:      requestId.String(),
				network: params.Network,
				entries: params.Entries,
			}

			go func() {
				err := ProcessRequest(request, db)
				if err != nil {
					log.Printf("Error processing request: %v", err)
					return
				}
			}()

			writeResult(w, requestId)
		case "getMerkleRoot":
			var params GetMerkleRootRequestParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				writeError(w, err)
				return
			}

			if params.RequestId == "" {
				writeError(w, errors.New("requestId cannot be empty"))
				return
			}

			var root sql.NullString
			if err := db.QueryRow("SELECT root FROM merkletree_requests WHERE id = $1", params.RequestId).Scan(&root); err != nil {
				if err == sql.ErrNoRows {
					writeError(w, errors.New("request not found"))
					return
				}

				writeError(w, err)
				return
			}

			if !root.Valid {
				writeResult(w, nil)
				return
			}

			writeResult(w, root.String)
		case "getMerkleProof":
			var params GetMerkleProofRequestParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				writeError(w, err)
				return
			}

			var encodedTree string
			if err := db.QueryRow("SELECT tree FROM merkletree_requests WHERE root = $1", params.Root).Scan(&encodedTree); err != nil {
				if err == sql.ErrNoRows {
					writeError(w, errors.New("request not found"))
					return
				}

				writeError(w, err)
				return
			}

			var tree MerkleTree
			if err := json.Unmarshal([]byte(encodedTree), &tree); err != nil {
				writeError(w, err)
				return
			}

			proof, err := tree.GetMerkleProof(params.Index)
			if err != nil {
				writeError(w, err)
				return
			}

			var strProof = make([]string, len(proof))
			for i, node := range proof {
				strProof[i] = fmt.Sprintf("0x%x", node)
			}

			writeResult(w, strProof)

		default:
			http.Error(w, "Method not found", http.StatusNotFound)
			return
		}
	})

	handler := cors.Default().Handler(mux)

	return handler
}
