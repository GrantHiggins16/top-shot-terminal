package service

import (
	"context"
	"fmt"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk"
)

type Event struct {
	id            uint64
	playId        int
	play          string
	setId         int
	setName       int
	serialNumber  int
	price         float32
	uri           string
	lowAsk        float32
	walletAddress flow.Address
}

func NewEvent(id uint64, walletAddress flow.Address, latestBlockHeight uint64, client *client.Client)  *Event {
	event := &Event{
		id:            id,
		walletAddress: walletAddress,
	}
	hydrateMetadata(event, latestBlockHeight, client)
	return event
}

func hydrateMetadata(e *Event, latestBlockHeight uint64, client *client.Client) {
	const getSaleMomentScript = `
	      import TopShot from 0x0b2a3299cc857e29
	      import Market from 0xc1e4f4f4c4257510

	      pub struct SaleMoment {
		pub var id: UInt64
		pub var playId: UInt32
		pub var play: {String: String}
		pub var setId: UInt32
		pub var setName: String
		pub var serialNumber: UInt32
		pub var price: UFix64
		init(moment: &TopShot.NFT, price: UFix64) {
		  self.id = moment.id
		  self.playId = moment.data.playID
		  self.play = TopShot.getPlayMetaData(playID: self.playId)!
		  self.setId = moment.data.setID
		  self.setName = TopShot.getSetName(setID: self.setId)!
		  self.serialNumber = moment.data.serialNumber
		  self.price = price
		}
	      }	
		
	      pub fun main(owner:Address, momentID:UInt64): SaleMoment {
		let acct = getAccount(owner)
		let collectionRef = acct.getCapability(/public/topshotSaleCollection)!.borrow<&{Market.SalePublic}>() ?? panic("Could not borrow capability from public collection")
		return SaleMoment(moment: collectionRef.borrowMoment(id: momentID)!,price: collectionRef.getPrice(tokenID: momentID)!)
	      }`
	metadataRes, err := client.ExecuteScriptAtBlockHeight(context.Background(), latestBlockHeight, []byte(getSaleMomentScript), []cadence.Value{
		cadence.BytesToAddress(e.walletAddress.Bytes()),
		cadence.UInt64(e.id),
	})
	if err != nil {
		fmt.Errorf("Failed to fetch metadata")
	}
	
}
