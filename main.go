package main

import (
	"context"
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/bjartek/flow-koinly-export/pkg/flowgraph"
	"github.com/sanity-io/litter"
)

func main() {

	ctx := context.Background()

	time := "2022-01-13T00:00:00.0Z"
	accountId := "0x886f3aeaf848c535"
	t, err := dateparse.ParseAny(time)
	if err != nil {
		panic(err)
	}

	fmt.Println(t)
	result, err := flowgraph.GetAccountTransfers(ctx, accountId, t)
	if err != nil {
		panic(err)
	}

	litter.Dump(result)
}
