package Block

import (
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"github.com/Poseidon/Config"
	"fmt"
	"strings"
)

type Block struct {
	Index	int
	TimeStamp	string
	Result	int
	Hash	string
	PrevHash	string
	Difficulty	int
	Nonce	string
	Validator string
}

// SHA256 hasing
func CalculateHash(block Block) string{
	record := strconv.Itoa(block.Index)+block.TimeStamp+strconv.Itoa(block.Result)+block.PrevHash+block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateGenesisBlock() Block{
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{Index:0,TimeStamp:t.String(),Result:0,Hash:CalculateHash(genesisBlock),PrevHash:"",Difficulty: Config.Difficulty,Nonce:""}
	return genesisBlock
}

func GenerateBlock(oldBlock Block, Result int) Block {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index+1
	newBlock.TimeStamp = t.String()
	newBlock.Result = Result
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = Config.Difficulty
	for i:=0; ;i++{
		hex := fmt.Sprintf("%x",i)
		newBlock.Nonce = hex
		if !isHashValid(CalculateHash(newBlock),newBlock.Difficulty){
			fmt.Println(CalculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		}else{
			fmt.Println(CalculateHash(newBlock), " work done!")
			newBlock.Hash = CalculateHash(newBlock)
			break
		}
	}
	return newBlock
}

func isHashValid(hash string,difficulty int) bool{
	prefix := strings.Repeat("0",difficulty)
	return strings.HasPrefix(hash,prefix)
}

func IsBlockValid(newBlock,oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index{
		return false
	}
	if oldBlock.Hash!=newBlock.PrevHash{
		return false
	}
	if CalculateHash(newBlock)!=newBlock.Hash{
		return false
	}
	return true
}