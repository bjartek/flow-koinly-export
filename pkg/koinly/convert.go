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
func Convert(address string, entry flowgraph.Entry, nftIdMapping *NFTId, packs *Packs) ([]Event, error) {

	//TODO: if a tx is not using blocto, we need to check if we have to remove the TokenTransfer that is the fee, or even add it as fee in the format

	fee := ""
	feeCurrency := ""

	transfers := []flowgraph.TokenTransfer{}
	//we remove the tokenTransfer that is a fee if it is here and add the fee field
	for _, transfer := range entry.Tokens {
		if transfer.Counterparty == "0xf919ee77447b7497" { //fee receiver
			fee = fmt.Sprintf("%f", transfer.Amount)
			feeCurrency = ConvertCurrency(transfer.Token)
		} else {
			transfers = append(transfers, transfer)
		}
	}
	entry.Tokens = transfers

	entries := []Event{}
	t := DateTime{Time: entry.Transaction.Time}

	event := Event{
		Date:        t,
		TxHash:      entry.Transaction.Hash,
		Description: fmt.Sprintf("url=https://f.dnz.dev/%s ", entry.Transaction.Hash),
		FeeAmount:   fee,
		FeeCurrency: feeCurrency,
	}

	numberOfFTTransfers := len(entry.Tokens)
	numberOfNFTTransfers := len(entry.NFT)

	if numberOfFTTransfers == 0 && numberOfNFTTransfers == 0 {
		//TODO consider saving as skipped
		return nil, nil
	}
	scriptHash := entry.Transaction.ScriptHash
	ignoreHashes := []string{
		//these can just be thrown away
		"15bbd08bc4c18fa30c9bcf0440cd93318cb447fbc9c14863e7465f087c8cf836", // versus bid, we handle this in settle
		"a120159c824203e71c7478314209c87fb4c18039f011000e6683951d79119b12", //starly pack purchase, impossible to correlate pack to item
		"0f8879f814e28abe0e38101d1865a0b05cf8fa59baa362b12c541409caff7e11", //flovatar attach component

		"7fd94abec6a05dfc3722ee8205943f8fa291b9b0a56dd2bc7279ad4084d10ecd", //versus list marketplace
		"f61c8f03b845500aa0baee2d659fae8719cb898e427f137263e2d974fe37dc3f", //flovatar list for sale component
		"d6cfd4842774c0618f209cfc6d5d17217464ed685fc9506ba16e890a6b1340ed", //flovatar mint, would be fun to implement this in tax sw :D
		"4ff04e92f3f7649b7c83dc03e5ea3da786768f7ef4c6de19a066d9f75ee36745", //versus bid, handled in settle!
		"37b2f71c2376e4946229a5a5583e01229ee7dcb5d2fd96e21d5e021129cdad83", //charity mint
		"00cf662d9cb1266d6add9707e1b85aaa54e0babcc7585663f2137b21d316cf4e", //flovatar listed for sale
		//todo these below we can readd
		"08c194dfd79cd979f8a6d37af91f5890b7f3278db179ea1bed004269ac85f3cc", //versus send NFT, if we chose to implement NFTID automatically here we can add this
		"896bf1ffd62c43b418d4923244fe9ac4b138c005f44e8fa2cc063aec40733ea8", //zay codes swap, if we chose NFTID then implement
		"fcdd17efba950df4e45d8885eb983433d8833bcd042d4e9d67d91cad94abe948", //zay codes swap
		"cb04c6a0531eb4e0c85f63625f13882ebdbf7152598fa7d6efbae717bade095b", //flovatar airdrop many (newer)
		"afb4f09ad0b295196bc13599b7d719de421c9252ad81e1b7971b3cd058031962", //flovatar airdrop
		"80c607d3c993d6617fe023400356e9d2fb86bbde04f2d24595a0831da54757c0", //add keys
		"5e3e570f99de92bf0beec6c83273359d9cfa9c80a278a403b8283e24dd2d3248", //TODO: swapping pool see tx fd87d4d86e145e26f41d21955527b719bc29d3aaae9767561757d57c99b01c2c
		"d0afe76af92b11ad43f26558ec3aecb1008a0429706aa8b6261ed08cc6b5f36b", //TODO: REVV defi tx b1bfb9b4c984625eeba828d55a37d0ffe146c1a32483fb4e93fce8ac1eddc953
		"abcdf6bb2842402cedec581bed263c03454e53439ed02ccdc052350db1b940c9", //TODO: more defi flow/usdt 209311220240e10c871ed6e63cc52148ba6e3187df7b4852d41abc66a0ea867e
		"dd217a2cb1aa7f4317bc4adc8cf1da6460915d2f36f411af783f9dbb71bd29b7", //TODO: more defi 5dfc4cc1b0721b4a54276b3dbdce5ba851b60c3ca02fe55218c27626fba9aff1

	}

	if slices.Contains(ignoreHashes, scriptHash) {
		return nil, nil
	}

	if scriptHash == "a10d1ab887ad772ba33a70878ff96c5b9420bf6ad6e78ee3d445c91a8fb949b7" {
		//airdropp Goobers, but we have paid for this before without adding value to the goober, so maybe find the goober send tx and remove that?
		for i, nft := range entry.NFT {
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			ev := event
			if i != 0 {
				//we only pay the fee once
				ev.FeeAmount = ""
				ev.FeeCurrency = ""
			}
			ev.Label = "airdrop"
			ev.ReceivedAmount = "1"
			ev.ReceivedCurrency = nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id))
			entries = append(entries, ev)
		}
		return entries, nil
	}

	if scriptHash == "444817259dec224209b32f97e190ba4e980545ffa9561b6e59c80ddc1ba48952" {
		//stake
		token := entry.Tokens[0]

		event.Label = "stake"
		event.SentAmount = fmt.Sprintf("%v", token.Amount)
		event.SentCurrency = ConvertCurrency(token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "56747555d52af593d3c56465852a3def2559fe36e5142e86edda4d91459266e6" {
		token := entry.Tokens[0]
		event.Label = "Trade"
		nftId := nftIdMapping.GetOrAdd("A.921ea449dffec68a.Flovatar.NFT", string(entry.Transaction.Arguments[1]))

		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftId
		}

		entries = append(entries, event)
		return entries, nil

	}

	if scriptHash == "47851586d962335e3f7d9e5d11a4c527ee4b5fd1c3895e3ce1b9c2821f60b166" ||
		scriptHash == "25717e66e70730e00440b2b9e52b581021825241cb540e46c0aa5cf4a9514b58" || //blocto
		scriptHash == "1f4921d504e24e11bd06e57feff2d6c3567893ab0e90aa6230b714d0dfad85aa" {
		token := entry.Tokens[0]
		event.Label = "gift"
		event.SentAmount = fmt.Sprintf("%v", token.Amount)
		event.SentCurrency = ConvertCurrency(token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "4a63987d53eab6579ee737a056f980a514b32edac558a67c5b951ac0696c7710" {
		token := entry.Tokens[0]
		event.Label = "reward"
		event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
		event.ReceivedCurrency = ConvertCurrency(token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	/*
		if scriptHash == "83ca36e3d8c492576e58540a4eb0fb321c1f56acb09b10bb3422527ed7e81e2c" {
			//this is a starly airdrop.. lets see what we get here
			litter.Dump(packs.Mappings)
			litter.Dump(entry)
			os.Exit(1)
		}

			if scriptHash == "a120159c824203e71c7478314209c87fb4c18039f011000e6683951d79119b12" {
				//starly packs are _not_ on chain...
				token := entry.Tokens[0]
				packs.Add(entry.Transaction.Arguments[1], token.Amount, token.Token)
				return nil, nil
			}
	*/

	if scriptHash == "c10d71b54483f24bae20db5109a748f792622dcaf170d61f6f8dfd37503a3a46" {
		//flovatar component market

		token := entry.Tokens[0]
		id := entry.Transaction.Arguments[1]
		event := Event{
			Date:        t,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("Flovatar.Component=%s url=https://f.dnz.dev/%s ", id, entry.Transaction.Hash),
		}

		nftId := nftIdMapping.GetOrAdd("A.921ea449dffec68a.FlovatarComponent.NFT", id)

		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftId
		}
		entries = append(entries, event)
		return entries, nil

	}
	if scriptHash == "2b26cb7784ee28b5747d3a04cff60ca3f9ae93cf555eca3d4b7572b54eb75f46" {
		//versus marketplace
		token := entry.Tokens[0]
		id := entry.Transaction.Arguments[1]
		event := Event{
			Date:        t,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("secondary Versus.Art=%s url=https://f.dnz.dev/%s ", id, entry.Transaction.Hash),
		}

		nftId := nftIdMapping.GetOrAdd("A.d796ff17107bbff6.Art.NFT", id)

		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftId
		}
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "cc36a9819d06d40062b621ec505d97ebd9238d9e096cdd2a4b84a3fd1e1df5e5" {
		//versus settle
		eventType := "A.d796ff17107bbff6.Auction.TokenPurchased"
		for _, event := range entry.Transaction.Events {
			to, _ := event.Fields["to"].(string)
			if event.Name == eventType && to == address {
				price, _ := event.Fields["price"].(string)
				artId, _ := event.Fields["artId"].(string)

				if artId == "" {
					litter.Dump(event.Fields)
					panic("foo2")
				}
				nftId := nftIdMapping.GetOrAdd("A.d796ff17107bbff6.Art.NFT", artId)

				event := Event{
					Date:             t,
					TxHash:           entry.Transaction.Hash,
					Description:      fmt.Sprintf("primary Versus.Art=%s url=https://f.dnz.dev/%s ", artId, entry.Transaction.Hash),
					SentAmount:       price,
					SentCurrency:     ConvertCurrency("A.1654653399040a61.FlowToken"),
					ReceivedAmount:   "1",
					ReceivedCurrency: nftId,
				}
				entries = append(entries, event)
			}
		}
		return entries, nil
	}

	//if scriptHash == "1ad96e0b57fb2a4fe61daa778111e2ce6eb84214bd915065b4e0f23ffedfa4f0" {
	//can we just do this here?
	if numberOfNFTTransfers > 0 && numberOfFTTransfers == 0 {
		//airdrops
		//TODO: here we might add something that allows you to specify mappings
		for _, nft := range entry.NFT {
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			event := Event{
				Date:             t,
				TxHash:           entry.Transaction.Hash,
				Description:      fmt.Sprintf(" url=https://f.dnz.dev/%s ", entry.Transaction.Hash),
				Label:            "airdrop",
				ReceivedAmount:   "1",
				ReceivedCurrency: nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id)),
			}
			entries = append(entries, event)
		}
		return entries, nil
	}
	//we spend some funds and get X nfts back
	if numberOfFTTransfers == 1 && numberOfNFTTransfers > 1 {
		token := entry.Tokens[0]
		eachSum := token.Amount / float64(len(entry.NFT))

		for _, nft := range entry.NFT {

			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			event := Event{
				Date:        t,
				TxHash:      entry.Transaction.Hash,
				Description: fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", token.Counterparty, entry.Transaction.Hash),
			}
			if token.Type == "Withdraw" {
				event.SentAmount = fmt.Sprintf("%v", eachSum)
				event.SentCurrency = ConvertCurrency(token.Token)
				event.ReceivedAmount = "1"
				event.ReceivedCurrency = nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id))

			} else if token.Type == "Deposit" {
				event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
				event.ReceivedCurrency = ConvertCurrency(token.Token)
				event.SentAmount = "1"
				event.SentCurrency = nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id))
			}
			entries = append(entries, event)
		}
		return entries, nil
	}
	if numberOfFTTransfers == 1 && numberOfNFTTransfers == 1 {
		token := entry.Tokens[0]
		nft := entry.NFT[0]

		eventName := fmt.Sprintf("%s.NFT", nft.Contract)
		event := Event{
			Date:        t,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", token.Counterparty, entry.Transaction.Hash),
		}
		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id))

		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftIdMapping.GetOrAdd(eventName, fmt.Sprint(nft.Id))
		}

		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "4968e16ef6c4b0fa5a16a321f3aaee98f202fe0f38759178d928a921fecd6ac4" {
		token := entry.Tokens[0]
		event := Event{
			Date:             t,
			TxHash:           entry.Transaction.Hash,
			Description:      fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", token.Counterparty, entry.Transaction.Hash),
			ReceivedAmount:   fmt.Sprintf("%v", token.Amount),
			ReceivedCurrency: ConvertCurrency(token.Token),
			Label:            "reward",
		}
		entries = append(entries, event)
		return entries, nil
	}

	if numberOfFTTransfers == 1 {
		token := entry.Tokens[0]
		event := Event{
			Date:        t,
			TxHash:      entry.Transaction.Hash,
			Description: fmt.Sprintf("counterparty=%s url=https://f.dnz.dev/%s ", token.Counterparty, entry.Transaction.Hash),
		}
		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.Label = "gift"

			event.SentCurrency = ConvertCurrency(token.Token)
		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.Label = "airdrop"
		} else {
			litter.Dump(entry)
			os.Exit(1)
		}
		entries = append(entries, event)
		return entries, nil
	}

	if numberOfFTTransfers == 2 {

		event := Event{
			Date:        t,
			TxHash:      entry.Transaction.Hash,
			Label:       "Trade",
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
		return entries, nil
	}

	if numberOfFTTransfers == 3 {
		//DEFI stuff
		return nil, nil
	}
	//versus market

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

	return entries, nil

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
