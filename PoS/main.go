package main

import (
	"github.com/Poseidon/Block"
	"sync"
	"log"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"
	"github.com/Poseidon/Blockchain"
	"os"
	"net"
	"time"
	"math/rand"
	"io"
	"bufio"
	"strconv"
	"fmt"
	"encoding/hex"
	"crypto/sha256"
	"encoding/json"
)

// 处理那些进来的区块进行验证
var candidateBlocks = make(chan Block.Block)

// 向所有节点广播胜利验证器
var announcements = make(chan string)

var mutex = &sync.Mutex{}

// 跟踪打开的确认器和余额
var validators = make(map[string]int)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	genesisBlock := Block.Block{}
	genesisBlock = Block.Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", 1,"",""}
	spew.Dump(genesisBlock)
	Blockchain.Blockchain = append(Blockchain.Blockchain,genesisBlock)
	httpPort := os.Getenv("PORT")

	// start TCP and serve TCP server
	server,err := net.Listen("tcp",":"+httpPort)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port :", httpPort)
	defer server.Close()
	go func(){
		for candidate := range candidateBlocks{
			mutex.Lock()
			Blockchain.TempBlocks = append(Blockchain.TempBlocks,candidate)
			mutex.Unlock()
		}
	}()
	go func(){
		for {
			pickWinner()
		}
	}()

	for {
		conn,err := server.Accept()
		if err != nil{
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn){
	defer conn.Close()
	go func(){
		msg := <- announcements
		io.WriteString(conn,msg)
	}()
	var address string
	// 让用户输入抵押的token数量，给出的token越多，获得出新区块的机会就越大
	io.WriteString(conn,"Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan(){
		balance , err := strconv.Atoi(scanBalance.Text())
		if err!=nil{
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}
	io.WriteString(conn,"\nEnter a new Result:")
	scanResult := bufio.NewScanner(conn)
	go func(){
		for{
			for scanResult.Scan(){
				_result,err := strconv.Atoi(scanResult.Text())
				if err!=nil{
					log.Printf("%v not a number: %v", scanResult.Text(), err)
					delete(validators,address)
					conn.Close()
				}

				mutex.Lock()
				oldLastIndex := Blockchain.Blockchain[len(Blockchain.Blockchain)-1]
				mutex.Unlock()

				newBlock,err := generateBlock(oldLastIndex,_result,address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock,oldLastIndex){
					candidateBlocks <- newBlock
				}
				io.WriteString(conn,"\nEnter a new Result:")
			}
		}
	}()
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output,err := json.Marshal(Blockchain.Blockchain)
		mutex.Unlock()
		if err != nil{
			log.Fatal(err)
		}
		io.WriteString(conn,string(output)+"\n")
	}
}
func isBlockValid(newBlock, oldBlock Block.Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func calculateBlockHash(block Block.Block) string{
	record := string(block.Index)+block.TimeStamp+string(block.Result)+block.PrevHash
	return calculateHash(record)
}

func generateBlock(oldBlock Block.Block,Result int,address string)(Block.Block,error){
	var newBlock Block.Block
	t := time.Now()
	newBlock.Index = oldBlock.Index+1
	newBlock.TimeStamp = t.String()
	newBlock.Result = Result
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address
	return newBlock,nil
}

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 从待候选的出块者中选出胜利者出块获得奖励
func pickWinner(){
	time.Sleep(15 * time.Second)
	mutex.Lock()
	temp := Blockchain.TempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp)>0{
		OUTER:
			for _,block := range temp{
				// 如果临时区块的验证者已经在验证者乐透池中，就跳过
				for _,node := range lotteryPool{
					if block.Validator==node{
						continue OUTER
					}
				}

				mutex.Lock()
				setValidators := validators
				mutex.Unlock()
				// 将区块的验证者放到setValidators这样一个队列，要考虑到Token押注大小k，放入k份区块的验证者到乐透池中
				k,ok:=setValidators[block.Validator]
				if ok{
					for i:=0;i<k;i++{
						lotteryPool = append(lotteryPool,block.Validator)
					}
				}
			}
			// 从验证者乐透池中随机选出胜利者
			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

			// 新增区块到区块链上，胜利者进行广播让其他节点知道
			for _,block := range temp{
				if block.Validator == lotteryWinner{
					mutex.Lock()
					Blockchain.Blockchain = append(Blockchain.Blockchain,block)
					mutex.Unlock()
					for _ = range validators{
						announcements <- "\nwinning validator: " + lotteryWinner + "\n"
					}
					break
				}
			}
		// 将临时区块列表清空
		mutex.Lock()
		Blockchain.TempBlocks = []Block.Block{}
		mutex.Unlock()
	}
}