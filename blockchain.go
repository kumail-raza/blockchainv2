package blockchainv2

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Block Basic block of a blochain linked list
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Nonce     string
}

//Blockchain Blockchain
type Blockchain struct {
	sync.Mutex
	Blocks        []Block
	Difficulty    int
	Channel       chan Blockchain
	BroadcastTime time.Duration
}

//Message Message
type Message struct {
	BPM int
}

//CalculateHash CalculateHash
func (bc *Blockchain) CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//IsBlockValid IsBlockValid
func (bc *Blockchain) IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if bc.CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

//ReplaceChain ReplaceChain
func (bc *Blockchain) ReplaceChain(newBlocks []Block) {

	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	if len(newBlocks) > len(bc.Blocks) {
		bc.Blocks = newBlocks
	}

}

//GenerateBlock GenerateBlock
func (bc *Blockchain) GenerateBlock(oldBlock Block, BPM int) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		calculatedHash := bc.CalculateHash(newBlock)
		if !bc.isHashValid(calculatedHash) {
			fmt.Println("invalid hash calculated.", calculatedHash)
			bc.CalculateHash(newBlock)
			time.Sleep(time.Second)
			continue

		} else {
			fmt.Println(calculatedHash)
			newBlock.Hash = bc.CalculateHash(newBlock)
			break
		}

	}
	return newBlock, nil
}

func (bc *Blockchain) isHashValid(hash string) bool {
	prefix := strings.Repeat("0", bc.Difficulty)
	return strings.HasPrefix(hash, prefix)
}
