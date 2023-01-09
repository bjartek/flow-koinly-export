package core

import (
	"time"

	"github.com/samber/lo"
)

type Event struct {
	Name   string
	Fields map[string]interface{}
}

//TODO: if the user has paid a fee here extract it out and put it here
type Entry struct {
	Transaction Transaction
	NFT         []NFTTransfer
	Tokens      []TokenTransfer
}

type TokenTransfer struct {
	Type         string
	Amount       float64
	Token        string
	Counterparty string
}

type NFTTransfer struct {
	From     string
	To       string
	Contract string
	Id       string
}

type Transaction struct {
	Hash       string
	Time       time.Time
	Script     string
	ScriptHash string
	Events     []Event
	Arguments  []string
}

type State struct {
	Packs           *Packs
	CompositeStatus *CompositeStatus
	NFTMappings     NFTIDGenerator
	RawEntries      []Entry
}

func (self *State) HasNFTID(contract, id string) bool {
	return self.NFTMappings.Contains(contract, id)
}

func (self *State) GetNFTID(contract, id string) string {
	return self.NFTMappings.Get(contract, id)
}
func (self *State) AddNFTID(contract, id string) string {
	//We might have owned this NFT before so we try to get it first
	return self.NFTMappings.GetOrAdd(contract, id)
}

func (self *State) AddPack(key string, amount float64, currency string) {
	self.Packs.Add(key, amount, currency)
}

func (self *State) RemoveCompositeComponent(nft string, component string) bool {
	return self.CompositeStatus.RemoveComponent(nft, component)
}

func (self *State) GetPack(key string) (PackDetail, bool) {
	value, ok := self.Packs.Mappings[key]
	return value, ok
}

func (self *State) GetCompositeComponent(id string) []string {
	return self.CompositeStatus.Component(id)
}

func (self *State) AddCompositeComponent(nft string, component string) {
	self.CompositeStatus.AddComponent(nft, component)
}

type CompositeStatus struct {
	NFTS map[string]map[string]bool
}

func (self *CompositeStatus) AddComponent(nft string, component string) {
	main, ok := self.NFTS[nft]
	if !ok {
		main = map[string]bool{}
	}
	main[component] = true

	self.NFTS[nft] = main
}
func (self *CompositeStatus) RemoveComponent(nft string, component string) bool {
	main, ok := self.NFTS[nft]
	if !ok {
		return false
	}

	delete(main, component)
	self.NFTS[nft] = main
	return true
}

func (self *CompositeStatus) Component(nft string) []string {
	return lo.Keys(self.NFTS[nft])
}

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

//An interface to get
type NFTIDGenerator interface {
	GetOrAdd(typ string, id string) string
	Get(typ string, id string) string
	Contains(typ string, id string) bool
}