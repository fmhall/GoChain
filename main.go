package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
)

// Data Structure to represent each block
type Block struct {
	Index     int
	Timestamp string
	Data      int
	Hash      string
	PrevHash  string
}

// Blockchain is a slice of the Block struct
var Blockchain []Block

func getHash(block Block) string {
	aggregate := string(block.Index) + block.Timestamp + string(block.Data)
	h := sha256.New()
	h.Write([]byte(aggregate))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, Data int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = getHash(newBlock)

	return newBlock, nil

}

func isBlockValid(newBlock, prevBlock Block) bool {
	if newBlock.Index != prevBlock.Index+1 {
		return false
	}
	if newBlock.Timestamp < prevBlock.Timestamp {
		return false
	}
	if newBlock.PrevHash != prevBlock.Hash {
		return false
	}
	if newBlock.Hash != getHash(newBlock) {
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func run() error {
	r := mux.makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	Data int
}



func main() {
	fmt.Println("Placeholder...")
}
