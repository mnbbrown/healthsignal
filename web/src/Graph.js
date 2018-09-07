import React, { Component } from "react";
import {
  styler,
  Charts,
  Resizable,
  ChartContainer,
  ChartRow,
  YAxis,
  LineChart
} from "react-timeseries-charts";
import { TimeSeries } from "pondjs";
import "./Graph.css";

class Graph extends Component {
  state = {
    data: {}
  };

  syncData = () => {
    const { endpoint } = this.props;
    fetch(`https://api.healthsignal.live/endpoints/${endpoint.id}/data`, {
      method: "GET"
    })
      .then(response => response.json())
      .then(data => {
        let locations = data.reduce((lcs, point) => {
          lcs[point["location"]] = lcs[point["location"]] || [];
          lcs[point["location"]].push([
            point["timestamp"],
            point["requestTime"],
            point["connectionTime"],
            point["status"]
          ]);
          return lcs;
        }, {});
        locations = Object.keys(locations).reduce((lcns, lcn) => {
          return {
            ...lcns,
            [lcn]: new TimeSeries({
              name: lcn,
              columns: ["time", "requestTime", "connectionTime", "status"],
              points: locations[lcn]
            })
          };
        }, {});
        const merged = Object.keys(locations).reduce((mlcs, lcn) => {
          if (this.state.data[lcn]) {
            return {
              ...mlcs,
              [lcn]: TimeSeries.timeSeriesListMerge({
                name: lcn,
                seriesList: [locations[lcn], this.state.data[lcn]]
              })
            };
          }
          return { ...mlcs, [lcn]: locations[lcn] };
        }, {});
        this.setState({ data: merged });
      })
      .catch(e => console.error(e));
  };

  componentDidMount() {
    this.syncData();
    this.poller = setInterval(() => {
      this.syncData();
    }, 5000);
  }

  componentWillUnmount() {
    clearInterval(this.poller);
  }

  render() {
    const { data } = this.state;
    const colors = ["steelblue", "green", "pink"];
    const generateLineStyle = index => {
      return styler([{ key: "requestTime", color: colors[index], width: 2 }]);
    };
    if (Object.keys(data).length === 0) {
      return null;
    }
    const locations = Object.keys(data);
    const first = data[locations[0]];
    return (
      data && (
        <div className="Graph">
          <h3>{this.props.endpoint.name}</h3>
          <Resizable>
            <ChartContainer
              timeRange={first.timerange()}
              enablePanZoom={true}
              width={800}
            >
              <ChartRow height="100">
                <YAxis
                  id="response"
                  min={0}
                  max={first.max("requestTime") * 2}
                  width="60"
                  label="ms"
                  type="linear"
                />
                <Charts>
                  {locations.map((location, index) => {
                    return (
                      <LineChart
                        key={location}
                        style={generateLineStyle(index)}
                        axis="response"
                        series={data[location]}
                        columns={["requestTime"]}
                      />
                    );
                  })}
                </Charts>
              </ChartRow>
            </ChartContainer>
          </Resizable>
        </div>
      )
    );
  }
}

export default Graph;
