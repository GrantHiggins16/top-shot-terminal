import React, { Component } from "react";

class RecentResults extends Component {

  constructor() {
    this.state = {moments: new Array()}
  }

  componentDidMount() {
    const columns = React.useMemo(
      () => [
        {
          Header: 'Player',
          accessor: 'player',
        },
        {
          Header: 'Listed Price',
          accessor: 'listedPrice',
        },
      ],
      []
    )
  }

  render() {

  }
}

export default RecentResults;