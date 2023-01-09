package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/araddon/dateparse"
	"github.com/bjartek/flow-koinly-export/pkg/flowgraph"
	"github.com/bjartek/flow-koinly-export/pkg/koinly"
)

func main() {

	state := koinly.NewKoinlyState()
	accountId := "0x886f3aeaf848c535"

	outputFile := fmt.Sprintf("%s.csv", accountId)

	time := "2019-01-01T00:00:00.0Z"
	//accountId := "0xdf868d4de6d2e0ab"
	ctx := context.Background()
	t, err := dateparse.ParseAny(time)
	if err != nil {
		panic(err)
	}

	result, err := flowgraph.GetAccountTransfers(ctx, accountId, t)
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

	//Here we can Unmarshal state so that we do not have to fetch it down every time
	entires := []koinly.Event{}
	for _, ev := range state.RawEntries {
		entry, err := koinly.Convert(accountId, ev, state)
		if err != nil {
			panic(err)
		}
		entires = append(entires, entry...)
	}

	fmt.Println("=================")

	err = koinly.Marshal(entires, outputFile)
	if err != nil {
		panic(err)
	}

}
