package flowgraph

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/bjartek/flow-koinly-export/pkg/core"
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

func GetTxNFTS(ctx context.Context, txId string) ([]core.NFTTransfer, error) {

	client := NewFlowgraphClientFromEnv()

	resp, err := TransactionNFTFirst(ctx, client, txId)
	if err != nil {
		return nil, err
	}
	root := resp.Transaction.NftTransfers
	nodes := lo.Map(root.Edges, TransactionNFTTransactionNftTransfersNFTTransferConnectionEdgesNFTTransferEdge.Convert)

	if !root.PageInfo.HasNextPage {
		return nodes, nil
	}

	after := root.PageInfo.EndCursor
	for {
		resp, err := TransactionNFT(ctx, client, txId, after)
		if err != nil {
			return nil, err
		}

		root := resp.Transaction.NftTransfers
		eventNodes := lo.Map(root.Edges, TransactionNFTTransactionNftTransfersNFTTransferConnectionEdgesNFTTransferEdge.Convert)
		nodes = append(nodes, eventNodes...)

		if !root.PageInfo.HasNextPage {
			return nodes, nil
		}

		after = root.PageInfo.EndCursor
	}
}

func GetTxTransfers(ctx context.Context, txId string) ([]core.TokenTransfer, error) {

	client := NewFlowgraphClientFromEnv()

	resp, err := TransactionTokensFirst(ctx, client, txId)
	if err != nil {
		return nil, err
	}
	root := resp.Transaction.TokenTransfers
	nodes := lo.Map(root.Edges, TransactionTokensTransactionTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge.Convert)

	if !root.PageInfo.HasNextPage {
		return nodes, nil
	}

	after := root.PageInfo.EndCursor
	for {
		resp, err := TransactionTokens(ctx, client, txId, after)
		if err != nil {
			return nil, err
		}

		root := resp.Transaction.TokenTransfers
		eventNodes := lo.Map(root.Edges, TransactionTokensTransactionTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge.Convert)
		nodes = append(nodes, eventNodes...)

		if !root.PageInfo.HasNextPage {
			return nodes, nil
		}

		after = root.PageInfo.EndCursor
	}
}

func GetEvents(ctx context.Context, hash string) ([]core.Event, error) {

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

func GetAccountTransfers(ctx context.Context, accountId string) ([]core.Entry, error) {
	client := NewFlowgraphClientFromEnv()

	resp, err := AccountTransfersFirst(ctx, client, accountId)
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
		resp, err = AccountTransfers(ctx, client, accountId, after)
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
