import React, { Component } from "react";
import Graph from "./Graph";
import "./App.css";

class App extends Component {
  state = {
    endpoints: []
  };

  componentDidMount() {
    fetch("http://localhost:8080/endpoints")
      .then(response => response.json())
      .then(data => this.setState({ endpoints: data }));
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
