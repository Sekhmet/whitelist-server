package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

func NewRpcMux() *http.ServeMux {
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

			request := &Request{
				id:      requestId.String(),
				network: params.Network,
				entries: params.Entries,
			}

			go ProcessRequest(request)

			writeResult(w, requestId)

		default:
			http.Error(w, "Method not found", http.StatusNotFound)
			return
		}
	})

	return mux
}
