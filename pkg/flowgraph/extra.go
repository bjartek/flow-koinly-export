package flowgraph

//go:generate go run github.com/Khan/genqlient
import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Khan/genqlient/graphql"
	"github.com/bjartek/flow-koinly-export/pkg/core"
	"github.com/samber/lo"
	"github.com/xeipuuv/gojsonpointer"
)

func AccountTransfersFirst(
	ctx context.Context,
	client graphql.Client,
	accountId string,
) (*AccountTransfersResponse, error) {
	req := &graphql.Request{
		OpName: "AccountTransfers",
		Query: `
query AccountTransfers ($accountId: ID!) {
	account(id: $accountId) {
		transferTransactions(first: 50, ordering: Ascending) {
			pageInfo {
				hasNextPage
				endCursor
			}
			edges {
				transaction {
					hash
					time
					script
					arguments
					events(first: 50) {
						pageInfo {
							hasNextPage
						}
						edges {
							node {
								fields
								type {
									id
									fields {
										identifier
									}
								}
							}
						}
					}
				}
				nftTransfers {
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
				tokenTransfers {
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
	}
}
`,
		Variables: &__AccountTransfersInput{
			AccountId: accountId,
		},
	}
	var err error

	var data AccountTransfersResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransaction) Convert() core.Transaction {

	events := lo.Map(self.Events.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdge.Convert)

	//TODO: if this is pay rewards do not do this :P
	//TODO: maybe even just here bake in that we only want certain events?
	if self.Events.PageInfo.HasNextPage {
		ev, err := GetEvents(context.Background(), self.Hash)
		if err != nil {
			panic(err)
		}
		events = ev
	}

	res, err := json.Marshal(self.Arguments)
	if err != nil {
		panic(err)
	}
	var args []RawArgument
	err = json.Unmarshal(res, &args)
	if err != nil {
		panic(err)
	}

	return core.Transaction{
		Hash:       self.Hash,
		Time:       self.Time,
		Script:     self.Script,
		ScriptHash: fmt.Sprintf("%x", sha256.Sum256([]byte(self.Script))),
		Events:     events,
		Arguments:  lo.Map(args, RawArgument.GetValue),
	}

}
func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdgeNodeEventTypeFieldsEventTypeField) Convert(_ int) string {

	return self.Identifier
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdge) Convert(_ int) core.Event {
	keys := lo.Map(self.Node.Type.Fields, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdgeNodeEventTypeFieldsEventTypeField.Convert)
	values := self.Node.Fields

	subValuePointer, _ := gojsonpointer.NewJsonPointer("/value/value")
	valuePointer, _ := gojsonpointer.NewJsonPointer("/value")

	valuesLength := len(values)
	fields := map[string]interface{}{}
	for i, key := range keys {

		if i >= valuesLength {
			continue
		}
		value, _, err := valuePointer.Get(values[i])
		if err != nil {
			panic(err)
		}

		if value != nil {
			fields[key] = value
		} else {
			optionalValue, _, err := subValuePointer.Get(values[i])
			if err != nil {
				fields[key] = nil
			} else {
				fields[key] = optionalValue
			}
		}
	}

	return core.Event{
		Name:   self.Node.Type.Id,
		Fields: fields,
	}
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeNftTransfersNFTTransferConnectionEdgesNFTTransferEdge) Convert(_ int) core.NFTTransfer {
	return core.NFTTransfer{
		From:     self.Node.From.Address,
		To:       self.Node.To.Address,
		Contract: self.Node.Nft.Contract.Id,
		Id:       self.Node.Nft.NftId,
	}
}

//Use a bindings marshaller here
func ConvertString(value string) float64 {

	number, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return float64(number) / 100_000_000

}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge) Convert(_ int) core.TokenTransfer {
	//500_000_000
	//5.00000000

	return core.TokenTransfer{
		Type:         string(self.Node.Type),
		Amount:       ConvertString(self.Node.Amount.Value),
		Token:        self.Node.Amount.Token.Id,
		Counterparty: self.Node.Counterparty.Address,
	}
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdge) Convert(_ int) core.Entry {
	tx := self.Transaction.Convert()
	nfts := lo.Map(self.NftTransfers.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeNftTransfersNFTTransferConnectionEdgesNFTTransferEdge.Convert)
	transfer := lo.Map(self.TokenTransfers.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge.Convert)

	var err error
	//if we have more then 10 transfers we have to refetch everyone from here, we cannot use the first and then append since it is not sorted the same!
	if self.NftTransfers.PageInfo.HasNextPage {
		nfts, err = GetTxNFTS(context.Background(), tx.Hash)
		if err != nil {
			panic(err)
		}
	}

	if self.TokenTransfers.PageInfo.HasNextPage {
		transfer, err = GetTxTransfers(context.Background(), tx.Hash)
		if err != nil {
			panic(err)
		}
	}

	return core.Entry{
		Transaction: tx,
		NFT:         nfts,
		Tokens:      transfer,
	}
}
