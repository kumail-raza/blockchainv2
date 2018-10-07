package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/minhajuddinkhan/blockchainv2/tcp"

	"github.com/joho/godotenv"

	"github.com/minhajuddinkhan/blockchainv2"
	"github.com/minhajuddinkhan/blockchainv2/routes"

	"github.com/gorilla/mux"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	httpAddr := os.Getenv("HTTP_ADDR")
	tcpAddr := os.Getenv("TCP_ADDR")

	r := mux.NewRouter()
	genesisBlock := blockchainv2.Block{
		Index:     1,
		BPM:       1,
		Nonce:     "1",
		PrevHash:  "",
		Hash:      "",
		Timestamp: time.Now().String(),
	}
	blockChain := blockchainv2.Blockchain{
		Blocks:        []blockchainv2.Block{genesisBlock},
		Difficulty:    1,
		BroadcastTime: 10 * time.Second,
	}

	r.HandleFunc("/", routes.GetBlockChain(&blockChain)).Methods("GET")
	r.HandleFunc("/", routes.CreateBlock(&blockChain)).Methods("POST")

	server, err := net.Listen("tcp", ":"+tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on tcp connections", tcpAddr)
	defer server.Close()
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go tcp.HandleTCPConn(conn, &blockChain)
		}

	}()
	log.Println("Listening on http requests", httpAddr)

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
