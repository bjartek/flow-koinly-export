package koinly

import (
	"fmt"
	"os"

	"github.com/bjartek/flow-koinly-export/pkg/core"
	"github.com/pkg/errors"
	"github.com/samber/lo"
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

const NameFindPack = "A.097bafa4e0b48eef.FindPack.NFT"
const NameVersusArt = "A.d796ff17107bbff6.Art.NFT"
const NameFlovatarNFT = "A.921ea449dffec68a.Flovatar.NFT"
const NameBl0xPack = "A.7620acf6d7f2468a.Bl0xPack.NFT"
const ZayTraderEvent = "A.4c577a03bc1a82e0.ZayTraderV2.TradeExecuted"

// atm this just converts if there are FT involved or not
func Convert(address string, entry core.Entry, state *core.State) ([]Event, error) {

	//TODO: starly staking bug aaf2ac69a5ea03a24294e8545fcfc5010151ca81d16b815b3a1ffa03d2b6998a
	//TODO; find names https://f.dnz.dev/39fe10d09d0457f9fc69c2148aa20fe1d435f218437ecc23e54551f415e13add

	fee := ""
	feeCurrency := ""

	transfers := []core.TokenTransfer{}
	//we remove the tokenTransfer that is a fee if it is here and add the fee field
	for _, transfer := range entry.Tokens {
		if transfer.Counterparty == "0xf919ee77447b7497" && transfer.Type == "Withdraw" { //fee receiver
			fee = fmt.Sprintf("%f", transfer.Amount)
			feeCurrency = ConvertCurrency(transfer.Token)
		} else {
			transfers = append(transfers, transfer)
		}
	}
	entry.Tokens = transfers

	entries := []Event{}
	t := DateTime{Time: entry.Transaction.Time}

	scriptHash := entry.Transaction.ScriptHash
	event := Event{
		Date:        t,
		TxHash:      entry.Transaction.Hash,
		Description: fmt.Sprintf("url=https://f.dnz.dev/%s script=%s", entry.Transaction.Hash, scriptHash),
		FeeAmount:   fee,
		FeeCurrency: feeCurrency,
	}

	isZayTrade := lo.ContainsBy(entry.Transaction.Events, func(e core.Event) bool {
		return e.Name == ZayTraderEvent
	})

	ftSend := []core.TokenTransfer{}
	ftReceived := []core.TokenTransfer{}
	for _, t := range entry.Tokens {
		if t.Type == "Deposit" {
			ftReceived = append(ftReceived, t)
		}
		ftSend = append(ftSend, t)
	}
	nftSend := []core.NFTTransfer{}
	nftReceived := []core.NFTTransfer{}
	for _, t := range entry.NFT {
		if t.To == address {
			nftReceived = append(nftReceived, t)
		}
		nftSend = append(nftSend, t)
	}

	numberOfFTTransfers := len(entry.Tokens)
	numberOfNFTTransfers := len(entry.NFT)

	if numberOfFTTransfers == 0 && numberOfNFTTransfers == 0 {
		//TODO consider saving as skipped
		return nil, nil
	}
	ignoreHashes := []string{
		//these can just be thrown away
		"1f177b71729f1a54ce62ca64d246f38418c8fb78189a1b223be8be479e15a350", //delist flovatar
		"7a7a69fdd932f4c47e6677c32f151edb9f70df19a8f2c31f3574f761a5b8ebe2", //versus list
		"209a21a38379d2322c72bebf7bb50060c70ccae0da1f6324bc02561278c89ffa", //open a find pack does not do anything
		"15bbd08bc4c18fa30c9bcf0440cd93318cb447fbc9c14863e7465f087c8cf836", // versus bid, we handle this in settle
		"a120159c824203e71c7478314209c87fb4c18039f011000e6683951d79119b12", //starly pack purchase, impossible to correlate pack to item
		"7fd94abec6a05dfc3722ee8205943f8fa291b9b0a56dd2bc7279ad4084d10ecd", //versus list marketplace
		"f61c8f03b845500aa0baee2d659fae8719cb898e427f137263e2d974fe37dc3f", //flovatar list for sale component
		"4ff04e92f3f7649b7c83dc03e5ea3da786768f7ef4c6de19a066d9f75ee36745", //versus bid, handled in settle!
		"00cf662d9cb1266d6add9707e1b85aaa54e0babcc7585663f2137b21d316cf4e", //flovatar listed for sale
		"926cf83e0b0de6bb32291b93b1070fc98224686b3c6f7d7962d8fd287b98596b", //flovatar listed for sale new
		"37b2f71c2376e4946229a5a5583e01229ee7dcb5d2fd96e21d5e021129cdad83", //charity mint
		"80c607d3c993d6617fe023400356e9d2fb86bbde04f2d24595a0831da54757c0", //add keys

		//todo these below we can readd
		"fcdd17efba950df4e45d8885eb983433d8833bcd042d4e9d67d91cad94abe948", //zay codes swap
		"5e3e570f99de92bf0beec6c83273359d9cfa9c80a278a403b8283e24dd2d3248", //TODO: swapping pool see tx fd87d4d86e145e26f41d21955527b719bc29d3aaae9767561757d57c99b01c2c
		"d0afe76af92b11ad43f26558ec3aecb1008a0429706aa8b6261ed08cc6b5f36b", //TODO: REVV defi tx b1bfb9b4c984625eeba828d55a37d0ffe146c1a32483fb4e93fce8ac1eddc953
		"abcdf6bb2842402cedec581bed263c03454e53439ed02ccdc052350db1b940c9", //TODO: more defi flow/usdt 209311220240e10c871ed6e63cc52148ba6e3187df7b4852d41abc66a0ea867e
		"dd217a2cb1aa7f4317bc4adc8cf1da6460915d2f36f411af783f9dbb71bd29b7", //TODO: more defi 5dfc4cc1b0721b4a54276b3dbdce5ba851b60c3ca02fe55218c27626fba9aff1

	}

	if slices.Contains(ignoreHashes, scriptHash) {
		return nil, nil
	}

	if scriptHash == "2a1e7927441136c24b1eadcb316abc96c8b32faf403c6bf2d5e412e3e71bf51a" {
		//find lease extended
		return nil, nil
	}

	if scriptHash == "581049d525f2c0c7a21eeb9f6277747c5f0291a1778bf4395423aa10c1e6abb3" {
		//flobot create
		return nil, nil
	}

	if scriptHash == "d6cfd4842774c0618f209cfc6d5d17217464ed685fc9506ba16e890a6b1340ed" {
		//flovatar create

		destroyedIds := []string{}
		eventType := "A.921ea449dffec68a.FlovatarComponent.Destroyed"
		for _, e := range entry.Transaction.Events {
			if e.Name == eventType {
				id := e.Fields["id"].(string)
				destroyedIds = append(destroyedIds, id)
			}
		}

		flovatarId := ""
		for i, nft := range entry.NFT {
			ev := event
			if i != 0 {
				//we only pay the fee once
				ev.FeeAmount = ""
				ev.FeeCurrency = ""
			}
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)

			//the components used to create the flovatar are cost
			if nft.Contract == "A.921ea449dffec68a.FlovatarComponent" && slices.Contains(destroyedIds, nft.Id) {
				//we have to have the NFTID for this component already if not something is wrong
				componentId := state.AddNFTID(eventName, nft.Id)
				ev.Label = "cost"
				ev.SentAmount = "1"
				ev.SentCurrency = componentId
				entries = append(entries, ev)
			}

			if nft.Contract == "A.921ea449dffec68a.Flovatar" {
				componentId := state.AddNFTID(eventName, nft.Id)
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = componentId
				flovatarId = componentId
				entries = append(entries, ev)
			}
		}

		for _, nft := range entry.NFT {
			if nft.Contract == "A.921ea449dffec68a.Flovatar" {
				continue
			}
			if nft.Contract == "A.921ea449dffec68a.FlovatarComponent" && slices.Contains(destroyedIds, nft.Id) {
				continue
			}

			//all components that where not destoryed as part of mint and is not the flovatar is on that flovatar
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			componentId, err := state.GetNFTID(eventName, nft.Id)

			if err != nil {
				return nil, err
			}
			//we add the subComponent to state to tell that if we sell this we also sell this component
			err = state.AddCompositeComponent(flovatarId, componentId)
			if err != nil {
				return nil, errors.Wrap(err, "mint_flovatar")
			}
		}
		return entries, nil

	}

	if scriptHash == "e765cca2e07d6d14232bcd3924192051c8e50ab0c4c121dad508fa3652103a79" {
		//TODO: jambb momenets
		return nil, nil
	}

	if slices.Contains([]string{
		"bcf4800902a1cdccbae4966b32f728c036a8e82e97dd477bc4e44e9ad1ff8d23",
		"0f8879f814e28abe0e38101d1865a0b05cf8fa59baa362b12c541409caff7e11",                //flovatar attach background
		"2206353dc476ab1fadec47fe1ea083a66bb0ec72e88f976a09aaeea2484b1f80",                //eyes //attach glasses
		"86bb52678dd47d7abcc457e4cbc473f7b8c237c0a882fc950327adc992a8f564",                //hat
		"96e4656bb3d929a909be26404d8a9cadcc1a3773a49e20af4addb6bc416c704e"}, scriptHash) { //flovatar attach accessory

		//somehow we have flovatars that are not here before?
		flovatar := state.AddNFTID(NameFlovatarNFT, entry.Transaction.Arguments[0])

		for _, nft := range entry.NFT {
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			//this might be a flovatarComponent we have never owned before
			componentId := state.AddNFTID(eventName, nft.Id)
			if nft.From == "" {
				present := state.RemoveCompositeComponent(flovatar, componentId)
				if !present {
					//we have removed a component from flovatar that we did not add it to, so we need to add it as an airdrop to ourself
					//TODO; this should not incur tax
					ev := event
					ev.Label = "gift"
					ev.ReceivedAmount = "1"
					ev.ReceivedCurrency = componentId
					entries = append(entries, event)
				}
			} else {
				err := state.AddCompositeComponent(flovatar, componentId)
				if err != nil {
					return nil, errors.Wrap(err, "equip flobit")
				}
			}
		}
		return entries, nil
	}

	if scriptHash == "a10d1ab887ad772ba33a70878ff96c5b9420bf6ad6e78ee3d445c91a8fb949b7" {
		//TODO This is standard airdrop and can go away
		//airdropp Goobers
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
			ev.ReceivedCurrency = state.AddNFTID(eventName, fmt.Sprint(nft.Id))
			entries = append(entries, ev)
		}
		return entries, nil
	}

	if scriptHash == "74b37352af81d750cb29f94e21ab449f09b2623e60a0074a75ac3f5aa09a10f7" || scriptHash == "a83c8ec471199a12b496f070ea6036962a49f9c1f7210ed01280ef1603deae71" {
		//buy flovatarPack

		//what do we do with packs that you bought but have not opened?
		token := entry.Tokens[0]
		event.Label = "Trade"
		event.SentAmount = fmt.Sprintf("%v", token.Amount)
		event.SentCurrency = ConvertCurrency(token.Token)
		event.ReceivedAmount = "1"
		event.ReceivedCurrency = state.AddNFTID("A.921ea449dffec68a.FlovatarPack", entry.Transaction.Arguments[1])
		state.AddPack(event.ReceivedCurrency, token.Amount, token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "ecc73476640233932aefa6aed80688e877907968a27337cf35af0a3ee86c6d98" || scriptHash == "d0ec403b0c44a4936d563834bdb322734fa514fb7d6fcc60843a663116cf1724" {
		//open flovatar pack

		// somehow I have flovar packs that i have never received
		packId := state.AddNFTID("A.921ea449dffec68a.FlovatarPack", entry.Transaction.Arguments[0])

		packPrice, _ := state.GetPack(packId)
		amountPerEntry := packPrice.Amount / float64(len(entry.NFT))
		for _, nft := range entry.NFT {

			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			nftId := state.AddNFTID(eventName, fmt.Sprint(nft.Id))
			ev := event
			//TODO: Pretty sure I am calculating gains here twice
			ev.Label = "Swap"
			ev.ReceivedAmount = "1"
			ev.ReceivedCurrency = nftId
			ev.SentCurrency = ConvertCurrency(packPrice.Currency)
			ev.SentAmount = fmt.Sprintf("%f", amountPerEntry)

			entries = append(entries, ev)
		}
		return entries, nil
	}

	if scriptHash == "444817259dec224209b32f97e190ba4e980545ffa9561b6e59c80ddc1ba48952" || scriptHash == "7c254512fa53e57314cb9070c3557db276167276c758af19ada207da2cdd3ffd" {
		//stake
		token := entry.Tokens[0]

		event.Label = "stake"
		event.SentAmount = fmt.Sprintf("%v", token.Amount)
		event.SentCurrency = ConvertCurrency(token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "56747555d52af593d3c56465852a3def2559fe36e5142e86edda4d91459266e6" {
		//trade flovatar
		token := entry.Tokens[0]
		event.Label = "Trade"

		ev := event
		ev.Label = "Trade"
		if token.Type == "Withdraw" {
			nftId := state.AddNFTID("A.921ea449dffec68a.Flovatar.NFT", string(entry.Transaction.Arguments[1]))
			ev.SentAmount = fmt.Sprintf("%v", token.Amount)
			ev.SentCurrency = ConvertCurrency(token.Token)
			ev.ReceivedAmount = "1"
			ev.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			nftId, err := state.GetNFTID("A.921ea449dffec68a.Flovatar.NFT", string(entry.Transaction.Arguments[1]))
			if err != nil {
				return nil, err
			}
			ev.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			ev.ReceivedCurrency = ConvertCurrency(token.Token)
			ev.SentAmount = "1"
			ev.SentCurrency = nftId

			components := state.GetCompositeComponent(nftId)
			for _, component := range components {
				ev2 := event
				ev2.SentCurrency = component
				ev2.SentAmount = "1"
				ev2.Label = "cost"
				ev2.Description = fmt.Sprintf("%s component sold as part of main sale", ev2.Description)
				entries = append(entries, ev2)
			}
		}
		entries = append(entries, ev)
		return entries, nil
	}

	//FT gifts
	if scriptHash == "47851586d962335e3f7d9e5d11a4c527ee4b5fd1c3895e3ce1b9c2821f60b166" ||
		scriptHash == "25717e66e70730e00440b2b9e52b581021825241cb540e46c0aa5cf4a9514b58" || //blocto
		scriptHash == "1f4921d504e24e11bd06e57feff2d6c3567893ab0e90aa6230b714d0dfad85aa" {
		token := entry.Tokens[0]
		event.Label = ""
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

	if scriptHash == "c10d71b54483f24bae20db5109a748f792622dcaf170d61f6f8dfd37503a3a46" || scriptHash == "dad0d9a3cd6c593a7b3a6c04a0ea227c3b93e4badf81ccd4b2d353bf6722a9fc" {
		//flovatar component market

		token := entry.Tokens[0]
		id := entry.Transaction.Arguments[1]

		event.Label = "Trade"
		if token.Type == "Withdraw" {
			nftId := state.AddNFTID("A.921ea449dffec68a.FlovatarComponent.NFT", id)
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			nftId, err := state.GetNFTID("A.921ea449dffec68a.FlovatarComponent.NFT", id)
			if err != nil {
				return nil, err
			}
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

		event.Label = "Trade"
		token := entry.Tokens[0]
		id := entry.Transaction.Arguments[1]

		if token.Type == "Withdraw" {
			nftId := state.AddNFTID(NameVersusArt, id)
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = nftId

		} else if token.Type == "Deposit" {
			//some versus is not there
			nftId := state.AddNFTID(NameVersusArt, id)
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftId
		}
		entries = append(entries, event)
		return entries, nil
	}

	if scriptHash == "cc36a9819d06d40062b621ec505d97ebd9238d9e096cdd2a4b84a3fd1e1df5e5" {
		/*
			//TODO: for the old pattern we have to keep track of bids...
			//for old patterns we cannot do anything.. we have to point them to the webpage, or actually run the script and check?
		*/
		//versus settle
		//TokenPurchased is not here in all versus sales...
		eventType := "A.d796ff17107bbff6.Auction.TokenPurchased"
		for _, e := range entry.Transaction.Events {
			to, _ := e.Fields["to"].(string)
			if e.Name == eventType && to == address {
				price, _ := e.Fields["price"].(string)
				artId, _ := e.Fields["artId"].(string)
				nftId := state.AddNFTID(NameVersusArt, artId)
				ev := event
				ev.Label = "Trade"
				ev.SentAmount = price
				ev.SentCurrency = ConvertCurrency("A.1654653399040a61.FlowToken")
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = nftId
				entries = append(entries, ev)
			}
		}
		if len(entries) == 0 {

			for _, nft := range entry.NFT {
				ev := event
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = state.AddNFTID(NameVersusArt, nft.Id)
				ev.Description = fmt.Sprintf("%s lookup drop to find price https://versus.acution/drops/%s ", ev.Description, entry.Transaction.Arguments[0])
				entries = append(entries, ev)
			}
		}
		return entries, nil
	}

	//This should be generalized to only check the label type against a scriptHash, since the logic is the same for other tx of the same pattern
	if scriptHash == "4968e16ef6c4b0fa5a16a321f3aaee98f202fe0f38759178d928a921fecd6ac4" || scriptHash == "a8e34487a5ebb460fb05dbcfb06e7f00b5dbe5e41462ddcbfd550e078a87d8f0" {
		token := entry.Tokens[0]
		event.Label = "reward"
		event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
		event.ReceivedCurrency = ConvertCurrency(token.Token)
		entries = append(entries, event)
		return entries, nil
	}

	/*
		if scriptHash == "c2080352d0b3a813a73f9d48b94e4371aef10a20416da310c0ecb3b69dc8acf8" {
			//FindPack buy

			if numberOfFTTransfers == 0 {
				//we got airdropped a pack so it is worth nothing for us, handle it in open pack
				return nil, nil
			}

			token := entry.Tokens[0]
			sumPerPack := token.Amount / float64(len(entry.NFT))

			for i, nft := range entry.NFT {
				ev := event

				if i != 0 {
					ev.FeeAmount = ""
					ev.FeeCurrency = ""
				}
				eventName := fmt.Sprintf("%s.NFT", nft.Contract)
				ev.Label = "Trade"
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = state.AddNFTID(eventName, fmt.Sprint(nft.Id))
				ev.SentAmount = fmt.Sprintf("%v", sumPerPack)
				ev.SentCurrency = ConvertCurrency(token.Token)
				state.AddPack(ev.ReceivedCurrency, sumPerPack, token.Token)

				entries = append(entries, ev)
			}
			return entries, nil

		}
	*/

	if scriptHash == "27c50fc75f5a8812748ed3cc39dacde12236012580e681c23bded52979defb5d" {
		//Bl0xPack  redeem lots of duplication from below

		//a map to hold all packs that where revealed to you in this tx
		packMappings := map[string][]string{}

		for _, ev := range entry.Transaction.Events {
			if ev.Name != "A.097bafa4e0b48eef.Bl0xPack.PackReveal" {
				continue
			}
			to, ok := ev.Fields["address"].(string)
			if !ok {
				continue
			}

			if to != address {
				continue
			}

			packId, _ := ev.Fields["packId"].(string)
			rewardId, _ := ev.Fields["rewardId"].(string)
			rewardType, _ := ev.Fields["rewardType"].(string)
			rewardNFT := state.AddNFTID(rewardType, rewardId)

			packMappingId := state.AddNFTID(NameBl0xPack, packId)

			packForId, ok := packMappings[packMappingId]
			if !ok {
				packForId = []string{}
			}
			packForId = append(packForId, rewardNFT)
			packMappings[packMappingId] = packForId
		}
		//we now have a multimap of NFTX -> NFTX,NFTX where the first is an pack NFT and the others are reward NFTS

		for packId, rewards := range packMappings {
			packPurchase, ok := state.GetPack(packId)
			if !ok {
				for _, reward := range rewards {
					ev := event
					ev.Label = "airdrop"
					ev.ReceivedAmount = "1"
					ev.ReceivedCurrency = reward
					entries = append(entries, ev)
				}
			} else {
				pricePerReward := packPurchase.Amount / float64(len(rewards))

				for _, reward := range rewards {
					ev := event
					ev.Label = "Swap"
					ev.SentAmount = fmt.Sprintf("%f", pricePerReward)
					ev.SentCurrency = ConvertCurrency(packPurchase.Currency)
					ev.ReceivedAmount = "1"
					ev.ReceivedCurrency = reward
					entries = append(entries, ev)
				}

			}
		}
		return entries, nil
	}

	if scriptHash == "80719f6a41daeb27f9ac5d7f49ea02b4adb45b1ee34272cd42603e6ca06aaeb3" {
		//FindPack redeem

		//a map to hold all packs that where revealed to you in this tx
		packMappings := map[string][]string{}

		for _, ev := range entry.Transaction.Events {
			if ev.Name != "A.097bafa4e0b48eef.FindPack.PackReveal" {
				continue
			}
			to, ok := ev.Fields["address"].(string)
			if !ok {
				continue
			}

			if to != address {
				continue
			}

			packId, _ := ev.Fields["packId"].(string)
			rewardId, _ := ev.Fields["rewardId"].(string)
			rewardType, _ := ev.Fields["rewardType"].(string)
			rewardNFT := state.AddNFTID(rewardType, rewardId)

			packMappingId := state.AddNFTID(NameFindPack, packId)

			packForId, ok := packMappings[packMappingId]
			if !ok {
				packForId = []string{}
			}
			packForId = append(packForId, rewardNFT)
			packMappings[packMappingId] = packForId
		}
		//we now have a multimap of NFTX -> NFTX,NFTX where the first is an pack NFT and the others are reward NFTS

		for packId, rewards := range packMappings {
			packPurchase, ok := state.GetPack(packId)
			if !ok {
				for _, reward := range rewards {
					ev := event
					ev.Label = "airdrop"
					ev.ReceivedAmount = "1"
					ev.ReceivedCurrency = reward
					entries = append(entries, ev)
				}
			} else {
				pricePerReward := packPurchase.Amount / float64(len(rewards))

				for _, reward := range rewards {
					ev := event
					ev.Label = "Swap"
					ev.SentAmount = fmt.Sprintf("%f", pricePerReward)
					ev.SentCurrency = ConvertCurrency(packPurchase.Currency)
					ev.ReceivedAmount = "1"
					ev.ReceivedCurrency = reward
					entries = append(entries, ev)
				}

			}
		}
		return entries, nil
	}

	//if scriptHash == "1ad96e0b57fb2a4fe61daa778111e2ce6eb84214bd915065b4e0f23ffedfa4f0" {
	//can we just do this here?
	if numberOfNFTTransfers > 0 && numberOfFTTransfers == 0 {
		//airdrops
		for i, nft := range entry.NFT {
			eventName := fmt.Sprintf("%s.NFT", nft.Contract)

			ev := event
			if i != 0 {
				ev.FeeAmount = ""
				ev.FeeCurrency = ""
			}

			if nft.From == address {
				nftId, err := state.GetNFTID(eventName, fmt.Sprint(nft.Id))
				if err != nil {
					nftId = "TODO-send-single-nft"
					//return nil, errors.Wrap(err, "airdrop single nft")
				}
				ev.SentAmount = "1"
				ev.SentCurrency = nftId
				//				ev.Label = "income"
			} else {
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = state.AddNFTID(eventName, fmt.Sprint(nft.Id))
				//				ev.Label = "airdrop"
			}
			entries = append(entries, ev)
		}
		return entries, nil
	}

	if isZayTrade {

		for _, nft := range entry.NFT {

			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			nftId := state.AddNFTID(eventName, fmt.Sprint(nft.Id))

			ev := event
			ev.Label = "ZaySwap! NB Review!"
			if nft.From == address {
				ev.SentAmount = "1"
				ev.SentCurrency = nftId
			} else {
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = nftId
			}
			entries = append(entries, ev)
		}

		for _, ft := range entry.Tokens {
			ev := event
			ev.Label = "ZaySwap! NB Review!"

			currency := ConvertCurrency(ft.Token)
			if ft.Type == "Deposit" {
				ev.SentCurrency = currency
				ev.SentAmount = fmt.Sprintf("%v", ft.Amount)
			} else {
				ev.ReceivedAmount = fmt.Sprintf("%v", ft.Amount)
				ev.ReceivedCurrency = currency
			}
			entries = append(entries, ev)
		}

		return entries, nil
	}

	if numberOfFTTransfers == 1 && numberOfNFTTransfers > 1 {
		token := entry.Tokens[0]
		eachSum := token.Amount / float64(len(entry.NFT))

		for i, nft := range entry.NFT {

			eventName := fmt.Sprintf("%s.NFT", nft.Contract)
			ev := event
			ev.Label = "Trade"
			if token.Type == "Withdraw" {
				if i != 0 {
					ev.FeeAmount = ""
					ev.FeeCurrency = ""
				}
				ev.SentAmount = fmt.Sprintf("%v", eachSum)
				ev.SentCurrency = ConvertCurrency(token.Token)
				ev.ReceivedAmount = "1"
				ev.ReceivedCurrency = state.AddNFTID(eventName, fmt.Sprint(nft.Id))

				//we buy bl0x packs so we add the pack price to the registry
				if eventName == NameBl0xPack || eventName == NameFindPack {
					state.Packs.Add(ev.ReceivedCurrency, eachSum, token.Token)
				}

			} else if token.Type == "Deposit" {
				//we only pay the fee once
				if i != 0 {
					ev.FeeAmount = ""
					ev.FeeCurrency = ""
				}
				nftId := state.AddNFTID(eventName, fmt.Sprint(nft.Id))
				ev.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
				ev.ReceivedCurrency = ConvertCurrency(token.Token)
				ev.SentAmount = "1"
				ev.SentCurrency = nftId
			}
			entries = append(entries, ev)
		}
		return entries, nil
	}

	if numberOfFTTransfers == 1 && numberOfNFTTransfers == 1 {
		token := entry.Tokens[0]
		nft := entry.NFT[0]

		if token.Type == "Deposit" && nft.To == address {
			//I think it is a good idea to pre run all tx where we classify transfers as incomming/outgoing. We can do the same with Events/arguments to create more transfers
			//3aef58f5eaac6d8809948e9bc343552365f00df5a8ae47392f026dfe71e73df4
			//TODO double airdrop?
			return nil, nil
		}
		//we have cases where this is double airdrop?
		eventName := fmt.Sprintf("%s.NFT", nft.Contract)
		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.ReceivedAmount = "1"
			event.ReceivedCurrency = state.AddNFTID(eventName, fmt.Sprint(nft.Id))

		} else if token.Type == "Deposit" {
			nftId, err := state.GetNFTID(eventName, fmt.Sprint(nft.Id))
			if err != nil {
				return nil, errors.Wrap(err, "buy single nft")
			}
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.SentAmount = "1"
			event.SentCurrency = nftId
		}

		entries = append(entries, event)
		return entries, nil
	}

	if numberOfFTTransfers == 1 {
		token := entry.Tokens[0]
		if token.Type == "Withdraw" {
			event.SentAmount = fmt.Sprintf("%v", token.Amount)
			event.SentCurrency = ConvertCurrency(token.Token)
			event.Label = "airdrop"
		} else if token.Type == "Deposit" {
			event.ReceivedAmount = fmt.Sprintf("%v", token.Amount)
			event.ReceivedCurrency = ConvertCurrency(token.Token)
			event.Label = "income"
		}
		entries = append(entries, event)
		return entries, nil
	}

	if numberOfFTTransfers == 2 {

		event.Label = "Trade"
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

		if len(ftSend) == 1 {
			//we are sending 1 in and getting two back
			amount := ftSend[0].Amount / 2

			for _, ft := range ftReceived {
				ev := event
				ev.Label = "defi"
				ev.ReceivedCurrency = ConvertCurrency(ft.Token)
				ev.ReceivedAmount = fmt.Sprintf("%v", ft.Amount)
				ev.SentCurrency = ConvertCurrency(ftSend[0].Token)
				ev.SentAmount = fmt.Sprintf("%v", amount)
				entries = append(entries, ev)
			}

		} else {
			amount := ftReceived[0].Amount / 2
			for _, ft := range ftSend {
				ev := event
				ev.Label = "defi"
				ev.SentCurrency = ConvertCurrency(ft.Token)
				ev.SentAmount = fmt.Sprintf("%v", ft.Amount)
				ev.ReceivedCurrency = ConvertCurrency(ftReceived[0].Token)
				ev.ReceivedAmount = fmt.Sprintf("%v", amount)
				entries = append(entries, ev)
			}
		}
		return entries, nil
	}

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
		"A.142fa6570b62fd97.StarlyToken":           "ID:773411",
		"A.c6c77b9f5c7a378f.FlowSwapPair":          "NULL1",
	}
	return currencyMap[currency]

}
