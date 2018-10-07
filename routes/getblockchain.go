package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/minhajuddinkhan/blockchainv2"
)

//GetBlockChain GetBlockChain
func GetBlockChain(blockChain *blockchainv2.Blockchain) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.MarshalIndent(blockChain.Blocks, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(bytes))
	}
}
