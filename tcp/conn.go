package tcp

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/minhajuddinkhan/blockchainv2"
)

//HandleTCPConn HandleTCPConn
func HandleTCPConn(conn net.Conn, blockChain *blockchainv2.Blockchain) {

	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	blockChainCh := make(chan blockchainv2.Blockchain)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := blockChain.GenerateBlock(blockChain.Blocks[len(blockChain.Blocks)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if blockChain.IsBlockValid(newBlock, blockChain.Blocks[len(blockChain.Blocks)-1]) {
				newBlockchain := append(blockChain.Blocks, newBlock)
				blockChain.ReplaceChain(newBlockchain)
			}

			blockChainCh <- *blockChain

			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(blockChain.BroadcastTime)
			blockChain.Mutex.Lock()
			output, err := json.Marshal(blockChain.Blocks)
			if err != nil {
				log.Fatal(err)
			}
			blockChain.Mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for range blockChainCh {
		spew.Dump(blockChain)
	}

}
