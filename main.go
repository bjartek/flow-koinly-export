package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bjartek/flow-koinly-export/pkg/core"
	"github.com/bjartek/flow-koinly-export/pkg/flowgraph"
	"github.com/bjartek/flow-koinly-export/pkg/koinly"
	"github.com/sanity-io/litter"
)

// TODO: found bug where I am not getting all nft/token transfers... this bloody api...
func main() {

	//accountId := "0xdf868d4de6d2e0ab"
	//me
	accountId := "0x886f3aeaf848c535"
	//accountId := "0x8e1231b8b045cf96"
	//accountId := "0x5b64854c16a96267"
	//c3
	//accountId := "0x16ae8f1cbfceaa9e"

	stateFile := fmt.Sprintf("%s.json", accountId)
	outputFile := fmt.Sprintf("%s.csv", accountId)

	ctx := context.Background()

	state := koinly.NewKoinlyState()

	//we try to fetch the old state file so that we do not have to fetch all interactions from flowgraph again
	bytes, err := os.ReadFile(stateFile)
	if err == nil {
		var oldState core.State
		err = json.Unmarshal(bytes, &oldState)
		if err != nil {
			panic(err)
		}
		state.RawEntries = oldState.RawEntries
		state.ManualPrices = oldState.ManualPrices

	} else {
		result, err := flowgraph.GetAccountTransfers(ctx, accountId)
		if err != nil {
			panic(err)
		}
		state.RawEntries = result
	}

	//Here we can Unmarshal state so that we do not have to fetch it down every time
	entires := []koinly.Event{}
	for _, ev := range state.RawEntries {
		entry, err := koinly.Convert(accountId, ev, state)
		if err != nil {
			litter.Dump(ev)
			panic(err)
		}
		entires = append(entires, entry...)
	}

	bytes2, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s.json", accountId), bytes2, 0644)
	if err != nil {
		panic(err)
	}

	err = koinly.Marshal(entires, outputFile)
	if err != nil {
		panic(err)
	}

}
