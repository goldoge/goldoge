package gd

import (
	"errors"
	"github.com/bitonicnl/verify-signed-message/pkg"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

type BlockChain struct {
	Chain  []Block
	Height int64
}

func (bc *BlockChain) SubmitStar(address string, message string, signature string, star string) (*Block, error) {
	submitedTime, _ := strconv.ParseInt(strings.Split(message, ":")[1], 10, 64)
	submitedTime = submitedTime * 1000
	currentTime := time.Now().UnixMilli()
	result, err := verifier.Verify(
		verifier.SignedMessage{
			Address:   address,
			Message:   message,
			Signature: signature,
		})
	if (err == nil) && result && (currentTime-submitedTime < 300000) {
		newBlock := bc.AddBlock(Block{
			Payload: Payload{
				Owner: address,
				Star:  star,
			},
		})
		return &newBlock, nil
	}
	return nil, errors.New("empty")
}

func (bc *BlockChain) GetBlockByHeight(height int64) (*Block, error) {
	idx := slices.IndexFunc(bc.Chain, func(bl Block) bool { return bl.Height == height })
	if idx == -1 {
		return nil, errors.New("empty")
	}
	return &bc.Chain[idx], nil
}

func (bc *BlockChain) GetBlockByHash(hash string) (Block, error) {
	idx := slices.IndexFunc(bc.Chain, func(bl Block) bool { return bl.Hash == hash })
	if idx == -1 {
		return Block{}, errors.New("empty")
	}
	return bc.Chain[idx], nil
}

func (bc *BlockChain) GetStarByWalletAddress(address string) []string {
	var stars []string
	for _, block := range bc.Chain {
		if block.Payload.Owner == address {
			stars = append(stars, block.Payload.Star)
		}
	}
	return stars
}

func (bc *BlockChain) InitializeChain() {
	if bc.Height == -1 {
		bl := Block{
			Hash:              "",
			Height:            0,
			Payload:           Payload{Owner: "0", Star: "Genesis Block"},
			Time:              0,
			PreviousBlockHash: "",
		}
		bc.AddBlock(bl)
	}
}

func (bc *BlockChain) AddBlock(bl Block) Block {
	currentHeight := int64(len(bc.Chain))
	bl.Height = currentHeight
	bl.Time = time.Now().UnixMilli()
	if currentHeight > 0 {
		bl.PreviousBlockHash = bc.Chain[currentHeight-1].Hash
	}
	bl.Hash = RecalculateHash(bl)
	bc.Chain = append(bc.Chain, bl)
	bc.Height++
	return bl
}

func (bc *BlockChain) RequestMessageOwnershipVerification(address string) string {
	return address + ":" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ":startRegistry"
}
