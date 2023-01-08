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

	accountId := "0x886f3aeaf848c535"
	packs := &koinly.Packs{Mappings: map[string]koinly.PackDetail{}}
	fileName := fmt.Sprintf("%s.json", accountId)
	outputFile := fmt.Sprintf("%s.csv", accountId)
	nftIds, err := koinly.ReadNFTMapping(fileName)
	if err != nil {
		nftIds = &koinly.NFTId{
			Next:     1,
			Mappings: map[string]int{},
		}
	}

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

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("all-tx.json", bytes, 0644)
	if err != nil {
		panic(err)
	}

	/*
		result := []flowgraph.Entry{}
		dropFile, err := os.ReadFile("all-tx.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(dropFile, &result)
		if err != nil {
			panic(err)
		}
	*/

	entires := []koinly.Event{}
	for _, ev := range result {
		entry, err := koinly.Convert(accountId, ev, nftIds, packs)
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

	err = koinly.WriteNFTMapping(nftIds, fileName)
	if err != nil {
		panic(err)
	}

	packBytes, err := json.MarshalIndent(packs, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("packs.json", packBytes, 0644)
	if err != nil {
		panic(err)
	}

}
