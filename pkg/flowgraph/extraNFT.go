package flowgraph

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/bjartek/flow-koinly-export/pkg/core"
)

func (self TransactionNFTTransactionNftTransfersNFTTransferConnectionEdgesNFTTransferEdge) Convert(_ int) core.NFTTransfer {

	return core.NFTTransfer{
		From:     self.Node.From.Address,
		To:       self.Node.To.Address,
		Contract: self.Node.Nft.Contract.Id,
		Id:       self.Node.Nft.NftId,
	}
}

func TransactionNFTFirst(
	ctx context.Context,
	client graphql.Client,
	txId string,
) (*TransactionNFTResponse, error) {
	req := &graphql.Request{
		OpName: "TransactionNFT",
		Query: `
query TransactionNFT ($txId: ID!){
	transaction(id: $txId) {
		nftTransfers(first: 50) {
			pageInfo {
				hasNextPage
				endCursor
			}
			edges {
				node {
					nft {
						nftId
						contract {
							id
						}
					}
					from {
						address
					}
					to {
						address
					}
				}
			}
		}
	}
}
`,
		Variables: &__TransactionNFTInput{
			TxId: txId,
		},
	}
	var err error

	var data TransactionNFTResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
