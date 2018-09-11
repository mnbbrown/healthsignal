import React, { Component } from "react";
import Graph from "./Graph";
import api from "./api";
import "./App.css";

class App extends Component {
  state = {
    endpoints: []
  };

  componentDidMount() {
    api.getEndpoints().then(data => this.setState({ endpoints: data }));
  }

  render() {
    return (
      <div className="App">
        <div className="App-container">
          {this.state.endpoints.map((d, i) => (
            <Graph endpoint={d} key={i} />
          ))}
        </div>
      </div>
    );
  }
}

export default App;
