package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bjartek/flow-koinly-export/pkg/core"
	"github.com/bjartek/flow-koinly-export/pkg/koinly"
	"github.com/sanity-io/litter"
)

//TODO: found bug where I am not getting all nft/token transfers... this bloody api...
func main() {

	accountId := "0x886f3aeaf848c535"

	outputFile := fmt.Sprintf("%s.csv", accountId)

	/*
		//accountId := "0xdf868d4de6d2e0ab"
		ctx := context.Background()

		state := koinly.NewKoinlyState()
		result, err := flowgraph.GetAccountTransfers(ctx, accountId)
		if err != nil {
			panic(err)
		}

		state.RawEntries = result

		bytes, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("state.json", bytes, 0644)
		if err != nil {
			panic(err)
		}
	*/

	bytes, err := os.ReadFile("state.json")
	if err != nil {
		panic(err)
	}
	var state core.State
	err = json.Unmarshal(bytes, &state)
	if err != nil {
		panic(err)
	}

	//Here we can Unmarshal state so that we do not have to fetch it down every time
	entires := []koinly.Event{}
	for _, ev := range state.RawEntries {
		entry, err := koinly.Convert(accountId, ev, &state)
		if err != nil {
			litter.Dump(ev)
			panic(err)
		}
		entires = append(entires, entry...)
	}

	bytes, err = json.MarshalIndent(state, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("state-done.json", bytes, 0644)
	if err != nil {
		panic(err)
	}

	err = koinly.Marshal(entires, outputFile)
	if err != nil {
		panic(err)
	}

}
