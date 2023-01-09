package koinly

import (
	"github.com/bjartek/flow-koinly-export/pkg/core"
)

//TODO: create an interface here
func NewKoinlyState() *core.State {
	return &core.State{
		Packs:           &core.Packs{Mappings: map[string]core.PackDetail{}},
		CompositeStatus: &core.CompositeStatus{NFTS: map[string]map[string]bool{}},
		NFTMappings:     &core.NFTId{Next: 1, Mappings: map[string]int{}},
	}
}
