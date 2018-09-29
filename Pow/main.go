package main

import (
	"sync"
	"github.com/joho/godotenv"
	"log"
	"github.com/Poseidon/Block"
	"github.com/davecgh/go-spew/spew"
	"github.com/Poseidon/Blockchain"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io"
	"github.com/Poseidon/Message"
	"os"
	"time"
)

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func(){
		genesisBlock := Block.GenerateGenesisBlock()
		spew.Dump(genesisBlock)
		mutex.Lock()
		Blockchain.Blockchain = append(Blockchain.Blockchain,genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}

// web server
func run() error{
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
		ReadTimeout:10*time.Second,
		WriteTimeout:10*time.Second,
		MaxHeaderBytes:1<<20,
	}
	if err:=s.ListenAndServe();err!=nil{
		return err
	}
	return nil
}

// create handlers
func makeMuxRouter() http.Handler{
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/",handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/",handleWriteBlock).Methods("POST")
	return muxRouter
}

// 当我们收到一个http请求的时候，写入blockchain
func handleGetBlockchain(w http.ResponseWriter,r *http.Request){
	bytes,err := json.MarshalIndent(Blockchain.Blockchain,"","  ")
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	io.WriteString(w,string(bytes))
}

// 将JSON有效负载作为结果的一个输入
func handleWriteBlock(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var m Message.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m);err!=nil{
		respondWithJSON(w,r,http.StatusBadRequest,r.Body)
		return
	}
	defer r.Body.Close()
	mutex.Lock()
	newBlock := Block.GenerateBlock(Blockchain.Blockchain[len(Blockchain.Blockchain)-1],m.Result)
	mutex.Unlock()
	if Block.IsBlockValid(newBlock, Blockchain.Blockchain[len(Blockchain.Blockchain)-1]){
		Blockchain.Blockchain = append(Blockchain.Blockchain,newBlock)
		spew.Dump(Blockchain.Blockchain)
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request,code int,payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	response,err := json.MarshalIndent(payload,"","")
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Inernal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}


















