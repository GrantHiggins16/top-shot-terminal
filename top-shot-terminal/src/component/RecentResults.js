import React, { useState, useEffect } from "react";
import "../config";

export default function RecentResults() {
  const [eventIds, setEventIds] = useState(new Set());
  const [eventsDictionary, setEventsDictionary] = useState({});

  const events = Array.from(eventIds);

  const ws = new WebSocket("ws://localhost:8080/ws");

  ws.onmessage = evt => {
    console.log(evt);
  }

  ws.onopen = () => {
    // on connecting, do nothing but log it to the console
    console.log('connected')
  }

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

