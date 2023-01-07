package main

import (
	"context"
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/bjartek/flow-koinly-export/pkg/flowgraph"
	"github.com/bjartek/flow-koinly-export/pkg/koinly"
)

func main() {

	ctx := context.Background()

	time := "2022-01-13T00:00:00.0Z"
	accountId := "0x886f3aeaf848c535"
	t, err := dateparse.ParseAny(time)
	if err != nil {
		panic(err)
	}

	result, err := flowgraph.GetAccountTransfers(ctx, accountId, t)
	if err != nil {
		panic(err)
	}

	entires := []koinly.Event{}
	for _, ev := range result {
		entry, err := koinly.Convert(accountId, ev)
		if err != nil {
			panic(err)
		}
		entires = append(entires, entry...)
	}

	fmt.Println("=================")

	err = koinly.Marshal(entires)
	if err != nil {
		panic(err)
	}
}
