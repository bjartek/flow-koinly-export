package flowgraph

import (
	"fmt"
	"time"
)

type Transaction struct {
	Hash       string
	Time       time.Time
	Script     string
	ScriptHash string
	Events     []Event
	Arguments  []string
}

type RawArgument struct {
	Value interface{}
	Type  string
}

func (self RawArgument) GetValue(_ int) string {
	return fmt.Sprintf("%v", self.Value)

}

type Event struct {
	Name   string
	Fields map[string]interface{}
}

type Entry struct {
	Transaction Transaction
	NFT         []NFTTransfer
	Tokens      []TokenTransfer
}

type NFTTransfer struct {
	From     string
	To       string
	Contract string
	Id       string
}
