package flowgraph

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/bjartek/flow-koinly-export/pkg/core"
	"github.com/samber/lo"
	"github.com/xeipuuv/gojsonpointer"
)

func TransactionEventsFirst(
	ctx context.Context,
	client graphql.Client,
	txID string,
) (*TransactionEventsResponse, error) {
	req := &graphql.Request{
		OpName: "TransactionEvents",
		Query: `
query TransactionEvents ($txID: ID!) {
	transaction(id: $txID) {
		events(first: 50) {
			pageInfo {
				endCursor
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
}
`,
		Variables: &__TransactionEventsInput{
			TxID: txID,
		},
	}
	var err error

	var data TransactionEventsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func (self TransactionEventsTransactionEventsEventConnectionEdgesEventEdgeNodeEventTypeFieldsEventTypeField) Convert(_ int) string {
	return self.Identifier
}

func (self TransactionEventsTransactionEventsEventConnectionEdgesEventEdge) Convert(_ int) core.Event {
	keys := lo.Map(self.Node.Type.Fields, TransactionEventsTransactionEventsEventConnectionEdgesEventEdgeNodeEventTypeFieldsEventTypeField.Convert)
	values := self.Node.Fields

	subValuePointer, _ := gojsonpointer.NewJsonPointer("/value/value")
	valuePointer, _ := gojsonpointer.NewJsonPointer("/value")

	valuesLength := len(values)
	fields := map[string]interface{}{}
	for i, key := range keys {
		if i >= valuesLength {
			continue
		}
		optionalValue, _, _ := subValuePointer.Get(values[i])
		value, _, _ := valuePointer.Get(values[i])
		if optionalValue == nil {
			fields[key] = value
		} else {
			fields[key] = optionalValue
		}
	}

	return core.Event{
		Name:   self.Node.Type.Id,
		Fields: fields,
	}

}
