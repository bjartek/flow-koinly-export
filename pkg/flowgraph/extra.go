package flowgraph

//go:generate go run github.com/Khan/genqlient
import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/samber/lo"
	"github.com/xeipuuv/gojsonpointer"
)

func AccountTransfersSince(
	ctx context.Context,
	client graphql.Client,
	accountId string,
	since time.Time,
) (*AccountTransfersResponse, error) {
	req := &graphql.Request{
		OpName: "AccountTransfers",
		Query: `
query AccountTransfers ($accountId: ID!, $since: Time!) {
	account(id: $accountId) {
		transferTransactions(first: 50, ordering: Ascending, since: $since) {
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
			Since:     since,
			After:     "",
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

func AccountTransfersAfter(
	ctx context.Context,
	client graphql.Client,
	accountId string,
	after string,
) (*AccountTransfersResponse, error) {
	req := &graphql.Request{
		OpName: "AccountTransfers",
		Query: `
query AccountTransfers ($accountId: ID!, $after: ID) {
	account(id: $accountId) {
		transferTransactions(first: 50, ordering: Ascending, after: $after) {
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
			After:     after,
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

type Transaction struct {
	Hash       string
	Time       time.Time
	Script     string
	ScriptHash string
	Events     []Event
	Arguments  []string
}

type RawArgument struct {
	Value interface{}
	Type  string
}

func (self RawArgument) GetValue(_ int) string {
	return fmt.Sprintf("%v", self.Value)

}

type Event struct {
	Name   string
	Fields map[string]interface{}
}

type Entry struct {
	Transaction Transaction
	NFT         []NFTTransfer
	Tokens      []TokenTransfer
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransaction) Convert() Transaction {
	events := lo.Map(self.Events.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdge.Convert)

	res, err := json.Marshal(self.Arguments)
	if err != nil {
		panic(err)
	}
	var args []RawArgument
	err = json.Unmarshal(res, &args)
	if err != nil {
		panic(err)
	}

	return Transaction{
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

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdge) Convert(_ int) Event {
	keys := lo.Map(self.Node.Type.Fields, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTransactionEventsEventConnectionEdgesEventEdgeNodeEventTypeFieldsEventTypeField.Convert)
	values := self.Node.Fields

	subValuePointer, _ := gojsonpointer.NewJsonPointer("/value/value")
	valuePointer, _ := gojsonpointer.NewJsonPointer("/value")

	fields := map[string]interface{}{}
	for i, key := range keys {
		optionalValue, _, _ := subValuePointer.Get(values[i])
		value, _, _ := valuePointer.Get(values[i])
		if optionalValue == nil {
			fields[key] = value
		} else {
			fields[key] = optionalValue
		}
	}

	return Event{
		Name:   self.Node.Type.Id,
		Fields: fields,
	}

}

type NFTTransfer struct {
	From     string
	To       string
	Contract string
	Id       string
}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeNftTransfersNFTTransferConnectionEdgesNFTTransferEdge) Convert(_ int) NFTTransfer {

	return NFTTransfer{
		From:     self.Node.From.Address,
		To:       self.Node.To.Address,
		Contract: self.Node.Nft.Contract.Id,
		Id:       self.Node.Nft.NftId,
	}
}

type TokenTransfer struct {
	Type         string
	Amount       float64
	Token        string
	Counterparty string
}

//Use a bindings marshaller here
func ConvertString(value string) float64 {

	number, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return float64(number) / 100_000_000

}

func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge) Convert(_ int) TokenTransfer {
	//500_000_000
	//5.00000000

	return TokenTransfer{
		Type:         string(self.Node.Type),
		Amount:       ConvertString(self.Node.Amount.Value),
		Token:        self.Node.Amount.Token.Id,
		Counterparty: self.Node.Counterparty.Address,
	}
}
func (self AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdge) Convert(_ int) Entry {
	return Entry{
		Transaction: self.Transaction.Convert(),
		NFT:         lo.Map(self.NftTransfers.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeNftTransfersNFTTransferConnectionEdgesNFTTransferEdge.Convert),
		Tokens:      lo.Map(self.TokenTransfers.Edges, AccountTransfersAccountTransferTransactionsAccountTransferConnectionEdgesAccountTransferEdgeTokenTransfersTokenTransferConnectionEdgesTokenTransferEdge.Convert),
	}
}
