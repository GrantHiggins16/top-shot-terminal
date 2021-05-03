package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
)

type ListEvent cadence.Struct

func NewListEvent(fluxEvent cadence.Event, latestBlockHeight uint64, client *client.Client) (*ListEvent, error) {

	id := fluxEvent.Fields[0].(cadence.UInt64)
	walletAddress := cadence.BytesToAddress((fluxEvent.Fields[2].(cadence.Optional)).Value.(cadence.Address).Bytes())

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
	metadata, err := client.ExecuteScriptAtBlockHeight(context.Background(), latestBlockHeight, []byte(getSaleMomentScript), []cadence.Value{
		walletAddress,
		id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata")
	}
	listEvent := ListEvent(metadata.(cadence.Struct))
	return &listEvent, nil
}

func (s ListEvent) ID() uint64 {
	return uint64(s.Fields[0].(cadence.UInt64))
}

func (s ListEvent) PlayID() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}

func (s ListEvent) SetName() string {
	return string(s.Fields[4].(cadence.String))
}

func (s ListEvent) SetID() uint32 {
	return uint32(s.Fields[3].(cadence.UInt32))
}

func (s ListEvent) Play() map[string]string {
	dict := s.Fields[2].(cadence.Dictionary)
	res := map[string]string{}
	for _, kv := range dict.Pairs {
		res[string(kv.Key.(cadence.String))] = string(kv.Value.(cadence.String))
	}
	return res
}

func (s ListEvent) SerialNumber() uint32 {
	return uint32(s.Fields[5].(cadence.UInt32))
}

func (s ListEvent) Price() float64 {
	return float64(s.Fields[6].(cadence.UFix64).ToGoValue().(uint64)) / 1e8
}

func (s ListEvent) PlayerName() string {
	return s.Play()["FullName"]
}

func (s ListEvent) String() string {
	return fmt.Sprintf("Listed Moment:\n\n serialNumber: %d,\n setID: %d,\n setName: %s,\n playID: %d,\n price: %f,\n playerName: %s\n\n",
		s.SerialNumber(), s.SetID(), s.SetName(), s.PlayID(), s.Price(), s.PlayerName())
}

func (s ListEvent) Bytes() []byte {
	b, _ := json.Marshal(s)
	return b
}
