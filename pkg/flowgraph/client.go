package flowgraph

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/samber/lo"
)

func NewFlowgraphClientFromEnv() graphql.Client {
	key := os.Getenv("FLOWGRAPH_KEY")
	if key == "" {
		err := fmt.Errorf("must set FLOWGRAPH_KEY")
		log.Fatal(err)
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     key,
			wrapped: http.DefaultTransport,
		},
	}

	return graphql.NewClient("https://query.flowgraph.co", &httpClient)

}

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

func GetEvents(ctx context.Context, hash string) ([]Event, error) {

	client := NewFlowgraphClientFromEnv()

	resp, err := TransactionEventsFirst(ctx, client, hash)
	if err != nil {
		return nil, err
	}
	root := resp.Transaction.Events
	nodes := lo.Map(root.Edges, TransactionEventsTransactionEventsEventConnectionEdgesEventEdge.Convert)

	if !root.PageInfo.HasNextPage {
		return nodes, nil
	}

	after := root.PageInfo.EndCursor
	for {
		resp, err = TransactionEvents(ctx, client, hash, after)
		if err != nil {
			return nil, err
		}

		root := resp.Transaction.Events
		eventNodes := lo.Map(root.Edges, TransactionEventsTransactionEventsEventConnectionEdgesEventEdge.Convert)
		nodes = append(nodes, eventNodes...)

		if !root.PageInfo.HasNextPage {
			return nodes, nil
		}

		after = root.PageInfo.EndCursor
	}
}

func GetAccountTransfers(ctx context.Context, accountId string, since time.Time) ([]Entry, error) {
	client := NewFlowgraphClientFromEnv()

	resp, err := AccountTransfersSince(ctx, client, accountId, since)
	if err != nil {
		return nil, err
	}
	root := resp.Account.TransferTransactions
	nodes := lo.Map(root.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdge.Convert)

	if !root.PageInfo.HasNextPage {
		return nodes, nil
	}

	after := root.PageInfo.EndCursor
	for {
		fmt.Print(".")
		resp, err = AccountTransfersAfter(ctx, client, accountId, after)
		if err != nil {
			return nil, err
		}
		root = resp.Account.TransferTransactions
		eventNodes := lo.Map(root.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdge.Convert)
		nodes = append(nodes, eventNodes...)

		if !root.PageInfo.HasNextPage {
			return nodes, nil
		}

		after = root.PageInfo.EndCursor
	}
}
