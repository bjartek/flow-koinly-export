package koinly

import (
	"encoding/json"
	"fmt"
	"os"
)

type Packs struct {
	Mappings map[string]PackDetail
}

type PackDetail struct {
	Amount   float64
	Currency string
}

func (self *Packs) Add(key string, amount float64, currency string) {
	self.Mappings[key] = PackDetail{
		Amount:   amount,
		Currency: currency,
	}
}

type NFTId struct {
	Next     int
	Mappings map[string]int
}

func (self *NFTId) GetOrAdd(typ string, id string) string {
	if id == "" {
		panic("foo")
	}
	key := fmt.Sprintf("%s-%s", typ, id)
	value, ok := self.Mappings[key]
	if ok {
		return GenerateTextId(value)
	}
	value = self.Next
	self.Mappings[key] = value
	self.Next = self.Next + 1
	return GenerateTextId(value)
}

func GenerateTextId(value int) string {
	if value > 5000 {
		value = value - 5000
		return fmt.Sprintf("NULL%d", value)
	}
	return fmt.Sprintf("NFT%d", value)

}

func ReadNFTMapping(fileName string) (*NFTId, error) {
	dropFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var drop NFTId
	err = json.Unmarshal(dropFile, &drop)
	if err != nil {
		return nil, err
	}
	return &drop, nil
}

func WriteNFTMapping(nft *NFTId, fileName string) error {

	bytes, err := json.MarshalIndent(nft, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, bytes, 0644)
	return err

}
