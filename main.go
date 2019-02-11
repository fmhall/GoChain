package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
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

func main() {
	fmt.Println("Placeholder...")
}
