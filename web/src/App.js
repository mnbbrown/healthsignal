import React, { Component } from "react";
import {
  Charts,
  ChartContainer,
  ChartRow,
  YAxis,
  LineChart
} from "react-timeseries-charts";
import { TimeSeries } from "pondjs";
import "./App.css";

class App extends Component {
  state = {
    data: null
  };

  componentDidMount() {
    this.poller = setInterval(() => {
      fetch("http://localhost:8080/?endpoint=1", {
        method: "GET"
      })
        .then(response => response.json())
        .then(data => {
          const timeseries = new TimeSeries({
            name: "1",
            columns: ["time", "requestTime", "connectionTime"],
            points: data.map(p => [
              p["timestamp"],
              p["requestTime"],
              p["connectionTime"]
            ])
          });
          this.setState({
            data: this.state.data
              ? TimeSeries.timeSeriesListMerge({
                  name: "1",
                  seriesList: [timeseries, this.state.data]
                })
              : timeseries
          });
        })
        .catch(e => console.error(e));
    }, 5000);
  }

  componentWillUnmount() {
    clearInterval(this.poller);
  }

  render() {
    const { data } = this.state;
    return (
      data && (
        <div className="App">
          <ChartContainer
            timeRange={data.timerange()}
            enablePanZoom={true}
            width={800}
          >
            <ChartRow height="200">
              <YAxis
                id="response"
                min={data.min("requestTime")}
                max={data.max("requestTime")}
                width="60"
                type="linear"
              />
              <Charts>
                <LineChart
                  axis="response"
                  series={data}
                  columns={["requestTime"]}
                />
              </Charts>
            </ChartRow>
          </ChartContainer>
        </div>
      )
    );
  }
}

export default App;
