package routes

import (
	"encoding/json"
	"net/http"

	"github.com/minhajuddinkhan/blockchainv2/respond"

	"github.com/davecgh/go-spew/spew"
	"github.com/minhajuddinkhan/blockchainv2"
)

//CreateBlock CreateBlock
func CreateBlock(BlockChain *blockchainv2.Blockchain) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var m blockchainv2.Message

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			respond.WithJSON(w, r, http.StatusBadRequest, r.Body)
			return
		}
		defer r.Body.Close()

		//ensure atomicity when creating new block
		BlockChain.Mutex.Lock()
		newBlock, _ := BlockChain.GenerateBlock(BlockChain.Blocks[len(BlockChain.Blocks)-1], m.BPM)
		BlockChain.Mutex.Unlock()

		if BlockChain.IsBlockValid(newBlock, BlockChain.Blocks[len(BlockChain.Blocks)-1]) {
			BlockChain.Blocks = append(BlockChain.Blocks, newBlock)
			spew.Dump(BlockChain.Blocks)
		}

		respond.WithJSON(w, r, http.StatusCreated, newBlock)

	}
}
