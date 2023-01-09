package flowgraph

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/bjartek/flow-koinly-export/pkg/core"
)

func (self TransactionTokensTransactionTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge) Convert(_ int) core.TokenTransfer {

	return core.TokenTransfer{
		Type:         string(self.Node.Type),
		Amount:       ConvertString(self.Node.Amount.Value),
		Token:        self.Node.Amount.Token.Id,
		Counterparty: self.Node.Counterparty.Address,
	}
}

func TransactionTokensFirst(
	ctx context.Context,
	client graphql.Client,
	txId string,
) (*TransactionTokensResponse, error) {
	req := &graphql.Request{
		OpName: "TransactionTokens",
		Query: `
query TransactionTokens ($txId: ID!) {
	transaction(id: $txId) {
		tokenTransfers(first: 50) {
			pageInfo {
				hasNextPage
				endCursor
			}
			edges {
				node {
					type
					amount {
						token {
							id
						}
						value
					}
					counterparty {
						address
					}
				}
			}
		}
	}
}
`,
		Variables: &__TransactionTokensInput{
			TxId: txId,
		},
	}
	var err error

	var data TransactionTokensResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
