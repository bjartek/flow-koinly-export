package koinly

import (
	"github.com/bjartek/flow-koinly-export/pkg/core"
)

func NewKoinlyState() *core.State {
	return &core.State{
		Packs:           &core.Packs{Mappings: map[string]core.PackDetail{}},
		CompositeStatus: &core.CompositeStatus{NFTS: map[string]map[string]bool{}},
		NFTMappings:     &core.NFTId{Next: 1, Mappings: map[string]int{}},
		ManualPrices:    map[string]core.Price{},
	}
}
