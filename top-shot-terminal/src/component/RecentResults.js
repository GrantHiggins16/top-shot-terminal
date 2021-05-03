import React, { useState, useEffect } from "react";
import "../config";

export default function RecentResults() {
  const [events, setEvents] = useState({
    eventIds: [],
    eventsDictionary: {}
  });


  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = evt => {
      evt = JSON.parse(evt.data);
      const eventId = evt.Fields[0]
      const newEvent = {
              price: evt.Fields[6],
              serial: evt.Fields[5]
      };
      events.eventsDictionary[eventId] = newEvent;
      events.eventIds = [...events.eventIds, eventId];
      setEvents({
        eventIds: events.eventIds,
        eventsDictionary: events.eventsDictionary
      });
    }

    ws.onopen = () => {
      // on connecting, do nothing but log it to the console
      console.log('connected')
    }
  }, []);
  

  return (
    <div className="App">
      <p>
        Latest processed block: <b>temp</b>
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
          {events.eventIds.map((eventId) => {
            const event = events.eventsDictionary[eventId];
            const momentId = eventId;
            const momentPrice = event.price;
            const momentSerial = event.serial;
            return (
              <tr>
                <td align="left">#{momentId}</td>
                <td align="right">#{momentSerial}</td>
                <td align="right">#{momentPrice}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  )
}

