package koinly

import (
	"fmt"

	"github.com/bjartek/flow-koinly-export/pkg/core"
)

func NewKoinlyState() *core.State {
	return &core.State{
		Packs:           &core.Packs{Mappings: map[string]core.PackDetail{}},
		CompositeStatus: &core.CompositeStatus{NFTS: map[string]map[string]bool{}},
		NFTMappings:     &NFTId{Next: 1, Mappings: map[string]int{}},
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
	self.Next = value + 1
	return GenerateTextId(value)
}

func (self *NFTId) Contains(typ string, id string) bool {
	key := fmt.Sprintf("%s-%s", typ, id)
	_, ok := self.Mappings[key]
	return ok
}

func (self *NFTId) Get(typ string, id string) string {
	key := fmt.Sprintf("%s-%s", typ, id)
	value, ok := self.Mappings[key]
	if ok {
		return GenerateTextId(value)
	}
	panic("Could not get existing NFTId")
}

func GenerateTextId(value int) string {
	if value > 5000 {
		value = value - 5000
		return fmt.Sprintf("NULL%d", value)
	}
	return fmt.Sprintf("NFT%d", value)

}
