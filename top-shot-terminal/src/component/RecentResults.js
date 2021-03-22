import React, { Component, useMemo, useState, useEffect } from "react";
import "../config";
import * as fcl from "@onflow/fcl";
import * as types from '@onflow/types';

const EVENT_MOMENT_LISTED = "A.c1e4f4f4c4257510.Market.MomentListed";
const FETCH_INTERVAL = 5000;

export default function RecentResults() {
  const [lastBlockHeight, setLastBlockHeight] = useState(0);
  const [eventIds, setEventIds] = useState(new Set());
  const [eventsDictionary, setEventsDictionary] = useState({});

  const fetchEvents = async () => {
    const latestBlock = await fcl.send([fcl.getLatestBlock(true)]);

    const latestHeight = latestBlock.block.height;
    let end = latestHeight;
    let start = lastBlockHeight;
    if (!lastBlockHeight) {
      start = latestHeight;
    }

    // fetch events
    const response = await fcl.send([fcl.getEvents(EVENT_MOMENT_LISTED, start, end)]);

    const { events } = response;

    if (events.length > 0) {
      const newSet = await Promise.all(new Set(
        events.map(async (event) => {
          const id = event.payload.value.fields[0].value.value;
          const accountAddress = event.payload.value.fields[2].value.value.value;
          eventsDictionary[id] = await fetchMetadata({ id: id, accountAddress: accountAddress });
          return id;
        })
      ));
      const newEvents = new Set([...eventIds, ...newSet]);
      setEventsDictionary(eventsDictionary);
      setEventIds(newEvents);
    }

    // update last processed block
    setLastBlockHeight(latestHeight);
  };

  const fetchMetadata = async (event) => {
    if (event.id in eventsDictionary) {
      return eventsDictionary[event.id];
    }
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
      }
    `;
    const metadata = await fcl.send([fcl.script(getSaleMomentScript), 
      fcl.args([fcl.arg(event.accountAddress, types.Address), fcl.arg(parseInt(event.id), types.UInt64)])
    ]);
    const decodedMetadata = await fcl.decode(metadata);
    console.log(decodedMetadata);
    return {
      id: decodedMetadata.id, playId: decodedMetadata.playId, price: decodedMetadata.price, serialNumber: decodedMetadata.serialNumber,
      setId: decodedMetadata.setId, setName: decodedMetadata.setName, seller: event.accountAddress
    };
  }

  useEffect(() => {
    const interval = setInterval(fetchEvents, FETCH_INTERVAL);

    return () => {
      clearInterval(interval);
    };
  });

  const events = Array.from(eventIds);


  return (
    <div className="App">
      <p>
        Latest processed block: <b>#{lastBlockHeight}</b>
      </p>
      <p>
        Events found: <b>{events.length}</b>
      </p>
      <table>
        <thead>
          <tr>
            <th align="left">Moment ID</th>
            <th align="right">Seller</th>
            <th align="right">Price</th>
          </tr>
        </thead>
        <tbody>
          {events.map((eventId) => {
            const event = eventsDictionary[eventId];
            const momentId = eventId;
            const momentPrice = event.price;
            const momentSeller = event.seller;
            return (
              <tr>
                <td align="left">#{momentId}</td>
                <td align="right">#{momentSeller}</td>
                <td align="right">#{momentPrice}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  )
}

