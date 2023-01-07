package koinly

import (
	"fmt"
	"os"

	"github.com/bjartek/flow-koinly-export/pkg/flowgraph"
	"github.com/sanity-io/litter"
	"golang.org/x/exp/slices"
)

/*
Labels

Labels can be added as appropriate. For regular deposits/withdrawals/trades, no label is required.
Koinly allows the following labels for outgoing transactions:
    gift
    lost
    cost
    margin fee
    realized gain
    stake

The following labels are allowed for incoming transactions:

    airdrop
    fork
    mining
    reward
    income (for other income)
    loan interest
    realized gain
    unstake


*/
//atm this just converts if there are FT involved or not
func Convert(address string, entry flowgraph.Entry) ([]Event, error) {

	entries := []Event{}
	time := DateTime{Time: entry.Transaction.Time}
	numberOfFTTransfers := len(entry.Tokens)

	scriptHash := entry.Transaction.ScriptHash
	ignoreHashes := []string{
		"7fd94abec6a05dfc3722ee8205943f8fa291b9b0a56dd2bc7279ad4084d10ecd", //versus list marketplace
		"08c194dfd79cd979f8a6d37af91f5890b7f3278db179ea1bed004269ac85f3cc", //versus send NFT, if we chose to implement NFTID automatically here we can add this
		"896bf1ffd62c43b418d4923244fe9ac4b138c005f44e8fa2cc063aec40733ea8", //zay codes swap, if we chose NFTID then implement
		"fcdd17efba950df4e45d8885eb983433d8833bcd042d4e9d67d91cad94abe948", //zay codes swap
		"f61c8f03b845500aa0baee2d659fae8719cb898e427f137263e2d974fe37dc3f", //flovatar list for sale component
		"afb4f09ad0b295196bc13599b7d719de421c9252ad81e1b7971b3cd058031962", //flovatar airdrop
		"d6cfd4842774c0618f209cfc6d5d17217464ed685fc9506ba16e890a6b1340ed", //flovatar mint, would be fun to implement this in tax sw :D
		"83ca36e3d8c492576e58540a4eb0fb321c1f56acb09b10bb3422527ed7e81e2c", //starly airdrop, NFTID implement?
		"4ff04e92f3f7649b7c83dc03e5ea3da786768f7ef4c6de19a066d9f75ee36745", //versus bid, handled in settle!
	}

	if slices.Contains(ignoreHashes, scriptHash) {
		return nil, nil
	}
	if numberOfFTTransfers == 1 {
		litter.Dump(entry)
		os.Exit(1)
		token := entry.Tokens[0]
		event := Event{
			Date:        time,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", token.Counterparty, entry.Transaction.Hash),
		}
		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
		} else {
			litter.Dump(entry)
			os.Exit(1)
		}
		entries = append(entries, event)
	} else if numberOfFTTransfers == 2 {

		event := Event{
			Date:        time,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("url=https://f.dnz.dev/%s ", entry.Transaction.Hash),
		}
		for _, token := range entry.Tokens {
			if token.Type == "Withdraw" {
				event.SentAmount = fmt.Sprintf("%v", token.Amount)
				event.SentCurrency = ConvertCurrency(token.Token)
			} else if token.Type == "Deposit" {
				event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
				event.ReceivedCurrency = ConvertCurrency(token.Token)
			}
		}
		entries = append(entries, event)
	} else {
		//versus settle
		if scriptHash == "cc36a9819d06d40062b621ec505d97ebd9238d9e096cdd2a4b84a3fd1e1df5e5" {

			/*
									"artId": "892",
				          "from": "0x1b945b52f416ddf9",
				          "id": "746",
				          "price": "22.00000000",
				          "to": "0x626e2a193431984d",
			*/

			eventType := "A.d796ff17107bbff6.Auction.TokenPurchased"
			for _, event := range entry.Transaction.Events {
				to, _ := event.Fields["to"].(string)
				price, _ := event.Fields["price"].(string)
				artId, _ := event.Fields["artId"].(string)
				if event.Name == eventType && to == address {
					event := Event{
						Date:         time,
						TxHash:       entry.Transaction.Hash,
						Description:  fmt.Sprintf("Versus.Art=%s url=https://f.dnz.dev/%s ", artId, entry.Transaction.Hash),
						SentAmount:   price,
						SentCurrency: ConvertCurrency("A.1654653399040a61.FlowToken"),
					}
					//TODO: consider adding ReceivedAmount=1 and NFTXX here
					entries = append(entries, event)
					return entries, nil
				}
			}
		}

		/*
			//versus send NFT, this might be handled generically I do not know
			if scriptHash == "08c194dfd79cd979f8a6d37af91f5890b7f3278db179ea1bed004269ac85f3cc" {

				event := Event{
					Date:         time,
					TxHash:       entry.Transaction.Hash,
					Description:  fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", entry.NFT[0].To, entry.Transaction.Hash),
					SentAmount:   price,
					SentCurrency: ConvertCurrency("A.1654653399040a61.FlowToken"),
				}
				entries = append(entries, event)
				return entries, nil

				/*
				 NFT: []flowgraph.NFTTransfer{
				    flowgraph.NFTTransfer{
				      From: "0x886f3aeaf848c535",
				      To: "0x4bcda1de73a17d95",
				      Contract: "A.d796ff17107bbff6.Art",
				      Id: "289",
				    },
				  },
				  Tokens: []flowgraph.TokenTransfer{},
		*/
		//}

		litter.Dump(entry)
		fmt.Println("Count not handle tx")
		os.Exit(0)
	}

	return entries, nil

	/*
	 Tokens: []flowgraph.TokenTransfer{
	      flowgraph.TokenTransfer{
	        Type: "Withdraw",
	        Amount: 5.0,
	        Token: "A.1654653399040a61.FlowToken",
	        Counterparty: "0xd796ff17107bbff6",
	      },
	    },
	*/

}

func ConvertCurrency(currency string) string {

	currencyMap := map[string]string{
		"A.3c5959b568896393.FUSD":                  "ID:3054",
		"A.1654653399040a61.FlowToken":             "ID:7961",
		"A.0f9df91c9121c460.BloctoToken":           "ID:35927",
		"A.cfdd90d4a00f7b5b.TeleportedTetherToken": "USDT",
	}
	return currencyMap[currency]

}
