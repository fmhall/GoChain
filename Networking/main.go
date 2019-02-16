package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	"bufio"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// Block Data Structure to represent each block
type Block struct {
	Index     int
	Timestamp string
	Data      int
	Hash      string
	PrevHash  string
}

// Blockchain is a slice of the Block struct
var Blockchain []Block

// Adding networking functionality - bcServer is initialized as a channel for the blockchain
var bcServer chan []Block

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

// func makeMuxRouter() http.Handler {
// 	muxRouter := mux.NewRouter()
// 	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
// 	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
// 	return muxRouter
// }

// func run() error {
// 	r := makeMuxRouter()
// 	httpAddr := os.Getenv("PORT")
// 	log.Println("Listening on ", os.Getenv("PORT"))
// 	s := &http.Server{
// 		Addr:           ":" + httpAddr,
// 		Handler:        r,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}

// 	if err := s.ListenAndServe(); err != nil {
// 		return err
// 	}

// 	return nil
// }

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// Message type to pass in data
type Message struct {
	Data int
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Enter new data:")

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			data, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}

			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], data)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter new data: ")
		}
	}()

	// Modelling broadcast retrieval
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}

}

// Uncomment for HTTP blockchain visualization without networking

// func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
// 	var m Message

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&m); err != nil {
// 		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
// 		return
// 	}
// 	defer r.Body.Close()

// 	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.Data)
// 	if err != nil {
// 		respondWithJSON(w, r, http.StatusInternalServerError, m)
// 		return
// 	}
// 	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
// 		newBlockchain := append(Blockchain, newBlock)
// 		replaceChain(newBlockchain)
// 		spew.Dump(Blockchain)
// 	}

// 	respondWithJSON(w, r, http.StatusCreated, newBlock)

// }

// Uncomment for HTTP blockchain visualization without networking

// func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
// 	response, err := json.MarshalIndent(payload, "", "  ")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("HTTP 500: Internal Server Error"))
// 		return
// 	}
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

	log.Fatal(run())
}
